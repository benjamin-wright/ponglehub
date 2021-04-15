#!/bin/bash

set -o errexit -o pipefail

function start_registry() {
  if k3d registry list -o json | jq '.[].name' -r | grep -q $REGISTRY_NAME; then
    echo "Skipping creating registry, already exists"
  else
    k3d registry create $REGISTRY_NAME --port $REGISTRY_PORT
  fi
}

function start_cluster() {
  if k3d cluster list -o json | jq '.[].name' | grep -q $CLUSTER_NAME; then
    echo "Skipping creating cluster, already exists"
  else
    k3d cluster create $CLUSTER_NAME \
      --registry-use $REGISTRY_NAME \
      --agents 3 \
      --servers 1 \
      --k3s-server-arg "--no-deploy=traefik" \
      --kubeconfig-update-default=false \
      --wait
    
    mkdir -p $SCRATCH_DIR
    k3d kubeconfig get $CLUSTER_NAME > $KUBECONFIG
  fi
}

function start_istio() {
  if kubectl get ns -o json | jq '.items[].metadata.name' -r | grep -q 'istio-system'; then
    echo "Skipping installing istio, already installed"
  else
    istioctl install --set profile=demo -y
  fi
}

function start_npm_registry() {
  docker run -d --restart always -p 4873:4873 verdaccio/verdaccio:4
  
  if [ ! -f ~/.npmrc ]; then
    touch ~/.npmrc
  fi

  if [ ! -f ~/.npmrc.bak ]; then
    cp ~/.npmrc ~/.npmrc.bak

    success="1"

    while [[ "$success" != "0" ]]; do
      npm ping --registry http://localhost:4873
      success="$?"
    done

    /usr/bin/expect <<EOD
spawn npm login --registry http://localhost:4873 --scope=pongle
expect {
  "Username:" {send "$NPM_USERNAME\r"; exp_continue}
  "Password:" {send "$NPM_PASSWORD\r"; exp_continue}
  "Email: (this IS public)" {send "$NPM_EMAIL\r"; exp_continue}
}
EOD

    npm config set registry $NPM_REGISTRY
  fi
}

start_registry
start_cluster
start_istio
start_npm_registry