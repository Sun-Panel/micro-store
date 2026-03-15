# 微应用审核系统说明

## 概述

微应用审核系统包含两个独立的审核流程：
1. **微应用主信息审核**：应用的基本信息（名称、图标、描述、分类等）审核
2. **微应用版本审核**：应用版本包的审核

---

## 一、微应用主信息审核

### 1.1 审核流程

```
开发者修改信息 → 自动提交审核 → 管理员审核 → 审核结果
                                  ├─ 通过：更新主表信息
                                  └─ 拒绝：保留原信息
```

**注意**：修改微应用信息时会自动提交审核，无需手动提交。

### 1.2 审核状态

| 状态码 | 状态名称 | 说明 |
|-------|---------|------|
| 0 | 无审核 | 没有待审核的记录 |
| 1 | 审核中 | 已提交，等待管理员审核 |
| 2 | 已通过 | 审核通过，信息已更新 |
| 3 | 已拒绝 | 审核被拒绝，需要修改后重新提交 |

### 1.3 API 接口

#### 开发者接口（需要登录）

**撤销审核**
```typescript
POST /admin/microApp/cancelReview
{
  "id": 1  // 微应用ID
}
```

**获取审核历史**
```typescript
POST /admin/microApp/getReviewHistory
{
  "appId": 1,  // 微应用ID
  "page": 1,
  "limit": 10
}
```

#### 管理员接口（需要管理员权限）

**获取待审核列表**
```typescript
POST /admin/microApp/getPendingReviewList
{
  "page": 1,
  "limit": 10
}
```

**获取审核详情**
```typescript
POST /admin/microApp/getReviewInfo
{
  "reviewId": 1  // 审核记录ID
}
```

**审核微应用**
```typescript
POST /admin/microApp/reviewApp
{
  "reviewId": 1,      // 审核记录ID
  "status": 1,        // 1-通过 2-拒绝
  "reviewNote": "审核备注"
}
```

### 1.4 数据表结构

#### micro_app 表（主表）
新增字段：
- `review_status` - 审核状态：0-无审核 1-审核中 2-已通过 3-已拒绝
- `review_id` - 当前审核记录ID
- `review_time` - 最后审核时间

#### micro_app_review 表（审核快照表）
| 字段 | 类型 | 说明 |
|------|------|------|
| id | bigint | 主键 |
| app_id | bigint | 微应用ID |
| app_name | varchar(100) | 应用名称快照 |
| app_icon | varchar(200) | 应用图标快照 |
| app_desc | varchar(500) | 应用简介快照 |
| category_id | int | 分类ID快照 |
| charge_type | tinyint | 收费方式快照 |
| price | decimal(10,2) | 价格快照 |
| screenshots | varchar(2000) | 图集快照 |
| lang_map | json | 多语言信息快照 |
| remark | varchar(500) | 备注快照 |
| status | tinyint(2) | 审核状态：0-待审核 1-已通过 2-已拒绝 |
| reviewer_id | int | 审核人ID |
| review_note | varchar(500) | 审核备注 |
| review_time | datetime | 审核时间 |

---

## 二、微应用版本审核

### 2.1 审核流程

```
开发者上传版本 → 草稿状态 → 提交审核 → 管理员审核 → 审核结果
                                              ├─ 通过：版本公开可见
                                              └─ 拒绝：版本不可见，可修改后重新提交
```

### 2.2 审核状态

| 状态码 | 状态名称 | 说明 |
|-------|---------|------|
| -1 | 草稿 | 版本已创建但未提交审核 |
| 0 | 待审核 | 已提交，等待管理员审核 |
| 1 | 通过 | 审核通过，版本公开可见 |
| 2 | 拒绝 | 审核被拒绝 |

### 2.3 API 接口

#### 开发者接口（需要登录）

**创建版本**
```typescript
POST /developer/version/create
{
  "appId": 1,
  "version": "1.0.0",
  "versionCode": 100,
  "packageUrl": "https://...",
  "packageHash": "md5...",
  "versionDesc": "版本说明",
  "config": {...}
}
```

**提交审核**
```typescript
POST /developer/version/submitReview
{
  "versionId": 1
}
```

**撤销审核**
```typescript
POST /developer/version/cancelReview
{
  "versionId": 1
}
```

**获取版本列表**
```typescript
POST /developer/version/getList
{
  "appId": 1,
  "page": 1,
  "limit": 10,
  "status": 0  // 可选筛选状态
}
```

#### 管理员接口（需要管理员权限）

**获取待审核版本列表**
```typescript
POST /admin/version/getPendingList
{
  "page": 1,
  "limit": 10
}
```

**审核版本**
```typescript
POST /admin/version/review
{
  "versionId": 1,
  "status": 1,        // 1-通过 2-拒绝
  "reviewNote": "审核备注"
}
```

### 2.4 数据表结构

#### micro_app_version 表（版本表）
已有字段：
- `status` - 审核状态：-1-草稿 0-待审核 1-通过 2-拒绝
- `reviewer_id` - 审核人ID
- `review_note` - 审核备注
- `review_time` - 审核时间

---

## 三、数据库建表

### 执行建表SQL

直接执行以下SQL语句创建表：

