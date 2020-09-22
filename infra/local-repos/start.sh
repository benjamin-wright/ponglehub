#!/bin/bash

set -o errexit

function npm-login() {
  /usr/bin/expect <<EOD
spawn npm login --registry "$NPM_REGISTRY" --scope=pongle --strict-ssl false
expect {
  "Username:" {send "$NPM_USERNAME\r"; exp_continue}
  "Password:" {send "$NPM_PASSWORD\r"; exp_continue}
  "Email: (this IS public)" {send "$NPM_EMAIL\r"; exp_continue}
}
EOD
}

if [ ! -f $PWD/infra/local-repos/ssl/certificate.crt ]; then
  docker run --rm -v $PWD/infra/local-repos/ssl:/work -it nginx \
    openssl req \
    -out /work/CSR.csr \
    -new \
    -newkey rsa:2048 \
    -nodes \
    -keyout /work/caKey.key \
    -subj "/C=UK/ST=Test/L=Test/O=Test/CN=ponglehub.co.uk"

  docker run --rm -v $PWD/infra/local-repos/ssl:/work -it nginx \
    openssl req \
    -x509 \
    -sha256 \
    -nodes \
    -days 365 \
    -newkey rsa:2048 \
    -keyout /work/private.key \
    -out /work/certificate.crt \
    -subj "/C=UK/ST=Test/L=Test/O=Test/CN=ponglehub.co.uk" \
    -addext "subjectAltName = DNS:ponglehub.co.uk"
fi

NETWORK_NAME=local-registries
PROXY_NAME=proxy
NPM_NAME=npm

docker network create $NETWORK_NAME

docker run \
  -d \
  --name $NPM_NAME \
  --network $NETWORK_NAME \
  docker.io/verdaccio/verdaccio

docker run \
  -d \
  --name $PROXY_NAME \
  --network $NETWORK_NAME \
  -p 80:80 \
  -p 443:443 \
  -v $PWD/infra/local-repos/default.conf:/etc/nginx/nginx.conf \
  -v $PWD/infra/local-repos/ssl/certificate.crt:/etc/nginx/ssl/certificate.crt \
  -v $PWD/infra/local-repos/ssl/private.key:/etc/nginx/ssl/private.key \
  docker.io/nginx

sleep 5

npm-login