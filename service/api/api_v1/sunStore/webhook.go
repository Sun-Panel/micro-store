package sunStore

import (
	"errors"
	"fmt"
	"net/http"
	"sun-panel/biz"
	"sun-panel/global"
	"sun-panel/lib/cmn"
	"sun-panel/lib/sunStore/openApi/clientApi"
	"sun-panel/lib/sunStore/webhook"
	"sun-panel/models"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"gorm.io/gorm"
)

type GoodsOrder struct {
}

func (a *GoodsOrder) Enter(c *gin.Context) {
	key := global.Config.GetValueString("sun_store_webhook_secret_key", "goods_order")
	if ok := webhook.VerifySignature(c.Request, key); !ok {
		c.String(http.StatusUnauthorized, "") // 401
		return
	}

	req := webhook.EventGoodsOrderDataReq{}
	if err := c.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		c.String(http.StatusInternalServerError, "")
		return
	}

	if req.Event != "goodsOrder.status" || req.Data.Status != "PAY_SUCCESS" {
		return
	}

	// 成功支付
	var dayNum float64
	if v, ok := req.Data.GoodsSnapshot.CustomData["day"]; ok {
		if day, intOk := v.(float64); intOk {
			dayNum = day
		}
		// 查找用户
		mUser := models.User{}
		userInfo, err := mUser.GetUserInfoByUsername(req.Data.User.Username)
		if err != nil {
			// 查找用户异常
			global.Logger.Errorln("error: recharge failed. find user error:", err)
			return
		}

		countDays := int(dayNum) * req.Data.Number

		global.Logger.Infoln(fmt.Sprintf("order finish orderNo:%s status:%s username:%s AuthDays:%d",
			req.Data.OrderNo, "FINISH", userInfo.Username, countDays))

		// 正式充值PRO授权天数
		// 查询是否已存在该订单记录，如果已存在将不再重复充值
		_, err = biz.ProAuthorize.GetChangeRecordByOrderNo(req.Data.OrderNo)
		if err != nil {
			err = biz.ProAuthorize.ChangeExpiredTimeByDayNum(userInfo.ID, countDays, "", "", req.Data.OrderNo)
			if err != nil {
				global.Logger.Errorln("error: recharge failed. biz.ProAuthorize.ChangeExpiredTimeByDayNum:", err)
				c.String(http.StatusInternalServerError, "")
			}
		}

		// 在此处调用SunStore的API，来修改订单状态
		clientId, clientSecret := biz.SunStore.GetClientIdAndSecret()
		host := biz.SunStore.ApiHost()
		tk, err := biz.SunStore.GetClientApiToken(biz.SunStore.ApiHost(), clientId, clientSecret)
		if err != nil {
			global.Logger.Errorln("error: recharge failed. find user error:", err)
			return
		}
		cApi := clientApi.NewGoodsOrder(clientApi.NewClientAPI(host, tk))
		err = cApi.UpdateStatus(clientApi.GoodsOrderUpdateAtatusReq{
			OrderNo: req.Data.OrderNo,
			Status:  "FINISH",
		})
		if err != nil {
			global.Logger.Errorln("error: recharge failed. API error:", err.Error())
		}

		c.String(http.StatusOK, "")
		return
	}

	c.String(http.StatusOK, "")
}

type User struct {
}

func (a *User) Enter(c *gin.Context) {
	key := global.Config.GetValueString("sun_store_webhook_secret_key", "user")
	if ok := webhook.VerifySignature(c.Request, key); !ok {
		c.String(http.StatusUnauthorized, "") // 401
		return
	}

	req := webhook.EventUserDataReq{}
	if err := c.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		c.String(http.StatusInternalServerError, "")
		return
	}

	var (
		err      error
		userInfo models.User
	)

	// 记录日志
	global.Logger.Infoln("webhook "+req.Event+":", cmn.AnyToJsonStr(req.Data))

	switch req.Event {
	case "user.update":
		fallthrough
	case "user.updatePassword":
		mUser := models.User{}
		userInfo, err = mUser.GetUserInfoByUsername(req.Data.Mail)
		if err != nil {
			global.Logger.Errorln("error:webhook user.update", err)
			c.String(http.StatusInternalServerError, "")
			return
		}
		err = mUser.UpdateUserInfoByUserId(userInfo.ID, map[string]interface{}{
			"head_image":  req.Data.HeadImage,
			"name":        req.Data.Name,
			"lang":        req.Data.Lang,
			"system_lang": req.Data.SystemLang,
			"time_zone":   req.Data.TimeZone,
			"password":    req.Data.Password,
		})
	case "user.create":
		mUser := models.User{}
		userInfo, err = mUser.GetUserInfoByUsername(req.Data.Mail)
		if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
			createUser := models.User{}
			createUser.Mail = req.Data.Mail
			createUser.Username = req.Data.Username
			createUser.Name = req.Data.Name
			createUser.HeadImage = req.Data.HeadImage
			createUser.Status = 1
			createUser.Role = 2
			createUser.Lang = req.Data.Lang
			createUser.SystemLang = req.Data.SystemLang
			createUser.TimeZone = req.Data.TimeZone
			createUser.Password = req.Data.Password
			err = global.Db.Create(&createUser).Error
			if err == nil {
				global.Logger.Infoln("create user:", cmn.AnyToJsonStr(createUser))
			}
		} else {
			err = mUser.UpdateUserInfoByUserId(userInfo.ID, map[string]interface{}{
				"head_image":  req.Data.HeadImage,
				"name":        req.Data.Name,
				"lang":        req.Data.Lang,
				"system_lang": req.Data.SystemLang,
				"time_zone":   req.Data.TimeZone,
				"password":    req.Data.Password,
			})
		}
	default:
		return
	}

	if err != nil {
		global.Logger.Errorln("webhook "+req.Event+":", err)
		c.String(http.StatusInternalServerError, "")
		return
	}

	c.String(http.StatusOK, "")
}
