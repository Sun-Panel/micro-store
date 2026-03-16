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

// 下架微应用
export function offline<T>(data: MicroApp.OfflineRequest) {
  return post<T>({
    url: '/admin/microApp/offline',
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

// ==================== 微应用主信息审核相关 ====================

// 撤销审核（与版本审核的 cancelReview 区分）
export function cancelAppReview<T>(data: MicroApp.CancelAppReviewRequest) {
  return post<T>({
    url: '/admin/microApp/cancelReview',
    data,
  })
}

// 撤销审核
export function cancelReview<T>(data: MicroApp.CancelAppReviewRequest) {
  return post<T>({
    url: '/admin/microApp/cancelReview',
    data,
  })
}

// 获取审核历史
export function getReviewHistory<T>(data: MicroApp.GetReviewHistoryRequest) {
  return post<T>({
    url: '/admin/microApp/getReviewHistory',
    data,
  })
}

// 获取审核详情
export function getReviewInfo<T>(reviewId: number) {
  return post<T>({
    url: '/admin/microApp/getReviewInfo',
    data: { reviewId },
  })
}

// 获取待审核列表（管理员）
export function getPendingReviewList<T>(data: MicroApp.GetPendingReviewListRequest) {
  return post<T>({
    url: '/admin/microApp/getPendingReviewList',
    data,
  })
}

// 审核微应用主信息（管理员）
export function reviewApp<T>(data: MicroApp.ReviewAppRequest) {
  return post<T>({
    url: '/admin/microApp/reviewApp',
    data,
  })
}
