#!/bin/bash

set -o errexit

if k3d cluster list | grep $CLUSTER_NAME -q; then
    k3d cluster delete $CLUSTER_NAME
    rm -rf $SCRATCH_DIR
fi