import { post } from '@/utils/request'

// ==================== 前台/公开 API ====================

// 获取微应用列表（公开接口）
export function getList<T>(data: MicroApp.GetListRequest) {
  return post<T>({
    url: '/microApp/getList',
    data,
  })
}

// 获取微应用详情（公开接口 - 新增）
export function getInfo<T>(id: number) {
  return post<T>({
    url: '/microApp/getInfo',
    data: { id },
  })
}

// 获取微应用版本列表（只返回审核通过的）
export function getVersionList<T>(data: MicroApp.GetVersionListRequest) {
  return post<T>({
    url: '/microApp/version/getList',
    data,
  })
}

export function getDownloadUrl<T>(microAppId: string, version?: string) {
  return post<T>({
    url: '/microApp/download/getUrl/',
    data: { microAppId, version },
  })
}
