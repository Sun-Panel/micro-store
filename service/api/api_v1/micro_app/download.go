package microapp

import (
	"os"
	"path/filepath"
	"sun-panel/api/api_v1/common/apiReturn"
	"sun-panel/api/api_v1/common/base"
	"sun-panel/biz"
	"sun-panel/global"
	"sun-panel/models"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type DownloadApi struct{}

func (a *DownloadApi) GetUrl(c *gin.Context) {
	req := DownloadGetUrlReq{}
	if err := c.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}

	url := biz.MicroAppPackage.BuildDownloadUrl(req.MicroAppId, req.Version)
	apiReturn.SuccessData(c, url)
}

// DownloadByVersionLatest 通过版本号下载/最新版本
func (a *DownloadApi) DownloadByVersionOrLatest(c *gin.Context) {
	microAppId := c.Param("microAppId")
	versionStr := c.Param("version")
	if microAppId == "" {
		apiReturn.ErrorParamFomat(c, "microAppId is required")
		return
	}

	// 查询应用信息
	microApp, err := biz.MicroApp.GetInfo(global.Db, microAppId)
	if err != nil {
		apiReturn.ErrorDataNotFound(c)
		return
	}

	// 检查应用状态（只有上架的应用才能下载）
	if microApp.Status != 1 {
		apiReturn.ErrorDataNotFound(c)
		return
	}

	var versionInfo models.MicroAppVersion

	if versionStr == "" {
		// 获取最新在线版本
		versionModel := models.MicroAppVersion{}
		latestVersion, err := versionModel.GetLatestOnlineByAppId(global.Db, microApp.ID)
		if err != nil {
			apiReturn.ErrorDataNotFound(c)
			return
		}
		versionInfo = latestVersion
	} else {
		// 获取指定版本信息
		v, err := biz.MicroAppVersion.GetInfoOnLineByVersion(global.Db, versionStr)
		if err != nil {
			apiReturn.ErrorDataNotFound(c)
			return
		}
		versionInfo = v
	}

	// 记录下载统计（异步）
	go func() {
		biz.MicroAppStatistics.IncrementDownload(
			versionInfo.AppRecordId,
			0,
			c.Query("clientId"),
			c.ClientIP(),
		)
	}()

	// 获取文件路径
	filePath := a.getFilePath(versionInfo.PackageSrc)
	if filePath == "" {
		apiReturn.ErrorDataNotFound(c)
		return
	}

	// 流式传输文件
	// a.serveFile(c, filePath, versionInfo.Version)
	// 非流式传输文件
	base.ServeFileNonStreaming(c, filePath)
}

// ===================================================================================================
// 辅助方法
// ===================================================================================================

// getFilePath 获取文件的完整路径
func (a *DownloadApi) getFilePath(packageSrc string) string {
	// 获取配置的上传路径
	uploadPath := global.Config.GetValueString("base", "micro_app_source_path")
	if uploadPath == "" {
		uploadPath = "./micro_app_upload"
	}

	// 拼接完整路径
	filePath := filepath.Join(uploadPath, packageSrc)

	// 检查文件是否存在
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return ""
	}

	return filePath
}
