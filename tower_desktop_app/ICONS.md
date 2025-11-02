# ICONS 图标对照表

本文件用于维护后端返回的 `icon` 字段值与前端 `tdesign_flutter` (TDIcons) 图标之间的映射。当前映射逻辑集中在：

`lib/core/icons/td_icon_mapper.dart`

---
## 1. 已支持映射列表

| 后端 icon 值 | 语义说明 | TDesign 图标 (TDIcons.*) | 兼容别名 | 备注 |
|--------------|----------|--------------------------|----------|------|
| setting | 系统设置 | `TDIcons.setting` | settings | — |
| user | 用户 | `TDIcons.user` |  | — |
| team | 用户组/团队 | `TDIcons.usergroup` | users | 组织/角色集合 |
| menu | 菜单 | `TDIcons.menu` |  | — |
| shop | 门店 | `TDIcons.shop` | store | 门店管理 |
| unordered-list | 列表 | `TDIcons.list` | list | 列表数据入口 |
| coffee | 餐饮/菜品 | `TDIcons.tea` | dish | 与菜品管理语义接近 |
| file-text | 文档/记录 | `TDIcons.file` | document | 可换成 file_copy 视需求 |
| bar-chart | 统计图 | `TDIcons.chart_bar` | chart | 数据统计、报表 |
| role | 角色/权限 | `TDIcons.secured` |  | 权限/分配相关 |
| (空 / 未匹配) | 默认回退 | `TDIcons.app` |  | fallback |

> 注意：具体哪些 TDIcons 存在可能会因 `tdesign_flutter` 版本变化；升级后若图标不存在请及时调整。

---
## 2. 代码入口说明

`TdIconMapper.build(String? name, { double size = 18, Color? color })`

- 负责：
  - 统一 icon 字符串 → 图标组件
  - 未匹配时使用 fallback (`TDIcons.app`)
  - 可在调用处传入统一的 `size` / `color`
- 主要被以下位置调用：
  - 菜单树：`lib/features/menu/menu_tree.dart`
  - 占位页面：`lib/core/routes/route_manager.dart`

---
## 3. 扩展映射步骤

1. 后端新增一个图标标识（例如：`order`）
2. 打开 `td_icon_mapper.dart`，在 `_iconData` 的 `switch` 中添加：
   ```dart
   case 'order':
     return TDIcons.order; // 选择最贴切的图标
   ```
3. 如果存在多个别名，直接多写几个 case：
   ```dart
   case 'order':
   case 'orders':
     return TDIcons.order;
   ```
4. 保存后热重载即可生效。

---
## 4. 建议的未来扩展分类（占位）
| 分类 | 建议后端值 | 可能图标 |
|------|------------|----------|
| 订单 | order / orders | `TDIcons.order` / `TDIcons.list` |
| 员工 | staff | `TDIcons.user` / `TDIcons.usergroup` |
| 通知 | notice / bell | `TDIcons.notification` |
| 报表 | report | `TDIcons.chart_bar` / `TDIcons.data` |
| 日志 | log | `TDIcons.time` / `TDIcons.file` |
| 配置 | config | `TDIcons.setting` |
| 分类 | category | `TDIcons.app` / `TDIcons.menu` |
| 退出 | logout | （若无可用，单独用 Material Icons） |

---
## 5. 可选增强（尚未实现）

### 5.1 智能 Fallback
根据关键字自动猜测：
```dart
if (name.contains('chart')) return TDIcons.chart_bar;
if (name.contains('user')) return TDIcons.user;
```

### 5.2 开发期未匹配警告
```dart
assert(() {
  if (icon == null) debugPrint('[ICON] 未映射图标: $name');
  return true;
}());
```

### 5.3 动态主题适配
让颜色自动继承：
```dart
final effectiveColor = color ?? Theme.of(context).iconTheme.color;
```

### 5.4 图标调试面板
创建一个页面循环展示所有后端已用/可能用的 icon，方便校对与 UI 评审。

---
## 6. 维护规范

| 规则 | 说明 |
|------|------|
| 命名统一 | 后端统一使用中横线或小写英文（如 `file-text` / `bar-chart`） |
| 多别名收敛 | 新旧命名都兼容时，最终在文档中保留一个“主命名” |
| 避免过度粒度 | 不要为极其相近的含义拆出多余图标（如 `user-info`, `user-detail`） |
| 升级评估 | 升级 `tdesign_flutter` 前 diff 一次 TDIcons 变化 |
| 回退策略 | 所有未匹配必须显示 fallback，不能抛异常 |

---
## 7. 快速自检
- [ ] 是否出现了新的后端 icon 值？
- [ ] 是否已经在 `_iconData` 添加映射？
- [ ] 是否有多余的过时 icon？
- [ ] Fallback 是否仍工作正常？

---
## 8. 变更记录

| 日期 | 变更人 | 内容 |
|------|--------|------|
| 2025-11-02 | 初始 | 建立基础映射与文档 |

---
如需我自动生成调试页面或加入智能 fallback，请直接说明。