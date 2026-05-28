package v1

import (
	"sun-panel/api/apiClientApp/common/apiReturn"
	"sun-panel/api/apiClientApp/common/base"
	"sun-panel/api/apiClientApp/common/types"
	"sun-panel/lib/cmn"

	"sun-panel/biz"
	"sun-panel/global"
	"sun-panel/models"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type Api struct {
}

// ==================== 接口实现 ====================

// BatchGetDeveloperInfo 批量获取作者昵称
func (a *Api) BatchGetDeveloperInfo(c *gin.Context) {
	param := types.BatchGetDeveloperInfoReq{}
	if err := c.ShouldBindBodyWith(&param, binding.JSON); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}

	// 去重
	nameSet := make(map[string]bool)
	uniqueNames := make([]string, 0, len(param.DeveloperNames))
	for _, name := range param.DeveloperNames {
		if !nameSet[name] {
			nameSet[name] = true
			uniqueNames = append(uniqueNames, name)
		}
	}

	// 使用 biz 层带缓存的方法查询开发者信息
	developers, err := biz.Developer.BatchGetByDeveloperNames(global.Db, uniqueNames)
	if err != nil {
		apiReturn.ErrorDatabase(c, err.Error())
		return
	}

	// 构建响应map
	infos := make(map[string]types.DeveloperInfo, len(developers))
	for name, dev := range developers {
		infos[name] = types.DeveloperInfo{
			DeveloperName: dev.DeveloperName,
			Name:          dev.Name,
		}
	}

	apiReturn.SuccessData(c, types.BatchGetDeveloperInfoResp{
		DeveloperInfos: infos,
	})
}

// RecordInstall 微应用安装触发记录
func (a *Api) RecordInstall(c *gin.Context) {
	param := types.MicroAppInstallRecordReq{}
	if err := c.ShouldBindBodyWith(&param, binding.JSON); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}

	global.Logger.Infoln("RecordInstall", "param:", cmn.AnyToJsonStr(param))
	global.Logger.Infoln("RecordInstall", "client info:", cmn.AnyToJsonStr(base.GetCurrentClient(c)))

	// 获取当前用户
	user, exists := base.GetCurrentUser(c)
	if !exists {
		global.Logger.Debugln("RecordInstall", "user not exists")
	}
	global.Logger.Infoln("RecordInstall", "user:", cmn.AnyToJsonStr(user))

	// // 根据 microAppId 获取应用信息
	// microApp := models.MicroApp{}
	// app, err := microApp.GetByMicroAppId(global.Db, param.MicroAppId)
	// if err != nil {
	// 	apiReturn.Error(c, "应用不存在")
	// 	return
	// }

	// // 根据版本号获取版本信息
	// version := models.MicroAppVersion{}
	// versionInfo, err := version.GetByVersion(global.Db, param.MicroAppVersion)
	// if err != nil {
	// 	apiReturn.Error(c, "版本不存在")
	// 	return
	// }

	// _ = app
	// _ = versionInfo

	// TODO: 获取客户端标识（ClientId），需要确认从哪里获取
	// 解析当前token是否已经绑定了用户
	// 保存安装记录。用户信息为空也可以保存记录
	// clientId := c.GetHeader("Client-Id")
	// if clientId == "" {
	// 	clientId = "unknown"
	// }

	// 获取用户ID（可选，可能未登录）
	// user, ok := base.GetCurrentUser(c)
	// if !ok {
	// 	apiReturn.ErrorNotLogin(c)
	// 	return
	// }

	// // 获取IP
	// publicIP := c.ClientIP()

	// // 创建安装记录
	// install := models.MicroAppInstall{
	// 	AppRecordId: app.ID,
	// 	VersionId:   versionInfo.ID,
	// 	UserId:      user.ID,
	// 	ClientId:    clientId,
	// 	PublicIp:    publicIP,
	// }

	// if err := install.Create(global.Db); err != nil {
	// 	apiReturn.ErrorDatabase(c, err.Error())
	// 	return
	// }

	// TODO: 是否需要更新应用的安装次数？与现有逻辑是否冲突？

	apiReturn.Success(c)
}

// BatchGetMicroAppInfo 批量查询微应用信息
func (a *Api) BatchGetMicroAppInfo(c *gin.Context) {
	param := types.BatchGetMicroAppInfoReq{}
	if err := c.ShouldBindBodyWith(&param, binding.JSON); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}

	// 去重
	idSet := make(map[string]bool)
	uniqueIds := make([]string, 0, len(param.MicroAppIds))
	for _, id := range param.MicroAppIds {
		if !idSet[id] {
			idSet[id] = true
			uniqueIds = append(uniqueIds, id)
		}
	}

	// 使用 biz 层带缓存的方法查询
	infos, err := biz.MicroApp.BatchGetMicroAppInfo(global.Db, uniqueIds)
	if err != nil {
		apiReturn.ErrorDatabase(c, err.Error())
		return
	}

	// 构建响应map
	result := make(map[string]*types.MicroAppInfo, len(infos))
	for id, info := range infos {
		result[id] = &types.MicroAppInfo{
			MicroAppId:     info.MicroAppId,
			AppIcon:        info.AppIcon,
			ChargeType:     info.ChargeType,
			Points:         info.Points,
			DeveloperName:  info.DeveloperName,
			DeveloperName2: info.DeveloperName2,
		}
	}

	apiReturn.SuccessData(c, types.BatchGetMicroAppInfoResp{
		MicroAppInfos: result,
	})
}

// CheckIsDeveloper 查询当前账号是否为开发者
func (a *Api) CheckIsDeveloper(c *gin.Context) {
	user, exists := base.GetCurrentUser(c)
	if !exists {
		apiReturn.ErrorNotLogin(c)
		return
	}

	developer := models.Developer{}
	isDeveloper, err := developer.CheckUserIsDeveloper(global.Db, user.ID)
	if err != nil {
		apiReturn.ErrorDatabase(c, err.Error())
		return
	}

	apiReturn.SuccessData(c, types.CheckIsDeveloperResp{
		IsDeveloper: isDeveloper,
	})
}
