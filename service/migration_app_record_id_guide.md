# micro_app_review 表 app_record_id 字段迁移指南

## 背景

项目支持软删除功能，并且 `micro_app_review` 表中可能存在多条审核记录（草稿、待审核、已通过、已拒绝）。为了精确关联和追溯，需要添加 `app_record_id` 字段。

## 字段设计说明

### 两个字段的区别

1. **micro_app_id**（业务ID）
   - 类型：`varchar(50)`
   - 用途：标识同一个应用的不同版本
   - 特点：同一个业务ID在 `micro_app` 表中可能有多条记录（因为有软删除）
   - 示例：`"abc123xyz"` - 表示"天气应用"这个业务应用

2. **app_record_id**（记录ID）
   - 类型：`int(11)`
   - 用途：精确指向 `micro_app` 表中某一条具体的数据库记录
   - 特点：即使记录被软删除，也能追溯到原始记录
   - 示例：`12345` - 指向 `micro_app` 表中 `id=12345` 的具体记录

### 为什么需要两个字段？

**场景1：软删除场景**
```
micro_app 表：
id=1, micro_app_id="app1", status=1, deleted_at=NULL  (当前生效版本)
id=2, micro_app_id="app1", status=0, deleted_at=2025-01-01  (已软删除的旧版本)

micro_app_review 表：
id=10, micro_app_id="app1", app_record_id=2, status=2  (关联到已删除的记录)
id=20, micro_app_id="app1", app_record_id=1, status=1  (关联到当前记录)
```

- 如果只用 `micro_app_id`，审核记录无法区分是针对哪个数据库版本
- 使用 `app_record_id`，即使 `id=2` 的记录被软删除，审核记录仍然能追溯到它

**场景2：审核通过后的更新**
```go
// 审核通过时，根据 app_record_id 更新具体的记录
tx.Model(&models.MicroApp{}).Where("id = ?", review.AppRecordId).Updates(...)
```

## 迁移步骤

### 1. 执行 SQL 迁移脚本

```bash
cd service
mysql -u your_username -p your_database < scripts/migration_add_app_record_id.sql
```

或者直接在数据库客户端中执行 `scripts/migration_add_app_record_id.sql` 文件中的 SQL 语句。

### 2. 验证迁移结果

```sql
-- 检查字段是否添加成功
SHOW COLUMNS FROM micro_app_review LIKE 'app_record_id';

-- 检查索引是否创建成功
SHOW INDEX FROM micro_app_review WHERE Key_name = 'idx_app_record_id';

-- 检查历史数据是否正确填充
SELECT r.id, r.micro_app_id, r.app_record_id, r.status, a.id as app_table_id
FROM micro_app_review r
LEFT JOIN micro_app a ON r.micro_app_id = a.micro_app_id
WHERE r.app_record_id != a.id
LIMIT 10;
```

如果最后一个查询返回空结果，说明数据迁移正确。

### 3. 重启后端服务

```bash
# 停止当前服务
# 重新编译并启动
cd service
go build && ./sun-panel-micro-store
```

重启后，GORM 的 AutoMigrate 会自动检查表结构，确保字段和索引正确。

## 代码层面的处理

### 创建审核记录时

```go
// GetOrCreateDraftApp 方法中
draft = models.MicroAppReview{
    MicroAppBaseInfo: models.MicroAppBaseInfo{
        MicroAppId: app.MicroAppId,
        // ... 其他字段
    },
    AppRecordId: app.ID,  // 关键：设置为当前 micro_app 记录的 ID
    LangMap:     string(langMapJson),
    Status:      -1,  // 草稿状态
}
```

### 审核通过时

```go
// 审核通过，根据 app_record_id 更新具体的记录
if param.Status == 1 {
    if err := tx.Model(&models.MicroApp{}).Where("id = ?", review.AppRecordId).Updates(map[string]interface{}{
        "app_name": review.AppName,
        // ... 其他字段
    }).Error; err != nil {
        return err
    }
}
```

### 查询审核记录时

```go
// 如果需要访问关联的 micro_app 记录，考虑软删除情况
var review models.MicroAppReview
db.Where("id = ?", reviewId).First(&review)

// 查询关联的记录（可能已被软删除）
var relatedApp models.MicroApp
db.Unscoped().Where("id = ?", review.AppRecordId).First(&relatedApp)
```

## 注意事项

1. **历史数据**：迁移脚本会为所有现有审核记录填充 `app_record_id`，基于 `micro_app_id` 关联找到的 `micro_app.id`

2. **新记录**：代码中已经在创建审核记录时正确设置 `app_record_id`，无需额外处理

3. **软删除**：如果使用 `app_record_id` 关联查询，注意可能需要使用 `Unscoped()` 来查询已软删除的记录

4. **数据一致性**：确保在创建审核记录时，`app_record_id` 总是指向正确的 `micro_app` 记录ID

## 回滚方案

如果迁移出现问题，可以执行以下 SQL 回滚：

```sql
ALTER TABLE micro_app_review DROP INDEX idx_app_record_id;
ALTER TABLE micro_app_review DROP COLUMN app_record_id;
```

## 完成标志

迁移成功的标志：
1. ✅ 数据库中 `micro_app_review` 表存在 `app_record_id` 字段
2. ✅ 存在索引 `idx_app_record_id`
3. ✅ 历史数据的 `app_record_id` 已正确填充
4. ✅ 后端服务正常启动，无报错
5. ✅ 审核流程（提交、审核通过、拒绝）正常工作
