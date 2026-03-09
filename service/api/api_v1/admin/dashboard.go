package admin

import (
	"encoding/json"
	"fmt"
	"sun-panel/api/api_v1/common/apiReturn"
	"sun-panel/api/api_v1/common/base"
	"sun-panel/global"
	"sun-panel/lib/cmn"
	"sun-panel/models"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"gorm.io/gorm"
)

type DashboardApi struct{}

type Dates struct {
	Start string
	End   string
}

// 0.1.* 版本
func (a *DashboardApi) GetStatisticsOld(c *gin.Context) {
	// const cacheName:="DashboardStatistics"
	location := base.GetUserTimezoneLocation(c)
	todayTime := time.Now().In(location)
	zeroTime := time.Date(todayTime.Year(), todayTime.Month(), todayTime.Day(), 0, 0, 0, 0, todayTime.Location())

	// fmt.Println("现在时间", zeroTime.Format(cmn.TimeFormatMode1))
	todayStartTimeStr := zeroTime.UTC().Format(cmn.TimeFormatMode1)
	todayEndTimeStr := zeroTime.UTC().Add(time.Hour * 24).Format(cmn.TimeFormatMode1)

	// 总用户/今日新增
	var userCount int64
	var userToday int64
	global.Db.Model(&models.User{}).Count(&userCount)
	global.Db.Model(&models.User{}).Where("created_at BETWEEN ? AND ?", todayStartTimeStr, todayEndTimeStr).Count(&userToday)

	// 总客户端/在线数
	var clientCount int64
	var clientIncreaseCountToday int64

	global.Db.Model(&models.SoftwareClient{}).Count(&clientCount)

	// 今日新增客户端数据
	global.Db.Model(&models.SoftwareClient{}).
		Where("created_at BETWEEN ? AND ?", todayStartTimeStr, todayEndTimeStr).
		Count(&clientIncreaseCountToday)

	var clientOnline24 int64
	var clientOnline48 int64
	var clientOnline72 int64

	// 24小时到现在 (用户时区计算)
	onlineStartTime := time.Now().In(location).Add(-24 * time.Hour).Format(cmn.TimeFormatMode1)
	onlineEndTime := time.Now().In(location).Local().Format(cmn.TimeFormatMode1)
	global.Db.Model(&models.SoftwareClient{}).Where("updated_at BETWEEN ? AND ?", onlineStartTime, onlineEndTime).Count(&clientOnline24)

	// 48小时到现在 (用户时区计算)
	onlineStartTime = time.Now().In(location).Add(-48 * time.Hour).Format(cmn.TimeFormatMode1)
	onlineEndTime = time.Now().In(location).Local().Format(cmn.TimeFormatMode1)
	global.Db.Model(&models.SoftwareClient{}).Where("updated_at BETWEEN ? AND ?", onlineStartTime, onlineEndTime).Count(&clientOnline48)

	// 72小时到现在 (用户时区计算)
	onlineStartTime = time.Now().In(location).Add(-72 * time.Hour).Format(cmn.TimeFormatMode1)
	onlineEndTime = time.Now().In(location).Local().Format(cmn.TimeFormatMode1)
	global.Db.Model(&models.SoftwareClient{}).Where("updated_at BETWEEN ? AND ?", onlineStartTime, onlineEndTime).Count(&clientOnline72)

	// 取出每个版本的总客户端数
	type VersionCount struct {
		Version       string `json:"version"`
		Count         int64  `json:"count"`
		OnlineCount48 int64  `json:"onlineCount48"`
	}
	var versionCounts []VersionCount
	// global.Db.Debug().Model(&models.SoftwareClient{}).
	// 	Select("version, count(*) as count").
	// 	Group("version").
	// 	Scan(&versionCounts)

	global.Db.Model(&models.SoftwareClient{}).
		Select(
			"version, count(*) as count, SUM(CASE WHEN updated_at >= ? THEN 1 ELSE 0 END) as OnlineCount48",
			time.Now().In(location).Add(-48*time.Hour),
		).
		Group("version").
		Scan(&versionCounts)

	resp := gin.H{
		"userCount": userCount,
		"userToday": userToday,

		"clientCount":              clientCount,
		"clientIncreaseCountToday": clientIncreaseCountToday,

		"clientOnline24": clientOnline24,
		"clientOnline48": clientOnline48,
		"clientOnline72": clientOnline72,

		"versionClientCount": versionCounts,
	}

	// global.SystemSetting.Set(cacheName,resp)
	apiReturn.SuccessData(c, resp)
}

