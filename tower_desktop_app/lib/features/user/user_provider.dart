import 'package:flutter/foundation.dart';
import 'user_api.dart';
import 'models.dart';
import '../store/store_api.dart';

class UserProvider extends ChangeNotifier {
  final UserApi _api;
  final StoreApi _storeApi;

  UserProvider(this._api, this._storeApi);

  List<User> _users = [];
  int _total = 0;
  int _page = 1;
  int _pageSize = 10;
  bool _loading = false;
  String? _error;
  String? _keyword;

  // 门店映射表 id -> name
  Map<int, String> _storeMap = {};

  List<User> get users => _users;
  int get total => _total;
  int get page => _page;
  int get pageSize => _pageSize;
  bool get loading => _loading;
  String? get error => _error;
  String? get keyword => _keyword;

  Future<void> loadUsers({
    int? page,
    int? pageSize,
    String? keyword,
  }) async {
    if (_loading) return;
    _loading = true;
    _error = null;
    notifyListeners();

    try {
      _page = page ?? _page;
      _pageSize = pageSize ?? _pageSize;
      _keyword = keyword;

      debugPrint(
          'Loading users: page=$_page, pageSize=$_pageSize, keyword=$_keyword');

      // 先加载门店列表以建立映射
      await _loadStoreMap();

      final resp = await _api.getUsers(
        page: _page,
        pageSize: _pageSize,
        keyword: _keyword,
      );

      debugPrint(
          'Received users: ${resp.list.length}, total=${resp.total}, page=${resp.page}, pageSize=${resp.pageSize}');

      // 填充 storeName
      _users = resp.list.map((user) {
        String? storeName;
        if (user.storeId != null && _storeMap.containsKey(user.storeId)) {
          storeName = _storeMap[user.storeId];
        }
        return user.copyWith(storeName: storeName);
      }).toList();

      _total = resp.total;
      // 不使用后端返回的 page 和 pageSize，保持前端的设置
    } catch (e) {
      _error = e.toString();
      debugPrint('Error loading users: $e');
    } finally {
      _loading = false;
      notifyListeners();
    }
  }

  Future<void> _loadStoreMap() async {
    try {
      final storesResp = await _storeApi.getStores(page: 1, pageSize: 100);
      _storeMap = {
        for (var store in storesResp.list) store.id: store.name,
      };
    } catch (e) {
      // 忽略门店加载失败的错误
      debugPrint('Failed to load stores: $e');
    }
  }

  Future<bool> createUser(CreateUserRequest req) async {
    try {
      await _api.createUser(req);
      await loadUsers(); // 刷新列表
      return true;
    } catch (e) {
      _error = e.toString();
      notifyListeners();
      return false;
    }
  }

  Future<bool> updateUser(int id, UpdateUserRequest req) async {
    try {
      await _api.updateUser(id, req);
      await loadUsers();
      return true;
    } catch (e) {
      _error = e.toString();
      notifyListeners();
      return false;
    }
  }

  Future<bool> deleteUser(int id) async {
    try {
      await _api.deleteUser(id);
      await loadUsers();
      return true;
    } catch (e) {
      _error = e.toString();
      notifyListeners();
      return false;
    }
  }

  Future<bool> resetPassword(int id, String newPassword) async {
    try {
      await _api.resetPassword(id, newPassword);
      return true;
    } catch (e) {
      _error = e.toString();
      notifyListeners();
      return false;
    }
  }

  void clearError() {
    _error = null;
    notifyListeners();
  }
}
