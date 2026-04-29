package biz

import (
	"sun-panel/models"
	"time"

	"gorm.io/gorm"
)

// DeveloperService 开发者业务服务
type DeveloperService struct{}

// GetByDeveloperName 根据开发者标识获取开发者信息
func (s *DeveloperService) GetByDeveloperName(db *gorm.DB, developerName string) (models.Developer, error) {
	m := models.Developer{}
	developer, err := m.GetByDeveloperName(db, developerName)
	if err != nil {
		return models.Developer{}, err
	}
	return developer, nil
}

// GetDeveloperInfo 根据ID获取开发者详情（业务层）
func (s *DeveloperService) GetDeveloperInfo(db *gorm.DB, id uint) (models.Developer, error) {
	m := models.Developer{}
	developer, err := m.GetById(db, id)
	if err != nil {
		return models.Developer{}, err
	}
	return developer, nil
}

// GetByUserId 根据用户ID获取开发者信息（业务层）
func (s *DeveloperService) GetByUserId(db *gorm.DB, userId uint) (models.Developer, error) {
	m := models.Developer{}
	developer, err := m.GetByUserId(db, userId)
	if err != nil {
		return models.Developer{}, err
	}
	return developer, nil
}

// GetDeveloperList 获取开发者列表（业务层）
func (s *DeveloperService) GetDeveloperList(db *gorm.DB, page, limit int, status *int, keyWord string) ([]models.Developer, int64, error) {
	m := models.Developer{}
	return m.GetList(db, page, limit, status, keyWord)
}

// UpdateDeveloperInfo 更新开发者信息（业务层，包含业务规则校验）
func (s *DeveloperService) UpdateDeveloperInfo(db *gorm.DB, id uint, updateFields models.DeveloperUpdateFields) error {
	// 如果要修改 Name，检查冷却期（180天）
	if updateFields.Name != nil {
		developer, err := s.GetDeveloperInfo(db, id)
		if err != nil {
			return err
		}

		// 如果有上次更新时间，检查是否满180天
		if developer.NameUpdatedAt != nil {
			daysSinceUpdate := time.Since(*developer.NameUpdatedAt).Hours() / 24
			if daysSinceUpdate < 180 {
				daysRemaining := 180 - int(daysSinceUpdate)
				return models.NewModelErrorWithData("E_DEVELOPER_NAME_COOLDOWN", map[string]any{
					"daysRemaining": daysRemaining,
				})
			}
		}
	}

	// 调用 Model 层执行数据库操作
	m := models.Developer{}
	return m.UpdateInfo(db, id, updateFields)
}
