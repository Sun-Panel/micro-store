import { post } from '@/utils/request'

export function get<T>(mdPageName: string) {
  return post<T>({
    url: '/mdPage/get',
    data: { mdPageName },
  })
}
