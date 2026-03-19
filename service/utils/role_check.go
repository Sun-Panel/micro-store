package utils

import (
	"fmt"
	"strings"
	"sun-panel/api/api_v1/common/apiReturn"
	"sun-panel/api/api_v1/common/base"
	"sun-panel/models"

	"github.com/gin-gonic/gin"
)

// CheckRole 检查用户是否拥有指定角色
func CheckRole(c *gin.Context, requiredRole int, roleName string) bool {
	currentUser, exists := base.GetCurrentUserInfo(c)
	if !exists {
		apiReturn.ErrorNotLogin(c)
		return false
	}

	if !models.HasRole(currentUser.Role, requiredRole) {
		apiReturn.ErrorByCodeAndMsg(c, 403, "需要"+roleName+"权限")
		return false
	}
	return true
}

// CheckRoles 检查用户是否拥有指定角色之一
// 使用位运算：requiredRolesMask = role1 | role2 | role3 ...
func CheckRoles(c *gin.Context, requiredRolesMask int) bool {
	currentUser, exists := base.GetCurrentUserInfo(c)
	if !exists {
		apiReturn.ErrorNotLogin(c)
		return false
	}

	// 使用位运算检查用户是否拥有任意一个所需角色
	if (currentUser.Role & requiredRolesMask) == 0 {
		// 根据角色掩码构建角色名称字符串
		roleNames := getRoleNamesByMask(requiredRolesMask)
		var namesBuilder strings.Builder
		for i, name := range roleNames {
			if i > 0 {
				namesBuilder.WriteString("或")
			}
			namesBuilder.WriteString(name)
		}
		apiReturn.ErrorByCodeAndMsg(c, 403, "需要"+namesBuilder.String()+"权限")
		return false
	}
	return true
}

// getRoleNameByValue 根据角色值获取角色名称
func getRoleNameByValue(roleValue int) string {
	for _, role := range models.AllRoles {
		if role.Value == roleValue {
			return role.Name
		}
	}
	return fmt.Sprintf("未知角色(%d)", roleValue)
}

// getRoleNamesByMask 根据角色掩码获取角色名称列表
func getRoleNamesByMask(roleMask int) []string {
	var roleNames []string
	for _, role := range models.AllRoles {
		if (roleMask & role.Value) != 0 {
			roleNames = append(roleNames, role.Name)
		}
	}
	return roleNames
}

// GetUserRoles 获取用户的所有角色列表
func GetUserRoles(c *gin.Context) []models.RoleInfo {
	currentUser, exists := base.GetCurrentUserInfo(c)
	if !exists {
		return []models.RoleInfo{}
	}
	return models.GetRoleList(currentUser.Role)
}

// GetUserRoleNames 获取用户的所有角色名称列表
func GetUserRoleNames(c *gin.Context) []string {
	currentUser, exists := base.GetCurrentUserInfo(c)
	if !exists {
		return []string{}
	}
	return models.GetRoleNames(currentUser.Role)
}

// HasAdminRole 检查用户是否是管理员
func HasAdminRole(c *gin.Context) bool {
	currentUser, exists := base.GetCurrentUserInfo(c)
	if !exists {
		return false
	}
	return models.HasAdmin(currentUser.Role)
}

// HasDeveloperRole 检查用户是否是开发者
func HasDeveloperRole(c *gin.Context) bool {
	return CheckRole(c, models.ROLE_DEVELOPER, "开发者")
}

// HasAuditorRole 检查用户是否是审核员
func HasAuditorRole(c *gin.Context) bool {
	return CheckRole(c, models.ROLE_AUDITOR, "审核员")
}

// HasOperatorRole 检查用户是否是运营
func HasOperatorRole(c *gin.Context) bool {
	return CheckRole(c, models.ROLE_OPERATOR, "运营")
}