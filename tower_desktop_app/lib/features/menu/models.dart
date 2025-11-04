import 'package:freezed_annotation/freezed_annotation.dart';
part 'models.freezed.dart';
part 'models.g.dart';

@freezed
class MenuItem with _$MenuItem {
  const factory MenuItem({
    required int id,
    @JsonKey(name: 'parent_id') required int parentId,
    required String name,
    required String title,
    String? icon,
    String? path,
    String? component,
    int? type,
    int? sort,
    int? visible,
    int? status,
    // 后端有的可能是单个 permission，也可能未来扩展成数组 perms
    String? permission,
    @JsonKey(name: 'perms') @Default([]) List<String> permissions, // 兼容数组形式
    // 是否外链: 1是 0否（假设后端字段 is_external）
    @JsonKey(name: 'is_external') int? isExternal,
    // 是否缓存: 1缓存 0不缓存（假设后端字段 cache 或 keep_alive，这里使用 cache）
    @JsonKey(name: 'cache') int? cache,
    // 路由参数（假设后端字段 route_params，使用字符串保存原始表达式）
    @JsonKey(name: 'route_params') String? routeParams,
    String? remark,
    @JsonKey(name: 'created_at') String? createdAt,
    @JsonKey(name: 'updated_at') String? updatedAt,
    @Default([]) List<MenuItem> children,
  }) = _MenuItem;

  factory MenuItem.fromJson(Map<String, dynamic> json) =>
      _$MenuItemFromJson(json);
}

/// 创建菜单请求
class CreateMenuRequest {
  final int parentId;
  final String name;
  final String title;
  final int type; // 1目录 2页面 3按钮
  final int? sort;
  final String? icon;
  final String? path;
  final String? component;
  final String? permission; // 按钮或页面的权限编码
  final int? visible; // 1显示 0隐藏
  final int? status; // 1启用 0禁用
  final int? isExternal; // 1外链 0非外链
  final int? cache; // 1缓存 0不缓存
  final String? routeParams; // 路由参数字符串

  CreateMenuRequest({
    required this.parentId,
    required this.name,
    required this.title,
    required this.type,
    this.sort,
    this.icon,
    this.path,
    this.component,
    this.permission,
    this.visible,
    this.status,
    this.isExternal,
    this.cache,
    this.routeParams,
  });

  Map<String, dynamic> toJson() => {
        'parent_id': parentId,
        'name': name,
        'title': title,
        'type': type,
        if (sort != null) 'sort': sort,
        if (icon != null && icon!.isNotEmpty) 'icon': icon,
        if (path != null && path!.isNotEmpty) 'path': path,
        if (component != null && component!.isNotEmpty) 'component': component,
        if (permission != null && permission!.isNotEmpty)
          'permission': permission,
        if (visible != null) 'visible': visible,
        if (status != null) 'status': status,
        if (isExternal != null) 'is_external': isExternal,
        if (cache != null) 'cache': cache,
        if (routeParams != null && routeParams!.isNotEmpty)
          'route_params': routeParams,
      };
}

/// 更新菜单请求（字段可选）
class UpdateMenuRequest {
  final int? parentId;
  final String? name;
  final String? title;
  final int? type;
  final int? sort;
  final String? icon;
  final String? path;
  final String? component;
  final String? permission;
  final int? visible;
  final int? status;
  final int? isExternal;
  final int? cache;
  final String? routeParams;

  UpdateMenuRequest({
    this.parentId,
    this.name,
    this.title,
    this.type,
    this.sort,
    this.icon,
    this.path,
    this.component,
    this.permission,
    this.visible,
    this.status,
    this.isExternal,
    this.cache,
    this.routeParams,
  });

  Map<String, dynamic> toJson() => {
        if (parentId != null) 'parent_id': parentId,
        if (name != null) 'name': name,
        if (title != null) 'title': title,
        if (type != null) 'type': type,
        if (sort != null) 'sort': sort,
        if (icon != null) 'icon': icon,
        if (path != null) 'path': path,
        if (component != null) 'component': component,
        if (permission != null) 'permission': permission,
        if (visible != null) 'visible': visible,
        if (status != null) 'status': status,
        if (isExternal != null) 'is_external': isExternal,
        if (cache != null) 'cache': cache,
        if (routeParams != null) 'route_params': routeParams,
      };
}
