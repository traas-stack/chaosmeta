#!/bin/bash

set -e
VERSION=$1
ROOT_PATH=$(dirname $(readlink -f $0))/chaosmeta_build
echo "${ROOT_PATH}"
mkdir -p "${ROOT_PATH}"/yamls
curl -o "${ROOT_PATH}"/yamls/chaosmeta.yaml https://raw.githubusercontent.com/traas-stack/chaosmeta/"${VERSION}"/chaosmeta-inject-operator/build/yamls/chaosmeta.yaml
kubectl delete --ignore-not-found=$(ignore-not-found) -f "${ROOT_PATH}"/yamls/chaosmeta.yaml
