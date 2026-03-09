package admin

import (
	"sun-panel/api/api_v1/common/apiData/commonApiStructs"
	"sun-panel/api/api_v1/common/apiReturn"
	"sun-panel/biz"
	"sun-panel/global"
	"sun-panel/lib/cmn"
	"sun-panel/models"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type VersionApi struct {
}

type VersionSetActiveReq struct {
	ID          uint   `json:"id"`
	VersionType string `json:"versionType"`
}

type VersionInfoResp struct {
	ID             uint   `json:"id"`
	Version        string `json:"version"`
	Type           string `json:"type"`
	ReleaseTime    string `json:"releaseTime"`
	Description    string `json:"description"`
	DownloadURL    string `json:"downloadURL"`
	PageUrl        string `json:"pageUrl"`
	IsActive       bool   `json:"isActive"`
	IsRolledBack   bool   `json:"isRolledBack"`
	AloneSecretKey int    `json:"aloneSecretKey"` // 0 不存在 1 存在 2 已停用
}

type VersionInfoReq struct {
	models.Version
	ReleaseTime string `json:"releaseTime"`
	SecretKey   string `json:"secretKey"`
}

// type ParamUserInfo struct {
// 	UserId    uint   `json:"userId"`
// 	Username  string `json:"username" validate:"required,email"`
// 	Password  string `json:"password" validate:"required"`
// 	Name      string `json:"name" `
// 	HeadImage string `json:"headImage" `
// 	Status    int    `json:"status" `
// 	Role      int    `json:"role" `
// 	Mail      string `json:"mail" `
// }

func (a VersionApi) Edit(c *gin.Context) {
	param := VersionInfoReq{}
	if err := c.ShouldBindBodyWith(&param, binding.JSON); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}

	ReleaseTime, _ := cmn.StrToTime(cmn.TimeFormatMode1, param.ReleaseTime)

	newData := models.Version{
		ID:           param.ID,
		Version:      param.Version.Version,
		Type:         models.VersionType(param.Type),
		ReleaseTime:  ReleaseTime,
		Description:  param.Description,
		DownloadURL:  param.DownloadURL,
		PageUrl:      param.PageUrl,
		IsActive:     false,
		IsRolledBack: param.IsRolledBack,
	}

	var err error

	if param.ID == 0 {
		err = global.Db.Create(&newData).Error
	} else {
		err = global.Db.Model(&models.Version{}).Where("id = ?", param.ID).Select(
			"Version",
			"Type",
			"ReleaseTime",
			"Description",
			"DownloadURL",
			"PageUrl",
			"IsActive",
			"IsRolledBack",
		).Updates(&newData).Error
	}

	if err != nil {
		apiReturn.ErrorDatabase(c, err.Error())
		return
	}
	apiReturn.Success(c)
}

func (a VersionApi) Deletes(c *gin.Context) {
	req := commonApiStructs.RequestDeleteIds[uint]{}
	if err := c.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		c.Abort()
		return
	}

	// 获取版本信息
	for _, v := range req.Ids {
		var versionInfo models.Version
		err := global.Db.First(&versionInfo, "id=?", v).Error
		if err != nil {
			apiReturn.ErrorDatabase(c, err.Error())
			return
		}

		secretInfo, getVersionSecretErr := biz.VersionSecret.GetVersionSecret(versionInfo.Version)
		if getVersionSecretErr != nil {
			global.Logger.Error("getVersionSecretErr", getVersionSecretErr.Error())
		}
		deleteVersionSecret := biz.VersionSecret.DeleteVersionSecret(secretInfo)
		if deleteVersionSecret != nil {
			global.Logger.Error("deleteVersionSecret", deleteVersionSecret.Error())
		}
	}

	if err := global.Db.Delete(&models.Version{}, &req.Ids).Error; err != nil {
		apiReturn.ErrorDatabase(c, err.Error())
		return
	}

	apiReturn.Success(c)
}

func (a VersionApi) SetActive(c *gin.Context) {
	req := VersionSetActiveReq{}
	if err := c.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}

	updateModels := models.Version{}

	// 给最新版增加标记
	err := global.Db.Model(&updateModels).Where("id=? AND type=?", req.ID, req.VersionType).Update("is_active", true).Error
	if err != nil {
		apiReturn.ErrorDatabase(c, err.Error())
		return
	}

	// 清除所有的最新版本标记
	err = global.Db.Model(&updateModels).Where("id!=? AND type=? AND is_active=?", req.ID, req.VersionType, true).Update("is_active", false).Error
	if err != nil {
		apiReturn.ErrorDatabase(c, err.Error())
		return
	}
	apiReturn.Success(c)
}

func (a VersionApi) GetList(c *gin.Context) {

	req := commonApiStructs.RequestPage{}
	if err := c.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		c.Abort()
		return
	}

	var (
		list  []models.Version
		count int64
	)
	db := global.Db

	// 查询条件
	if req.Keyword != "" {
		db = db.Where("version LIKE ?", "%"+req.Keyword+"%")
	}

	if err := db.Preload("VersionSecret").Order("release_time Desc").Limit(req.Limit).Offset((req.Page - 1) * req.Limit).Find(&list).Limit(-1).Offset(-1).Count(&count).Error; err != nil {
		apiReturn.ErrorDatabase(c, err.Error())
		return
	}

	resMap := []VersionInfoResp{}
	for _, v := range list {
		aloneSecretKey := 0
		if v.VersionSecret.Version == "" {
			aloneSecretKey = 0
		} else {
			if v.VersionSecret.Status {
				aloneSecretKey = 1
			} else {
				aloneSecretKey = 2
			}
		}

		resMap = append(resMap, VersionInfoResp{
			ID:             v.ID,
			Version:        v.Version,
			Type:           string(v.Type),
			ReleaseTime:    v.ReleaseTime.Format(cmn.TimeFormatMode1),
			Description:    v.Description,
			DownloadURL:    v.DownloadURL,
			PageUrl:        v.PageUrl,
			IsActive:       v.IsActive,
			IsRolledBack:   v.IsRolledBack,
			AloneSecretKey: aloneSecretKey,
		})
	}

	apiReturn.SuccessListData(c, resMap, count)
}
