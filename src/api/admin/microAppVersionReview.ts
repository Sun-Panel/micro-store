import { post } from '@/utils/request'

// ==================== 审核员专用接口（应用审核） ====================

// 获取待审核应用列表（审核员专用）
export function getPendingAppList<T>(data: Common.ListRequest) {
  return post<T>({
    url: '/admin/review/microApp/getPendingList',
    data,
  })
}

// 获取审核详情（审核员专用）
export function getAppReviewInfo<T>(reviewId: number) {
  return post<T>({
    url: '/admin/review/getInfo',
    data: { reviewId },
  })
}

// 审核应用（审核员专用）
export function approveApp<T>(data: MicroApp.ReviewAppRequest) {
  return post<T>({
    url: '/admin/review/approve',
    data,
  })
}

// ==================== 审核员专用接口（版本审核） ====================

// 获取待审核版本列表（审核员专用）
export function getPendingVersionList<T>(data: Common.ListRequest) {
  return post<T>({
    url: '/admin/review/microApp/version/getPendingList',
    data,
  })
}

// 兼容旧版本的别名导出
export function getPendingList<T>(data: Common.ListRequest) {
  return getPendingVersionList<T>(data)
}

// 审核版本（审核员专用）
export function reviewVersion<T>(data: MicroApp.ReviewVersionRequest) {
  return post<T>({
    url: '/admin/reviewVersion/review',
    data,
  })
}

// 兼容旧版本的别名导出
export function review<T>(data: MicroApp.ReviewVersionRequest) {
  return reviewVersion<T>(data)
}

// 下架版本（审核员专用）
export function offlineVersion<T>(data: MicroApp.OfflineVersionRequest) {
  return post<T>({
    url: '/admin/reviewVersion/offline',
    data,
  })
}

// 主动触发安全审核
export function triggerSecurityAudit<T>(versionId: number) {
  return post<T>({
    url: '/admin/reviewVersion/triggerSecurityAudit',
    data: { versionId },
  })
}

// 获取版本列表（管理员专用）
export function getVersionList<T>(data: MicroApp.GetVersionListRequest) {
  return post<T>({
    url: '/admin/version/getList',
    data,
  })
}
