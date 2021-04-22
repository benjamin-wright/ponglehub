#!/bin/bash

set -o errexit -o pipefail

function start_npm_registry() {
  if docker ps --format '{{ .Names }}' | grep -q $NPM_CONTAINER; then
    echo "Skipping installing npm registry, already installed"
  else
    docker run -d --restart always -p 4873:4873 --name $NPM_CONTAINER verdaccio/verdaccio:4
  
    if [ ! -f ~/.npmrc ]; then
      touch ~/.npmrc
    fi

    if [ ! -f ~/.npmrc.bak ]; then
      cp ~/.npmrc ~/.npmrc.bak

      success="1"

      while [[ "$success" != "0" ]]; do
        npm ping --registry http://localhost:4873
        success="$?"
      done

      /usr/bin/expect <<EOD
spawn npm login --registry http://localhost:4873 --scope=pongle
expect {
  "Username:" {send "local\r"; exp_continue}
  "Password:" {send "password\r"; exp_continue}
  "Email: (this IS public)" {send "local@example.com\r"; exp_continue}
}
EOD

      npm config set registry http://localhost:4873
    fi
  fi
}

start_npm_registry