package biz

import (
	"errors"
	"fmt"
	"sun-panel/global"
	"sun-panel/lib/cache"
	"sun-panel/models"
	"time"

	"gorm.io/gorm"
)

type ProAuthorizeType struct {
	UserProAuthCache cache.Cacher[models.ProAuthorize]
}

func (a *ProAuthorizeType) Init() {
	a.UserProAuthCache = global.NewCache[models.ProAuthorize](48*time.Hour, time.Hour, "user_pro_auth_cache")
}

// 根据天数改变结束时间
func (a *ProAuthorizeType) ChangeExpiredTimeByDayNum(userId uint, changeDayNum int, note string, adminNote string, orderNo string) error {
	a.UserProAuthCache.Delete(a.proAuthorizeCacheKey(userId)) // 删除缓存
	return global.Db.Transaction(func(tx *gorm.DB) error {
		// 查询用户的最后到期时间
		// 如果到期时间大于等于当前的时间，则使用这个时间作为本次的起始时间
		// 否则使用当前的时间
		userProAuthorize := models.ProAuthorize{}
		currentExpiredTime := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.Local).Add(24 * time.Hour)
		if err := tx.Order("created_at DESC").First(&userProAuthorize, "user_id=?", userId).Error; err != nil {
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				return err
			}
			currentExpiredTime = currentExpiredTime.Add(time.Duration(changeDayNum) * 24 * time.Hour)

			userProAuthorize.UserId = userId
			userProAuthorize.ExpiredTime = currentExpiredTime
			// 没有记录创建记录
			if err := tx.Create(&userProAuthorize).Error; err != nil {
				return err
			}
		} else {

			// 当前的时间在会员期之后
			// 将当前时间重置为未来的时间
			if userProAuthorize.ExpiredTime.After(currentExpiredTime) {
				currentExpiredTime = userProAuthorize.ExpiredTime
			}

			// 将用户的到期时间加（减）上本次增加的值
			currentExpiredTime = currentExpiredTime.Add(time.Duration(changeDayNum) * 24 * time.Hour)
			userProAuthorize.ExpiredTime = currentExpiredTime

			// 更新 PRO 的过期时间
			if err := tx.Updates(&userProAuthorize).Error; err != nil {
				return err
			}
		}

		// 添加到ChangeProRecord记录中
		newRecord := models.ChangeProRecord{
			UserId:      userId,
			ExpiredTime: currentExpiredTime,
			DayNum:      changeDayNum,
			Note:        note,
			OrderNo:     orderNo,
			AdminNote:   adminNote,
		}
		if err := tx.Create(&newRecord).Error; err != nil {
			return err
		}

		return nil
	})
}

func (a *ProAuthorizeType) GetAuthorizeHistoryRecordByUserId(userId uint) ([]models.ChangeProRecord, error) {
	list := []models.ChangeProRecord{}
	if err := global.Db.Order("created_at DESC").Find(&list, "user_id=?", userId).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (a *ProAuthorizeType) GetChangeRecordByOrderNo(orderNo string) (models.ChangeProRecord, error) {
	info := models.ChangeProRecord{}
	if err := global.Db.Order("created_at DESC").First(&info, "order_no=?", orderNo).Error; err != nil {
		return info, err
	}
	return info, nil
}

func (a *ProAuthorizeType) GetProAuthorizeByUserIdFromDb(userId uint) (models.ProAuthorize, error) {
	info := models.ProAuthorize{}
	if err := global.Db.First(&info, "user_id=?", userId).Error; err != nil {
		return info, err
	}
	return info, nil
}

func (a *ProAuthorizeType) proAuthorizeCacheKey(userId uint) string {
	return fmt.Sprintf("u_%d", userId)
}

func (a *ProAuthorizeType) GetProAuthorizeByUserId(userId uint) (models.ProAuthorize, error) {
	info, exists := a.UserProAuthCache.Get(a.proAuthorizeCacheKey(userId))
	if !exists {
		infoDb, err := a.GetProAuthorizeByUserIdFromDb(userId)
		if err != nil {
			return infoDb, err
		}
		info = infoDb
		a.UserProAuthCache.SetDefault(a.proAuthorizeCacheKey(userId), infoDb)
	}
	return info, nil
}
