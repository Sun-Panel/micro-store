package microapp

import (
	"fmt"
	"os"
	"sun-panel/api/api_v1/common/apiReturn"
	"sun-panel/api/api_v1/common/base"
	"sun-panel/biz"
	"sun-panel/global"
	"sun-panel/lib/cmn"
	"sun-panel/lib/file"
	"sun-panel/models"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type MicroAppApi struct{}

// GetInfo 获取微应用详情（管理员专用）
func (a *MicroAppApi) GetInfo(c *gin.Context) {
	req := models.BaseModel{}
	if err := c.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}

	lang := base.GetCurrentUserLang(c)

	info, err := biz.MicroApp.GetByIdWithLang(global.Db, req.ID, lang, "Developer")
	if err != nil {
		base.HandleBizErrorAndReturn(c, err)
		return
	}

	if info.Status == 0 {
		apiReturn.ErrorByCode(c, apiReturn.ErrCodeAppNotFound)
		return
	}

	apiReturn.SuccessData(c, info)
}

// 获取微应用列表
func (a *MicroAppApi) GetList(c *gin.Context) {
	req := MicroAppVersionGetListReq{}
	if err := c.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}

	// 参数校验和默认值设置
	if req.Page < 1 {
		req.Page = 1
	}
	if req.Limit < 1 {
		req.Limit = 10 // 默认每页10条
	}
	if req.Limit > 100 {
		req.Limit = 100 // 最大每页100条
	}

	// 调用业务层获取列表
	list, total, err := biz.MicroApp.GetList(global.Db, biz.GetListOptions{
		Page:       req.Page,
		Limit:      req.Limit,
		Order:      req.Order,
		CategoryId: req.CategoryId,
		Keyword:    req.Keyword,
	})
	if err != nil {
		apiReturn.ErrorDatabase(c, err.Error())
		return
	}

	apiReturn.SuccessListData(c, list, total)
}

// GetVersionList 获取微应用版本列表
func (a *MicroAppApi) GetVersionList(c *gin.Context) {
	req := MicroAppVersionGetVersionListReq{}
	if err := c.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}

	if req.Page < 1 {
		req.Page = 1
	}
	if req.Limit < 1 {
		req.Limit = 1
	}

	m := models.MicroAppVersion{}
	status := 1
	list, total, err := m.GetList(global.Db.Debug(), req.Page, req.Limit, &req.AppRecordId, &status)
	if err != nil {
		apiReturn.ErrorDatabase(c, err.Error())
		return
	}

	for k, _ := range list {
		list[k].PackageSrc = "" // 置空
		// list[k].PackageUrl = biz.MicroAppPackage.GenerateDownloadURL(v.PackageSrc)
	}

	apiReturn.SuccessListData(c, list, total)
}

func (a *MicroAppApi) Test(c *gin.Context) {
	// packageFolderName := "/Users/sunjingliang/my_code/sun-panel_group/sun-panel-micro-store/service/micro_app_upload/2026/04/03/YM-music-player-free_2.1.0_b823b477ed6febf9"
	// // 异步审核
	// go func() {
	// 	securityAuditResult, err := biz.MicroAppAudit.CodeSecurityAudit(packageFolderName, biz.SecurityAuditConfig{
	// 		PlatformURL: "http://127.0.0.1:3025",
	// 		// APIKey:          "sunapi",
	// 		APISecret:       "hYWxxDCCcM5Ma8Mt3h2H0RemTn9bTG6Q",
	// 		Timeout:         60 * time.Second,                                       // 60秒
	// 		AllowedFileExts: []string{".js", ".ts", ".jsx", ".tsx", ".mjs", ".cjs"}, // 只发送支持的文件类型
	// 		MaxFileSize:     1024 * 1024 * 10,
	// 	})
	// 	if err != nil {
	// 		global.Logger.Errorln("安全审核失败:", err)
	// 		return
	// 	}

	// 	global.Logger.Infoln("安全审核结果:", cmn.AnyToJsonStr(securityAuditResult))
	// }()

}

