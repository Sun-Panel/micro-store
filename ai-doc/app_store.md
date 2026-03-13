# 微应用商店

这是sun-panel的微应用商店，主要功能有，微应用发布、审核、下载、安装等。审核包含机器审核。后期会增加评价、举报等功能。

## 数据表

### 微应用列表 `micro_app`

字段|类型|必填|描述
---|---|---|---
id|int|是|主键
app_name|string|是|应用名称
app_icon|string|是|应用图标URL
app_desc|string|否|应用简介
category_id|int|是|应用分类ID
charge_type|int|是|收费方式：0-免费 1-付费 2-订阅
price|decimal|否|价格（付费时）
author_id|int|是|开发者ID
permission_level|int|否|应用权限等级
status|int|是|状态：0-下架 1-上架 2-审核中
create_time|datetime|是|创建时间
update_time|datetime|是|更新时间

### 微应用版本列表 `micro_app_version`

字段|类型|必填|描述
---|---|---|---
id|int|是|主键
app_id|int|是|微应用ID
version|string|是|版本号（如 1.0.0）
version_code|int|是|版本号（数字）
package_url|string|是|应用包下载地址
package_hash|string|是|版本校验值（MD5/SHA）
status|int|是|审核状态：0-待审核 1-通过 2-拒绝
review_time|datetime|否|审核时间
reviewer_id|int|否|审核人ID
review_note|string|否|审核备注
create_time|datetime|是|创建时间

### 下载记录表 `micro_app_download`

字段|类型|必填|描述
---|---|---|---
id|int|是|主键
app_id|int|是|微应用ID
version_id|int|是|版本ID
user_id|int|否|用户ID（匿名下载为空）
client_id|string|是|客户端标识
download_ip|string|是|下载IP
download_device|string|否|下载设备信息
download_client|string|否|下载客户端类型
download_time|datetime|是|下载时间

### 用户表 `user`

字段|类型|必填|描述
---|---|---|---
id|int|是|主键
username|string|是|用户名
password|string|是|密码（加密）
mail|string|否|邮箱
head_image|string|否|头像URL
role|int|是|角色（位运算，可叠加）
status|int|是|状态：0-禁用 1-正常
create_time|datetime|是|创建时间
is_bind_sun_store|bool|否|是否绑定商店

**方案对比**

| 方案 | 优点 | 缺点 | 适用场景 |
|------|------|------|----------|
| 位运算 | 存储高效、查询快、代码简洁 | 最多32个角色、角色固定 | 角色少且稳定 |
| RBAC权限表 | 灵活扩展、动态分配、支持复杂权限 | 表关联复杂、查询稍慢 | 角色多且易变 |

**推荐方案：位运算（当前阶段）**
- 微应用商店角色简单，短期内不会超过10个
- 位运算满足需求，开发成本低
- 未来如需复杂权限，可平滑迁移到RBAC

**角色定义（位运算）**
```
ROLE_USER      = 1   // 普通用户   (二进制: 00001)
ROLE_DEVELOPER = 2   // 开发者     (二进制: 00010)
ROLE_ADMIN     = 4   // 管理员     (二进制: 00100)
// 预留扩展
ROLE_AUDITOR   = 8   // 审核员     (二进制: 01000)
ROLE_OPERATOR  = 16  // 运营       (二进制: 10000)
```

**新增角色步骤**
1. 定义新角色值（下一个2的幂次方）
2. 在代码中添加常量定义
3. 数据库无需改动，直接使用

**角色叠加示例**
- 普通用户: 1
- 开发者: 2
- 管理员: 4
- 开发者+管理员: 6 (2+4, 二进制: 110)
- 审核员+管理员: 12 (8+4, 二进制: 1100)

**权限判断（使用位与运算）**
```javascript
// 判断是否包含某个角色
function hasRole(userRole, role) {
  return (userRole & role) !== 0
}

// 判断是否是管理员
isAdmin = (user.role & 4) !== 0 // 或 user.role & 4 === 4

// 添加角色
user.role |= 2 // 添加开发者角色

// 移除角色
user.role &= ~2 // 移除开发者角色
```

