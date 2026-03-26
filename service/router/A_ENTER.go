package router

import (
	"sun-panel/global"

	"sun-panel/router/admin"
	microApp "sun-panel/router/micro_app"
	"sun-panel/router/oauth2"
	"sun-panel/router/openness"
	"sun-panel/router/panel"
	sunstore "sun-panel/router/sunStore"
	"sun-panel/router/system"

	"github.com/gin-gonic/gin"
)

// 初始化总路由
func InitRouters(addr string) error {
	router := gin.Default()
	rootRouter := router.Group("/")
	routerGroup := rootRouter.Group("api")

	// 接口
	system.Init(routerGroup)
	admin.Init(routerGroup)
	panel.Init(routerGroup)
	microApp.Init(routerGroup)
	openness.Init(routerGroup)
	oauth2.Init(routerGroup)

	// SunStore
	sunStoreRouterGroup := rootRouter.Group("sunStore/api")
	sunstore.InitWebhook(sunStoreRouterGroup)
	sunstore.InitApi(sunStoreRouterGroup)

	// WEB文件服务
	{
		webPath := "./web"
		// 上传
		sourcePath := global.Config.GetValueString("base", "source_path")
		router.Static("/uploads", sourcePath)
		// 自定义风格文件夹
		customPath := global.Config.GetValueString("base", "custom_style_path")
		// 兼容旧版
		if customPath == "" {
			router.Static("/custom", webPath+"/custom")
		} else {
			router.Static("/custom", customPath)
		}
		router.StaticFile("/", webPath+"/index.html")
		router.Static("/assets", webPath+"/assets")
		router.StaticFile("/favicon.ico", webPath+"/favicon.ico")
		router.StaticFile("/favicon.svg", webPath+"/favicon.svg")

		router.NoRoute(func(c *gin.Context) {
			c.File(webPath + "/index.html")
		})
	}

	global.Logger.Info("Sun-Panel-Server is Started.  Listening and serving HTTP on 0.0.0.0", addr)
	return router.Run(addr)
}
