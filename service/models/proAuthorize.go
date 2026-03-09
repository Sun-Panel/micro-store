package models

import (
	"time"
)

// PRO
type ProAuthorize struct {
	BaseModel
	UserId      uint
	ExpiredTime time.Time // 过期日期
}
