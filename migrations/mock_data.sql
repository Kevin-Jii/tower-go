-- ============================================
-- Tower Go Mock数据
-- 执行方式: mysql -u用户名 -p密码 数据库名 < migrations/mock_data.sql
-- ============================================

-- 1. 供应商数据
INSERT INTO suppliers (id, supplier_code, supplier_name, contact_person, contact_phone, contact_email, supplier_address, remark, status, created_at, updated_at) VALUES
(1, 'GYS001', '杭州鲜蔬农产品有限公司', '张经理', '13800001001', 'zhang@xiansu.com', '杭州市余杭区农贸市场A区', '主营蔬菜类', 1, NOW(), NOW()),
(2, 'GYS002', '浙江海鲜水产批发中心', '李总', '13800001002', 'li@haixian.com', '杭州市萧山区水产市场', '主营海鲜水产', 1, NOW(), NOW()),
(3, 'GYS003', '金华火腿食品厂', '王师傅', '13800001003', 'wang@jinhua.com', '金华市婺城区工业园', '主营肉类制品', 1, NOW(), NOW()),
(4, 'GYS004', '温州调味品贸易公司', '陈老板', '13800001004', 'chen@tiaoweipin.com', '温州市鹿城区调味品市场', '主营调味品', 1, NOW(), NOW()),
(5, 'GYS005', '宁波粮油批发商行', '刘经理', '13800001005', 'liu@liangyou.com', '宁波市海曙区粮油市场', '主营粮油干货', 1, NOW(), NOW())
ON DUPLICATE KEY UPDATE supplier_name=VALUES(supplier_name);

-- 2. 供应商分类数据
-- 供应商1: 蔬菜类
INSERT INTO supplier_categories (id, supplier_id, name, sort, status, created_at, updated_at) VALUES
(1, 1, '叶菜类', 1, 1, NOW(), NOW()),
(2, 1, '根茎类', 2, 1, NOW(), NOW()),
(3, 1, '瓜果类', 3, 1, NOW(), NOW()),
(4, 1, '菌菇类', 4, 1, NOW(), NOW());

-- 供应商2: 海鲜类
INSERT INTO supplier_categories (id, supplier_id, name, sort, status, created_at, updated_at) VALUES
(5, 2, '鱼类', 1, 1, NOW(), NOW()),
(6, 2, '虾蟹类', 2, 1, NOW(), NOW()),
(7, 2, '贝类', 3, 1, NOW(), NOW());

-- 供应商3: 肉类
INSERT INTO supplier_categories (id, supplier_id, name, sort, status, created_at, updated_at) VALUES
(8, 3, '猪肉', 1, 1, NOW(), NOW()),
(9, 3, '牛肉', 2, 1, NOW(), NOW()),
(10, 3, '禽类', 3, 1, NOW(), NOW());

-- 供应商4: 调味品
INSERT INTO supplier_categories (id, supplier_id, name, sort, status, created_at, updated_at) VALUES
(11, 4, '酱油醋', 1, 1, NOW(), NOW()),
(12, 4, '香料', 2, 1, NOW(), NOW()),
(13, 4, '酱料', 3, 1, NOW(), NOW());

-- 供应商5: 粮油
INSERT INTO supplier_categories (id, supplier_id, name, sort, status, created_at, updated_at) VALUES
(14, 5, '大米', 1, 1, NOW(), NOW()),
(15, 5, '食用油', 2, 1, NOW(), NOW()),
(16, 5, '干货', 3, 1, NOW(), NOW());

-- 3. 供应商商品数据
-- 叶菜类商品
INSERT INTO supplier_products (id, supplier_id, category_id, name, unit, price, spec, remark, status, created_at, updated_at) VALUES
(1, 1, 1, '大白菜', '斤', 2.50, '新鲜', '当季蔬菜', 1, NOW(), NOW()),
(2, 1, 1, '小白菜', '斤', 3.00, '嫩叶', '', 1, NOW(), NOW()),
(3, 1, 1, '菠菜', '斤', 4.50, '新鲜', '', 1, NOW(), NOW()),
(4, 1, 1, '生菜', '斤', 5.00, '罗马生菜', '', 1, NOW(), NOW()),
(5, 1, 1, '油麦菜', '斤', 4.00, '新鲜', '', 1, NOW(), NOW()),
(6, 1, 1, '空心菜', '斤', 3.50, '新鲜', '', 1, NOW(), NOW());

