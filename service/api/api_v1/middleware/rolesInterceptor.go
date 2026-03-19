package middleware

import (
	"fmt"
	"strings"
	"sun-panel/api/api_v1/common/apiReturn"
	"sun-panel/api/api_v1/common/base"
	"sun-panel/models"

	"github.com/gin-gonic/gin"
)

// 角色验证配置
type RoleConfig struct {
	RequiredRole int    // 需要的角色
	RoleName     string // 角色名称（用于错误信息）
}

// RolesInterceptor 通用角色拦截器
// 使用方式：middleware.RolesInterceptor(requiredRole, roleName)
func RolesInterceptor(requiredRole int, roleName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		currentUser, exists := base.GetCurrentUserInfo(c)
		if !exists {
			apiReturn.ErrorNotLogin(c)
			c.Abort()
			return
		}

		// 检查用户是否拥有所需角色
		if !models.HasRole(currentUser.Role, requiredRole) {
			apiReturn.ErrorByCodeAndMsg(c, 403, "需要"+roleName+"权限")
			c.Abort()
			return
		}
	}
}

// 预定义的角色拦截器（方便使用）
var (
	// UserInterceptor 需要普通用户角色
	UserInterceptor = RolesInterceptor(models.ROLE_USER, "普通用户")
	
	// DeveloperInterceptor 需要开发者角色
	DeveloperInterceptor = RolesInterceptor(models.ROLE_DEVELOPER, "开发者")
	
	// AdminInterceptor 需要管理员角色（已存在，但保持兼容）
	// AdminInterceptor = RolesInterceptor(models.ROLE_ADMIN, "管理员")
	
	// AuditorInterceptor 需要审核员角色
	AuditorInterceptor = RolesInterceptor(models.ROLE_AUDITOR, "审核员")
	
	// OperatorInterceptor 需要运营角色
	OperatorInterceptor = RolesInterceptor(models.ROLE_OPERATOR, "运营")
	
	// AdminOrAuditorInterceptor 需要管理员或审核员角色
	AdminOrAuditorInterceptor = MultiRolesInterceptor(models.ROLE_ADMIN | models.ROLE_AUDITOR)
	
	// AdminOrDeveloperInterceptor 需要管理员或开发者角色
	AdminOrDeveloperInterceptor = MultiRolesInterceptor(models.ROLE_ADMIN | models.ROLE_DEVELOPER)
	
	// DeveloperOrAuditorInterceptor 需要开发者或审核员角色
	DeveloperOrAuditorInterceptor = MultiRolesInterceptor(models.ROLE_DEVELOPER | models.ROLE_AUDITOR)
)

// MultiRolesInterceptor 多角色拦截器（满足任意一个角色即可）
// 使用位运算：requiredRolesMask = role1 | role2 | role3 ...
func MultiRolesInterceptor(requiredRolesMask int) gin.HandlerFunc {
	return func(c *gin.Context) {
		currentUser, exists := base.GetCurrentUserInfo(c)
		if !exists {
			apiReturn.ErrorNotLogin(c)
			c.Abort()
			return
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
			c.Abort()
			return
		}
	}
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
