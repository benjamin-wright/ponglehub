#!/bin/bash

set -o errexit

if [ ! -f $PWD/infra/local-repos/ssl/certificate.crt ]; then
  docker run --rm -v $PWD/infra/local-repos/ssl:/work -it nginx \
    openssl req \
    -out /work/CSR.csr \
    -new \
    -newkey rsa:2048 \
    -nodes \
    -keyout /work/privateKey.key \
    -subj "/C=UK/ST=Test/L=Test/O=Test/CN=ponglehub.co.uk"

  docker run --rm -v $PWD/infra/local-repos/ssl:/work -it nginx \
    openssl req \
    -x509 \
    -sha256 \
    -nodes \
    -days 365 \
    -newkey rsa:2048 \
    -keyout privateKey.key \
    -out /work/certificate.crt \
    -subj "/C=UK/ST=Test/L=Test/O=Test/CN=ponglehub.co.uk" \
    -addext "subjectAltName = DNS:ponglehub.co.uk"
fi

docker run -d --name proxy \
  -v $PWD/infra/local-repos/default.conf:/etc/nginx/nginx.conf \
  -v $PWD/infra/local-repos/ssl/certificate.crt:/etc/nginx/ssl/ponglehub.crt \
  -v $PWD/infra/local-repos/ssl/private.key:/etc/nginx/ssl/private.key \
  nginx