import { post } from '@/utils/request'

// 获取待审核列表（审核员专用）
export function getPendingList<T>(data: MicroApp.GetPendingReviewListRequest) {
  return post<T>({
    url: '/admin/review/microApp/getPendingList',
    data,
  })
}

// 获取审核详情（审核员专用）
export function getReviewInfo<T>(reviewId: number) {
  return post<T>({
    url: '/admin/review/getReviewInfo',
    data: { reviewId },
  })
}

// 获取审核详情（审核员专用）
export function getMicroAppInfo<T>(id: number) {
  return post<T>({
    url: '/admin/review/getMicroAppInfo',
    data: { id },
  })
}

// 审核微应用主信息（审核员专用）
export function approve<T>(data: MicroApp.ReviewAppRequest) {
  return post<T>({
    url: '/admin/review/approve',
    data,
  })
}
