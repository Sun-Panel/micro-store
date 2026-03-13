package clientHistoryData

func Start() {
	StartStatisticsHourAgo(1)
}

// 从 hour 小时前开始统计，如果当前为5:05，那么一小时前的数据就是4-5点，开始时间=当前小时-1小时，结束时间=当前小时的整点.
// hourAgo计算时间为当前小时向前数hourAgo个小时
// 如果当前时间为3:05,hourAgo=1的时候，那么统计时间范围就是从2-3点
func StartStatisticsHourAgo(hourAgo int) {

	// // 以当前时间为准，向前获取24个小时的整点数据数组
	// now := time.Now()

	// // 遍历24小时，获取每个整点时间
	// for i := hourAgo; i > 0; i-- {

	// 	// 计算当前整点时间
	// 	hour := now.Add(-time.Duration(i-1) * time.Hour)
	// 	// 只保留整点，去除分钟和秒
	// 	hour = time.Date(hour.Year(), hour.Month(), hour.Day(), hour.Hour(), 0, 0, 0, hour.Location())

	// 	startHour := hour.Add(-1 * time.Hour)
	// 	endHour := hour

	// 	global.Logger.Infoln("==== Start statistics history clients:", startHour, "-", endHour)

	// 	// 执行前查询每个小时的数据是否存在，存在将不执行统计了
	// 	if isStatisticsCompleted(startHour) {
	// 		global.Logger.Infoln("统计数据已存在，跳过统计：", hour)
	// 		continue
	// 	}

	// 	// ==========
	// 	// 开始统计
	// 	// ==========
	// 	global.Logger.Infoln("Start sync data:", startHour, "-", endHour)
	// 	// 同步数据
	// 	syncCacheDataAndGetNewClientNum(startHour, endHour)

	// 	global.Logger.Infoln("Start statistics history clients:", startHour, "-", endHour)
	// 	// 保存历史客户端统计数据
	// 	saveHistoryClientStatistics(startHour, endHour)

	// 	global.Logger.Infoln("Start statistics history clients version:", startHour, "-", endHour)
	// 	// 保存历史客户端版本统计的数据
	// 	savehistoryClientVersionStatistics(startHour, endHour)

	// 	global.Logger.Infoln("==== End statistics history clients:", startHour, "-", endHour)
	// }

}
