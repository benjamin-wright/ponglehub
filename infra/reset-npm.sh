#!/bin/bash

set -o errexit

ROOT_DIR=$( cd "$(dirname "`realpath $0`")" ; pwd -P )

helm upgrade --install infra $ROOT_DIR/chart \
    --wait \
    --namespace infra \
    --set-file "grafana.dashboards.default.metrics.json=$ROOT_DIR/dashboards/default.json" \
    --set "secrets.admin.password=password" \
    --set-file global.identityTrustAnchorsPEM=$ROOT_DIR/certs/ca.crt \
    --set-file linkerd2.identity.issuer.tls.crtPEM=$ROOT_DIR/certs/issuer.crt \
    --set-file linkerd2.identity.issuer.tls.keyPEM=$ROOT_DIR/certs/issuer.key \
    --set linkerd2.identity.issuer.crtExpiry=$(date -v+8760H +"%Y-%m-%dT%H:%M:%SZ") \
    --set verdaccio.enabled=false

helm upgrade --install infra $ROOT_DIR/chart \
    --wait \
    --namespace infra \
    --set-file "grafana.dashboards.default.metrics.json=$ROOT_DIR/dashboards/default.json" \
    --set "secrets.admin.password=password" \
    --set-file global.identityTrustAnchorsPEM=$ROOT_DIR/certs/ca.crt \
    --set-file linkerd2.identity.issuer.tls.crtPEM=$ROOT_DIR/certs/issuer.crt \
    --set-file linkerd2.identity.issuer.tls.keyPEM=$ROOT_DIR/certs/issuer.key \
    --set linkerd2.identity.issuer.crtExpiry=$(date -v+8760H +"%Y-%m-%dT%H:%M:%SZ")

/usr/bin/expect <<EOD
spawn npm login --registry "$NPM_REGISTRY" --scope=pongle --strict-ssl false
expect {
  "Username:" {send "$NPM_USERNAME\r"; exp_continue}
  "Password:" {send "$NPM_PASSWORD\r"; exp_continue}
  "Email: (this IS public)" {send "$NPM_EMAIL\r"; exp_continue}
}
EOD