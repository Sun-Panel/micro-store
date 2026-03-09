package clientCache

import (
	"sun-panel/global"
	"sun-panel/lib/cache"
	"time"
)

type PingRateCache struct {
	Cache cache.Cacher[PingRateCacheRechord]
}

type PingRateCacheRechord struct {
	Time time.Time
}

func NewPingRateCache() *PingRateCache {
	interval := global.Config.GetValueInt("base", "client_ping_processing_interval")
	global.Logger.Debugln("client_ping_processing_interval:", interval)
	return &PingRateCache{
		Cache: global.NewCache[PingRateCacheRechord](
			time.Duration(interval)*time.Hour,
			time.Hour*24*30,
			"ping_rate_cache"),
	}
}

func (pr *PingRateCache) GetLastPingRechord(clientId string, lanIp string) bool {
	_, ok := pr.Cache.Get(pr.Key(clientId, lanIp))
	return ok
}

func (pr *PingRateCache) SetLastPingRechord(clientId string, lanIp string) {
	pr.Cache.SetDefault(pr.Key(clientId, lanIp), PingRateCacheRechord{Time: time.Now()})
}

func (pr *PingRateCache) Key(clientId string, lanIp string) string {
	return clientId + "-lip" + lanIp
}
