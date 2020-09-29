#!/bin/bash

set -o errexit

ROOT_DIR=$( cd "$(dirname "`realpath $0`")" ; pwd -P )

function create-network() {
    if docker network ls | grep $NETWORK_NAME -q; then
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

    local nodes=$(kubectl get nodes -o go-template --template='{{range .items}}{{printf "%s\n" .metadata.name}}{{end}}')
    for node in $nodes; do
        kubectl annotate node "${node}" tilt.dev/registry=localhost:${REGISTRY_PORT};
        kubectl annotate node "${node}" tilt.dev/registry-from-cluster=${REGISTRY_NAME}:${REGISTRY_PORT}
        docker exec "$node" sysctl fs.inotify.max_user_watches=524288
        docker exec "$node" sysctl fs.inotify.max_user_instances=512
    done

    echo "waiting for kube services to be ready"
    looping=true
    while $looping; do
      value=$(kubectl get apiservices v1beta1.metrics.k8s.io -o json | jq '.status.conditions[] | .status' -r | tr -d '\n')
      if [ "$value" == "True" ]; then
        looping=false
      else
        sleep 0.5
      fi
    done
}

function make-linkerd-certs() {
  if [ -f $ROOT_DIR/ssl/linkerd.crt ]; then
    echo "Certs already exist, skipping..."
    return
  fi

  mkdir -p $ROOT_DIR/ssl

  step certificate create identity.linkerd.cluster.local $ROOT_DIR/ssl/linkerdCA.crt $ROOT_DIR/ssl/linkerdCA.key \
    --profile root-ca \
    --no-password \
    --insecure

  step certificate create identity.linkerd.cluster.local $ROOT_DIR/ssl/linkerd.crt $ROOT_DIR/ssl/linkerd.key \
    --ca $ROOT_DIR/ssl/linkerdCA.crt \
    --ca-key $ROOT_DIR/ssl/linkerdCA.key \
    --profile intermediate-ca \
    --not-after 8760h \
    --no-password \
    --insecure
}

function make-ingress-certs() {
  local SSL_PATH=$PWD/infra/cluster/ssl
  local CA_NAME=ponglehubCA
  local DOMAIN=ponglehub.co.uk
  if [ ! -f $SSL_PATH/$DOMAIN.crt ]; then
    docker run --rm -v $SSL_PATH:/work -it nginx \
      openssl genrsa -out /work/$CA_NAME.key 2048

    docker run --rm -v $SSL_PATH:/work -it nginx \
      openssl req \
      -x509 \
      -new \
      -nodes \
      -key /work/$CA_NAME.key \
      -sha256 \
      -days 1825 \
      -out /work/$CA_NAME.crt \
      -subj "/C=UK/ST=Test/L=Test/O=Test/CN=$DOMAIN"

    docker run --rm -v $SSL_PATH:/work -it nginx \
      openssl genrsa -out /work/$DOMAIN.key 2048

    docker run --rm -v $SSL_PATH:/work -it nginx \
      openssl req \
      -new \
      -key /work/$DOMAIN.key \
      -out /work/$DOMAIN.csr \
      -subj "/C=UK/ST=Test/L=Test/O=Test/CN=$DOMAIN"

    cat > $SSL_PATH/$DOMAIN.ext << EOF
authorityKeyIdentifier=keyid,issuer
basicConstraints=CA:FALSE
keyUsage = digitalSignature, nonRepudiation, keyEncipherment, dataEncipherment
subjectAltName = @alt_names
[alt_names]
DNS.1 = *.$DOMAIN
EOF

    docker run --rm -v $SSL_PATH:/work -it nginx \
      openssl x509 \
      -req \
      -in /work/$DOMAIN.csr \
      -CA /work/$CA_NAME.crt \
      -CAkey /work/$CA_NAME.key \
      -CAcreateserial \
      -out /work/$DOMAIN.crt -days 825 -sha256 -extfile /work/$DOMAIN.ext

    sudo security add-trusted-cert -d -r trustRoot -k /Library/Keychains/System.keychain $SSL_PATH/$CA_NAME.crt
    npm config set -g cafile $SSL_PATH/$CA_NAME.crt
  fi
}

