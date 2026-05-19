/**
 * iframe 通信 composable
 * 封装 useIframeCommunication，自动判断登录状态并发送登录事件
 * 已登录时传递账号信息，未登录时账号为空
 */

import type { UseIframeCommunicationReturn } from './useIframeCommunication'
import { logout } from '@/api'
import { iframeCaptchaLogin } from '@/api/login'

import { useAuthStore } from '@/store'
import { debug } from '@/utils/logger'
import { MICRO_APP_EVENTS, MICRO_APP_STORE_EVENTS } from './iframeEvents'
import { getIframeAllUrlParam } from './useGetIframeUrlParam'
import { useIframeCommunication } from './useIframeCommunication'

export interface LoginEventData {
  /** 是否已登录 */
  loggedIn: boolean
  /** 账号，未登录时为空 */
  captcha: string
}

export type LoginEventHandler = (data: LoginEventData) => void

export interface UseIframeConfig {
  /** 目标 origin，默认 '*' */
  targetOrigin?: string
  /** 来源标识，默认 'microAppStore' */
  sourceId?: string
  /** postMessage 安全密钥 */
  postMessageKey?: string
  /** 是否启用调试模式 */
  debug?: boolean
}

export interface UseIframeReturn extends UseIframeCommunicationReturn {
  /** 发送登录事件，根据当前登录状态自动填充账号信息 */
  sendLoginEvent: () => void
}

/**
 * 使用 iframe 通信
 * 自动判断登录状态，发送登录事件时已登录传账号，未登录传空
 */
export function useIframe(config: UseIframeConfig = {}): UseIframeReturn {
  const urlParam = getIframeAllUrlParam()
  const authStore = useAuthStore()
  let loginIn = false // 本次是否已登录

  const communication = useIframeCommunication({
    targetOrigin: config.targetOrigin || '*',
    sourceId: config.sourceId || 'microAppStore',
    isIframe: true,
    debug: config.debug ?? import.meta.env.DEV,
    autoDestroy: true,
    postMessageKey: urlParam?.postMessageKey,
  })

  // iframe 端自动注册父窗口为 'main' 目标
  if (window.parent && window.parent !== window) {
    communication.channel.registerTarget('main', window.parent)
  }

  /**
   * 判断当前是否已登录
   */
  function isLoggedIn(): boolean {
    return !!authStore.token
  }

  /**
   * 发送登录事件
   * 已登录时传递用户账号信息，未登录时账号为空
   */
  function sendLoginEvent() {
    debug(`触发事件:${MICRO_APP_EVENTS.COMMUNICATION_READY}`, 'authStore:', authStore)
    const loggedIn = isLoggedIn()
    communication.sendMessage('main', MICRO_APP_EVENTS.COMMUNICATION_READY, {
      loggedIn,
      account: loggedIn ? (authStore.userInfo?.username || authStore.userInfo?.name || '') : '',
    })
  }

  // 监听主面板发送的登录事件，自动同步登录状态
  communication.on(MICRO_APP_STORE_EVENTS.LOGIN, async (data: LoginEventData) => {
    debug(`监听事件:${MICRO_APP_STORE_EVENTS.LOGIN}`, data)
    if (!loginIn) {
      try {
        const res = await iframeCaptchaLogin<Login.LoginResponse>(data.captcha)
        if (res.data) {
          debug(`请求登录:`, res.data)
          authStore.setUserInfo({
            ...authStore.userInfo,
            ...res.data,
          })
          authStore.setToken(res.data.token)
          loginIn = true
        }
      }
      catch (error) {
        console.error('Iframe登录同步失败:', error)
      }
    }
    else {
      authStore.removeToken()
    }
  })

  // 监听主面板发送的登出事件
  communication.on(MICRO_APP_STORE_EVENTS.LOGOUT, async () => {
    debug(`监听事件:${MICRO_APP_STORE_EVENTS.LOGOUT}`)
    await logout()
    // userStore.resetUserInfo()
    authStore.removeToken()
  })

  return {
    ...communication,
    sendLoginEvent,
  }
}
