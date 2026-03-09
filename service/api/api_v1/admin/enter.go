package admin

type ApiAdmin struct {
	UsersApi                   UsersApi
	SystemSettingApi           SystemSettingApi
	AboutSettingApi            AboutSettingApi
	EmailApi                   EmailApi
	ProAuthorizeApi            ProAuthorizeApi
	SystemVariableApi          SystemVariableApi
	MdPageApi                  MdPageApi
	DashboardApi               DashboardApi
	VersionApi                 VersionApi
	VersionSecretApi           VersionSecretApi
	ClientBlackListIPApi       ClientBlackListIPApi
	ClientCreateOnlineCacheApi ClientCreateOnlineCacheApi
	RedeemCodeApi              RedeemCodeApi
}
