/** 微应用状态枚举 */
export enum MicroAppStatus {
  /** 下架 */
  OFFLINE = 0,
  /** 上架 */
  ONLINE = 1,
  /** 审核中 */
  PENDING = 2,
}

/** 微应用状态映射 */
export const microAppStatusMap: Record<number, string> = {
  [MicroAppStatus.OFFLINE]: '下架',
  [MicroAppStatus.ONLINE]: '上架',
  [MicroAppStatus.PENDING]: '审核中',
}

/** 微应用收费方式枚举 */
export enum MicroAppChargeType {
  /** 免费 */
  FREE = 0,
  /** 积分 */
  POINTS = 1,
  /** PRO免费 */
  PRO_FREE = 2,
}

/** 微应用收费方式映射 */
export const microAppChargeTypeMap: Record<number, string> = {
  [MicroAppChargeType.FREE]: '免费',
  [MicroAppChargeType.POINTS]: '积分',
  [MicroAppChargeType.PRO_FREE]: 'PRO免费',
}
