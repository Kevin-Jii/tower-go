import 'package:freezed_annotation/freezed_annotation.dart';
part 'models.freezed.dart';
part 'models.g.dart';

@freezed
class LoginRequest with _$LoginRequest {
  const factory LoginRequest({
    required String phone,
    required String password,
  }) = _LoginRequest;
  factory LoginRequest.fromJson(Map<String, dynamic> json) =>
      _$LoginRequestFromJson(json);
}

@freezed
class LoginResponse with _$LoginResponse {
  const factory LoginResponse({
    @Default('') String token,
    @JsonKey(name: 'token_type') @Default('Bearer') String tokenType,
    @JsonKey(name: 'expires_in') @Default(0) int expiresIn,
    @JsonKey(name: 'user_info') required UserInfo userInfo,
  }) = _LoginResponse;
  factory LoginResponse.fromJson(Map<String, dynamic> json) =>
      _$LoginResponseFromJson(json);
}

@freezed
class UserInfo with _$UserInfo {
  const factory UserInfo({
    required int id,
    @Default('') String phone,
    @Default('') String username,
    @Default('') String nickname,
    @Default('') String email,
    @JsonKey(name: 'store_id') @Default(0) int storeId,
    @JsonKey(name: 'role_id') @Default(0) int roleId,
    Role? role,
    @Default(1) int status,
    @JsonKey(name: 'last_login_at') @Default('') String lastLoginAt,
    @JsonKey(name: 'created_at') @Default('') String createdAt,
    @JsonKey(name: 'updated_at') @Default('') String updatedAt,
  }) = _UserInfo;
  factory UserInfo.fromJson(Map<String, dynamic> json) =>
      _$UserInfoFromJson(json);
}

@freezed
class Role with _$Role {
  const factory Role({
    @Default(0) int id,
    @Default('') String name,
    @Default('') String code,
    @Default('') String description,
    @JsonKey(name: 'created_at') @Default('') String createdAt,
    @JsonKey(name: 'updated_at') @Default('') String updatedAt,
  }) = _Role;
  factory Role.fromJson(Map<String, dynamic> json) => _$RoleFromJson(json);
}
