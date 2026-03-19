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
	MicroAppApi                MicroAppApi
	DeveloperVersionApi        DeveloperVersionApi
	MicroAppVersionUploadApi   MicroAppVersionUploadApi
}
