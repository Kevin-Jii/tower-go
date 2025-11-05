import 'package:flutter/foundation.dart';
import '../../../core/constants/error_texts.dart';
import '../../store/store_api.dart';
import '../../store/models.dart';

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
}
