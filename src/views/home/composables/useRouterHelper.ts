import { computed } from 'vue'
import { useRoute } from 'vue-router'

/**
 * 路由辅助 composable
 * 用于处理主路由和 iframe 路由的跳转逻辑
 * @param version iframe 版本号，默认为 'v1'
 */
export function useRouterHelper(version: string = 'v1') {
  const route = useRoute()

  /** iframe 路由前缀 */
  const iframePrefix = `/iframe/${version}`

  /** 当前是否在 iframe 路由下 */
  const isIframeRoute = computed(() => route.path.startsWith(iframePrefix))

  /**
   * 获取首页路径
   */
  const getHomePath = (): string => {
    return isIframeRoute.value ? `${iframePrefix}/` : '/'
  }

  /**
   * 获取详情页路径
   * @param id 微应用 ID
   */
  const getDetailPath = (id: string | number | undefined): string => {
    if (id === undefined) {
      return getHomePath()
    }
    return isIframeRoute.value ? `${iframePrefix}/microApp/${id}` : `/microApp/${id}`
  }

  /**
   * 获取路由路径（通用）
   * @param basePath 基础路径，如 '/' 或 '/microApp/123'
   */
  const getRoutePath = (basePath: string): string => {
    if (isIframeRoute.value) {
      return basePath === '/' ? `${iframePrefix}/` : `${iframePrefix}${basePath}`
    }
    return basePath
  }

  return {
    isIframeRoute,
    getHomePath,
    getDetailPath,
    getRoutePath,
  }
}
