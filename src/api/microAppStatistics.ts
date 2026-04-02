import { post } from '@/utils/request'

// ==================== 微应用统计接口 ====================

/**
 * 记录下载
 * @param data 下载信息
 */
export function incrementDownload<T>(data: {
  appId: number
  versionId?: number
  userId?: number
  clientId: string
  downloadIp?: string
}) {
  return post<T>({
    url: '/api/microApp/download/increment',
    data,
  })
}

/**
 * 获取下载统计
 * @param appId 应用ID
 */
export function getDownloadStatistics<T>(appId: number) {
  return post<T>({
    url: '/api/microApp/statistics/get',
    data: { appId },
  })
}

/**
 * 批量获取下载统计
 * @param appIds 应用ID列表
 */
export function getBatchDownloadStatistics<T>(appIds: number[]) {
  return post<T>({
    url: '/api/microApp/statistics/batchGet',
    data: { appIds },
  })
}

/**
 * 同步 Redis 统计到数据库（管理员接口）
 */
export function syncRedisStatistics<T>() {
  return post<T>({
    url: '/api/admin/microApp/statistics/sync',
    data: {},
  })
}
