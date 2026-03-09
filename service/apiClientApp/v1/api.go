package v1

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sun-panel/apiClientApp/v1/common/apiReturn"
	"sun-panel/apiClientApp/v1/common/base"
	"sun-panel/apiClientApp/v1/common/types"
	"sun-panel/biz"
	"sun-panel/biz/clientCache"
	"sun-panel/global"
	"sun-panel/lib/cmn"
	"sun-panel/lib/cmn/systemSetting"
	"sun-panel/lib/sunStore"
	"sun-panel/models"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Api struct {
}

// 注册
func (a *Api) Register(c *gin.Context) {
	req := types.RegisterReq{}
	secretKey := base.GetVersionSecretKey(c)

	if err := base.GetRequestParam(c, &req); err != nil {
		global.Logger.Error(apiReturn.ErrorWithNumReturnErrStr(c, err))
		return
	}

	// 校验时间戳
	if ok := base.CheckTimestamp(req.Timestamp); !ok {
		global.Logger.Errorln(apiReturn.ErrorWithNumReturnErrStr(c, errors.New("abnormal timestamp")), "IP:", c.ClientIP(), "param:", cmn.AnyToJsonStr(req))
		return
	}

	global.Logger.Info(a.log("Register", c.ClientIP(), req))

	// 黑名单IP，返回404
	if biz.ClientCache.BlacklistIP.CheckAndUpdate(c.ClientIP()) {
		global.Logger.Info("Blacklist IP:", c.ClientIP())
		c.JSON(404, gin.H{})
		return
	}

	// 判断是否存在clientId,不存在返回一个clientId
	clientInfo := models.SoftwareClient{}

	clientInfo.Version = req.Version
	clientInfo.MacAddress = req.MacAddress
	clientInfo.LastLanIp = req.LanIP
	clientInfo.LastIp = c.ClientIP()
	clientInfo.JoinTime = time.Now()
	clientInfo.LastOnlineTime = time.Now()

	if req.ClientId == nil || *req.ClientId == "" {
		// 创建一个clientId数据
		newClientIdIsExist := true
		newClientId := ""
		for newClientIdIsExist {
			newClientId = cmn.BuildRandCode(10, cmn.RAND_CODE_MODE2)
			if ok, err := a.clientIDIsExist(newClientId); !ok {
				global.Logger.Debugln("newClientId:", newClientId)
				// 库中不存在此clientId
				newClientIdIsExist = false
				break
			} else if err != nil {
				global.Logger.Error(apiReturn.ErrorWithNumReturnErrStr(c, err))
				return
			}
		}

		clientInfo.ClientId = newClientId
		// =================== 注册新设备在数据库中(已经修改为缓存逻辑)- 2024年12月30日
		// if err := global.Db.Create(&clientInfo).Error; err != nil {
		// 	global.Logger.Error(apiReturn.ErrorWithNumReturnErrStr(c, err))
		// 	return
		// }
		// ===================

		// 写入到缓存中，定时任务同步客户端数据库 - 2024年12月30日
		biz.ClientCache.HistoryCache.WriteToNewClientsCache(newClientId, clientInfo)

		result, paramError := base.GetRequestResp(secretKey, &types.RegisterResp{
			ClientId: newClientId,
			RespBase: base.GetRespBase(""),
		})
		if paramError != nil {
			global.Logger.Error(apiReturn.ErrorWithNumReturnErrStr(c, paramError))
			return
		}
		apiReturn.SuccessData(c, result)
	} else {
		// ========================= 更新设备的数据在数据库中 - 已废弃，将采用全新使用缓存逻辑 2024年12月30日
		// if err := global.Db.First(&models.SoftwareClient{}, "client_id=?", req.ClientId).Error; err != nil {
		// 	if errors.Is(gorm.ErrRecordNotFound, err) {
		// 		// 客户端不存在，自动创建客户端
		// 		if errors.Is(err, gorm.ErrRecordNotFound) {
		// 			// 客户端id不存在、为空、长度大于40全部不允许登录
		// 			if req.ClientId == nil || *req.ClientId == "" || len(*req.ClientId) > 40 {
		// 				global.Logger.Error(apiReturn.ErrorWithNumReturnErrStr(c, errors.New("clientID incorrect")))
		// 				return
		// 			}
		// 			clientInfo.ClientId = *req.ClientId
		// 			global.Db.Create(&clientInfo)
		// 			global.Logger.Infoln("client id not exist，auto create:", cmn.AnyToJsonStr(clientInfo))
		// 		} else {
		// 			global.Logger.Error(apiReturn.ErrorWithNumReturnErrStr(c, err))
		// 			return
		// 		}
		// 	}
		// 	global.Logger.Error(apiReturn.ErrorWithNumReturnErrStr(c, err))
		// 	return
		// } else {
		// 	updateFields := []string{
		// 		"Version",
		// 		"MacAddress",
		// 		"LastIp",
		// 		"LastLanUrl",
		// 	}
		// 	if err := global.Db.Model(&models.SoftwareClient{}).Where("client_id=?", req.ClientId).Select(updateFields).Updates(clientInfo).Error; err != nil {
		// 		global.Logger.Error(apiReturn.ErrorWithNumReturnErrStr(c, err))
		// 		return
		// 	}
		// }
		// ===========================

		clientInfo.ClientId = *req.ClientId

		// 全新使用缓存逻辑 2024年12月30日
		if err := a.checkRequestClientIdAndCreate(req.ClientId, clientInfo); err != nil {
			global.Logger.Error(apiReturn.ErrorWithNumReturnErrStr(c, err))
			return
		}

		// 更新数据
		biz.ClientCache.HistoryCache.WriteToHourUpdateClientsCache(clientInfo.ClientId, clientInfo)

		result, paramError := base.GetRequestResp(secretKey, &types.RegisterResp{
			ClientId: clientInfo.ClientId,
			RespBase: base.GetRespBase(""),
		})
		if paramError != nil {
			global.Logger.Error(apiReturn.ErrorWithNumReturnErrStr(c, paramError))
			return
		}
		apiReturn.SuccessData(c, result)
	}

	global.Logger.Infoln("client register", "IP:", c.ClientIP(), "clientId:", clientInfo.ClientId)

}

