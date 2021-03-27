#!/bin/sh

set -o errexit

CHART=$1

helm repo update

patch_version=$(helm show chart local/$CHART | yq e '.version | split(".") | .[2] | . tag= "!!int" | . + 1' -)

if [[ "$patch_version" == "" ]]; then
  patch_version=0
fi

echo "New version $patch_version"
yq e -i ".version = \"1.0.$patch_version\"" ./Chart.yaml

helm push . local

echo "patch_version" > marker