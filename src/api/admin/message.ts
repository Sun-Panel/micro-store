import { post } from '@/utils/request'

// 获取列表
export function getTemplateList<T>() {
  return post<T>({
    url: '/admin/message/tempalte/getList',
  })
}

// 删除
export function templateDeletes<T>(ids: number[]) {
  return post<T>({
    url: '/admin/message/tempalte/deletes',
    data: { ids },
  })
}

// 修改
export function templateEdit<T>(info: Message.TemplateInfo) {
  return post<T>({
    url: '/admin/message/tempalte/edit',
    data: info,
  })
}
