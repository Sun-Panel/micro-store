package biz

import (
	"sun-panel/global"
	"sun-panel/lib/cache"
	"sun-panel/models"
	"time"

	"gorm.io/gorm"
)

// DeveloperService 开发者业务服务
type DeveloperService struct {
	// 开发者信息缓存，key: developerName，value: models.Developer
	InfoCache cache.Cacher[models.Developer]
}

// Init 初始化开发者服务（含缓存）
func (s *DeveloperService) Init() {
	s.InfoCache = global.NewCache[models.Developer](1*time.Hour, 10*time.Minute, "DeveloperInfoCache")
}

// invalidateCache 清除指定 developerName 的缓存
func (s *DeveloperService) invalidateCache(developerName string) {
	if s.InfoCache != nil {
		s.InfoCache.Delete(developerName)
	}
}

// RegisterParams 注册开发者参数
type RegisterParams struct {
	UserId        uint
	DeveloperName string
	ContactMail   string
	PaymentName   string
	PaymentQrcode string
	PaymentMethod string
	Name          string
}

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

// Register 注册成为开发者（业务层，包含用户权限更新）
func (s *DeveloperService) Register(db *gorm.DB, p RegisterParams) (uint, error) {
	m := models.Developer{}
	id, err := m.Register(db, p.UserId, p.DeveloperName, p.ContactMail, p.PaymentName, p.PaymentQrcode, p.PaymentMethod, p.Name)
	if err != nil {
		return 0, err
	}

	user := models.User{}
	currentUser, err := user.GetUserInfoByUid(p.UserId)
	if err != nil {
		return id, nil
	}

	newRole := models.AddRole(currentUser.Role, models.ROLE_DEVELOPER)
	if err := user.UpdateUserInfoByUserId(p.UserId, map[string]interface{}{"role": newRole}); err != nil {
		return id, err
	}

	return id, nil
}

// BatchGetByDeveloperNames 批量获取开发者信息（带缓存）
// 1. 先从缓存获取已有的
// 2. 未命中的再批量查库
// 3. 查库结果写入缓存
func (s *DeveloperService) BatchGetByDeveloperNames(db *gorm.DB, developerNames []string) (map[string]models.Developer, error) {
	result := make(map[string]models.Developer, len(developerNames))

	if s.InfoCache == nil {
		// 缓存未初始化，降级为直接查库
		var developers []models.Developer
		if err := db.Where("developer_name IN ?", developerNames).Find(&developers).Error; err != nil {
			return nil, err
		}
		for _, dev := range developers {
			result[dev.DeveloperName] = dev
		}
		return result, nil
	}

	// 1. 读缓存
	var missNames []string
	for _, name := range developerNames {
		if dev, ok := s.InfoCache.Get(name); ok {
			result[name] = dev
		} else {
			missNames = append(missNames, name)
		}
	}

	// 2. 缓存未命中，批量查库
	if len(missNames) > 0 {
		var developers []models.Developer
		if err := db.Where("developer_name IN ?", missNames).Find(&developers).Error; err != nil {
			return nil, err
		}
		for _, dev := range developers {
			result[dev.DeveloperName] = dev
			s.InfoCache.SetDefault(dev.DeveloperName, dev) // 写入缓存
		}
	}

	return result, nil
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

		// 更新后清除缓存
		defer s.invalidateCache(developer.DeveloperName)
	}

	// 调用 Model 层执行数据库操作
	m := models.Developer{}
	return m.UpdateInfo(db, id, updateFields)
}
