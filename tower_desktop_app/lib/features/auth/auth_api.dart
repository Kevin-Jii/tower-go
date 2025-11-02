import 'package:dio/dio.dart';
import '../../core/network/api_client.dart';
import 'models.dart';

class AuthApi {
  final Dio _dio = ApiClient().dio;

  Future<LoginResponse> login(LoginRequest req) async {
    try {
      final resp = await _dio.post('/auth/login', data: req.toJson());
      final data = resp.data;
      // 假设后端统一响应: {"code":0,"message":"OK","data":{...}}
      final wrapped = data['data'] ?? data; // 兼容两种
      final loginResp = LoginResponse.fromJson(wrapped);
      ApiClient().setToken(loginResp.token);
      return loginResp;
    } on DioException catch (e) {
      throw ApiException(e.message ?? '网络错误', statusCode: e.response?.statusCode);
    }
  }
}
