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

    k3d cluster create $CLUSTER_NAME \
        --volume $ROOT_DIR/k3s/registries.yaml:/etc/rancher/k3s/registries.yaml:cached \
        --volume $ROOT_DIR/k3s/traefik.yaml:/var/lib/rancher/k3s/server/manifests/traefik.yaml:cached \
        --agents 2 \
        --network "$NETWORK_NAME" \
        -p 80:80@loadbalancer \
        -p 443:443@loadbalancer

    local nodes=$(kubectl get nodes -o go-template --template='{{range .items}}{{printf "%s\n" .metadata.name}}{{end}}')
    for node in $nodes; do
        kubectl annotate node "${node}" tilt.dev/registry=localhost:${REGISTRY_PORT};
        docker exec "$node" sysctl fs.inotify.max_user_watches=524288
        docker exec "$node" sysctl fs.inotify.max_user_instances=512
    done
}

function install-linkerd() {
    if kubectl get ns | grep -q linkerd; then
        echo "Linkerd already installed, skipping..."
        return
    fi

	linkerd install | kubectl apply -f -
}

function deploy-infra() {
    echo "Deploying/upgrading standard infrastructure..."
    kubectl get ns | grep infra || kubectl create ns infra
    kubectl annotate namespace infra linkerd.io/inject=enabled --overwrite
    helm dep update $ROOT_DIR/chart
    helm upgrade --install infra $ROOT_DIR/chart \
        --namespace infra \
        --set-file "grafana.dashboards.default.metrics.json=$ROOT_DIR/dashboards/default.json" \
        --set "secrets.admin.password=$(echo "password" | tr -d '\n' | base64)"
}

create-network
start-registry
start-cluster
install-linkerd
deploy-infra