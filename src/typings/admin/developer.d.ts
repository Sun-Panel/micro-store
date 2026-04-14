declare namespace Developer {
  // 开发者信息
  interface DeveloperInfo {
    id: number
    userId: number
    developerName: string
    contactMail: string
    paymentName: string
    paymentQrcode: string
    paymentMethod: string
    status: number
    createTime: string
    updateTime: string
  }

  // 获取列表请求
  interface GetListRequest {
    page: number
    limit: number
    status?: number
    keyWord?: string
  }

  // 更新请求
  interface UpdateRequest {
    id: number
    developerName: string
    contactMail?: string
    paymentName?: string
    paymentQrcode?: string
    paymentMethod?: string
    status: number
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

  // 注册请求（前台）
  interface RegisterRequest {
    developerName: string
    contactMail?: string
    paymentName?: string
    paymentQrcode?: string
    paymentMethod?: string
    name: string
  }

  // 更新请求（前台）
  interface UpdateMyInfoRequest {
    developerName: string
    contactMail?: string
    paymentName?: string
    paymentQrcode?: string
    paymentMethod?: string
  }

  // 检查是否是开发者响应
  interface CheckIsDeveloperResponse {
    isDeveloper: boolean
  }
}
