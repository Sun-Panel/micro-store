import { post } from '@/utils/request'

// 获取列表
export function getList<T>() {
  return post<T>({
    url: '/admin/notice/getList',
  })
}

// 删除
export function deletes<T>(ids: number[]) {
  return post<T>({
    url: '/admin/notice/deletes',
    data: { ids },
  })
}

// 修改
export function edit<T>(info: Notice.NoticeInfo) {
  return post<T>({
    url: '/admin/notice/edit',
    data: info,
  })
}
