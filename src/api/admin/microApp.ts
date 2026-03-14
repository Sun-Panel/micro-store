import { post } from '@/utils/request'

// 获取微应用列表
export function getList<T>(data: MicroApp.GetListRequest) {
  return post<T>({
    url: '/admin/microApp/getList',
    data,
  })
}

// 获取微应用详情
export function getInfo<T>(id: number) {
  return post<T>({
    url: '/admin/microApp/getInfo',
    data: { id },
  })
}

// 创建微应用
export function create<T>(data: MicroApp.CreateRequest) {
  return post<T>({
    url: '/admin/microApp/create',
    data,
  })
}

// 更新微应用
export function update<T>(data: MicroApp.UpdateRequest) {
  return post<T>({
    url: '/admin/microApp/update',
    data,
  })
}

// 删除微应用
export function deletes<T>(ids: number[]) {
  return post<T>({
    url: '/admin/microApp/deletes',
    data: { ids },
  })
}

// 更新状态
export function updateStatus<T>(data: MicroApp.UpdateStatusRequest) {
  return post<T>({
    url: '/admin/microApp/updateStatus',
    data,
  })
}

// 更新语言
export function updateLang<T>(data: MicroApp.UpdateLangRequest) {
  return post<T>({
    url: '/admin/microApp/updateLang',
    data,
  })
}
