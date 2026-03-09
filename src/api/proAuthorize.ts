import { post } from '@/utils/request'

export function getAuthorizeHistoryRecord<T>() {
  return post<T>({
    url: '/proAuthorize/getAuthorizeHistoryRecord',
  })
}

export function getAuthorize<T>() {
  return post<T>({
    url: '/proAuthorize/getAuthorize',
  })
}
