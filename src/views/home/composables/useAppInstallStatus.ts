import { useLocalAppStore } from '@/store'

/**
 * 应用状态信息接口
 */
export interface AppStatusInfo {
  /** 按钮文本：安装/已安装/升级 */
  text: string
  /** 按钮类型：primary/default/warning */
  type: 'primary' | 'default' | 'warning' | 'success'
  /** 是否禁用 */
  disabled: boolean
  /** 是否已安装 */
  isInstalled: boolean
  /** 是否有更新 */
  hasUpdate: boolean
}

/**
 * 获取微应用的按钮状态
 * @param microAppId - 微应用 ID
 * @param latestVersion - 最新版本号（可选）
 * @returns 应用状态信息
 */
export function getAppButtonStatus(microAppId: string | undefined, latestVersion?: string): AppStatusInfo {
  const localAppStore = useLocalAppStore()
  
  if (!microAppId) {
    return { text: '安装', type: 'primary', disabled: false, isInstalled: false, hasUpdate: false }
  }
  
  const isInstalled = localAppStore.isInstalled(microAppId)
  if (!isInstalled) {
    return { text: '安装', type: 'primary', disabled: false, isInstalled: false, hasUpdate: false }
  }
  
  // 已安装，检查是否有更新
  if (latestVersion && localAppStore.hasUpdate(microAppId, latestVersion)) {
    return { text: '升级', type: 'warning', disabled: false, isInstalled: true, hasUpdate: true }
  }
  
  // 已安装且是最新版本
  return { text: '已安装', type: 'default', disabled: true, isInstalled: true, hasUpdate: false }
}

/**
 * Vue Composable：获取应用安装状态相关的工具函数
 */
export function useAppInstallStatus() {
  return {
    getAppButtonStatus,
  }
}
