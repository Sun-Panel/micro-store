package admin

import (
	"encoding/json"
	"sun-panel/api/api_v1/common/apiReturn"
	"sun-panel/api/api_v1/common/base"
	"sun-panel/biz"
	"sun-panel/global"
	"sun-panel/lib/cmn"
	"sun-panel/models"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

// 兑换码
type RedeemCodeApi struct {
}

type RedeemCodeInfo struct {
	biz.RedeemCodeParam
	Code         string      `json:"code"` // 兑换码
	CreateTime   string      `json:"createTime"`
	WriteOffTime string      `json:"writeOffTime"`
	UserInfo     models.User `json:"userInfo"`
}

// 获取兑换码列表
func (a *RedeemCodeApi) GetList(c *gin.Context) {

	type ParamsStruct struct {
		models.User
		Limit   int
		Page    int
		Keyword string `json:"keyword"`
		Status  int    `json:"status"`
	}

	param := ParamsStruct{}
	if err := c.ShouldBindBodyWith(&param, binding.JSON); err != nil {
		apiReturn.Error(c, global.Lang.GetAndInsert("common.api_error_param_format", "[", err.Error(), "]"))
		c.Abort()
		return
	}

	var (
		list  []models.RedeemCode
		count int64
	)
	db := global.Db.Preload("User")

	// 查询条件
	if param.Keyword != "" {
		db = db.Where("title LIKE ? OR code LIKE ? OR note LIKE ?",
			"%"+param.Keyword+"%",
			"%"+param.Keyword+"%",
			"%"+param.Keyword+"%")
	}

	if param.Status != 0 {
		db = db.Where("status = ?", param.Status)
	}

	sqlError := db.Limit(param.Limit).Order("status ASC,expire_time DESC").
		Offset((param.Page - 1) * param.Limit).
		Find(&list).Limit(-1).
		Offset(-1).Count(&count).Error
	if sqlError != nil {
		apiReturn.ErrorDatabase(c, sqlError.Error())
		return
	}

	userLocal := base.GetUserTimezoneLocation(c)

	resp := []RedeemCodeInfo{}

	for _, v := range list {
		extendData := map[string]interface{}{}
		json.Unmarshal([]byte(v.ExtendData), &extendData)

		info := RedeemCodeInfo{
			Code:         v.Code,
			CreateTime:   v.CreatedAt.In(userLocal).Format(cmn.TimeFormatMode1),
			WriteOffTime: base.ConvertSQLNullTimeUserTimeToString(userLocal, v.WriteOffTime, cmn.TimeFormatMode1),
			UserInfo:     v.User,
			RedeemCodeParam: biz.RedeemCodeParam{
				ExpireTime:  base.ConvertTimeToUserTime(c, v.ExpireTime).Format("2006-01-02 15:04:05"),
				ReleaseType: v.ReleaseType,
				RedeemType:  v.RedeemType,
				Status:      v.Status,
				Title:       v.Title,
				ExtendData:  extendData,
			},
		}

		resp = append(resp, info)
	}

	apiReturn.SuccessListData(c, resp, count)
}

type RedeemCodeApiCreateReq struct {
	biz.RedeemCodeParam
	Number int    `json:"number"` // 创建数量
	Prefix string `json:"prefix"` // 前缀
}

// 创建兑换码
func (a *RedeemCodeApi) Create(c *gin.Context) {

	userLocal := base.GetUserTimezoneLocation(c)

	req := RedeemCodeApiCreateReq{}

	if err := c.ShouldBindJSON(&req); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}

	expireTime, expireTimeErr := time.ParseInLocation(cmn.TimeFormatMode1, req.ExpireTime, userLocal)

	if expireTimeErr != nil {
		apiReturn.ErrorParamFomat(c, "expireTime invalid")
		return
	}

	if req.Number > 100 {
		apiReturn.ErrorParamFomat(c, "number must be less than 100")
		return
	}

	param := biz.RedeemCodeParam{}
	param.ExpireTime = base.ConvertUserTimeToSystemTime(userLocal, expireTime).Format(cmn.TimeFormatMode1)
	param.ReleaseType = 1
	param.RedeemType = req.RedeemType
	param.Status = 1
	param.Title = req.Title
	param.Note = req.Note
	param.ExtendData = req.ExtendData

	lists := []RedeemCodeInfo{}
	for i := 0; i < req.Number; i++ {
		info, err := biz.RedeemCode.Add(req.Prefix, param)
		if err != nil {
			apiReturn.Error(c, err.Error())
			return
		}

		lists = append(lists, RedeemCodeInfo{
			Code: info.Code,
		})
	}

	apiReturn.SuccessListData(c, lists, int64(len(lists)))
}

type RedeemCodeApiSetInvalidReq struct {
	Codes []string `json:"codes"` // 兑换码列表
}

// 设置兑换码失效
func (a *RedeemCodeApi) SetInvalid(c *gin.Context) {

	req := RedeemCodeApiSetInvalidReq{}

	if err := c.ShouldBindJSON(&req); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}

	if err := biz.RedeemCode.SetInvalid(req.Codes...); err != nil {
		apiReturn.Error(c, err.Error())
		return
	}

	apiReturn.Success(c)
}
