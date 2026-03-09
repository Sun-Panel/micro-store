import { get } from '@/utils/request'

export function getLoginConfig<T>() {
  return get<T>({
    url: '/openness/loginConfig',
  })
}

// 获取免责声明
export function getDisclaimer<T>() {
  return get<T>({
    url: '/openness/getDisclaimer',
  })
}

// 获取关于的描述信息
export function getAboutDescription<T>() {
  return get<T>({
    url: '/openness/getAboutDescription',
  })
}

export function getHomeBase<T>() {
  return get<T>({
    url: '/openness/getHomeBase',
  })
}

export function getProDescription<T>() {
  return get<T>({
    url: '/openness/getProDescription',
  })
}

export function getRootPageDescription<T>() {
  return get<T>({
    url: '/openness/getRootPageDescription',
  })
}
