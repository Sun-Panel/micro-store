package admin

import (
	"sun-panel/api/api_v1/common/apiReturn"
	"sun-panel/biz"
	"sun-panel/models"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type VersionSecretApi struct {
}

func (a *VersionSecretApi) Edit(c *gin.Context) {
	req := models.VersionSecret{}
	if err := c.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}

	if req.ID == 0 {
		// 判断版本号是否已经存在
		_, err := biz.VersionSecret.GetVersionSecretFromDb(req.Version)
		if err == nil {
			apiReturn.ErrorDatabase(c, "版本号已经存在")
			return
		}

		// 创建版本密钥
		createErr := biz.VersionSecret.CreateVersionSecret(models.VersionSecret{
			Version:   req.Version,
			SecretKey: req.SecretKey,
			Status:    true,
		})
		if createErr != nil {
			apiReturn.ErrorDatabase(c, createErr.Error())
			return
		}

		apiReturn.Success(c)
		return
	}

	// 更新版本密钥
	versionSecret, _ := biz.VersionSecret.GetVersionSecretByIdFromDb(req.ID)
	versionSecret.Version = req.Version
	versionSecret.SecretKey = req.SecretKey
	versionSecret.Status = req.Status

	// 判断版本号是否已经存在
	foundData, err := biz.VersionSecret.GetVersionSecretFromDb(req.Version)
	if err == nil && foundData.ID != versionSecret.ID {
		apiReturn.ErrorDatabase(c, "版本号已经存在")
		return
	}

	if err := biz.VersionSecret.UpdateVersionSecret(versionSecret); err != nil {
		apiReturn.ErrorDatabase(c, err.Error())
		return
	}
	apiReturn.Success(c)
}

// 获取某个版本的信息根据版本号
func (a *VersionSecretApi) GetVersionSecretInfoByVersion(c *gin.Context) {
	req := models.VersionSecret{}
	if err := c.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}
	versionSecret, err := biz.VersionSecret.GetVersionSecretFromDb(req.Version)
	if err != nil {
		apiReturn.ErrorDatabase(c, err.Error())
		return
	}
	apiReturn.SuccessData(c, versionSecret)
}