func (a *DashboardApi) GetStatistics(c *gin.Context) {
	// const cacheName:="DashboardStatistics"
	location := base.GetUserTimezoneLocation(c)
	todayTime := time.Now().In(location)
	zeroTime := time.Date(todayTime.Year(), todayTime.Month(), todayTime.Day(), 0, 0, 0, 0, todayTime.Location())

	// fmt.Println("现在时间", zeroTime.Format(cmn.TimeFormatMode1))
	todayStartTimeStr := zeroTime.UTC().Format(cmn.TimeFormatMode1)
	todayEndTimeStr := zeroTime.UTC().Add(time.Hour * 24).Format(cmn.TimeFormatMode1)

	// 总用户/今日新增
	var userCount int64
	var userToday int64
	global.Db.Model(&models.User{}).Count(&userCount)
	global.Db.Model(&models.User{}).Where("created_at BETWEEN ? AND ?", todayStartTimeStr, todayEndTimeStr).Count(&userToday)

	// 累计安装次数（总客户端）
	var installCount int64
	global.Db.Model(&models.HistoryClientStatistics{}).Select("SUM(hour_new_client_num)").Scan(&installCount)

	// 半年内活跃客户端总数
	var activeClientCount int64
	var historyClientStatistics models.HistoryClientStatistics
	global.Db.Model(&models.HistoryClientStatistics{}).Order("date_time Desc").First(&historyClientStatistics)
	activeClientCount = historyClientStatistics.ActiveClientTotalNum

	// 今日新增客户端数据
	var clientIncreaseCountToday int64 = 0
	if todayStartTimeStr != todayEndTimeStr {
		global.Db.Model(&models.HistoryClientStatistics{}).
			Select("SUM(hour_new_client_num)").
			Where("date_time BETWEEN ? AND ?", todayStartTimeStr, todayEndTimeStr).
			Scan(&clientIncreaseCountToday)
	}

	var clientOnline24 int64 = historyClientStatistics.OnlineNum24h
	var clientOnline48 int64 = historyClientStatistics.OnlineNum48h
	var clientOnline72 int64 = historyClientStatistics.OnlineNum72h

	resp := gin.H{
		"userCount": userCount,
		"userToday": userToday,

		"installCount": installCount,

		"activeClientCount":        activeClientCount,
		"clientIncreaseCountToday": clientIncreaseCountToday,

		"clientOnline24": clientOnline24,
		"clientOnline48": clientOnline48,
		"clientOnline72": clientOnline72,
	}

	apiReturn.SuccessData(c, resp)
}

// 获取活跃客户端的版本统计数据
func (a *DashboardApi) GetActiveClientVersionStatistics(c *gin.Context) {
	data := models.HistoryClientVersionStatistics{}
	global.Db.Model(&models.HistoryClientVersionStatistics{}).Order("date_time Desc").First(&data)

	var (
		activeClientNum map[string]int64
		onlineNum24h    map[string]int64
		onlineNum48h    map[string]int64
		onlineNum72h    map[string]int64
	)

	json.Unmarshal([]byte(data.ActiveClientNum), &activeClientNum)
	json.Unmarshal([]byte(data.OnlineNum24h), &onlineNum24h)
	json.Unmarshal([]byte(data.OnlineNum48h), &onlineNum48h)
	json.Unmarshal([]byte(data.OnlineNum72h), &onlineNum72h)
	apiReturn.SuccessData(c, gin.H{
		"activeClientNum": activeClientNum,
		"onlineNum24h":    onlineNum24h,
		"onlineNum48h":    onlineNum48h,
		"onlineNum72h":    onlineNum72h,
		"dateTime":        base.ConvertTimeToUserTime(c, data.DateTime.Add(1*time.Hour)),
	})
}

func (a *DashboardApi) GetVersions(c *gin.Context) {

	var (
		list  []models.Version
		count int64
	)

	if err := global.Db.Order("release_time Desc").Find(&list).Error; err != nil {
		apiReturn.ErrorDatabase(c, err.Error())
		return
	}

	resMap := []VersionInfoResp{}
	for _, v := range list {
		resMap = append(resMap, VersionInfoResp{
			ID:           v.ID,
			Version:      v.Version,
			Type:         string(v.Type),
			ReleaseTime:  v.ReleaseTime.Format(cmn.TimeFormatMode1),
			IsActive:     v.IsActive,
			IsRolledBack: v.IsRolledBack,
		})
	}

	apiReturn.SuccessListData(c, resMap, count)
}

type SQLCount struct {
	Count int `json:"count"`
}

