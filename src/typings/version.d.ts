declare namespace Version {
  interface Info {
    id?: number
    version: string
    type: 'release' | 'beta' | 'alpha' | 'rc' | 'dev'
    releaseTime: string // 时间戳字符串
    description: string
    downloadURL: string
    pageUrl: string
    isActive?: boolean
    isRolledBack: boolean
    aloneSecretKey: int
  }

  interface SecretInfo {
    id?: number
    version: string
    secretKey: string
    status: boolean
  }
}
