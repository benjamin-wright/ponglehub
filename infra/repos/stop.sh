#!/bin/bash

set -o errexit

if docker network ls | grep local-registries -q; then
    helm repo remove local || true

    docker rm $(docker stop helm) || true
    docker rm $(docker stop npm) || true
    docker rm $(docker stop proxy) || true

    docker network rm local-registries || true

    sudo security remove-trusted-cert -d infra/repos/ssl/ponglehubCA.crt || true
    npm config delete -g cafile
    rm -rf ./infra/repos/ssl
    rm -rf ./chart-cache
fi