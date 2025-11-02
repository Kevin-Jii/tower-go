import '../../core/network/api_client.dart';
import 'models.dart';

class SessionManager {
  static final SessionManager _instance = SessionManager._internal();
  factory SessionManager() => _instance;
  SessionManager._internal();

  String? _token;
  UserInfo? _userInfo;
  List<String> _permissions = [];

  String? get token => _token;
  UserInfo? get userInfo => _userInfo;
  List<String> get permissions => _permissions;

  bool get isLoggedIn => _token != null && _userInfo != null;

  void updateSession(
      {required String token,
      required UserInfo userInfo,
      required List<String> permissions}) {
    _token = token;
    _userInfo = userInfo;
    _permissions = permissions;
    ApiClient().setToken(token);
  }

  void clear() {
    _token = null;
    _userInfo = null;
    _permissions = [];
    ApiClient().setToken(null);
  }

  bool hasPermission(String code) => _permissions.contains(code);
}
