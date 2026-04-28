#!/bin/bash

# =============================================================================
# website-probe 多平台 Docker 镜像构建脚本
# 支持平台: linux/amd64, linux/arm64
# =============================================================================

set -e

# 配置
IMAGE_NAME="${IMAGE_NAME:-website-probe}"
IMAGE_REGISTRY="${IMAGE_REGISTRY:-}"
VERSION="${VERSION:-latest}"
PLATFORMS="${PLATFORMS:-linux/amd64,linux/arm64}"
PUSH="${PUSH:-false}"
BUILD_DATE="$(date -u +'%Y-%m-%dT%H:%M:%SZ')"
GIT_COMMIT="$(git rev-parse --short HEAD 2>/dev/null || echo 'unknown')"

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

log_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

log_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# 使用帮助
usage() {
    cat << EOF
用法: $0 [选项]

选项:
    -n, --name          镜像名称 (默认: website-probe)
    -r, --registry      镜像仓库地址 (如: docker.io, ghcr.io/user)
    -v, --version       版本号 (默认: latest)
    -p, --platforms    目标平台 (默认: linux/amd64,linux/arm64)
    --push              构建完成后推送到仓库
    --no-cache          不使用缓存构建
    -h, --help          显示帮助信息

示例:
    # 构建本地镜像
    $0

    # 构建并推送
    $0 -r docker.io/myrepo -v 1.0.0 --push

    # 仅构建 amd64 平台
    $0 -p linux/amd64
EOF
    exit 0
}

# 解析参数
while [[ $# -gt 0 ]]; do
    case $1 in
        -n|--name)
            IMAGE_NAME="$2"
            shift 2
            ;;
        -r|--registry)
            IMAGE_REGISTRY="$2"
            shift 2
            ;;
        -v|--version)
            VERSION="$2"
            shift 2
            ;;
        -p|--platforms)
            PLATFORMS="$2"
            shift 2
            ;;
        --push)
            PUSH="true"
            shift
            ;;
        --no-cache)
            NO_CACHE="--no-cache"
            shift
            ;;
        -h|--help)
            usage
            ;;
        *)
            log_error "未知参数: $1"
            usage
            ;;
    esac
done

# 完整镜像标签
if [[ -n "$IMAGE_REGISTRY" ]]; then
    FULL_IMAGE="${IMAGE_REGISTRY}/${IMAGE_NAME}:${VERSION}"
else
    FULL_IMAGE="${IMAGE_NAME}:${VERSION}"
fi

# 检查 Docker Buildx 是否可用
check_buildx() {
    if ! docker buildx version &> /dev/null; then
        log_error "Docker Buildx 未安装或不可用"
        log_info "安装方法: docker buildx install"
        exit 1
    fi
}

# 初始化 buildx builder
init_builder() {
    log_info "初始化 Docker Buildx..."

    # 创建 (如果不存在) 并使用 builder
    docker buildx create --use --name multiarch-builder 2>/dev/null || \
    docker buildx use multiarch-builder

    # 验证 builder 可用
    docker buildx inspect --bootstrap &> /dev/null || {
        log_warn "正在创建新的 builder..."
        docker buildx create --name multiarch-builder --driver docker-container --driver-opt "image=moby/buildkit:buildx-stable-1"
        docker buildx use multiarch-builder
        docker buildx inspect --bootstrap
    }
}

# 构建镜像
build_image() {
    log_info "开始构建多平台镜像..."
    log_info "镜像: ${FULL_IMAGE}"
    log_info "平台: ${PLATFORMS}"
    log_info "Git Commit: ${GIT_COMMIT}"
    log_info "构建时间: ${BUILD_DATE}"

    # 构建参数
    BUILD_ARGS=(
        --platform "$PLATFORMS"
        --builder multiarch-builder
        -t "${FULL_IMAGE}"
        --build-arg "BUILD_DATE=${BUILD_DATE}"
        --build-arg "VERSION=${VERSION}"
        --build-arg "VCS_REF=${GIT_COMMIT}"
    )

    # 添加额外 tag (latest)
    if [[ "$VERSION" != "latest" ]]; then
        if [[ -n "$IMAGE_REGISTRY" ]]; then
            BUILD_ARGS+=("-t" "${IMAGE_REGISTRY}/${IMAGE_NAME}:latest")
        else
            BUILD_ARGS+=("-t" "${IMAGE_NAME}:latest")
        fi
    fi

    # 添加 no-cache 参数
    if [[ -n "$NO_CACHE" ]]; then
        BUILD_ARGS+=("$NO_CACHE")
    fi

    # 执行构建
    if [[ "$PUSH" == "true" ]]; then
        log_info "构建并推送到远程仓库..."
        docker buildx build "${BUILD_ARGS[@]}" --push .
    else
        log_info "构建本地镜像..."
        docker buildx build "${BUILD_ARGS[@]}" --load .
    fi
}

# 显示构建结果
show_result() {
    log_info "构建完成!"
    echo ""
    echo "镜像信息:"
    echo "  - 镜像名称: ${FULL_IMAGE}"
    echo "  - 支持平台: ${PLATFORMS}"
    echo "  - Git Commit: ${GIT_COMMIT}"
    echo ""

    if [[ "$PUSH" != "true" ]]; then
        log_info "本地镜像列表:"
        docker images "${IMAGE_NAME}" --format "table {{.Repository}}\t{{.Tag}}\t{{.Architecture}}\t{{.Size}}"
    fi
}

# 主流程
main() {
    check_buildx
    init_builder
    build_image
    show_result
}

main "$@"
