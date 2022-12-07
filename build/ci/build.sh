#!/bin/bash

# Copyright 2022-2023 Chaos Meta Authors.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# See the License for the specific language governing permissions and
# limitations under the License.

set -eu

OS_NAME=$1
if [ "${OS_NAME}" != "darwin" ] && [ "${OS_NAME}" != "windows" ] && [ "${OS_NAME}" != "linux" ]; then
  echo "please add args os：darwin/windows/linux. eg: sh build.sh linux"
  exit 1
fi

# base info
BUILD_NAME="chaosmetad"
VERSION="1.0.0"
BUILD_DATE=$(date "+%Y-%m-%d %H:%M:%S")

# env var
GO_TOOL="go"
ARCH_NAME="amd64"

# constant
BUILD_DATE_FLAG="@DATE@"
VERSION_FLAG="@VERSION@"

# tool name
CPU_BURN="chaosmeta_cpuburn"
CPU_LOAD="chaosmeta_cpuload"
DISK_BURN="chaosmeta_diskburn"
MEM_FILL="chaosmeta_memfill"
FD_FULL="chaosmeta_fd"
NPROC="chaosmeta_nproc"
NET_OCCUPY="chaosmeta_occupy"

TOOL_EXECNS="chaosmeta_execns"

# file path
CI_DIR=$(
  cd $(dirname $0)
  pwd
)

cd ${CI_DIR}/..
BUILD_DIR=$(pwd)
cd ${BUILD_DIR}/..
PROJECT_DIR=$(pwd)
PACKAGE_DIR=${BUILD_DIR}/package
OUTPUT_DIR=${BUILD_DIR}/${BUILD_NAME}
VERSION_DIR=${PROJECT_DIR}/pkg/version

# set version and date
cp ${VERSION_DIR}/version.go ${VERSION_DIR}/version.bak
sed -i "s/${VERSION_FLAG}/${VERSION}/g" ${VERSION_DIR}/version.go
sed -i "s/${BUILD_DATE_FLAG}/${BUILD_DATE}/g" ${VERSION_DIR}/version.go

# build main
mkdir -p ${OUTPUT_DIR}
cd ${PROJECT_DIR}/cmd
go mod tidy
CGO_ENABLED=1 GOOS=${OS_NAME} GOARCH=${ARCH_NAME} ${GO_TOOL} build -o ${OUTPUT_DIR}/${BUILD_NAME} ${PROJECT_DIR}/cmd/main.go
rm -rf ${VERSION_DIR}/version.go && mv ${VERSION_DIR}/version.bak ${VERSION_DIR}/version.go

# build tool
mkdir -p ${PACKAGE_DIR}/${OS_NAME}/tools
if [ "${OS_NAME}" == "linux" ]; then
  g++ ${PROJECT_DIR}/tools/${CPU_LOAD}.cpp -o ${PACKAGE_DIR}/${OS_NAME}/tools/${CPU_LOAD}
fi
CGO_ENABLED=1 GOOS=${OS_NAME} GOARCH=${ARCH_NAME} ${GO_TOOL} build -o ${PACKAGE_DIR}/${OS_NAME}/tools/${CPU_BURN} ${PROJECT_DIR}/tools/${CPU_BURN}.go
CGO_ENABLED=1 GOOS=${OS_NAME} GOARCH=${ARCH_NAME} ${GO_TOOL} build -o ${PACKAGE_DIR}/${OS_NAME}/tools/${DISK_BURN} ${PROJECT_DIR}/tools/${DISK_BURN}.go
CGO_ENABLED=1 GOOS=${OS_NAME} GOARCH=${ARCH_NAME} ${GO_TOOL} build -o ${PACKAGE_DIR}/${OS_NAME}/tools/${MEM_FILL} ${PROJECT_DIR}/tools/${MEM_FILL}.go
CGO_ENABLED=1 GOOS=${OS_NAME} GOARCH=${ARCH_NAME} ${GO_TOOL} build -o ${PACKAGE_DIR}/${OS_NAME}/tools/${NET_OCCUPY} ${PROJECT_DIR}/tools/${NET_OCCUPY}.go
CGO_ENABLED=1 GOOS=${OS_NAME} GOARCH=${ARCH_NAME} ${GO_TOOL} build -o ${PACKAGE_DIR}/${OS_NAME}/tools/${FD_FULL} ${PROJECT_DIR}/tools/${FD_FULL}.go
CGO_ENABLED=1 GOOS=${OS_NAME} GOARCH=${ARCH_NAME} ${GO_TOOL} build -o ${PACKAGE_DIR}/${OS_NAME}/tools/${NPROC} ${PROJECT_DIR}/tools/${NPROC}.go
CGO_ENABLED=1 GOOS=${OS_NAME} GOARCH=${ARCH_NAME} ${GO_TOOL} build -o ${PACKAGE_DIR}/${OS_NAME}/tools/${TOOL_EXECNS} ${PROJECT_DIR}/tools/${TOOL_EXECNS}.go

cp -R ${PACKAGE_DIR}/${OS_NAME}/tools ${OUTPUT_DIR}/
