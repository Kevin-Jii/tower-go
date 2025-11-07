import 'package:flutter/foundation.dart';
import '../../../core/constants/error_texts.dart';
import '../../store/store_api.dart';
import '../../store/models.dart';
import '../../auth/session_manager.dart';

/// 管理门店列表与当前选中门店
class StoreSelectorProvider extends ChangeNotifier {
  final StoreApi _storeApi;
  StoreSelectorProvider(this._storeApi);

  List<Store> _stores = [];
  bool _loading = false;
  String? _error;
  int? _selectedStoreId;

  List<Store> get stores => _stores;
  bool get loading => _loading;
  String? get error => _error;
  int? get selectedStoreId => _selectedStoreId;

  /// 初始化门店:如果是 admin 则加载所有门店,否则使用用户的 storeId
  Future<void> initializeStore() async {
    final session = SessionManager();
    final userInfo = session.userInfo;

    if (userInfo == null) return;

    // 检查是否是 admin 角色
    final isAdmin = userInfo.role?.code == 'admin';

    if (isAdmin) {
      // Admin 用户:加载所有门店列表
      await loadStores();
    } else {
      // 非 admin 用户:直接使用用户的 storeId
      if (userInfo.storeId > 0) {
        _selectedStoreId = userInfo.storeId;
        notifyListeners();
      }
    }
  }

  Future<void> loadStores() async {
    _loading = true;
    _error = null;
    notifyListeners();
    try {
      _stores = await _storeApi.listStores();
      if (_selectedStoreId == null && _stores.isNotEmpty) {
        _selectedStoreId = _stores.first.id;
      }
    } catch (e) {
      _error = ErrorTexts.loadStores + ': ' + e.toString();
    } finally {
      _loading = false;
      notifyListeners();
    }
  }

  void selectStore(int id) {
    if (_selectedStoreId == id) return;
    _selectedStoreId = id;
    notifyListeners();
  }

  /// 检查当前用户是否是 admin
  bool get isAdmin {
    final session = SessionManager();
    return session.userInfo?.role?.code == 'admin';
  }
}
