#!/bin/bash

set -o errexit

docker rm $(docker stop proxy)