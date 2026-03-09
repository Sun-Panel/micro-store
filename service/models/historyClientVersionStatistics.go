package models

import "time"

type HistoryClientVersionStatistics struct {
	BaseModel
	DateTime        time.Time `json:"dateTime"` // 整小时时间
	OnlineNum24h    string    `json:"onlineNum24h" gorm:"column:online_num_24h"`
	OnlineNum48h    string    `json:"onlineNum48h" gorm:"column:online_num_48h"`
	OnlineNum72h    string    `json:"onlineNum72h" gorm:"column:online_num_72h"`
	ActiveClientNum string    `json:"activeClientNum"` // 活跃客户端的各个版本的统计数
}
