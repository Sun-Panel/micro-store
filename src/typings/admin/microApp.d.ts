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
    permissionLevel: number
    status: number
    screenshots: string
    createTime: string
    updateTime: string
    langList?: Array<{
      lang: string
      appName: string
      appDesc: string
    }>
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
}
