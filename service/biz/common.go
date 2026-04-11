package biz

import (
	"os"
	"path/filepath"
	"strings"
	"sun-panel/lib/cmn"
)

func ExtractMicroAppPackageZipToTemp(zipPath string) (string, error) {
	// 计算文件的哈希值，用于创建唯一的临时目录名
	// 去掉文件后缀作为目录名
	uniqueName := strings.TrimSuffix(filepath.Base(zipPath), filepath.Ext(zipPath))

	tempPath := Config.GetTempPath()

	// 创建临时解压目录
	tempDir := filepath.Join(tempPath, "micro_app_extract", uniqueName+"_"+cmn.BuildRandCode(10, cmn.RAND_CODE_MODE1))
	if err := os.MkdirAll(tempDir, 0755); err != nil {
		return "", err
	}

	// 解压文件
	err := cmn.UnzipFile(zipPath, tempDir)
	if err != nil {
		os.RemoveAll(tempDir) // 解压失败，清理临时目录
		return "", err
	}

	return tempDir, nil
}
