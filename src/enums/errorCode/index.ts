export enum ErrorCode {
  // 认证与权限类 (1000-1009)
  NotLoggedIn = 1000,
  IncorrectUsernameOrPassword = 1003,
  AccountDisabledOrNotActivated = 1004,
  NoCurrentPermission = 1005,
  AccountDoesNotExist = 1006,
  OldPasswordError = 1007,
  NoPROAuthorization = 1008,
  CaptchaError = 1009,

  // 数据操作类 (1200-1203)
  DatabaseError = 1200,
  PleaseKeepAtLeastOne = 1201,
  NoDataRecordFound = 1202,
  DataAlreadyExists = 1203,

  // 文件上传类 (1300-1301)
  UploadFailed = 1300,
  UnsupportedFileFormat = 1301,

  // 参数与业务类 (1400-1402)
  ParameterFormatError = 1400,
  OrderCreateFailed = 1401,
  GoodsNoUsePayPlatform = 1402,

  // 微应用版本业务类 (2000-2007)
  AppNotFound = 2000,
  VersionNotFound = 2001,
  VersionExists = 2002,
  VersionCodeExists = 2003,
  StatusNotAllowed = 2004,
  ApprovedCannotDelete = 2005,
  NotPendingReview = 2006,
  NoUpdateContent = 2007,
}
