#!/bin/bash

# =============================================================================
# 多平台 Docker 镜像构建脚本
# 支持: 本地构建、远程推送、离线 OCI 导出、中转站模式
# 整合 build-docker.sh + buildx_image.sh 优点
# =============================================================================

set -e

# ------------------------------ 默认配置 ----------------------------------
IMAGE_NAME="sun-panel"
IMAGE_VERSION="latest"
IMAGE_REGISTRY=""
PLATFORMS="linux/amd64,linux/arm64,linux/arm/v7"
PUSH_MODE=""          # ""=本地加载, "push"=推送, "save"=导出OCI
EXTRA_TAGS=()
CACHE_ENABLED=false
AUTO_LATEST=true
BUILDER_NAME="multiarch-builder"
PROVENANCE=false
VIA_REGISTRY=""       # 中转仓库
CACHE_DIR=""          # 缓存目录 (默认通过环境变量或相对路径)
BUILD_DATE="$(date -u +'%Y-%m-%dT%H:%M:%SZ')"
GIT_COMMIT="$(git rev-parse --short HEAD 2>/dev/null || echo 'unknown')"
CACHE_DIR="${DOCKER_BUILDX_CACHE_DIR:-$(cd "$(dirname "$0")/.." && pwd)/.buildx-cache}"

# -------------------------------- 颜色定义 ----------------------------------
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
NC='\033[0m'

# -------------------------------- 日志函数 ----------------------------------
log_info()    { echo -e "${GREEN}[INFO]${NC} $1"; }
log_warn()    { echo -e "${YELLOW}[WARN]${NC} $1"; }
log_error()   { echo -e "${RED}[ERROR]${NC} $1"; }
log_step()    { echo -e "${BLUE}[STEP]${NC} $1"; }
log_success() { echo -e "${GREEN}[OK]${NC} $1"; }

# -------------------------------- 帮助信息 ----------------------------------
show_help() {
    cat <<- 'ENDHELP'
	========================================
	  多平台 Docker 镜像构建脚本
	========================================

	用法: build.sh [选项]

	镜像配置 (必选):
	  -n, --name <name>        镜像名称 (默认: sun-panel)
	  -v, --version <ver>     版本标签 (默认: latest)

	输出模式 (互斥，默认本地加载):
	  --push                   构建并推送到仓库
	  --save                   构建并导出 OCI 到本地目录

	可选配置:
	  -r, --registry <addr>    镜像仓库地址
	                           例: docker.io/user, registry.cn-hangzhou.aliyuncs.com/ns
	  -p, --platforms <list>   目标平台，逗号分隔 (默认: linux/amd64,linux/arm64,linux/arm/v7)
	  -t, --tag <tag>          额外 tag，可多次使用
	                           例: -t latest -t stable
	  --cache                  启用构建缓存
	  --no-cache               禁用缓存 (默认)
	  --cache-dir <path>       缓存目录 (默认: .buildx-cache)
	  --no-latest              不自动打 latest 标签 (测试版本用)
	  --builder <name>         指定 buildx builder (默认: multiarch-builder)
	  --via <registry>          中转仓库，save 模式后可从中转站同步
	                           例: --via docker.io/user
	  --provenance             生成 provenance (默认关闭)
	  -h, --help               显示帮助

	示例:
	  # 本地构建
	  ./build.sh -n myapp -v 1.0.0

	  # 推送到阿里云
	  ./build.sh -n myapp -v 1.0.0 -r registry.cn-hangzhou.aliyuncs.com/ns --push

	  # 离线导出 OCI
	  ./build.sh -n myapp -v 1.0.0 --save

	  # 测试版本，不更新 latest
	  ./build.sh -n myapp -v 1.0.0-beta --no-latest

	  # 单平台构建
	  ./build.sh -n myapp -v 1.0.0 -p linux/amd64

	  # 推送并打多个 tag
	  ./build.sh -n myapp -v 1.0.0 -r docker.io/user --push -t latest

	  # 中转站模式 (save 后可手动同步)
	  ./build.sh -n myapp -v 1.0.0 --save --via docker.io/user

	环境准备:
	  - 确保已安装 docker buildx: docker buildx version
	  - 推送前请先登录: docker login <registry>
	  - 离线推送需安装 skopeo: sudo apt-get install skopeo

	ENDHELP
    exit 0
}

# -------------------------------- 参数解析 ----------------------------------
parse_args() {
    while [[ $# -gt 0 ]]; do
        case $1 in
            -n|--name)
                IMAGE_NAME="$2"; shift 2 ;;
            -v|--version)
                IMAGE_VERSION="$2"; shift 2 ;;
            -r|--registry)
                IMAGE_REGISTRY="$2"; shift 2 ;;
            -p|--platforms)
                PLATFORMS="$2"; shift 2 ;;
            -t|--tag)
                EXTRA_TAGS+=("$2"); shift 2 ;;
            --push)
                PUSH_MODE="push"; shift ;;
            --save)
                PUSH_MODE="save"; shift ;;
            --cache)
                CACHE_ENABLED=true; shift ;;
            --no-cache)
                CACHE_ENABLED=false; shift ;;
            --no-latest)
                AUTO_LATEST=false; shift ;;
            --builder)
                BUILDER_NAME="$2"; shift 2 ;;
            --via)
                VIA_REGISTRY="$2"; shift 2 ;;
            --cache-dir)
                CACHE_DIR="$2"; shift 2 ;;
            --provenance)
                PROVENANCE=true; shift ;;
            -h|--help)
                show_help ;;
            *)
                log_error "未知参数: $1"; show_help ;;
        esac
    done
}

