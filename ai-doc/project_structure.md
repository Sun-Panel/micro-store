# 项目目录结构说明

## 根目录

```
sun-panel-micro-store/
├── ai-doc/              # AI 文档目录（需求、设计文档等）
├── cache/               # 缓存相关
├── config/              # 配置文件目录
├── dist/                # 前端构建输出目录
├── doc/                 # 项目文档
├── pkg/                 # 公共包
├── public/              # 静态资源目录
├── service/             # 后端服务（Go）
├── src/                 # 前端源码（Vue + TypeScript）
├── build.sh             # 构建脚本
├── docker-compose.yml   # Docker 编排文件
├── Dockerfile           # Docker 镜像构建文件
├── package.json         # 前端依赖配置
├── vite.config.ts       # Vite 配置文件
└── tsconfig.json        # TypeScript 配置文件
```

---

## 后端目录结构（service/）

### 核心目录

```
service/
├── main.go              # 程序入口
├── go.mod               # Go 模块依赖
├── adapter/             # 适配器（数据转换）
├── api/                 # API 控制器（处理HTTP请求）
│   └── api_v1/          # API v1 版本控制器
├── apiClientApp/        # API 客户端应用
├── biz/                 # 业务逻辑层
├── conf/                # 配置文件
├── global/              # 全局变量
├── initialize/          # 初始化模块
│   ├── authService/     # 认证服务初始化
│   ├── config/          # 配置初始化
│   ├── database/        # 数据库初始化（创建表）
│   ├── lang/            # 语言初始化
│   ├── redis/           # Redis 初始化
│   └── ...              # 其他初始化
├── lang/                # 多语言文件
├── lib/                 # 工具库
│   ├── AES/             # AES 加密
│   ├── cache/           # 缓存工具
│   ├── captcha/         # 验证码
│   ├── cmn/             # 通用工具
│   ├── language/        # 语言处理
│   ├── mail/            # 邮件发送
│   ├── pay/             # 支付相关
│   ├── queue/           # 队列
│   ├── sunStore/        # 商店相关工具
│   └── ...              # 其他工具
├── models/              # 数据模型（GORM模型）
│   ├── base.go          # 基础模型和数据库连接
│   ├── User.go          # 用户模型
│   ├── microApp.go      # 微应用模型
│   ├── microAppCategory.go # 应用分类模型
│   ├── developer.go     # 开发者模型
│   └── ...              # 其他模型
├── router/              # 路由配置
│   ├── admin/           # 管理员路由
│   ├── oauth2/          # OAuth2 路由
│   ├── openness/        # 开放接口路由
│   ├── panel/           # 面板路由
│   ├── sunStore/        # 商店路由
│   └── system/          # 系统路由
├── runtime/             # 运行时文件（日志等）
├── scheduler/           # 定时任务
├── scripts/             # 脚本文件
└── structs/             # 结构体定义（DTO等）
```

### 后端关键文件说明

- **`models/`** - 数据表模型，每个模型对应一个数据表，包含表结构和增删改查方法
- **`api/`** - API 控制器，处理 HTTP 请求，调用 biz 层处理业务
- **`biz/`** - 业务逻辑层，处理复杂业务逻辑
- **`router/`** - 路由配置，定义 API 路径和对应的处理器
- **`initialize/database/`** - 数据库初始化，AutoMigrate 创建数据表
- **`lib/`** - 工具库，包含各种通用功能

---

## 前端目录结构（src/）

