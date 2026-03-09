import { post } from '@/utils/request'

export function getMultiple<T>(data: string[]) {
  return post<T>({
    url: '/systemVariableApi/getMultiple',
    data,
  })
}
