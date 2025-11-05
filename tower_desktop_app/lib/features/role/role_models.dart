import 'package:freezed_annotation/freezed_annotation.dart';
import '../../core/utils/value_parsers.dart';
part 'role_models.freezed.dart';
part 'role_models.g.dart';

@freezed
class RoleItem with _$RoleItem {
  const factory RoleItem({
    @JsonKey(fromJson: parseInt) required int id,
    required String name,
    required String code,
    @JsonKey(name: 'description') String? description,
    @JsonKey(fromJson: parseIntNullable) int? status, // 1启用 0禁用
    @JsonKey(name: 'created_at', fromJson: parseStringNullable)
    String? createdAt,
    @JsonKey(name: 'updated_at', fromJson: parseStringNullable)
    String? updatedAt,
  }) = _RoleItem;
  factory RoleItem.fromJson(Map<String, dynamic> json) =>
      _$RoleItemFromJson(json);
}

class CreateRoleRequest {
  final String name;
  final String code;
  final int? status;
  final String? description;
  CreateRoleRequest({
    required this.name,
    required this.code,
    this.status,
    this.description,
  });
  Map<String, dynamic> toJson() => {
        'name': name,
        'code': code,
        if (status != null) 'status': status,
        if (description != null && description!.isNotEmpty)
          'description': description,
      };
}

class UpdateRoleRequest {
  final String? name;
  final String? code;
  final int? status;
  final String? description;
  UpdateRoleRequest({this.name, this.code, this.status, this.description});
  Map<String, dynamic> toJson() => {
        if (name != null) 'name': name,
        if (code != null) 'code': code,
        if (status != null) 'status': status,
        if (description != null) 'description': description,
      };
}
