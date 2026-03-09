package models

type BrowserExtensionClient struct {
	BaseModel
	ClientId      string `gorm:"type:varchar(100);not null;index" json:"clientId"` // 客户端id
	Version       string `gorm:"type:varchar(100);not null" json:"version"`        // 软件版本号
	LastIp        string `gorm:"type:varchar(50);not null" json:"lastIp"`          // 最后登录ip
	SystemLang    string `gorm:"type:char(10)" json:"systemLang"`                  // 系统语言 zh_cn 手动设置、注册和登录更新(只允许系统支持的语言)
	Lang          string `gorm:"type:char(10)" json:"lang"`                        // 语言 zh_cn 注册和登录更新
	TimeZone      string `gorm:"type:char(50)" json:"timeZone"`
	LastTimeStamp int    `json:"lastTimeStamp"` // 上次更新时间戳
	// UserId uint `gorm:"not null;index" json:"userId"` // 绑定的用户
	// User   User `json:"user"`
}
