import { post } from '@/utils/request'

// 添加keys
export function getAll<T>() {
  return post<T>({
    url: '/admin/clientSetting/blacklistIP/getList',
  })
}

export function deletes<T>(ips: string[]) {
  return post<T>({
    url: '/admin/clientSetting/blacklistIP/deletes',
    data: ips,
  })
}

export function set<T>(ips: string[], second: number) {
  return post<T>({
    url: '/admin/clientSetting/blacklistIP/set',
    data: { ips, second },
  })
}
