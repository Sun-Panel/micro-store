package models

// 角色常量定义 - 使用位运算支持角色叠加
const (
	ROLE_USER      int = 1  // 普通用户   (二进制: 00001)
	ROLE_DEVELOPER int = 2  // 开发者     (二进制: 00010)
	ROLE_ADMIN     int = 4  // 管理员     (二进制: 00100)
	ROLE_AUDITOR   int = 8  // 审核员     (二进制: 01000)
	ROLE_OPERATOR  int = 16 // 运营       (二进制: 10000)
)

// 角色信息结构体
type RoleInfo struct {
	Value       int    `json:"value"`
	Key         string `json:"key"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// 所有角色列表
var AllRoles = []RoleInfo{
	{Value: ROLE_USER, Key: "user", Name: "普通用户", Description: "浏览、下载、安装微应用"},
	{Value: ROLE_DEVELOPER, Key: "developer", Name: "开发者", Description: "创建和管理微应用"},
	{Value: ROLE_ADMIN, Key: "admin", Name: "管理员", Description: "系统管理、用户管理"},
	{Value: ROLE_AUDITOR, Key: "auditor", Name: "审核员", Description: "审核微应用版本"},
	{Value: ROLE_OPERATOR, Key: "operator", Name: "运营", Description: "数据统计、推广管理"},
}

// HasRole 判断用户是否拥有某个角色
func HasRole(userRole int, role int) bool {
	return (userRole & role) != 0
}

// HasAdmin 判断用户是否是管理员
// 兼容旧系统：旧 Role=1 是管理员，新 Role=4 是管理员
func HasAdmin(userRole int) bool {
	// 新系统位运算判断
	if HasRole(userRole, ROLE_ADMIN) {
		return true
	}
	// 兼容旧系统：Role=1 是管理员
	if userRole == 1 {
		return true
	}
	return false
}

// AddRole 添加角色
func AddRole(userRole int, role int) int {
	return userRole | role
}

// RemoveRole 移除角色
func RemoveRole(userRole int, role int) int {
	return userRole &^ role
}

// GetRoleList 获取用户的所有角色列表
func GetRoleList(userRole int) []RoleInfo {
	roles := []RoleInfo{}
	for _, role := range AllRoles {
		if HasRole(userRole, role.Value) {
			roles = append(roles, role)
		}
	}
	return roles
}

// GetRoleNames 获取用户所有角色名称列表
func GetRoleNames(userRole int) []string {
	names := []string{}
	for _, role := range AllRoles {
		if HasRole(userRole, role.Value) {
			names = append(names, role.Name)
		}
	}
	return names
}

// MigrateOldRole 迁移旧角色值到新角色值
// 旧: Role=1(管理员), Role=2(普通用户)
// 新: Role=4(管理员), Role=1(普通用户)
func MigrateOldRole(oldRole int) int {
	switch oldRole {
	case 1: // 旧管理员
		return ROLE_ADMIN
	case 2: // 旧普通用户
		return ROLE_USER
	default:
		return oldRole
	}
}
