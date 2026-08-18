import type { AppInfo, State } from './helper'
import { defineStore } from 'pinia'
import { isVersionEqual, isVersionGreater } from '@/utils/functions'
import { defaultSetting } from './helper'

/**
 * 本地已安装微应用状态管理
 * 管理本地已安装的微应用列表，提供安装状态查询和版本比较功能
 */
export const useLocalAppStore = defineStore('local-app-store', {
  state: (): State => defaultSetting(),
  actions: {
    /**
     * 设置已安装的微应用列表
     * @param apps - 已安装的微应用对象，key 为 microAppId
     */
    setInstalledApps(apps: Record<string, AppInfo>) {
      this.installedApps = apps
    },

    /**
     * 检查指定微应用是否已安装
     * @param microAppId - 微应用的唯一标识符
     * @returns 是否已安装
     */
    isInstalled(microAppId: string) {
      microAppId = microAppId.toLowerCase()
      return this.installedApps[microAppId] !== undefined
    },

    /**
     * 检查指定微应用是否为最新版本
     * 使用语义化版本比较，支持从右向左比较版本号
     * @param microAppId - 微应用的唯一标识符
     * @param latestVersion - 最新版本号
     * @returns 是否为最新版本
     */
    isLatestVersion(microAppId: string, latestVersion: string) {
      microAppId = microAppId.toLowerCase()
      const app = this.installedApps[microAppId]
      if (!app?.version)
        return false
      return isVersionEqual(app.version, latestVersion)
    },

    /**
     * 检查指定微应用是否有可用更新
     * @param microAppId - 微应用的唯一标识符
     * @param latestVersion - 最新版本号
     * @returns 是否有可用更新
     */
    hasUpdate(microAppId: string, latestVersion: string) {
      microAppId = microAppId.toLowerCase()
      const app = this.installedApps[microAppId]
      if (!app?.version)
        return false
      return isVersionGreater(latestVersion, app.version)
    },

  },
})
