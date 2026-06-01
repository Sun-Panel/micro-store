import { getDownloadUrl as getDownloadUrlApi } from '@/api/microApp'

/**
 * 获取微应用下载地址
 * @param microAppId 微应用ID
 * @param version 版本号（可选）
 * @returns 下载地址
 */
export async function getDownloadUrl(microAppId: string, version?: string): Promise<string> {
  let url = ''
  await getDownloadUrlApi<string>(microAppId, version).then(({ data }) => {
    url = data
  }).catch((res) => {
    console.error(`get url error: ${res}`)
  })
  // 如果是相对路径，拼接当前域名
  if (url && !url.startsWith('http://') && !url.startsWith('https://')) {
    url = `${window.location.origin}${url.startsWith('/') ? '' : '/'}${url}`
  }
  return url
}
