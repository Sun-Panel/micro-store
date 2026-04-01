package biz

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"sun-panel/global"
	"sun-panel/models"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// MicroAppStatistics 微应用统计
type MicroAppStatistics struct {
	// RedisCache *cache.RedisCacheStruct[int]
	redisClient  *redis.Client
	redisContext context.Context
}

func (b *MicroAppStatistics) Init(rc *redis.Client) {
	b.redisClient = rc
	b.redisContext = context.Background()
}

// IncrementDownload 增加下载计数（使用 Redis Hash 缓存）
func (b *MicroAppStatistics) IncrementDownload(appId uint, userId uint, clientId string, downloadIp string) error {
	// 1. 记录下载明细到数据库
	download := models.MicroAppDownload{
		AppRecordId:    appId,
		UserId:         userId,
		ClientId:       clientId,
		DownloadIp:     downloadIp,
		DownloadDevice: "", // 可选
		DownloadClient: "", // 可选
	}

	if err := download.Create(global.Db); err != nil {
		return err
	}

	// 2. 使用 Redis Hash 计数器（高性能）
	if b.redisClient != nil {
		// Redis Hash: Key = micro_app:download, Field = app_id
		key := "micro_app:download"
		field := fmt.Sprintf("%d", appId)
		if err := b.redisClient.HIncrBy(b.redisContext, key, field, 1).Err(); err != nil {
			// Redis 失败时降级到直接更新数据库
			return b.incrementDownloadInDB(appId)
		}
	} else {
		// 无 Redis 配置，直接更新数据库
		return b.incrementDownloadInDB(appId)
	}

	return nil
}

// incrementDownloadInDB 在数据库中增加下载计数（降级方案）
func (b *MicroAppStatistics) incrementDownloadInDB(appId uint) error {
	return global.Db.Model(&models.MicroApp{}).
		Where("id = ?", appId).
		UpdateColumn("download_count", gorm.Expr("download_count + 1")).Error
}

// IncrementInstall 增加安装计数（使用 Redis Hash 缓存）
func (b *MicroAppStatistics) IncrementInstall(appId uint, versionId uint, userId uint, clientId string, publicIp string) error {
	// 1. 记录安装明细到数据库
	install := models.MicroAppInstall{
		AppRecordId: appId,
		VersionId:   versionId,
		UserId:      userId,
		ClientId:    clientId,
		PublicIp:    publicIp,
		IntranetIp:  "", // 可选
		UserIsPro:   false,
		PointValue:  0,
	}

	if err := install.Create(global.Db); err != nil {
		return err
	}

	// 2. 使用 Redis Hash 计数器（高性能）
	if b.redisClient != nil {
		// Redis Hash: Key = micro_app:install, Field = app_id
		key := "micro_app:install"
		field := fmt.Sprintf("%d", appId)
		if err := b.redisClient.HIncrBy(b.redisContext, key, field, 1).Err(); err != nil {
			// Redis 失败时降级到直接更新数据库
			return b.incrementInstallInDB(appId)
		}
	} else {
		// 无 Redis 配置，直接更新数据库
		return b.incrementInstallInDB(appId)
	}

	return nil
}

// incrementInstallInDB 在数据库中增加安装计数（降级方案）
func (b *MicroAppStatistics) incrementInstallInDB(appId uint) error {
	return global.Db.Model(&models.MicroApp{}).
		Where("id = ?", appId).
		UpdateColumn("install_count", gorm.Expr("install_count + 1")).Error
}

// SyncRedisCountersToDB 同步 Redis Hash 计数器到数据库（定时任务调用）
func (b *MicroAppStatistics) SyncRedisCountersToDB() error {
	if b.redisClient == nil {
		return errors.New("Redis not configured")
	}

	ctx := b.redisContext

	// 同步下载计数（从 Hash 中获取所有字段）
	downloadHashKey := "micro_app:download"
	downloadFields, err := b.redisClient.HKeys(ctx, downloadHashKey).Result()
	if err != nil {
		return err
	}

	// 批量获取下载计数
	downloadCounts, err := b.redisClient.HMGet(ctx, downloadHashKey, downloadFields...).Result()
	if err != nil {
		return err
	}

	// 同步到数据库
	for i, field := range downloadFields {
		appId, err := strconv.ParseUint(field, 10, 32)
		if err != nil {
			continue
		}

		countStr, ok := downloadCounts[i].(string)
		if !ok || countStr == "" {
			continue
		}

		count, err := strconv.Atoi(countStr)
		if err != nil || count == 0 {
			continue
		}

		// 更新数据库
		if err := global.Db.Model(&models.MicroApp{}).
			Where("id = ?", appId).
			UpdateColumn("download_count", gorm.Expr("download_count + ?", count)).Error; err == nil {
			// 清除 Redis Hash 字段
			b.redisClient.HDel(ctx, downloadHashKey, field)
		}
	}

	// 同步安装计数（从 Hash 中获取所有字段）
	installHashKey := "micro_app:install"
	installFields, err := b.redisClient.HKeys(ctx, installHashKey).Result()
	if err != nil {
		return err
	}

	// 批量获取安装计数
	installCounts, err := b.redisClient.HMGet(ctx, installHashKey, installFields...).Result()
	if err != nil {
		return err
	}

	// 同步到数据库
	for i, field := range installFields {
		appId, err := strconv.ParseUint(field, 10, 32)
		if err != nil {
			continue
		}

		countStr, ok := installCounts[i].(string)
		if !ok || countStr == "" {
			continue
		}

		count, err := strconv.Atoi(countStr)
		if err != nil || count == 0 {
			continue
		}

		// 更新数据库
		if err := global.Db.Model(&models.MicroApp{}).
			Where("id = ?", appId).
			UpdateColumn("install_count", gorm.Expr("install_count + ?", count)).Error; err == nil {
			// 清除 Redis Hash 字段
			b.redisClient.HDel(ctx, installHashKey, field)
		}
	}

	return nil
}

