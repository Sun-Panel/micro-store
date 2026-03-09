import { post } from '@/utils/request'

export function getList<T>(req: Admin.RedeemCode.GetListRequest) {
  return post<T>({
    url: '/admin/redeemCode/getList',
    data: req,
  })
}

export function create<T>(data: RedeemCode.CreateReq) {
  return post<T>({
    url: '/admin/redeemCode/create',
    data,
  })
}

export function setInvalid<T>(codes: string[]) {
  return post<T>({
    url: '/admin/redeemCode/setInvalid',
    data: { codes },
  })
}
