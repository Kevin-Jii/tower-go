// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'models.dart';

// **************************************************************************
// JsonSerializableGenerator
// **************************************************************************

_$RoleImpl _$$RoleImplFromJson(Map<String, dynamic> json) => _$RoleImpl(
      id: (json['id'] as num).toInt(),
      name: json['name'] as String,
      code: json['code'] as String,
      description: json['description'] as String?,
      createdAt: json['created_at'] as String?,
      updatedAt: json['updated_at'] as String?,
    );

Map<String, dynamic> _$$RoleImplToJson(_$RoleImpl instance) =>
    <String, dynamic>{
      'id': instance.id,
      'name': instance.name,
      'code': instance.code,
      'description': instance.description,
      'created_at': instance.createdAt,
      'updated_at': instance.updatedAt,
    };

_$UserImpl _$$UserImplFromJson(Map<String, dynamic> json) => _$UserImpl(
      id: (json['id'] as num).toInt(),
      username: json['username'] as String,
      phone: json['phone'] as String,
      email: json['email'] as String?,
      nickname: json['nickname'] as String?,
      avatar: json['avatar'] as String?,
      gender: (json['gender'] as num?)?.toInt(),
      roleId: (json['role_id'] as num?)?.toInt(),
      role: json['role'] == null
          ? null
          : Role.fromJson(json['role'] as Map<String, dynamic>),
      storeId: (json['store_id'] as num?)?.toInt(),
      storeName: json['store_name'] as String?,
      status: (json['status'] as num?)?.toInt(),
      remark: json['remark'] as String?,
      createdAt: json['created_at'] as String?,
      updatedAt: json['updated_at'] as String?,
    );

Map<String, dynamic> _$$UserImplToJson(_$UserImpl instance) =>
    <String, dynamic>{
      'id': instance.id,
      'username': instance.username,
      'phone': instance.phone,
      'email': instance.email,
      'nickname': instance.nickname,
      'avatar': instance.avatar,
      'gender': instance.gender,
      'role_id': instance.roleId,
      'role': instance.role,
      'store_id': instance.storeId,
      'store_name': instance.storeName,
      'status': instance.status,
      'remark': instance.remark,
      'created_at': instance.createdAt,
      'updated_at': instance.updatedAt,
    };

_$UserListResponseImpl _$$UserListResponseImplFromJson(
        Map<String, dynamic> json) =>
    _$UserListResponseImpl(
      list: (json['list'] as List<dynamic>)
          .map((e) => User.fromJson(e as Map<String, dynamic>))
          .toList(),
      total: (json['total'] as num).toInt(),
      page: (json['page'] as num?)?.toInt() ?? 1,
      pageSize: (json['page_size'] as num?)?.toInt() ?? 10,
    );

Map<String, dynamic> _$$UserListResponseImplToJson(
        _$UserListResponseImpl instance) =>
    <String, dynamic>{
      'list': instance.list,
      'total': instance.total,
      'page': instance.page,
      'page_size': instance.pageSize,
    };

_$CreateUserRequestImpl _$$CreateUserRequestImplFromJson(
        Map<String, dynamic> json) =>
    _$CreateUserRequestImpl(
      username: json['username'] as String,
      phone: json['phone'] as String,
      password: json['password'] as String,
      email: json['email'] as String?,
      nickname: json['nickname'] as String?,
      gender: (json['gender'] as num?)?.toInt(),
      roleId: (json['role_id'] as num?)?.toInt(),
      storeId: (json['store_id'] as num?)?.toInt(),
      status: (json['status'] as num?)?.toInt(),
      remark: json['remark'] as String?,
    );

Map<String, dynamic> _$$CreateUserRequestImplToJson(
        _$CreateUserRequestImpl instance) =>
    <String, dynamic>{
      'username': instance.username,
      'phone': instance.phone,
      'password': instance.password,
      'email': instance.email,
      'nickname': instance.nickname,
      'gender': instance.gender,
      'role_id': instance.roleId,
      'store_id': instance.storeId,
      'status': instance.status,
      'remark': instance.remark,
    };

_$UpdateUserRequestImpl _$$UpdateUserRequestImplFromJson(
        Map<String, dynamic> json) =>
    _$UpdateUserRequestImpl(
      phone: json['phone'] as String?,
      email: json['email'] as String?,
      nickname: json['nickname'] as String?,
      password: json['password'] as String?,
      roleId: (json['role_id'] as num?)?.toInt(),
      storeId: (json['store_id'] as num?)?.toInt(),
      status: (json['status'] as num?)?.toInt(),
      remark: json['remark'] as String?,
    );

Map<String, dynamic> _$$UpdateUserRequestImplToJson(
        _$UpdateUserRequestImpl instance) =>
    <String, dynamic>{
      'phone': instance.phone,
      'email': instance.email,
      'nickname': instance.nickname,
      'password': instance.password,
      'role_id': instance.roleId,
      'store_id': instance.storeId,
      'status': instance.status,
      'remark': instance.remark,
    };
