import '../../core/network/api_client.dart';
import 'models.dart';

class AuthApi {
  final ApiClient _client = ApiClient();

  Future<LoginResponse> login(LoginRequest req) async {
    // 使用通用 post，并通过 converter 从 data 字段创建 LoginResponse
    final loginResp = await _client.post<LoginResponse>(
      '/auth/login',
      data: req.toJson(),
      converter: (json) {
        final data = json['data'];
        if (data == null) {
          throw Exception('登录响应缺少 data 字段');
        }
        return LoginResponse.fromJson(data as Map<String, dynamic>);
      },
    );
    _client.setToken(loginResp.token);
    return loginResp;
  }
}
