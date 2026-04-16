package file

import (
	"errors"
	"fmt"
	"mime/multipart"
	"os"
	"path"
	"strings"
	"time"

	"sun-panel/lib/cmn"

	"github.com/gin-gonic/gin"
)

var (
	ErrUploadExtensionNameNotAllowed = errors.New("extension not allowed") // "上传文件类型不允许"
	ErrUploadExceedMaxSize           = errors.New("exceeds maximum size")  // "上传文件超出最大尺寸"
)

type FormFileUploadOptions struct {
	// 表单文件字段名
	FormFileName string
	// 同意的文件扩展名,为空不检查
	AgreeExtNames []string
	// 允许的最大文件大小，单位字节，为0不检查
	MaxSize int64
	// 文件保存目录
	SaveDir string
}

type FormFileUploadFileInfo struct {
	// 文件储存路径
	FileSavePath string
	// 文件原始名称
	FileOriginalName string
	// 文件信息
	FormFile *multipart.FileHeader
	// 文件扩展名 eg: .zip
	Ext string
	// 图片宽度（仅图片文件）
	Width int
	// 图片高度（仅图片文件）
	Height int
}

// 表单上传文件
func FormFileUpload(c *gin.Context, options FormFileUploadOptions) (uploadInfo FormFileUploadFileInfo, err error) {

	f, err := c.FormFile(options.FormFileName)
	if err != nil {
		return
	} else {
		fileExt := strings.ToLower(path.Ext(f.Filename))

		if len(options.AgreeExtNames) != 0 && !cmn.InArray(options.AgreeExtNames, fileExt) {
			err = ErrUploadExtensionNameNotAllowed
			return
		}

		if options.MaxSize != 0 && f.Size > options.MaxSize {
			err = ErrUploadExceedMaxSize
			return
		}

		fileName := cmn.Md5(fmt.Sprintf("%s%s", f.Filename, time.Now().String()))
		fildNamePrefix := fmt.Sprintf("%s/%s-%s-", options.SaveDir, time.Now().Format("20060102150405"), cmn.BuildRandCode(10, cmn.RAND_CODE_MODE2))
		isExist, _ := cmn.PathExists(options.SaveDir)
		if !isExist {
			os.MkdirAll(options.SaveDir, os.ModePerm)
		}
		filepath := fmt.Sprintf("%s%s%s", fildNamePrefix, fileName, fileExt)
		c.SaveUploadedFile(f, filepath)

		// 尝试获取图片尺寸
		width, height, _ := GetImageDimensions(filepath, fileExt)

		uploadInfo = FormFileUploadFileInfo{
			FileSavePath:     filepath,
			FileOriginalName: f.Filename,
			FormFile:         f,
			Ext:              fileExt,
			Width:            width,
			Height:           height,
		}
	}

	return
}
