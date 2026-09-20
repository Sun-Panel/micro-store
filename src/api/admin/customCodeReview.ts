import { post } from '@/utils/request'

export function getList<T>(param: Common.ListRequest) {
  return post<T>({
    url: '/admin/customCodeReview/getList',
    data: param,
  })
}

export function getInfo<T>(id: number) {
  return post<T>({
    url: '/admin/customCodeReview/getInfo',
    data: { id },
  })
}

export function approve<T>(id: number, note?: string) {
  return post<T>({
    url: '/admin/customCodeReview/approve',
    data: { id, note },
  })
}

export function reject<T>(id: number, note: string) {
  return post<T>({
    url: '/admin/customCodeReview/reject',
    data: { id, note },
  })
}

export function getHistory<T>(id: number, page: number, limit: number) {
  return post<T>({
    url: '/admin/customCodeReview/getHistory',
    data: { id, page, limit },
  })
}
