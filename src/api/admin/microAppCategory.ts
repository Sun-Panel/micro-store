import { post } from '@/utils/request'

// 获取分类列表
export function getList<T>(data: MicroAppCategory.GetListRequest) {
  return post<T>({
    url: '/admin/microAppCategory/getList',
    data,
  })
}

// 获取分类详情
export function getInfo<T>(data: MicroAppCategory.GetInfoRequest) {
  return post<T>({
    url: '/admin/microAppCategory/getInfo',
    data,
  })
}

// 创建分类
export function create<T>(data: MicroAppCategory.CreateRequest) {
  return post<T>({
    url: '/admin/microAppCategory/create',
    data,
  })
}

// 更新分类
export function update<T>(data: MicroAppCategory.UpdateRequest) {
  return post<T>({
    url: '/admin/microAppCategory/update',
    data,
  })
}

// 删除分类
export function deletes<T>(data: MicroAppCategory.DeletesRequest) {
  return post<T>({
    url: '/admin/microAppCategory/deletes',
    data,
  })
}

// 更新状态
export function updateStatus<T>(data: MicroAppCategory.UpdateStatusRequest) {
  return post<T>({
    url: '/admin/microAppCategory/updateStatus',
    data,
  })
}

// 获取启用的分类列表
export function getEnabledList<T>() {
  return post<T>({
    url: '/admin/microAppCategory/getEnabledList',
  })
}
