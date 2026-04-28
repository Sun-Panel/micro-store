#!/bin/bash

# ============================================
# 多平台 Docker 镜像构建与推送脚本
# 支持网络不稳定环境下分步执行构建和推送
# ============================================

# 参数说明:
# $1: TAG - 镜像标签 (必需)
# $2: MODE - 构建模式 (可选，默认：push_aly)
#   - push_aly: 构建并推送到阿里云仓库 (中转站)
#   - push_dh:  构建并直接推送到 Docker Hub
#   - save:     仅构建并保存到本地 OCI 目录，不推送
# $3: DH_TAG - Docker Hub 标签 (可选，默认与 TAG 相同)

TAG=$1
MODE=${2:-push_aly}
DH_TAG=${3:-$TAG}

# 设置脚本在任何命令失败时立即退出
set -e

# 缓存目录（项目本地，跨平台兼容）
CACHE_DIR="$(cd "$(dirname "$0")" && pwd)/.buildx-cache"

build_image() {
  alyLib="registry.cn-hangzhou.aliyuncs.com/hslr/"
  dhLib="docker.io/hslr/"  # 替换为你的 Docker Hub 用户名
  name=sun-panel
  tag=$1
  mode=$2
  
  echo "============================================"
  echo "=== Docker Build: $name"
  echo "=== Tag: $tag"
  echo "=== Mode: $mode"
  echo "============================================"

  case "$mode" in
    "save")
      # 方案一：保存为本地 OCI 目录 (适合网络不稳定时先构建)
      # 注意：buildx 的 type=oci 会输出 tar 文件，需要解压
      
      echo "📁 开始构建并保存 OCI 镜像..."
      echo ""
      
      # 构建 OCI tar 文件
      # --provenance=false: 禁用 provenance attestation 生成，避免 skopeo 压缩错误
      docker buildx build \
        --cache-to type=local,dest="$CACHE_DIR",mode=max \
        --cache-from type=local,src="$CACHE_DIR" \
        --platform=linux/amd64 \
        --platform=linux/arm64 \
        --platform=linux/arm/v7 \
        --provenance=false \
        --output type=oci,dest=./${name}-oci.tar \
        -t ${name}:${tag} .
      
      # 创建目录并解压
      rm -rf ./${name}-oci
      mkdir -p ./${name}-oci
      tar -xf ./${name}-oci.tar -C ./${name}-oci/
      
      # 删除 tar 文件（可选）
      rm -f ./${name}-oci.tar
      
      echo ""
      echo "=== ✅ OCI 保存成功：./${name}-oci/"
      echo ""
      echo "=== 📦 验证 OCI 结构:"
      ls -lh ./${name}-oci/ | head -10
      echo ""
      echo "=== 📦 后续推送操作指引:"
      echo ""
      echo "1️⃣  推送到 Docker Hub (使用 skopeo):"
      echo "   skopeo copy --all \\"
      echo "     oci:./${name}-oci/ \\"
      echo "     docker://${dhLib}${name}:${DH_TAG}"
      echo ""
      echo "2️⃣  推送到阿里云 (如果需要):"
      echo "   skopeo copy --all \\"
      echo "     oci:./${name}-oci/ \\"
      echo "     docker://${alyLib}${name}:${tag}"
      echo ""
      ;;
    
    "push_dh")
      # 方案二：直接推送到 Docker Hub
      echo "⚠️  注意：直接推送 Docker Hub 可能因网络原因失败..."
      echo "💡 建议：如网络不稳定，请使用 'save' 模式先构建，再手动推送"
      echo ""
      
      docker buildx build \
        --cache-to type=local,dest="$CACHE_DIR",mode=max \
        --cache-from type=local,src="$CACHE_DIR" \
        --platform=linux/amd64 \
        --platform=linux/arm64 \
        --platform=linux/arm/v7 \
        --provenance=false \
        --push -t ${dhLib}${name}:${DH_TAG} .
      
      echo ""
      echo "=== ✅ 已推送到 Docker Hub: ${dhLib}${name}:${DH_TAG}"
      echo ""
      ;;
    
    "push_aly"|"push")
      # 方案三：推送到阿里云仓库 (推荐作为中转站)
      docker buildx build \
        --cache-to type=local,dest="$CACHE_DIR",mode=max \
        --cache-from type=local,src="$CACHE_DIR" \
        --platform=linux/amd64 \
        --platform=linux/arm64 \
        --platform=linux/arm/v7 \
        --provenance=false \
        --push -t ${alyLib}${name}:${tag} .
      
      echo ""
      echo "=== ✅ 已推送到阿里云：${alyLib}${name}:${tag}"
      echo ""
      echo "=== 📦 从中转站复制到 Docker Hub (当网络稳定时执行):"
      echo ""
      echo "skopeo copy --all --dest-creds <docker_hub_user>:<password> \\"
      echo "  docker://${alyLib}${name}:${tag} \\"
      echo "  docker://${dhLib}${name}:${DH_TAG}"
      echo ""
      ;;
    
    *)
      echo "❌ 错误：未知的模式 '$mode'"
      echo ""
      echo "可用模式:"
      echo "  push_aly  - 推送到阿里云 (默认)"
      echo "  push_dh   - 直接推送到 Docker Hub"
      echo "  save      - 保存到本地 OCI 目录"
      exit 1
      ;;
  esac
}

# 显示帮助信息
show_help() {
  echo "用法：$0 <TAG> [MODE] [DH_TAG]"
  echo ""
  echo "参数说明:"
  echo "  TAG       镜像标签 (例如：v1.0.0, latest)"
  echo "  MODE      构建模式 (可选，默认：push_aly)"
  echo "            - push_aly: 构建并推送到阿里云仓库"
  echo "            - push_dh:  构建并直接推送到 Docker Hub"
  echo "            - save:     仅构建并保存到本地 OCI 目录"
  echo "  DH_TAG    Docker Hub 标签 (可选，默认与 TAG 相同)"
  echo ""
  echo "示例:"
  echo "  $0 1.0.0                    # 推送到阿里云，标签 1.0.0"
  echo "  $0 1.0.0 save              # 仅保存到本地 OCI"
  echo "  $0 1.0.0 push_dh           # 直接推送到 Docker Hub"
  echo "  $0 1.0.0 push_aly latest   # 阿里云用 1.0.0，Docker Hub 用 latest"
  echo ""
  echo "环境准备:"
  echo "  1. 确保已登录对应仓库：docker login"
  echo "  2. 安装 skopeo 工具 (用于后续推送): sudo apt-get install skopeo"
  echo ""
}

# 主逻辑
if [ "$1" == "-h" ] || [ "$1" == "--help" ] || [ -z "$1" ]; then
  show_help
  exit 0
fi

echo "🚀 开始构建..."
echo ""

build_image $TAG $MODE