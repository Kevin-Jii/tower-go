import 'package:dio/dio.dart';
import '../../core/network/api_client.dart';
import '../../core/utils/value_parsers.dart';
import '../../core/network/response_utils.dart';
import '../../core/constants/error_texts.dart';

/// Data model for dish category returned by backend.
class DishCategory {
  final int id;
  final String name;
  final int storeId; // not in payload; injected from request context

  DishCategory({required this.id, required this.name, required this.storeId});

  factory DishCategory.fromJson(Map<String, dynamic> json,
      {required int storeId}) {
    return DishCategory(
      id: parseInt(json['id']),
      name: (json['name'] ?? '').toString(),
      storeId: storeId,
    );
  }
}

class DishCategoryApi {
  final ApiClient _client;
  DishCategoryApi(this._client);

  String _basePath(int storeId) => '/stores/$storeId/dish-categories';

  Future<ApiResult<List<DishCategory>>> list(int storeId) async {
    try {
      final Response raw = await _client.dio.get(_basePath(storeId));
      final maps = ResponseUtils.payloadToList(raw.data);
      // 后端若返回 {data: []} 或直接 [] 都能被处理；若结构异常返回空列表而不是抛出未知错误
      if (maps.isEmpty) {
        final extracted = ResponseUtils.extractData(raw.data);
        if (extracted is List && extracted.isEmpty) {
          return ApiResult.success([]);
        }
      }
      final items =
          maps.map((e) => DishCategory.fromJson(e, storeId: storeId)).toList();
      return ApiResult.success(items);
    } on DioException catch (e) {
      // 将错误简化为统一文案 + 状态码，减少 UI 噪音
      return ApiResult.error(
          ErrorTexts.loadDishCategories, e.response?.statusCode);
    } catch (_) {
      return ApiResult.error(ErrorTexts.loadDishCategories, null);
    }
  }

  Future<ApiResult<DishCategory?>> create(int storeId, String name) async {
    try {
      final Response raw = await _client.dio.post(_basePath(storeId), data: {
        'name': name,
      });
      final payload = ResponseUtils.extractData(raw.data);
      // Backend may return null for success-only response or actual category data
      if (payload == null) {
        return ApiResult.success(null);
      }
      final item = DishCategory.fromJson(Map<String, dynamic>.from(payload),
          storeId: storeId);
      return ApiResult.success(item);
    } catch (e) {
      return ApiResult.error('${ErrorTexts.createCategory}: $e', null);
    }
  }

  Future<ApiResult<DishCategory?>> update(
      int storeId, int categoryId, String name) async {
    try {
      final Response raw =
          await _client.dio.put('${_basePath(storeId)}/$categoryId', data: {
        'name': name,
      });
      final payload = ResponseUtils.extractData(raw.data);
      // Backend may return null for success-only response or actual category data
      if (payload == null) {
        return ApiResult.success(null);
      }
      final item = DishCategory.fromJson(Map<String, dynamic>.from(payload),
          storeId: storeId);
      return ApiResult.success(item);
    } catch (e) {
      return ApiResult.error('${ErrorTexts.updateCategory}: $e', null);
    }
  }

  Future<ApiResult<void>> delete(int storeId, int categoryId) async {
    try {
      final Response raw =
          await _client.dio.delete('${_basePath(storeId)}/$categoryId');
      // Backend returns standard response format, extract to verify success
      ResponseUtils.extractData(raw.data);
      return ApiResult.success(null);
    } on DioException catch (e) {
      return ApiResult.error(ErrorTexts.deleteCategory, e.response?.statusCode);
    } catch (e) {
      return ApiResult.error('${ErrorTexts.deleteCategory}: $e', null);
    }
  }
}
