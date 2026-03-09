import { post } from '@/utils/request'

// 添加keys
export function addKeys<T>(note?: string, keys?: Admin.ApiKeys.ApiKey[]) {
  return post<T>({
    url: '/admin/apiKeys/addKeys',
    data: {
      note,
      keys,
    },
  })
}

// 获取列表
export function getList<T>() {
  return post<T>({
    url: '/admin/apiKeys/getList',
  })
}

// 删除
export function deletes<T>(ids: number[]) {
  return post<T>({
    url: '/admin/apiKeys/deletes',
    data: { ids },
  })
}

// 修改
export function update<T>(info: Admin.ApiKeys.ApiKey) {
  return post<T>({
    url: '/admin/apiKeys/update',
    data: info,
  })
}
