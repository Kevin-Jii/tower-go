// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'models.dart';

// **************************************************************************
// JsonSerializableGenerator
// **************************************************************************

_$MenuItemImpl _$$MenuItemImplFromJson(Map<String, dynamic> json) =>
    _$MenuItemImpl(
      id: (json['id'] as num).toInt(),
      parentId: (json['parent_id'] as num).toInt(),
      name: json['name'] as String,
      title: json['title'] as String,
      icon: json['icon'] as String?,
      path: json['path'] as String?,
      component: json['component'] as String?,
      type: (json['type'] as num?)?.toInt(),
      sort: (json['sort'] as num?)?.toInt(),
      visible: (json['visible'] as num?)?.toInt(),
      status: (json['status'] as num?)?.toInt(),
      permission: json['permission'] as String?,
      permissions:
          (json['perms'] as List<dynamic>?)?.map((e) => e as String).toList() ??
              const [],
      isExternal: (json['is_external'] as num?)?.toInt(),
      cache: (json['cache'] as num?)?.toInt(),
      routeParams: json['route_params'] as String?,
      remark: json['remark'] as String?,
      createdAt: json['created_at'] as String?,
      updatedAt: json['updated_at'] as String?,
      children: (json['children'] as List<dynamic>?)
              ?.map((e) => MenuItem.fromJson(e as Map<String, dynamic>))
              .toList() ??
          const [],
    );

Map<String, dynamic> _$$MenuItemImplToJson(_$MenuItemImpl instance) =>
    <String, dynamic>{
      'id': instance.id,
      'parent_id': instance.parentId,
      'name': instance.name,
      'title': instance.title,
      'icon': instance.icon,
      'path': instance.path,
      'component': instance.component,
      'type': instance.type,
      'sort': instance.sort,
      'visible': instance.visible,
      'status': instance.status,
      'permission': instance.permission,
      'perms': instance.permissions,
      'is_external': instance.isExternal,
      'cache': instance.cache,
      'route_params': instance.routeParams,
      'remark': instance.remark,
      'created_at': instance.createdAt,
      'updated_at': instance.updatedAt,
      'children': instance.children,
    };