// GetRealtimeStatistics 获取实时统计（Redis Hash + 数据库）
func (b *MicroAppStatistics) GetRealtimeStatistics(appId uint) (downloadCount, installCount int, err error) {
	// 先从 Redis Hash 获取增量
	var redisDownloadCount, redisInstallCount int64

	if b.redisClient != nil {
		downloadHashKey := "micro_app:download"
		installHashKey := "micro_app:install"
		field := fmt.Sprintf("%d", appId)

		redisDownloadCount, _ = b.redisClient.HGet(b.redisContext, downloadHashKey, field).Int64()
		redisInstallCount, _ = b.redisClient.HGet(b.redisContext, installHashKey, field).Int64()
	}

	// 从数据库获取基础数据
	var app models.MicroApp
	if err := global.Db.Select("download_count, install_count").Where("id = ?", appId).First(&app).Error; err != nil {
		return 0, 0, err
	}

	// 合并 Redis 和数据库数据
	downloadCount = int(app.DownloadCount) + int(redisDownloadCount)
	installCount = int(app.InstallCount) + int(redisInstallCount)

	return downloadCount, installCount, nil
}

// GetBatchRealtimeStatistics 批量获取实时统计（使用 Hash 结构，性能更好）
func (b *MicroAppStatistics) GetBatchRealtimeStatistics(appIds []uint) (map[uint][2]int, error) {
	result := make(map[uint][2]int)

	if b.redisClient == nil {
		// 无 Redis，直接从数据库获取
		var apps []models.MicroApp
		if err := global.Db.Select("id, download_count, install_count").Where("id IN ?", appIds).Find(&apps).Error; err != nil {
			return nil, err
		}

		for _, app := range apps {
			result[app.ID] = [2]int{app.DownloadCount, app.InstallCount}
		}

		return result, nil
	}

	// 准备字段列表和映射
	fields := make([]string, len(appIds))
	fieldToIndex := make(map[string]int) // field -> index
	for i, appId := range appIds {
		field := fmt.Sprintf("%d", appId)
		fields[i] = field
		fieldToIndex[field] = i
	}

	// 批量从 Redis Hash 获取下载计数
	downloadHashKey := "micro_app:download"
	downloadValues, err := b.redisClient.HMGet(b.redisContext, downloadHashKey, fields...).Result()
	if err != nil {
		return nil, err
	}

	// 批量从 Redis Hash 获取安装计数
	installHashKey := "micro_app:install"
	installValues, err := b.redisClient.HMGet(b.redisContext, installHashKey, fields...).Result()
	if err != nil {
		return nil, err
	}

	// 从数据库获取基础数据
	var apps []models.MicroApp
	if err := global.Db.Select("id, download_count, install_count").Where("id IN ?", appIds).Find(&apps).Error; err != nil {
		return nil, err
	}

	// 合并数据
	for _, app := range apps {
		field := fmt.Sprintf("%d", app.ID)
		index, exists := fieldToIndex[field]

		// 获取 Redis 中的增量
		downloadVal := downloadValues[index]
		installVal := installValues[index]

		redisDownloadCount := int64(0)
		redisInstallCount := int64(0)

		if downloadVal != nil && exists {
			if val, ok := downloadVal.(string); ok && val != "" {
				redisDownloadCount, _ = strconv.ParseInt(val, 10, 64)
			}
		}

		if installVal != nil && exists {
			if val, ok := installVal.(string); ok && val != "" {
				redisInstallCount, _ = strconv.ParseInt(val, 10, 64)
			}
		}

		result[app.ID] = [2]int{
			app.DownloadCount + int(redisDownloadCount),
			app.InstallCount + int(redisInstallCount),
		}
	}

	return result, nil
}
