import 'package:dio/dio.dart';
import '../../core/network/api_client.dart';
import 'models.dart';

class StoreApi {
  final Dio _dio = ApiClient().dio;

  Future<List<Store>> listStores() async {
    try {
      final resp = await _dio.get('/stores');
      final data = resp.data['data'] ?? resp.data;
      final List list = data is List ? data : (data['list'] ?? []);
      return list.map((e) => Store.fromJson(Map<String, dynamic>.from(e))).toList();
    } on DioException catch (e) {
      throw ApiException(e.message ?? '加载门店失败', statusCode: e.response?.statusCode);
    }
  }
}
