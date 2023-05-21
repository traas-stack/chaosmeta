#!/bin/bash

set -e
ROOT_PATH=$(dirname $(readlink -f $0))/chaosmeta_build
echo "${ROOT_PATH}"
mkdir -p "${ROOT_PATH}"/build
mkdir -p "${ROOT_PATH}"/config
mkdir -p "${ROOT_PATH}"/yamls

curl -o "${ROOT_PATH}"/config/chaosmeta-inject.json https://raw.githubusercontent.com/traas-stack/chaosmeta/master/chaosmeta-inject-operator/config/chaosmeta-inject.json
curl -o "${ROOT_PATH}"/build/build.sh https://raw.githubusercontent.com/traas-stack/chaosmeta/master/chaosmeta-inject-operator/build/build.sh
curl -o "${ROOT_PATH}"/yamls/chaosmeta.yaml https://raw.githubusercontent.com/traas-stack/chaosmeta/master/chaosmeta-inject-operator/build/yamls/chaosmeta.yaml
curl -o "${ROOT_PATH}"/yamls/chaosmeta-daemonset.yaml https://raw.githubusercontent.com/traas-stack/chaosmeta/master/chaosmeta-inject-operator/build/yamls/chaosmeta-daemonset.yaml

kubectl apply -f "${ROOT_PATH}"/yamls/chaosmeta.yaml
kubectl apply -f "${ROOT_PATH}"/yamls/chaosmeta-daemonset.yaml
sh "${ROOT_PATH}"/build/build.sh
