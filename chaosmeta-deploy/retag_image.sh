#!/bin/bash

BASE_DIR=$(
  cd $(dirname $0)
  pwd
)
source_registry="registry.cn-hangzhou.aliyuncs.com/chaosmeta"
target_registry=""
REGISTRY_REPLACE="DEPLOYREGISTRY"
function process_args() {
  while getopts ":r:h" opt; do
    case $opt in
      r)
        target_registry=$OPTARG
        ;;
      h)
        echo "-r：target image registry"
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
if [[ -z "$target_registry" ]]; then
  echo "please provide args: -r"
  exit 1
fi

echo "target registry: $target_registry"
image_list=$(cat templates/* | grep 'image:' | awk -F'image: ' '{print $2}' | tr '\n' ',')
IFS=','
image_array=($image_list)
source_registry=${source_registry//\//\\/}
target_registry=${target_registry//\//\\/}

for image in "${image_array[@]}"; do
  source_image=$(echo "$image" | sed "s/$REGISTRY_REPLACE/$source_registry/g")
  dst_image=$(echo "$image" | sed "s/$REGISTRY_REPLACE/$target_registry/g")
  echo "source image: $source_image"
  echo "target image: $dst_image"
  docker pull $source_image
  docker tag $source_image $dst_image
  docker push $dst_image
done
