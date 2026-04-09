package biz

import (
	"strings"
	"sun-panel/global"
)

type ConfigType struct {
}

func (c *ConfigType) GetCustomStylePath() string {
	// 绑定自定义风格的路由
	customPath := global.Config.GetValueString("base", "custom_style_path")
	// 兼容旧版 2024-02-06 开发
	if customPath == "" {
		customPath = strings.TrimSuffix(c.GetWebPath(), "/") + "/custom"
	}

	return customPath
}

func (c *ConfigType) GetWebPath() string {
	return global.Config.GetValueStringOrDefault("base", "web_path")
}

// 获取微应用根目录
func (c *ConfigType) GetMicroAppRootPath() string {
	path := global.Config.GetValueStringOrDefault("base", "micro_app_path")
	if path == "" {
		path = "./micro_app"
	}
	return strings.TrimSuffix(path, "/")
}

func (c *ConfigType) GetMicroAppResourcePath() string {
	return c.GetMicroAppRootPath() + "/res"
}

// // 获取微应用运行根目录
// func (c *ConfigType) GetMicroAppRunRootPath() string {
// 	return c.GetMicroAppRootPath() + "/run"
// }

// // 获取微应用运行路径
// func (c *ConfigType) GetMicroAppRunPathByMicroAppId(microAppId string) string {
// 	return c.GetMicroAppRootPath() + "/run/" + microAppId
// }

// // 获取微应用安装包根路径
// func (c *ConfigType) GetMicroAppInstallPackageRootPath() string {
// 	return c.GetMicroAppRootPath() + "/install_package/"
// }

func (c *ConfigType) GetUploadsPath() string {
	return global.Config.GetValueStringOrDefault("base", "source_path")
}

func (c *ConfigType) GetSqliteFilePath() string {
	return global.Config.GetValueStringOrDefault("sqlite", "file_path")
}

func (c *ConfigType) GetTempPath() string {
	return global.Config.GetValueString("base", "source_temp_path")
}

// 根据路径获取临时路径
func (c *ConfigType) GetTempPathByPath(path string) string {
	return c.GetTempPath() + "/" + strings.Trim(path, "/")
}