#### 1. 创建微应用审核快照表

```sql
CREATE TABLE IF NOT EXISTS `micro_app_review` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `app_id` bigint(20) unsigned NOT NULL COMMENT '微应用ID',
  `app_name` varchar(100) NOT NULL COMMENT '应用名称',
  `app_icon` varchar(200) NOT NULL COMMENT '应用图标URL',
  `app_desc` varchar(500) DEFAULT NULL COMMENT '应用简介',
  `category_id` int(11) NOT NULL COMMENT '应用分类ID',
  `charge_type` tinyint(1) NOT NULL DEFAULT '0' COMMENT '收费方式：0-免费 1-积分 2-订阅PRO免费',
  `price` decimal(10,2) DEFAULT '0.00' COMMENT '价格',
  `screenshots` varchar(2000) DEFAULT NULL COMMENT '图集',
  `lang_map` json DEFAULT NULL COMMENT '多语言信息JSON',
  `remark` varchar(500) DEFAULT NULL COMMENT '应用备注',
  `status` tinyint(2) NOT NULL DEFAULT '0' COMMENT '审核状态：0-待审核 1-已通过 2-已拒绝',
  `reviewer_id` int(11) unsigned DEFAULT NULL COMMENT '审核人ID',
  `review_note` varchar(500) DEFAULT NULL COMMENT '审核备注',
  `review_time` datetime(3) DEFAULT NULL COMMENT '审核时间',
  PRIMARY KEY (`id`),
  KEY `idx_micro_app_review_app_id` (`app_id`),
  KEY `idx_micro_app_review_status` (`status`),
  KEY `idx_micro_app_review_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='微应用审核快照表';
```

#### 2. 为微应用主表添加审核相关字段

```sql
ALTER TABLE `micro_app` 
ADD COLUMN `review_status` tinyint(2) NOT NULL DEFAULT '0' COMMENT '审核状态：0-无审核 1-审核中 2-已通过 3-已拒绝' AFTER `status`,
ADD COLUMN `review_id` bigint(20) unsigned DEFAULT NULL COMMENT '当前审核记录ID' AFTER `review_status`,
ADD COLUMN `review_time` datetime(3) DEFAULT NULL COMMENT '最后审核时间' AFTER `review_id`;

ALTER TABLE `micro_app` 
ADD INDEX `idx_micro_app_review_status` (`review_status`);
```

---

## 四、使用示例

### 前端示例（开发者提交审核）

```typescript
import { submitReview, cancelReview, getReviewHistory } from '@/api/admin/microApp'
import { MicroAppReviewStatus } from '@/enums/panel/microApp'

// 提交审核
const handleSubmitReview = async (appId: number) => {
  const res = await submitReview<number>({ id: appId })
  if (res.code === 0) {
    // 提交成功，reviewId 为审核记录ID
    console.log('审核ID:', res.data)
  }
}

// 撤销审核
const handleCancelReview = async (appId: number) => {
  const res = await cancelReview({ id: appId })
  if (res.code === 0) {
    console.log('撤销成功')
  }
}

// 获取审核历史
const getHistory = async (appId: number) => {
  const res = await getReviewHistory<MicroApp.MicroAppReviewInfo>({
    appId,
    page: 1,
    limit: 10
  })
  if (res.code === 0) {
    console.log('审核历史:', res.data.list)
  }
}
```

### 前端示例（管理员审核）

```typescript
import { getPendingReviewList, reviewApp } from '@/api/admin/microApp'

// 获取待审核列表
const getPendingList = async () => {
  const res = await getPendingReviewList<MicroApp.MicroAppReviewInfo>({
    page: 1,
    limit: 10
  })
  if (res.code === 0) {
    console.log('待审核列表:', res.data.list)
  }
}

// 审核通过
const approveReview = async (reviewId: number) => {
  const res = await reviewApp({
    reviewId,
    status: 1,  // 1-通过
    reviewNote: '审核通过'
  })
  if (res.code === 0) {
    console.log('审核通过')
  }
}

// 审核拒绝
const rejectReview = async (reviewId: number) => {
  const res = await reviewApp({
    reviewId,
    status: 2,  // 2-拒绝
    reviewNote: '审核拒绝：信息不完整'
  })
  if (res.code === 0) {
    console.log('审核拒绝')
  }
}
```

---

## 五、注意事项

1. **审核快照机制**：主信息审核采用快照机制，审核通过后才会更新主表数据
2. **并发控制**：同一应用同时只能有一条待审核记录
3. **审核历史**：所有审核记录都会保存在 `micro_app_review` 表中，可追溯历史
4. **版本独立**：版本审核与主信息审核相互独立，互不影响
5. **权限控制**：开发者端只能提交和撤销审核，管理员端可以执行审核操作
6. **数据一致性**：所有审核操作都在事务中执行，确保数据一致性

---

## 六、后续优化建议

1. **邮件通知**：审核结果通过邮件通知开发者
2. **审核日志**：记录审核操作日志，便于审计
3. **批量审核**：管理员支持批量审核功能
4. **审核模板**：预设常用的审核备注模板
5. **审核统计**：统计审核通过率、平均审核时间等指标
