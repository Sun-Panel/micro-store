package license

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/net"
)

// GetMachineInfo 获取机器信息
func GetMachineInfo() (*MachineInfo, error) {
	info := &MachineInfo{
		Platform: runtime.GOOS,
		Arch:     runtime.GOARCH,
	}

	// 获取主机名
	hostname, err := os.Hostname()
	if err == nil {
		info.Hostname = hostname
	}

	// 获取主机信息
	hostInfo, err := host.Info()
	if err == nil {
		info.Hostname = hostInfo.Hostname
		info.Platform = hostInfo.OS
	}

	// 获取 Docker 容器 ID（从 /proc/self/cgroup 读取）
	containerID := getContainerID()
	if containerID != "" {
		info.ContainerID = containerID
	}

	// 获取 MAC 地址
	interfaces, err := net.Interfaces()
	if err == nil {
		for _, iface := range interfaces {
			// 跳过回环接口和虚拟接口
			if len(iface.HardwareAddr) > 0 && !strings.Contains(iface.Name, "lo") && !strings.Contains(iface.Name, "docker") {
				info.MACAddress = append(info.MACAddress, iface.HardwareAddr)
			}
		}
	}

	// 获取 CPU 信息
	cpuInfo, err := cpu.Info()
	if err == nil && len(cpuInfo) > 0 {
		info.CPUID = fmt.Sprintf("%s-%s", cpuInfo[0].VendorID, cpuInfo[0].ModelName)
	}

	// 获取磁盘 ID（使用根目录所在磁盘）
	diskID := getDiskID()
	if diskID != "" {
		info.DiskID = diskID
	}

	return info, nil
}

// GetMachineID 生成机器码
func GetMachineID(salt string) (string, error) {
	info, err := GetMachineInfo()
	if err != nil {
		return "", err
	}

	// 组合多个特征生成唯一标识
	data := fmt.Sprintf("%s|%s|%s|%v|%s|%s|%s",
		info.Hostname,
		info.ContainerID,
		info.MACAddress,
		info.CPUID,
		info.DiskID,
		info.Platform,
		salt,
	)

	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:]), nil
}

// GetMachineIDWithFeatures 使用指定特征生成机器码
func GetMachineIDWithFeatures(features []string, salt string) (string, error) {
	info, err := GetMachineInfo()
	if err != nil {
		return "", err
	}

	var dataParts []string
	for _, feature := range features {
		switch feature {
		case "hostname":
			dataParts = append(dataParts, info.Hostname)
		case "container":
			dataParts = append(dataParts, info.ContainerID)
		case "mac":
			dataParts = append(dataParts, fmt.Sprintf("%v", info.MACAddress))
		case "cpu":
			dataParts = append(dataParts, info.CPUID)
		case "disk":
			dataParts = append(dataParts, info.DiskID)
		case "platform":
			dataParts = append(dataParts, info.Platform)
		}
	}
	dataParts = append(dataParts, salt)

	data := strings.Join(dataParts, "|")
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:]), nil
}

// getContainerID 获取 Docker 容器 ID
func getContainerID() string {
	// 尝试从 /proc/self/cgroup 读取容器 ID
	data, err := os.ReadFile("/proc/self/cgroup")
	if err != nil {
		return ""
	}

	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		// Docker 容器的 cgroup 行通常包含容器 ID
		if strings.Contains(line, "docker") || strings.Contains(line, "kubepods") {
			parts := strings.Split(line, "/")
			for _, part := range parts {
				// 容器 ID 通常是 64 位十六进制字符串
				if len(part) >= 12 {
					// 取前 12 位作为容器 ID
					return part[:12]
				}
			}
		}
	}

	// 尝试从环境变量获取（Docker 可能通过环境变量传递）
	if containerID := os.Getenv("HOSTNAME"); len(containerID) == 12 {
		return containerID
	}

	// 尝试从 /proc/1/cpuset 读取
	data, err = os.ReadFile("/proc/1/cpuset")
	if err == nil {
		content := strings.TrimSpace(string(data))
		parts := strings.Split(content, "/")
		if len(parts) > 0 {
			lastPart := parts[len(parts)-1]
			if len(lastPart) >= 12 {
				return lastPart[:12]
			}
		}
	}

	return ""
}

// getDiskID 获取磁盘 ID（简化版，使用磁盘序列号或标识）
func getDiskID() string {
	// 在 Linux 下尝试读取磁盘序列号
	// 这里简化处理，使用主机 ID 或生成一个标识

	// 尝试读取 /etc/machine-id（Linux 系统唯一标识）
	if data, err := os.ReadFile("/etc/machine-id"); err == nil {
		return strings.TrimSpace(string(data))
	}

	// 尝试读取 /var/lib/dbus/machine-id
	if data, err := os.ReadFile("/var/lib/dbus/machine-id"); err == nil {
		return strings.TrimSpace(string(data))
	}

	// Windows 下可以读取磁盘序列号（需要 WMI，这里暂不实现）
	return ""
}

// VerifyMachineID 验证机器码
func VerifyMachineID(expectedMachineID string, salt string) (bool, error) {
	currentMachineID, err := GetMachineID(salt)
	if err != nil {
		return false, err
	}
	return currentMachineID == expectedMachineID, nil
}

// VerifyMachineIDWithFeatures 使用指定特征验证机器码
func VerifyMachineIDWithFeatures(expectedMachineID string, features []string, salt string) (bool, error) {
	currentMachineID, err := GetMachineIDWithFeatures(features, salt)
	if err != nil {
		return false, err
	}
	return currentMachineID == expectedMachineID, nil
}
