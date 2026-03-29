package admin

import (
	"io"
	"sun-panel/api/api_v1/common/apiReturn"
	"sun-panel/biz"

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

	// 返回结果
	apiReturn.SuccessData(c, MicroAppVersionUploadResp{
		URL:        result.URL,
		Hash:       result.Hash,
		Config:     result.Config,
		FileName:   result.FileName,
		FileSize:   result.FileSize,
		FolderName: result.FolderName,
		IconURL:    result.IconURL,
	})
}
