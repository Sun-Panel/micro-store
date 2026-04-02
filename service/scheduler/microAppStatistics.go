package scheduler

import (
	"sun-panel/biz"
	"sun-panel/global"
)

// SyncMicroAppStatisticsToDB 同步微应用统计数据到数据库的定时任务
func SyncMicroAppStatisticsToDB() {
	global.Logger.Infoln("开始同步微应用统计数据到数据库")
	
	err := biz.MicroAppStatistics.SyncRedisCountersToDB()
	if err != nil {
		global.Logger.Errorln("同步微应用统计数据失败:", err)
	} else {
		global.Logger.Infoln("同步微应用统计数据成功")
	}
}
