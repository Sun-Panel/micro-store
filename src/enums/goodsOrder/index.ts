export enum PayPlatform {
  'ALI' = 1, // 支付平台 阿里
  'WEIXIN' = 2, // 支付平台 微信
  'PADDLE' = 3, // 支付平台 Paddle
}

export enum OrderStatus {
  'PAY_WAIT' = 1, // 支付状态 等待付款/未付款
  'PAY_SUCCESS' = 2, // 支付状态 已付款
  'CLOSE' = 3, // 支付状态 关闭订单
  'CANCEL' = 4, // 支付状态 取消订单
  'PAY_REVIEW' = 5, // 付款后等待审核
  'FINISH' = 6, // 支付状态 已完成订单
}