// 登录 - 不依赖授权平台单独运行（减少无用请求） 2024-8-27
// 优先验证本平台的账号密码，如果匹配成功则结束
// 如果密码不存在，证明旧版本没有保存，将走授权平台并获得最新的密码进行保存
// 如果账号无法再本平台找到，将直接走授权平台进行登录，成功后在本平台创建账号信息
func (a *Api) Login(c *gin.Context) {
	req := types.LoginReq{}

	if err := base.GetRequestParam(c, &req); err != nil {
		global.Logger.Error(apiReturn.ErrorWithNumReturnErrStr(c, err))
		return
	}

	logReq := req
	logReq.Password = "***"
	global.Logger.Info(a.log("Login", c.ClientIP(), logReq))

	if !a.isCheckPassedClientIdExist(c, req.ClientId) {
		return
	}

	// 黑名单IP，返回404
	if biz.ClientCache.BlacklistIP.CheckAndUpdate(c.ClientIP()) {
		global.Logger.Info("Blacklist IP:", c.ClientIP())
		c.JSON(404, gin.H{})
		return
	}

	// 账号或密码为空
	if req.Password == "" || req.Username == "" {
		apiReturn.ErrorByCode(c, apiReturn.ErrAccountNotExist)
		return
	}

	clientiIp := c.ClientIP() + "_" + req.LanIP // 改为外网IP + 内网IP 组合 - 2024-9-18 20:38:28

	// 查询ip是否为尝试次数过多的封禁状态
	if num, ok := global.ClientLoginAttemptsCacheCache.Get(clientiIp); ok && num > 5 {
		apiReturn.ErrorByCode(c, apiReturn.ErrLoginTooMuchFail)
		return
	}

	// 校验时间戳
	if ok := base.CheckTimestamp(req.Timestamp); !ok {
		global.Logger.Errorln(apiReturn.ErrorWithNumReturnErrStr(c, errors.New("abnormal timestamp")), "IP:", c.ClientIP(), "param:", cmn.AnyToJsonStr(req))
		return
	}

	// 账号密码验证
	mUser := models.User{}
	loginUserInfo, err := mUser.GetUserInfoByUsername(req.Username)
	userInfo := loginUserInfo

	if err != nil {
		// 没有账号记录，将走授权平台登录
		if !errors.Is(gorm.ErrRecordNotFound, err) {
			apiReturn.Error(c, "server error")
			return
		}

		u, ok := a.authProformLogin(c, req)
		if !ok {
			a.clicentIPLoginAttemptsAdd(clientiIp, 1)
			// authProformLogin 直接返回，无需再返回
			// apiReturn.ErrorByCode(c, apiReturn.ErrAccountNotExist)
			return
		}
		userInfo = u
	} else if userInfo.Password == "" {
		// 适配旧版密码为空的问题
		u, ok := a.authProformLogin(c, req)
		if !ok {
			a.clicentIPLoginAttemptsAdd(clientiIp, 1)
			// authProformLogin 直接返回，无需再返回
			// apiReturn.ErrorByCode(c, apiReturn.ErrAccountNotExist)
			return
		}
		userInfo = u
	} else {
		// 停用或未激活或者密码错误
		if loginUserInfo.Status != 1 {
			a.clicentIPLoginAttemptsAdd(clientiIp, 1)
			apiReturn.ErrorByCode(c, apiReturn.ErrAccountNotExist)
			return
		}

		// 密码错误
		if loginUserInfo.Password != cmn.PasswordEncryption(req.Password) {
			a.clicentIPLoginAttemptsAdd(clientiIp, 1)
			apiReturn.ErrorByCode(c, apiReturn.ErrIncorrectUsernameOrPassword)
			return
		}
	}

	// 登录成功，获取授权信息
	proExpiration, _ := base.GetProAuthExpiration(userInfo.ID)

	// 更新客户端信息到数据库
	clientInfo := models.SoftwareClient{}
	clientInfo.Version = req.Version
	clientInfo.UserId = userInfo.ID
	clientInfo.MacAddress = req.MacAddress
	clientInfo.LastLanIp = req.LanIP
	clientInfo.LastIp = c.ClientIP()
	clientInfo.JoinTime = time.Now()
	clientInfo.LastOnlineTime = time.Now()
	clientInfo.ClientId = *req.ClientId

	// ==================== 更新客户端信息 废弃 2024年12月30日
	// updateFields := []string{
	// 	"Version",
	// 	"MacAddress",
	// 	"LastIp",
	// 	"LastLanIp",
	// 	"UserId",
	// }
	// if err := global.Db.Model(&models.SoftwareClient{}).Where("client_id=?", req.ClientId).Select(updateFields).Updates(clientInfo).Error; err != nil {
	// 	// 客户端不存在，自动创建客户端
	// 	if errors.Is(err, gorm.ErrRecordNotFound) {
	// 		// 客户端id不存在、为空、长度大于40全部不允许登录
	// 		if req.ClientId == nil || *req.ClientId == "" || len(*req.ClientId) > 40 {
	// 			global.Logger.Error(apiReturn.ErrorWithNumReturnErrStr(c, errors.New("clientID incorrect")))
	// 			return
	// 		}
	// 		clientInfo.ClientId = *req.ClientId
	// 		global.Db.Create(&clientInfo)
	// 		global.Logger.Infoln("client id not exist，auto create:", cmn.AnyToJsonStr(clientInfo))
	// 	} else {
	// 		global.Logger.Error(apiReturn.ErrorWithNumReturnErrStr(c, err))
	// 		return
	// 	}
	// }
	// ======================

	// 全新使用缓存逻辑 2024年12月30日
	if err := a.checkRequestClientIdUpdateAndCreate(&clientInfo.ClientId, clientInfo); err != nil {
		global.Logger.Error(apiReturn.ErrorWithNumReturnErrStr(c, err))
		return
	}

	// 更新并返回token
	cToken := uuid.NewString() + "-" + cmn.Md5(cmn.Md5("userId"+strconv.Itoa(int(userInfo.ID))))

	// 改为通用业务处理 - 简单化 - 2024-9-19
	biz.ClientCache.AccountOnline.AddCtokenByUserInfo(cToken, userInfo)
	biz.ClientCache.AccountOnline.SetClient(userInfo.ID, *req.ClientId, clientCache.AccountOnlineCacheInfo{
		Ctoken:    cToken,
		IP:        c.ClientIP(),
		Timestamp: time.Now().Unix(),
		LanIP:     req.LanIP,
		Mac:       req.MacAddress,
	})

	global.Logger.Debugln("返回数据", types.LoginResp{
		Token:         cToken,
		ProExpiration: proExpiration,
		RespBase:      base.GetRespBase(cToken),
		UserInfo: types.UserInfo{
			Username: userInfo.Username,
			Name:     userInfo.Name,
		},
	})

	apiReturn.SuccessDataDw(c, types.LoginResp{
		Token:         cToken,
		ProExpiration: proExpiration,
		RespBase:      base.GetRespBase(cToken),
		UserInfo: types.UserInfo{
			Username: userInfo.Username,
			Name:     userInfo.Name,
		},
	})
}

