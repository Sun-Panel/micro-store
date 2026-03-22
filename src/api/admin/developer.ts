import { post } from '@/utils/request'

// ==================== 管理员专用接口（开发者用户管理） ====================

// 获取开发者列表
export function getList<T>(data: Developer.GetListRequest) {
  return post<T>({
    url: '/admin/developer/user/getList',
    data,
  })
}

// 获取开发者详情
export function getInfo<T>(id: number) {
  return post<T>({
    url: '/admin/developer/user/getInfo',
    data: { id },
  })
}

// 根据开发者标识获取开发者信息
export function getByDeveloperName<T>(developerName: string) {
  return post<T>({
    url: '/admin/developer/user/getByDeveloperName',
    data: { developerName },
  })
}

// 更新开发者
export function update<T>(data: Developer.UpdateRequest) {
  return post<T>({
    url: '/admin/developer/user/update',
    data,
  })
}

// 删除开发者
export function deletes<T>(ids: number[]) {
  return post<T>({
    url: '/admin/developer/user/deletes',
    data: { ids },
  })
}

// 更新状态
export function updateStatus<T>(data: Developer.UpdateStatusRequest) {
  return post<T>({
    url: '/admin/developer/user/updateStatus',
    data,
  })
}
