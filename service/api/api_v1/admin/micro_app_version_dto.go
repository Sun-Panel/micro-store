package admin

import (
	"encoding/json"
	"sun-panel/models"
)

// ==================== 微应用配置 JSON 类型定义 ====================

// MicroAppConfig 微应用配置文件结构
type MicroAppConfig struct {
	AppJsonVersion string                    `json:"appJsonVersion"` // JSON 版本
	MicroAppId     string                    `json:"microAppId"`     // 应用唯一标识
	Version        string                    `json:"version"`        // 版本号
	APIVersion     string                    `json:"apiVersion"`     // API 版本
	Author         string                    `json:"author"`         // 作者
	Entry          string                    `json:"entry"`          // 入口文件
	Icon           string                    `json:"icon"`           // 图标
	Debug          bool                      `json:"debug"`          // 调试模式
	Components     *ComponentsConfig         `json:"components"`     // 组件配置
	Permissions    []string                  `json:"permissions"`    // 权限列表
	DataNodes      map[string]DataNodeConfig `json:"dataNodes"`      // 数据节点配置
	NetworkDomains []string                  `json:"networkDomains"` // 网络域名白名单
	AppInfo        map[string]AppInfo        `json:"appInfo"`        // 应用信息
}

// ComponentsConfig 组件配置
type ComponentsConfig struct {
	Pages   map[string]PageWidgetConfig `json:"pages"`   // 页面组件
	Widgets map[string]WidgetConfig     `json:"widgets"` // 小部件组件
}

// PageWidgetConfig 页面/小部件配置
type PageWidgetConfig struct {
	Component         interface{} `json:"component"`         // 组件类（JS/TS）
	Background        string      `json:"background"`        // 背景
	HeaderTextColor   string      `json:"headerTextColor"`   // 头部文字颜色
	Width             interface{} `json:"width"`             // 宽度
	Height            interface{} `json:"height"`            // 高度
	ShowFullscreenBtn bool        `json:"showFullscreenBtn"` // 显示全屏按钮
	Resize            bool        `json:"resize"`            // 允许调整大小
	Move              bool        `json:"move"`              // 允许移动
	Type              string      `json:"type"`              // 页面类型
}

// WidgetConfig 小部件配置
type WidgetConfig struct {
	Component           interface{} `json:"component"`           // 组件类
	ConfigComponentName string      `json:"configComponentName"` // 配置页面名称
	Size                []string    `json:"size"`                // 支持的尺寸
	Background          string      `json:"background"`          // 背景
	IsModifyBackground  bool        `json:"isModifyBackground"`  // 是否可修改背景
}

// DataNodeConfig 数据节点配置
type DataNodeConfig struct {
	Scope    string `json:"scope"`    // 作用域: app/user
	IsPublic bool   `json:"isPublic"` // 是否公开
}

// AppInfo 应用信息
type AppInfo struct {
	AppName            string `json:"appName"`            // 应用名称
	Description        string `json:"description"`        // 描述
	NetworkDescription string `json:"networkDescription"` // 网络请求说明
}

// ==================== 版本上传响应 ====================

// MicroAppVersionUploadResp 版本上传响应
type MicroAppVersionUploadResp struct {
	URL           string                       `json:"url"`           // 文件访问 URL
	Hash          string                       `json:"hash"`          // 文件 MD5 校验值
	Config        models.MicroAppVersionConfig `json:"config"`        // 解析出的配置文件
	FileName      string                       `json:"fileName"`      // 文件名
	FileSize      int64                        `json:"fileSize"`      // 文件大小
	FolderName    string                       `json:"folderName"`    // 缓存文件夹名（不含路径）
	IconURL       string                       `json:"iconURL"`       // 图标访问 URL
	UploadCacheId string                       `json:"uploadCacheId"` // 上传缓存ID
}

// MarshalJSON 自定义 JSON 序列化
func (m *MicroAppConfig) MarshalJSON() ([]byte, error) {
	type Alias MicroAppConfig
	return json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(m),
	})
}
