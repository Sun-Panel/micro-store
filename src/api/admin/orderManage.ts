import { post } from '@/utils/request'

// 获取列表
export function getList<T>(data: Common.ListRequest) {
  return post<T>({
    url: '/admin/orderManage/getList',
    data,
  })
}

export function updateAdminNoteByOrderNo<T>(orderNo: string, adminNote: string) {
  return post<T>({
    url: '/admin/orderManage/updateAdminNoteByOrderNo',
    data: {
      orderNo,
      adminNote,
    },
  })
}

export function orderManageUpdateStatusByOrderNo<T>(orderNo: string, status: number) {
  return post<T>({
    url: '/admin/orderManage/orderManageUpdateStatusByOrderNo',
    data: {
      orderNo,
      status,
    },
  })
}

export function getOrderInfo<T>(orderNo: string) {
  return post<T>({
    url: '/admin/orderManage/getOrderInfo',
    data: {
      orderNo,
    },
  })
}
