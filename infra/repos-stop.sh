#!/bin/bash

set -o errexit -o pipefail

function stop_npm_registry() {
  if docker ps --format '{{ .Names }}' | grep -q $NPM_CONTAINER; then
    docker stop $NPM_CONTAINER
    docker rm $NPM_CONTAINER

    if [ -f ~/.npmrc.bak ]; then
      mv ~/.npmrc.bak ~/.npmrc
    fi

    npm config delete registry
  else
    echo "npm registry already down"
  fi
}

stop_npm_registry