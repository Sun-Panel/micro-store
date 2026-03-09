import { post } from '@/utils/request'

export function deletes<T>(ids: number[]) {
  return post<T>({
    url: '/admin/goodsManage/deletes',
    data: { ids },
  })
}

export function getList<T>(param: Common.ListRequest) {
  return post<T>({
    url: '/admin/goodsManage/getList',
    data: param,
  })
}

export function createSnapshot<T>(goodId: number) {
  return post<T>({
    url: '/admin/goodsManage/createSnapshot',
    data: { id: goodId },
  })
}

export function add<T>(info: Admin.GoodsManage.GoodsInfo) {
  return post<T>({
    url: '/admin/goodsManage/add',
    data: info,
  })
}

export function update<T>(info: Admin.GoodsManage.GoodsInfo) {
  return post<T>({
    url: '/admin/goodsManage/update',
    data: info,
  })
}

export function updateSale<T>(info: Admin.GoodsManage.UpdateSaleReq) {
  return post<T>({
    url: '/admin/goodsManage/updateSale',
    data: info,
  })
}
