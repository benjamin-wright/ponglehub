#!/bin/sh

set -o errexit -o pipefail

npm version $(npm view . version)
npm version patch
npm publish

jq '.version' package.json -r > marker