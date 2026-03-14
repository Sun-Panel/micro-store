import { post } from '@/utils/request'

// 获取开发者列表
export function getList<T>(data: Developer.GetListRequest) {
  return post<T>({
    url: '/admin/developer/getList',
    data,
  })
}

// 获取开发者详情
export function getInfo<T>(id: number) {
  return post<T>({
    url: '/admin/developer/getInfo',
    data: { id },
  })
}

// 更新开发者
export function update<T>(data: Developer.UpdateRequest) {
  return post<T>({
    url: '/admin/developer/update',
    data,
  })
}

// 删除开发者
export function deletes<T>(ids: number[]) {
  return post<T>({
    url: '/admin/developer/deletes',
    data: { ids },
  })
}

// 更新状态
export function updateStatus<T>(data: Developer.UpdateStatusRequest) {
  return post<T>({
    url: '/admin/developer/updateStatus',
    data,
  })
}
