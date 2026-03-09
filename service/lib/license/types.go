package license

import "time"

// License 授权许可证结构
type License struct {
	LicenseID string    `json:"licenseId"`           // 唯一标识
	Product   string    `json:"product"`             // 产品名称
	Version   string    `json:"version"`             // 授权版本
	IssuedTo  string    `json:"issuedTo"`            // 授权用户/组织
	IssuedAt  time.Time `json:"issuedAt"`            // 颁发时间
	ExpiresAt time.Time `json:"expiresAt"`           // 过期时间
	Features  []string  `json:"features"`            // 授权的功能列表
	MachineID string    `json:"machineId,omitempty"` // 绑定的机器码（可选）
	MaxUsers  int       `json:"maxUsers,omitempty"`  // 最大用户数（可选）
	MaxNodes  int       `json:"maxNodes,omitempty"`  // 最大节点数（可选）
	Extra     string    `json:"extra,omitempty"`     // 扩展信息
	Signature string    `json:"signature"`           // 数字签名
	Type      string    `json:"type"`                // 授权类型：trial, standard, professional, enterprise
	Status    string    `json:"status"`              // 状态：active, expired, revoked
}

// LicenseOptions 生成 License 的选项
type LicenseOptions struct {
	LicenseID string   `json:"licenseId"`
	Product   string   `json:"product"`
	Version   string   `json:"version"`
	IssuedTo  string   `json:"issuedTo"`
	Duration  int      `json:"duration"` // 有效期（天）
	Features  []string `json:"features"`
	MachineID string   `json:"machineId,omitempty"`
	MaxUsers  int      `json:"maxUsers,omitempty"`
	MaxNodes  int      `json:"maxNodes,omitempty"`
	Extra     string   `json:"extra,omitempty"`
	Type      string   `json:"type"` // trial, standard, professional, enterprise
}

// ValidateResult 验证结果
type ValidateResult struct {
	Valid       bool      `json:"valid"`
	Message     string    `json:"message"`
	LicenseInfo *License  `json:"licenseInfo,omitempty"`
	ExpiresAt   time.Time `json:"expiresAt,omitempty"`
	Remaining   int       `json:"remaining"` // 剩余天数
}

// HeartbeatResult 心跳检测结果
type HeartbeatResult struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// InitResult 初始化结果
type InitResult struct {
	Success bool              `json:"success"`
	Message string            `json:"message"`
	Info    map[string]string `json:"info,omitempty"`
}

// LicenseStatus License 状态
type LicenseStatus struct {
	LicenseID   string    `json:"licenseId"`
	Status      string    `json:"status"`
	ActivatedAt time.Time `json:"activatedAt,omitempty"`
	LastSeenAt  time.Time `json:"lastSeenAt,omitempty"`
	InstanceID  string    `json:"instanceId,omitempty"`
}

// MachineInfo 机器信息
type MachineInfo struct {
	Hostname    string   `json:"hostname"`
	ContainerID string   `json:"containerId,omitempty"`
	MACAddress  []string `json:"macAddress"`
	CPUID       string   `json:"cpuId"`
	DiskID      string   `json:"diskId"`
	Platform    string   `json:"platform"`
	Arch        string   `json:"arch"`
}

// TimeValidation 时间验证结果
type TimeValidation struct {
	Valid     bool   `json:"valid"`
	Message   string `json:"message"`
	RealTime  int64  `json:"realTime"`  // 真实时间戳
	LocalTime int64  `json:"localTime"` // 本地时间戳
	Drift     int64  `json:"drift"`     // 时间偏差（秒）
	Tampered  bool   `json:"tampered"`  // 是否被篡改
}

// Config License 验证器配置
type Config struct {
	PublicKey       string        `json:"publicKey"`       // RSA 公钥
	PrivateKey      string        `json:"privateKey"`      // RSA 私钥（仅生成器需要）
	LicenseFilePath string        `json:"licenseFilePath"` // License 文件路径
	HeartbeatURL    string        `json:"heartbeatUrl"`    // 心跳验证 URL（可选）
	HeartbeatPeriod time.Duration `json:"heartbeatPeriod"` // 心跳周期
	OfflineLimit    int           `json:"offlineLimit"`    // 离线有效期限（天）
	TimeDriftLimit  int64         `json:"timeDriftLimit"`  // 时间偏差限制（秒）
	MachineBind     bool          `json:"machineBind"`     // 是否绑定机器码
	Salt            string        `json:"salt"`            // 自定义盐值
}

// DefaultConfig 默认配置
func DefaultConfig() *Config {
	return &Config{
		HeartbeatPeriod: 5 * time.Minute,
		OfflineLimit:    7,    // 7 天离线期限
		TimeDriftLimit:  300,  // 5 分钟时间偏差
		MachineBind:     true, // 默认绑定机器码
		Salt:            "sun-panel-license-salt-2024",
	}
}
