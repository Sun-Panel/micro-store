import { post } from '@/utils/request'

// ==================== 前台/公开 API ====================

// 获取微应用详情（公开接口）
export function getInfo<T>(id: number) {
  return post<T>({
    url: '/admin/microApp/getInfo',
    data: { id },
  })
}

// 获取微应用版本列表（只返回审核通过的）
export function getVersionList<T>(data: MicroApp.GetVersionListRequest) {
  return post<T>({
    url: '/developer/version/getList',
    data,
  })
}
