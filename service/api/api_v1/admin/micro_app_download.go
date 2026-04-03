package admin

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"sun-panel/api/api_v1/common/apiReturn"
	"sun-panel/biz"
	"sun-panel/global"
	"sun-panel/lib/cmn"
	"sun-panel/models"
	"time"

	"github.com/gin-gonic/gin"
)

type MicroAppDownloadApi struct {
}

func (a *MicroAppDownloadApi) GetUrl(c *gin.Context) {

	versionId := c.Param("versionId")
	if versionId == "" {
		apiReturn.ErrorParamFomat(c, "versionId is required")
		return
	}

	// 生成下载URL
	token := cmn.BuildRandCode(32, cmn.RAND_CODE_MODE1) + time.Now().Format("20060102150405")
	downloadUrl := fmt.Sprintf("/api/admin/microApp/download/tk/%s", token)

	// 保存token到缓存
	// a.downloadCache.SetDefault(token, versionId)
	biz.MicroAppDownload.AdminDownloadCacheVersionId.SetDefault(token, versionId)

	apiReturn.SuccessData(c, downloadUrl)
}

func (a *MicroAppDownloadApi) DownloadByVersionId(c *gin.Context) {

	downloadToken := c.Param("downloadToken")
	if downloadToken == "" {
		apiReturn.ErrorParamFomat(c, "token is required")
		return
	}

	versionId, exist := biz.MicroAppDownload.AdminDownloadCacheVersionId.Get(downloadToken)
	if !exist {
		apiReturn.ErrorParamFomat(c, "token is not exist")
		return
	}

	// 转换为 uint
	versionIdUint, err := strconv.ParseUint(versionId, 10, 32)
	if err != nil {
		apiReturn.ErrorParamFomat(c, "versionId invalid")
		return
	}

	// 查询版本信息
	m := models.MicroAppVersion{}
	version, err := m.GetById(global.Db, uint(versionIdUint))
	if err != nil {
		apiReturn.ErrorDataNotFound(c)
		return
	}
	global.Logger.Debugln("PackageSrc: ", version.PackageSrc)

	// 获取文件路径
	filePath := a.getFilePath(version.PackageSrc)
	if filePath == "" {
		apiReturn.ErrorDataNotFound(c)
		return
	}

	global.Logger.Debugln("download file path: ", filePath)

	// 流式传输文件（不统计下载次数）
	a.serveFile(c, filePath, version.Version)
}

// getFilePath 获取文件的完整路径
func (a *MicroAppDownloadApi) getFilePath(packageSrc string) string {
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
func (a *MicroAppDownloadApi) serveFile(c *gin.Context, filePath, version string) {
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
	fileName := fileInfo.Name()
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
