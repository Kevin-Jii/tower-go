import 'package:flutter/foundation.dart';
import '../../core/network/api_client.dart';
import 'report_api.dart';
import 'models.dart';

class ReportProvider with ChangeNotifier {
  final ReportApi _api;

  ReportProvider({ReportApi? api}) : _api = api ?? ReportApi(ApiClient());

  List<MenuReport> _reports = [];
  List<MenuReport> get reports => _reports;

  bool _loading = false;
  bool get loading => _loading;

  String? _error;
  String? get error => _error;

  // 筛选条件
  int? _selectedStoreId;
  int? _selectedDishId;
  String? _startDate;
  String? _endDate;

  int? get selectedStoreId => _selectedStoreId;
  int? get selectedDishId => _selectedDishId;
  String? get startDate => _startDate;
  String? get endDate => _endDate;

  /// 设置筛选条件并加载数据
  Future<void> setFilter({
    int? storeId,
    int? dishId,
    String? startDate,
    String? endDate,
  }) async {
    _selectedStoreId = storeId;
    _selectedDishId = dishId;
    _startDate = startDate;
    _endDate = endDate;
    await loadReports();
  }

  /// 清除筛选条件
  void clearFilter() {
    _selectedStoreId = null;
    _selectedDishId = null;
    _startDate = null;
    _endDate = null;
    notifyListeners();
  }

  /// 加载报菜记录列表
  Future<void> loadReports() async {
    _loading = true;
    _error = null;
    notifyListeners();

    final result = await _api.list(
      storeId: _selectedStoreId,
      dishId: _selectedDishId,
      startDate: _startDate,
      endDate: _endDate,
    );

    result.when(
      success: (data) {
        _reports = data;
        _error = null;
      },
      error: (message, code) {
        _error = message;
        debugPrint('加载报菜记录失败: $message (code: $code)');
      },
    );

    _loading = false;
    notifyListeners();
  }

  /// 创建报菜记录
  Future<bool> createReport(CreateMenuReportRequest payload) async {
    _loading = true;
    _error = null;
    notifyListeners();

    final result = await _api.create(payload);

    bool success = false;
    result.when(
      success: (data) {
        success = true;
        _error = null;
        // 刷新列表
        loadReports();
      },
      error: (message, code) {
        _error = message;
        debugPrint('创建报菜失败: $message (code: $code)');
        _loading = false;
        notifyListeners();
      },
    );

    return success;
  }

  /// 更新报菜记录
  Future<bool> updateReport(int id, UpdateMenuReportRequest payload) async {
    _loading = true;
    _error = null;
    notifyListeners();

    final result = await _api.update(id, payload);

    bool success = false;
    result.when(
      success: (data) {
        success = true;
        _error = null;
        // 刷新列表
        loadReports();
      },
      error: (message, code) {
        _error = message;
        debugPrint('更新报菜失败: $message (code: $code)');
        _loading = false;
        notifyListeners();
      },
    );

    return success;
  }

  /// 删除报菜记录
  Future<bool> deleteReport(int id) async {
    _loading = true;
    _error = null;
    notifyListeners();

    final result = await _api.delete(id);

    bool success = false;
    result.when(
      success: (_) {
        success = true;
        _error = null;
        // 刷新列表
        loadReports();
      },
      error: (message, code) {
        _error = message;
        debugPrint('删除报菜失败: $message (code: $code)');
        _loading = false;
        notifyListeners();
      },
    );

    return success;
  }

  /// 批量删除报菜记录
  Future<bool> batchDeleteReports(List<int> ids) async {
    if (ids.isEmpty) return false;

    _loading = true;
    _error = null;
    notifyListeners();

    final result = await _api.batchDelete(ids);

    bool success = false;
    result.when(
      success: (_) {
        success = true;
        _error = null;
        // 刷新列表
        loadReports();
      },
      error: (message, code) {
        _error = message;
        debugPrint('批量删除失败: $message (code: $code)');
        _loading = false;
        notifyListeners();
      },
    );

    return success;
  }

  /// 清除错误状态
  void clearError() {
    _error = null;
    notifyListeners();
  }
}
