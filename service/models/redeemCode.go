package models

import (
	"database/sql"
	"time"
)

type RedeemCodeReleaseType int
type RedeemCodeStatus int

// 兑换码发布类型
var (
	RedeemCodeReleaseTypeWeChatPay RedeemCodeReleaseType = 1 // 微信支付
	RedeemCodeReleaseTypeAliPay    RedeemCodeReleaseType = 2 // 支付宝支付
	RedeemCodeReleaseTypeRebate    RedeemCodeReleaseType = 3 // 贡献回赠
	RedeemCodeReleaseTypeOther     RedeemCodeReleaseType = 4 // 其他
)

// 兑换码状态
var (
	RedeemCodeStatusNotUsed RedeemCodeStatus = 1 // 未使用
	RedeemCodeStatusUsed    RedeemCodeStatus = 2 // 已使用
	RedeemCodeStatusExpired RedeemCodeStatus = 3 // 已过期
	RedeemCodeStatusInvalid RedeemCodeStatus = 4 // 已作废
)

// 兑换码
type RedeemCode struct {
	BaseModel
	Title        string                `gorm:"type:varchar(50);not null" json:"title"`            // 兑换码标题
	Code         string                `gorm:"type:varchar(50);not null;index" json:"code"`       // 兑换码
	ExpireTime   time.Time             `gorm:"index:idx_expire_time,sort:desc" json:"expireTime"` // 过期时间
	RedeemType   int                   `gorm:"type:tinyint(1);not null" json:"redeemType"`        // 兑换码类型,默认：1.PRO 授权
	ReleaseType  RedeemCodeReleaseType `gorm:"type:tinyint(1);not null" json:"releaseType"`       // 发布类型：参考：RedeemCodeRedeemTypeWeChatPay
	Status       RedeemCodeStatus      `gorm:"type:tinyint(1);not null" json:"status"`            // 兑换码状态：1.未使用 2.已使用 3.已过期 4.已作废
	Note         string                `gorm:"type:varchar(255);not null" json:"note"`            // 备注
	ExtendData   string                `json:"extendData"`                                        // 扩展数据,json 格式
	WriteOffTime sql.NullTime          `json:"writeOffTime"`                                      // 核验/兑换时间

	UserId uint `json:"userId"` // 绑定的用户(兑换后会绑定用户ID)
	User   User
}
