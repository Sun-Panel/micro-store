import { post } from '@/utils/request'

// 下发重置密码的验证码到邮箱
export function sendResetPasswordVCode<T>(email: string, vcode: string) {
  return post<T>({
    url: '/login/sendResetPasswordVCode',
    data: { email, vcode },
  })
}

// 下发重置密码的验证码到邮箱
export function resetPasswordByVCode<T>(data: Login.ResetPasswordByVCodeReqest) {
  return post<T>({
    url: '/login/resetPasswordByVCode',
    data,
  })
}

export function oAuth2CodeLogin<T>(code: string) {
  return post<T>({
    url: '/oAuth2CodeLogin',
    data: { code },
  })
}

export function oAuth2CodeBind<T>(code: string) {
  return post<T>({
    url: '/oAuth2CodeBind',
    data: { code },
  })
}
