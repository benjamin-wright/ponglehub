#!/bin/bash

set -o errexit

helm repo remove local || true

docker rm $(docker stop helm) || true
docker rm $(docker stop npm) || true
docker rm $(docker stop proxy) || true

docker network rm local-registries || true

rm -rf ./infra/local-repos/ssl
rm -rf ./chart-cache