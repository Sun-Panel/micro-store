/**
 * iframe通信事件类型定义
 * 简化版本，定义微应用商店与主面板之间的通信事件
 */

// 微应用商店事件（从主面板发送到微应用商店）
export const MICRO_APP_STORE_EVENTS = {
  // 应用安装事件
  INSTALL_APP: 'microAppStore:installApp',
  // 应用卸载事件
  UNINSTALL_APP: 'microAppStore:uninstallApp',
  // 应用更新事件
  UPDATE_APP: 'microAppStore:updateApp',
  // 获取应用列表
  GET_APP_LIST: 'microAppStore:getAppList',
  // 获取应用详情
  GET_APP_DETAIL: 'microAppStore:getAppDetail',
  // 获取用户信息
  GET_USER_INFO: 'microAppStore:getUserInfo',
  // 登录事件（主面板通知微应用登录状态变化）
  LOGIN: 'microAppStore:login',
  LOGOUT: 'microAppStore:logout',
  // 获取已安装应用列表
  INSTALLED_APPS: 'microAppStore:installedApps',
} as const

// 微应用事件（从微应用商店发送到主面板）
export const MICRO_APP_EVENTS = {
  // 应用安装完成
  APP_INSTALLED: 'microApp:appInstalled',
  // 应用卸载完成
  APP_UNINSTALLED: 'microApp:appUninstalled',
  // 应用更新完成
  APP_UPDATED: 'microApp:appUpdated',
  // 应用列表数据
  APP_LIST_DATA: 'microApp:appListData',
  // 应用详情数据
  APP_DETAIL_DATA: 'microApp:appDetailData',
  // 用户信息数据
  USER_INFO_DATA: 'microApp:userInfoData',
  // 通信就绪事件
  COMMUNICATION_READY: 'microApp:communicationReady',
  // 获取已安装应用列表
  GET_INSTALLED_APPS: 'microApp:getInstalledApps',
  // 触发安装应用
  INSTALL_APP: 'microApp:installApp',
} as const

// 事件类型联合类型
export type MicroAppStoreEvent = typeof MICRO_APP_STORE_EVENTS[keyof typeof MICRO_APP_STORE_EVENTS]
export type MicroAppEvent = typeof MICRO_APP_EVENTS[keyof typeof MICRO_APP_EVENTS]
