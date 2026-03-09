package license

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"
)

// TimeValidator 时间验证器
type TimeValidator struct {
	config        *Config
	ntpServers    []string
	lastNTPCheck  time.Time
	lastKnownTime time.Time
	totalRunTime  time.Duration
	lastStartTime time.Time
	ntpFailCount  int
	mu            sync.RWMutex
}

// NewTimeValidator 创建时间验证器
func NewTimeValidator(config *Config) *TimeValidator {
	return &TimeValidator{
		config: config,
		ntpServers: []string{
			"time.google.com",
			"time.cloudflare.com",
			"time.windows.com",
			"pool.ntp.org",
		},
		lastStartTime: time.Now(),
	}
}

// ValidateTime 验证时间是否被篡改
func (tv *TimeValidator) ValidateTime() *TimeValidation {
	tv.mu.Lock()
	defer tv.mu.Unlock()

	localTime := time.Now()

	// 1. 尝试从 NTP 服务器获取真实时间
	realTime, err := tv.getNTPTime()
	if err != nil {
		tv.ntpFailCount++
		// NTP 失败时，使用累计运行时间估算
		return tv.validateWithEstimatedTime(localTime)
	}

	tv.ntpFailCount = 0
	tv.lastNTPCheck = time.Now()
	tv.lastKnownTime = realTime

	// 2. 计算时间偏差
	drift := localTime.Unix() - realTime.Unix()
	driftSeconds := drift
	if driftSeconds < 0 {
		driftSeconds = -driftSeconds
	}

	// 3. 检查偏差是否超过限制
	tampered := driftSeconds > tv.config.TimeDriftLimit

	message := "time is valid"
	if tampered {
		message = fmt.Sprintf("time drift detected: %d seconds (limit: %d)", driftSeconds, tv.config.TimeDriftLimit)
	}

	return &TimeValidation{
		Valid:     !tampered,
		Message:   message,
		RealTime:  realTime.Unix(),
		LocalTime: localTime.Unix(),
		Drift:     drift,
		Tampered:  tampered,
	}
}

// validateWithEstimatedTime 使用估算时间验证（NTP 失败时的降级方案）
func (tv *TimeValidator) validateWithEstimatedTime(localTime time.Time) *TimeValidation {
	// 如果之前有 NTP 同步，使用累计运行时间估算
	if !tv.lastKnownTime.IsZero() {
		elapsed := time.Since(tv.lastStartTime)
		estimatedTime := tv.lastKnownTime.Add(elapsed)

		drift := localTime.Unix() - estimatedTime.Unix()
		driftSeconds := drift
		if driftSeconds < 0 {
			driftSeconds = -driftSeconds
		}

		// 放宽限制，允许更大的偏差
		tampered := driftSeconds > tv.config.TimeDriftLimit*2

		message := "time validation using estimated time (NTP unavailable)"
		if tampered {
			message = fmt.Sprintf("time tampering detected (estimated): drift %d seconds", driftSeconds)
		}

		return &TimeValidation{
			Valid:     !tampered,
			Message:   message,
			RealTime:  estimatedTime.Unix(),
			LocalTime: localTime.Unix(),
			Drift:     drift,
			Tampered:  tampered,
		}
	}

	// 无法验证，返回通过（降级）
	return &TimeValidation{
		Valid:     true,
		Message:   "time validation skipped (NTP unavailable and no reference)",
		RealTime:  localTime.Unix(),
		LocalTime: localTime.Unix(),
		Drift:     0,
		Tampered:  false,
	}
}

// getNTPTime 从 NTP 服务器获取时间
func (tv *TimeValidator) getNTPTime() (time.Time, error) {
	// 简化版本：使用 HTTP 方式获取时间
	// 实际生产环境应该使用 NTP 协议

	for _, server := range tv.ntpServers {
		url := fmt.Sprintf("https://%s", server)

		client := &http.Client{Timeout: 3 * time.Second}
		resp, err := client.Head(url)
		if err != nil {
			continue
		}
		resp.Body.Close()

		// 从响应头获取时间
		dateStr := resp.Header.Get("Date")
		if dateStr != "" {
			t, err := time.Parse(time.RFC1123, dateStr)
			if err == nil {
				return t, nil
			}
		}
	}

	return time.Time{}, fmt.Errorf("all NTP servers failed")
}