```
src/
├── main.ts              # 前端入口文件
├── App.vue              # 根组件
├── api/                 # API 接口定义
│   ├── admin/           # 管理员 API
│   ├── login.ts         # 登录 API
│   └── ...              # 其他 API
├── assets/              # 静态资源（图片、图标等）
├── components/          # 公共组件
├── enums/               # 枚举定义
├── hooks/               # Vue Composables（组合式函数）
├── icons/               # 图标组件
├── locales/             # 国际化文件
│   ├── zh-CN.json       # 中文语言包
│   └── en-US.json       # 英文语言包
├── plugins/             # Vue 插件
├── router/              # 路由配置
├── store/               # 状态管理（Pinia/Vuex）
│   └── modules/         # 状态模块
├── styles/              # 样式文件
├── typings/             # TypeScript 类型定义
│   ├── admin/           # 管理员相关类型
│   ├── openness/        # 开放接口类型
│   └── ...              # 其他类型
├── utils/               # 工具函数
└── views/               # 页面组件
    ├── admin/           # 管理员页面
    ├── home/            # 首页
    ├── login/           # 登录页
    ├── platform/        # 平台页面
    └── ...              # 其他页面
```

### 前端关键文件说明

- **`api/`** - API 接口封装，定义与后端交互的方法
- **`typings/`** - TypeScript 类型定义，定义数据结构
- **`views/`** - 页面组件，每个路由对应的页面
- **`components/`** - 可复用的组件
- **`store/`** - 全局状态管理
- **`router/`** - 前端路由配置

---

## 数据库模型命名规范

### 模型文件命名
- 文件名采用驼峰命名：`microApp.go`, `microAppCategory.go`
- 结构体名采用大驼峰：`MicroApp`, `MicroAppCategory`

### 数据表命名
- 表名采用下划线命名：`micro_app`, `micro_app_category`
- 通过 `TableName()` 方法定义表名

### 模型方法规范
```go
// 所有方法的第一个参数都是 *gorm.DB
func (m *Model) GetList(db *gorm.DB, ...) ([]Model, int64, error)
func (m *Model) GetById(db *gorm.DB, id uint) (Model, error)
func (m *Model) Create(db *gorm.DB) error
func (m *Model) Update(db *gorm.DB, id uint, data map[string]interface{}) error
func (m *Model) Delete(db *gorm.DB, ids []uint) error
```

---

## 开发流程

### 1. 新增数据表
1. 在 `service/models/` 创建模型文件
2. 在 `service/initialize/database/connect.go` 的 `AutoMigrate` 中注册模型
3. 程序启动时自动创建表

### 2. 新增 API 接口
1. 在 `service/models/` 创建或修改模型，添加数据库操作方法
2. 在 `service/api/api_v1/` 创建或修改 API 控制器
3. 在 `service/router/` 注册路由
4. 在 `src/api/` 创建前端 API 接口
5. 在 `src/typings/` 定义类型

### 3. 新增页面
1. 在 `src/views/` 创建页面组件
2. 在 `src/router/` 注册路由
3. 在 `src/locales/` 添加国际化文本

---

## 技术栈

### 后端
- **语言**: Go
- **Web框架**: Gin
- **ORM**: GORM
- **数据库**: MySQL / SQLite
- **缓存**: Redis
- **认证**: OAuth2

### 前端
- **框架**: Vue 3
- **语言**: TypeScript
- **构建工具**: Vite
- **UI框架**: Element Plus / Tailwind CSS
- **状态管理**: Pinia
- **路由**: Vue Router

---

## 重要配置文件

- `service/conf/` - 后端配置文件目录
- `config/` - 前端配置文件目录
- `.env` - 环境变量配置（如有）
- `docker-compose.yml` - Docker 编排配置

---

## 微应用商店相关

### 数据表（已创建）
- `micro_app` - 微应用列表
- `micro_app_version` - 微应用版本列表
- `micro_app_category` - 应用分类
- `micro_app_lang` - 微应用多语言信息
- `micro_app_download` - 下载记录
- `micro_app_install` - 安装记录
- `developer` - 开发者表

### 模型文件
- `service/models/microApp.go` - 微应用模型（支持多语言查询）
- `service/models/microAppCategory.go` - 分类模型
- `service/models/microAppVersion.go` - 版本模型
- `service/models/microAppDownload.go` - 下载记录模型
- `service/models/microAppInstall.go` - 安装记录模型
- `service/models/microAppLang.go` - 多语言模型
- `service/models/developer.go` - 开发者模型
