import { post } from '@/utils/request'

// ==================== 开发者端 API ====================

// 上传版本包
export function uploadVersionPackage<T>(file: File) {
  const formData = new FormData()
  formData.append('file', file)
  return post<T>({
    url: '/developer/version/upload',
    data: formData,
    headers: { 'Content-Type': 'multipart/form-data' },
  })
}

// 获取版本列表
export function getVersionList<T>(data: MicroApp.GetVersionListRequest) {
  return post<T>({
    url: '/developer/version/getList',
    data,
  })
}

// 获取版本详情
export function getVersionInfo<T>(id: number) {
  return post<T>({
    url: '/developer/version/getInfo',
    data: { id },
  })
}

// 创建版本
export function createVersion<T>(data: MicroApp.CreateVersionRequest) {
  return post<T>({
    url: '/developer/version/create',
    data,
  })
}

// 更新版本信息
export function updateVersion<T>(data: MicroApp.UpdateVersionRequest) {
  return post<T>({
    url: '/developer/version/update',
    data,
  })
}

// 提交审核
export function submitReview<T>(data: MicroApp.SubmitReviewRequest) {
  return post<T>({
    url: '/developer/version/submitReview',
    data,
  })
}

// 撤销审核
export function cancelReview<T>(data: MicroApp.CancelReviewRequest) {
  return post<T>({
    url: '/developer/version/cancelReview',
    data,
  })
}

// 删除版本
export function deleteVersion<T>(ids: number[]) {
  return post<T>({
    url: '/developer/version/delete',
    data: { ids },
  })
}

// ==================== 管理员端 API ====================

// 管理员获取版本列表
export function adminGetVersionList<T>(data: MicroApp.GetVersionListRequest & Common.ListRequest) {
  return post<T>({
    url: '/admin/version/getList',
    data,
  })
}

// 管理员获取待审核版本列表
export function adminGetPendingVersionList<T>(data: Common.ListRequest) {
  return post<T>({
    url: '/admin/version/getPendingList',
    data,
  })
}

// 管理员审核版本
export function adminReviewVersion<T>(data: MicroApp.ReviewVersionRequest) {
  return post<T>({
    url: '/admin/version/review',
    data,
  })
}

// 管理员下架版本
export function adminOfflineVersion<T>(data: MicroApp.OfflineVersionRequest) {
  return post<T>({
    url: '/admin/version/offline',
    data,
  })
}
