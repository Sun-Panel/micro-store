package microapp

import (
	"sun-panel/api/api_v1/common/apiReturn"
	"sun-panel/api/api_v1/common/base"
	"sun-panel/biz"
	"sun-panel/global"
	"sun-panel/models"

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
