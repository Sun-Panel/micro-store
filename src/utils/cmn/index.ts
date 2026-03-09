import moment from 'moment'
import { h } from 'vue'
import type { NotificationReactive } from 'naive-ui'
import { NButton, createDiscreteApi } from 'naive-ui'
import { apiRespErrMsg as apiRespErrMsgApiMessage } from './apiMessage'
import { useAuthStore, useNoticeStore, useUserStore } from '@/store'
import { getAuthInfo } from '@/api/system/user'
import type { VisitMode } from '@/enums/auth'
import { getListByDisplayType as getListByDisplayTypeApi } from '@/api/notice'
import { OrderStatus, PayPlatform } from '@/enums/goodsOrder'

const noticeStore = useNoticeStore()
const userStore = useUserStore()
const authStore = useAuthStore()

const { notification } = createDiscreteApi(['notification'])
/**
 * 生成指定时间格式
 * @param format 时间格式 默认：'YYYY-MM-DD HH:mm:ss'
 * @returns string
 */
export function buildTimeString(timeString?: string, format?: string): string | null {
  if (timeString === '')
    return null

  if (!format)
    format = 'YYYY-MM-DD HH:mm:ss'

  return moment(timeString).format(format)
}

export function timeFormat(timeString?: string) {
  return moment(timeString).format('YYYY-MM-DD HH:mm:ss')
}

export function getCurrencySymbol(currencyCode: string): string | null {
  // 创建一个对象来存储常见的货币代码和对应的符号
  const currencySymbols: { [key: string]: string } = {
    CNY: '¥',
    USD: '$',
    EUR: '€',
    GBP: '£',
    JPY: '¥',
    // 在这里添加更多的货币代码和对应的符号
  }

  return currencySymbols[currencyCode] || null
}

/**
 * 创建新的公告
 * @param timeString
 */
export function noticeCreate(info: Notice.NoticeInfo) {
  const option: any = {
    title: info.title,
    content: info.content,
    meta: info.createTime ? timeFormat(info.createTime) : '',
  }

  const btns: any = []

  let n: NotificationReactive
  // 链接按钮
  if (info.url !== '') {
    btns.push(
      h(
        NButton,
        {
          text: true,
          type: 'info',
          onClick: () => {
            window.open(info.url, '_blank')
            n.destroy()
          },
        },
        {
          default: () => '打开链接',
        },
      ),
    )
  }
  if (info.oneRead === 1) {
    btns.push(
      h(
        NButton,
        {
          text: true,
          type: 'primary',
          style: { marginLeft: '20px' },
          onClick: () => {
            if (info.id) {
              if (info.isLogin === 1 && userStore.userInfo.username) {
                noticeStore.setReadByUsername(userStore.userInfo.username, info.id)
                console.log('设置用户已读', info.id)
              }
              else {
                noticeStore.setReadByGlobal(info.id)
                console.log('设置全局已读', info.id)
              }
            }
            n.destroy()
          },
        },
        {
          default: () => '不再提醒',
        },
      ),
    )
  }
  option.action = () => btns
  n = notification.create(option)
}

export function setTitle(titile: string) {
  document.title = titile
}

export function getTitle(titile: string) {
  document.title = titile
}

//
export async function updateLocalUserInfo() {
  interface Req {
    user: User.Info
    visitMode: VisitMode
  }

  const { data } = await getAuthInfo<Req>()
  userStore.updateUserInfo({ headImage: data.user.headImage, name: data.user.name })
  authStore.setUserInfo(data.user)
  authStore.setVisitMode(data.visitMode)
}

export async function getNotice(displayType: number | number[]) {
  let param: number[]
  if (typeof displayType === 'number')
    param = [displayType]
  else
    param = displayType

  const { data } = await getListByDisplayTypeApi<Common.ListResponse<Notice.NoticeInfo[]>>(param)

  for (let i = 0; i < data.list.length; i++) {
    const element = data.list[i]
    if (element.id && !noticeStore.getReadByNoticeId(element.id, userStore.userInfo.username))
      noticeCreate(element)
  }
}

export function getFaviconUrl(url: string): string {
  // 获取网址的域名
  const { protocol, host } = new URL(url)
  const domain = `${protocol}//${host}`
  // 构建 favicon URL
  return `${domain}/favicon.ico`
}

/**
 * @description: 获取随机码
 * @param {number} size
 * @param {array} seed ["a","b"m"c]
 * @return {string}
 */
export function randomCode(size: number, seed?: Array<string>) {
  seed = seed || ['A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z',
    'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'm', 'n', 'p', 'Q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z',
    '2', '3', '4', '5', '6', '7', '8', '9',
  ]// 数组
  const seedlength = seed.length// 数组长度
  let createPassword = ''
  for (let i = 0; i < size; i++) {
    const j = Math.floor(Math.random() * seedlength)
    createPassword += seed[j]
  }
  return createPassword
}

// 复制文字到剪切板
export async function copyToClipboard(text: string): Promise<boolean> {
  if (navigator.clipboard) {
    // 使用 Clipboard API
    try {
      await navigator.clipboard.writeText(text)
      return true
    }
    catch (err) {
      console.error('copy fail', err)
      return false
    }
  }
  else {
    // 兼容旧版浏览器
    const textArea = document.createElement('textarea')
    textArea.value = text
    document.body.appendChild(textArea)
    textArea.select()

    try {
      document.execCommand('copy')
      return true
    }
    catch (err) {
      console.error('copy fail', err)
      return false
    }
    finally {
      document.body.removeChild(textArea)
    }
  }
}

export function bytesToSize(bytes: number) {
  const sizes = ['B', 'KB', 'MB', 'GB', 'TB', 'PB', 'EB']
  if (bytes === 0)
    return '0B'
  const i = parseInt(String(Math.floor(Math.log(bytes) / Math.log(1024))))
  return `${(bytes / 1024 ** i).toFixed(1)} ${sizes[i]}`
}

export function getPayPlatformText(code: number): string {
  switch (code) {
    case PayPlatform.ALI:
      return '支付宝'
      break

    case PayPlatform.WEIXIN:
      return '微信'
      break
    case PayPlatform.PADDLE:
      return 'PayPal/信用卡'
      break

    default:
      return '未知'
      break
  }
}

export function getPayStatusText(code: number): string {
  switch (code) {
    case OrderStatus.CLOSE:
      return '关闭'
      break

    case OrderStatus.PAY_WAIT:
      return '等待支付'
      break

    case OrderStatus.PAY_SUCCESS:
      return '支付成功'
      break

    default:
      return '未知'
      break
  }
}

export function formatAmount(amount: number): string {
  // 定义单位后缀
  const suffixes = ['', 'k', 'M', 'B', 'T']

  // 寻找合适的单位后缀
  let suffixIndex = 0
  while (amount >= 1000 && suffixIndex < suffixes.length - 1) {
    amount /= 1000
    suffixIndex++
  }

  // 格式化金额
  let formattedAmount: string
  if (suffixIndex === 0)
    formattedAmount = amount.toFixed(2) // 如果没有单位后缀，则保留两位小数
  else
    formattedAmount = amount.toFixed(1) // 使用一个小数点来表示k、M、B、T等单位

  // 添加单位后缀
  formattedAmount += suffixes[suffixIndex]

  return formattedAmount
}

export const apiRespErrMsg = apiRespErrMsgApiMessage
