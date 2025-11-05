import 'package:flutter/foundation.dart';
import '../dish_api.dart';
import '../models.dart';
import 'dish_category_provider.dart';
import 'store_selector_provider.dart';

/// 菜品列表 Provider：依赖门店与分类选择
class DishListProvider extends ChangeNotifier {
  final DishApi _api;
  final StoreSelectorProvider _storeSelector;
  final DishCategoryProvider _categoryProvider;
  DishListProvider(this._api, this._storeSelector, this._categoryProvider) {
    _storeSelector.addListener(_triggerReload);
    _categoryProvider.addListener(_triggerReload);
  }

  List<Dish> _dishes = [];
  bool _loading = false;
  String? _error;
  String _keyword = '';

  List<Dish> get dishes => _filtered();
  bool get loading => _loading;
  String? get error => _error;
  String get keyword => _keyword;

  List<Dish> _filtered() {
    if (_keyword.isEmpty) return _dishes;
    final k = _keyword.toLowerCase();
    return _dishes.where((d) => d.name.toLowerCase().contains(k)).toList();
  }

  void setKeyword(String v) {
    _keyword = v;
    notifyListeners();
  }

  void _triggerReload() {
    loadDishes();
  }

  Future<void> loadDishes() async {
    final storeId = _storeSelector.selectedStoreId;
    final category = _categoryProvider.selectedCategory;
    if (storeId == null || category == null || category.id <= 0) {
      _dishes = [];
      _loading = false;
      _error = null;
      notifyListeners();
      return;
    }
    _loading = true;
    _error = null;
    notifyListeners();
    final result = await _api.list(storeId: storeId, categoryId: category.id);
    result.when(
      success: (list) {
        _dishes = list;
      },
      error: (msg, code) {
        _error = msg;
      },
    );
    _loading = false;
    notifyListeners();
  }

  Future<bool> createDish(CreateDishRequest payload) async {
    final storeId = _storeSelector.selectedStoreId;
    final category = _categoryProvider.selectedCategory;
    if (storeId == null || category == null) return false;
    final result = await _api.create(
        storeId: storeId, categoryId: category.id, payload: payload);
    bool ok = false;
    result.when(
      success: (item) {
        _error = null;
        ok = true;
        loadDishes();
      },
      error: (msg, code) {
        _error = msg;
      },
    );
    if (!ok) notifyListeners();
    return ok;
  }

  Future<bool> updateDish(Dish item, UpdateDishRequest payload) async {
    final storeId = _storeSelector.selectedStoreId;
    final category = _categoryProvider.selectedCategory;
    if (storeId == null || category == null) return false;
    final result = await _api.update(
        storeId: storeId,
        categoryId: category.id,
        dishId: item.id,
        payload: payload);
    bool ok = false;
    result.when(
      success: (updated) {
        _error = null;
        ok = true;
        loadDishes();
      },
      error: (msg, code) {
        _error = msg;
      },
    );
    if (!ok) notifyListeners();
    return ok;
  }

  Future<bool> deleteDish(Dish item) async {
    final storeId = _storeSelector.selectedStoreId;
    final category = _categoryProvider.selectedCategory;
    if (storeId == null || category == null) return false;
    final result = await _api.delete(
        storeId: storeId, categoryId: category.id, dishId: item.id);
    bool ok = false;
    result.when(
      success: (_) {
        _dishes.removeWhere((d) => d.id == item.id);
        _error = null; // Clear error on success
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
    _storeSelector.removeListener(_triggerReload);
    _categoryProvider.removeListener(_triggerReload);
    super.dispose();
  }
}