// UploadMicroAppIcon 上传微应用图标
// 要求：
// 1. 图标必须是正方形（宽高相等）
// 2. 文件大小不超过 512KB
// 3. 支持格式：.png, .jpg, .jpeg, .webp, .svg, .ico
func (a *MicroAppApi) UploadMicroAppIcon(c *gin.Context) {
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
		FormFileName: "iconfile",
		AgreeExtNames: []string{
			".png",
			".jpg",
			".jpeg",
			".webp",
			".svg",
			".ico",
		},
		MaxSize: 512 * 1024, // 512KB
		SaveDir: fildDir,
	})

	if err != nil {
		switch err {
		case file.ErrUploadExceedMaxSize:
			apiReturn.Error(c, "图标文件不能大于512KB")
			return
		case file.ErrUploadExtensionNameNotAllowed:
			apiReturn.ErrorByCode(c, 1301)
			return
		default:
			apiReturn.ErrorByCode(c, 1300)
			return
		}
	}

	// 检查图片是否为正方形
	if uploadInfo.Width > 0 && uploadInfo.Height > 0 && uploadInfo.Width != uploadInfo.Height {
		// 删除已上传的文件
		os.Remove(uploadInfo.FileSavePath)
		apiReturn.Error(c, "图标宽高比例必须为 1:1")
		return
	}

	pureFilePath := uploadInfo.FileSavePath[len(configUpload):]
	// 向数据库添加记录
	mFile := models.File{}
	mFile.AddFile(userInfo.ID, uploadInfo.FileOriginalName, uploadInfo.Ext, pureFilePath)

	responseData := gin.H{
		"imageUrl": global.UPLOAD_ROUTE + pureFilePath,
	}

	apiReturn.SuccessData(c, responseData)
}

// UploadMicroAppScreenshot 上传微应用截图
// 要求：
// 1. 截图必须是横屏（宽度 > 高度）
// 2. 截图比例必须为 4:3（误差 ±5%）
// 3. 文件大小不超过 2MB
// 4. 支持格式：.png, .jpg, .jpeg, .webp, .gif
func (a *MicroAppApi) UploadMicroAppScreenshot(c *gin.Context) {
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
		FormFileName: "screenshotfile",
		AgreeExtNames: []string{
			".png",
			".jpg",
			".jpeg",
			".webp",
			".gif",
		},
		MaxSize: 2 * 1024 * 1024, // 2MB
		SaveDir: fildDir,
	})

	if err != nil {
		switch err {
		case file.ErrUploadExceedMaxSize:
			apiReturn.Error(c, "截图文件不能大于2MB")
			return
		case file.ErrUploadExtensionNameNotAllowed:
			apiReturn.ErrorByCode(c, 1301)
			return
		default:
			apiReturn.ErrorByCode(c, 1300)
			return
		}
	}
	global.Logger.Debugln("文件上传成功:", cmn.AnyToJsonStr(uploadInfo))
	// 检查图片尺寸是否获取成功
	if uploadInfo.Width <= 0 || uploadInfo.Height <= 0 {
		// 删除已上传的文件
		os.Remove(uploadInfo.FileSavePath)
		apiReturn.Error(c, "无法获取图片尺寸，请检查图片格式")
		return
	}
	// 检查图片是否为横屏（宽度 > 高度）
	if uploadInfo.Width <= uploadInfo.Height {
		// 删除已上传的文件
		os.Remove(uploadInfo.FileSavePath)
		apiReturn.Error(c, "截图必须是横屏图片（宽度 > 高度）")
		return
	}

	// // 检查图片比例是否为 4:3（允许 ±5% 的误差）
	// ratio := float64(uploadInfo.Width) / float64(uploadInfo.Height)
	// targetRatio := 4.0 / 3.0 // 4:3 的理论比例
	// errorMargin := 0.05      // 5% 的误差范围

	// if ratio < targetRatio*(1-errorMargin) || ratio > targetRatio*(1+errorMargin) {
	// 	// 删除已上传的文件
	// 	os.Remove(uploadInfo.FileSavePath)
	// 	apiReturn.Error(c, "截图宽高比例必须为 4:3（误差 ±5%）")
	// 	return
	// }

	pureFilePath := uploadInfo.FileSavePath[len(configUpload):]
	// 向数据库添加记录
	mFile := models.File{}
	mFile.AddFile(userInfo.ID, uploadInfo.FileOriginalName, uploadInfo.Ext, pureFilePath)

	responseData := gin.H{
		"imageUrl": global.UPLOAD_ROUTE + pureFilePath,
	}

	// 添加尺寸信息
	if uploadInfo.Width > 0 && uploadInfo.Height > 0 {
		responseData["width"] = uploadInfo.Width
		responseData["height"] = uploadInfo.Height
	}

	apiReturn.SuccessData(c, responseData)
}
