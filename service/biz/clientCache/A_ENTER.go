package clientCache

type ClientCacheType struct {
	BlacklistIP   *BlacklistIP
	AccountOnline *AccountOnline
	HistoryCache  *HistoryCacheStruct
	PingRateCache *PingRateCache
}

// 初始化所有的客户端缓存
func (c *ClientCacheType) Init() {
	c.BlacklistIP = NewBlacklistIP()
	c.AccountOnline = NewAccountOnline()
	c.HistoryCache = NewHistoryCache()
	c.PingRateCache = NewPingRateCache()
}
