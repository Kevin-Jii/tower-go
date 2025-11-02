import '../../core/network/api_client.dart';
import 'models.dart';

class AuthApi {
  final ApiClient _client = ApiClient();

  Future<LoginResponse> login(LoginRequest req) async {
    // 使用通用 post，并通过 converter 创建 LoginResponse
    final loginResp = await _client.post<LoginResponse>(
      '/auth/login',
      data: req.toJson(),
      converter: (json) =>
          LoginResponse.fromJson(Map<String, dynamic>.from(json as Map)),
    );
    _client.setToken(loginResp.token);
    return loginResp;
  }
}
