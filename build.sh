#!/bin/bash

# APPNAME=$1
TAG=$1
# # 打印脚本名称
# echo "脚本名称：$0"

# # 打印参数数量
# echo "参数数量：$#"

# # 打印所有参数
# echo "所有参数：$@"

# # 循环打印每个参数
# echo "逐个参数打印："
# for arg in "$@"; do
#     echo "$arg"
# done

# 设置脚本在任何命令失败时立即退出
set -e

# $1 镜像库地址（最后不能带斜杠） $2 镜像名 $3 标签
build_push() {
  library=$1
  name=$2
  tag=$3
  

  # 推送最新标签
  echo "=== docker tag $name"
  docker tag $name ${library}$name
  echo "=== docker push $name:latest"
  docker push ${library}$name

  # 推送设定标签
  echo "=== docker tag $name:$tag"
  docker tag $name ${library}$name:$tag
  echo "=== docker push $name:$tag"
  docker push ${library}$name:$tag
}


build_sunPayAuthServer() {
  lib=enianteam-docker.pkg.coding.net/sun-panel/pay/
  name=sun-pay-auth-server
  tag=$1
  echo "=== docker build $name"
  docker build -t $name .
  build_push $lib $name $tag
}



build_sunPayAuthServer $TAG
