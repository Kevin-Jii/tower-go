import 'package:dio/dio.dart';
import '../../core/network/api_client.dart';
import '../../core/constants/app_constants.dart';
import '../../core/network/response_utils.dart';
import '../../core/constants/error_texts.dart';
import 'models.dart';

class DishApi {
  final ApiClient _client;
  DishApi(this._client);
  Dio get _dio => _client.dio;

  // 旧的全量菜品列表（可能不再使用）
  Future<List<Dish>> listDishes() async {
    try {
      final resp = await _dio.get(ApiPaths.dishes);
      final list = ResponseUtils.payloadToList(resp.data, listKey: 'list');
      return list.map(Dish.fromJson).toList();
    } on DioException catch (e) {
      throw ApiException(e.message ?? ErrorTexts.loadDishes,
          statusCode: e.response?.statusCode);
    }
  }

  Future<Dish> createDish(CreateDishRequest req) async {
    try {
      final resp = await _dio.post(ApiPaths.dishes, data: req.toJson());
      final payload = ResponseUtils.extractData(resp.data);
      return Dish.fromJson(Map<String, dynamic>.from(payload));
    } on DioException catch (e) {
      throw ApiException(e.message ?? ErrorTexts.createDish,
          statusCode: e.response?.statusCode);
    }
  }

  String _nestedBase(int storeId, int categoryId) =>
      '/stores/$storeId/dish-categories/$categoryId/dishes';

  Future<ApiResult<List<Dish>>> list(
      {required int storeId, required int categoryId}) async {
    try {
      final resp = await _dio.get(_nestedBase(storeId, categoryId));
      final list = ResponseUtils.payloadToList(resp.data);
      final items = list.map(Dish.fromJson).toList();
      return ApiResult.success(items);
    } on DioException catch (e) {
      // 后端若分类存在但无菜品应返回 200 + []，若暂时返回 404,前端容错为空列表
      if (e.response?.statusCode == 404) {
        return ApiResult.success([]);
      }
      return ApiResult.error(ErrorTexts.loadDishes, e.response?.statusCode);
    } catch (e) {
      return ApiResult.error('${ErrorTexts.loadDishes}: $e', null);
    }
  }

  Future<ApiResult<Dish?>> create(
      {required int storeId,
      required int categoryId,
      required CreateDishRequest payload}) async {
    try {
      final resp = await _dio.post(_nestedBase(storeId, categoryId),
          data: payload.toJson());
      final data = ResponseUtils.extractData(resp.data);
      // Backend may return null for success-only response or actual dish data
      if (data == null) {
        return ApiResult.success(null);
      }
      final item = Dish.fromJson(Map<String, dynamic>.from(data));
      return ApiResult.success(item);
    } on DioException catch (e) {
      return ApiResult.error(ErrorTexts.createDish, e.response?.statusCode);
    } catch (e) {
      return ApiResult.error(ErrorTexts.createDish, null);
    }
  }

  Future<ApiResult<Dish?>> update(
      {required int storeId,
      required int categoryId,
      required int dishId,
      required UpdateDishRequest payload}) async {
    try {
      final resp = await _dio.put('${_nestedBase(storeId, categoryId)}/$dishId',
          data: payload.toJson());
      final data = ResponseUtils.extractData(resp.data);
      // Backend may return null for success-only response or actual dish data
      if (data == null) {
        return ApiResult.success(null);
      }
      final item = Dish.fromJson(Map<String, dynamic>.from(data));
      return ApiResult.success(item);
    } on DioException catch (e) {
      return ApiResult.error(ErrorTexts.updateDish, e.response?.statusCode);
    } catch (e) {
      return ApiResult.error('${ErrorTexts.updateDish}: $e', null);
    }
  }

  Future<ApiResult<void>> delete(
      {required int storeId,
      required int categoryId,
      required int dishId}) async {
    try {
      final resp =
          await _dio.delete('${_nestedBase(storeId, categoryId)}/$dishId');
      // Backend returns standard response format, extract to verify success
      ResponseUtils.extractData(resp.data);
      return ApiResult.success(null);
    } on DioException catch (e) {
      return ApiResult.error(ErrorTexts.deleteDish, e.response?.statusCode);
    } catch (e) {
      return ApiResult.error(ErrorTexts.deleteDish, null);
    }
  }
}
