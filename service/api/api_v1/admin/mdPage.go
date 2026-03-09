package admin

import (
	"encoding/json"
	"sun-panel/api/api_v1/common/apiData/commonApiStructs"
	"sun-panel/api/api_v1/common/apiData/systemApiStructs"
	"sun-panel/api/api_v1/common/apiReturn"
	"sun-panel/global"
	"sun-panel/models"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type MdPageApi struct {
}

var MdPageVarPrefix = "mdPage."

func (a *MdPageApi) GetList(c *gin.Context) {
	req := commonApiStructs.RequestPage{}

	if err := c.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}

	var count int64

	varList := []models.SystemSetting{}
	dbErr := global.Db.Order("id asc").
		Limit(req.Limit).Offset((req.Page-1)*req.Limit).
		Where(
			"config_name like ? OR (config_name = ? AND description like ?) ",
			MdPageVarPrefix+req.Keyword+"%",
			MdPageVarPrefix,
			"%"+req.Keyword+"%",
		).
		Find(&varList).
		Limit(-1).Offset(-1).Count(&count).Error

	if dbErr != nil {
		apiReturn.ErrorDatabase(c, dbErr.Error())
		return
	}

	respList := []systemApiStructs.MdPageListItem{}
	for _, v := range varList {
		mdPageInfo := systemApiStructs.MdPageInfo{}
		json.Unmarshal([]byte(v.ConfigValue), &mdPageInfo)
		respList = append(respList, systemApiStructs.MdPageListItem{
			MdPageName:        v.ConfigName[len(MdPageVarPrefix):],
			MdPageDescription: v.Description,
			MdPageInfo:        mdPageInfo,
		})
	}
	apiReturn.SuccessListData(c, respList, count)
}

func (a *MdPageApi) Edit(c *gin.Context) {
	req := systemApiStructs.MdPageInfoEditReq{}

	if err := c.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}

	mgName := MdPageVarPrefix + req.MdPageName

	if err := global.SystemSetting.Set(mgName, req.MdPageInfo); err != nil {
		apiReturn.Error(c, err.Error())
	} else {
		global.Db.Model(&models.SystemSetting{}).Where("config_name=?", mgName).Update("description", req.MdPageDescription)
		apiReturn.Success(c)
	}

}

func (a *MdPageApi) Delete(c *gin.Context) {
	req := systemApiStructs.MdPageInfoEditReq{}

	if err := c.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}
	mgName := MdPageVarPrefix + req.MdPageName
	if err := global.SystemSetting.Delete(mgName); err != nil {
		apiReturn.Error(c, err.Error())
		return
	}

	apiReturn.Success(c)
}