**未来扩展：RBAC权限表（备选方案）**
```
用户角色关联表 user_role
- id, user_id, role_id, create_time

角色表 role
- id, role_name, role_key, description, status

权限表 permission
- id, permission_name, permission_key, description

角色权限关联表 role_permission
- id, role_id, permission_id
```

### 开发者表 `developer`

字段|类型|必填|描述
---|---|---|---
id|int|是|主键
user_id|int|是|用户ID
developer_name|string|是|开发者名称
contact_mail|string|否|联系邮箱
status|int|是|状态：0-禁用 1-正常
create_time|datetime|是|创建时间

### 应用安装记录表 `micro_app_install`

字段|类型|必填|描述
---|---|---|---
id|int|是|主键
app_id|int|是|微应用ID
version_id|int|是|版本ID
user_id|int|否|用户ID
client_id|string|是|客户端标识
intranet_ip|string|否|内网IP
public_ip|string|是|公网IP
install_time|datetime|是|安装时间
user_is_pro|bool|否|安装时用户是否为PRO
point_value|int|否|本次积分值
author_point_rule|string|否|作者当前积分规则JSON

### 应用分类表 `micro_app_category`

字段|类型|必填|描述
---|---|---|---
id|int|是|主键
name|string|是|分类名称
icon|string|否|分类图标
sort|int|否|排序
status|int|是|状态：0-禁用 1-正常

### 用户积分表 `user_points`（暂不开发）

字段|类型|必填|描述
---|---|---|---
id|int|是|主键
user_id|int|是|用户ID
points|int|是|积分值
source_type|int|是|来源类型：1-安装 2-购买 3-充值
source_id|int|否|来源ID
app_id|int|否|关联应用ID
version_id|int|否|关联版本ID
create_time|datetime|是|积分时间

### 用户已购应用表 `user_purchased_apps`（暂不开发）

字段|类型|必填|描述
---|---|---|---
id|int|是|主键
user_id|int|是|用户ID
app_id|int|是|应用ID
purchase_time|datetime|是|购买时间
expire_time|datetime|否|过期时间（订阅制）

## API 列表

### 微应用相关
- `POST /microApp/list` - 获取微应用列表（支持分类、搜索筛选）
- `POST /microApp/detail` - 获取微应用详情
- `POST /microApp/download` - 下载微应用包
- `POST /microApp/install` - 记录安装信息

### 开发者相关
- `POST /developer/register` - 注册成为开发者
- `POST /developer/apps` - 获取我的应用列表
- `POST /developer/app/create` - 创建微应用
- `POST /developer/app/update` - 更新微应用信息
- `POST /developer/version/create` - 创建新版本
- `POST /developer/version/list` - 获取版本列表

### 管理员相关
- `POST /admin/microApp/list` - 管理应用列表
- `POST /admin/microApp/review` - 审核应用
- `POST /admin/version/review` - 审核版本
- `POST /admin/category/list` - 分类列表
- `POST /admin/category/create` - 创建分类
- `POST /admin/category/update` - 更新分类
- `POST /admin/developer/list` - 开发者列表
- `POST /admin/developer/updateStatus` - 更新开发者状态

### 用户相关
- `POST /user/myApps` - 获取已安装应用列表
- `POST /user/purchasedApps` - 获取已购应用列表

## 页面功能

### 用户端
- **应用市场首页** - 应用分类展示、推荐应用、搜索
- **应用详情页** - 应用信息、版本历史、下载安装
- **我的应用** - 已安装应用列表、更新管理
- **开发者中心** - 应用管理、版本管理、数据统计

### 管理端
- **应用审核页** - 待审核应用列表、审核操作、审核历史
- **版本审核页** - 待审核版本列表、版本详情、审核操作
- **分类管理页** - 分类CRUD、排序管理
- **开发者管理页** - 开发者列表、状态管理
- **数据统计页** - 下载量、安装量、收入统计

### 机器审核
- 自动化版本包扫描（病毒、敏感内容）
- 自动化代码安全检测
- 审核结果回调通知
