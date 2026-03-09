import { post } from '@/utils/request'

// 保存
export function save<T>(content: string) {
  return post<T>({
    url: '/admin/aboutSetting/save',
    data: {
      content,
    },
  })
}
