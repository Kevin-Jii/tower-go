import 'package:freezed_annotation/freezed_annotation.dart';
part 'role_models.freezed.dart';
part 'role_models.g.dart';

@freezed
class RoleItem with _$RoleItem {
  const factory RoleItem({
    required int id,
    required String name,
    required String code,
    String? remark,
    int? status, // 1启用 0禁用
    @JsonKey(name: 'created_at') String? createdAt,
    @JsonKey(name: 'updated_at') String? updatedAt,
  }) = _RoleItem;
  factory RoleItem.fromJson(Map<String, dynamic> json) =>
      _$RoleItemFromJson(json);
}

class CreateRoleRequest {
  final String name;
  final String code;
  final int? status;
  final String? remark;
  CreateRoleRequest({
    required this.name,
    required this.code,
    this.status,
    this.remark,
  });
  Map<String, dynamic> toJson() => {
        'name': name,
        'code': code,
        if (status != null) 'status': status,
        if (remark != null && remark!.isNotEmpty) 'remark': remark,
      };
}

class UpdateRoleRequest {
  final String? name;
  final String? code;
  final int? status;
  final String? remark;
  UpdateRoleRequest({this.name, this.code, this.status, this.remark});
  Map<String, dynamic> toJson() => {
        if (name != null) 'name': name,
        if (code != null) 'code': code,
        if (status != null) 'status': status,
        if (remark != null) 'remark': remark,
      };
}
