import 'package:flutter/foundation.dart';
import 'role_api.dart';
import 'role_models.dart';

class RoleProvider extends ChangeNotifier {
  final RoleApi _api;
  RoleProvider(this._api);

  List<RoleItem> _list = [];
  bool _loading = false;
  String? _error;

  List<RoleItem> get list => _list;
  bool get loading => _loading;
  String? get error => _error;

  Future<void> load({String? keyword}) async {
    if (_loading) return;
    _loading = true;
    _error = null;
    notifyListeners();
    try {
      _list = await _api.getRoles(keyword: keyword);
    } catch (e) {
      _error = e.toString();
    } finally {
      _loading = false;
      notifyListeners();
    }
  }

  Future<bool> create(CreateRoleRequest req) async {
    try {
      await _api.createRole(req);
      await load();
      return true;
    } catch (e) {
      _error = e.toString();
      notifyListeners();
      return false;
    }
  }

  Future<bool> update(int id, UpdateRoleRequest req) async {
    try {
      await _api.updateRole(id, req);
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
      await _api.deleteRole(id);
      await load();
      return true;
    } catch (e) {
      _error = e.toString();
      notifyListeners();
      return false;
    }
  }
}
