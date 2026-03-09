import { post } from '@/utils/request'

export function getThirdAppInfo<T>(appid: string) {
  return post<T>({
    url: 'thirdApp/thirdApp/getThirdAppInfo',
    data: { appid },
  })
}
