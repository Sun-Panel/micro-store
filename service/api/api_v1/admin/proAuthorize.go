package admin

import (
	"sun-panel/api/api_v1/common/apiData/adminApiStructs"
	"sun-panel/api/api_v1/common/apiData/commonApiStructs"
	"sun-panel/api/api_v1/common/apiData/proAuthorizeApiStructs"
	"sun-panel/api/api_v1/common/apiReturn"

	"sun-panel/biz"
	"sun-panel/global"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type ProAuthorizeApi struct {
}

func (a *ProAuthorizeApi) UpdateUserExpiredTimeByDay(c *gin.Context) {
	param := adminApiStructs.ProAuthorizeUpdateUserExpiredTimeByDayReq{}

	if err := c.ShouldBindWith(&param, binding.JSON); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}

	if err := biz.ProAuthorize.ChangeExpiredTimeByDayNum(param.UserId, param.DayNum, param.Note, param.AdminNote, ""); err != nil {
		apiReturn.Error(c, err.Error())
		return
	}
	apiReturn.Success(c)

}

func (a *ProAuthorizeApi) GetUserProAuthorizeList(c *gin.Context) {

	req := commonApiStructs.RequestPage{}

	if err := c.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}

	var count int64

	result := []adminApiStructs.ProAuthorizeGetUserProAuthorizeListItemResp{}
	db := global.Db.Table("user as u").
		Select("u.username, u.name,u.id as user_id,pa.expired_time").
		Joins("LEFT JOIN pro_authorize AS pa ON u.id = pa.user_id").
		Order("pa.created_at DESC")

	// 条件
	if req.Keyword != "" {
		db.Where("u.username like ? OR name like ?", "%"+req.Keyword+"%", "%"+req.Keyword+"%")
	}

	dbErr := db.Limit(req.Limit).Offset((req.Page - 1) * req.Limit).
		Scan(&result).
		Limit(-1).Offset(-1).Count(&count).Error

	if dbErr != nil {
		apiReturn.ErrorDatabase(c, dbErr.Error())
		return
	}

	apiReturn.SuccessListData(c, result, count)

}

func (a *ProAuthorizeApi) GetAuthorizeHistoryRecordByUserId(c *gin.Context) {
	type Request struct {
		UserId uint `json:"userId"`
	}
	req := Request{}

	if err := c.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}

	resp := []proAuthorizeApiStructs.GetAuthorizeHistoryRecordResp{}
	if resList, err := biz.ProAuthorize.GetAuthorizeHistoryRecordByUserId(req.UserId); err != nil {
		apiReturn.ErrorDatabase(c, err.Error())
		return
	} else {
		for _, v := range resList {
			resp = append(resp, proAuthorizeApiStructs.GetAuthorizeHistoryRecordResp{
				ChangeTime:  v.CreatedAt,
				ExpiredTime: v.ExpiredTime,
				DayNum:      v.DayNum,
				Note:        v.Note,
				OrderNo:     v.OrderNo,
				AdminNote:   v.AdminNote,
			})
		}
	}

	apiReturn.SuccessListData(c, resp, 0)

}
