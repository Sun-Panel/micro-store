import { post } from '@/utils/request'

export function getList<T>(data: Common.ListRequest) {
  return post<T>({
    url: '/admin/systemVariable/getList',
    data,
  })
}

export function edit<T>(data: SystemVariable.SystemVariableEditReq) {
  return post<T>({
    url: '/admin/systemVariable/edit',
    data,
  })
}

export function set<T>(name: string, value: string) {
  return post<T>({
    url: '/admin/systemVariable/set',
    data: {
      name,
      value,
    },
  })
}

export function deleteByName<T>(name: string) {
  return post<T>({
    url: '/admin/systemVariable/delete',
    data: {
      name,
    },
  })
}

export function clearCache<T>(varName: string) {
  return post<T>({
    url: '/admin/systemVariable/clearCache',
    data: { Name: varName },
  })
}

export function getByCache<T>(varName: string) {
  return post<T>({
    url: '/admin/systemVariable/getByCache',
    data: { Name: varName },
  })
}
