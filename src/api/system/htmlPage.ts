import { post } from '@/utils/request'

export function get<T>(pageName: string) {
  return post<T>({
    url: '/htmlPage/get',
    data: { pageName },
  })
}
