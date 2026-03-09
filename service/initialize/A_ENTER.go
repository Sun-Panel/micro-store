package initialize

import (
	"fmt"
	"sun-panel/biz"
	"sun-panel/global"
	"sun-panel/initialize/authService"
	"sun-panel/initialize/cUserToken"
	"sun-panel/initialize/config"
	"sun-panel/initialize/database"
	"sun-panel/initialize/lang"
	"sun-panel/initialize/redis"
	"sun-panel/initialize/runlog"
	"sun-panel/initialize/systemSettingCache"
	"sun-panel/initialize/userToken"
	"sun-panel/lib/cmn"
	"sun-panel/models"
	"sun-panel/scheduler"
	"sun-panel/structs"

	"log"

	"github.com/gin-gonic/gin"
)

func InitAuthServiceApp() error {
	Logo()
	gin.SetMode(global.RUNCODE) // GIN 运行模式

	// 日志
	if logger, err := runlog.InitRunlog(global.RUNCODE, "runtime/runlog/"); err != nil {
		log.Panicln("Log initialization error", err)
		panic(err)
	} else {
		global.Logger = logger
	}

	// 配置初始化
	{
		if config, err := config.ConfigInit(); err != nil {
			global.Logger.Errorln("Configuration initialization error", err)
			return err
		} else {
			global.Config = config
		}
	}

	// 多语言初始化
	lang.LangInit("zh-cn") // en-us

	DatabaseConnect()

	// Redis 连接
	{
		// 判断是否有使用redis的驱动，没有将不连接
		cacheDrive := global.Config.GetValueString("base", "cache_drive")
		queueDrive := global.Config.GetValueString("base", "queue_drive")
		if cacheDrive == "redis" || queueDrive == "redis" {
			redisConfig := structs.IniConfigRedis{}
			global.Config.GetSection("redis", &redisConfig)
			rdb, err := redis.InitRedis(redis.Options{
				Addr:     redisConfig.Address,
				Password: redisConfig.Password,
				DB:       redisConfig.Db,
			})

			if err != nil {
				log.Panicln("Redis initialization error", err)
				panic(err)
				// return err
			}
			global.RedisDb = rdb
		}
	}

	// 初始化WEB站用户token
	global.UserToken = userToken.InitUserToken()
	global.CUserToken = cUserToken.InitCUserToken()
	global.CUserApiTokenAccessToken = cUserToken.InitCUserApiTokenAccessToken()
	global.CUserAccessTokenApiToken = cUserToken.InitCUserAccessTokenApiToken()

	// 初始化客户端api用户token
	global.UserAuthServiceClientToken = authService.InitUserAuthServiceClientToken()
	global.CUserAuthServiceClientToken = authService.InitCUserAuthServiceClientToken()
	global.ClientAccountOnlineCache = authService.InitClientAccountOnlineCache()
	global.ClientLoginAttemptsCacheCache = authService.InitClientLoginAttemptsCache()

	biz.InitBIZ()
	global.SystemSetting = systemSettingCache.InItSystemSettingCache()

	// 开启调度中心(根据配置是否开启)
	if global.Config.GetValueString("base", "scheduler_enable") == "true" {
		go scheduler.Start()
	}

	return nil
}

func DatabaseConnect() {
	// 数据库连接 - 开始
	var dbClientInfo database.DbClient
	databaseDrive := global.Config.GetValueStringOrDefault("base", "database_drive")
	if databaseDrive == database.MYSQL {
		dbClientInfo = &database.MySQLConfig{
			Username:    global.Config.GetValueStringOrDefault("mysql", "username"),
			Password:    global.Config.GetValueStringOrDefault("mysql", "password"),
			Host:        global.Config.GetValueStringOrDefault("mysql", "host"),
			Port:        global.Config.GetValueStringOrDefault("mysql", "port"),
			Database:    global.Config.GetValueStringOrDefault("mysql", "db_name"),
			WaitTimeout: global.Config.GetValueInt("mysql", "wait_timeout"),
		}
	} else {
		dbClientInfo = &database.SQLiteConfig{
			Filename: global.Config.GetValueStringOrDefault("sqlite", "file_path"),
		}
	}

	if db, err := database.DbInit(dbClientInfo); err != nil {
		log.Panicln("Database initialization error", err)
		panic(err)
	} else {
		global.Db = db
		models.Db = global.Db
	}

	database.CreateDatabase(databaseDrive, global.Db)

	database.NotFoundAndCreateUser(global.Db)
}

func Logo() {
	fmt.Println("===============================")
	fmt.Println("Sun-Panel-Server")
	fmt.Println("-------------------------------")
	versionInfo := cmn.GetSysVersionInfo()
	fmt.Println("Version:", versionInfo.Version)
	fmt.Println("Welcome to the Sun-Panel-Server.")
	fmt.Println("===============================")
}
