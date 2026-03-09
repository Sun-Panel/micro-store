import { post } from '@/utils/request'

// 添加keys
export function getAll<T>() {
  return post<T>({
    url: '/admin/clientSetting/createOnlineCache/getAll',
  })
}

export function setAll<T>(data: any) {
  return post<T>({
    url: '/admin/clientSetting/createOnlineCache/setAll',
    data,
  })
}
