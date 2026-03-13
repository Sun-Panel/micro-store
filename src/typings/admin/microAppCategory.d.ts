declare namespace MicroAppCategory {
  // 分类信息
  interface CategoryInfo {
    id: number
    name: string
    icon: string
    sort: number
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

  // 创建请求
  interface CreateRequest {
    name: string
    icon?: string
    sort?: number
    status?: number
  }

  // 更新请求
  interface UpdateRequest {
    id: number
    name: string
    icon?: string
    sort?: number
    status?: number
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

  // 获取详情请求
  interface GetInfoRequest {
    id: number
  }
}
