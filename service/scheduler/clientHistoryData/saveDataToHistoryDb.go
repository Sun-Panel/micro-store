package clientHistoryData

import (
	"sun-panel/biz"
	"sun-panel/global"
	"sun-panel/lib/cmn"
	"sun-panel/models"
	"sun-panel/scheduler/clientHistoryData/syncClientData"
	"time"

	"gorm.io/gorm"
)

var activeClientMonth = 6

// 同步缓存数据到数据库并获得新注册的客户端数量
func syncCacheDataAndGetNewClientNum(startHour, endHour time.Time) int64 {
	// 同步保存新注册的客户端数据到数据库
	newClientNum := syncClientData.SyncNewClientsCacheByHourTime(startHour, endHour)

	// 同步更新的客户端数据到数据库
	updateClientNum := syncClientData.SyncHourUpdateClientsCacheByHourCacheKey(biz.ClientCache.HistoryCache.GetHourKey(startHour))
	_ = updateClientNum // 暂时没有用

	// 清理过期的小时缓存，防止内存泄露
	biz.ClientCache.HistoryCache.CleanupOldHourCaches()

	return newClientNum

}

// 保存历史数据到数据库
func saveHistoryClientStatistics(startHourTime, endHourTime time.Time) {

	var (
		newClientNum int64 = 0

		onlineNum24h int64 = -1
		onlineNum48h int64 = -1
		onlineNum72h int64 = -1

		activeClientTotalNum int64 = 0
	)

	// 查询一小时内新增的客户端数量
	global.Db.Model(&models.SoftwareClient{}).
		Where("join_time >=? and join_time <?", startHourTime, endHourTime).
		Count(&newClientNum)
	global.Logger.Debugln("一小时内新增条数:", newClientNum)

	// 查询活跃客户端数据， 默认为 activeClientMonth 个月内上过线统计
	global.Db.Model(&models.SoftwareClient{}).
		Where("join_time >= ? and join_time <?", endHourTime.AddDate(0, -activeClientMonth, 0), endHourTime).
		Count(&activeClientTotalNum)
	global.Logger.Debugln(activeClientMonth, "个月内总条数:", activeClientTotalNum)

	// 判断是否为前1小时的整点时间，是则统计以下数据
	anHourAgo := getOnHourTime(time.Now()).Add(-1 * time.Hour)
	global.Logger.Debug("startHourTime:", startHourTime, "anHourAgo:", anHourAgo)
	if startHourTime == anHourAgo {
		// 查询24小时内的数据
		startTime24h := endHourTime.Add(-24 * time.Hour)
		global.Db.Model(&models.SoftwareClient{}).
			Where("last_online_time >= ? and last_online_time < ? ", startTime24h, endHourTime).
			Count(&onlineNum24h)

		// 查询48小时内的数据
		startTime48h := endHourTime.Add(-48 * time.Hour)
		global.Db.Model(&models.SoftwareClient{}).
			Where("last_online_time >= ? and last_online_time < ? ", startTime48h, endHourTime).
			Count(&onlineNum48h)

		// 查询72小时内的数据
		startTime72h := endHourTime.Add(-72 * time.Hour)
		global.Db.Model(&models.SoftwareClient{}).
			Where("last_online_time >= ? and last_online_time < ? ", startTime72h, endHourTime).
			Count(&onlineNum72h)
	}

	newData := models.HistoryClientStatistics{
		DateTime:             startHourTime,
		OnlineNum24h:         onlineNum24h,
		OnlineNum48h:         onlineNum48h,
		OnlineNum72h:         onlineNum72h,
		HourNewClientNum:     newClientNum,
		ActiveClientTotalNum: activeClientTotalNum,
	}

	global.Db.Model(&models.HistoryClientStatistics{}).Create(&newData)
}

// 保存各个版本在某时间到现在的客户端数量
func savehistoryClientVersionStatistics(startHourTime, endHourTime time.Time) {

	var (
		onlineNum24h    = ""
		onlineNum48h    = ""
		onlineNum72h    = ""
		activeClientNum = ""
		// activeClientTotalNum = ""
	)

	db := global.Db

	// 判断是否为当前小时的整点时间，是则统计以下数据
	anHourAgo := getOnHourTime(time.Now()).Add(-1 * time.Hour)
	global.Logger.Debug("startHourTime:", startHourTime, "anHourAgo:", anHourAgo)
	if startHourTime == anHourAgo {
		mapOnlineNum24h, _ := getClientsCountByVersionWithinDuration(db, endHourTime.Add(-24*time.Hour), endHourTime)
		mapOnlineNum48h, _ := getClientsCountByVersionWithinDuration(db, endHourTime.Add(-48*time.Hour), endHourTime)
		mapOnlineNum72h, _ := getClientsCountByVersionWithinDuration(db, endHourTime.Add(-72*time.Hour), endHourTime)

		onlineNum24h = cmn.AnyToJsonStr(mapOnlineNum24h)
		onlineNum48h = cmn.AnyToJsonStr(mapOnlineNum48h)
		onlineNum72h = cmn.AnyToJsonStr(mapOnlineNum72h)
	}

	mapActiveClientNum, _ := getClientsCountByVersionWithinDuration(db, endHourTime.Add(-30*24*time.Hour), endHourTime)
	activeClientNum = cmn.AnyToJsonStr(mapActiveClientNum)

	newData := models.HistoryClientVersionStatistics{
		DateTime:        startHourTime,
		OnlineNum24h:    onlineNum24h,
		OnlineNum48h:    onlineNum48h,
		OnlineNum72h:    onlineNum72h,
		ActiveClientNum: activeClientNum,
	}

	global.Db.Create(&newData)
}

// 查询出各个版本在某时间到现在的客户端数量
func getClientsCountByVersionWithinDuration(db *gorm.DB, startHourTime, endHourTime time.Time) (map[string]int64, error) {
	var results []struct {
		Version string
		Count   int64
	}

	if err := db.Model(&models.SoftwareClient{}).
		Select("version, COUNT(*) as count").
		Where("last_online_time >= ? AND last_online_time < ?", startHourTime, endHourTime).
		Group("version").
		Find(&results).Error; err != nil {
		return nil, err
	}

	// 转为 JSON 格式
	counts := make(map[string]int64)
	for _, result := range results {
		counts[result.Version] = result.Count
	}

	return counts, nil
}

// 是否统计完成
func isStatisticsCompleted(hourTime time.Time) bool {
	var count int64
	global.Db.Model(&models.HistoryClientStatistics{}).Where("date_time =?", hourTime).Count(&count)
	return count > 0
}

// 取时间整点
func getOnHourTime(hourTime time.Time) time.Time {
	return time.Date(hourTime.Year(), hourTime.Month(), hourTime.Day(), hourTime.Hour(), 0, 0, 0, hourTime.Location())
}
