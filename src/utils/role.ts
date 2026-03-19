// 角色常量定义
export const ROLE_USER = 1      // 普通用户
export const ROLE_DEVELOPER = 2 // 开发者
export const ROLE_ADMIN = 4     // 管理员
export const ROLE_AUDITOR = 8   // 审核员
export const ROLE_OPERATOR = 16 // 运营

// 角色信息接口
export interface RoleInfo {
  value: number
  key: string
  name: string
  description: string
}

// 所有角色列表
export const ALL_ROLES: RoleInfo[] = [
  { value: ROLE_USER, key: 'user', name: '普通用户', description: '浏览、下载、安装微应用' },
  { value: ROLE_DEVELOPER, key: 'developer', name: '开发者', description: '创建和管理微应用' },
  { value: ROLE_ADMIN, key: 'admin', name: '管理员', description: '系统管理、用户管理' },
  { value: ROLE_AUDITOR, key: 'auditor', name: '审核员', description: '审核微应用版本' },
  { value: ROLE_OPERATOR, key: 'operator', name: '运营', description: '数据统计、推广管理' },
]

/**
 * 判断用户是否拥有某个角色
 */
export function hasRole(userRole: number, role: number): boolean {
  return (userRole & role) !== 0
}

/**
 * 判断用户是否是管理员
 */
export function hasAdmin(userRole: number): boolean {
  return hasRole(userRole, ROLE_ADMIN)
}

/**
 * 添加角色
 */
export function addRole(userRole: number, role: number): number {
  return userRole | role
}

/**
 * 移除角色
 */
export function removeRole(userRole: number, role: number): number {
  return userRole & ~role
}

/**
 * 获取用户的所有角色列表
 */
export function getRoleList(userRole: number): RoleInfo[] {
  return ALL_ROLES.filter(role => hasRole(userRole, role.value))
}

/**
 * 获取用户所有角色名称列表
 */
export function getRoleNames(userRole: number): string[] {
  return ALL_ROLES.filter(role => hasRole(userRole, role.value)).map(role => role.name)
}

/**
 * 根据角色值获取角色标签颜色
 */
export function getRoleTagType(role: number): 'default' | 'error' | 'primary' | 'success' | 'warning' {
  switch (role) {
    case ROLE_ADMIN:
      return 'error'
    case ROLE_DEVELOPER:
      return 'primary'
    case ROLE_AUDITOR:
      return 'success'
    case ROLE_OPERATOR:
      return 'warning'
    default:
      return 'default'
  }
}

/**
 * 切换角色
 */
export function toggleRole(userRole: number, role: number): number {
  if (hasRole(userRole, role)) {
    return removeRole(userRole, role)
  }
  return addRole(userRole, role)
}

/**
 * 计算选中的角色列表（用于复选框组）
 */
export function getSelectedRoles(userRole: number): number[] {
  return ALL_ROLES.filter(role => hasRole(userRole, role.value)).map(role => role.value)
}

/**
 * 从选中列表计算角色值
 */
export function calculateRoleFromSelected(selectedRoles: number[]): number {
  return selectedRoles.reduce((acc, role) => acc | role, 0)
}

// ==================== 菜单权限 ====================

/**
 * 后台管理菜单项配置
 * - key: 路由名称（叶子节点必填）
 * - label: 菜单显示文本
 * - children: 子菜单
 * - roles: 需要的角色（多个角色用位运算组合），父级必填，子项可选（继承父级）
 */

// 叶子菜单项（无children）
export interface LeafMenuItem {
  key: string
  label: string
  roles?: number
}

// 分组菜单项（有children）
export interface GroupMenuItem {
  key: string
  label: string
  children: AdminMenuItem[]
  roles: number
}

// 联合类型
export type AdminMenuItem = LeafMenuItem | GroupMenuItem

/**
 * 递归过滤菜单：检查用户角色权限，子项未设置roles时继承父项
 */
export function filterAdminMenuByRole(
  userRole: number,
  menu: AdminMenuItem,
  parentRoles?: number
): AdminMenuItem | undefined {
  // 当前项的角色：自身roles或继承父项
  const currentRoles = menu.roles ?? parentRoles

  // 有children的是分组菜单，递归处理子项
  if ('children' in menu && menu.children && menu.children.length > 0) {
    const filteredChildren = menu.children
      .map(child => filterAdminMenuByRole(userRole, child, currentRoles))
      .filter(Boolean) as AdminMenuItem[]

    // 有可见子菜单时才返回
    if (filteredChildren.length > 0) {
      return { ...menu, children: filteredChildren }
    }
    return undefined
  }

  // 叶子菜单项，检查权限
  return hasRole(userRole, currentRoles!) ? menu : undefined
}

/**
 * 过滤并转换菜单为 naive-ui 格式
 */
export function filterAndConvertMenu(
  userRole: number,
  menuConfig: AdminMenuItem[]
): any[] {
  return menuConfig
    .map(item => filterAdminMenuByRole(userRole, item))
    .filter(Boolean)
}
