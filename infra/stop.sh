#!/bin/bash

set -o errexit -o pipefail

function stop_registry() {
  if k3d registry list -o json | jq '.[].name' -r | grep -q $REGISTRY_NAME; then
    k3d registry delete $REGISTRY_NAME
  else
    echo "Skipping deleting registry, doesn't exist"
  fi
}

function stop_cluster() {
  if k3d cluster list -o json | jq '.[].name' | grep -q $CLUSTER_NAME; then
    k3d cluster delete $CLUSTER_NAME
    rm $KUBECONFIG
  else
    echo "Skipping deleting cluster, doesn't exist"
  fi
}

stop_cluster
stop_registry
