package biz

import (
	"database/sql"
	"errors"
	"sun-panel/global"
	"sun-panel/lib/cmn"
	"sun-panel/models"
	"time"

	"gorm.io/gorm"
)

type RedeemCodeType struct {
}

type RedeemCodeParam struct {
	Title       string                       `gorm:"type:varchar(50);not null" json:"title"`            // 兑换码标题
	ExpireTime  string                       `gorm:"index:idx_expire_time,sort:desc" json:"expireTime"` // 过期时间
	RedeemType  int                          `gorm:"type:tinyint(1);not null" json:"redeemType"`        // 兑换码类型,默认：1.PRO 授权
	ReleaseType models.RedeemCodeReleaseType `gorm:"type:tinyint(1);not null" json:"releaseType"`       // 发布类型：参考：RedeemCodeRedeemTypeWeChatPay
	Status      models.RedeemCodeStatus      `gorm:"type:tinyint(1);not null" json:"status"`            // 兑换码状态：1.未使用 2.已使用 3.已过期 4.已作废
	Note        string                       `gorm:"type:varchar(255);not null" json:"note"`            // 备注
	ExtendData  map[string]interface{}       `json:"extendData"`                                        // 扩展数据,json 格式
}

// 定义错误
var (
	// 过期
	ErrRedeemCodeExpired = errors.New("expired")
	// 已使用
	ErrRedeemCodeUsed = errors.New("used")
	// 无效
	ErrRedeemCodeInvalid = errors.New("invalid")
)

// 添加兑换码
func (r *RedeemCodeType) Add(prefix string, param RedeemCodeParam) (models.RedeemCode, error) {
	code := ""
	i := 0

	// 循环查询是否重复
	for {
		code = cmn.BuildRandCode(10, "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
		code = prefix + code

		// 查重
		err := global.Db.First(&models.RedeemCode{}, "code=?", code).Error
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				break
			}
			return models.RedeemCode{}, err
		}

		if i > 100 {
			return models.RedeemCode{}, errors.New("unable to generate redemption code")
		}

		i++
	}

	expireTime, err := cmn.StrToTime(cmn.TimeFormatMode1, param.ExpireTime)
	if err != nil {
		return models.RedeemCode{}, err
	}

	redeemCode := models.RedeemCode{
		Title:       param.Title,
		Code:        code,
		ExpireTime:  expireTime,
		RedeemType:  param.RedeemType,
		ReleaseType: param.ReleaseType,
		Status:      param.Status,
		Note:        param.Note,
		ExtendData:  cmn.AnyToJsonStr(param.ExtendData),
	}
	err = global.Db.Create(&redeemCode).Error
	return redeemCode, err
}

// 核销
func (r *RedeemCodeType) WriteOff(code string, userId uint) error {
	redeemCode := models.RedeemCode{}
	err := global.Db.First(&redeemCode, "code=?", code).Error
	if err != nil {
		return err
	}

	// 验证是否过期
	if redeemCode.ExpireTime.Before(time.Now()) {
		return ErrRedeemCodeExpired
	}

	// 验证是否已使用
	if redeemCode.Status == models.RedeemCodeStatusUsed {
		return ErrRedeemCodeUsed
	}

	// 验证是否已作废
	if redeemCode.Status == models.RedeemCodeStatusInvalid {
		return ErrRedeemCodeInvalid
	}

	redeemCode.Status = models.RedeemCodeStatusUsed
	redeemCode.UserId = userId
	redeemCode.WriteOffTime = sql.NullTime{
		Time:  time.Now(),
		Valid: true,
	}
	return global.Db.Save(&redeemCode).Error
}

// 查询兑换码详情，如果过期直接修改状态
func (r *RedeemCodeType) Query(code string) (models.RedeemCode, error) {
	redeemCode := models.RedeemCode{}
	err := global.Db.First(&redeemCode, "code=?", code).Error
	if err != nil {
		return redeemCode, err
	}

	// 验证是否过期
	if redeemCode.ExpireTime.Before(time.Now()) {
		redeemCode.Status = models.RedeemCodeStatusExpired
		global.Db.Save(&redeemCode)
		return redeemCode, ErrRedeemCodeExpired
	}

	// 验证是否已使用
	if redeemCode.Status == models.RedeemCodeStatusUsed {
		return redeemCode, ErrRedeemCodeUsed
	}

	// 验证是否已作废
	if redeemCode.Status == models.RedeemCodeStatusInvalid {
		return redeemCode, ErrRedeemCodeInvalid
	}

	return redeemCode, nil
}

// 查询兑换码列表
func (r *RedeemCodeType) QueryList(param map[string]interface{}) ([]models.RedeemCode, error) {
	var redeemCodes []models.RedeemCode
	// 未兑换的排在靠前的位置
	err := global.Db.Preload("User").Where(param).Order("status asc,created_at asc").Find(&redeemCodes).Error
	return redeemCodes, err
}

// 设作废
func (r *RedeemCodeType) SetInvalidOld(code string) error {
	redeemCode := models.RedeemCode{}
	err := global.Db.First(&redeemCode, "code=?", code).Error
	if err != nil {
		return err
	}

	// 验证是否已作废
	if redeemCode.Status == models.RedeemCodeStatusInvalid {
		return ErrRedeemCodeInvalid
	}

	redeemCode.Status = models.RedeemCodeStatusInvalid
	return global.Db.Save(&redeemCode).Error
}

// 设作废
func (r *RedeemCodeType) SetInvalid(codes ...string) error {
	redeemCode := models.RedeemCode{}
	redeemCode.Status = models.RedeemCodeStatusInvalid
	return global.Db.Where("code in ?", codes).Select("Status").Updates(&redeemCode).Error
}
