import { post } from '@/utils/request'
import { useAuthStore } from '@/store'

/** 上传块预览图（单张 ≤512K），返回图片地址 */
export async function uploadPreviewImage(file: File): Promise<string> {
  const authStore = useAuthStore()
  const formData = new FormData()
  formData.append('imgfile', file)

  const response = await fetch('/api/admin/customCode/uploadPreviewImage', {
    method: 'POST',
    body: formData,
    headers: new Headers({ token: authStore.token || '' }),
  })
  const res = await response.json()

  if (res.code === 0 && res.data?.imageUrl)
    return res.data.imageUrl

  throw new Error(res.msg || '上传失败')
}

export function getMyList<T>(param: Common.ListRequest) {
  return post<T>({
    url: '/admin/customCode/getMyList',
    data: param,
  })
}

export function getInfo<T>(id: number) {
  return post<T>({
    url: '/admin/customCode/getInfo',
    data: { id },
  })
}

/** 获取当前开发者标识（用于编辑页唯一标识前缀展示） */
export function getAuthorKeyPrefix<T>() {
  return post<T>({
    url: '/admin/customCode/getAuthorKeyPrefix',
    data: {},
  })
}

export function edit<T>(info: CustomCode.EditReq) {
  return post<T>({
    url: '/admin/customCode/edit',
    data: info,
  })
}

export function withdraw<T>(id: number) {
  return post<T>({
    url: '/admin/customCode/withdraw',
    data: { id },
  })
}

export function offline<T>(id: number) {
  return post<T>({
    url: '/admin/customCode/offline',
    data: { id },
  })
}

export function deleteById<T>(id: number) {
  return post<T>({
    url: '/admin/customCode/delete',
    data: { id },
  })
}

export function forceOffline<T>(id: number, reason: string) {
  return post<T>({
    url: '/admin/customCode/forceOffline',
    data: { id, reason },
  })
}
