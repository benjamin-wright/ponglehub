#!/bin/bash

set -o errexit

function npm-login() {
  cp ~/.npmrc ~/.npmrc.bak

  /usr/bin/expect <<EOD
spawn npm login --registry "$NPM_REGISTRY" --scope=pongle
expect {
  "Username:" {send "$NPM_USERNAME\r"; exp_continue}
  "Password:" {send "$NPM_PASSWORD\r"; exp_continue}
  "Email: (this IS public)" {send "$NPM_EMAIL\r"; exp_continue}
}
EOD

  npm config set registry $NPM_REGISTRY
}

function helm-login() {
  helm repo add local https://helm.ponglehub.co.uk --insecure-skip-tls-verify
}

SSL_PATH=$PWD/infra/repos/ssl
CA_NAME=ponglehubCA
DOMAIN=ponglehub.co.uk
if [ ! -f $SSL_PATH/$DOMAIN.crt ]; then
  docker run --rm -v $SSL_PATH:/work -it nginx \
    openssl genrsa -out /work/$CA_NAME.key 2048

  docker run --rm -v $SSL_PATH:/work -it nginx \
    openssl req \
    -x509 \
    -new \
    -nodes \
    -key /work/$CA_NAME.key \
    -sha256 \
    -days 1825 \
    -out /work/$CA_NAME.crt \
    -subj "/C=UK/ST=Test/L=Test/O=Test/CN=$DOMAIN"

  docker run --rm -v $SSL_PATH:/work -it nginx \
    openssl genrsa -out /work/$DOMAIN.key 2048

  docker run --rm -v $SSL_PATH:/work -it nginx \
    openssl req \
    -new \
    -key /work/$DOMAIN.key \
    -out /work/$DOMAIN.csr \
    -subj "/C=UK/ST=Test/L=Test/O=Test/CN=$DOMAIN"

  cat > $SSL_PATH/$DOMAIN.ext << EOF
authorityKeyIdentifier=keyid,issuer
basicConstraints=CA:FALSE
keyUsage = digitalSignature, nonRepudiation, keyEncipherment, dataEncipherment
subjectAltName = @alt_names
[alt_names]
DNS.1 = *.$DOMAIN
EOF

  docker run --rm -v $SSL_PATH:/work -it nginx \
    openssl x509 \
    -req \
    -in /work/$DOMAIN.csr \
    -CA /work/$CA_NAME.crt \
    -CAkey /work/$CA_NAME.key \
    -CAcreateserial \
    -out /work/$DOMAIN.crt -days 825 -sha256 -extfile /work/$DOMAIN.ext

  sudo security add-trusted-cert -d -r trustRoot -k /Library/Keychains/System.keychain infra/repos/ssl/$CA_NAME.crt
  npm config set -g cafile $SSL_PATH/$CA_NAME.crt
fi

NETWORK_NAME=local-registries
PROXY_NAME=proxy
NPM_NAME=npm
CHART_MUSEUM_NAME=helm

docker network create $NETWORK_NAME

docker run \
  -d \
  --name $NPM_NAME \
  --network $NETWORK_NAME \
  docker.io/verdaccio/verdaccio

docker run \
  -d \
  --name $CHART_MUSEUM_NAME \
  --network $NETWORK_NAME \
  -v $(pwd)/chart-cache:/charts \
  -e STORAGE=local \
  -e STORAGE_LOCAL_ROOTDIR=/charts \
  chartmuseum/chartmuseum:latest

docker run \
  -d \
  --name $PROXY_NAME \
  --network $NETWORK_NAME \
  -p 80:80 \
  -p 443:443 \
  -v $PWD/infra/repos/default.conf:/etc/nginx/nginx.conf \
  -v $PWD/infra/repos/ssl/ponglehub.co.uk.crt:/etc/nginx/ssl/certificate.crt \
  -v $PWD/infra/repos/ssl/ponglehub.co.uk.key:/etc/nginx/ssl/private.key \
  docker.io/nginx

echo "waiting for things to start..."
sleep 5

npm-login
helm-login