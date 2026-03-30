declare namespace MicroApp {
  // 微应用基础信息（公共字段）
  interface BaseInfo {
    microAppId: string
    appName: string // 应用名称（默认语言）
    appIcon: string // 应用图标URL
    appDesc: string // 应用简介（默认语言）
    remark?: string // 应用备注
    categoryId: number // 应用分类ID
    chargeType: number // 收费方式：0-免费 1-积分 2-订阅PRO免费
    points: number // 价格（积分数值）
    screenshots?: string // 图集（多个图片URL用逗号分隔）
  }

  // 单个语言的信息
  interface LangInfo {
    appName: string
    appDesc: string
  }

  // 版本信息
  interface VersionInfo {
    id: number
    appRecordId: number
    microAppId: string
    version: string
    versionCode: number
    packageUrl: string
    packageHash: string
    versionDesc: string
    config?: VersionConfig
    status: number // 审核状态：-1-草稿 0-待审核 1-已通过 2-已拒绝 3-已下架
    reviewTime?: string
    reviewerId?: number
    reviewNote?: string
    offlineType: number // 下架类型：0-正常 1-作者下架 2-平台下架
    offlineReason?: string
    appName?: string
    createdAt?: string
    updatedAt?: string
    microApp?: MicroApp.Info
    createTime?: string
    updateTime?: string
  }

  // 版本配置信息
  interface VersionConfig {
    appJsonVersion?: string
    microAppId?: string
    version?: string
    apiVersion?: string
    author?: string
    entry?: string
    icon?: string
    debug?: boolean
    components?: Record<string, any>
    permissions?: string[]
    dataNodes?: Record<string, any>
    networkDomains?: string[]
    appInfo?: Record<string, AppInfo>
  }

  // 开发者信息
  interface DeveloperInfo {
    id: number
    name: string
    avatar?: string
  }

  // 微应用信息
  interface Info extends BaseInfo {
    id?: number
    microAppId: string
    developerId?: number
    status: number
    offlineType?: number // 下架类型：0-正常 1-作者下架 2-平台下架 3-首次创建
    offlineReason?: string // 下架原因
    createdAt?: string // 创建时间
    updatedAt?: string // 更新时间
    deletedAt?: string // 软删除时间
    // 审核相关字段（从审核表关联）
    reviewStatus?: number // 审核状态：0-无审核 1-审核中 2-已通过 3-已拒绝
    reviewId?: number // 当前审核记录ID
    reviewTime?: string // 最后审核时间
    // 多语言信息
    langList?: Array<{
      lang: string
      appName: string
      appDesc: string
    }>
    defaultLangInfo?: LangInfo // 默认语言信息
    developer?: DeveloperInfo // 开发者详细信息
    createTime?: string // 创建时间
  }

  // 带语言信息的微应用信息（用于前端展示）
  interface InfoWithLang extends BaseInfo {
    id?: number
    microAppId: string
    developerId: number
    developerName?: string // 开发者名称（用于前端展示）
    developerAvatar?: string // 开发者头像（用于前端展示）
    permissionLevel?: number
    status: number
    offlineType?: number
    offlineReason?: string
    createdAt?: string
    updatedAt?: string
    deletedAt?: string
    lang: string // 当前语言
    // 审核相关字段
    reviewStatus?: number
    reviewId?: number
    reviewTime?: string
  }

  // 微应用审核信息
  interface MicroAppReviewInfo extends BaseInfo {
    id: number
    appRecordId: number
    langMap?: Record<string, LangInfo>
    status: number // 审核状态：0-待审核 1-已通过 2-已拒绝
    reviewerId?: number
    reviewerName?: string
    reviewNote?: string
    reviewTime?: string
    createTime: string
    updateTime: string
  }

  // 审核状态枚举
  enum ReviewStatus {
    /** 无审核 */
    NO_REVIEW = 0,
    /** 审核中 */
    REVIEWING = 1,
    /** 已通过 */
    APPROVED = 2,
    /** 已拒绝 */
    REJECTED = 3,
  }

  // 应用状态枚举
  enum AppStatus {
    /** 下架 */
    OFFLINE = 0,
    /** 上架 */
    ONLINE = 1,
  }

  // 收费方式枚举
  enum ChargeType {
    /** 免费 */
    FREE = 0,
    /** 积分 */
    POINTS = 1,
    /** 订阅PRO免费 */
    PRO_FREE = 2,
  }

  // 下架类型枚举
  enum OfflineType {
    /** 正常 */
    NORMAL = 0,
    /** 作者下架 */
    AUTHOR = 1,
    /** 平台下架 */
    PLATFORM = 2,
    /** 首次创建 */
    FIRST_CREATE = 3,
  }

  // 获取微应用详情请求
  interface GetInfoRequest {
    id: number
  }

  // 获取微应用版本列表请求
  interface GetVersionListRequest {
    appRecordId: number
    page?: number
    limit?: number
  }

  // 创建版本请求
  interface CreateVersionRequest {
    // appRecordId: number
    // version: string
    // versionCode: number
    // packageUrl: string
    // packageHash?: string
    versionDesc?: string
    // config?: MicroAppVersionConfig
    uploadCacheId?: string
  }

  // 更新版本请求
  interface UpdateVersionRequest {
    id: number
    version: string
    type: 'release' | 'beta' | 'alpha' | 'rc' | 'dev'
    releaseTime: string
    description: string
    downloadURL: string
    pageUrl: string
  }

  // 提交版本审核请求
  interface SubmitReviewRequest {
    versionId: number
  }

  // 撤销版本审核请求
  interface CancelReviewRequest {
    versionId: number
  }

  // 版本下架请求
  interface OfflineVersionRequest {
    id: number
    type: number
    reason?: string
  }

  // 获取列表请求
  interface GetListRequest {
    page: number
    limit: number
    categoryId?: number
    status?: number
    keyWord?: string
    authorId?: number // 开发者ID筛选
    lang?: string // 可选，用于多语言查询
    includeDeveloper?: boolean // 是否包含开发者信息
  }

  // 获取列表响应
  interface GetListResponse {
    list: Info[]
    total: number
    page: number
    limit: number
  }

  // 获取带语言信息的列表响应
  interface GetListWithLangResponse {
    list: InfoWithLang[]
    total: number
    page: number
    limit: number
  }

  // 创建请求
  interface CreateRequest extends Partial<BaseInfo> {
    microAppId?: string
    langMap?: Record<string, LangInfo> // 多语言信息
  }

  // 更新请求
  interface UpdateRequest extends Partial<BaseInfo> {
    id: number
    langMap?: Record<string, LangInfo>
  }

  // 删除请求
  interface DeletesRequest {
    ids: number[]
  }

  // 更新状态请求
  interface UpdateStatusRequest {
    id: number
    status: number
  }

  // 下架请求
  interface OfflineRequest {
    id: number
    offlineType: number // 下架类型：1-作者下架 2-平台下架
    reason?: string // 下架原因
  }

  // 更新语言请求
  interface UpdateLangRequest {
    id: number
    langMap: Record<string, LangInfo>
  }

  // 撤销审核请求（微应用主信息）
  interface CancelAppReviewRequest {
    reviewId: number
  }

  // 提交审核请求（微应用主信息）
  interface SubmitAppReviewRequest {
    reviewId: number
  }

  // 撤销审核（微应用主信息）
  interface CancelAppRequest {
    id: number
  }

  // 获取微应用详情（带审核信息）响应
  interface GetInfoWithReviewResponse {
    microApp: Info
    microAppReview: MicroAppReviewInfo | null
  }

  // 管理员审核微应用主信息请求
  interface ReviewAppRequest {
    reviewId: number
    status: number // 1-通过 2-拒绝
    reviewNote?: string
  }

  // 获取审核历史请求
  interface GetReviewHistoryRequest {
    appRecordId: number
    page?: number
    limit?: number
  }

  // 获取待审核列表请求
  interface GetPendingReviewListRequest {
    page?: number
    limit?: number
  }

  // 微应用版本上传响应
  interface MicroAppVersionUploadResp {
    url: string // 文件访问 URL
    hash: string // 文件 MD5 校验值
    config: VersionConfig // 解析出的配置文件
    fileName: string // 文件名
    fileSize: number // 文件大小
    folderName: string // 缓存文件夹名（不含路径）
    iconURL: string // 图标访问 URL
    uploadCacheId: string // 上传缓存ID
  }
}