-- 根茎类商品
INSERT INTO supplier_products (id, supplier_id, category_id, name, unit, price, spec, remark, status, created_at, updated_at) VALUES
(7, 1, 2, '土豆', '斤', 2.00, '黄心', '', 1, NOW(), NOW()),
(8, 1, 2, '红薯', '斤', 2.50, '红心', '', 1, NOW(), NOW()),
(9, 1, 2, '胡萝卜', '斤', 3.00, '新鲜', '', 1, NOW(), NOW()),
(10, 1, 2, '白萝卜', '斤', 1.50, '新鲜', '', 1, NOW(), NOW()),
(11, 1, 2, '莲藕', '斤', 6.00, '九孔', '', 1, NOW(), NOW()),
(12, 1, 2, '山药', '斤', 8.00, '铁棍山药', '', 1, NOW(), NOW());

-- 瓜果类商品
INSERT INTO supplier_products (id, supplier_id, category_id, name, unit, price, spec, remark, status, created_at, updated_at) VALUES
(13, 1, 3, '黄瓜', '斤', 3.00, '新鲜', '', 1, NOW(), NOW()),
(14, 1, 3, '西红柿', '斤', 4.00, '普罗旺斯', '', 1, NOW(), NOW()),
(15, 1, 3, '茄子', '斤', 3.50, '长茄', '', 1, NOW(), NOW()),
(16, 1, 3, '青椒', '斤', 4.00, '新鲜', '', 1, NOW(), NOW()),
(17, 1, 3, '冬瓜', '斤', 1.50, '新鲜', '', 1, NOW(), NOW()),
(18, 1, 3, '南瓜', '斤', 2.00, '贝贝南瓜', '', 1, NOW(), NOW());

-- 菌菇类商品
INSERT INTO supplier_products (id, supplier_id, category_id, name, unit, price, spec, remark, status, created_at, updated_at) VALUES
(19, 1, 4, '香菇', '斤', 12.00, '干香菇', '', 1, NOW(), NOW()),
(20, 1, 4, '平菇', '斤', 6.00, '新鲜', '', 1, NOW(), NOW()),
(21, 1, 4, '金针菇', '包', 4.00, '150g/包', '', 1, NOW(), NOW()),
(22, 1, 4, '杏鲍菇', '斤', 8.00, '新鲜', '', 1, NOW(), NOW());

-- 鱼类商品
INSERT INTO supplier_products (id, supplier_id, category_id, name, unit, price, spec, remark, status, created_at, updated_at) VALUES
(23, 2, 5, '草鱼', '斤', 12.00, '活鱼', '', 1, NOW(), NOW()),
(24, 2, 5, '鲈鱼', '斤', 25.00, '活鱼', '', 1, NOW(), NOW()),
(25, 2, 5, '黄鱼', '斤', 35.00, '野生', '', 1, NOW(), NOW()),
(26, 2, 5, '带鱼', '斤', 18.00, '新鲜', '', 1, NOW(), NOW()),
(27, 2, 5, '三文鱼', '斤', 45.00, '挪威进口', '', 1, NOW(), NOW());

-- 虾蟹类商品
INSERT INTO supplier_products (id, supplier_id, category_id, name, unit, price, spec, remark, status, created_at, updated_at) VALUES
(28, 2, 6, '基围虾', '斤', 38.00, '活虾', '', 1, NOW(), NOW()),
(29, 2, 6, '大闸蟹', '只', 35.00, '3两母蟹', '', 1, NOW(), NOW()),
(30, 2, 6, '皮皮虾', '斤', 45.00, '活虾', '', 1, NOW(), NOW()),
(31, 2, 6, '小龙虾', '斤', 28.00, '4-6钱', '', 1, NOW(), NOW());

