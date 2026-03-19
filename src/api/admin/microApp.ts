import { post } from '@/utils/request'

// ==================== 管理员专用接口 ====================

// 获取微应用列表（管理员专用）
export function getList<T>(data: MicroApp.GetListRequest) {
  return post<T>({
    url: '/admin/microApp/getList',
    data,
  })
}

// 获取微应用详情（管理员专用）
export function getInfo<T>(id: number) {
  return post<T>({
    url: '/admin/microApp/getInfo',
    data: { id },
  })
}

// 删除微应用（管理员专用）
export function deletes<T>(ids: number[]) {
  return post<T>({
    url: '/admin/microApp/deletes',
    data: { ids },
  })
}

// 更新状态（管理员专用）
export function updateStatus<T>(data: MicroApp.UpdateStatusRequest) {
  return post<T>({
    url: '/admin/microApp/updateStatus',
    data,
  })
}

// 下架微应用（管理员专用）
export function offline<T>(data: MicroApp.OfflineRequest) {
  return post<T>({
    url: '/admin/microApp/offline',
    data,
  })
}
