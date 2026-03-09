package clientCache

import (
	"fmt"
	"sort"
	"sun-panel/global"
	"sun-panel/lib/cache"
	"sun-panel/lib/cmn"
	"sun-panel/models"
	"sync"
	"time"
)

type AccountOnlineCacheInfo struct {
	IP        string
	LanIP     string
	Timestamp int64
	Ctoken    string // 临时token
	ClientID  string
	Mac       string
}

type AccountOnline struct {
	mu    sync.Mutex
	Cache cache.Cacher[MapAccountOnlineCacheInfo]
}

type MapAccountOnlineCacheInfo map[string]AccountOnlineCacheInfo

func NewAccountOnline() *AccountOnline {
	return &AccountOnline{
		Cache: global.NewCache[MapAccountOnlineCacheInfo](-1, time.Hour*24*30, "ClientAccountOnlineDb"),
	}
}

// Deprecated: 请使用 SetClient 来替代.
// 添加登录客户端
func (c *AccountOnline) AddClient(userId uint, clientId string, loginInfo AccountOnlineCacheInfo) error {
	return c.SetClient(userId, clientId, loginInfo)
}

// 添加/设置登录客户端
func (c *AccountOnline) SetClient(userId uint, clientId string, loginInfo AccountOnlineCacheInfo) error {
	cachekey := fmt.Sprintf("user_id_%d", userId)

	global.Logger.Debugln("添加登录客户端", userId, clientId, loginInfo)

	c.mu.Lock()         // 在修改缓存之前加锁
	defer c.mu.Unlock() // 确保在函数返回时解锁

	cacheValue, _ := c.Cache.Get(cachekey)
	if cacheValue == nil {
		global.Logger.Debugln("获取缓存-失败，创建")
		cacheValue = make(MapAccountOnlineCacheInfo)
	}

	global.Logger.Debugln("获取缓存", cmn.AnyToJsonStr(cacheValue))

	if oldLofinInfo, ok := cacheValue[clientId]; ok {
		// 存在则删除旧 ctoken 的记录，阻止再次使用旧ctoken
		c.DeleteCtoken(oldLofinInfo.Ctoken)
	}
	cacheValue[clientId] = loginInfo

	// 限制同时登录数量
	maxNum := 2 // 暂时写死，后期可以配置
	if len(cacheValue) > maxNum {
		global.Logger.Debugln("超出最大登录数量，开始退出登录")
		global.Logger.Debugln("增加本次登录信息后的总数据", cmn.AnyToJsonStr(cacheValue))
		sortSlice := c.SortToKeySlice(cacheValue)
		global.Logger.Debugln("排序结果", sortSlice)
		for _, v := range sortSlice {
			// 退出登录
			c.DeleteClient(userId, v)
			delete(cacheValue, v)
			global.Logger.Debugln("退出客户端", userId, v)
			// 退出登录后，判断是否小于最大登录数量，小于则退出循环
			global.Logger.Debugln("是否终止循环", len(cacheValue) <= maxNum)
			if len(cacheValue) <= maxNum {
				break
			}
		}
	}

	global.Logger.Debugln("保存的值", cmn.AnyToJsonStr(cacheValue))
	c.Cache.SetDefault(cachekey, cacheValue)
	return nil
}

func (c *AccountOnline) GetClient(userId uint, clientId string) (AccountOnlineCacheInfo, bool) {
	enptyV := AccountOnlineCacheInfo{}
	cachekey := fmt.Sprintf("user_id_%d", userId)
	cacheValue, _ := c.Cache.Get(cachekey)

	if cacheValue == nil {
		return enptyV, false
	}

	if v, ok := cacheValue[clientId]; ok {
		return v, true
	}

	return enptyV, false
}

// 删除客户端，同时删除ctoken
func (c *AccountOnline) DeleteClient(userId uint, clientId string) error {
	cachekey := fmt.Sprintf("user_id_%d", userId)

	// c.mu.Lock()         // 在修改缓存之前加锁
	// defer c.mu.Unlock() // 确保在函数返回时解锁

	cacheValue, _ := c.Cache.Get(cachekey)

	if cacheValue == nil {
		return nil
	}

	if _, ok := cacheValue[clientId]; !ok {
		return nil
	}
	c.DeleteCtoken(cacheValue[clientId].Ctoken)
	c.Cache.SetDefault(cachekey, cacheValue)

	return nil
}

// 根据时间戳把map的key从最早到最新排序返回切片
func (c *AccountOnline) SortToKeySlice(value map[string]AccountOnlineCacheInfo) []string {
	var keys []string
	for k := range value {
		keys = append(keys, k)
	}

	sort.SliceStable(keys, func(i, j int) bool {
		return value[keys[i]].Timestamp < value[keys[j]].Timestamp
	})

	return keys
}

// 添加ctoken
func (c *AccountOnline) AddCtokenByUserInfo(cToken string, userInfo models.User) {
	bToken := fmt.Sprintf("user_id_%d", userInfo.ID)

	global.UserAuthServiceClientToken.SetDefault(bToken, global.AuthServiceClientTokenUser{
		User:   userInfo,
		Ctoken: cToken,
	})

	global.CUserAuthServiceClientToken.SetDefault(cToken, bToken)
}

// 删除ctoken，让ctoken失效
func (c *AccountOnline) DeleteCtoken(cToken string) {
	global.CUserAuthServiceClientToken.Delete(cToken)
}

// // 登录
// func (c *AccountOnline) Login(ctoken string, userId uint, loginInfo AccountOnlineCacheInfo) error {
// 	cachekey := fmt.Sprintf("user_id_%d", userId)

// 	c.mu.Lock()         // 在修改缓存之前加锁
// 	defer c.mu.Unlock() // 确保在函数返回时解锁

// 	cacheValue, _ :=c.Cache.Get(cachekey)
// 	// 限制同时登录数量
// 	maxNum := 2 // 暂时写死，后期可以配置
// 	if len(cacheValue) >= maxNum+1 {
// 		sortSlice := c.SortToKeySlice(cacheValue)
// 		for _, v := range sortSlice {
// 			// 退出登录
// 			delete(cacheValue, v)
// 			global.CUserAuthServiceClientToken.Delete(v)

// 			// 退出登录后，判断是否小于最大登录数量，小于则退出循环
// 			if len(cacheValue) < maxNum {
// 				break
// 			}
// 		}
// 	}
// 	cacheValue[ctoken] = loginInfo
// 	global.ClientAccountOnlineCache.SetDefault(cachekey, cacheValue)
// 	return nil
// }

// // 退出登录
// func (c *AccountOnline) Logout(ctoken string, userId uint) error {
// 	cachekey := fmt.Sprintf("user_id_%d", userId)

// 	c.mu.Lock()         // 在修改缓存之前加锁
// 	defer c.mu.Unlock() // 确保在函数返回时解锁

// 	cacheValue, _ :=c.Cache.Get(cachekey)
// 	if _, ok := cacheValue[ctoken]; !ok {
// 		return nil
// 	}
// 	delete(cacheValue, ctoken)
// 	global.ClientAccountOnlineCache.SetDefault(cachekey, cacheValue)
// 	return nil
// }
