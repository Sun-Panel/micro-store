package license

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"
)

// HeartbeatManager 心跳管理器
type HeartbeatManager struct {
	validator       *Validator
	config          *Config
	instanceID      string
	lastHeartbeat   time.Time
	stopChan        chan struct{}
	mu              sync.RWMutex
	running         bool
	heartbeatErrors int // 连续心跳错误次数
}

// NewHeartbeatManager 创建心跳管理器
func NewHeartbeatManager(validator *Validator, config *Config) *HeartbeatManager {
	return &HeartbeatManager{
		validator: validator,
		config:    config,
		stopChan:  make(chan struct{}),
	}
}

// Start 启动心跳检测
func (h *HeartbeatManager) Start() error {
	h.mu.Lock()
	defer h.mu.Unlock()

	if h.running {
		return nil
	}

	if h.config.HeartbeatURL == "" {
		return fmt.Errorf("heartbeat URL not configured")
	}

	// 生成实例 ID
	if h.instanceID == "" {
		machineID, _ := GetMachineID(h.config.Salt)
		h.instanceID = fmt.Sprintf("%s-%d", machineID[:8], time.Now().Unix())
	}

	h.running = true
	go h.heartbeatLoop()

	return nil
}

// Stop 停止心跳检测
func (h *HeartbeatManager) Stop() {
	h.mu.Lock()
	defer h.mu.Unlock()

	if !h.running {
		return
	}

	close(h.stopChan)
	h.running = false
}

// heartbeatLoop 心跳循环
func (h *HeartbeatManager) heartbeatLoop() {
	ticker := time.NewTicker(h.config.HeartbeatPeriod)
	defer ticker.Stop()

	// 首次立即执行一次
	h.doHeartbeat()

	for {
		select {
		case <-ticker.C:
			h.doHeartbeat()
		case <-h.stopChan:
			return
		}
	}
}

// doHeartbeat 执行心跳
func (h *HeartbeatManager) doHeartbeat() *HeartbeatResult {
	license := h.validator.GetLicense()
	if license == nil {
		return &HeartbeatResult{
			Success: false,
			Message: "no license loaded",
		}
	}

	// 构建心跳请求
	req := map[string]interface{}{
		"licenseId":   license.LicenseID,
		"instanceId":  h.instanceID,
		"timestamp":   time.Now().Unix(),
		"product":     license.Product,
		"version":     license.Version,
		"type":        license.Type,
		"machineInfo": nil, // 可选：发送机器信息
	}

	// 获取机器信息
	machineInfo, _ := GetMachineInfo()
	req["machineInfo"] = machineInfo

	// 发送心跳请求
	result, err := h.sendHeartbeat(req)
	if err != nil {
		h.mu.Lock()
		h.heartbeatErrors++
		h.mu.Unlock()

		return &HeartbeatResult{
			Success: false,
			Message: fmt.Sprintf("heartbeat failed: %v", err),
		}
	}

	h.mu.Lock()
	h.lastHeartbeat = time.Now()
	h.heartbeatErrors = 0
	h.mu.Unlock()

	return result
}

// sendHeartbeat 发送心跳请求
func (h *HeartbeatManager) sendHeartbeat(req map[string]interface{}) (*HeartbeatResult, error) {
	reqBody, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("marshal request failed: %w", err)
	}

	// 创建 HTTP 请求
	httpReq, err := http.NewRequest("POST", h.config.HeartbeatURL, bytes.NewReader(reqBody))
	if err != nil {
		return nil, fmt.Errorf("create request failed: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("X-License-ID", fmt.Sprintf("%v", req["licenseId"]))
	httpReq.Header.Set("X-Instance-ID", fmt.Sprintf("%v", req["instanceId"]))

	// 发送请求
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("send request failed: %w", err)
	}
	defer resp.Body.Close()

	// 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response failed: %w", err)
	}

	// 解析响应
	var result HeartbeatResult
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("parse response failed: %w", err)
	}

	return &result, nil
}

// CheckHeartbeat 检查心跳状态
func (h *HeartbeatManager) CheckHeartbeat() *HeartbeatResult {
	h.mu.RLock()
	defer h.mu.RUnlock()

	if h.config.HeartbeatURL == "" {
		return &HeartbeatResult{
			Success: true,
			Message: "heartbeat not required",
		}
	}

	if !h.running {
		return &HeartbeatResult{
			Success: false,
			Message: "heartbeat not started",
		}
	}

	// 检查是否超过离线期限
	if h.config.OfflineLimit > 0 {
		offlineHours := time.Since(h.lastHeartbeat).Hours() / 24
		if offlineHours > float64(h.config.OfflineLimit) {
			return &HeartbeatResult{
				Success: false,
				Message: fmt.Sprintf("offline for %.0f days, exceed limit %d days", offlineHours, h.config.OfflineLimit),
			}
		}
	}

	// 检查连续错误次数
	if h.heartbeatErrors >= 3 {
		return &HeartbeatResult{
			Success: false,
			Message: fmt.Sprintf("heartbeat failed %d times consecutively", h.heartbeatErrors),
		}
	}

	return &HeartbeatResult{
		Success: true,
		Message: "heartbeat ok",
	}
}

// GetLastHeartbeat 获取上次心跳时间
func (h *HeartbeatManager) GetLastHeartbeat() time.Time {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return h.lastHeartbeat
}

// GetInstanceID 获取实例 ID
func (h *HeartbeatManager) GetInstanceID() string {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return h.instanceID
}

// GetHeartbeatErrors 获取连续心跳错误次数
func (h *HeartbeatManager) GetHeartbeatErrors() int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return h.heartbeatErrors
}

// ForceHeartbeat 强制执行一次心跳
func (h *HeartbeatManager) ForceHeartbeat() *HeartbeatResult {
	return h.doHeartbeat()
}

// HeartbeatStatus 心跳状态
type HeartbeatStatus struct {
	Running         bool      `json:"running"`
	InstanceID      string    `json:"instanceId"`
	LastHeartbeat   time.Time `json:"lastHeartbeat"`
	HeartbeatErrors int       `json:"heartbeatErrors"`
}

// GetStatus 获取心跳状态
func (h *HeartbeatManager) GetStatus() *HeartbeatStatus {
	h.mu.RLock()
	defer h.mu.RUnlock()

	return &HeartbeatStatus{
		Running:         h.running,
		InstanceID:      h.instanceID,
		LastHeartbeat:   h.lastHeartbeat,
		HeartbeatErrors: h.heartbeatErrors,
	}
}
