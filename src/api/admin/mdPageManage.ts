import { post } from '@/utils/request'

export function deleteByMdPageName<T>(mdPageName: string) {
  return post<T>({
    url: '/admin/mdPage/delete',
    data: { mdPageName },
  })
}

export function getList<T>(param: Common.ListRequest) {
  return post<T>({
    url: '/admin/mdPage/getList',
    data: param,
  })
}

export function edit<T>(info: MdPage.EditReq) {
  return post<T>({
    url: '/admin/mdPage/edit',
    data: info,
  })
}
