#!/bin/sh

set -o errexit

(cd infra/terraform/cluster && terraform destroy -auto-approve)
(cd infra/terraform/registries && terraform destroy -auto-approve)

docker stop $(docker ps -q) || true
docker rm $(docker ps -aq) || true
docker system prune --all --volumes --force

rm -rf **/package-lock.json
