package clientCache

import (
	"sun-panel/global"
	"sun-panel/lib/cache"
	"sun-panel/models"
	"time"
)

type HistoryCacheStruct struct {
	newClientsCache        cache.Cacher[HistoryCacheDataStruct]            // 记录新的客户端缓存
	hourUpdateClientsCache map[string]cache.Cacher[HistoryCacheDataStruct] // 记录每小时更新的客户端缓存
	hourKeyPrefix          string
}

// 缓存过期时间配置
const (
	// 新客户端缓存过期时间：5小时（足够定时任务同步）
	newClientsCacheExpiration = 5 * time.Hour
	// 小时更新缓存过期时间：3小时
	hourUpdateCacheExpiration = 3 * time.Hour
	// 清理间隔
	cacheCleanupInterval = 30 * time.Minute
	// 保留的最大小时缓存数量（超过此数量的旧缓存将被清理）
	maxHourCacheKeep = 3
)

type HistoryCacheDataStruct struct {
	ClientId       string    `json:"clientId"`
	Version        string    `json:"version"`
	MacAddress     string    `json:"macAddress"`
	LastIp         string    `json:"lastIp"`
	LastLanIp      string    `json:"lastLanIp"`
	JoinTime       time.Time `json:"joinTime"`
	LastOnlineTime time.Time `json:"lastOnlineTime"`

	UserId uint `json:"userId"`
}

func NewHistoryCache() *HistoryCacheStruct {
	// 当前小时的缓存区
	currentCacheKey := getCurrentHourKey()
	hourUpdateClientsCache := make(map[string]cache.Cacher[HistoryCacheDataStruct])
	hourUpdateClientsCache[currentCacheKey] = createUpdateClientCacher(currentCacheKey)
	return &HistoryCacheStruct{
		// 修复：设置合理的过期时间，避免内存泄露
		// 原来是 -1（永不过期），现在设置为 2 小时过期
		newClientsCache:        global.NewCache[HistoryCacheDataStruct](newClientsCacheExpiration, cacheCleanupInterval, "NewClientsCache"),
		hourUpdateClientsCache: hourUpdateClientsCache,
		hourKeyPrefix:          getHourKeyPrefix(),
	}
}

// 新注册的客户端，写入到新客户端的缓存区
func (h *HistoryCacheStruct) WriteToNewClientsCache(clientId string, client models.SoftwareClient) {
	saveData := HistoryCacheDataStruct{
		ClientId:       clientId,
		Version:        client.Version,
		MacAddress:     client.MacAddress,
		LastIp:         client.LastIp,
		LastLanIp:      client.LastLanIp,
		JoinTime:       client.JoinTime,
		LastOnlineTime: client.LastOnlineTime,
		UserId:         client.UserId,
	}
	h.newClientsCache.SetDefault(clientId, saveData)
}

func (h *HistoryCacheStruct) ReadNewClientsCacheClientInfo(clientId string) (models.SoftwareClient, bool) {
	getdata, ok := h.newClientsCache.Get(clientId)
	if !ok {
		return models.SoftwareClient{}, false
	}
	return models.SoftwareClient{
		ClientId:       getdata.ClientId,
		Version:        getdata.Version,
		MacAddress:     getdata.MacAddress,
		LastIp:         getdata.LastIp,
		LastLanIp:      getdata.LastLanIp,
		JoinTime:       getdata.JoinTime,
		LastOnlineTime: getdata.LastOnlineTime,
		UserId:         getdata.UserId,
	}, true
}

func (h *HistoryCacheStruct) ReadNewClientsCacheAllClientIds() ([]string, error) {
	return h.newClientsCache.GetKeys()
}

// 客户端ID是否存在客户端的缓存中
func (h *HistoryCacheStruct) ClientIdInNewClientCache(clientId string) bool {
	if ok, _ := h.newClientsCache.KeyExist(clientId); ok {
		return true
	}
	return false
}

// 删除新客户端的缓存区的key(redis中指field)，一般是把数据全部同步到数据库之后的调用
func (h *HistoryCacheStruct) DeleteNewClientIDKeys(keys ...string) {
	for _, k := range keys {
		h.newClientsCache.Delete(k)
	}
}

func (h *HistoryCacheStruct) GetHourKey(hourTime time.Time) string {
	return h.hourKeyPrefix + hourTime.Format("2006010215")
}

func (h *HistoryCacheStruct) HourUpdateClientsCacheExist(hourCacheKey string) bool {
	_, ok := h.hourUpdateClientsCache[hourCacheKey]
	return ok
}

func (h *HistoryCacheStruct) ReadHourUpdateClientsCacheClientInfo(hourCacheKey, clientId string) (models.SoftwareClient, bool) {
	if !h.HourUpdateClientsCacheExist(hourCacheKey) {
		return models.SoftwareClient{}, false
	}
	getdata, ok := h.hourUpdateClientsCache[hourCacheKey].Get(clientId)
	if !ok {
		return models.SoftwareClient{}, false
	}
	return models.SoftwareClient{
		ClientId:       getdata.ClientId,
		Version:        getdata.Version,
		MacAddress:     getdata.MacAddress,
		LastIp:         getdata.LastIp,
		LastLanIp:      getdata.LastLanIp,
		JoinTime:       getdata.JoinTime,
		LastOnlineTime: getdata.LastOnlineTime,
		UserId:         getdata.UserId,
	}, true
}

