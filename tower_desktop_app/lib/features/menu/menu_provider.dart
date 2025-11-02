import 'package:flutter/foundation.dart';
import 'menu_api.dart';
import 'models.dart';
import '../../core/constants/menu_types.dart';
import '../auth/permission_provider.dart';

class MenuProvider extends ChangeNotifier {
  final MenuApi _api;
  MenuProvider(this._api);

  List<MenuItem> _tree = [];
  bool _loading = false;
  String? _error;
  MenuItem? _selected;

  List<MenuItem> get tree => _tree;
  bool get loading => _loading;
  String? get error => _error;
  MenuItem? get selected => _selected;

  Future<void> load({PermissionProvider? permissionProvider}) async {
    if (_loading) return;
    _loading = true;
    _error = null;
    notifyListeners();
    try {
      // 后端已返回树结构，直接使用并递归过滤掉按钮 (type == 3)
      final raw = await _api.getUserMenus();

      final List<String> collectedPerms = [];

      List<MenuItem> filterTree(List<MenuItem> nodes) {
        final result = <MenuItem>[];
        for (final n in nodes) {
          if (n.type == MenuType.button) {
            // 收集按钮权限（目前仅支持单 permission 字段；perms 数组待代码生成更新后启用）
            if (n.permission != null && n.permission!.isNotEmpty) {
              collectedPerms.add(n.permission!);
            }
            continue; // 不放入侧边栏
          }
          result.add(n.copyWith(children: filterTree(n.children)));
        }
        return result;
      }

      _tree = filterTree(raw.where((n) => n.parentId == 0).toList());

      // 设置权限集合
      permissionProvider?.setPermissions(collectedPerms);
      if (_tree.isNotEmpty) {
        // 默认选择第一个“可点击”的菜单（type == 2），否则第一个目录
        MenuItem? firstMenu;
        void findFirst(List<MenuItem> list) {
          for (final m in list) {
            if (m.type == MenuType.page) {
              firstMenu = m;
              return;
            }
            if (m.children.isNotEmpty) findFirst(m.children);
            if (firstMenu != null) return;
          }
        }

        findFirst(_tree);
        _selected = firstMenu ?? _tree.first;
      }
    } catch (e) {
      _error = e.toString();
    } finally {
      _loading = false;
      notifyListeners();
    }
  }

  void select(MenuItem item) {
    _selected = item;
    notifyListeners();
  }
}
