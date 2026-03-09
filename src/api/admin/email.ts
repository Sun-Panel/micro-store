import { post } from '@/utils/request'

export function send<T>(data: Email.SendEmailReq) {
  return post<T>({
    url: '/admin/email/send',
    data,
  })
}