function load-dashboard() {
  local dashboard=$1

  if [ -f $ROOT_DIR/dashboards/$dashboard ]; then
    return
  fi

  curl -L https://raw.githubusercontent.com/linkerd/linkerd2/main/grafana/dashboards/$dashboard > $ROOT_DIR/dashboards/$dashboard
}

function load-from-url() {
  local name=$1
  local url=$2

  if [ -f $ROOT_DIR/dashboards/$name ]; then
    return
  fi

  curl -L $url > $ROOT_DIR/dashboards/$name

  sed -i '' 's/${DS_PROMETHEUS}/prometheus/g' $ROOT_DIR/dashboards/$name
}

function deploy-infra() {
  echo "downloading dashboards"
  mkdir -p $ROOT_DIR/dashboards
  load-dashboard top-line.json
  load-dashboard namespace.json
  load-dashboard pod.json
  load-dashboard health.json
  load-from-url kubernetes.json https://grafana.com/api/dashboards/8588/revisions/1/download
  load-from-url kubernetes-cluster.json https://grafana.com/api/dashboards/11802/revisions/4/download

  echo "installing linkerd..."
  helm upgrade -i linkerd linkerd2 \
    --wait \
    --repo https://helm.linkerd.io/edge \
    --set grafana.enabled=false \
    --set ingress.hostname=linkerd.ponglehub.co.uk \
    --set enforcedHostRegexp="^linkerd\\.ponglehub\\.co\\.uk\$" \
    --set identity.issuer.crtExpiry=$(date -v+8760H +"%Y-%m-%dT%H:%M:%SZ") \
    --set-file global.identityTrustAnchorsPEM=$ROOT_DIR/ssl/linkerdCA.crt \
    --set-file identity.issuer.tls.crtPEM=$ROOT_DIR/ssl/linkerd.crt \
    --set-file identity.issuer.tls.keyPEM=$ROOT_DIR/ssl/linkerd.key \

  echo "Deploying/upgrading standard infrastructure..."
  kubectl get ns | grep infra || kubectl create ns infra
  kubectl annotate namespace infra linkerd.io/inject=enabled --overwrite
  helm dep update $ROOT_DIR/chart
  helm upgrade --install infra $ROOT_DIR/chart \
    --wait \
    --timeout 10m0s \
    --namespace infra \
    --set-file "grafana.dashboards.default.top-line.json=$ROOT_DIR/dashboards/top-line.json" \
    --set-file "grafana.dashboards.default.namespace.json=$ROOT_DIR/dashboards/namespace.json" \
    --set-file "grafana.dashboards.default.pod.json=$ROOT_DIR/dashboards/pod.json" \
    --set-file "grafana.dashboards.default.health.json=$ROOT_DIR/dashboards/health.json" \
    --set-file "grafana.dashboards.default.kubernetes.json=$ROOT_DIR/dashboards/kubernetes.json" \
    --set-file "grafana.dashboards.default.kubernetes-cluster.json=$ROOT_DIR/dashboards/kubernetes-cluster.json" \
    --set "secrets.admin.password=password" \
    --set-file secrets.ssl.key=$ROOT_DIR/ssl/ponglehub.co.uk.key \
    --set-file secrets.ssl.crt=$ROOT_DIR/ssl/ponglehub.co.uk.crt
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

function npm-login() {
  /usr/bin/expect <<EOD
spawn npm login --registry "$NPM_REGISTRY" --scope=pongle
expect {
  "Username:" {send "$NPM_USERNAME\r"; exp_continue}
  "Password:" {send "$NPM_PASSWORD\r"; exp_continue}
  "Email: (this IS public)" {send "$NPM_EMAIL\r"; exp_continue}
}
EOD
}

function helm-login() {
  helm repo add local https://helm.ponglehub.co.uk
}

make-linkerd-certs
make-ingress-certs

create-network
start-registry
start-cluster
deploy-infra

overwrite-traefik-config

echo "waiting for things to start..."
sleep 5

npm-login
helm-login