// 自动登录
// 自动登录采用账号和密码登录，但是对比登录接口做了相关限制，避免频繁迁移服务器导致缓存丢失，从而导致用户需要重新授权登录。所以没有考虑refreshToken的方式
func (a *Api) AutoLogin(c *gin.Context) {
	req := types.LoginReq{}

	if err := base.GetRequestParam(c, &req); err != nil {
		global.Logger.Error(apiReturn.ErrorWithNumReturnErrStr(c, err))
		return
	}

	// 校验时间戳
	if ok := base.CheckTimestamp(req.Timestamp); !ok {
		global.Logger.Error(apiReturn.ErrorWithNumReturnErrStr(c, errors.New("abnormal timestamp")), "IP:", c.ClientIP(), "param:", cmn.AnyToJsonStr(req))
		return
	}

	logReq := req
	logReq.Password = "***"
	global.Logger.Info(a.log("AutoLogin", c.ClientIP(), logReq))

	mUser := models.User{}
	userInfo, err := mUser.GetUserInfoByUsername(req.Username)

	// 密码错误
	if userInfo.Password != cmn.PasswordEncryption(req.Password) {
		apiReturn.ErrorByCode(c, apiReturn.ErrIncorrectUsernameOrPassword)
		return
	}

	// 停用或未激活或者密码错误
	if userInfo.Status != 1 || err != nil {
		// a.clicentIPLoginAttemptsAdd(clicentIP, 1)
		apiReturn.ErrorByCode(c, apiReturn.ErrAccountNotExist)
		return
	}

	if !a.isCheckPassedClientIdExist(c, req.ClientId) {
		return
	}

	// 获取授权信息
	proExpiration, _ := base.GetProAuthExpiration(userInfo.ID)

	// 获取上次客户端的信息，此方式影响集成镜像导致出问题
	// 将使用账号+ip的方式来判断是否超出范围

	// lastClientInfo := models.SoftwareClient{}
	// if err := global.Db.Where("client_id=?", req.ClientId).Order("updated_at Desc").First(&lastClientInfo).Error; err != nil {
	// 	global.Logger.Errorln("find client id fail, request:", cmn.AnyToJsonStr(req))
	// 	apiReturn.ErrorByCode(c, apiReturn.ErrManualLogin)
	// 	return
	// }

	clientInfo := models.SoftwareClient{}
	clientInfo.Version = req.Version
	clientInfo.UserId = userInfo.ID
	clientInfo.MacAddress = req.MacAddress
	clientInfo.LastLanIp = req.LanIP
	clientInfo.LastIp = c.ClientIP()
	clientInfo.ClientId = *req.ClientId
	clientInfo.LastOnlineTime = time.Now()
	clientInfo.JoinTime = time.Now()

	// ==================更新客户端信息 废弃 采用全新使用缓存逻辑 2024年12月30日
	// updateFields := []string{
	// 	"Version",
	// 	"LastIp",
	// 	"LastLanIp",
	// }
	// if err := global.Db.Model(&models.SoftwareClient{}).Where("client_id=?", req.ClientId).Select(updateFields).Updates(clientInfo).Error; err != nil {
	// 	global.Logger.Errorln("upload client info fail.", "param:", cmn.AnyToJsonStr(clientInfo), ",err:", err)
	// 	apiReturn.ErrorByCode(c, apiReturn.ErrManualLogin)
	// 	return
	// }
	// ==================

	// 全新使用缓存逻辑 2024年12月30日
	biz.ClientCache.HistoryCache.WriteToHourUpdateClientsCache(clientInfo.ClientId, clientInfo)

	// 更新并返回token
	cToken := uuid.NewString() + "-" + cmn.Md5(cmn.Md5("userId"+strconv.Itoa(int(userInfo.ID))))

	// 改为通用业务处理 - 简单化 - 2024-9-19

	getClientInfoAndNotExistCreateOption := GetClientInfoAndNotExistCreateOptions{
		ClientId: *req.ClientId,
		UsserId:  userInfo.ID,
		NewClientInfo: clientCache.AccountOnlineCacheInfo{
			Ctoken:    cToken,
			IP:        c.ClientIP(),
			Timestamp: time.Now().Unix(),
			LanIP:     req.LanIP,
			ClientID:  *req.ClientId,
			Mac:       req.MacAddress,
		},
		FnName: `AutoLogin`,
	}

	// 获取账号+客户端缓存信息
	oldClientInfo, oldClientInfoOk := a.getClientInfoAndNotExistCreate(getClientInfoAndNotExistCreateOption)
	if !oldClientInfoOk {
		biz.ClientCache.AccountOnline.DeleteClient(userInfo.ID, *req.ClientId)
		apiReturn.ErrorByCode(c, apiReturn.ErrManualLogin)
		return
	}

	// token为空 || mac地址不一致，需要重新授权
	if oldClientInfo.Ctoken == "" || req.MacAddress != oldClientInfo.Mac {
		apiReturn.ErrorByCode(c, apiReturn.ErrManualLogin)
		return
	}

	// 删除旧的token
	if oldClientInfo.Ctoken != cToken {
		biz.ClientCache.AccountOnline.DeleteCtoken(oldClientInfo.Ctoken)
	}
	biz.ClientCache.AccountOnline.AddCtokenByUserInfo(cToken, userInfo)
	biz.ClientCache.AccountOnline.SetClient(userInfo.ID, *req.ClientId, getClientInfoAndNotExistCreateOption.NewClientInfo)

	apiReturn.SuccessDataDw(c, types.LoginResp{
		Token:         cToken,
		ProExpiration: proExpiration,
		RespBase:      base.GetRespBase(cToken),
		UserInfo: types.UserInfo{
			Username: userInfo.Username,
			Name:     userInfo.Name,
		},
	})
}

