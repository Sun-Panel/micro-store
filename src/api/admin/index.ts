import { post } from '@/utils/request'

// ===================
// 管理系统 管理员
// ===================

// 用户相关
export function AdminUserManageGetList<T>(param: AdminUserManage.GetListRequest) {
  return post<T>({
    url: '/admin/users/getList',
    data: param,
  })
}

export function AdminUserManageEdit<T>(param: User.Info) {
  let url = '/admin/users/create'
  if (param.id)
    url = '/admin/users/update'

  return post<T>({
    url,
    data: param,
  })
}

export function AdminUserManageDelete<T>(userIds: number[]) {
  return post<T>({
    url: '/admin/users/deletes',
    data: { userIds },
  })
}

// 系统设置
export function AdminSystemSettingSetEmail<T>(email: AdminSystemSetting.Email) {
  return post<T>({
    url: '/admin/systemSetting/setEmail',
    data: email,
  })
}

export function AdminSystemSettingGetEmail<T>() {
  return post<T>({
    url: '/admin/systemSetting/getEmail',
  })
}

export function AdminSystemSettingGetWebsiteSetting<T>() {
  return post<T>({
    url: '/admin/systemSetting/getApplicationSetting',
  })
}

export function AdminSystemSettingSetWebsiteSetting<T>(data: AdminSystemSetting.Website) {
  return post<T>({
    url: '/admin/systemSetting/setApplicationSetting',
    data,
  })
}

export function adminSystemSettingRoleManageGetSystemList<T>(data: Common.ListRequest) {
  return post<T>({
    url: '/admin/roleManage/getSystemList',
    data,
  })
}

export function adminSystemSettingRoleManageGetInfo<T>(aiRoleId: number) {
  return post<T>({
    url: '/admin/roleManage/getInfo',
    data: { aiRoleId },
  })
}

export function adminSystemSettingRoleManageDeletes<T>(aiRoleIds: number[]) {
  return post<T>({
    url: '/admin/roleManage/deletes',
    data: { aiRoleIds },
  })
}
