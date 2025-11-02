import 'package:dio/dio.dart';
import '../../core/network/api_client.dart';
import 'models.dart';

class MenuApi {
  final Dio _dio = ApiClient().dio;

  Future<List<MenuItem>> getUserMenus() async {
    try {
      final resp = await _dio.get('/menus/user-menus');
      final data = resp.data['data'] ?? resp.data;
      final List list = data is List ? data : (data['menus'] ?? []);
      return list.map((e) => MenuItem.fromJson(Map<String, dynamic>.from(e))).toList();
    } on DioException catch (e) {
      throw ApiException(e.message ?? '加载菜单失败', statusCode: e.response?.statusCode);
    }
  }

  Future<List<String>> getUserPermissions() async {
    try {
      final resp = await _dio.get('/menus/user-permissions');
      final data = resp.data['data'] ?? resp.data;
      final List list = data is List ? data : (data['permissions'] ?? []);
      return list.map((e) => e.toString()).toList();
    } on DioException catch (e) {
      throw ApiException(e.message ?? '加载权限失败', statusCode: e.response?.statusCode);
    }
  }
}
