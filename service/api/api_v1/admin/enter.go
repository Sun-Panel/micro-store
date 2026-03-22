package admin

type ApiAdmin struct {
	UsersApi                   UsersApi
	SystemSettingApi           SystemSettingApi
	AboutSettingApi            AboutSettingApi
	EmailApi                   EmailApi
	SystemVariableApi          SystemVariableApi
	MdPageApi                  MdPageApi
	DashboardApi               DashboardApi
	ClientBlackListIPApi       ClientBlackListIPApi
	ClientCreateOnlineCacheApi ClientCreateOnlineCacheApi
	RedeemCodeApi              RedeemCodeApi
	MicroAppCategoryApi        MicroAppCategoryApi
	DeveloperApi               DeveloperApi
	// 微应用 API（按角色划分）
	MicroAppAdminApi          MicroAppAdminApi
	MicroAppAuditorApi        MicroAppAuditorApi
	MicroAppDeveloperApi      MicroAppDeveloperApi
	// 版本 API（按角色划分）
	DeveloperVersionApi       DeveloperVersionApi
	MicroAppVersionAdminApi   MicroAppVersionAdminApi
	MicroAppVersionUploadApi  MicroAppVersionUploadApi
}
