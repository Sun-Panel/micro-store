import { post } from '@/utils/request'

// ==================== 开发者专用接口（微应用） ====================

// 获取开发者的微应用列表
export function getList<T>(data: MicroApp.GetListRequest) {
  return post<T>({
    url: '/admin/developer/myMicroApp/list',
    data,
  })
}

// 获取开发者的微应用详情
export function getInfo<T>(id: number) {
  return post<T>({
    url: '/admin/developer/myMicroApp/info',
    data: { id },
  })
}
export function getMicroInfoAndReviewInfoByMicroAppModelId<T>(id: number) {
  return post<T>({
    url: '/admin/developer/myMicroApp/getMicroInfoAndReviewInfoByMicroAppModelId',
    data: { id },
  })
}

// 创建微应用
export function create<T>(data: MicroApp.CreateRequest) {
  return post<T>({
    url: '/admin/developer/myMicroApp/create',
    data,
  })
}

// 更新微应用
export function update<T>(data: MicroApp.UpdateRequest) {
  return post<T>({
    url: '/admin/developer/myMicroApp/update',
    data,
  })
}

// 更新语言
export function updateLang<T>(data: MicroApp.UpdateLangRequest) {
  return post<T>({
    url: '/admin/developer/myMicroApp/updateLang',
    data,
  })
}

// 撤销审核
export function cancelReview<T>(data: MicroApp.CancelAppReviewRequest) {
  return post<T>({
    url: '/admin/developer/myMicroApp/cancelReview',
    data,
  })
}

// 提交审核
export function submitReview<T>(data: MicroApp.SubmitAppReviewRequest) {
  return post<T>({
    url: '/admin/developer/myMicroApp/submitReview',
    data,
  })
}

// 获取审核历史
export function getReviewHistory<T>(data: MicroApp.GetReviewHistoryRequest) {
  return post<T>({
    url: '/admin/developer/myMicroApp/getReviewHistory',
    data,
  })
}

// 删除微应用（共享接口）
export function deletes<T>(ids: number[]) {
  return post<T>({
    url: '/admin/microApp/deletes',
    data: { ids },
  })
}

// 下架微应用（共享接口）
export function offline<T>(data: MicroApp.OfflineRequest) {
  return post<T>({
    url: '/admin/microApp/offline',
    data,
  })
}
