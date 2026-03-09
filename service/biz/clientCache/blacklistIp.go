package clientCache

import (
	"sun-panel/global"
	"sun-panel/lib/cache"
	"time"
)

type CacheParams struct {
	LastVisitTimestamp  int64 // 最后访问时间戳
	ExpirationTimestamp int64 // 过期时间戳
}

type BlacklistIP struct {
	Cache cache.Cacher[CacheParams]
}

func NewBlacklistIP() *BlacklistIP {
	return &BlacklistIP{
		Cache: global.NewCache[CacheParams](-1, time.Hour*500, "ClientBlacklistIpDb"),
	}
}

func (c *BlacklistIP) Check(ip string) (isBlacklist bool) {
	p, ok := c.Cache.Get(ip)

	if !ok {
		return false
	}

	if p.ExpirationTimestamp < time.Now().Unix() {
		c.Cache.Delete(ip)
		return false
	}

	return true
}

func (c *BlacklistIP) GetInfo(ip string) (CacheParams, bool) {
	return c.Cache.Get(ip)
}

func (c *BlacklistIP) CheckAndUpdate(ip string) (isBlacklist bool) {
	p, ok := c.Cache.Get(ip)
	if !ok {
		return false
	}

	if p.ExpirationTimestamp < time.Now().Unix() {
		c.Cache.Delete(ip)
		return false
	}
	p.LastVisitTimestamp = time.Now().Unix()
	c.Cache.SetDefault(ip, p)
	return true
}

func (c *BlacklistIP) GetAll() map[string]CacheParams {
	d := map[string]CacheParams{}
	keys, _ := c.Cache.GetKeys()

	for _, ip := range keys {
		p, ok := c.Cache.Get(ip)
		if !ok || p.ExpirationTimestamp < time.Now().Unix() {
			c.Cache.Delete(ip)
			continue
		}

		d[ip] = p
	}

	return d
}

func (c *BlacklistIP) Set(expiration time.Time, ip ...string) (err error) {

	for _, v := range ip {
		var lastVisitTimestamp int64 = 0
		if p, ok := c.Cache.Get(v); ok {
			lastVisitTimestamp = p.LastVisitTimestamp
		}
		c.Cache.SetDefault(v, CacheParams{
			LastVisitTimestamp:  lastVisitTimestamp,
			ExpirationTimestamp: expiration.Unix(),
		})
	}
	return
}

func (c *BlacklistIP) Del(ip ...string) (err error) {
	for _, v := range ip {
		c.Cache.Delete(v)
	}
	return
}
