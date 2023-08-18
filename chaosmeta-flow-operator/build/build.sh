#!/bin/bash

set -e

BASE_DIR=`cd $(dirname $0); pwd`
ROOT_PATH=${BASE_DIR}/..
echo "${ROOT_PATH}"
ns="chaosmeta-flow"

BUILD_DIR="/tmp/chaosmeta_build"
mkdir -p ${BUILD_DIR}/ssl && cd ${BUILD_DIR}/ssl
docker run --mount type=bind,source=$(pwd),destination=${BUILD_DIR}/ssl registry.cn-hangzhou.aliyuncs.com/chaosmeta/chaosmeta-openssl:v1.0.0 openssl req -x509 -newkey rsa:4096 -keyout ${BUILD_DIR}/ssl/tls.key -out ${BUILD_DIR}/ssl/tls.crt -days 3650 -nodes -subj "/CN=chaosmeta-flow-webhook-service.${ns}.svc" -addext "subjectAltName=DNS:chaosmeta-flow-webhook-service.${ns}.svc"
caBundle=""
if [ "$(uname -s)" = "Linux" ]; then
    caBundle=$(cat tls.crt | base64 -w 0)
elif [ "$(uname -s)" = "Darwin" ]; then
    caBundle=$(base64 -i tls.crt -o - | tr -d '\n')
else
    echo "Unknown environment"
    exit 1
fi

kubectl create secret tls webhook-server-cert --cert=tls.crt --key=tls.key -n ${ns}
kubectl patch MutatingWebhookConfiguration chaosmeta-flow-mutating-webhook-configuration --type='json' -p='[{"op": "add", "path": "/webhooks/0/clientConfig/caBundle", "value": "'"${caBundle}"'"}]'
kubectl patch ValidatingWebhookConfiguration chaosmeta-flow-validating-webhook-configuration --type='json' -p='[{"op": "add", "path": "/webhooks/0/clientConfig/caBundle", "value": "'"${caBundle}"'"}]'
