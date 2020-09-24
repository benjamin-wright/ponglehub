#!/bin/bash

set -o errexit

ROOT_DIR=$( cd "$(dirname "`realpath $0`")" ; pwd -P )

if k3d cluster list | grep $CLUSTER_NAME -q; then
    sudo security remove-trusted-cert -d infra/cluster/ssl/ponglehubCA.crt || true
    npm config delete -g cafile

    k3d cluster delete $CLUSTER_NAME
    docker rm $(docker stop $REGISTRY_NAME)
    docker network rm $NETWORK_NAME
    rm -rf $ROOT_DIR/k3s
    rm -rf $ROOT_DIR/ssl
fi