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

-- 3. 供应商商品数据
INSERT INTO supplier_products (id, supplier_id, category_id, name, unit, price, spec, remark, status, created_at, updated_at) VALUES
-- 蔬菜类商品
(1, 1, 1, '大白菜', '斤', 1.50, '新鲜', '当日采摘', 1, NOW(), NOW()),
(2, 1, 1, '小青菜', '斤', 2.00, '嫩叶', '有机种植', 1, NOW(), NOW()),
(3, 1, 1, '菠菜', '斤', 3.00, '新鲜', '', 1, NOW(), NOW()),
(4, 1, 2, '土豆', '斤', 2.50, '中等', '云南土豆', 1, NOW(), NOW()),
(5, 1, 2, '胡萝卜', '斤', 2.00, '中等', '', 1, NOW(), NOW()),
(6, 1, 2, '白萝卜', '斤', 1.80, '大', '', 1, NOW(), NOW()),
(7, 1, 3, '黄瓜', '斤', 3.50, '新鲜', '本地黄瓜', 1, NOW(), NOW()),
(8, 1, 3, '西红柿', '斤', 4.00, '新鲜', '自然熟', 1, NOW(), NOW()),
-- 海鲜类商品
(9, 2, 4, '草鱼', '斤', 12.00, '活鱼', '现杀现卖', 1, NOW(), NOW()),
(10, 2, 4, '鲈鱼', '斤', 25.00, '活鱼', '海鲈鱼', 1, NOW(), NOW()),
(11, 2, 4, '带鱼', '斤', 18.00, '冰鲜', '', 1, NOW(), NOW()),
(12, 2, 5, '基围虾', '斤', 45.00, '活虾', '南美白虾', 1, NOW(), NOW()),
(13, 2, 5, '大闸蟹', '只', 35.00, '4两公', '阳澄湖', 1, NOW(), NOW()),
(14, 2, 6, '花蛤', '斤', 8.00, '新鲜', '', 1, NOW(), NOW()),
-- 火腿类商品
(15, 3, 7, '金华火腿整腿', '只', 580.00, '3年陈', '传统工艺', 1, NOW(), NOW()),
(16, 3, 8, '火腿切片', '包', 68.00, '200g/包', '即食', 1, NOW(), NOW()),
-- 调味品商品
(17, 4, 9, '生抽', '瓶', 12.00, '500ml', '海天', 1, NOW(), NOW()),
(18, 4, 9, '老抽', '瓶', 14.00, '500ml', '海天', 1, NOW(), NOW()),
(19, 4, 9, '香醋', '瓶', 8.00, '500ml', '镇江香醋', 1, NOW(), NOW()),
(20, 4, 10, '八角', '袋', 15.00, '100g', '', 1, NOW(), NOW()),
-- 粮油商品
(21, 5, 11, '东北大米', '袋', 65.00, '10kg', '五常大米', 1, NOW(), NOW()),
(22, 5, 11, '泰国香米', '袋', 85.00, '10kg', '进口', 1, NOW(), NOW()),
(23, 5, 12, '花生油', '桶', 128.00, '5L', '鲁花', 1, NOW(), NOW()),
(24, 5, 12, '菜籽油', '桶', 98.00, '5L', '金龙鱼', 1, NOW(), NOW())
ON DUPLICATE KEY UPDATE name=VALUES(name);

-- 4. 门店绑定供应商商品（示例门店1绑定部分商品）
INSERT INTO store_supplier_products (id, store_id, product_id, is_default, created_at, updated_at) VALUES
(1, 1, 1, 1, NOW(), NOW()),
(2, 1, 2, 0, NOW(), NOW()),
(3, 1, 4, 1, NOW(), NOW()),
(4, 1, 7, 1, NOW(), NOW()),
(5, 1, 9, 1, NOW(), NOW()),
(6, 1, 12, 0, NOW(), NOW()),
(7, 1, 17, 1, NOW(), NOW()),
(8, 1, 21, 1, NOW(), NOW()),
(9, 1, 23, 1, NOW(), NOW())
ON DUPLICATE KEY UPDATE is_default=VALUES(is_default);

-- 5. 菜品分类数据
INSERT INTO dish_categories (id, store_id, name, sort, status, created_at, updated_at) VALUES
(1, 1, '凉菜', 1, 1, NOW(), NOW()),
(2, 1, '热菜', 2, 1, NOW(), NOW()),
(3, 1, '汤类', 3, 1, NOW(), NOW()),
(4, 1, '主食', 4, 1, NOW(), NOW()),
(5, 2, '凉菜', 1, 1, NOW(), NOW()),
(6, 2, '热菜', 2, 1, NOW(), NOW())
ON DUPLICATE KEY UPDATE name=VALUES(name);

-- 6. 菜品数据
INSERT INTO dishes (id, store_id, category_id, name, price, unit, status, remark, created_at, updated_at) VALUES
(1, 1, 1, '凉拌黄瓜', 12.00, '份', 1, '清爽开胃', NOW(), NOW()),
(2, 1, 1, '皮蛋豆腐', 18.00, '份', 1, '', NOW(), NOW()),
(3, 1, 2, '红烧肉', 48.00, '份', 1, '招牌菜', NOW(), NOW()),
(4, 1, 2, '清蒸鲈鱼', 68.00, '条', 1, '新鲜活鱼', NOW(), NOW()),
(5, 1, 2, '宫保鸡丁', 38.00, '份', 1, '', NOW(), NOW()),
(6, 1, 3, '番茄蛋汤', 18.00, '份', 1, '', NOW(), NOW()),
(7, 1, 4, '米饭', 3.00, '碗', 1, '', NOW(), NOW()),
(8, 2, 5, '凉拌木耳', 15.00, '份', 1, '', NOW(), NOW()),
(9, 2, 6, '糖醋排骨', 52.00, '份', 1, '', NOW(), NOW())
ON DUPLICATE KEY UPDATE name=VALUES(name);

-- ============================================
-- 模拟数据初始化完成
-- 包含：5个供应商、12个分类、24个商品、9个门店绑定、6个菜品分类、9个菜品
-- ============================================
