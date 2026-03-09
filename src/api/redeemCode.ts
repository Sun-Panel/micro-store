import { post } from '@/utils/request'

export function getRedeemCodeInfo<T>(code: string) {
  return post<T>({
    url: '/proAuthorize/getRedeemCodeInfo',
    data: { code },
  })
}

export function redeemCodeWriteOff<T>(code: string, vcode: string) {
  return post<T>({
    url: '/proAuthorize/redeemCodeWriteOff',
    data: { code, vcode },
  })
}
