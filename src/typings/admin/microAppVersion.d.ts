declare namespace MicroApp {
  // 应用多语言信息（对应后端 AppInfo）
  interface AppLangInfo {
    appName: string
    description: string
    networkDescription?: string
  }

  // 版本配置信息（对应后端 MicroAppConfig）
  interface VersionConfig {
    appJsonVersion?: string
    microAppId?: string
    version?: string
    apiVersion?: string
    author?: string
    entry?: string
    icon?: string
    debug?: boolean
    permissions?: string[]
    networkDomains?: string[]
    appInfo?: Record<string, AppLangInfo>
    // 其他字段
    components?: Record<string, any>
    dataNodes?: Record<string, any>
  }

  // 微应用版本信息
  interface VersionInfo {
    id: number
    appId: number
    version: string
    versionCode: number
    packageUrl: string
    packageHash: string
    versionDesc?: string
    config?: VersionConfig
    status: number
    reviewTime?: string
    reviewerId: number
    reviewNote: string
    createTime: string
    // 兼容旧字段（用于显示）
    iconURL?: string
    author?: string
    appName?: string
    appIcon?: string
  }

  // 版本状态枚举
  enum VersionStatus {
    /** 草稿 */
    DRAFT = -1,
    /** 待审核 */
    PENDING = 0,
    /** 通过 */
    APPROVED = 1,
    /** 拒绝 */
    REJECTED = 2,
  }

  // 创建版本请求
  interface CreateVersionRequest {
    appId: number
    version: string
    versionCode: number
    packageUrl: string
    packageHash: string
    versionDesc?: string
    config?: VersionConfig
  }

  // 更新版本请求
  interface UpdateVersionRequest {
    id: number
    version?: string
    versionCode?: number
    versionDesc?: string
  }

  // 提交审核请求
  interface SubmitReviewRequest {
    versionId: number
  }

  // 撤销审核请求
  interface CancelReviewRequest {
    versionId: number
  }

  // 管理员审核请求
  interface ReviewVersionRequest {
    versionId: number
    status: number
    reviewNote?: string
  }

  // 管理员下架版本请求
  interface OfflineVersionRequest {
    id: number
    type: number       // 下架类型：1-作者下架 2-平台下架
    reason?: string   // 下架原因
  }

  // 获取版本列表请求
  interface GetVersionListRequest {
    appId: number
    page?: number
    limit?: number
    status?: number
  }
}
