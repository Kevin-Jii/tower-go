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
    String? remark,
    @JsonKey(name: 'created_at') String? createdAt,
    @JsonKey(name: 'updated_at') String? updatedAt,
    @Default([]) List<MenuItem> children,
  }) = _MenuItem;

  factory MenuItem.fromJson(Map<String, dynamic> json) =>
      _$MenuItemFromJson(json);
}
