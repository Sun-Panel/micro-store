package api_v1

import (
	"sun-panel/api/api_v1/admin"
	"sun-panel/api/api_v1/oAuth2"
	"sun-panel/api/api_v1/openness"
	"sun-panel/api/api_v1/panel"
	"sun-panel/api/api_v1/proAuthorize"
	"sun-panel/api/api_v1/sunStore"
	"sun-panel/api/api_v1/system"
)

type ApiGroup struct {
	ApiSystem       system.ApiSystem // 系统功能api
	ApiOpen         openness.ApiPpenness
	ApiAdmin        admin.ApiAdmin
	ApiPanel        panel.ApiPanel
	ApiProAuthorize proAuthorize.ApiProAuthorize
	ApiOAuth2       oAuth2.OAuth2
	ApiSunStore     sunStore.ApiSunStore
}

var ApiGroupApp = new(ApiGroup)
