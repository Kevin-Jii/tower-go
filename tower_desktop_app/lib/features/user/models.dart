import 'package:freezed_annotation/freezed_annotation.dart';

part 'models.freezed.dart';
part 'models.g.dart';

@freezed
class Role with _$Role {
  const factory Role({
    required int id,
    required String name,
    required String code,
    String? description,
    @JsonKey(name: 'created_at') String? createdAt,
    @JsonKey(name: 'updated_at') String? updatedAt,
  }) = _Role;

  factory Role.fromJson(Map<String, dynamic> json) => _$RoleFromJson(json);
}

@freezed
class User with _$User {
  const factory User({
    required int id,
    required String username,
    required String phone,
    String? email,
    String? nickname,
    String? avatar,
    int? gender, // 1=男 2=女
    @JsonKey(name: 'role_id') int? roleId,
    Role? role, // 嵌套的角色对象
    @JsonKey(name: 'store_id') int? storeId,
    @JsonKey(name: 'store_name') String? storeName, // 门店名称（前端填充）
    int? status, // 0=禁用 1=启用
    String? remark,
    @JsonKey(name: 'created_at') String? createdAt,
    @JsonKey(name: 'updated_at') String? updatedAt,
  }) = _User;

  factory User.fromJson(Map<String, dynamic> json) => _$UserFromJson(json);
}

@freezed
class UserListResponse with _$UserListResponse {
  const factory UserListResponse({
    required List<User> list,
    required int total,
    @Default(1) int page,
    @Default(10) @JsonKey(name: 'page_size') int pageSize,
  }) = _UserListResponse;

  factory UserListResponse.fromJson(Map<String, dynamic> json) =>
      _$UserListResponseFromJson(json);
}

@freezed
class CreateUserRequest with _$CreateUserRequest {
  const factory CreateUserRequest({
    required String username,
    required String phone,
    required String password,
    String? email,
    String? nickname,
    int? gender, // 1=男 2=女
    @JsonKey(name: 'role_id') int? roleId,
    @JsonKey(name: 'store_id') int? storeId,
    int? status,
    String? remark,
  }) = _CreateUserRequest;

  factory CreateUserRequest.fromJson(Map<String, dynamic> json) =>
      _$CreateUserRequestFromJson(json);
}

@freezed
class UpdateUserRequest with _$UpdateUserRequest {
  const factory UpdateUserRequest({
    String? phone,
    String? email,
    String? nickname,
    String? password, // 可选,为空则不修改
    @JsonKey(name: 'role_id') int? roleId,
    @JsonKey(name: 'store_id') int? storeId,
    int? status,
    String? remark,
  }) = _UpdateUserRequest;

  factory UpdateUserRequest.fromJson(Map<String, dynamic> json) =>
      _$UpdateUserRequestFromJson(json);
}
