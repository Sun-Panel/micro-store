import { post } from '@/utils/request'

export function edit<T>(version: Version.Info) {
  return post<T>({
    url: '/admin/version/edit',
    data: version,
  })
}

export function setActive<T>(id: number, type: 'release' | 'beta' | 'alpha' | 'rc' | string) {
  return post<T>({
    url: '/admin/version/setActive',
    data: {
      id,
      versionType: type,
    },
  })
}

export function deletes<T>(ids: number[]) {
  return post<T>({
    url: '/admin/version/deletes',
    data: { ids },
  })
}

export function getList<T>(data: Common.ListRequest) {
  return post<T>({
    url: '/admin/version/getList',
    data,
  })
}

export function editSecret<T>(secret: Version.SecretInfo) {
  return post<T>({
    url: '/admin/versionSecret/edit',
    data: secret,
  })
}

export function getSecretByVersion<T>(version: string) {
  return post<T>({
    url: '/admin/versionSecret/getByVersion',
    data: { version },
  })
}
