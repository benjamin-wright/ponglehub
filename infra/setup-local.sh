#!/bin/bash

set -o errexit

/usr/bin/expect <<EOD
spawn npm login --registry "$NPM_REGISTRY" --scope=pongle
expect {
  "Username:" {send "$NPM_USERNAME\r"; exp_continue}
  "Password:" {send "$NPM_PASSWORD\r"; exp_continue}
  "Email: (this IS public)" {send "$NPM_EMAIL\r"; exp_continue}
}
EOD

npm config set registry $NPM_REGISTRY