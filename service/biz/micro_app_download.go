package biz

import (
	"sun-panel/global"
	"sun-panel/lib/cache"
	"time"
)

type microAppDownload struct {
	// token:versionId
	AdminDownloadCacheVersionId cache.Cacher[string]
}

func (m *microAppDownload) Init() {
	m.AdminDownloadCacheVersionId = global.NewCache[string](20*time.Minute, 30*time.Minute, "micro_app_admin_download_version_id")
}

func (m *microAppDownload) Download(appId string) (string, error) {

	return "", nil
}
