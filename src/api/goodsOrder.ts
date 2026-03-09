import { post } from '@/utils/request'

// 商品

export function create<T>(data: GoodsOrder.CreateRequest) {
  return post<T>({
    url: '/goodsOrder/create',
    data,
  })
}

export function orderDetail<T>(orderNo: string) {
  return post<T>({
    url: '/goodsOrder/details',
    data: { orderNo },
  })
}

export function pay<T>(goodsOrderNo: string, payPlatform: number) {
  return post<T>({
    url: '/goodsOrder/pay',
    data: { goodsOrderNo, payPlatform },
  })
}

export function getOrderList<T>() {
  return post<T>({
    url: '/goodsOrder/getOrderList',
  })
}

export function getOrderInfo<T>(orderNo: string) {
  return post<T>({
    url: '/goodsOrder/getOrderInfo',
    data: { orderNo },
  })
}

export function queryPayStatus<T>(orderNo: string) {
  return post<T>({
    url: '/goodsOrder/queryPayStatus',
    data: { orderNo },
  })
}
