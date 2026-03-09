package admin

import (
	"sun-panel/api/api_v1/common/apiData/commonApiStructs"
	"sun-panel/api/api_v1/common/apiData/systemApiStructs"
	"sun-panel/api/api_v1/common/apiReturn"
	"sun-panel/global"
	"sun-panel/models"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type SystemVariableApi struct {
}

func (a *SystemVariableApi) GetList(c *gin.Context) {
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
			"config_name like ? OR (config_name = ? AND description like ?)",
			"var.%"+req.Keyword+"%",
			req.Keyword,
			"%"+req.Keyword+"%",
		).
		Find(&varList).
		Limit(-1).Offset(-1).Count(&count).Error

	if dbErr != nil {
		apiReturn.ErrorDatabase(c, dbErr.Error())
		return
	}

	respList := []systemApiStructs.SystemVariableListItem{}
	for _, v := range varList {
		respList = append(respList, systemApiStructs.SystemVariableListItem{
			Id:          v.ID,
			ConfigName:  v.ConfigName[4:],
			ConfigValue: v.ConfigValue,
			Description: v.Description,
		})
	}
	apiReturn.SuccessListData(c, respList, count)
}

func (a *SystemVariableApi) Edit(c *gin.Context) {
	prefix := "var."
	req := systemApiStructs.SystemVariableEditReq{}

	if err := c.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}

	if req.Id == 0 {
		findRes := models.SystemSetting{}
		if err := global.Db.First(&findRes, "config_name = ?", prefix+req.Name).Error; err == nil {
			apiReturn.Error(c, "config_name already exists")
			return
		}
		findRes.ConfigName = prefix + req.Name
		// findRes.ConfigValue = req.Value
		findRes.Description = req.Description
		if err := global.Db.Model(&models.SystemSetting{}).Create(&findRes).Error; err != nil {
			apiReturn.ErrorDatabase(c, err.Error())
			return
		}
	} else {
		findRes := models.SystemSetting{}
		if err := global.Db.First(&findRes, "config_name = ? AND id!=?", prefix+req.Name, req.Id).Error; err == nil {
			apiReturn.Error(c, "config_name already exists")
			return
		}

		findRes.ConfigName = prefix + req.Name
		// findRes.ConfigValue = req.Value
		findRes.Description = req.Description
		if err := global.Db.Select("ConfigName", "Description").Where("id=?", req.Id).Updates(&findRes).Error; err != nil {
			apiReturn.ErrorDatabase(c, err.Error())
			return
		}
	}

	if err := global.SystemSetting.SetVariable(req.Name, req.Value); err != nil {
		apiReturn.Error(c, err.Error())
	} else {
		apiReturn.Success(c)
	}

}

func (a *SystemVariableApi) Set(c *gin.Context) {
	req := systemApiStructs.SystemVariable{}

	if err := c.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}

	if err := global.SystemSetting.SetVariable(req.Name, req.Value); err != nil {
		apiReturn.Error(c, err.Error())
		return
	} else {
		apiReturn.Success(c)
	}
}

func (a *SystemVariableApi) Delete(c *gin.Context) {
	req := systemApiStructs.SystemVariable{}
	prefix := "var."

	if err := c.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}

	if err := global.Db.Delete(&models.SystemSetting{}, "config_name=?", prefix+req.Name).Error; err != nil {
		apiReturn.Error(c, err.Error())
		return
	}
	apiReturn.Success(c)
}

type SystemVariableApiGet struct {
	Name string
}

// 清理变量的缓存
func (a *SystemVariableApi) ClearCache(c *gin.Context) {
	req := SystemVariableApiGet{}

	if err := c.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}

	global.SystemSetting.ClearVaribaleCache(req.Name)

	apiReturn.Success(c)
}

func (a *SystemVariableApi) GetByCache(c *gin.Context) {
	req := SystemVariableApiGet{}

	if err := c.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}

	reap := ""
	if v, err := global.SystemSetting.GetVariableString(req.Name); err != nil {
		apiReturn.Error(c, err.Error())
		return
	} else {
		reap = v
	}

	apiReturn.SuccessData(c, reap)
}