# -------------------------------- 镜像标签 ----------------------------------
get_full_image() {
    local tag="$1"
    if [[ -n "$IMAGE_REGISTRY" ]]; then
        echo "${IMAGE_REGISTRY}/${IMAGE_NAME}:${tag}"
    else
        echo "${IMAGE_NAME}:${tag}"
    fi
}

# -------------------------------- 构建参数 ----------------------------------
build_args=()
EXTRA_TAGS_FULL=()

prepare_build_args() {
    # 基础构建参数
    build_args=(
        --platform "$PLATFORMS"
        --builder "$BUILDER_NAME"
    )

    # 镜像标签
    build_args+=(-t "$(get_full_image "$IMAGE_VERSION")")

    # 额外 tags
    for tag in "${EXTRA_TAGS[@]}"; do
        local full_tag="$(get_full_image "$tag")"
        build_args+=(-t "$full_tag")
        EXTRA_TAGS_FULL+=("$full_tag")
    done

    # 非 latest 版本额外打 latest tag (可通过 --no-latest 禁用)
    if [[ "$AUTO_LATEST" == "true" ]] && [[ "$IMAGE_VERSION" != "latest" ]] && [[ ! " ${EXTRA_TAGS[*]} " =~ " latest " ]]; then
        build_args+=(-t "$(get_full_image "latest")")
        EXTRA_TAGS_FULL+=("$(get_full_image "latest")")
    fi

    # 构建参数
    build_args+=(
        --build-arg "BUILD_DATE=${BUILD_DATE}"
        --build-arg "VERSION=${IMAGE_VERSION}"
        --build-arg "VCS_REF=${GIT_COMMIT}"
    )

    # Provenance
    if [[ "$PROVENANCE" == "true" ]]; then
        build_args+=(--provenance=true)
    else
        build_args+=(--provenance=false)
    fi

    # 缓存
    if [[ "$CACHE_ENABLED" == "true" ]]; then
        local cache_dest="${CACHE_DIR:-${DOCKER_BUILDX_CACHE_DIR:-$(cd "$(dirname "$0")/.." && pwd)/.buildx-cache}}"
        build_args+=(
            --cache-to "type=local,dest=${cache_dest},mode=max"
            --cache-from "type=local,src=${cache_dest}"
        )
    fi
}

# -------------------------------- 检查 ----------------------------------
check_buildx() {
    if ! docker buildx version &> /dev/null; then
        log_error "Docker Buildx 未安装或不可用"
        log_info "安装方法: docker buildx install"
        exit 1
    fi
    log_success "Docker Buildx 已就绪"
}

# -------------------------------- Builder 初始化 ---------------------------
init_builder() {
    log_step "初始化 Buildx Builder..."

    # 尝试使用已有 builder 或创建
    if ! docker buildx inspect "$BUILDER_NAME" &> /dev/null; then
        log_info "创建新 builder: $BUILDER_NAME"
        docker buildx create --name "$BUILDER_NAME" \
            --driver docker-container \
            --driver-opt "image=moby/buildkit:buildx-stable-1" \
            --bootstrap
    fi

    docker buildx use "$BUILDER_NAME"

    # Bootstrap
    docker buildx inspect --bootstrap &> /dev/null || {
        log_warn "正在 Bootstrap builder..."
        docker buildx inspect --bootstrap
    }

    log_success "Builder '$BUILDER_NAME' 已就绪"
}

# -------------------------------- 构建本地 ----------------------------------
build_load() {
    log_step "构建镜像 (本地加载模式)..."
    log_info "镜像: $(get_full_image "$IMAGE_VERSION")"
    log_info "平台: $PLATFORMS"
    log_info "缓存: $CACHE_ENABLED"

    prepare_build_args
    build_args+=(--load)

    echo ""
    docker buildx build "${build_args[@]}" .
}

# -------------------------------- 构建推送 ----------------------------------
build_push() {
    log_step "构建并推送镜像..."
    log_info "镜像: $(get_full_image "$IMAGE_VERSION")"
    log_info "仓库: ${IMAGE_REGISTRY:-<本地>}"
    log_info "平台: $PLATFORMS"

    prepare_build_args
    build_args+=(--push)

    echo ""
    docker buildx build "${build_args[@]}" .
}

