import 'package:dio/dio.dart';
import '../../core/network/api_client.dart';
import '../../core/constants/app_constants.dart';
import 'models.dart';

class DishApi {
  final Dio _dio = ApiClient().dio;

  Future<List<Dish>> listDishes() async {
    try {
      final resp = await _dio.get(ApiPaths.dishes);
      final data = resp.data['data'] ?? resp.data;
      final List list = data is List ? data : (data['list'] ?? []);
      return list
          .map((e) => Dish.fromJson(Map<String, dynamic>.from(e)))
          .toList();
    } on DioException catch (e) {
      throw ApiException(e.message ?? '加载菜品失败',
          statusCode: e.response?.statusCode);
    }
  }

  Future<Dish> createDish(CreateDishRequest req) async {
    try {
      final resp = await _dio.post(ApiPaths.dishes, data: req.toJson());
      final data = resp.data['data'] ?? resp.data;
      return Dish.fromJson(Map<String, dynamic>.from(data));
    } on DioException catch (e) {
      throw ApiException(e.message ?? '新增菜品失败',
          statusCode: e.response?.statusCode);
    }
  }
}
