import 'package:flutter/foundation.dart';

/// 通用列表 Provider 抽象：处理 loading / error / list 基础状态。
/// 子类只需实现 fetchList 来提供实际数据来源，必要时可扩展分页等。
abstract class ListProvider<T> extends ChangeNotifier {
  List<T> _items = [];
  bool _loading = false;
  String? _error;

  List<T> get items => _items;
  bool get loading => _loading;
  String? get error => _error;

  /// 可选：最近一次查询参数（如 keyword 等）
  Map<String, dynamic>? lastQuery;

  /// 子类实现：基于 query 拉取数据。
  Future<List<T>> fetchList(Map<String, dynamic>? query);

  Future<void> load({Map<String, dynamic>? query}) async {
    if (_loading) return;
    _loading = true;
    _error = null;
    notifyListeners();
    try {
      lastQuery = query;
      _items = await fetchList(query);
    } catch (e) {
      _error = e.toString();
    } finally {
      _loading = false;
      notifyListeners();
    }
  }

  /// 重新加载使用上一次查询条件。
  Future<void> reload() => load(query: lastQuery);

  /// 子类在新增/修改/删除成功后可调用，或直接重新 load。
  void replaceAll(List<T> data) {
    _items = data;
    notifyListeners();
  }

  /// 内部帮助：添加单条并刷新。
  void addItem(T item) {
    _items.add(item);
    notifyListeners();
  }

  void updateItem(bool Function(T) test, T Function(T) mapper) {
    _items = _items.map((e) => test(e) ? mapper(e) : e).toList();
    notifyListeners();
  }

  void removeWhere(bool Function(T) test) {
    _items.removeWhere(test);
    notifyListeners();
  }
}
