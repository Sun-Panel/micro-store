import { post } from '@/utils/request'

// 获取列表
export function getList<T>() {
  return post<T>({
    url: '/admin/emailTemplate/getList',
  })
}

// 删除
export function deletes<T>(ids: number[]) {
  return post<T>({
    url: '/admin/emailTemplate/deletes',
    data: { ids },
  })
}

// 修改
export function edit<T>(info: EmailTemplate.Info) {
  return post<T>({
    url: '/admin/emailTemplate/edit',
    data: info,
  })
}