// 刷新信息
func (a *Api) RefreshInfo(c *gin.Context) {
	req := types.RefreshInfoReq{}
	if err := base.GetRequestParam(c, &req); err != nil {
		global.Logger.Error(apiReturn.ErrorWithNumReturnErrStr(c, err))
		return
	}

	// 校验时间戳
	if ok := base.CheckTimestamp(req.Timestamp); !ok {
		global.Logger.Errorln(apiReturn.ErrorWithNumReturnErrStr(c, errors.New("abnormal timestamp")), "IP:", c.ClientIP(), "param:", cmn.AnyToJsonStr(req))
		return
	}

	global.Logger.Info(a.log("RefreshInfo", c.ClientIP(), req))
	if !a.isCheckPassedClientIdExist(c, req.ClientId) {
		return
	}

	var (
		proExpiration *time.Time
		resp          types.RefreshInfoResp
		userInfo      models.User
	)

	if v, err := base.GetUserInfoWithCheckByCToken(c, req.Token); err != nil {
		global.Logger.Debugln("token expires 1", err)
		apiReturn.ErrorByCode(c, apiReturn.ErrTokenExpires)
		return
	} else {
		userInfo = *v
		resp.UserInfo.Username = v.Username
		resp.UserInfo.Name = v.Name

		mUser := models.User{}
		u, _ := mUser.GetUserInfoByUid(v.ID)
		if u.Name != "" {
			resp.UserInfo.Name = u.Name
		}

		if u.Username != "" {
			resp.UserInfo.Username = u.Username
		}

		// global.Logger.Debugln("userInfo", resp)
		proExpiration, _ = base.GetProAuthExpiration(v.ID)
	}

	getClientInfoAndNotExistCreateOption := GetClientInfoAndNotExistCreateOptions{
		ClientId: *req.ClientId,
		UsserId:  userInfo.ID,
		NewClientInfo: clientCache.AccountOnlineCacheInfo{
			Ctoken:    req.Token,
			IP:        c.ClientIP(),
			Timestamp: time.Now().Unix(),
			LanIP:     req.LanIP,
			ClientID:  *req.ClientId,
			Mac:       req.MacAddress,
		},
		FnName: `false`,
	}
	// 获取账号+客户端缓存信息
	oldClientInfo, oldClientInfoOk := a.getClientInfoAndNotExistCreate(getClientInfoAndNotExistCreateOption)
	if !oldClientInfoOk {
		global.Logger.Debugln("token expires 2")
		apiReturn.ErrorByCode(c, apiReturn.ErrTokenExpires)
		return
	}

	if !a.isCheckPassedcTokenAndMacChamge(c, userInfo.ID, *req.ClientId, req.Token, req.MacAddress, oldClientInfo) {
		global.Logger.Debugln("token expires 3")
		return
	}

	resp.ProExpiration = proExpiration
	resp.RespBase = base.GetRespBase(req.Token)

	apiReturn.SuccessDataDw(c, resp)
}

