import 'package:freezed_annotation/freezed_annotation.dart';
part 'models.freezed.dart';
part 'models.g.dart';

@Freezed()
class LoginRequest with _$LoginRequest {
  const factory LoginRequest({
    required String phone,
    required String password,
  }) = _LoginRequest;
  factory LoginRequest.fromJson(Map<String, dynamic> json) =>
      _$LoginRequestFromJson(json);
}

@Freezed()
class LoginResponse with _$LoginResponse {
  const factory LoginResponse({
    required String token,
    @JsonKey(name: 'token_type') required String tokenType,
    @JsonKey(name: 'expires_in') required int expiresIn,
    @JsonKey(name: 'user_info') required UserInfo userInfo,
  }) = _LoginResponse;
  factory LoginResponse.fromJson(Map<String, dynamic> json) =>
      _$LoginResponseFromJson(json);
}

@Freezed()
class UserInfo with _$UserInfo {
  const factory UserInfo({
    required int id,
    required String phone,
    required String username,
    @Default('') String nickname,
    @Default('') String email,
    @JsonKey(name: 'store_id') required int storeId,
    @JsonKey(name: 'role_id') required int roleId,
    required Role role,
    required int status,
    @JsonKey(name: 'last_login_at') required String lastLoginAt,
    @JsonKey(name: 'created_at') required String createdAt,
    @JsonKey(name: 'updated_at') required String updatedAt,
  }) = _UserInfo;
  factory UserInfo.fromJson(Map<String, dynamic> json) =>
      _$UserInfoFromJson(json);
}

@Freezed()
class Role with _$Role {
  const factory Role({
    required int id,
    required String name,
    required String code,
    required String description,
    @JsonKey(name: 'created_at') required String createdAt,
    @JsonKey(name: 'updated_at') required String updatedAt,
  }) = _Role;
  factory Role.fromJson(Map<String, dynamic> json) => _$RoleFromJson(json);
}
