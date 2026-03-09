import { post } from '@/utils/request'

export function updateUserExpiredTimeByDay<T>(data: Admin.ProAuthorize.ProAuthorizeUpdateUserExpiredTimeByDayReq) {
  return post<T>({
    url: '/admin/proAuthorize/updateUserExpiredTimeByDay',
    data,
  })
}

export function getUserProAuthorizeList<T>(data: Common.ListRequest) {
  return post<T>({
    url: '/admin/proAuthorize/getUserProAuthorizeList',
    data,
  })
}

export function getAuthorizeHistoryRecordByUserId<T>(userId: number) {
  return post<T>({
    url: '/admin/proAuthorize/getAuthorizeHistoryRecordByUserId',
    data: { userId },
  })
}
