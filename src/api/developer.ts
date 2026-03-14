import { post } from '@/utils/request'

// 申请成为开发者
export function register<T>(data: Developer.RegisterRequest) {
  return post<T>({
    url: '/developer/register',
    data,
  })
}

// 获取当前用户的开发者信息
export function getInfo<T>() {
  return post<T>({
    url: '/developer/getInfo',
  })
}

// 更新开发者信息
export function updateMyInfo<T>(data: Developer.UpdateMyInfoRequest) {
  return post<T>({
    url: '/developer/update',
    data,
  })
}

// 检查是否是开发者
export function checkIsDeveloper<T>() {
  return post<T>({
    url: '/developer/checkIsDeveloper',
  })
}