func (a *DashboardApi) GetClientLine(c *gin.Context) {
	req := []string{}

	if err := c.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}

	location := base.GetUserTimezoneLocation(c)
	dates := ConvertToUserTimeZoneDates(location, req)
	resp := []int{}

	for _, d := range dates {
		count := SQLCount{}
		_ = d

		global.Db.Model(&models.SoftwareClient{}).
			Select("count(*) as count").
			Where("created_at BETWEEN ? AND ?", d.Start, d.End).
			Scan(&count)
		resp = append(resp, count.Count)
	}
	apiReturn.SuccessData(c, resp)
}

func (a *DashboardApi) GetUserLine(c *gin.Context) {
	req := []string{}

	if err := c.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		apiReturn.ErrorParamFomat(c, err.Error())
		return
	}

	location := base.GetUserTimezoneLocation(c)
	dates := ConvertToUserTimeZoneDates(location, req)
	resp := []int{}

	for _, d := range dates {
		count := SQLCount{}
		_ = d

		global.Db.Model(&models.User{}).
			Select("count(*) as count").
			Where("created_at BETWEEN ? AND ?", d.Start, d.End).
			Scan(&count)
		resp = append(resp, count.Count)
	}
	apiReturn.SuccessData(c, resp)
}

// 获取版本的历史记录（30天每天00点的数据）+当前最新的时间数据
func (a *DashboardApi) GetVersionHistory(c *gin.Context) {
	// 获取用户的时区，并检查是否成功
	userLocation := base.GetUserTimezoneLocation(c)

	// 获取用户时区的00点时间
	now := time.Now().In(userLocation)
	zeroTime := time.Date(
		now.Year(),
		now.Month(),
		now.Day(),
		0, 0, 0, 0,
		userLocation, // 直接指定用户时区
	)

	// 计算出服务器时区的00点时间
	serverZeroTime := zeroTime.In(time.Local)

	// 格式化为字符串
	serverZeroTimeStr := serverZeroTime.Format("15:04:05")

	// 获取历史数据并检查错误
	data, _ := GetHistoryClientVersionStatisticsByHour(global.Db, 30, serverZeroTimeStr)

	// 获取最新的一条
	lastData := models.HistoryClientVersionStatistics{}
	global.Db.Limit(1).Order("date_time desc").First(&lastData)
	data = append(data, lastData)

	respList := []map[string]any{}
	for i := 0; i < len(data); i++ {
		item := data[i]
		var (
			activeClientNum map[string]int64
			onlineNum24h    map[string]int64
			onlineNum48h    map[string]int64
			onlineNum72h    map[string]int64
		)

		json.Unmarshal([]byte(item.ActiveClientNum), &activeClientNum)
		json.Unmarshal([]byte(item.OnlineNum24h), &onlineNum24h)
		json.Unmarshal([]byte(item.OnlineNum48h), &onlineNum48h)
		json.Unmarshal([]byte(item.OnlineNum72h), &onlineNum72h)

		respList = append(respList, map[string]any{
			"activeClientNum": activeClientNum,
			"onlineNum24h":    onlineNum24h,
			"onlineNum48h":    onlineNum48h,
			"onlineNum72h":    onlineNum72h,
			"dateTime":        base.ConvertTimeToUserTime(c, item.DateTime),
		})

	}
	// 返回成功的结果
	apiReturn.SuccessListData(c, respList, int64(len(respList)))
}

// 获取历史客户端版本统计数据（指定某个小时）
func GetHistoryClientVersionStatisticsByHour(db *gorm.DB, days int, hour string) ([]models.HistoryClientVersionStatistics, error) {
	var logs []models.HistoryClientVersionStatistics
	startDate := time.Now().AddDate(0, 0, -days)
	err := db.Where("date_time >= ?", startDate).
		Where("TIME(date_time) = ?", hour).
		Order("date_time ASC").
		Find(&logs).Error
	return logs, err
}

func ConvertToUserTimeZoneDates(location *time.Location, dateStrings []string) []Dates {
	var dates []Dates

	for _, dateString := range dateStrings {
		startTime, err := time.ParseInLocation("2006-01-02 15:04:05", dateString, location)
		if err != nil {
			fmt.Println("Error parsing start time:", err)
			continue
		}
		// global.Logger.Debugln("user zone", startTime, startTime.Unix())
		startTime = startTime.Local()
		// global.Logger.Debugln("local zone", startTime, startTime.Unix())
		endTime := startTime.Add(24 * time.Hour)

		date := Dates{
			Start: startTime.Format("2006-01-02 15:04:05"),
			End:   endTime.Format("2006-01-02 15:04:05"),
		}

		dates = append(dates, date)
	}

	return dates
}
