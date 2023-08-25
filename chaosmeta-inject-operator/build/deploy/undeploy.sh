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
    mkdir -p "${ROOT_PATH}"/yamls
    filename="${ROOT_PATH}"/yamls/chaosmeta-inject.yaml

    if [ ! -e "${filename}" ]; then
        downloadfile https://raw.githubusercontent.com/traas-stack/chaosmeta/"${VERSION}"/chaosmeta-inject-operator/build/yamls/chaosmeta-inject.yaml "${filename}"
    fi
fi

kubectl delete -f "${ROOT_PATH}"/yamls/chaosmeta-inject.yaml
