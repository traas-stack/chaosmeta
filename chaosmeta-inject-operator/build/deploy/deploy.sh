#!/bin/bash

set -e

NOW_DIR=`cd $(dirname $0); pwd`
VERSION=$1
if [ -z "$VERSION" ]; then
    echo "version is empty"
    exit 1
fi
echo "version: ${VERSION}"

function downloadfile() {
    url=$1
    filename=$2
    echo "download from: ${url}"
    if [ ! -e "${filename}" ]; then
        curl -o "${filename}" "${url}"
    fi
    echo "download success: ${filename}"
}

ROOT_PATH=${NOW_DIR}/../
if [ "${VERSION}" != "local" ]; then
    ROOT_PATH=${NOW_DIR}/chaosmeta_build
    echo "${ROOT_PATH}"
    mkdir -p "${ROOT_PATH}"/build "${ROOT_PATH}"/config "${ROOT_PATH}"/yamls
    preurl=https://raw.githubusercontent.com/traas-stack/chaosmeta/"${VERSION}"/chaosmeta-inject-operator

    downloadfile "${preurl}"/config/chaosmeta-inject.json "${ROOT_PATH}"/config/chaosmeta-inject.json
    downloadfile "${preurl}"/build/build.sh "${ROOT_PATH}"/build/build.sh
    downloadfile "${preurl}"/build/yamls/chaosmeta.yaml "${ROOT_PATH}"/yamls/chaosmeta.yaml
    downloadfile "${preurl}"/build/yamls/chaosmeta-daemonset.yaml "${ROOT_PATH}"/yamls/chaosmeta-daemonset.yaml
fi

kubectl apply -f "${ROOT_PATH}"/yamls/chaosmeta.yaml
kubectl apply -f "${ROOT_PATH}"/yamls/chaosmeta-daemonset.yaml
if [ "${VERSION}" != "local" ]; then
    sh "${ROOT_PATH}"/build/build.sh
else
    sh "${ROOT_PATH}"/build.sh
fi
