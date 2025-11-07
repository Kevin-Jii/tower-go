import 'package:flutter/foundation.dart';
import 'dingtalk_api.dart';
import 'models.dart';

class DingTalkProvider extends ChangeNotifier {
  final DingTalkApi _api;
  DingTalkProvider(this._api);

  List<DingTalkRobot> _robots = [];
  bool _loading = false;
  String? _error;

  List<DingTalkRobot> get robots => _robots;
  bool get loading => _loading;
  String? get error => _error;

  Future<void> loadRobots({int? storeId}) async {
    print('🔄 开始加载钉钉机器人...');
    _loading = true;
    _error = null;
    notifyListeners();

    final result = await _api.list(storeId: storeId);
    result.when(
      success: (list) {
        print('✅ 加载成功,机器人数量: ${list.length}');
        _robots = list;
        _error = null;
      },
      error: (msg, code) {
        print('❌ 加载失败: $msg (code: $code)');
        _error = msg;
      },
    );

    _loading = false;
    print(
        '📊 状态更新: loading=$_loading, robots=${_robots.length}, error=$_error');
    notifyListeners();
  }

  Future<bool> createRobot(CreateDingTalkRobotRequest payload) async {
    final result = await _api.create(payload);
    bool ok = false;

    result.when(
      success: (item) {
        _error = null;
        ok = true;
        loadRobots();
      },
      error: (msg, code) {
        _error = msg;
      },
    );

    if (!ok) notifyListeners();
    return ok;
  }

  Future<bool> updateRobot(
      DingTalkRobot robot, UpdateDingTalkRobotRequest payload) async {
    final result = await _api.update(robot.id, payload);
    bool ok = false;

    result.when(
      success: (item) {
        _error = null;
        ok = true;
        loadRobots();
      },
      error: (msg, code) {
        _error = msg;
      },
    );

    if (!ok) notifyListeners();
    return ok;
  }

  Future<bool> deleteRobot(DingTalkRobot robot) async {
    final result = await _api.delete(robot.id);
    bool ok = false;

    result.when(
      success: (_) {
        _robots.removeWhere((r) => r.id == robot.id);
        _error = null;
        ok = true;
      },
      error: (msg, code) {
        _error = msg;
      },
    );

    notifyListeners();
    return ok;
  }

  Future<bool> testPush(DingTalkRobot robot) async {
    final result = await _api.testRobot(robot);
    bool ok = false;

    result.when(
      success: (_) {
        _error = null;
        ok = true;
      },
      error: (msg, code) {
        _error = msg;
      },
    );

    if (!ok) notifyListeners();
    return ok;
  }
}
