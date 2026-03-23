-- 迁移脚本：为 micro_app_review 表添加 app_record_id 字段
-- 日期：2026-03-23

-- 1. 添加 app_record_id 字段
ALTER TABLE micro_app_review ADD COLUMN app_record_id INT(11) NOT NULL DEFAULT 0 AFTER lang_map;

-- 2. 添加索引
ALTER TABLE micro_app_review ADD INDEX idx_app_record_id (app_record_id);

-- 3. 为历史数据填充 app_record_id
-- 根据 micro_app_id 关联到 micro_app 表的 id
UPDATE micro_app_review r
INNER JOIN micro_app a ON r.micro_app_id = a.micro_app_id
SET r.app_record_id = a.id
WHERE r.app_record_id = 0;

-- 说明：
-- - app_record_id 字段用于精确关联 micro_app 表的具体记录（而不是业务ID）
-- - 这样设计是为了支持软删除：即使 micro_app 记录被软删除，审核记录也能追溯到具体记录
-- - 对于新创建的审核记录，代码中已经正确设置 app_record_id
-- - 对于历史数据，此脚本会根据 micro_app_id 找到对应的记录ID
