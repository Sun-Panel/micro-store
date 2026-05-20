export interface AppInfo {
  microAppId: string
  author: string
  icon: string
  iconBase64?: string // 图标的base64编码，用于临时预览，仅前端使用
  appJson: Record<string, any>
  version: string
  apiVersion: string
  src: string
  dev: boolean
  debug: boolean
  runPath: string
  // componentInfo: components | null
  // permissionInfo: Permissions
  runPathUrl: string
  appInfo: {
    [key: string]: AppInfo
  }
}

export interface State {
  installedApps: Record<string, AppInfo>

}

export function defaultSetting(): State {
  return { installedApps: {} }
}
