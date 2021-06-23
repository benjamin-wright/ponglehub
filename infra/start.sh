#!/bin/bash

set -o errexit -o pipefail

# Start the docker registry that will be visible across localhost and the kube cluster
function start_registry() {
  if k3d registry list -o json | jq '.[].name' -r | grep -q $REGISTRY_NAME; then
    echo "Skipping creating registry, already exists"
  else
    k3d registry create $REGISTRY_NAME --port $REGISTRY_PORT
  fi
}

# Start the k3d cluster
function start_cluster() {
  if k3d cluster list -o json | jq '.[].name' | grep -q $CLUSTER_NAME; then
    echo "Skipping creating cluster, already exists"
  else
    # Disable the automatic kubeconfig update so that we can write it to a specific file, to avoid
    # contaminating the global kubeconfig with our development cluster
    k3d cluster create $CLUSTER_NAME \
      --registry-use $REGISTRY_NAME \
      --k3s-server-arg "--no-deploy=traefik" \
      --kubeconfig-update-default=false \
      --volume $(pwd)/infra/manifests:/var/lib/rancher/k3s/server/manifests/preload \
      -p "80:80@loadbalancer" \
      --wait
    
    mkdir -p $SCRATCH_DIR
    k3d kubeconfig get $CLUSTER_NAME > $KUBECONFIG
  fi
}

# Start the knative components
function start_knative() {
  kubectl apply --wait -f https://github.com/knative/serving/releases/download/v0.22.0/serving-crds.yaml
  kubectl apply --wait -f https://github.com/knative/serving/releases/download/v0.22.0/serving-core.yaml
  istioctl install --set profile=demo -y
  kubectl apply --wait -f https://github.com/knative/net-istio/releases/download/v0.22.0/net-istio.yaml
}

# Knative serving has some issues with the private docker registry, so patch the coredns config to
# redirect calls to the private registry host to the docker host IP (i.e. localhost on host system)
function update_coredns() {
  local file_name=tmp_configmap.yaml
  local backup_file_name=tmp_configmap.yaml.bak

  # grab the core-dns configuration from the cluster
  kubectl get configmap -n kube-system coredns -o yaml > $file_name

  if cat $file_name | grep k3d-$REGISTRY_NAME -q; then
    echo "hosts entry for private registry already exists."
  else
    # Find the IP associated with host.k3d.internal, this will be the IP of your host machine on the docker network
    local registry_ip=$(cat $file_name | grep host.k3d.internal | xargs | cut -d " " -f1)
    local line_number=$(cat tmp_configmap.yaml | grep host.k3d.internal -n | cut -f1 -d: | tr -d '\n')

    # Add the docker registry hostname with the host machine IP address to the core-dns config file
    sed -i.bak "${line_number}i\\
    $registry_ip k3d-$REGISTRY_NAME
" $file_name

    kubectl replace -n kube-system -f $file_name --wait

    # roll the core-dns pods to reload the config
    kubectl -n kube-system rollout restart deployment coredns
  fi

  rm $file_name
  rm $backup_file_name
}

start_registry
start_cluster
# start_knative
# update_coredns