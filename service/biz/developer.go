package biz

import (
	"sun-panel/models"

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