// getNTPTimeReal 真实的 NTP 时间获取（需要实现 NTP 协议）
// 这只是示例，实际使用需要引入 NTP 库或自己实现 NTP 协议
func (tv *TimeValidator) getNTPTimeReal() (time.Time, error) {
	// TODO: 实现 NTP 协议获取时间
	// 可以使用 github.com/beevik/ntp 库
	return time.Time{}, fmt.Errorf("not implemented")
}

// CheckTimeRollback 检查时间是否回拨
func (tv *TimeValidator) CheckTimeRollback() bool {
	tv.mu.Lock()
	defer tv.mu.Unlock()

	currentTime := time.Now()

	// 如果当前时间早于上次记录的时间，说明时间被回拨
	if !tv.lastKnownTime.IsZero() && currentTime.Before(tv.lastKnownTime) {
		return true
	}

	tv.lastKnownTime = currentTime
	return false
}

// UpdateLastCheckTime 更新最后检查时间
func (tv *TimeValidator) UpdateLastCheckTime() {
	tv.mu.Lock()
	defer tv.mu.Unlock()
	tv.lastNTPCheck = time.Now()
}

// GetLastNTPCheck 获取上次 NTP 检查时间
func (tv *TimeValidator) GetLastNTPCheck() time.Time {
	tv.mu.RLock()
	defer tv.mu.RUnlock()
	return tv.lastNTPCheck
}

// GetNTPFailCount 获取 NTP 失败次数
func (tv *TimeValidator) GetNTPFailCount() int {
	tv.mu.RLock()
	defer tv.mu.RUnlock()
	return tv.ntpFailCount
}

// ValidateTimeWithTolerance 带容差的验证（允许一定误差）
func (tv *TimeValidator) ValidateTimeWithTolerance(toleranceSeconds int64) *TimeValidation {
	result := tv.ValidateTime()
	if !result.Valid && result.Drift <= toleranceSeconds {
		result.Valid = true
		result.Message = fmt.Sprintf("time drift within tolerance: %d seconds", result.Drift)
	}
	return result
}

// CheckTimeWithLicense 检查时间并与 License 过期时间对比
func (tv *TimeValidator) CheckTimeWithLicense(license *License) (*TimeValidation, bool) {
	// 1. 验证时间是否被篡改
	timeResult := tv.ValidateTime()

	// 2. 使用真实时间检查过期
	var realTime time.Time
	if timeResult.RealTime > 0 {
		realTime = time.Unix(timeResult.RealTime, 0)
	} else {
		realTime = time.Now()
	}

	// 3. 检查是否过期
	isExpired := realTime.After(license.ExpiresAt)

	return timeResult, isExpired
}

// GetTrueTime 获取真实时间（优先使用 NTP 时间）
func (tv *TimeValidator) GetTrueTime() time.Time {
	tv.mu.RLock()
	defer tv.mu.RUnlock()

	if !tv.lastKnownTime.IsZero() {
		// 使用最后已知的真实时间 + 经过的时长
		elapsed := time.Since(tv.lastNTPCheck)
		return tv.lastKnownTime.Add(elapsed)
	}

	return time.Now()
}

// HTTPTimeClient 使用 HTTP 获取时间的客户端
type HTTPTimeClient struct {
	timeAPIs []string
}

// NewHTTPTimeClient 创建 HTTP 时间客户端
func NewHTTPTimeClient() *HTTPTimeClient {
	return &HTTPTimeClient{
		timeAPIs: []string{
			"https://worldtimeapi.org/api/timezone/Asia/Shanghai",
			"https://timeapi.io/api/Time/current/zone?timeZone=Asia/Shanghai",
		},
	}
}

// GetTime 从 HTTP API 获取时间
func (c *HTTPTimeClient) GetTime() (time.Time, error) {
	for _, api := range c.timeAPIs {
		resp, err := http.Get(api)
		if err != nil {
			continue
		}

		body, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			continue
		}

		// 解析时间响应
		var result map[string]interface{}
		if err := json.Unmarshal(body, &result); err != nil {
			continue
		}

		// 提取时间字段
		if datetime, ok := result["datetime"].(string); ok {
			t, err := time.Parse(time.RFC3339, datetime)
			if err == nil {
				return t, nil
			}
		}
	}

	return time.Time{}, fmt.Errorf("all time APIs failed")
}
