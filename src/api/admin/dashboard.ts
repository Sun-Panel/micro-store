import { post } from '@/utils/request'

// 获取统计信息
export function getStatistics<T>() {
  return post<T>({
    url: '/admin/dashboard/getStatistics',
  })
}

export function getUserLine<T>(data: string[]) {
  return post<T>({
    url: '/admin/dashboard/getUserLine', data,
  })
}

export function getClientLine<T>(data: string[]) {
  return post<T>({
    url: '/admin/dashboard/getClientLine', data,
  })
}

export function getActiveClientVersionStatistics<T>() {
  return post<T>({
    url: '/admin/dashboard/getActiveClientVersionStatistics',
  })
}

export function getVersions<T>() {
  return post<T>({
    url: '/admin/dashboard/getVersions',
  })
}

export function getVersionHistory<T>() {
  return post<T>({
    url: '/admin/dashboard/getVersionHistory',
  })
}
