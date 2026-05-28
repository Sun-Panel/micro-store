package apiClientApp

import (
	v1 "sun-panel/api/apiClientApp"
	"sun-panel/api/apiClientApp/middleware"

	"github.com/gin-gonic/gin"
)

func Init(routerGroup *gin.RouterGroup) {
	InitApiRouter(routerGroup)
}

// InitApiRouter 初始化客户端API路由
func InitApiRouter(router *gin.RouterGroup) {
	api := v1.ApiGroupApp.Api

	// 客户端API路由 中间件：生成responseId并记录错误日志，解析版本和秘钥，尝试解析用户信息，但不拦截
	r := router.Group("clientApp", middleware.ResponseLogger, middleware.ParseVersionWithSecretKey)
	{
		// ==========================
		// 公开接口（尝试解析用户信息，但不拦截）
		// ==========================
		public := r.Group("", middleware.ResolveCurrentUser)
		{
			public.POST("developer/batchGetInfo", api.BatchGetDeveloperInfo)
			public.POST("microApp/batchGetInfo", api.BatchGetMicroAppInfo)
		}

		// ==========================
		// 涉及需要登录但不强制登录的接口
		// ==========================
		auth := r.Group("", middleware.ResolveCurrentUser)
		{
			// 微应用相关
			auth.POST("microApp/installRecord", api.RecordInstall)
		}

		// ==========================
		// 涉及需要登录的接口，不登录将拦截
		// ==========================
		loginRouter := r.Group("", middleware.LoginInterceptor)
		{
			// 开发者相关
			loginRouter.POST("developer/checkIsDeveloper", api.CheckIsDeveloper)
		}

	}
}
