declare namespace MicroAppStatistics {
  // 下载统计信息
  interface DownloadStatistics {
    appId: number
    downloadCount: number
    installCount: number
  }

  // 批量下载统计信息
  interface BatchDownloadStatistics {
    [appId: number]: [number, number]
  }

  // 记录下载请求
  interface IncrementDownloadRequest {
    appId: number
    versionId?: number
    userId?: number
    clientId: string
    downloadIp?: string
  }

  // 获取统计请求
  interface GetStatisticsRequest {
    appId: number
  }

  // 批量获取统计请求
  interface GetBatchStatisticsRequest {
    appIds: number[]
  }

  // 获取统计响应
  interface GetStatisticsResponse {
    appId: number
    downloadCount: number
    installCount: number
  }

  // 批量获取统计响应
  interface GetBatchStatisticsResponse {
    [appId: number]: [number, number]
  }
}
