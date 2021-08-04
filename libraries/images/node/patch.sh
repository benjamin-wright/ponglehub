#!/bin/sh

set -o errexit -o pipefail

echo "getting npm version"
npm version $(npm view . version) --allow-same-version

echo "patching version"
npm version patch

echo "publishing"
npm publish