import { post } from '@/utils/request'

// ==================== 开发者专用接口（微应用） ====================

// 创建微应用
export function create<T>(data: MicroApp.CreateRequest) {
  return post<T>({
    url: '/admin/developer/microApp/create',
    data,
  })
}

// 更新微应用
export function update<T>(data: MicroApp.UpdateRequest) {
  return post<T>({
    url: '/admin/developer/microApp/update',
    data,
  })
}

// 更新语言
export function updateLang<T>(data: MicroApp.UpdateLangRequest) {
  return post<T>({
    url: '/admin/developer/microApp/updateLang',
    data,
  })
}

// 撤销审核
export function cancelAppReview<T>(data: MicroApp.CancelAppReviewRequest) {
  return post<T>({
    url: '/admin/developer/microApp/cancelReview',
    data,
  })
}

// 获取审核历史
export function getReviewHistory<T>(data: MicroApp.GetReviewHistoryRequest) {
  return post<T>({
    url: '/admin/developer/microApp/getReviewHistory',
    data,
  })
}

// 获取开发者的微应用列表
export function getList<T>(data: MicroApp.GetListRequest) {
  return post<T>({
    url: '/admin/developer/microApp/list',
    data,
  })
}

// 获取开发者的微应用详情
export function getInfo<T>(id: number) {
  return post<T>({
    url: '/admin/developer/microApp/info',
    data: { id },
  })
}

// 删除微应用
export function deletes<T>(ids: number[]) {
  return post<T>({
    url: '/admin/microApp/deletes',
    data: { ids },
  })
}

// 下架微应用
export function offline<T>(data: MicroApp.OfflineRequest) {
  return post<T>({
    url: '/admin/microApp/offline',
    data,
  })
}

// ==================== 开发者专用接口（版本管理） ====================

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

// 撤销版本审核
export function cancelVersionReview<T>(data: MicroApp.CancelReviewRequest) {
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
