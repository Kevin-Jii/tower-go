import 'dart:convert';
import '../../core/network/api_client.dart';
import '../../core/constants/app_constants.dart';
import '../../core/utils/value_parsers.dart';
import '../../core/utils/map_utils.dart';
import '../../core/constants/error_texts.dart';
import '../../core/network/crud_api.dart';
import 'role_models.dart';

/// 角色 API，继承通用 CrudApi 并保留自定义的健壮 list 解析逻辑。
class RoleApi extends CrudApi<RoleItem> {
  RoleApi(ApiClient client) : super(client);

  @override
  String get basePath => ApiPaths.roles;

  /// 将后端 JSON Map 转为 RoleItem，处理类型兼容与 remark -> description 迁移。
  @override
  RoleItem fromJson(Map<String, dynamic> map) {
    if (!map.containsKey('description') && map.containsKey('remark')) {
      map['description'] = map['remark'];
    }
    return RoleItem(
      id: parseInt(map['id']),
      name: map['name']?.toString() ?? '',
      code: map['code']?.toString() ?? '',
      description: map['description']?.toString(),
      status: parseIntNullable(map['status']),
      createdAt: parseStringNullable(map['created_at']),
      updatedAt: parseStringNullable(map['updated_at']),
    );
  }

  /// 自定义健壮列表获取：兼容字符串 JSON / Map 包裹 / 直接 List。
  Future<List<RoleItem>> getRoles({String? keyword}) async {
    try {
      final resp = await client.dio.get(basePath, queryParameters: {
        if (keyword != null && keyword.isNotEmpty) 'keyword': keyword,
      });
      dynamic body = resp.data;
      if (body is String) {
        try {
          body = body.trim().isEmpty ? [] : jsonDecode(body);
        } catch (_) {
          throw ApiException(ErrorTexts.loadRoles);
        }
      }
      List listData = [];
      if (body is List) {
        listData = body;
      } else if (body is Map) {
        final data = body['data'];
        if (data is List) {
          listData = data;
        } else if (data == null) {
          listData = [];
        } else {
          throw ApiException(ErrorTexts.loadRoles);
        }
      } else {
        throw ApiException(ErrorTexts.loadRoles);
      }
      return listData.map((raw) {
        try {
          final map = Map<String, dynamic>.from(raw as Map);
          return fromJson(map);
        } catch (e) {
          throw ApiException('${ErrorTexts.parseRoleItem}: $e 原始条目: $raw');
        }
      }).toList();
    } on ApiException {
      rethrow;
    } catch (e) {
      throw ApiException('${ErrorTexts.loadRoles}: $e');
    }
  }

  /// 兼容旧调用：创建角色
  Future<void> createRole(CreateRoleRequest req) async {
    await create(req.toJson());
  }

  /// 更新角色（过滤 null）
  Future<void> updateRole(int id, UpdateRoleRequest req) async {
    await update(id, compact(req.toJson()));
  }

  /// 删除角色
  Future<void> deleteRole(int id) async {
    await delete(id);
  }
}
