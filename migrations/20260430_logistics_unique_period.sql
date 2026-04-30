-- =============================================================================
-- 已有库：物流单唯一约束从「路线 + 保存日」改为「路线 + 导入起止日」
-- 与代码 SaveLogisticsSheet 按 (route_id, start_date, end_date) 覆盖保存一致。
--
-- 执行前：
-- 1) 备份表数据；
-- 2) 若存在多条相同 (route_id, start_date, end_date)，先合并/删除重复行，否则 ADD UNIQUE 会失败。
-- =============================================================================

-- 1) 导入区间列改为 NOT NULL（空串表示未填，与 Go 校验一致）
ALTER TABLE `third_party_logistics_sheets`
  MODIFY `start_date` VARCHAR(10) NOT NULL DEFAULT '',
  MODIFY `end_date` VARCHAR(10) NOT NULL DEFAULT '';

-- 2) 删除「仅由 route_id + sheet_date 两列组成的唯一索引」
--    不硬编码索引名：GORM/手工建表可能叫 idx_route_sheet_day 或其它名字。
SET @dn = DATABASE();

SELECT s.INDEX_NAME INTO @drop_name
FROM information_schema.STATISTICS s
WHERE s.TABLE_SCHEMA = @dn
  AND s.TABLE_NAME = 'third_party_logistics_sheets'
  AND s.NON_UNIQUE = 0
  AND s.INDEX_NAME <> 'PRIMARY'
GROUP BY s.INDEX_NAME
HAVING COUNT(*) = 2
  AND MAX(CASE WHEN s.SEQ_IN_INDEX = 1 THEN s.COLUMN_NAME END) = 'route_id'
  AND MAX(CASE WHEN s.SEQ_IN_INDEX = 2 THEN s.COLUMN_NAME END) = 'sheet_date'
LIMIT 1;

SET @sql = IF(
  IFNULL(@drop_name, '') <> '',
  CONCAT('ALTER TABLE `third_party_logistics_sheets` DROP INDEX `', @drop_name, '`'),
  'SELECT 1 AS `_skip_drop_old_route_sheet_unique`'
);
PREPARE stmt FROM @sql;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- 3) 新唯一索引（若已存在则跳过）
SELECT COUNT(*) INTO @has_period
FROM information_schema.STATISTICS
WHERE TABLE_SCHEMA = @dn
  AND TABLE_NAME = 'third_party_logistics_sheets'
  AND INDEX_NAME = 'uk_route_logistics_period';

SET @sql = IF(
  @has_period = 0,
  'ALTER TABLE `third_party_logistics_sheets` ADD UNIQUE KEY `uk_route_logistics_period` (`route_id`, `start_date`, `end_date`)',
  'SELECT 1 AS `_skip_add_uk_route_logistics_period`'
);
PREPARE stmt FROM @sql;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

-- 4) 保存日普通索引（若已存在则跳过）
SELECT COUNT(*) INTO @has_sheet_date
FROM information_schema.STATISTICS
WHERE TABLE_SCHEMA = @dn
  AND TABLE_NAME = 'third_party_logistics_sheets'
  AND INDEX_NAME = 'idx_logistics_sheets_sheet_date';

SET @sql = IF(
  @has_sheet_date = 0,
  'ALTER TABLE `third_party_logistics_sheets` ADD KEY `idx_logistics_sheets_sheet_date` (`sheet_date`)',
  'SELECT 1 AS `_skip_add_idx_logistics_sheets_sheet_date`'
);
PREPARE stmt FROM @sql;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;
