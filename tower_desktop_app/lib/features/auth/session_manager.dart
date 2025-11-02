import 'dart:convert';
import 'package:shared_preferences/shared_preferences.dart';
import '../../core/network/api_client.dart';
import 'models.dart';

class SessionManager {
  static final SessionManager _instance = SessionManager._internal();
  factory SessionManager() => _instance;
  SessionManager._internal();

  String? _token;
  UserInfo? _userInfo;
  List<String> _permissions = [];
  int? _expiresIn; // 秒
  int? _issuedAt; // epoch seconds

  String? get token => _token;
  UserInfo? get userInfo => _userInfo;
  List<String> get permissions => _permissions;
  int? get expiresIn => _expiresIn;
  int? get issuedAt => _issuedAt;

  bool get isLoggedIn => _token != null && _userInfo != null;
  bool get isExpired {
    if (_token == null || _expiresIn == null || _issuedAt == null) return true;
    final now = DateTime.now().millisecondsSinceEpoch ~/ 1000;
    return now >= _issuedAt! + _expiresIn!;
  }

  static const _kToken = 'session.token';
  static const _kUser = 'session.user';
  static const _kPerms = 'session.permissions';
  static const _kExpiresIn = 'session.expires_in';
  static const _kIssuedAt = 'session.issued_at';

  Future<void> init() async {
    final prefs = await SharedPreferences.getInstance();
    final savedToken = prefs.getString(_kToken);
    final savedUser = prefs.getString(_kUser);
    final savedPerms = prefs.getStringList(_kPerms) ?? [];
    final savedExpiresIn = prefs.getInt(_kExpiresIn);
    final savedIssuedAt = prefs.getInt(_kIssuedAt);
    if (savedToken != null && savedUser != null) {
      try {
        final userMap = jsonDecode(savedUser) as Map<String, dynamic>;
        final ui = UserInfo.fromJson(userMap);
        _token = savedToken;
        _userInfo = ui;
        _permissions = savedPerms;
        _expiresIn = savedExpiresIn;
        _issuedAt = savedIssuedAt;
        ApiClient().setToken(savedToken);
        // 过期检查
        if (isExpired) {
          await clear(persist: true);
        }
      } catch (_) {
        await clear(persist: true);
      }
    }
  }

  void updateSession(
      {required String token,
      required UserInfo userInfo,
      required List<String> permissions,
      int? expiresIn}) {
    _token = token;
    _userInfo = userInfo;
    _permissions = permissions;
    _expiresIn = expiresIn;
    _issuedAt = DateTime.now().millisecondsSinceEpoch ~/ 1000;
    ApiClient().setToken(token);
    _persist();
  }

  Future<void> _persist() async {
    final prefs = await SharedPreferences.getInstance();
    if (_token == null || _userInfo == null) return;
    await prefs.setString(_kToken, _token!);
    await prefs.setString(_kUser, jsonEncode(_userInfo!.toJson()));
    await prefs.setStringList(_kPerms, _permissions);
    if (_expiresIn != null) await prefs.setInt(_kExpiresIn, _expiresIn!);
    if (_issuedAt != null) await prefs.setInt(_kIssuedAt, _issuedAt!);
  }

  Future<void> clear({bool persist = false}) async {
    _token = null;
    _userInfo = null;
    _permissions = [];
    _expiresIn = null;
    _issuedAt = null;
    ApiClient().setToken(null);
    if (persist) {
      final prefs = await SharedPreferences.getInstance();
      await prefs.remove(_kToken);
      await prefs.remove(_kUser);
      await prefs.remove(_kPerms);
      await prefs.remove(_kExpiresIn);
      await prefs.remove(_kIssuedAt);
    }
  }

  bool hasPermission(String code) => _permissions.contains(code);
}
