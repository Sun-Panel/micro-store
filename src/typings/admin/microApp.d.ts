declare namespace MicroApp {
  // 单个语言的信息
  interface LangInfo {
    appName: string
    appDesc: string
  }

  // 微应用信息
  interface MicroAppInfo {
    id: number
    microAppId: string
    appName: string
    appIcon: string
    appDesc: string
    remark?: string
    categoryId: number
    chargeType: number
    price: number
    authorId: number
    authorName?: string
    permissionLevel: number
    status: number
    screenshots: string
    createTime: string
    updateTime: string
    // 审核相关字段
    reviewStatus: number // 审核状态：0-无审核 1-审核中 2-已通过 3-已拒绝
    reviewId?: number // 当前审核记录ID
    reviewTime?: string // 最后审核时间
    langList?: Array<{
      lang: string
      appName: string
      appDesc: string
    }>
  }

  // 微应用审核信息
  interface MicroAppReviewInfo {
    id: number
    appId: number
    appName: string
    appIcon: string
    appDesc: string
    categoryId: number
    chargeType: number
    price: number
    screenshots: string
    langMap?: Record<string, LangInfo>
    remark?: string
    status: number // 审核状态：0-待审核 1-已通过 2-已拒绝
    reviewerId?: number
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

  // 获取列表请求
  interface GetListRequest {
    page: number
    limit: number
    categoryId?: number
    status?: number
    keyWord?: string
    authorId?: number
  }

  // 创建请求
  interface CreateRequest {
    microAppId?: string
    appName: string
    appIcon: string
    appDesc?: string
    remark?: string
    categoryId: number
    chargeType?: number
    price?: number
    authorId: number
    screenshots?: string
    langMap?: Record<string, LangInfo>
  }

  // 更新请求
  interface UpdateRequest {
    id: number
    appName: string
    appIcon: string
    appDesc?: string
    remark?: string
    categoryId: number
    chargeType?: number
    price?: number
    screenshots?: string
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

  // 更新语言请求
  interface UpdateLangRequest {
    id: number
    langMap: Record<string, LangInfo>
  }

  // 撤销审核请求（微应用主信息）
  interface CancelAppReviewRequest {
    id: number
  }

  // 撤销审核（微应用主信息）
  interface CancelAppRequest {
    id: number
  }

  // 管理员审核微应用主信息请求
  interface ReviewAppRequest {
    reviewId: number
    status: number // 1-通过 2-拒绝
    reviewNote?: string
  }

  // 获取审核历史请求
  interface GetReviewHistoryRequest {
    appId: number
    page?: number
    limit?: number
  }

  // 获取待审核列表请求
  interface GetPendingReviewListRequest {
    page?: number
    limit?: number
  }
}