// 记录设备的在线状态，如果为登录状态还会更新授权时间
func (a *Api) Ping(c *gin.Context) {
	req := types.PingReq{}
	if err := base.GetRequestParam(c, &req); err != nil {
		global.Logger.Error(apiReturn.ErrorWithNumReturnErrStr(c, err))
		return
	}

	// 校验时间戳
	if ok := base.CheckTimestamp(req.Timestamp); !ok {
		global.Logger.Errorln(apiReturn.ErrorWithNumReturnErrStr(c, errors.New("abnormal timestamp")), "IP:", c.ClientIP(), "param:", cmn.AnyToJsonStr(req))
		return
	}

	// 降低Ping请求的日志级别，避免高频日志写入消耗CPU
	// global.Logger.Info(a.log("Ping", c.ClientIP(), req))

	if !a.isCheckPassedClientIdExist(c, req.ClientId) {
		return
	}

	var proExpiration *time.Time
	clientInfo := models.SoftwareClient{}
	clientInfo.Version = req.Version
	clientInfo.MacAddress = req.MacAddress
	clientInfo.LastLanIp = req.LanIP
	clientInfo.LastIp = c.ClientIP()
	clientInfo.ClientId = *req.ClientId
	clientInfo.LastOnlineTime = time.Now()
	clientInfo.JoinTime = time.Now()

	// 检查频率，规定范围内的ping请求将被忽略
	{
		if biz.ClientCache.PingRateCache.GetLastPingRechord(*req.ClientId, req.LanIP) {
			// 超出设置频率，忽略
			// global.Logger.Debugln("超出设置频率，忽略")
			return
		}

		// 记录ping的频率
		biz.ClientCache.PingRateCache.SetLastPingRechord(*req.ClientId, req.LanIP)
		// global.Logger.Debugln("记录ping的频率")
	}

	// ======================= 更新客户端信息 废弃 采用全新使用缓存逻辑 2024年12月30日
	// updateFields := []string{
	// 	"Version",
	// 	"LastIp",
	// 	"LastLanIp",
	// }
	// // 更新数据库中的客户端信息
	// if err := global.Db.Model(&models.SoftwareClient{}).Where("client_id=?", req.ClientId).Select(updateFields).Updates(clientInfo).Error; err != nil {
	// 	global.Logger.Errorln(err)
	// }
	// =======================

	// 全新使用缓存逻辑 2024年12月30日
	biz.ClientCache.HistoryCache.WriteToHourUpdateClientsCache(clientInfo.ClientId, clientInfo)

	// token 存在，可能代表已经登录了用户
	if req.Token != "" {
		userInfo := models.User{}
		// 获取cToken的用户信息，如果不存在返回token过期，否则获取授权到期时间
		if v, err := base.GetUserInfoWithCheckByCToken(c, req.Token); err != nil {
			apiReturn.ErrorByCode(c, apiReturn.ErrTokenExpires)
			return
		} else {
			userInfo = *v
			proExpiration, _ = base.GetProAuthExpiration(v.ID)

			clientInfo.UserId = userInfo.ID

			getClientInfoAndNotExistCreateOption := GetClientInfoAndNotExistCreateOptions{
				ClientId: *req.ClientId,
				UsserId:  userInfo.ID,
				NewClientInfo: clientCache.AccountOnlineCacheInfo{
					Ctoken:    req.Token,
					IP:        c.ClientIP(),
					Timestamp: time.Now().Unix(),
					LanIP:     req.LanIP,
					ClientID:  *req.ClientId,
					Mac:       req.MacAddress,
				},
				FnName: `PING`,
			}

			// 获取账号+客户端缓存信息
			oldClientInfo, oldClientInfoOk := a.getClientInfoAndNotExistCreate(getClientInfoAndNotExistCreateOption)
			if !oldClientInfoOk {
				// global.Logger.Debug("过期...")
				apiReturn.ErrorByCode(c, apiReturn.ErrTokenExpires)
				return
			}

			if !a.isCheckPassedcTokenAndMacChamge(c, userInfo.ID, *req.ClientId, req.Token, req.MacAddress, oldClientInfo) {
				// global.Logger.Debug("过期...111")
				return
			}
		}
	}

	// 已在前面合并更新，删除重复的缓存更新操作
	// biz.ClientCache.HistoryCache.WriteToHourUpdateClientsCache(clientInfo.ClientId, clientInfo)

	// 这里后期可以判断版本增加修改客户端的 ClientID 功能

	resp := types.PingResp{}
	resp.ProExpiration = proExpiration
	resp.RespBase = base.GetRespBase(req.Token)

	apiReturn.SuccessDataDw(c, resp)
}

