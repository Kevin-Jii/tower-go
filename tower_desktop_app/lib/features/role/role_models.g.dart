// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'role_models.dart';

// **************************************************************************
// JsonSerializableGenerator
// **************************************************************************

_$RoleItemImpl _$$RoleItemImplFromJson(Map<String, dynamic> json) =>
    _$RoleItemImpl(
      id: parseInt(json['id']),
      name: json['name'] as String,
      code: json['code'] as String,
      description: json['description'] as String?,
      status: parseIntNullable(json['status']),
      createdAt: parseStringNullable(json['created_at']),
      updatedAt: parseStringNullable(json['updated_at']),
    );

Map<String, dynamic> _$$RoleItemImplToJson(_$RoleItemImpl instance) =>
    <String, dynamic>{
      'id': instance.id,
      'name': instance.name,
      'code': instance.code,
      'description': instance.description,
      'status': instance.status,
      'created_at': instance.createdAt,
      'updated_at': instance.updatedAt,
    };