func (h *HistoryCacheStruct) ReadHourUpdateClientsCacheAllClientIds(hourCacheKey string) ([]string, error) {
	// 刚启动时无法得知（redis）是否有缓冲区，所以先创建缓冲区
	if !h.HourUpdateClientsCacheExist(hourCacheKey) {
		h.hourUpdateClientsCache[hourCacheKey] = createUpdateClientCacher(hourCacheKey)
	}

	return h.hourUpdateClientsCache[hourCacheKey].GetKeys()
}

// 更新客户端数据,写入到更新缓存区
func (h *HistoryCacheStruct) WriteToHourUpdateClientsCache(clientId string, client models.SoftwareClient) {
	// global.Logger.Debugln("WriteToHourUpdateClientsCache:", client)
	currentCacheKey := getCurrentHourKey()

	// 缓存不存在创建新的缓存区
	if !h.HourUpdateClientsCacheExist(currentCacheKey) {
		h.hourUpdateClientsCache[currentCacheKey] = createUpdateClientCacher(currentCacheKey)
	}

	saveData := HistoryCacheDataStruct{
		ClientId:       clientId,
		Version:        client.Version,
		MacAddress:     client.MacAddress,
		LastIp:         client.LastIp,
		LastLanIp:      client.LastLanIp,
		JoinTime:       client.JoinTime,
		LastOnlineTime: client.LastOnlineTime,
		UserId:         client.UserId,
	}

	h.hourUpdateClientsCache[currentCacheKey].SetDefault(clientId, saveData)
}

// 删除某个小时的更新缓存区的clientId，一般是把数据全部同步到数据库之后的调用
func (h *HistoryCacheStruct) DeleteHourUpdateClients(hourCacheKey string, clientIds ...string) {
	if !h.HourUpdateClientsCacheExist(hourCacheKey) {
		return
	}

	for _, k := range clientIds {
		h.hourUpdateClientsCache[hourCacheKey].Delete(k)
	}
}

// 删除整个小时的更新缓存区，一般是吧数据全部同步到数据库之后的调用
func (h *HistoryCacheStruct) DeleteHourUpdateClientsCache(hourCacheKey string) {
	if _, ok := h.hourUpdateClientsCache[hourCacheKey]; ok {
		global.Logger.Debugln("缓存区所有的clientids", hourCacheKey)
		h.hourUpdateClientsCache[hourCacheKey].Flush()
		delete(h.hourUpdateClientsCache, hourCacheKey)
	}
}

func (h *HistoryCacheStruct) GetCurrentHourKey() string {
	return getCurrentHourKey()
}

func getCurrentHourKey() string {
	return getHourKeyPrefix() + time.Now().Format("2006010215")
}

func getHourKeyPrefix() string {
	return "HourUpdateClientsCache:"
}

func createUpdateClientCacher(hourCacheKey string) cache.Cacher[HistoryCacheDataStruct] {
	opt := global.CacherOption{
		CacheAreaExpired: 48 * time.Hour, // 缓存区设置48小时过期（仅对Redis有效）
	}
	// 修复：设置合理的过期时间，避免内存泄露
	// 原来是 -1（永不过期），现在设置为 3 小时过期
	return global.NewCache[HistoryCacheDataStruct](hourUpdateCacheExpiration, cacheCleanupInterval, hourCacheKey, opt)
}

// CleanupOldHourCaches 清理超过指定小时数的旧缓存
// 防止 hourUpdateClientsCache map 无限增长导致内存泄露
func (h *HistoryCacheStruct) CleanupOldHourCaches() {
	// 获取当前小时的整点时间
	now := time.Now()
	currentHour := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), 0, 0, 0, now.Location())

	// 需要删除的缓存 key
	keysToDelete := make([]string, 0)

	for key := range h.hourUpdateClientsCache {
		// 解析 key 中的时间（格式：HourUpdateClientsCache:2006010215）
		if len(key) <= len(h.hourKeyPrefix) {
			continue
		}
		timeStr := key[len(h.hourKeyPrefix):]
		cacheTime, err := time.Parse("2006010215", timeStr)
		if err != nil {
			global.Logger.Warnln("解析缓存时间失败:", timeStr, err)
			continue
		}

		// 如果缓存时间超过保留的最大小时数，则标记删除
		hoursDiff := currentHour.Sub(cacheTime).Hours()
		if hoursDiff > float64(maxHourCacheKeep) {
			keysToDelete = append(keysToDelete, key)
		}
	}

	// 删除旧缓存
	for _, key := range keysToDelete {
		if cacher, ok := h.hourUpdateClientsCache[key]; ok {
			cacher.Flush()
			delete(h.hourUpdateClientsCache, key)
			global.Logger.Debugln("清理过期小时缓存:", key)
		}
	}

	if len(keysToDelete) > 0 {
		global.Logger.Infoln("清理了", len(keysToDelete), "个过期的小时缓存")
	}
}
