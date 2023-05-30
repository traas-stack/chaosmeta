#!/bin/bash

set -e

ROOT_PATH=$(dirname $(readlink -f $0))/..
echo "${ROOT_PATH}"

kubectl create configmap chaosmeta-inject-config --from-file="${ROOT_PATH}"/config/chaosmeta-inject.json -n chaosmeta

BUILD_DIR="/tmp/chaosmeta_build"
mkdir -p ${BUILD_DIR}/ssl && cd ${BUILD_DIR}/ssl
docker run --mount type=bind,source=$(pwd),destination=/data registry.cn-hangzhou.aliyuncs.com/chaosmeta/chaosmeta-openssl:v1.0.0 openssl req -x509 -newkey rsa:4096 -keyout /data/tls.key -out /data/tls.crt -days 3650 -nodes -subj "/CN=chaosmeta-inject-webhook-service.chaosmeta.svc" -addext "subjectAltName=DNS:chaosmeta-inject-webhook-service.chaosmeta.svc"
caBundle=$(cat tls.crt | base64 -w 0)
kubectl create secret tls webhook-server-cert --cert=tls.crt --key=tls.key -n chaosmeta
kubectl patch MutatingWebhookConfiguration chaosmeta-inject-mutating-webhook-configuration --type='json' -p='[{"op": "add", "path": "/webhooks/0/clientConfig/caBundle", "value": "'"${caBundle}"'"}]'
kubectl patch ValidatingWebhookConfiguration chaosmeta-inject-validating-webhook-configuration --type='json' -p='[{"op": "add", "path": "/webhooks/0/clientConfig/caBundle", "value": "'"${caBundle}"'"}]'