-- 贝类商品
INSERT INTO supplier_products (id, supplier_id, category_id, name, unit, price, spec, remark, status, created_at, updated_at) VALUES
(32, 2, 7, '花蛤', '斤', 8.00, '新鲜', '', 1, NOW(), NOW()),
(33, 2, 7, '蛏子', '斤', 15.00, '新鲜', '', 1, NOW(), NOW()),
(34, 2, 7, '扇贝', '斤', 12.00, '新鲜', '', 1, NOW(), NOW()),
(35, 2, 7, '生蚝', '个', 5.00, '大号', '', 1, NOW(), NOW());

-- 猪肉商品
INSERT INTO supplier_products (id, supplier_id, category_id, name, unit, price, spec, remark, status, created_at, updated_at) VALUES
(36, 3, 8, '五花肉', '斤', 18.00, '带皮', '', 1, NOW(), NOW()),
(37, 3, 8, '里脊肉', '斤', 22.00, '精选', '', 1, NOW(), NOW()),
(38, 3, 8, '排骨', '斤', 28.00, '肋排', '', 1, NOW(), NOW()),
(39, 3, 8, '猪蹄', '斤', 20.00, '前蹄', '', 1, NOW(), NOW());

-- 牛肉商品
INSERT INTO supplier_products (id, supplier_id, category_id, name, unit, price, spec, remark, status, created_at, updated_at) VALUES
(40, 3, 9, '牛腩', '斤', 45.00, '新鲜', '', 1, NOW(), NOW()),
(41, 3, 9, '牛腱子', '斤', 55.00, '新鲜', '', 1, NOW(), NOW()),
(42, 3, 9, '肥牛卷', '斤', 48.00, '火锅用', '', 1, NOW(), NOW());

-- 禽类商品
INSERT INTO supplier_products (id, supplier_id, category_id, name, unit, price, spec, remark, status, created_at, updated_at) VALUES
(43, 3, 10, '三黄鸡', '只', 35.00, '约2斤', '', 1, NOW(), NOW()),
(44, 3, 10, '鸡翅', '斤', 18.00, '新鲜', '', 1, NOW(), NOW()),
(45, 3, 10, '鸭腿', '斤', 15.00, '新鲜', '', 1, NOW(), NOW());

-- 调味品商品
INSERT INTO supplier_products (id, supplier_id, category_id, name, unit, price, spec, remark, status, created_at, updated_at) VALUES
(46, 4, 11, '生抽', '瓶', 12.00, '500ml', '海天', 1, NOW(), NOW()),
(47, 4, 11, '老抽', '瓶', 10.00, '500ml', '海天', 1, NOW(), NOW()),
(48, 4, 11, '陈醋', '瓶', 8.00, '500ml', '山西老陈醋', 1, NOW(), NOW()),
(49, 4, 12, '八角', '斤', 35.00, '干货', '', 1, NOW(), NOW()),
(50, 4, 12, '花椒', '斤', 45.00, '四川汉源', '', 1, NOW(), NOW()),
(51, 4, 12, '桂皮', '斤', 25.00, '干货', '', 1, NOW(), NOW()),
(52, 4, 13, '豆瓣酱', '瓶', 15.00, '500g', '郫县豆瓣', 1, NOW(), NOW()),
(53, 4, 13, '甜面酱', '瓶', 10.00, '300g', '', 1, NOW(), NOW());

-- 粮油商品
INSERT INTO supplier_products (id, supplier_id, category_id, name, unit, price, spec, remark, status, created_at, updated_at) VALUES
(54, 5, 14, '东北大米', '袋', 65.00, '10kg', '五常大米', 1, NOW(), NOW()),
(55, 5, 14, '泰国香米', '袋', 85.00, '10kg', '进口', 1, NOW(), NOW()),
(56, 5, 15, '花生油', '桶', 120.00, '5L', '鲁花', 1, NOW(), NOW()),
(57, 5, 15, '菜籽油', '桶', 80.00, '5L', '金龙鱼', 1, NOW(), NOW()),
(58, 5, 16, '干木耳', '斤', 45.00, '东北', '', 1, NOW(), NOW()),
(59, 5, 16, '干腐竹', '斤', 18.00, '优质', '', 1, NOW(), NOW()),
(60, 5, 16, '粉丝', '包', 8.00, '500g', '龙口粉丝', 1, NOW(), NOW());

-- ============================================
-- Mock数据初始化完成
-- 供应商: 5家
-- 分类: 16个
-- 商品: 60个
-- ============================================
