package file

import (
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"strings"

	"golang.org/x/image/webp"
)

// GetImageDimensions 获取图片尺寸
// 参数:
//   - filepath: 图片文件路径
//   - ext: 图片扩展名（如 .png, .jpg, .webp）
//
// 返回:
//   - width: 图片宽度（单位：像素/px）
//   - height: 图片高度（单位：像素/px）
//   - err: 错误信息，非图片文件或解码失败时返回错误
func GetImageDimensions(filepath, ext string) (width, height int, err error) {
	// width: 图片宽度（单位：像素/px）
	// height: 图片高度（单位：像素/px）
	file, err := os.Open(filepath)
	if err != nil {
		return 0, 0, err
	}
	defer file.Close()

	var img image.Image
	switch strings.ToLower(ext) {
	case ".webp":
		img, err = webp.Decode(file)
	default:
		img, _, err = image.Decode(file)
	}

	if err != nil {
		return 0, 0, err
	}

	bounds := img.Bounds()
	return bounds.Dx(), bounds.Dy(), nil
}

// IsImageExt 检查是否为支持的图片扩展名
func IsImageExt(ext string) bool {
	allowedExts := []string{".jpg", ".jpeg", ".png", ".gif", ".webp", ".bmp", ".svg"}
	for _, allowed := range allowedExts {
		if strings.EqualFold(ext, allowed) {
			return true
		}
	}
	return false
}
