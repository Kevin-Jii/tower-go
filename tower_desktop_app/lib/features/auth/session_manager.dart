import '../../core/network/api_client.dart';
import 'models.dart';

class SessionManager {
  static final SessionManager _instance = SessionManager._internal();
  factory SessionManager() => _instance;
  SessionManager._internal();

  String? _token;
  UserInfo? _user;
  List<String> _permissions = [];

  String? get token => _token;
  UserInfo? get user => _user;
  List<String> get permissions => _permissions;

  bool get isLoggedIn => _token != null && _user != null;

  void updateSession({required String token, required UserInfo user, required List<String> permissions}) {
    _token = token;
    _user = user;
    _permissions = permissions;
    ApiClient().setToken(token);
  }

  void clear() {
    _token = null;
    _user = null;
    _permissions = [];
    ApiClient().setToken(null);
  }

  bool hasPermission(String code) => _permissions.contains(code);
}
