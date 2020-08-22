#!/bin/bash

set -o errexit

ROOT_DIR=$( cd "$(dirname "`realpath $0`")" ; pwd -P )

k3d cluster delete $CLUSTER_NAME
docker rm $(docker stop $REGISTRY_NAME)
docker network rm $NETWORK_NAME
rm -rf $ROOT_DIR/k3s
rm -rf $ROOT_DIR/certs