package license

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// Middleware Gin 中间件
type Middleware struct {
	validator      *Validator
	timeValidator  *TimeValidator
	heartbeat      *HeartbeatManager
	config         *Config
	skipPaths      []string
	featureCheck   bool
	timeCheck      bool
	heartbeatCheck bool
}

// MiddlewareOption 中间件配置选项
type MiddlewareOption func(*Middleware)

// NewMiddleware 创建 License 中间件
func NewMiddleware(validator *Validator, config *Config, opts ...MiddlewareOption) *Middleware {
	m := &Middleware{
		validator:      validator,
		timeValidator:  NewTimeValidator(config),
		config:         config,
		featureCheck:   true,
		timeCheck:      true,
		heartbeatCheck: false, // 默认不检查心跳（需要手动开启）
	}

	for _, opt := range opts {
		opt(m)
	}

	return m
}

// WithSkipPaths 设置跳过验证的路径
func WithSkipPaths(paths []string) MiddlewareOption {
	return func(m *Middleware) {
		m.skipPaths = paths
	}
}

// WithFeatureCheck 设置是否检查功能权限
func WithFeatureCheck(enable bool) MiddlewareOption {
	return func(m *Middleware) {
		m.featureCheck = enable
	}
}

// WithTimeCheck 设置是否检查时间篡改
func WithTimeCheck(enable bool) MiddlewareOption {
	return func(m *Middleware) {
		m.timeCheck = enable
	}
}

// WithHeartbeatCheck 设置是否检查心跳
func WithHeartbeatCheck(enable bool) MiddlewareOption {
	return func(m *Middleware) {
		m.heartbeatCheck = enable
	}
}

// WithHeartbeatManager 设置心跳管理器
func WithHeartbeatManager(hb *HeartbeatManager) MiddlewareOption {
	return func(m *Middleware) {
		m.heartbeat = hb
	}
}

// Handler 返回 Gin 中间件处理函数
func (m *Middleware) Handler() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 检查是否跳过路径
		if m.shouldSkip(c.Request.URL.Path) {
			c.Next()
			return
		}

		// 1. 验证 License 基本有效性
		result := m.validator.Validate()
		if !result.Valid {
			c.JSON(http.StatusForbidden, gin.H{
				"code":    403,
				"message": "license invalid: " + result.Message,
				"error":   "LICENSE_INVALID",
			})
			c.Abort()
			return
		}

		// 2. 检查时间篡改
		if m.timeCheck {
			timeResult := m.timeValidator.ValidateTime()
			if !timeResult.Valid {
				c.JSON(http.StatusForbidden, gin.H{
					"code":    403,
					"message": "time validation failed: " + timeResult.Message,
					"error":   "TIME_TAMPERED",
				})
				c.Abort()
				return
			}
		}

		// 3. 检查心跳（如果启用）
		if m.heartbeatCheck && m.heartbeat != nil {
			hbResult := m.heartbeat.CheckHeartbeat()
			if !hbResult.Success {
				c.JSON(http.StatusForbidden, gin.H{
					"code":    403,
					"message": "heartbeat check failed: " + hbResult.Message,
					"error":   "HEARTBEAT_FAILED",
				})
				c.Abort()
				return
			}
		}

		// 将 License 信息注入到上下文
		c.Set("license", m.validator.GetInfo())
		c.Set("licenseID", m.validator.GetLicenseID())
		c.Set("licenseType", m.validator.GetType())

		c.Next()
	}
}

