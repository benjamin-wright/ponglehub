#!/bin/bash

set -o errexit

docker rm $(docker stop npm) || true
docker rm $(docker stop proxy) || true

docker network rm local-registries || true

rm -rf ./infra/local-repos/ssl