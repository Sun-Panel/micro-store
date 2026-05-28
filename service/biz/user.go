package biz

import (
	"fmt"
	"sun-panel/global"
	"sun-panel/lib/cache"
	"sun-panel/models"
	"time"
)

type userType struct {
	// 用户信息缓存，key 格式：username:{name} 或 uid:{id}
	InfoCache cache.Cacher[models.User]
}

// Init 初始化用户服务（含缓存）
func (u *userType) Init() {
	u.InfoCache = global.NewCache[models.User](30*time.Minute, 10*time.Minute, "UserInfoCache")
}

// 缓存 key 生成
func (u *userType) usernameKey(username string) string {
	return fmt.Sprintf("username:%s", username)
}

func (u *userType) uidKey(uid uint) string {
	return fmt.Sprintf("uid:%d", uid)
}

// invalidateByUsername 清除指定用户名的缓存
func (u *userType) invalidateByUsername(username string) {
	if u.InfoCache != nil {
		u.InfoCache.Delete(u.usernameKey(username))
	}
}

// invalidateByUid 清除指定用户ID的缓存
func (u *userType) invalidateByUid(uid uint) {
	if u.InfoCache != nil {
		u.InfoCache.Delete(u.uidKey(uid))
	}
}

// GetUser 根据用户名获取用户信息（带缓存）
func (u *userType) GetUser(username string) (models.User, error) {
	if u.InfoCache != nil {
		if user, ok := u.InfoCache.Get(u.usernameKey(username)); ok {
			return user, nil
		}
	}

	user := models.User{}
	result, err := user.GetUserInfoByUsername(username)
	if err != nil {
		return models.User{}, err
	}

	// 写入缓存
	if u.InfoCache != nil {
		u.InfoCache.SetDefault(u.usernameKey(username), result)
		u.InfoCache.SetDefault(u.uidKey(result.ID), result) // 同时缓存 uid 映射
	}

	return result, nil
}

// GetUserByUId 根据用户ID获取用户信息（带缓存）
func (u *userType) GetUserByUId(uid uint) (models.User, error) {
	if u.InfoCache != nil {
		if user, ok := u.InfoCache.Get(u.uidKey(uid)); ok {
			return user, nil
		}
	}

	user := models.User{}
	result, err := user.GetUserInfoByUid(uid)
	if err != nil {
		return models.User{}, err
	}

	// 写入缓存
	if u.InfoCache != nil {
		u.InfoCache.SetDefault(u.uidKey(uid), result)
		u.InfoCache.SetDefault(u.usernameKey(result.Username), result) // 同时缓存 username 映射
	}

	return result, nil
}