// FeatureRequired 要求特定功能权限的中间件
func (m *Middleware) FeatureRequired(feature string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !m.validator.HasFeature(feature) {
			c.JSON(http.StatusForbidden, gin.H{
				"code":    403,
				"message": "feature not licensed: " + feature,
				"error":   "FEATURE_NOT_LICENSED",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}

// FeaturesRequired 要求多个功能权限的中间件（任意一个）
func (m *Middleware) FeaturesRequired(features ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !m.validator.HasAnyFeature(features...) {
			c.JSON(http.StatusForbidden, gin.H{
				"code":    403,
				"message": "none of the features licensed: " + strings.Join(features, ", "),
				"error":   "FEATURES_NOT_LICENSED",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}

// AllFeaturesRequired 要求所有功能权限的中间件
func (m *Middleware) AllFeaturesRequired(features ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !m.validator.HasAllFeatures(features...) {
			c.JSON(http.StatusForbidden, gin.H{
				"code":    403,
				"message": "not all features licensed: " + strings.Join(features, ", "),
				"error":   "ALL_FEATURES_REQUIRED",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}

// TypeRequired 要求特定授权类型的中间件
func (m *Middleware) TypeRequired(types ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		licenseType := m.validator.GetType()
		for _, t := range types {
			if licenseType == t {
				c.Next()
				return
			}
		}
		c.JSON(http.StatusForbidden, gin.H{
			"code":    403,
			"message": "license type not allowed",
			"error":   "LICENSE_TYPE_NOT_ALLOWED",
		})
		c.Abort()
	}
}

// ExpiringWarning 即将过期警告中间件
func (m *Middleware) ExpiringWarning(days int) gin.HandlerFunc {
	return func(c *gin.Context) {
		if m.validator.CheckExpiringSoon(days) {
			c.Header("X-License-Warning", "expiring soon")
			c.Header("X-License-Remaining-Days", string(rune(m.validator.GetRemainingDays())))
		}
		c.Next()
	}
}

// shouldSkip 检查是否跳过路径
func (m *Middleware) shouldSkip(path string) bool {
	for _, skipPath := range m.skipPaths {
		if path == skipPath || strings.HasPrefix(path, skipPath) {
			return true
		}
	}
	return false
}

// GetLicenseFromContext 从 Gin 上下文获取 License 信息
func GetLicenseFromContext(c *gin.Context) *License {
	if license, exists := c.Get("license"); exists {
		return license.(*License)
	}
	return nil
}

// GetLicenseIDFromContext 从 Gin 上下文获取 License ID
func GetLicenseIDFromContext(c *gin.Context) string {
	if licenseID, exists := c.Get("licenseID"); exists {
		return licenseID.(string)
	}
	return ""
}

// GetLicenseTypeFromContext 从 Gin 上下文获取 License 类型
func GetLicenseTypeFromContext(c *gin.Context) string {
	if licenseType, exists := c.Get("licenseType"); exists {
		return licenseType.(string)
	}
	return ""
}

// RequireLicense 辅助函数：快速创建 License 验证中间件
func RequireLicense(validator *Validator, config *Config) gin.HandlerFunc {
	return NewMiddleware(validator, config).Handler()
}

// RequireFeature 辅助函数：快速创建功能权限验证中间件
func RequireFeature(validator *Validator, feature string) gin.HandlerFunc {
	m := &Middleware{validator: validator}
	return m.FeatureRequired(feature)
}

// LicenseInfoResponse License 信息响应
type LicenseInfoResponse struct {
	LicenseID string   `json:"licenseId"`
	IssuedTo  string   `json:"issuedTo"`
	Type      string   `json:"type"`
	Features  []string `json:"features"`
	ExpiresAt string   `json:"expiresAt"`
	Remaining int      `json:"remaining"`
	Status    string   `json:"status"`
}

// GetLicenseInfoHandler 获取 License 信息的 API 处理器
func GetLicenseInfoHandler(validator *Validator) gin.HandlerFunc {
	return func(c *gin.Context) {
		license := validator.GetInfo()
		if license == nil {
			c.JSON(http.StatusNotFound, gin.H{
				"code":    404,
				"message": "license not found",
			})
			return
		}

		response := LicenseInfoResponse{
			LicenseID: license.LicenseID,
			IssuedTo:  license.IssuedTo,
			Type:      license.Type,
			Features:  license.Features,
			ExpiresAt: license.ExpiresAt.Format("2006-01-02 15:04:05"),
			Remaining: validator.GetRemainingDays(),
			Status:    license.Status,
		}

		c.JSON(http.StatusOK, gin.H{
			"code": 200,
			"data": response,
		})
	}
}
