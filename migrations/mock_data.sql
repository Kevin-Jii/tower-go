-- ============================================
-- Tower Go 模拟数据
-- 执行方式: mysql -u用户名 -p密码 数据库名 < migrations/mock_data.sql
-- 注意: 请先执行 init_seed_data.sql
-- ============================================

-- 1. 供应商数据
INSERT INTO suppliers (id, supplier_code, supplier_name, contact_person, contact_phone, contact_email, supplier_address, remark, status, created_at, updated_at) VALUES
(1, '9990001', '杭州蔬菜批发中心', '王大明', '13800001001', 'wang@vegcenter.com', '杭州市余杭区勾庄农贸市场A区101', '主营各类新鲜蔬菜', 1, NOW(), NOW()),
(2, '9990002', '浙江海鲜水产', '李海波', '13800001002', 'li@seafood.com', '杭州市江干区农副产品物流中心B栋', '主营海鲜水产，每日新鲜直达', 1, NOW(), NOW()),
(3, '9990003', '金华火腿食品厂', '张金华', '13800001003', 'zhang@jinhuaham.com', '金华市婺城区火腿产业园18号', '正宗金华火腿，品质保证', 1, NOW(), NOW()),
(4, '9990004', '温州调味品公司', '陈调味', '13800001004', 'chen@seasoning.com', '温州市龙湾区食品工业园区', '各类调味品批发', 1, NOW(), NOW()),
(5, '9990005', '嘉兴粮油贸易', '刘粮油', '13800001005', 'liu@grain.com', '嘉兴市南湖区粮油批发市场', '大米、食用油批发', 1, NOW(), NOW())
ON DUPLICATE KEY UPDATE supplier_name=VALUES(supplier_name);

-- 2. 供应商分类数据
INSERT INTO supplier_categories (id, supplier_id, name, sort, status, created_at, updated_at) VALUES
-- 蔬菜供应商分类
(1, 1, '叶菜类', 1, 1, NOW(), NOW()),
(2, 1, '根茎类', 2, 1, NOW(), NOW()),
(3, 1, '瓜果类', 3, 1, NOW(), NOW()),
-- 海鲜供应商分类
(4, 2, '鱼类', 1, 1, NOW(), NOW()),
(5, 2, '虾蟹类', 2, 1, NOW(), NOW()),
(6, 2, '贝类', 3, 1, NOW(), NOW()),
-- 火腿供应商分类
(7, 3, '整腿', 1, 1, NOW(), NOW()),
(8, 3, '切片', 2, 1, NOW(), NOW()),
-- 调味品分类
(9, 4, '酱油醋', 1, 1, NOW(), NOW()),
(10, 4, '香料', 2, 1, NOW(), NOW()),
-- 粮油分类
(11, 5, '大米', 1, 1, NOW(), NOW()),
(12, 5, '食用油', 2, 1, NOW(), NOW())
ON DUPLICATE KEY UPDATE name=VALUES(name);



-- ============================================
-- 模拟数据初始化完成
-- 包含：5个供应商、12个分类
-- ============================================