// 续期短期授权的授权码
func (a *Api) RenewTempAuth(c *gin.Context) {
	req := types.PingReq{}
	if err := base.GetRequestParam(c, &req); err != nil {
		global.Logger.Error(apiReturn.ErrorWithNumReturnErrStr(c, err))
		return
	}

	// 校验时间戳
	if ok := base.CheckTimestamp(req.Timestamp); !ok {
		global.Logger.Errorln(apiReturn.ErrorWithNumReturnErrStr(c, errors.New("abnormal timestamp")), "IP:", c.ClientIP(), "param:", cmn.AnyToJsonStr(req))
		return
	}

	global.Logger.Info(a.log("RenewTempAuth", c.ClientIP(), req))

	if !a.isCheckPassedClientIdExist(c, req.ClientId) {
		return
	}

	var (
		proExpiration *time.Time
		userInfo      models.User
	)

	if v, err := base.GetUserInfoWithCheckByCToken(c, req.Token); err != nil {
		apiReturn.ErrorByCode(c, apiReturn.ErrTokenExpires)
		return
	} else {
		// global.Logger.Debugln("userInfo", v)
		userInfo = *v
		proExpiration, _ = base.GetProAuthExpiration(v.ID)
	}

	getClientInfoAndNotExistCreateOption := GetClientInfoAndNotExistCreateOptions{
		ClientId: *req.ClientId,
		UsserId:  userInfo.ID,
		NewClientInfo: clientCache.AccountOnlineCacheInfo{
			Ctoken:    req.Token,
			IP:        c.ClientIP(),
			Timestamp: time.Now().Unix(),
			LanIP:     req.LanIP,
			ClientID:  *req.ClientId,
			Mac:       req.MacAddress,
		},
		FnName: `RenewTempAuth`,
	}
	// 获取账号+客户端缓存信息
	oldClientInfo, oldClientInfoOk := a.getClientInfoAndNotExistCreate(getClientInfoAndNotExistCreateOption)
	if !oldClientInfoOk {
		apiReturn.ErrorByCode(c, apiReturn.ErrTokenExpires)
		return
	}

	if !a.isCheckPassedcTokenAndMacChamge(c, userInfo.ID, *req.ClientId, req.Token, req.MacAddress, oldClientInfo) {
		return
	}

	var tempExpiration string
	// 如果有过期时间并在今天之后，并且过期时间差小于七天
	// 将临时授权设置为pro的过期时间
	if proExpiration != nil && proExpiration.After(time.Now()) {
		if !base.IsTimeDifferenceGreaterThan(7, time.Now(), *proExpiration) {
			tempExpiration = proExpiration.Format(cmn.TimeFormatMode1)
		} else {
			tempExpiration = time.Now().Add(7 * 24 * time.Hour).Format(cmn.TimeFormatMode1)
		}
	}

	resp := types.RenewTempAuthResp{}
	resp.Timestamp = time.Now().Unix()
	resp.TempExpiration = tempExpiration
	resp.Token = req.Token
	resp.ProExpiration = proExpiration
	resp.RespBase = base.GetRespBase(req.Token)

	apiReturn.SuccessDataDw(c, resp)
}

// 检查更新版本
func (a *Api) CheckVersion(c *gin.Context) {
	req := types.CheckVersionReq{}
	resp := types.CheckVersionResp{}
	if err := base.GetRequestParam(c, &req); err != nil {
		global.Logger.Error(apiReturn.ErrorWithNumReturnErrStr(c, err))
		return
	}

	// 校验时间戳
	if ok := base.CheckTimestamp(req.Timestamp); !ok {
		global.Logger.Errorln(apiReturn.ErrorWithNumReturnErrStr(c, errors.New("abnormal timestamp")), "IP:", c.ClientIP(), "param:", cmn.AnyToJsonStr(req))
		return
	}

	// if req.ClientId == nil {
	// 	*req.ClientId = ""
	// }
	// global.Logger.Info(a.log("CheckVersion", c.ClientIP(), *req.ClientId, req))

	version, err := biz.Version.GetCurrentVersion(req.Version)

	versionType := req.VersionType
	// 如果版本号中包含 -dev，则将 VersionType 改为 dev
	if strings.Contains(req.Version, "-dev") {
		versionType = models.VersionTypeDev
	}

	var latestVersion *models.Version
	if err != nil {
		v, err := biz.Version.GetLatest(models.VersionType(versionType))
		if err != nil {
			if !errors.Is(gorm.ErrRecordNotFound, err) {
				global.Logger.Error(apiReturn.ErrorWithNumReturnErrStr(c, err))
				return
			}
		} else {
			latestVersion = &v
		}

	} else {
		v, err := biz.Version.GetLatestByAfterTime(version.ReleaseTime, models.VersionType(versionType))
		if err != nil {
			if !errors.Is(gorm.ErrRecordNotFound, err) {
				global.Logger.Error(apiReturn.ErrorWithNumReturnErrStr(c, err))
				return
			}
		} else {
			latestVersion = &v
		}

	}

	if latestVersion != nil {
		versionInfo := types.VersionInfo{}
		versionInfo.Version = latestVersion.Version
		versionInfo.ReleaseTime = latestVersion.ReleaseTime
		versionInfo.PageUrl = latestVersion.PageUrl
		versionInfo.IsActive = latestVersion.IsActive
		versionInfo.DownloadURL = latestVersion.DownloadURL
		versionInfo.IsRolledBack = latestVersion.IsRolledBack
		versionInfo.Type = latestVersion.Type
		resp.VersionInfo = &versionInfo
	}

	resp.Timestamp = time.Now().Unix()

	// jb, _ := json.Marshal(resp)
	// global.Logger.Debug("查询新版本数据-结果", string(jb))

	// jbrwq, _ := json.Marshal(req)
	// global.Logger.Debug("请求数据", string(jbrwq))

	// respData, respErr := base.EncryptParam(resp)
	// if respErr != nil {
	// 	global.Logger.Error(apiReturn.ErrorWithNumReturnErrStr(c, respErr))
	// 	return
	// }
	apiReturn.SuccessDataDw(c, resp)
}

