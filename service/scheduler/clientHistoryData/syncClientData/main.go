package syncClientData

import (
	"sun-panel/biz"
	"sun-panel/global"
	"sun-panel/lib/cmn"
	"sun-panel/models"
	"time"

	"gorm.io/gorm"
)

// 同步指定某个小时的新客户端历史数据
func SyncNewClientsCacheByHourTime(startHourTime, endHourTime time.Time) int64 {
	clientIds, err := biz.ClientCache.HistoryCache.ReadNewClientsCacheAllClientIds()
	global.Logger.Debugln("clientIds:", clientIds, err)
	if err != nil {
		return 0
	}

	var count int64 = 0

	for _, clientId := range clientIds {
		clientInfo, ok := biz.ClientCache.HistoryCache.ReadNewClientsCacheClientInfo(clientId)
		global.Logger.Debugln("clientInfo:", clientInfo, ok)
		if !ok {
			continue
		}
		global.Logger.Debugln("clientInfo:", cmn.AnyToJsonStr(clientInfo))

		// 判断当前时间是否在指定小时内，是则保存到数据库中并删除此条记录
		if IsTimeBetween(clientInfo.JoinTime, startHourTime, endHourTime) {
			global.Logger.Debugln("在该时间范围更新:", clientInfo.JoinTime, startHourTime.Add(time.Hour*1), cmn.AnyToJsonStr(clientInfo))
			// 保存到数据库中，如果存在将采用更新
			if addOrUpdateErr := addClientInfoToDb(clientInfo); addOrUpdateErr != nil {
				global.Logger.Errorln("syncNewClientsCacheByHourTime to db:", addOrUpdateErr, " | ", "data:", cmn.AnyToJsonStr(clientInfo))
			}
			// 删除此条缓存记录
			biz.ClientCache.HistoryCache.DeleteNewClientIDKeys(clientId)
			count++
		}
	}

	return count
}

// 同步指定某个小时的更新客户端历史数据
func SyncHourUpdateClientsCacheByHourCacheKey(hourCacheKey string) int64 {

	clientIds, err := biz.ClientCache.HistoryCache.ReadHourUpdateClientsCacheAllClientIds(hourCacheKey)
	if err != nil {
		return 0
	}

	var count int64 = 0

	for _, clientId := range clientIds {
		clientInfo, ok := biz.ClientCache.HistoryCache.ReadHourUpdateClientsCacheClientInfo(hourCacheKey, clientId)
		if !ok {
			continue
		}

		global.Logger.Debugln("update--clientInfo:", cmn.AnyToJsonStr(clientInfo))

		// // 保存到数据库中，如果存在将采用更新
		// if addOrUpdateErr := addOrUpdateClientInfo(clientInfo); addOrUpdateErr != nil {
		// 	global.Logger.Errorln("syncNewClientsCacheByHourTime to db:", addOrUpdateErr, " | ", "data:", cmn.AnyToJsonStr(clientInfo))
		// }

		global.Logger.Debugln("准备更新:", cmn.AnyToJsonStr(clientInfo))

		// 更新到数据库
		err := global.Db.Model(&models.SoftwareClient{}).
			Where("client_id=?", clientId).
			Omit("JoinTime").
			Updates(clientInfo).Error

		if err != nil {
			global.Logger.Errorln(err)
		}

		// 删除此条缓存记录
		biz.ClientCache.HistoryCache.DeleteHourUpdateClients(hourCacheKey, clientId)
		count++
	}

	// 删除所有并删除整个缓存区
	biz.ClientCache.HistoryCache.DeleteHourUpdateClientsCache(hourCacheKey)

	if !biz.ClientCache.HistoryCache.HourUpdateClientsCacheExist(hourCacheKey) {
		global.Logger.Debugln("已成功删除整个缓存区", hourCacheKey)
	}

	return count
}

func IsTimeBetween(t, start, end time.Time) bool {
	return (t.After(start) || t.Equal(start)) && (t.Before(end) || t.Equal(end))
}

func addClientInfoToDb(clientInfo models.SoftwareClient) error {
	// 首先尝试找到记录
	var existingClientInfo models.SoftwareClient
	result := global.Db.Where("client_id = ?", clientInfo.ClientId).First(&existingClientInfo)

	// 如果记录存在，则更新(不更新，创建就是创建，已存在就忽略)
	// if result.Error == nil {
	// 	// 如果你需要更新特定字段，可以在这里设置
	// 	// result = global.Db.Omit("JoinTime").Updates(&existingClientInfo)
	// } else
	if result.Error == gorm.ErrRecordNotFound {
		// 如果记录不存在，则创建
		result = global.Db.Create(&clientInfo)
	}

	return result.Error

}
