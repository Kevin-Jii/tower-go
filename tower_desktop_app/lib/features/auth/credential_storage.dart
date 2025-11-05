import 'dart:convert';
import 'package:shared_preferences/shared_preferences.dart';

/// 简单的本地凭据存储。
/// 注意：出于安全考虑，生产环境不要明文/可逆方式保存密码；
/// 可替换为本机安全存储（Windows Credential Manager / macOS Keychain / Linux Secret Service）。
class CredentialStorage {
  static const _kAccount = 'login.account'; // 兼容旧版本最后一次账号
  static const _kPassword = 'login.password'; // 兼容旧版本最后一次密码
  static const _kRemember = 'login.remember';
  static const _kAccountsHistory = 'login.accounts.history';
  static const _kAccountPwds = 'login.account.password.map';
  static const _kLastAccount = 'login.last.account';
  static const int _maxAccounts = 5;

  static final CredentialStorage _instance = CredentialStorage._internal();
  factory CredentialStorage() => _instance;
  CredentialStorage._internal();

  Future<void> save({required String account, required String password}) async {
    final prefs = await SharedPreferences.getInstance();
    // 兼容老字段保留
    await prefs.setString(_kAccount, account);
    final encoded = base64Encode(utf8.encode(password));
    await prefs.setString(_kPassword, encoded);
    await prefs.setBool(_kRemember, true);

    // 多账号扩展：更新历史 & 密码映射
    await _addAccountToHistory(account);
    await _savePasswordForAccount(account, encoded);
  }

  Future<(String? account, String? password, bool remember)> read() async {
    final prefs = await SharedPreferences.getInstance();
    // 新逻辑优先 last account
    final last = prefs.getString(_kLastAccount) ?? prefs.getString(_kAccount);
    final remember = prefs.getBool(_kRemember) ?? false;
    String? pwd;
    if (last != null && remember) {
      pwd = await getPasswordFor(last);
    } else {
      final legacy = prefs.getString(_kPassword);
      if (legacy != null) {
        try {
          pwd = utf8.decode(base64Decode(legacy));
        } catch (_) {}
      }
    }
    return (last, pwd, remember);
  }

  Future<List<String>> getAccountsHistory() async {
    final prefs = await SharedPreferences.getInstance();
    return prefs.getStringList(_kAccountsHistory) ?? const [];
  }

  Future<String?> getPasswordFor(String account) async {
    final prefs = await SharedPreferences.getInstance();
    final mapJson = prefs.getString(_kAccountPwds);
    if (mapJson == null) return null;
    try {
      final map = Map<String, dynamic>.from(jsonDecode(mapJson));
      final encoded = map[account];
      if (encoded is String) {
        return utf8.decode(base64Decode(encoded));
      }
    } catch (_) {}
    return null;
  }

  Future<void> selectAccount(String account) async {
    final prefs = await SharedPreferences.getInstance();
    await prefs.setString(_kLastAccount, account);
    await _addAccountToHistory(account); // 选择也刷新顺序
  }

  Future<void> removeAccount(String account) async {
    final prefs = await SharedPreferences.getInstance();
    final history = prefs.getStringList(_kAccountsHistory) ?? [];
    history.remove(account);
    await prefs.setStringList(_kAccountsHistory, history);
    final mapJson = prefs.getString(_kAccountPwds);
    if (mapJson != null) {
      try {
        final map = Map<String, dynamic>.from(jsonDecode(mapJson));
        map.remove(account);
        await prefs.setString(_kAccountPwds, jsonEncode(map));
      } catch (_) {}
    }
    if ((prefs.getString(_kLastAccount)) == account) {
      await prefs.remove(_kLastAccount);
    }
  }

  Future<void> _addAccountToHistory(String account) async {
    final prefs = await SharedPreferences.getInstance();
    final list = prefs.getStringList(_kAccountsHistory) ?? [];
    list.remove(account);
    list.insert(0, account);
    while (list.length > _maxAccounts) {
      list.removeLast();
    }
    await prefs.setStringList(_kAccountsHistory, list);
    await prefs.setString(_kLastAccount, account);
  }

  Future<void> _savePasswordForAccount(
      String account, String encodedPwd) async {
    final prefs = await SharedPreferences.getInstance();
    Map<String, dynamic> map = {};
    final mapJson = prefs.getString(_kAccountPwds);
    if (mapJson != null) {
      try {
        map = Map<String, dynamic>.from(jsonDecode(mapJson));
      } catch (_) {}
    }
    map[account] = encodedPwd;
    await prefs.setString(_kAccountPwds, jsonEncode(map));
  }

  Future<void> forgetPasswordFor(String account) async {
    final prefs = await SharedPreferences.getInstance();
    final mapJson = prefs.getString(_kAccountPwds);
    if (mapJson == null) return;
    try {
      final map = Map<String, dynamic>.from(jsonDecode(mapJson));
      map.remove(account);
      await prefs.setString(_kAccountPwds, jsonEncode(map));
    } catch (_) {}
  }

  Future<void> clear() async {
    final prefs = await SharedPreferences.getInstance();
    await prefs.remove(_kAccount);
    await prefs.remove(_kPassword);
    await prefs.remove(_kRemember);
    await prefs.remove(_kLastAccount);
    // 不清除历史列表与其它账号密码映射，保持其它账号仍可下拉可选
  }
}
