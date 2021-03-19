#!/bin/bash

set -o errexit

if [ -f ~/.npmrc.bak ]; then
  mv ~/.npmrc.bak ~/.npmrc
fi

