package admin

import (
	"io"
	"strconv"
	"sun-panel/api/api_v1/common/apiReturn"
	"sun-panel/biz"
	"sun-panel/global"

	"github.com/gin-gonic/gin"
)

// MicroAppVersionUploadApi 微应用版本上传 API
type MicroAppVersionUploadApi struct{}

// Upload 上传微应用版本包
func (a *MicroAppVersionUploadApi) Upload(c *gin.Context) {
	// 获取上传的文件
	f, err := c.FormFile("file")
	if err != nil {
		apiReturn.Error(c, "请选择要上传的文件")
		return
	}
	appRecordIdStr := c.PostForm("appRecordId")
	if appRecordIdStr == "" {
		apiReturn.Error(c, "not appRecordId")
		return
	}

	// 检查文件扩展名
	fileExt := f.Filename[len(f.Filename)-4:]
	if fileExt != ".zip" {
		apiReturn.Error(c, "只支持 .zip 格式的文件")
		return
	}

	// 读取文件内容
	file, err := f.Open()
	if err != nil {
		apiReturn.Error(c, "打开文件失败")
		return
	}
	defer file.Close()

	fileData, err := io.ReadAll(file)
	if err != nil {
		apiReturn.Error(c, "读取文件失败")
		return
	}

	// 调用业务层处理
	result, err := biz.MicroAppPackage.UploadMicroAppPackage(fileData, f.Filename)
	if err != nil {
		apiReturn.Error(c, err.Error())
		return
	}

	appRecordId, err := strconv.Atoi(appRecordIdStr)
	if err != nil {
		apiReturn.Error(c, "not appRecordId")
		return
	}

	microApp, err := biz.MicroApp.GetById(global.Db, uint(appRecordId), "Developer")
	if err != nil {
		apiReturn.Error(c, "not appRecordId")
		return
	}

	// global.Logger.Debugln("MicroAppVersionUploadApi", "microApp:", cmn.AnyToJsonStr(microApp), "result:", cmn.AnyToJsonStr(result))

	// 检测基础信息
	if err := biz.MicroAppAudit.BasicCheck(microApp, result.Config); err != nil {
		apiReturn.ErrorByCodeAndMsg(c, -2, err.Error())
		return
	}

	// 缓存起来，用户手动点击确定创建的时候读取
	cacheKey := biz.MicroAppPackage.SetUploadCache(uint(appRecordId), "none", biz.MicroAppPackageUploadCache{
		PackageResult: result,
		AppRecordId:   uint(appRecordId),
	})

	// 返回结果
	apiReturn.SuccessData(c, MicroAppVersionUploadResp{
		// URL:           result.Src,
		Hash:          result.Hash,
		Config:        result.Config,
		FileName:      result.FileName,
		FileSize:      result.FileSize,
		FullFilePath:  result.FullFilePath,
		IconURL:       result.IconURL,
		UploadCacheId: cacheKey,
	})
}
