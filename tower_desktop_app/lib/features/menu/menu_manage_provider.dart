import 'package:flutter/foundation.dart';
import 'menu_api.dart';
import 'models.dart';

/// 菜单管理 Provider (与运行时用户侧边栏菜单分离)
/// 修正：不再依赖后端返回的 children 嵌套，统一根据 parentId 重建树，保证页面与按钮层级可展示。
class MenuManageProvider extends ChangeNotifier {
  final MenuApi _api;
  MenuManageProvider(this._api);

  List<MenuItem> _all = []; // 原始列表（可能是扁平）
  bool _loading = false;
  String? _error;
  final Set<int> _expanded = {}; // 展开中的节点ID
  final Map<int, List<MenuItem>> _byParent = {}; // parentId -> children

  List<MenuItem> get tree => _all; // 暴露原始数据（如编辑选择父级）
  bool get loading => _loading;
  String? get error => _error;
  bool isExpanded(int id) => _expanded.contains(id);

  void toggleExpand(int id) {
    if (_expanded.contains(id)) {
      _expanded.remove(id);
    } else {
      _expanded.add(id);
    }
    notifyListeners();
  }

  Future<void> load() async {
    if (_loading) return;
    _loading = true;
    _error = null;
    notifyListeners();
    try {
      _all = await _api.getMenuTree();
      // 重建 parent -> children 映射
      _byParent.clear();
      for (final m in _all) {
        _byParent.putIfAbsent(m.parentId, () => []).add(m);
      }
      // 默认展开：根层 + 所有有子节点的非按钮（目录/页面）
      _expanded.clear();
      final roots = _byParent[0] ?? [];
      for (final r in roots) {
        final children = _byParent[r.id] ?? [];
        if (children.isNotEmpty) _expanded.add(r.id);
      }
      for (final m in _all) {
        final children = _byParent[m.id] ?? [];
        if (children.isNotEmpty && (m.type == null || m.type != 3)) {
          _expanded.add(m.id);
        }
      }
    } catch (e) {
      _error = e.toString();
    } finally {
      _loading = false;
      notifyListeners();
    }
  }

  /// 树扁平化用于管理表格展示
  List<FlatMenuNode> flattened() {
    final result = <FlatMenuNode>[];
    void walk(List<MenuItem> nodes, int depth) {
      for (final n in nodes) {
        final children = _byParent[n.id] ?? [];
        result.add(FlatMenuNode(n, depth, children.isNotEmpty));
        if (children.isNotEmpty && _expanded.contains(n.id)) {
          walk(children, depth + 1);
        }
      }
    }

    final roots = _byParent[0] ?? [];
    walk(roots, 0);
    return result;
  }

  Future<bool> create(CreateMenuRequest req) async {
    try {
      await _api.createMenu(req);
      await load();
      return true;
    } catch (e) {
      _error = e.toString();
      notifyListeners();
      return false;
    }
  }

  Future<bool> update(int id, UpdateMenuRequest req) async {
    try {
      await _api.updateMenu(id, req);
      await load();
      return true;
    } catch (e) {
      _error = e.toString();
      notifyListeners();
      return false;
    }
  }

  Future<bool> remove(int id) async {
    try {
      await _api.deleteMenu(id);
      await load();
      return true;
    } catch (e) {
      _error = e.toString();
      notifyListeners();
      return false;
    }
  }
}

class FlatMenuNode {
  final MenuItem item;
  final int depth;
  final bool hasChildren;
  FlatMenuNode(this.item, this.depth, this.hasChildren);
}
