import { post } from '@/utils/request'

// 商品

export function goodsGetList<T>() {
  return post<T>({
    url: '/goods/getList',
  })
}

export function getInfo<T>(id: number) {
  return post<T>({
    url: '/goods/getInfo',
    data: { id },
  })
}
