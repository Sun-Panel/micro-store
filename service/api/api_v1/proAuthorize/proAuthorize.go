package proAuthorize

import (
	"encoding/json"
	"errors"
	"sun-panel/api/api_v1/common/apiData/proAuthorizeApiStructs"
	"sun-panel/api/api_v1/common/apiReturn"
	"sun-panel/api/api_v1/common/base"
	"sun-panel/biz"
	"sun-panel/global"
	"sun-panel/lib/cmn"
	"sun-panel/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ProAuthorizeApi struct {
}

func (a *ProAuthorizeApi) GetAuthorize(c *gin.Context) {
	userInfo, _ := base.GetCurrentUserInfo(c)
	authorize := models.ProAuthorize{}
	resp := proAuthorizeApiStructs.GetAuthorizeResp{}
	if err := global.Db.First(&authorize, "user_id=?", userInfo.ID).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			apiReturn.ErrorDatabase(c, err.Error())
			return
		}
		resp.ExpiredTime = ""
	} else {
		resp.ExpiredTime = authorize.ExpiredTime.Format(cmn.TimeFormatMode1)
	}

	apiReturn.SuccessData(c, resp)

}

func (a *ProAuthorizeApi) GetAuthorizeHistoryRecord(c *gin.Context) {
	userInfo, _ := base.GetCurrentUserInfo(c)
	resp := []proAuthorizeApiStructs.GetAuthorizeHistoryRecordResp{}
	if resList, err := biz.ProAuthorize.GetAuthorizeHistoryRecordByUserId(userInfo.ID); err != nil {
		apiReturn.ErrorDatabase(c, err.Error())
		return
	} else {
		for _, v := range resList {
			resp = append(resp, proAuthorizeApiStructs.GetAuthorizeHistoryRecordResp{
				ChangeTime:  base.ConvertTimeToUserTime(c, v.CreatedAt),
				ExpiredTime: base.ConvertTimeToUserTime(c, v.ExpiredTime),
				DayNum:      v.DayNum,
				Note:        v.Note,
				OrderNo:     v.OrderNo,
			})
		}
	}

	apiReturn.SuccessListData(c, resp, 0)

}

type RedeemCodeInfoGetInfoReq struct {
	Code  string `json:"code"`  // 兑换码
	Vcode string `json:"vcode"` // 验证码
}

type RedeemCodeInfo struct {
	biz.RedeemCodeParam
	Code string `json:"code"` // 兑换码
}

func (a *ProAuthorizeApi) GetRedeemCodeInfo(c *gin.Context) {

	req := RedeemCodeInfoGetInfoReq{}

	if err := c.ShouldBindJSON(&req); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}

	// userInfo, _ := base.GetCurrentUserInfo(c)
	redeemCodeInfo, ok := getRedeemCodeInfoAndErrorReturn(c, req.Code)
	if !ok {
		return
	}

	extendData := map[string]interface{}{}
	json.Unmarshal([]byte(redeemCodeInfo.ExtendData), &extendData)

	respInfo := RedeemCodeInfo{
		Code: redeemCodeInfo.Code,
		RedeemCodeParam: biz.RedeemCodeParam{
			ExpireTime:  base.ConvertTimeToUserTime(c, redeemCodeInfo.ExpireTime).Format("2006-01-02 15:04:05"),
			ReleaseType: redeemCodeInfo.ReleaseType,
			RedeemType:  redeemCodeInfo.RedeemType,
			Status:      redeemCodeInfo.Status,
			Title:       redeemCodeInfo.Title,
			ExtendData:  extendData,
		},
	}

	apiReturn.SuccessData(c, respInfo)

}

// 兑换码扩展数据-PRO授权
type RedeemCodeExtendDataPro struct {
	Days int `json:"days"` // 天数
}

// 核销兑换码
func (a *ProAuthorizeApi) RedeemCodeWriteOff(c *gin.Context) {
	req := RedeemCodeInfoGetInfoReq{}

	if err := c.ShouldBindJSON(&req); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}

	// 验证验证码
	if err := biz.Captcha.CaptchaVerifyHandle(c, req.Vcode, true); err != nil {
		apiReturn.ErrorByCode(c, apiReturn.ErrCodeCaptchaError)
		return
	}

	userInfo, _ := base.GetCurrentUserInfo(c)
	redeemCodeInfo, ok := getRedeemCodeInfoAndErrorReturn(c, req.Code)
	if !ok {
		return
	}

	extendData := RedeemCodeExtendDataPro{}
	if err := json.Unmarshal([]byte(redeemCodeInfo.ExtendData), &extendData); err != nil {
		apiReturn.Error(c, "redemption code parsing failed")
		global.Logger.Errorln("redemption code parsing failed:", err, "data:", redeemCodeInfo.ExtendData)
		return
	}

	// 增加过期时间
	ChangeExpiredTimeByDayNumErr := biz.ProAuthorize.ChangeExpiredTimeByDayNum(userInfo.ID, extendData.Days, "Redeem code:"+req.Code, "", "")
	if ChangeExpiredTimeByDayNumErr != nil {
		global.Logger.Errorln("ChangeExpiredTimeByDayNum failed:", ChangeExpiredTimeByDayNumErr)
		apiReturn.Error(c, "redemption code failed")
		return
	}

	// 核销
	WriteOffErr := biz.RedeemCode.WriteOff(req.Code, userInfo.ID)
	if WriteOffErr != nil {
		global.Logger.Errorln("WriteOff failed:", WriteOffErr)
		apiReturn.Error(c, "redemption code failed")
		return
	}

	respInfo := RedeemCodeInfo{
		Code: redeemCodeInfo.Code,
		RedeemCodeParam: biz.RedeemCodeParam{
			ExpireTime:  base.ConvertTimeToUserTime(c, redeemCodeInfo.ExpireTime).Format("2006-01-02 15:04:05"),
			ReleaseType: redeemCodeInfo.ReleaseType,
			RedeemType:  redeemCodeInfo.RedeemType,
			Status:      redeemCodeInfo.Status,
			Title:       redeemCodeInfo.Title,
			ExtendData: map[string]interface{}{
				"days": extendData.Days,
			},
		},
	}

	apiReturn.SuccessData(c, respInfo)
}

func getRedeemCodeInfoAndErrorReturn(c *gin.Context, code string) (models.RedeemCode, bool) {
	redeemCodeInfo, err := biz.RedeemCode.Query(code)
	if err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			apiReturn.ErrorByCode(c, apiReturn.ErrCodeNoDataRecordFound)
			return redeemCodeInfo, false
		case biz.ErrRedeemCodeExpired:
			apiReturn.ErrorByCodeAndMsg(c, -2, "兑换码已过期")
			return redeemCodeInfo, false
		case biz.ErrRedeemCodeUsed:
			apiReturn.ErrorByCodeAndMsg(c, -3, "兑换码已使用")
			return redeemCodeInfo, false
		case biz.ErrRedeemCodeInvalid:
			apiReturn.ErrorByCodeAndMsg(c, -4, "兑换码已作废")
			return redeemCodeInfo, false
		}
	}

	return redeemCodeInfo, true
}
