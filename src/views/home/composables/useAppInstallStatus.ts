import { useLocalAppStore } from '@/store'
import { extractBaseVersion, isVersionGreaterOrEqual } from '@/utils/functions/version'
import { getHostAppBaseVersion, isIframeMode } from './useGetIframeUrlParam'

/**
 * 应用状态信息接口
 */
export interface AppStatusInfo {
  /** 按钮文本：安装/已安装/升级/不兼容 */
  text: string
  /** 按钮类型：primary/default/warning */
  type: 'primary' | 'default' | 'warning' | 'success'
  /** 是否禁用 */
  disabled: boolean
  /** 是否已安装 */
  isInstalled: boolean
  /** 是否有更新 */
  hasUpdate: boolean
  /** 是否版本不兼容 */
  incompatible: boolean
  /** 不兼容的提示信息 */
  incompatibleMsg?: string
}

/**
 * 版本兼容性检查结果
 */
export interface VersionCompatibilityResult {
  /** 是否兼容 */
  compatible: boolean
  /** 提示信息 */
  message: string
  /** 主应用版本 */
  hostVersion: string
  /** 最低要求版本 */
  requiredVersion: string
}

/**
 * 检查当前主应用版本是否兼容微应用指定的最低版本
 * @param requiredVersion 微应用要求的最低主应用版本（如 "2.0.0"）
 * @returns 兼容性检查结果
 */
export function checkVersionCompatibility(requiredVersion?: string): VersionCompatibilityResult {
  const defaultResult: VersionCompatibilityResult = {
    compatible: true,
    message: '',
    hostVersion: '',
    requiredVersion: '',
  }

  // 非 iframe 模式下不做版本兼容性检查（应用商店独立使用时无需判断）
  if (!isIframeMode()) {
    return defaultResult
  }

  const hostBaseVersion = getHostAppBaseVersion()
  defaultResult.hostVersion = hostBaseVersion

  if (!requiredVersion) {
    return defaultResult
  }

  const hostBase = extractBaseVersion(hostBaseVersion)
  const requiredBase = extractBaseVersion(requiredVersion)

  if (!isVersionGreaterOrEqual(hostBase, requiredBase)) {
    return {
      compatible: false,
      message: `当前主应用版本 ${hostBase}，该微应用需要主应用版本 >= ${requiredBase}`,
      hostVersion: hostBase,
      requiredVersion: requiredBase,
    }
  }

  return {
    compatible: true,
    message: '',
    hostVersion: hostBase,
    requiredVersion: requiredBase,
  }
}

/**
 * 获取微应用的按钮状态
 * @param microAppId - 微应用 ID
 * @param latestVersion - 最新版本号（可选）
 * @param lowVersionOrConfig - 最低主应用版本号字符串或版本配置对象（可选，用于版本兼容性检查）
 * @returns 应用状态信息
 */
export function getAppButtonStatus(
  microAppId: string | undefined,
  latestVersion?: string,
  lowVersionOrConfig?: string | MicroApp.VersionConfig,
): AppStatusInfo {
  const localAppStore = useLocalAppStore()

  const baseStatus: Omit<AppStatusInfo, 'incompatible' | 'incompatibleMsg'> = {
    text: '安装',
    type: 'primary',
    disabled: false,
    isInstalled: false,
    hasUpdate: false,
  }

  if (!microAppId) {
    return { ...baseStatus, incompatible: false }
  }

  const isInstalled = localAppStore.isInstalled(microAppId)

  // 版本兼容性检查（仅未安装时检查）
  if (!isInstalled) {
    // 提取 lowVersion：仅使用 lowVersion 字段（最低主应用版本要求）
    // 不 fallback 到 apiVersion/appJsonVersion，因为那是微应用自身版本，不是主应用版本要求
    let requiredVersion: string | undefined
    if (typeof lowVersionOrConfig === 'string') {
      requiredVersion = lowVersionOrConfig
    }
    else if (lowVersionOrConfig) {
      requiredVersion = lowVersionOrConfig.lowVersion
    }

    // lowVersion 为空时跳过兼容性检查
    const compatResult = checkVersionCompatibility(requiredVersion)
    if (!compatResult.compatible) {
      return {
        ...baseStatus,
        text: '安装',
        type: 'default',
        disabled: true,
        incompatible: true,
        incompatibleMsg: compatResult.message,
      }
    }
    return { ...baseStatus, incompatible: false }
  }

  // 已安装，检查是否有更新
  if (latestVersion && localAppStore.hasUpdate(microAppId, latestVersion)) {
    // 有更新时，检查新版本的兼容性
    let requiredVersion: string | undefined
    if (typeof lowVersionOrConfig === 'string') {
      requiredVersion = lowVersionOrConfig
    }
    else if (lowVersionOrConfig) {
      requiredVersion = lowVersionOrConfig.lowVersion
    }

    // 检查新版本兼容性
    const compatResult = checkVersionCompatibility(requiredVersion)
    if (!compatResult.compatible) {
      return {
        text: '升级',
        type: 'default',
        disabled: true,
        isInstalled: true,
        hasUpdate: true,
        incompatible: true,
        incompatibleMsg: compatResult.message,
      }
    }

    return { text: '升级', type: 'warning', disabled: false, isInstalled: true, hasUpdate: true, incompatible: false }
  }

  // 已安装且是最新版本
  return { text: '已安装', type: 'default', disabled: true, isInstalled: true, hasUpdate: false, incompatible: false }
}

/**
 * Vue Composable：获取应用安装状态相关的工具函数
 */
export function useAppInstallStatus() {
  return {
    getAppButtonStatus,
    checkVersionCompatibility,
  }
}
