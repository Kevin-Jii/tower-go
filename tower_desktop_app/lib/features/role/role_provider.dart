import '../../core/providers/list_provider.dart';
import 'role_api.dart';
import 'role_models.dart';

class RoleProvider extends ListProvider<RoleItem> {
  final RoleApi _api;
  RoleProvider(this._api);

  List<RoleItem> get list => items;

  // 分页状态（后端暂未提供分页接口时，前端模拟切片）
  int _page = 1;
  int _pageSize = 20; // 可调整默认每页数量
  int _total = 0; // 后端若返回 total 可覆盖
  String? _keyword; // 当前搜索关键字

  int get page => _page;
  int get pageSize => _pageSize;
  int get total => _total;
  String? get keyword => _keyword;

  List<RoleItem> get currentPageItems {
    final start = (_page - 1) * _pageSize;
    if (start >= list.length) return const [];
    final end =
        (start + _pageSize) > list.length ? list.length : (start + _pageSize);
    return list.sublist(start, end);
  }

  void setPage(int newPage) {
    if (newPage < 1) return;
    final totalPages =
        _total == 0 ? 1 : ((_total + _pageSize - 1) / _pageSize).ceil();
    if (newPage > totalPages) return;
    _page = newPage;
    notifyListeners();
  }

  void setPageSize(int newSize) {
    if (newSize <= 0) return;
    _pageSize = newSize;
    _page = 1; // 重置到第一页
    notifyListeners();
  }

  @override
  Future<List<RoleItem>> fetchList(Map<String, dynamic>? query) {
    final keyword = query?['keyword'] as String?;
    return _api.getRoles(keyword: keyword);
  }

  Future<bool> create(CreateRoleRequest req) async {
    try {
      await _api.createRole(req);
      await reload();
      return true;
    } catch (e) {
      return false;
    }
  }

  Future<bool> update(int id, UpdateRoleRequest req) async {
    try {
      await _api.updateRole(id, req);
      await reload();
      return true;
    } catch (e) {
      return false;
    }
  }

  Future<bool> remove(int id) async {
    try {
      await _api.deleteRole(id);
      await reload();
      return true;
    } catch (e) {
      return false;
    }
  }

  // 便捷调用：与旧签名保持一致，但避免与父类冲突，改名为 loadRoles。
  Future<void> loadRoles({String? keyword}) => super.load(query: {
        if (keyword != null && keyword.isNotEmpty) 'keyword': keyword,
      }).then((_) {
        _keyword = keyword?.trim().isEmpty == true ? null : keyword?.trim();
        _total = list.length; // 若后端未来返回 total，可在此替换
        _page = 1; // 搜索后回到第一页
        notifyListeners();
      });
}