// 授权平台登录
func (a *Api) authProformLogin(c *gin.Context, req types.LoginReq) (models.User, bool) {
	thirdAppClientId, clientSecret := biz.SunStore.GetClientIdAndSecret()
	apiHost := biz.SunStore.ApiHost()

	sunapi := sunStore.NewSunStoreApi(apiHost, thirdAppClientId, clientSecret)
	accessToken, err := sunapi.PasswordAuth(sunStore.PasswordAuthRequest{
		Username: req.Username,
		Password: req.Password,
	})

	if err != nil {
		// 服务器错误
		if errors.Is(err, sunStore.ErrorServer) {
			global.Logger.Errorln("auth Api login failed", err.Error())
			apiReturn.Error(c, "server error")
			return models.User{}, false
		}

		// a.clicentIPLoginAttemptsAdd(c.ClientIP(), 1)
		apiReturn.ErrorByCode(c, apiReturn.ErrAccountNotExist)
		global.Logger.Errorln("sunapi.PasswordAut", err)
		return models.User{}, false
	}

	global.Logger.Infoln("client user login", "IP:", c.ClientIP(), "clientId:", a.clientIdToStr(req.ClientId), "username:", req.Username)

	openUser, err := biz.SunStore.GetMainPlatformUserInfo(apiHost, accessToken.AccessToken)
	if err != nil {
		// 可能过期了,暂时返回账号过密码错误
		global.Logger.Errorln("authProformLogin：biz.SunStore.GetMainPlatformUserInfo ", err)
		apiReturn.ErrorByCode(c, apiReturn.ErrIncorrectUsernameOrPassword)
		return models.User{}, false
	}

	global.Logger.Infoln("client user login,openuser:", openUser)

	userInfo := models.User{}
	// 查询主平台邮箱是否在本平台有同邮箱账号
	// 不存在创建新账号，否则行绑定
	if err := global.Db.First(&userInfo, "username=?", openUser.Mail).Error; err != nil {
		// 创建新账号
		userInfo.Username = openUser.Username
		userInfo.Mail = openUser.Mail
		userInfo.Name = openUser.Name
		userInfo.Role = 2
		userInfo.Status = 1
		userInfo.Lang = openUser.Lang
		userInfo.SystemLang = openUser.SystemLang
		userInfo.TimeZone = openUser.TimeZone
		userInfo.Password = openUser.Password

		// 无法创建平台账号
		if errCreate := global.Db.Create(&userInfo).Error; errCreate != nil {
			global.Logger.Errorln("Unable to create platform account", apiReturn.ErrorWithNumReturnErrStr(c, err))
			apiReturn.ErrorByCode(c, apiReturn.ErrIncorrectUsernameOrPassword)
			return models.User{}, false
		}
	}

	// 绑定主平台账号
	mUser := models.User{}
	mUser.UpdateUserInfoByUserId(userInfo.ID, map[string]interface{}{
		"store_email": openUser.Mail,
		"lang":        openUser.Lang,
		"system_lang": openUser.SystemLang,
		"time_zone":   openUser.TimeZone,
		"password":    openUser.Password,
	})

	return userInfo, true
}

// 增加登录尝试次数
func (a *Api) clicentIPLoginAttemptsAdd(clicentIP string, num int) {
	// 将登录尝试错误次数+num
	timesNum := 1
	if num, ok := global.ClientLoginAttemptsCacheCache.Get(clicentIP); ok {
		timesNum = num + num
	}

	// global.Logger.Debugln("login failed ip:", clicentIP, "times:", timesNum)
	global.ClientLoginAttemptsCacheCache.SetDefault(clicentIP, timesNum)
}

func (a *Api) clientIdToStr(clientId *string) string {

	if clientId == nil {
		return ""
	}
	return *clientId
}

func (a *Api) log(funcName, ip, reqData interface{}) string {
	content := fmt.Sprintf(
		"Client API - FuncName:%s, IP:%s, ReqData:%s",
		funcName, ip, cmn.AnyToJsonStr(reqData))

	return content
}

// 检查MAC地址是否与上次一致
func (a *Api) isCheckPassedMacAddressChange(funcName, ip, reqData, respData interface{}) string {
	content := fmt.Sprintf(
		"Client API - FuncName:%s, IP:%s, ReqData:%s, RespData:%s",
		funcName, ip, cmn.AnyToJsonStr(reqData), cmn.AnyToJsonStr(respData))

	return content
}

// clientId 不存在或者为空返回404
func (a *Api) isCheckPassedClientIdExist(c *gin.Context, clientId *string) bool {
	if clientId == nil || *clientId == "" {
		c.JSON(http.StatusNotFound, gin.H{"error": "clientId is nil"})
		return false
	}

	return true

}

type GetClientInfoAndNotExistCreateOptions struct {
	UsserId       uint
	ClientId      string
	NewClientInfo clientCache.AccountOnlineCacheInfo
	FnName        string // 函数名称
	CToken        string
}

