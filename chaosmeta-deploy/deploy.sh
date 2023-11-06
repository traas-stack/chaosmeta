#!/bin/bash

# base var
COMPONENT_ALL="all"
COMPONENT_PLATFORM="platform"
COMPONENT_WORKFLOW="workflow"
COMPONENT_INJECT="inject"
COMPONENT_DAEMON="daemon"
COMPONENT_MEASURE="measure"
COMPONENT_FLOW="flow"
OP_INSTALL="install"
OP_UNINSTALL="uninstall"
NAMESPACE_REPLACE="DEPLOYNAMESPACE"
REGISTRY_REPLACE="DEPLOYREGISTRY"
NAMESPACE="chaosmeta"
OP="install"
REGISTRY="registry.cn-hangzhou.aliyuncs.com/chaosmeta"
COMPONENT="all"

BASE_DIR=$(
  cd $(dirname $0)
  pwd
)

# function define
function installComponent() {
  TARGET_REGISTRY=${REGISTRY//\//\\/}
  sed "s/${NAMESPACE_REPLACE}/$NAMESPACE/g; s/${REGISTRY_REPLACE}/$TARGET_REGISTRY/g" ${BASE_DIR}/templates/chaosmeta-$1-template.yaml >${BASE_DIR}/yamls/chaosmeta-$1.yaml
  if [[ "$OP" == "$OP_INSTALL" ]]; then
    kubectl apply -f ${BASE_DIR}/yamls/chaosmeta-$1.yaml
    if [[ "$1" != "$COMPONENT_PLATFORM" && "$1" != "$COMPONENT_DAEMON" && "$1" != "$COMPONENT_WORKFLOW" ]]; then
      sh tools/build.sh $1 $NAMESPACE
    fi
  else
    kubectl delete -f ${BASE_DIR}/yamls/chaosmeta-$1.yaml
    if [[ "$1" != "$COMPONENT_PLATFORM" && "$1" != "$COMPONENT_DAEMON" && "$1" != "$COMPONENT_WORKFLOW" ]]; then
      kubectl delete secret chaosmeta-$1-webhook-server-cert -n $NAMESPACE
    fi
  fi
}

function process_args() {
  while getopts ":n:o:c:r:h" opt; do
    case $opt in
      n)
        NAMESPACE=$OPTARG
        ;;
      o)
        OP=$OPTARG
        ;;
      c)
        COMPONENT=$OPTARG
        ;;
      r)
        REGISTRY=$OPTARG
        ;;
      h)
        echo "-n：existed namespace of kubernetes, default: ${NAMESPACE}"
        echo "-r：image registry, default: ${REGISTRY}"
        echo "-o：${OP_INSTALL}、${OP_UNINSTALL}, default: ${OP}"
        echo "-c：${COMPONENT_ALL}、${COMPONENT_PLATFORM}、${COMPONENT_WORKFLOW}、${COMPONENT_INJECT}、${COMPONENT_DAEMON}、${COMPONENT_MEASURE}、${COMPONENT_FLOW}, default: ${COMPONENT_ALL}"
        exit 1
        ;;
      \?)
        echo "invalid option： -$OPTARG" >&2
        exit 1
        ;;
      :)
        echo "option -$OPTARG need an args" >&2
        exit 1
        ;;
    esac
  done
}

process_args "$@"
echo "namespace：$NAMESPACE"
echo "operation：$OP"
echo "image registry：$REGISTRY"
echo "component：$COMPONENT"
echo "================================================================"
sleep 3s

mkdir -p ${BASE_DIR}/yamls
# param check

if [[ "$OP" != "$OP_INSTALL" && "$OP" != "$OP_UNINSTALL" ]]; then
  echo "not support operation: $OP"
  echo "operation support: ${OP_INSTALL},${OP_UNINSTALL}"
  exit 1
fi

if [[ -z "$NAMESPACE" ]]; then
  echo "please provide namespace"
  exit 1
fi

# check namespace exist
if [[ $(kubectl get ns $NAMESPACE | wc -l) -eq 0 ]]; then
  echo "namespace[$NAMESPACE] not exist, please create ns first"
  exit 2
fi

if [[ "$COMPONENT" == "$COMPONENT_ALL" ]]; then
  installComponent $COMPONENT_PLATFORM
  installComponent $COMPONENT_WORKFLOW
  installComponent $COMPONENT_INJECT
  installComponent $COMPONENT_DAEMON
  installComponent $COMPONENT_MEASURE
  installComponent $COMPONENT_FLOW
else
  installComponent $COMPONENT
fi
