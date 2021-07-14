#!/bin/sh

set -o errexit -o pipefail

/usr/bin/expect <<EOD 
spawn npm login --registry http://docker.for.mac.localhost:4873 --scope=pongle
expect {
  "Username:" {send "bot\r"; exp_continue}
  "Password:" {send "password\r"; exp_continue}
  "Email: (this IS public)" {send "bot@example.com\r"; exp_continue}
}
EOD

if [ ! -f ~/.npmrc ]; then
  echo "didn't create npmrc file"
  exit 1
fi