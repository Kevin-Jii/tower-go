import 'package:flutter/foundation.dart';
import 'user_api.dart';
import 'models.dart';

class UserProvider extends ChangeNotifier {
  final UserApi _api;
  UserProvider(this._api);

  List<User> _users = [];
  int _total = 0;
  int _page = 1;
  int _pageSize = 10;
  bool _loading = false;
  String? _error;
  String? _keyword;

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

      final resp = await _api.getUsers(
        page: _page,
        pageSize: _pageSize,
        keyword: _keyword,
      );

      _users = resp.list;
      _total = resp.total;
    } catch (e) {
      _error = e.toString();
    } finally {
      _loading = false;
      notifyListeners();
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
