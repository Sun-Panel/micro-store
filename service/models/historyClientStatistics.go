package models

import "time"

type HistoryClientStatistics struct {
	BaseModel
	DateTime             time.Time `json:"dateTime"` // 整小时时间
	OnlineNum24h         int64     `json:"onlineNum24h" gorm:"column:online_num_24h"`
	OnlineNum48h         int64     `json:"onlineNum48h" gorm:"column:online_num_48h"`
	OnlineNum72h         int64     `json:"onlineNum72h" gorm:"column:online_num_72h"`
	HourNewClientNum     int64     `json:"hourNewClientNum" `    // 小时新增客户端数据
	ActiveClientTotalNum int64     `json:"activeClientTotalNum"` // 活跃客户端数量
}
