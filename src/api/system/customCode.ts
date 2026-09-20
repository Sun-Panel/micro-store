import { post } from '@/utils/request'

export function getList<T>(param: CustomCode.ListRequest) {
  return post<T>({
    url: '/customCode/getList',
    data: param,
  })
}

export function get<T>(id: number) {
  return post<T>({
    url: '/customCode/get',
    data: { id },
  })
}
