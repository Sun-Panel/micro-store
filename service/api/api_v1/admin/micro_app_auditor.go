package admin

import (
	"encoding/json"
	"sun-panel/api/api_v1/common/apiReturn"
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

// ReviewApp 审核微应用主信息（审核员专用）
func (a *MicroAppAuditorApi) ReviewApp(c *gin.Context) {
	param := MicroAppReviewAppReq{}
	if err := c.ShouldBindBodyWith(&param, binding.JSON); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}

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

	// 获取当前审核员ID（从token中获取）
	auditorId := c.GetUint("adminId")

	// 开启事务
	err = global.Db.Transaction(func(tx *gorm.DB) error {
		// 更新审核记录状态
		now := time.Now()
		if err := tx.Model(&models.MicroAppReview{}).Where("id = ?", param.ReviewId).Updates(map[string]interface{}{
			"status":       param.Status,
			"reviewer_id":  auditorId,
			"review_note":  param.ReviewNote,
			"review_time":  now,
		}).Error; err != nil {
			return err
		}

		// 如果审核通过，将审核快照数据更新到主表
		if param.Status == 1 {
			// 反序列化多语言信息
			var langMap map[string]MicroAppLangInfo
			json.Unmarshal([]byte(review.LangMap), &langMap)

			// 更新主表
			if err := tx.Model(&models.MicroApp{}).Where("id = ?", review.AppId).Updates(map[string]interface{}{
				"app_name":      review.AppName,
				"app_icon":      review.AppIcon,
				"app_desc":      review.AppDesc,
				"category_id":   review.CategoryId,
				"charge_type":   review.ChargeType,
				"price":         review.Price,
				"screenshots":   review.Screenshots,
				"remark":        review.Remark,
				"review_status": 2, // 已通过
				"review_time":   &now,
			}).Error; err != nil {
				return err
			}

			// 更新多语言信息
			m := models.MicroApp{}
			app, _ := m.GetById(tx, review.AppId)
			if app.MicroAppId != "" {
				// 先删除当前的多语言记录
				if err := tx.Where("micro_app_id = ?", app.MicroAppId).Delete(&models.MicroAppLang{}).Error; err != nil {
					return err
				}

				// 然后插入新的多语言记录
				for lang, langInfo := range langMap {
					// 先尝试查找已存在的记录（包括软删除的）
					var existingLang models.MicroAppLang
					err := tx.Unscoped().Where("micro_app_id = ? AND lang = ?", app.MicroAppId, lang).First(&existingLang).Error

					if err == nil {
						// 记录已存在，更新它（包括软删除的记录）
						existingLang.AppName = langInfo.AppName
						existingLang.AppDesc = langInfo.AppDesc
						existingLang.DeletedAt = gorm.DeletedAt{} // 取消软删除标记
						if err := tx.Save(&existingLang).Error; err != nil {
							return err
						}
					} else {
						// 记录不存在，创建新记录
						langModel := models.MicroAppLang{
							MicroAppId: app.MicroAppId,
							Lang:       lang,
							AppName:    langInfo.AppName,
							AppDesc:    langInfo.AppDesc,
						}
						if err := tx.Create(&langModel).Error; err != nil {
							return err
						}
					}
				}
			}
		} else {
			// 审核拒绝
			if err := tx.Model(&models.MicroApp{}).Where("id = ?", review.AppId).Updates(map[string]interface{}{
				"review_status": 3, // 已拒绝
				"review_time":   &now,
			}).Error; err != nil {
				return err
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
