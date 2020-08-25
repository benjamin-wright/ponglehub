#!/bin/bash

set -o errexit

ROOT_DIR=$( cd "$(dirname "`realpath $0`")" ; pwd -P )

function create-network() {
    if docker network ls | grep $NETWORK_NAME; then
        echo "Network already exists, skipping..."
        return
    fi

    echo "Creating network..."
    docker network create $NETWORK_NAME
}

function start-registry() {
    if [[ "$(docker ps --filter name=$REGISTRY_NAME -q)" != "" ]]; then
        echo "Registry already exists, skipping..."
        return
    fi

    docker run \
        -d --restart=always -p "$REGISTRY_PORT:5000" --name "$REGISTRY_NAME" --network "$NETWORK_NAME" \
        registry:2
}

function start-cluster() {
    if k3d cluster list | grep -q pongle; then
        echo "Cluster already exists, skipping..."
        return
    fi

    echo "Creating k3s cluster..."

    mkdir -p $ROOT_DIR/k3s
    cat >$ROOT_DIR/k3s/registries.yaml <<EOL
mirrors:
  "localhost:$REGISTRY_PORT":
    endpoint:
    - http://$REGISTRY_NAME:$REGISTRY_PORT
EOL

    echo "" > $ROOT_DIR/k3s/traefik.yaml

    k3d cluster create $CLUSTER_NAME \
        --volume $ROOT_DIR/k3s/registries.yaml:/etc/rancher/k3s/registries.yaml:cached \
        --volume $ROOT_DIR/k3s/traefik.yaml:/var/lib/rancher/k3s/server/manifests/traefik.yaml:cached \
        --agents $WORKER_NODES \
        --network "$NETWORK_NAME" \
        -p 80:80@loadbalancer \
        -p 443:443@loadbalancer

    sleep 5

    local nodes=$(kubectl get nodes -o go-template --template='{{range .items}}{{printf "%s\n" .metadata.name}}{{end}}')
    for node in $nodes; do
        kubectl annotate node "${node}" tilt.dev/registry=localhost:${REGISTRY_PORT};
        docker exec "$node" sysctl fs.inotify.max_user_watches=524288
        docker exec "$node" sysctl fs.inotify.max_user_instances=512
    done
}

function make-certs() {
  if [ -d $ROOT_DIR/certs ]; then
    echo "Certs already exist, skipping..."
    return
  fi

  mkdir -p $ROOT_DIR/certs

  step certificate create identity.linkerd.cluster.local $ROOT_DIR/certs/ca.crt $ROOT_DIR/certs/ca.key \
    --profile root-ca \
    --no-password \
    --insecure

  step certificate create identity.linkerd.cluster.local $ROOT_DIR/certs/issuer.crt $ROOT_DIR/certs/issuer.key \
    --ca $ROOT_DIR/certs/ca.crt \
    --ca-key $ROOT_DIR/certs/ca.key \
    --profile intermediate-ca \
    --not-after 8760h \
    --no-password \
    --insecure
}

function deploy-infra() {
  echo "Deploying/upgrading standard infrastructure..."
  kubectl get ns | grep infra || kubectl create ns infra
  kubectl annotate namespace infra linkerd.io/inject=enabled --overwrite
  helm dep update $ROOT_DIR/chart
  helm upgrade --install infra $ROOT_DIR/chart \
    --wait \
    --namespace infra \
    --set-file "grafana.dashboards.default.metrics.json=$ROOT_DIR/dashboards/default.json" \
    --set "secrets.admin.password=password" \
    --set-file global.identityTrustAnchorsPEM=$ROOT_DIR/certs/ca.crt \
    --set-file linkerd2.identity.issuer.tls.crtPEM=$ROOT_DIR/certs/issuer.crt \
    --set-file linkerd2.identity.issuer.tls.keyPEM=$ROOT_DIR/certs/issuer.key \
    --set linkerd2.identity.issuer.crtExpiry=$(date -v+8760H +"%Y-%m-%dT%H:%M:%SZ")
}

function npm-login() {
  /usr/bin/expect <<EOD
spawn npm login --registry "$NPM_REGISTRY" --scope=pongle --strict-ssl false
expect {
  "Username:" {send "$NPM_USERNAME\r"; exp_continue}
  "Password:" {send "$NPM_PASSWORD\r"; exp_continue}
  "Email: (this IS public)" {send "$NPM_EMAIL\r"; exp_continue}
}
EOD
}

function overwrite-traefik-config() {
    cat >$ROOT_DIR/k3s/traefik.yaml <<EOL
apiVersion: helm.cattle.io/v1
kind: HelmChart
metadata:
  name: traefik
  namespace: kube-system
spec:
  chart: https://%{KUBERNETES_API}%/static/charts/traefik-1.81.0.tgz
  valuesContent: |-
    rbac:
      enabled: true
    ssl:
      enabled: true
    dashboard:
      enabled: true
      domain: traefik.ponglehub.co.uk
      ingress:
        tls:
        - hosts:
          - traefik.ponglehub.co.uk
    metrics:
      prometheus:
        enabled: true
    kubernetes:
      ingressEndpoint:
        useDefaultPublishedService: true
    image: "rancher/library-traefik"
    tolerations:
      - key: "CriticalAddonsOnly"
        operator: "Exists"
      - key: "node-role.kubernetes.io/master"
        operator: "Exists"
        effect: "NoSchedule"
EOL
}

create-network
start-registry
start-cluster
make-certs
deploy-infra
npm-login
overwrite-traefik-config