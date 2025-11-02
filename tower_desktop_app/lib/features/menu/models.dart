import 'package:freezed_annotation/freezed_annotation.dart';
part 'models.freezed.dart';
part 'models.g.dart';

@freezed
class MenuItem with _$MenuItem {
  const factory MenuItem({
    required int id,
    required int parentId,
    required String name,
    required String title,
    String? icon,
    String? path,
    String? component,
    int? type,
    int? sort,
    int? visible,
    int? status,
    String? permission,
    @Default([]) List<MenuItem> children,
  }) = _MenuItem;

  factory MenuItem.fromJson(Map<String, dynamic> json) => _$MenuItemFromJson(json);
}
