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
BASE_DIR=$(
  cd $(dirname $0)
  pwd
)

# param check
OP=$1
NAMESPACE=$2
COMPONENT_LIST=()
for ((i = 3; i <= $#; i++)); do
  COMPONENT_LIST+=("${!i}")
done

if [ "$OP" == "-h" ] || [ "$OP" == "--help" ] || [ -z "$OP" ] || [ ${#COMPONENT_LIST[@]} -eq 0 ]; then
  echo "format: sh deploy.sh [operation] [namespace] [component list]"
  echo "operation support: ${OP_INSTALL},${OP_UNINSTALL}"
  echo "component support: ${COMPONENT_ALL} ${COMPONENT_PLATFORM} ${COMPONENT_WORKFLOW} ${COMPONENT_INJECT} ${COMPONENT_MEASURE} ${COMPONENT_FLOW} ${COMPONENT_DAEMON}"
  echo "example: sh deploy.sh install default ${COMPONENT_PLATFORM} ${COMPONENT_WORKFLOW} ${COMPONENT_INJECT} ${COMPONENT_DAEMON}"
  exit 1
fi

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

for COMPONENT in "${COMPONENT_LIST[@]}"; do
  if [ ${#COMPONENT_LIST[@]} -ne 1 ] && [ "$COMPONENT" == "$COMPONENT_ALL" ]; then
    echo "component[$COMPONENT_ALL] should not in list with other components"
    exit 3
  fi

  if [[ "$COMPONENT" != "$COMPONENT_WORKFLOW" && "$COMPONENT" != "$COMPONENT_DAEMON" && "$COMPONENT" != "$COMPONENT_ALL" && "$COMPONENT" != "$COMPONENT_PLATFORM" && "$COMPONENT" != "$COMPONENT_INJECT" && "$COMPONENT" != "$COMPONENT_MEASURE" && "$COMPONENT" != "$COMPONENT_FLOW" ]]; then
    echo "not support component: $COMPONENT"
    echo "component support: ${COMPONENT_ALL} ${COMPONENT_PLATFORM} ${COMPONENT_WORKFLOW} ${COMPONENT_INJECT} ${COMPONENT_MEASURE} ${COMPONENT_FLOW} ${COMPONENT_DAEMON}"
    echo "example: sh deploy.sh install default ${COMPONENT_PLATFORM} ${COMPONENT_WORKFLOW} ${COMPONENT_INJECT} ${COMPONENT_DAEMON}"
    exit 1
  fi
done

function installComponent() {
  sed "s/${NAMESPACE_REPLACE}/$2/g" ${BASE_DIR}/templates/chaosmeta-$1-template.yaml >${BASE_DIR}/yamls/chaosmeta-$1.yaml
  if [[ "$OP" == "$OP_INSTALL" ]]; then
    kubectl apply -f ${BASE_DIR}/yamls/chaosmeta-$1.yaml
    if [[ "$1" != "$COMPONENT_PLATFORM" && "$1" != "$COMPONENT_DAEMON" && "$1" != "$COMPONENT_WORKFLOW" ]]; then
      sh tools/build.sh $1 $2
    fi
  else
    kubectl delete -f ${BASE_DIR}/yamls/chaosmeta-$1.yaml
    if [[ "$1" != "$COMPONENT_PLATFORM" && "$1" != "$COMPONENT_DAEMON" && "$1" != "$COMPONENT_WORKFLOW" ]]; then
      kubectl delete secret chaosmeta-$1-webhook-server-cert -n $2
    fi
  fi
}

# execute install
for COMPONENT in "${COMPONENT_LIST[@]}"; do
  if [[ "$COMPONENT" == "$COMPONENT_ALL" ]]; then
    installComponent $COMPONENT_PLATFORM $NAMESPACE
    installComponent $COMPONENT_WORKFLOW $NAMESPACE
    installComponent $COMPONENT_INJECT $NAMESPACE
    installComponent $COMPONENT_DAEMON $NAMESPACE
    installComponent $COMPONENT_MEASURE $NAMESPACE
    installComponent $COMPONENT_FLOW $NAMESPACE
  else
    installComponent $COMPONENT $NAMESPACE
  fi
done