// 获取账号+客户端缓存信息，如果不存在根据配置判断调用方是否创建
func (a *Api) getClientInfoAndNotExistCreate(options GetClientInfoAndNotExistCreateOptions) (clientCache.AccountOnlineCacheInfo, bool) {
	// global.Logger.Debugln("getClientInfoAndNotExistCreate", cmn.AnyToJsonStr(options))

	// 获取账号+客户端缓存信息
	clientInfo, clientInfoOk := biz.ClientCache.AccountOnline.GetClient(options.UsserId, options.ClientId)

	if !clientInfoOk {
		createOnlineConfig := systemSetting.ClientCreateOnlineClientCache{}
		global.SystemSetting.GetValueByInterface(systemSetting.ClientCreateOnlineClientCacheKey, &createOnlineConfig)

		if global.RUNCODE == "debug" {
			global.Logger.Debugln("读取配置看是否需要创建"+systemSetting.ClientCreateOnlineClientCacheKey, cmn.AnyToJsonStr(createOnlineConfig))
		}

		fnName := strings.ToLower(options.FnName)
		switch {
		case fnName == "ping" && createOnlineConfig.Ping:
			fallthrough
		case fnName == "renewtempauth" && createOnlineConfig.RenewTempAuth:
			fallthrough
		case fnName == "autologin" && createOnlineConfig.AutoLogin:
			// global.Logger.Debug("创建缓存-", fnName)
			clientInfo = options.NewClientInfo
			biz.ClientCache.AccountOnline.SetClient(options.UsserId, options.ClientId, options.NewClientInfo)
		default:
			return clientInfo, false
		}
	}

	return clientInfo, true
}

// 验证 token / MAC 是否与缓存一致，不一致将token置空
func (a *Api) isCheckPassedcTokenAndMacChamge(c *gin.Context, userId uint, clientId, cToken, mac string, oldClientInfo clientCache.AccountOnlineCacheInfo) bool {
	// global.Logger.Debugln("isCheckPassedcTokenAndMacChamge", "cToken:", cToken, "mac:", mac)
	// global.Logger.Debugln("oldClientInfo", "cToken:", oldClientInfo.Ctoken, "mac:", oldClientInfo.Mac)
	// 验证 token / MAC 是否与缓存一致，不一致将token置空
	if oldClientInfo.Ctoken != cToken || oldClientInfo.Mac != mac {
		biz.ClientCache.AccountOnline.DeleteCtoken(oldClientInfo.Ctoken) // 删除cToken 缓存
		oldClientInfo.Ctoken = ""
		biz.ClientCache.AccountOnline.SetClient(userId, clientId, oldClientInfo)
		apiReturn.ErrorByCode(c, apiReturn.ErrTokenExpires)
		return false
	}

	return true
}

// 检查clientID是否存在，优先查询缓存，不存在再查询数据库
func (a *Api) clientIDIsExist(clientId string) (bool, error) {
	var exist bool = false
	exist = biz.ClientCache.HistoryCache.ClientIdInNewClientCache(clientId)

	// 查询数据库
	if !exist {
		if ok, err := models.IsClientIdExist(global.Db, clientId); err != nil {
			return false, err
		} else {
			exist = ok
		}
	}

	return exist, nil
}

// 检查客户端传来的clientId是否合法，合法如果不存在则创建 - 2024年12月30日
func (a *Api) checkRequestClientIdAndCreate(requestClientId *string, clientInfo models.SoftwareClient) error {
	// 客户端id不存在、为空、长度大于40全部不允许登录
	if requestClientId == nil || *requestClientId == "" || len(*requestClientId) > 40 {
		// global.Logger.Error(apiReturn.ErrorWithNumReturnErrStr(c, errors.New("clientID incorrect")))
		return errors.New("clientID incorrect")
	}

	clientInfo.ClientId = *requestClientId

	if ok, err := a.clientIDIsExist(*requestClientId); err != nil {
		return err
	} else if !ok {
		// 创建客户端信息
		// 写入到缓存中，定时任务同步客户端数据库 - 2024年12月30日
		biz.ClientCache.HistoryCache.WriteToNewClientsCache(*requestClientId, clientInfo)
		global.Logger.Infoln("client id not exist，auto create:", cmn.AnyToJsonStr(clientInfo))
	}

	return nil
}

// 检查客户端传来的clientId是否合法，合法并判断是否存在，合法如果不存在则创建 - 2024年12月30日
// 如果客户端id不存在创建缓存和数据库中，将合法的clientId直接写入新客户端缓存中，存在的话写入到更新缓存中
func (a *Api) checkRequestClientIdUpdateAndCreate(requestClientId *string, clientInfo models.SoftwareClient) error {
	// 客户端id不存在、为空、长度大于40全部不允许登录
	if requestClientId == nil || *requestClientId == "" || len(*requestClientId) > 40 {
		// global.Logger.Error(apiReturn.ErrorWithNumReturnErrStr(c, errors.New("clientID incorrect")))
		return errors.New("clientID incorrect")
	}

	clientInfo.ClientId = *requestClientId

	if ok, err := a.clientIDIsExist(*requestClientId); err != nil {
		return err
	} else if ok {
		biz.ClientCache.HistoryCache.WriteToHourUpdateClientsCache(*requestClientId, clientInfo)
	} else if !ok {
		// 创建客户端信息
		// 写入到缓存中，定时任务同步客户端数据库 - 2024年12月30日
		biz.ClientCache.HistoryCache.WriteToNewClientsCache(*requestClientId, clientInfo)
		global.Logger.Infoln("client id not exist，auto create:", cmn.AnyToJsonStr(clientInfo))
	}

	return nil
}
