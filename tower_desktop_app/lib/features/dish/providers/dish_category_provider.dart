import 'package:flutter/foundation.dart';
import '../dish_category_api.dart';
import 'store_selector_provider.dart';
import '../../../core/constants/error_texts.dart';

/// 分类列表 Provider：监听门店选择变化自动刷新
class DishCategoryProvider extends ChangeNotifier {
  final DishCategoryApi _api;
  final StoreSelectorProvider _storeSelector;
  DishCategoryProvider(this._api, this._storeSelector) {
    _storeSelector.addListener(_onStoreChanged);
  }

  List<DishCategory> _categories = [];
  bool _loading = false;
  String? _error;
  DishCategory? _selected;

  List<DishCategory> get categories => _categories;
  bool get loading => _loading;
  String? get error => _error;
  DishCategory? get selectedCategory => _selected;

  void _onStoreChanged() {
    // 门店改变时自动加载分类
    loadCategories();
  }

  Future<void> loadCategories() async {
    final storeId = _storeSelector.selectedStoreId;
    if (storeId == null) {
      _categories = [];
      _selected = null;
      notifyListeners();
      return;
    }
    _loading = true;
    _error = null;
    notifyListeners();
    final result = await _api.list(storeId);
    result.when(
      success: (list) {
        _categories = list;
        if (_categories.isNotEmpty) {
          // 保持之前选中如果存在于新列表
          if (_selected != null) {
            final keep = _categories.firstWhere(
              (c) => c.id == _selected!.id,
              orElse: () => _categories.first,
            );
            _selected = keep;
          } else {
            _selected = null; // 不自动选中，交给 UI
          }
        } else {
          _selected = null;
        }
      },
      error: (msg, code) {
        _error = msg;
      },
    );
    _loading = false;
    notifyListeners();
  }

  void selectCategory(DishCategory c) {
    _selected = c;
    _error = null; // Clear any previous errors when selecting a category
    notifyListeners();
  }

  Future<bool> createCategory(String name) async {
    final storeId = _storeSelector.selectedStoreId;
    if (storeId == null) return false;
    final result = await _api.create(storeId, name);
    bool ok = false;
    result.when(
      success: (item) {
        _error = null; // Clear error on success
        ok = true;
        // Reload the entire category list to get fresh data from backend
        loadCategories();
      },
      error: (msg, code) {
        _error = msg;
      },
    );
    // Don't call notifyListeners here if error, loadCategories will handle success case
    if (!ok) notifyListeners();
    return ok;
  }

  Future<bool> updateCategory(DishCategory c, String name) async {
    final storeId = _storeSelector.selectedStoreId;
    if (storeId == null) return false;
    final result = await _api.update(storeId, c.id, name);
    bool ok = false;
    result.when(
      success: (item) {
        _error = null;
        ok = true;
        loadCategories();
      },
      error: (msg, code) {
        _error = msg;
      },
    );
    if (!ok) notifyListeners();
    return ok;
  }

  Future<bool> deleteCategory(DishCategory c,
      {required bool Function(DishCategory) hasDishes}) async {
    final storeId = _storeSelector.selectedStoreId;
    if (storeId == null) return false;

    // Check if category has dishes before deletion
    if (hasDishes(c)) {
      _error = ErrorTexts.deleteCategoryHasDishes;
      notifyListeners();
      return false;
    }

    final result = await _api.delete(storeId, c.id);
    bool ok = false;
    result.when(
      success: (_) {
        _categories.removeWhere((e) => e.id == c.id);
        if (_selected?.id == c.id) {
          _selected = null;
          if (_categories.isNotEmpty) {
            _selected = _categories.first;
          }
        }
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

  @override
  void dispose() {
    _storeSelector.removeListener(_onStoreChanged);
    super.dispose();
  }
}
