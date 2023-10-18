#!/bin/bash

set -e

COMPONENT=$1
NAMESPACE=$2

BUILD_DIR=`cd $(dirname $0); pwd`
BUILD_DIR=${BUILD_DIR}/ssl/${COMPONENT}
echo ${BUILD_DIR}
mkdir -p ${BUILD_DIR} && cd ${BUILD_DIR}
docker run --mount type=bind,source=$(pwd),destination=${BUILD_DIR} registry.cn-hangzhou.aliyuncs.com/chaosmeta/chaosmeta-openssl:v1.0.0 openssl req -x509 -newkey rsa:4096 -keyout ${BUILD_DIR}/tls.key -out ${BUILD_DIR}/tls.crt -days 3650 -nodes -subj "/CN=chaosmeta-${COMPONENT}-webhook-service.${NAMESPACE}.svc" -addext "subjectAltName=DNS:chaosmeta-${COMPONENT}-webhook-service.${NAMESPACE}.svc"
caBundle=""
if [ "$(uname -s)" = "Linux" ]; then
    caBundle=$(cat tls.crt | base64 -w 0)
elif [ "$(uname -s)" = "Darwin" ]; then
    caBundle=$(base64 -i tls.crt -o - | tr -d '\n')
else
    echo "Unknown environment"
    exit 1
fi

kubectl create secret tls chaosmeta-${COMPONENT}-webhook-server-cert --cert=tls.crt --key=tls.key -n ${NAMESPACE}
kubectl patch MutatingWebhookConfiguration chaosmeta-${COMPONENT}-mutating-webhook-configuration --type='json' -p='[{"op": "add", "path": "/webhooks/0/clientConfig/caBundle", "value": "'"${caBundle}"'"}]'
kubectl patch ValidatingWebhookConfiguration chaosmeta-${COMPONENT}-validating-webhook-configuration --type='json' -p='[{"op": "add", "path": "/webhooks/0/clientConfig/caBundle", "value": "'"${caBundle}"'"}]'
