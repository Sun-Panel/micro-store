package admin

import (
	"sun-panel/api/api_v1/common/apiReturn"
	"sun-panel/api/api_v1/common/base"
	"sun-panel/biz"
	"sun-panel/global"
	"sun-panel/models"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"gorm.io/gorm"
)

// MicroAppAuditorApi 审核员专用微应用 API
type MicroAppAuditorApi struct{}

// GetReviewPendingList 获取待审核应用列表（审核员专用）
func (a *MicroAppAuditorApi) GetReviewPendingList(c *gin.Context) {
	param := MicroAppGetPendingReviewListReq{}
	if err := c.ShouldBindBodyWith(&param, binding.JSON); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}

	reviewModel := models.MicroAppReview{}
	if param.Page <= 0 {
		param.Page = 1
	}
	if param.Limit <= 0 {
		param.Limit = 10
	}

	list, total, err := reviewModel.GetPendingList(global.Db, param.Page, param.Limit)
	if err != nil {
		apiReturn.ErrorDatabase(c, err.Error())
		return
	}

	apiReturn.SuccessListData(c, list, total)
}

// GetReviewInfo 获取审核详情（审核员专用）
func (a *MicroAppAuditorApi) GetReviewInfo(c *gin.Context) {
	param := MicroAppGetReviewInfoReq{}
	if err := c.ShouldBindBodyWith(&param, binding.JSON); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}

	reviewModel := models.MicroAppReview{}
	review, err := reviewModel.GetById(global.Db, param.ReviewId)
	if err != nil {
		apiReturn.ErrorDataNotFound(c)
		return
	}

	apiReturn.SuccessData(c, review)
}

// GetMicroAppInfo 获取微应用信息
func (a *MicroAppAuditorApi) GetMicroAppInfo(c *gin.Context) {
	req := models.MicroApp{}
	if err := c.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}

	microApp, err := biz.MicroApp.GetById(global.Db, req.ID, "Developer")
	if err != nil {
		base.HandleBizErrorAndReturn(c, err)
		return
	}

	apiReturn.SuccessData(c, microApp)
}

// ReviewApp 审核微应用主信息（审核员专用）
func (a *MicroAppAuditorApi) ReviewApp(c *gin.Context) {
	param := MicroAppReviewAppReq{}
	if err := c.ShouldBindBodyWith(&param, binding.JSON); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}

	userInfo, _ := base.GetCurrentUserInfo(c)

	// 验证审核状态
	if param.Status != 1 && param.Status != 2 {
		apiReturn.ErrorParamFomat(c, "审核状态无效")
		return
	}

	// 获取审核记录
	reviewModel := models.MicroAppReview{}
	review, err := reviewModel.GetById(global.Db, param.ReviewId)
	if err != nil {
		apiReturn.ErrorDataNotFound(c)
		return
	}

	// 检查是否已审核
	if review.Status != 0 {
		apiReturn.Error(c, "该记录已审核")
		return
	}

	// 开启事务
	err = global.Db.Transaction(func(tx *gorm.DB) error {
		// 更新审核记录状态
		now := time.Now()
		if err := tx.Model(&models.MicroAppReview{}).Where("id = ?", param.ReviewId).Updates(map[string]interface{}{
			"status":      param.Status,
			"reviewer_id": userInfo.ID,
			"review_note": param.ReviewNote,
			"review_time": now,
		}).Error; err != nil {
			return err
		}

		// 如果审核通过，更新 micro_app 表,并自动上架
		if param.Status == 1 {
			// 更新生效版本
			if err := tx.Model(&models.MicroApp{}).Where("id = ?", review.AppRecordId).Updates(map[string]interface{}{
				// "app_name":    review.AppName,
				"app_icon": review.AppIcon,
				// "app_desc":    review.AppDesc,
				"category_id":  review.CategoryId,
				"charge_type":  review.ChargeType,
				"points":       review.Points,
				"screenshots":  review.Screenshots,
				"remark":       review.Remark,
				"status":       1, // 上架
				"third_charge": review.ThirdCharge,
				"have_iframe":  review.HaveIframe,
			}).Error; err != nil {
				return err
			}

			// 更新多语言信息
			// review.LangMap 现在是 datatype.MapJson (map[string]interface{})
			for lang, langInfo := range review.LangMap {
				// 类型断言获取 map[string]interface{}
				infoMap, ok := langInfo.(map[string]interface{})
				if !ok {
					continue
				}
				appName, _ := infoMap["appName"].(string)
				appDesc, _ := infoMap["appDesc"].(string)

				// 查找或创建语言记录
				var existLang models.MicroAppLang
				err := tx.Where("micro_app_id = ? AND lang = ?", review.MicroAppId, lang).First(&existLang).Error
				if err == gorm.ErrRecordNotFound {
					if appName != "" || appDesc != "" {
						langModel := models.MicroAppLang{
							MicroAppId: review.MicroAppId,
							Lang:       lang,
							AppName:    appName,
							AppDesc:    appDesc,
						}
						if err := tx.Create(&langModel).Error; err != nil {
							return err
						}
					}
				} else if err == nil {
					if err := tx.Model(&models.MicroAppLang{}).Where("id = ?", existLang.ID).Updates(map[string]interface{}{
						"app_name": appName,
						"app_desc": appDesc,
					}).Error; err != nil {
						return err
					}
				} else {
					return err
				}
			}
		}

		return nil
	})

	if err != nil {
		apiReturn.ErrorDatabase(c, err.Error())
		return
	}

	apiReturn.Success(c)
}