# -------------------------------- 构建 OCI ----------------------------------
build_save() {
    local oci_dir="$(cd "$(dirname "$0")/.." && pwd)/output/${IMAGE_NAME}-${IMAGE_VERSION}"
    local oci_tar="${IMAGE_NAME}-oci.tar"

    log_step "构建并导出 OCI 镜像..."
    log_info "镜像: $(get_full_image "$IMAGE_VERSION")"
    log_info "平台: $PLATFORMS"
    log_info "输出目录: $oci_dir"

    prepare_build_args
    build_args+=(--output "type=oci,dest=${oci_tar}")

    echo ""
    docker buildx build "${build_args[@]}" .

    # 解压 tar 为目录
    log_step "解压 OCI 镜像..."
    rm -rf "$oci_dir"
    mkdir -p "$oci_dir"
    tar -xf "$oci_tar" -C "$oci_dir/"
    rm -f "$oci_tar"

    log_success "OCI 导出完成: $oci_dir"
    echo ""

    # 推送指令
    print_push_instructions "$oci_dir"
}

# -------------------------------- 推送指引 ----------------------------------
print_push_instructions() {
    local oci_dir="$1"

    echo -e "${CYAN}========================================${NC}"
    echo -e "${CYAN}  后续推送操作指引${NC}"
    echo -e "${CYAN}========================================${NC}"
    echo ""

    local target_registry="${IMAGE_REGISTRY:-docker.io/library}"
    local target_image="${target_registry}/${IMAGE_NAME}:${IMAGE_VERSION}"

    echo -e "${GREEN}OCI 镜像位置:${NC} $oci_dir"
    echo ""

    echo -e "${YELLOW}1. 推送到当前配置仓库:${NC}"
    echo "   skopeo copy --all \\"
    echo "     oci:$oci_dir \\"
    echo "     docker://${target_image}"
    echo ""

    if [[ -n "$VIA_REGISTRY" ]]; then
        local via_image="${VIA_REGISTRY}/${IMAGE_NAME}:${IMAGE_VERSION}"
        echo -e "${YELLOW}2. 通过中转站推送到 Docker Hub:${NC}"
        echo "   skopeo copy --all \\"
        echo "     oci:$oci_dir \\"
        echo "     docker://${via_image}"
        echo ""
        echo -e "${BLUE}   或使用 --dest-creds 指定认证:${NC}"
        echo "   skopeo copy --all --dest-creds <user>:<password> \\"
        echo "     oci:$oci_dir \\"
        echo "     docker://${via_image}"
        echo ""
    fi

    echo -e "${CYAN}提示:${NC} 确保已安装 skopeo: sudo apt-get install skopeo"
    echo ""
}

# -------------------------------- 结果展示 ----------------------------------
show_result() {
    echo ""
    echo -e "${CYAN}========================================${NC}"
    echo -e "${CYAN}  构建完成${NC}"
    echo -e "${CYAN}========================================${NC}"
    echo ""
    echo "  镜像名称: ${IMAGE_NAME}"
    echo "  版本: ${IMAGE_VERSION}"
    echo "  仓库: ${IMAGE_REGISTRY:-<本地>}"
    echo "  平台: ${PLATFORMS}"
    echo "  Git: ${GIT_COMMIT}"
    echo "  构建时间: ${BUILD_DATE}"
    echo ""

    if [[ -n "${EXTRA_TAGS_FULL[*]}" ]]; then
        echo "  额外 Tags:"
        for tag in "${EXTRA_TAGS_FULL[@]}"; do
            echo "    - $tag"
        done
        echo ""
    fi

    if [[ "$PUSH_MODE" != "push" ]]; then
        log_info "本地镜像列表:"
        docker images "${IMAGE_NAME}" --format "table {{.Repository}}\t{{.Tag}}\t{{.Architecture}}\t{{.Size}}"
    fi
}

# -------------------------------- 平台数量检测 ----------------------------------
count_platforms() {
    echo "$PLATFORMS" | tr ',' '\n' | wc -l | tr -d ' '
}

# -------------------------------- 主流程 ------------------------------------
main() {
    echo -e "${CYAN}========================================${NC}"
    echo -e "${CYAN}  Docker 多平台镜像构建${NC}"
    echo -e "${CYAN}========================================${NC}"
    echo ""

    parse_args "$@"

    # 验证
    if [[ -z "$IMAGE_NAME" ]]; then
        log_error "镜像名称不能为空"
        exit 1
    fi

    # 多平台 + 本地模式时自动切换到 save
    local platform_count=$(count_platforms)
    if [[ "$PUSH_MODE" != "push" ]] && [[ "$PUSH_MODE" != "save" ]] && [[ $platform_count -gt 1 ]]; then
        log_warn "Docker --load 不支持多平台构建，将自动切换到 --save 模式"
        log_info "如需推送，请使用 --push 选项"
        echo ""
        PUSH_MODE="save"
    fi

    # 显示配置
    log_info "镜像: $(get_full_image "$IMAGE_VERSION")"
    log_info "平台: $PLATFORMS"
    log_info "模式: ${PUSH_MODE:-local}"
    echo ""

    # 检查并初始化
    check_buildx
    init_builder

    # 执行构建
    case "$PUSH_MODE" in
        push)
            build_push
            ;;
        save)
            build_save
            ;;
        *)
            build_load
            ;;
    esac

    show_result
}

main "$@"
