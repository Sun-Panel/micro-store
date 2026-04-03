import { post } from '@/utils/request'

// ==================== 开发者专用接口（版本管理） ====================

// 上传版本包
export function uploadVersionPackage<T>(file: File, appRecordId: string) {
  const formData = new FormData()
  formData.append('file', file)
  formData.append('appRecordId', appRecordId)
  return post<T>({
    url: '/admin/developer/version/upload',
    data: formData,
    headers: { 'Content-Type': 'multipart/form-data' },
  })
}

// 获取版本列表
export function getVersionList<T>(data: MicroApp.GetVersionListRequest) {
  return post<T>({
    url: '/admin/developer/version/getList',
    data,
  })
}

// 获取版本详情
export function getVersionInfo<T>(id: number) {
  return post<T>({
    url: '/admin/developer/version/getInfo',
    data: { id },
  })
}

// 创建版本
export function createVersion<T>(data: MicroApp.CreateVersionRequest) {
  return post<T>({
    url: '/admin/developer/version/create',
    data,
  })
}

// 更新版本信息
export function updateVersion<T>(data: MicroApp.UpdateVersionRequest) {
  return post<T>({
    url: '/admin/developer/version/update',
    data,
  })
}

// 提交审核
export function submitReview<T>(data: MicroApp.SubmitReviewRequest) {
  return post<T>({
    url: '/admin/developer/version/submitReview',
    data,
  })
}

// 撤销版本审核
export function cancelReview<T>(data: MicroApp.CancelReviewRequest) {
  return post<T>({
    url: '/admin/developer/version/cancelReview',
    data,
  })
}

// 删除版本
export function deleteVersion<T>(ids: number[]) {
  return post<T>({
    url: '/admin/developer/version/delete',
    data: { ids },
  })
}

// 下架版本
export function offlineVersion<T>(data: { id: number, type: number, reason?: string }) {
  return post<T>({
    url: '/admin/developer/version/offline',
    data,
  })
}

export function getLatestOnlineByAppModelId<T>(id: number) {
  return post<T>({
    url: '/admin/review/microApp/version/getLatestOnlineByAppModelId',
    data: { id },
  })
}

export function getDownloadUrl<T>(versionId: number) {
  return post<T>({
    url: `/admin/microApp/download/getUrl/${versionId}`,
  })
}
