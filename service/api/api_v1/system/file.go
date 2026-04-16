package system

import (
	"fmt"
	"os"
	"path"
	"strings"
	"sun-panel/api/api_v1/common/apiData/commonApiStructs"
	"sun-panel/api/api_v1/common/apiReturn"
	"sun-panel/api/api_v1/common/base"
	"sun-panel/global"
	"sun-panel/lib/cmn"
	"sun-panel/lib/file"
	"sun-panel/models"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"gorm.io/gorm"
)

type FileApi struct{}

func (a *FileApi) UploadImg(c *gin.Context) {
	userInfo, _ := base.GetCurrentUserInfo(c)
	configUpload := global.Config.GetValueString("base", "source_path")

	// 创建按日期分类的子目录
	fildDir := fmt.Sprintf("%s/%d/%d/%d/", configUpload, time.Now().Year(), time.Now().Month(), time.Now().Day())
	isExist, _ := cmn.PathExists(fildDir)
	if !isExist {
		os.MkdirAll(fildDir, os.ModePerm)
	}

	// 使用 lib/file/upload 进行文件上传
	uploadInfo, err := file.FormFileUpload(c, file.FormFileUploadOptions{
		FormFileName: "imgfile",
		AgreeExtNames: []string{
			".png",
			".jpg",
			".gif",
			".jpeg",
			".webp",
			".svg",
			".ico",
		},
		MaxSize: 1024 * 1024 * 5, // 5MB
		SaveDir: fildDir,
	})

	if err != nil {
		switch err {
		case file.ErrUploadExceedMaxSize:
			apiReturn.Error(c, "文件尺寸不能大于5M")
			return
		case file.ErrUploadExtensionNameNotAllowed:
			apiReturn.ErrorByCode(c, 1301)
			return
		default:
			apiReturn.ErrorByCode(c, 1300)
			return
		}
	}

	pureFilePath := uploadInfo.FileSavePath[len(configUpload):]
	// 向数据库添加记录
	mFile := models.File{}
	mFile.AddFile(userInfo.ID, uploadInfo.FileOriginalName, uploadInfo.Ext, pureFilePath)

	responseData := gin.H{
		"imageUrl": global.UPLOAD_ROUTE + pureFilePath,
	}

	// 如果是图片文件，添加尺寸信息
	if uploadInfo.Width > 0 && uploadInfo.Height > 0 {
		responseData["width"] = uploadInfo.Width
		responseData["height"] = uploadInfo.Height
	}

	apiReturn.SuccessData(c, responseData)
}

// 一直未使用有待测试
func (a *FileApi) UploadFiles(c *gin.Context) {
	userInfo, _ := base.GetCurrentUserInfo(c)
	configUpload := global.Config.GetValueString("base", "source_path")

	form, err := c.MultipartForm()
	if err != nil {
		apiReturn.ErrorByCode(c, 1300)
		return
	}
	files := form.File["files[]"]
	errFiles := []string{}
	succMap := map[string]string{}
	for _, f := range files {
		fileExt := strings.ToLower(path.Ext(f.Filename))
		fileName := cmn.Md5(fmt.Sprintf("%s%s", f.Filename, time.Now().String()))
		fildDir := fmt.Sprintf("%s/%d/%d/%d/", configUpload, time.Now().Year(), time.Now().Month(), time.Now().Day())
		isExist, _ := cmn.PathExists(fildDir)
		if !isExist {
			os.MkdirAll(fildDir, os.ModePerm)
		}
		filepath := fmt.Sprintf("%s%s%s", fildDir, fileName, fileExt)
		if c.SaveUploadedFile(f, filepath) != nil {
			errFiles = append(errFiles, f.Filename)
		} else {
			// 成功
			pureFilePath := filepath[len(configUpload):]
			// 向数据库添加记录
			mFile := models.File{}
			mFile.AddFile(userInfo.ID, f.Filename, fileExt, pureFilePath)
			succMap[f.Filename] = global.UPLOAD_ROUTE + pureFilePath
		}
	}

	apiReturn.SuccessData(c, gin.H{
		"succMap":  succMap,
		"errFiles": errFiles,
	})
}

func (a *FileApi) GetList(c *gin.Context) {
	list := []models.File{}
	userInfo, _ := base.GetCurrentUserInfo(c)
	var count int64
	if err := global.Db.Order("created_at desc").Find(&list, "user_id=?", userInfo.ID).Count(&count).Error; err != nil {
		apiReturn.ErrorDatabase(c, err.Error())
		return
	}

	data := []map[string]interface{}{}
	for _, v := range list {
		// 兼容旧版
		if len(v.Src) >= 0 && v.Src[0] == '.' {
			v.Src = v.Src[1:]
		} else {
			v.Src = global.UPLOAD_ROUTE + v.Src
		}
		data = append(data, map[string]interface{}{
			"src":        v.Src,
			"fileName":   v.FileName,
			"id":         v.ID,
			"createTime": v.CreatedAt,
			"updateTime": v.UpdatedAt,
			"path":       v.Src,
		})
	}
	apiReturn.SuccessListData(c, data, count)
}

func (a *FileApi) Deletes(c *gin.Context) {
	req := commonApiStructs.RequestDeleteIds[uint]{}
	userInfo, _ := base.GetCurrentUserInfo(c)
	if err := c.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}

	global.Db.Transaction(func(tx *gorm.DB) error {
		files := []models.File{}

		if err := tx.Order("created_at desc").Find(&files, "user_id=? AND id in ?", userInfo.ID, req.Ids).Error; err != nil {
			return err
		}

		for _, v := range files {
			os.Remove(v.Src)
		}

		if err := tx.Order("created_at desc").Delete(&files, "user_id=? AND id in ?", userInfo.ID, req.Ids).Error; err != nil {
			return err
		}

		return nil
	})

	apiReturn.Success(c)

}
