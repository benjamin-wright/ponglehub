#!/bin/sh

set -o errexit

rm -rf **/package-lock.json

(cd infra/terraform/registries && terraform apply -auto-approve)
(cd infra/terraform/cluster && terraform apply -auto-approve)

earthly +init
helm repo update
