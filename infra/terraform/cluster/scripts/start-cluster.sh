#!/bin/bash

set -o errexit

function start-cluster() {
  echo "Creating k3s cluster..."

  mkdir -p $SCRATCH_DIR/k3s
  cat >$SCRATCH_DIR/k3s/registries.yaml <<EOL
mirrors:
  "$REGISTRY_NAME:$REGISTRY_PORT":
    endpoint:
    - http://$REGISTRY_NAME:$REGISTRY_PORT
EOL

  echo "" > $SCRATCH_DIR/k3s/traefik.yaml

  k3d cluster create $CLUSTER_NAME \
    --wait \
    --image docker.io/rancher/k3s:v1.18.9-k3s1 \
    --volume $SCRATCH_DIR/k3s/registries.yaml:/etc/rancher/k3s/registries.yaml:cached \
    --volume $SCRATCH_DIR/k3s/traefik.yaml:/var/lib/rancher/k3s/server/manifests/traefik.yaml:cached \
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

  echo "waiting for kube services to be ready..."
  local looping=true
  while $looping; do
    local value=$(kubectl get apiservices v1beta1.metrics.k8s.io -o json 2> /dev/null | jq '.status.conditions[] | .status' -r | tr -d '\n')
    if [ "$value" == "True" ]; then
      looping=false
    else
      sleep 0.5
    fi
  done
}

function wait-for-traefik() {
  echo "waiting for traefik ingress controller..."
  local looping=true
  while $looping; do
    local value=$(kubectl get deployment -n kube-system traefik -o json | jq '.status.readyReplicas' -r | tr -d '\n')
    if [ "$value" == "1" ]; then
      looping=false
    else
      sleep 0.5
    fi
  done
}

function overwrite-traefik-config() {
  local old_id=$(kubectl get pods -n kube-system -l app=traefik -o json | jq '.items[0].metadata.name' -r | tr -d '\n')

  cat >$SCRATCH_DIR/k3s/traefik.yaml <<EOL
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

  echo "waiting for replacement traefik ingress controller..."
  local looping=true
  while $looping; do
    local phase=$(kubectl get pods -n kube-system $old_id -o json | jq '.status.phase' -r | tr -d '\n')
    if [ "$phase" == "Running" ]; then
      sleep 0.5
    else
      echo "Phase = $phase"
      looping=false
    fi
  done
}



if k3d cluster list | grep -q pongle; then
  echo "Cluster already exists, skipping..."
  exit 0
fi

start-cluster
wait-for-traefik
overwrite-traefik-config