/** 微应用状态枚举 */
export enum MicroAppStatus {
  /** 下架 */
  OFFLINE = 0,
  /** 上架 */
  ONLINE = 1,
  /** 审核中 */
  PENDING = 2,
}

/** 微应用审核状态枚举（用于区分版本状态） */
export enum MicroAppReviewStatus {
  /** 审核通过（生效版本） */
  APPROVED = 0,
  /** 审核中 */
  REVIEWING = 1,
  /** 审核拒绝 */
  REJECTED = 2,
  /** 草稿 */
  DRAFT = 3,
}

/** 微应用审核状态映射 */
export const microAppReviewStatusMap: Record<number, string> = {
  [MicroAppReviewStatus.APPROVED]: '已通过',
  [MicroAppReviewStatus.REVIEWING]: '审核中',
  [MicroAppReviewStatus.REJECTED]: '已拒绝',
  [MicroAppReviewStatus.DRAFT]: '草稿',
}

/** 微应用状态映射 */
export const microAppStatusMap: Record<number, string> = {
  [MicroAppStatus.OFFLINE]: '已下架',
  [MicroAppStatus.ONLINE]: '已上架',
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

/** 微应用版本状态枚举 */
export enum MicroAppVersionStatus {
  /** 草稿 */
  DRAFT = -1,
  /** 待审核 */
  PENDING = 0,
  /** 通过 */
  APPROVED = 1,
  /** 拒绝 */
  REJECTED = 2,
  /** 已下架 */
  OFFLINE = 3,
}

/** 微应用版本状态映射 */
export const microAppVersionStatusMap: Record<number, string> = {
  [MicroAppVersionStatus.DRAFT]: '草稿',
  [MicroAppVersionStatus.PENDING]: '待审核',
  [MicroAppVersionStatus.APPROVED]: '已通过',
  [MicroAppVersionStatus.REJECTED]: '已拒绝',
  [MicroAppVersionStatus.OFFLINE]: '已下架',
}
