import 'package:flutter/foundation.dart';

/// 负责维护当前用户拥有的按钮权限（后端返回的 permission / perms）
class PermissionProvider extends ChangeNotifier {
  final Set<String> _permissions = {};

  Set<String> get all => _permissions;

  void setPermissions(Iterable<String> perms) {
    _permissions
      ..clear()
      ..addAll(perms.where((e) => e.trim().isNotEmpty));
    notifyListeners();
  }

  bool has(String code) => _permissions.contains(code);
  bool hasAny(Iterable<String> codes) => codes.any(has);
  bool hasAll(Iterable<String> codes) => codes.every(has);
}
