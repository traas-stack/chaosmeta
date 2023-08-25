#!/bin/bash

# base var
COMPONENT_ALL="all"
COMPONENT_PLATFORM="platform"
COMPONENT_INJECT="inject"
COMPONENT_DAEMON="daemon"
COMPONENT_MEASURE="measure"
COMPONENT_FLOW="flow"

OP_INSTALL="install"
OP_UNINSTALL="uninstall"

NAMESPACE_REPLACE="DEPLOYNAMESPACE"
BASE_DIR=$(
  cd $(dirname $0)
  pwd
)

# param check
OP=$1
COMPONENT=$2
NAMESPACE=$3

if [[ -z "$OP" || -z "$COMPONENT" ]]; then
  echo "sh deploy.sh [operation] [component] [namespace]\n  operation support: ${OP_INSTALL},${OP_UNINSTALL}\n  component support: ${COMPONENT_ALL},${COMPONENT_PLATFORM},${COMPONENT_INJECT},${COMPONENT_MEASURE},${COMPONENT_FLOW},${COMPONENT_DAEMON}"
  exit 1
fi

if [[ "$OP" != "$OP_INSTALL" && "$OP" != "$OP_UNINSTALL" ]]; then
  echo "not support operation: $OP\n  operation support: ${OP_INSTALL},${OP_UNINSTALL}"
  exit 1
fi

if [[ "$COMPONENT" != "$COMPONENT_DAEMON" && "$COMPONENT" != "$COMPONENT_ALL" && "$COMPONENT" != "$COMPONENT_PLATFORM" && "$COMPONENT" != "$COMPONENT_INJECT" && "$COMPONENT" != "$COMPONENT_MEASURE" && "$COMPONENT" != "$COMPONENT_FLOW" ]]; then
  echo "not support component: $COMPONENT\n  component support: ${COMPONENT_ALL},${COMPONENT_PLATFORM},${COMPONENT_INJECT},${COMPONENT_MEASURE},${COMPONENT_FLOW},${COMPONENT_DAEMON}"
  exit 1
fi

if [[ -z "$NAMESPACE" ]]; then
  echo "not provide namespace, set to \"default\""
  NAMESPACE="default"
fi

# check namespace exist
if [[ $(kubectl get ns $NAMESPACE | wc -l) -eq 0 ]]; then
  echo "namespace[$NAMESPACE] not exist, please create ns first"
  exit 2
fi

function installComponent() {
  sed "s/${NAMESPACE_REPLACE}/$2/g" ${BASE_DIR}/templates/chaosmeta-$1-template.yaml >${BASE_DIR}/yamls/chaosmeta-$1.yaml
  if [[ "$OP" == "$OP_INSTALL" ]]; then
    kubectl apply -f ${BASE_DIR}/yamls/chaosmeta-$1.yaml
    if [[ "$COMPONENT" != "$COMPONENT_PLATFORM" && "$COMPONENT" != "$COMPONENT_DAEMON" ]]; then
      sh tools/build.sh $1 $2
    fi
  else
    kubectl delete -f ${BASE_DIR}/yamls/chaosmeta-$1.yaml
    if [[ "$COMPONENT" != "$COMPONENT_PLATFORM" && "$COMPONENT" != "$COMPONENT_DAEMON" ]]; then
      kubectl delete secret chaosmeta-$1-webhook-server-cert -n $2
    fi
  fi
}

# execute install
if [[ "$COMPONENT" == "$COMPONENT_ALL" ]]; then
  installComponent $COMPONENT_PLATFORM $NAMESPACE
  installComponent $COMPONENT_INJECT $NAMESPACE
  installComponent $COMPONENT_DAEMON $NAMESPACE
  installComponent $COMPONENT_MEASURE $NAMESPACE
  installComponent $COMPONENT_FLOW $NAMESPACE
else
  installComponent $COMPONENT $NAMESPACE
fi
