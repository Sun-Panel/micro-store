package models

import "time"

// 更改记录
// (增加天数，或者减少天数将被记录)
type ChangeProRecord struct {
	BaseModel
	UserId      uint
	ExpiredTime time.Time // 当前的过期时间
	DayNum      int       `gorm:"int(11)"`      // 变动天数（可以为负数）
	Note        string    `gorm:"varcar(2000)"` // （购买[显示订单号]、人工[充值说明]、退款）
	AdminNote   string    `gorm:"varcar(2000)"` // 管理员备注
	OrderNo     string    `gorm:"varcar(50)"`
}
