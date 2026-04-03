package microapp

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"sun-panel/api/api_v1/common/apiReturn"
	"sun-panel/biz"
	"sun-panel/global"
	"sun-panel/models"

	"github.com/gin-gonic/gin"
)

type DownloadApi struct{}

func (a *DownloadApi) GetUrl(c *gin.Context) {
	microAppId := c.Param("microAppId")
	version := c.Param("version")
	if microAppId == "" {
		apiReturn.ErrorParamFomat(c, "microAppId is required")
		return
	}

	url := biz.MicroAppPackage.BuildDownloadUrl(microAppId, version)
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
	a.serveFileNonStreaming(c, filePath, versionInfo.Version)
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

// serveFile 流式传输文件（支持断点续传）
func (a *DownloadApi) serveFile(c *gin.Context, filePath, version string) {
	// 打开文件
	file, err := os.Open(filePath)
	if err != nil {
		apiReturn.ErrorDataNotFound(c)
		return
	}
	defer file.Close()

	// 获取文件信息
	fileInfo, err := file.Stat()
	if err != nil {
		apiReturn.ErrorDataNotFound(c)
		return
	}

	// 设置响应头
	fileName := fmt.Sprintf("micro_app_v%s.zip", version)
	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", fileName))
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Expires", "0")
	c.Header("Cache-Control", "must-revalidate")
	c.Header("Pragma", "public")
	c.Header("Content-Length", strconv.FormatInt(fileInfo.Size(), 10))

	// 支持 Range 请求（断点续传）
	http.ServeContent(c.Writer, c.Request, fileName, fileInfo.ModTime(), file)
}

// serveFileNonStreaming 非流式传输文件（一次性读取，可检测完整传输）
func (a *DownloadApi) serveFileNonStreaming(c *gin.Context, filePath, version string) {
	// 读取整个文件到内存
	data, err := os.ReadFile(filePath)
	if err != nil {
		apiReturn.ErrorDataNotFound(c)
		return
	}

	// 设置响应头
	fileName := fmt.Sprintf("micro_app_v%s.zip", version)
	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", fileName))
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Expires", "0")
	c.Header("Cache-Control", "must-revalidate")
	c.Header("Pragma", "public")
	c.Header("Content-Length", strconv.Itoa(len(data)))

	// 写入响应（完成后可以执行token过期等操作）
	c.Data(http.StatusOK, "application/octet-stream", data)
}
