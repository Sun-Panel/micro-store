import { post } from '@/utils/request'

export function send<T>(data: Message.SendReq) {
  return post<T>({
    url: '/message/send',
    data,
  })
}

// 获取列表
export function getTemplateList<T>() {
  return post<T>({
    url: '/message/template/getList',
  })
}

export function getSendList<T>() {
  return post<T>({
    url: '/message/getSendList',
  })
}

export function getReceiveList<T>() {
  return post<T>({
    url: '/message/getReceiveList',
  })
}

export function getMessageInfo<T>(messageId: number) {
  return post<T>({
    url: '/message/getMessageInfo',
    data: { messageId },
  })
}

export function updateReadStatus<T>(messageId: number, readStatus: boolean) {
  return post<T>({
    url: '/message/updateReadStatus',
    data: { messageId, readStatus },
  })
}

export function getUnReadCount<T>() {
  return post<T>({
    url: '/message/getUnReadCount',
  })
}

export function getTemplateListByFlag<T>(flags: string[]) {
  return post<T>({
    url: '/message/getTemplateListByFlag',
    data: { flags },
  })
}

export function deleteByMessageId<T>(messageId: number, isSend: boolean) {
  return post<T>({
    url: '/message/delete',
    data: {
      messageId,
      isSend,
    },
  })
}
