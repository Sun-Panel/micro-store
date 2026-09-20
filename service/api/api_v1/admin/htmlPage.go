package admin

import (
	"encoding/json"
	"errors"
	"regexp"
	"strings"
	"sun-panel/api/api_v1/common/apiData/commonApiStructs"
	"sun-panel/api/api_v1/common/apiData/systemApiStructs"
	"sun-panel/api/api_v1/common/apiReturn"
	"sun-panel/global"
	"sun-panel/models"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type HtmlPageApi struct {
}

var HtmlPageVarPrefix = "htmlPage."

// 页面名会作为预览页的 URL 路径使用，不允许为空、包含空白字符或斜杠
var htmlPageNameRegexp = regexp.MustCompile(`^[^\s/\\]+$`)

func checkPageName(pageName string) (string, error) {
	pageName = strings.TrimSpace(pageName)
	if pageName == "" {
		return pageName, errors.New("pageName is required")
	}
	if !htmlPageNameRegexp.MatchString(pageName) {
		return pageName, errors.New("pageName cannot contain whitespace or slash")
	}
	return pageName, nil
}

func (a *HtmlPageApi) GetList(c *gin.Context) {
	req := commonApiStructs.RequestPage{}

	if err := c.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}

	// 归一化分页参数，避免传入非法值时生成 OFFSET 为负数的 SQL
	if req.Page < 1 {
		req.Page = 1
	}
	if req.Limit < 1 || req.Limit > 200 {
		req.Limit = 10
	}

	var count int64

	varList := []models.SystemSetting{}
	dbErr := global.Db.Order("id asc").
		Limit(req.Limit).Offset((req.Page-1)*req.Limit).
		Where(
			"config_name like ? OR (config_name = ? AND description like ?) ",
			HtmlPageVarPrefix+req.Keyword+"%",
			HtmlPageVarPrefix,
			"%"+req.Keyword+"%",
		).
		Find(&varList).
		Limit(-1).Offset(-1).Count(&count).Error

	if dbErr != nil {
		apiReturn.ErrorDatabase(c, dbErr.Error())
		return
	}

	respList := []systemApiStructs.HtmlPageListItem{}
	for _, v := range varList {
		pageInfo := systemApiStructs.HtmlPageInfo{}
		json.Unmarshal([]byte(v.ConfigValue), &pageInfo)
		respList = append(respList, systemApiStructs.HtmlPageListItem{
			PageName:        v.ConfigName[len(HtmlPageVarPrefix):],
			PageDescription: v.Description,
			HtmlPageInfo:    pageInfo,
		})
	}
	apiReturn.SuccessListData(c, respList, count)
}

func (a *HtmlPageApi) GetInfo(c *gin.Context) {
	req := systemApiStructs.HtmlPageInfoEditReq{}

	if err := c.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}

	validatedName, err := checkPageName(req.PageName)
	if err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}

	pageName := HtmlPageVarPrefix + validatedName

	info := models.SystemSetting{}
	dbErr := global.Db.Order("id asc").
		First(&info, "config_name=?", pageName).Error

	if dbErr != nil {
		apiReturn.ErrorDatabase(c, dbErr.Error())
		return
	}

	pageInfo := systemApiStructs.HtmlPageInfo{}
	json.Unmarshal([]byte(info.ConfigValue), &pageInfo)
	resp := systemApiStructs.HtmlPageListItem{
		PageName:        info.ConfigName[len(HtmlPageVarPrefix):],
		PageDescription: info.Description,
		HtmlPageInfo:    pageInfo,
	}

	apiReturn.SuccessData(c, resp)
}

func (a *HtmlPageApi) Edit(c *gin.Context) {
	req := systemApiStructs.HtmlPageInfoEditReq{}

	if err := c.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}

	validatedName, err := checkPageName(req.PageName)
	if err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}

	mgName := HtmlPageVarPrefix + validatedName

	if err := global.SystemSetting.Set(mgName, req.HtmlPageInfo); err != nil {
		apiReturn.Error(c, err.Error())
	} else {
		global.Db.Model(&models.SystemSetting{}).Where("config_name=?", mgName).Update("description", req.PageDescription)
		apiReturn.Success(c)
	}

}

func (a *HtmlPageApi) Delete(c *gin.Context) {
	req := systemApiStructs.HtmlPageInfoEditReq{}

	if err := c.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}
	validatedName, err := checkPageName(req.PageName)
	if err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}

	mgName := HtmlPageVarPrefix + validatedName
	if err := global.SystemSetting.Delete(mgName); err != nil {
		apiReturn.Error(c, err.Error())
		return
	}

	apiReturn.Success(c)
}
