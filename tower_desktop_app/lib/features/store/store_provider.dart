import 'package:flutter/foundation.dart';
import 'store_api.dart';
import 'models.dart';

class StoreProvider with ChangeNotifier {
  final StoreApi _api;

  StoreProvider(this._api);

  List<Store> _stores = [];
  List<Store> get stores => _stores;

  bool _loading = false;
  bool get loading => _loading;

  String? _error;
  String? get error => _error;

  int _page = 1;
  int get page => _page;

  int _total = 0;
  int get total => _total;

  int _pageSize = 10;
  int get pageSize => _pageSize;

  /// 加载门店列表
  Future<void> loadStores({
    int? page,
    int? pageSize,
    String? keyword,
  }) async {
    if (page != null) _page = page;
    if (pageSize != null) _pageSize = pageSize;

    _loading = true;
    _error = null;
    notifyListeners();

    try {
      final response = await _api.getStores(
        page: _page,
        pageSize: _pageSize,
        keyword: keyword,
      );
      _stores = response.list;
      _total = response.total;
      _page = response.page;
      _pageSize = response.pageSize;
      _error = null;
    } catch (e) {
      _error = e.toString();
      _stores = [];
    } finally {
      _loading = false;
      notifyListeners();
    }
  }

  /// 创建门店
  Future<bool> createStore(CreateStoreRequest request) async {
    try {
      await _api.createStore(request);
      await loadStores(page: 1); // 刷新列表
      return true;
    } catch (e) {
      _error = e.toString();
      notifyListeners();
      return false;
    }
  }

  /// 更新门店（静默更新，不显示全局loading）
  Future<bool> updateStore(int id, UpdateStoreRequest request) async {
    try {
      await _api.updateStore(id, request);
      // 静默刷新：不显示loading状态
      await _refreshStoresSilently();
      return true;
    } catch (e) {
      _error = e.toString();
      notifyListeners();
      return false;
    }
  }

  /// 静默刷新门店列表（不显示loading动画）
  Future<void> _refreshStoresSilently() async {
    try {
      final response = await _api.getStores(
        page: _page,
        pageSize: _pageSize,
      );
      _stores = response.list;
      _total = response.total;
      _page = response.page;
      _pageSize = response.pageSize;
      _error = null;
      notifyListeners(); // 只通知数据更新，不改变loading状态
    } catch (e) {
      _error = e.toString();
      notifyListeners();
    }
  }

  /// 删除门店
  Future<bool> deleteStore(int id) async {
    try {
      await _api.deleteStore(id);
      await loadStores(); // 刷新当前页
      return true;
    } catch (e) {
      _error = e.toString();
      notifyListeners();
      return false;
    }
  }
}
