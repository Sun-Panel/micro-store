package models

import (
	"time"

	"gorm.io/gorm"
)

type SoftwareClient struct {
	BaseModel
	ClientId       string    `gorm:"type:varchar(50);not null;index" json:"clientId"`             // 客户端id
	Version        string    `gorm:"type:varchar(100);not null;index:idx_version" json:"version"` // 软件版本号
	MacAddress     string    `gorm:"type:varchar(20)" json:"macAddress"`                          // mac地址
	LastIp         string    `gorm:"type:varchar(50);not null" json:"lastIp"`                     // 最后登录ip
	LastLanIp      string    `gorm:"type:varchar(50);not null" json:"lastLanIp"`                  // 最后登录的内网ip
	JoinTime       time.Time `gorm:"index:idx_join_time" json:"joinTime"`                         // 加入时间
	LastOnlineTime time.Time `gorm:"index:idx_last_online_time,sort:desc" json:"lastOnlineTime"`  // 最后一次在线时间

	UserId uint `gorm:"not null;index" json:"userId"` // 绑定的用户
	User   User `json:"user"`
}

// 查询是否存在某个clientid，如果不存在则自动生成一个新的
// func GetOrCreateSoftwareClientByClientId(db *gorm.DB, clientId string) (*SoftwareClient, error) {
// 	var softwareClient SoftwareClient
// 	if err := db.Where("clicent_id = ?", clientId).First(&softwareClient).Error; err != nil {
// 		if gorm.IsRecordNotFoundError(err) {
// 			// 如果记录不存在，则创建新的软件客户端
// 			var newClientId string
// 			for {
// 				// 生成唯一的ID
// 				newClientId = uuid.New().String()
// 				// 检查生成的ID是否已存在
// 				var count int
// 				if err := db.Model(&SoftwareClient{}).Where("id = ?", newClientId).Count(&count).Error; err != nil {
// 					return nil, err
// 				}
// 				if count == 0 {
// 					break
// 				}
// 			}
// 			softwareClient = SoftwareClient{
// 				ClientId: clientId,
// 				// 这里可以设置其他默认值
// 			}
// 			if err := db.Create(&softwareClient).Error; err != nil {
// 				return nil, err
// 			}
// 		} else {
// 			return nil, err
// 		}
// 	}
// 	return &softwareClient, nil
// }

// 更新软件客户端信息
func UpdateSoftwareClient(db *gorm.DB, softwareClient *SoftwareClient) error {
	return db.Save(softwareClient).Error
}

// 查询clientId是否存在
func IsClientIdExist(db *gorm.DB, clientId string) (bool, error) {
	var exists bool
	if err := db.Model(&SoftwareClient{}).
		Select("1").
		Where("client_id = ?", clientId).
		Limit(1).
		Scan(&exists).Error; err != nil {
		return false, err
	}
	return exists, nil
}
