import { post } from '@/utils/request'

export function deleteByPageName<T>(pageName: string) {
  return post<T>({
    url: '/admin/htmlPageManage/delete',
    data: { pageName },
  })
}

export function getInfo<T>(pageName: string) {
  return post<T>({
    url: '/admin/htmlPageManage/getInfo',
    data: { pageName },
  })
}

export function getList<T>(param: Common.ListRequest) {
  return post<T>({
    url: '/admin/htmlPageManage/getList',
    data: param,
  })
}

export function edit<T>(info: HtmlPage.EditReq) {
  return post<T>({
    url: '/admin/htmlPageManage/edit',
    data: info,
  })
}
