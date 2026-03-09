import { post } from '@/utils/request'

export function authLogin<T>(appid: string, response_type: string, redirect_uri: string) {
  return post<T>({
    url: '/thirdApp/oAuth2/authLogin',
    data: { appid, response_type, redirect_uri },
  })
}
