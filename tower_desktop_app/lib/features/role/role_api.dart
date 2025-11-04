import 'dart:convert';
import '../../core/network/api_client.dart';
import '../../core/constants/app_constants.dart';
import 'role_models.dart';

class RoleApi {
  final ApiClient _client;
  RoleApi(this._client);

  /// 获取角色列表，兼容两种后端响应格式：
  /// 1. 直接返回 List
  /// 2. 包装为 { code, message, data: [] }
  /// 如果是字符串响应则尝试 JSON 解析。
  Future<List<RoleItem>> getRoles({String? keyword}) async {
    try {
      final resp = await _client.dio.get(ApiPaths.roles, queryParameters: {
        if (keyword != null && keyword.isNotEmpty) 'keyword': keyword,
      });
      dynamic body = resp.data;
      // 字符串尝试转 JSON
      if (body is String) {
        try {
          body = body.trim().isEmpty ? [] : jsonDecode(body);
        } catch (_) {
          throw ApiException('角色列表响应不是合法 JSON');
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
          // data 为空视为无数据
          listData = [];
        } else {
          throw ApiException('角色列表 data 字段格式错误，期望 List');
        }
      } else {
        throw ApiException('角色列表响应格式错误，期望 List 或包裹的 Map');
      }
      return listData
          .map((e) => RoleItem.fromJson(Map<String, dynamic>.from(e)))
          .toList();
    } on ApiException {
      rethrow;
    } catch (e) {
      throw ApiException('获取角色列表失败: $e');
    }
  }

  Future<void> createRole(CreateRoleRequest req) async {
    await _client.post(ApiPaths.roles, data: req.toJson());
  }

  Future<void> updateRole(int id, UpdateRoleRequest req) async {
    final json = req.toJson();
    json.removeWhere((k, v) => v == null);
    await _client.request('${ApiPaths.roles}/$id', method: 'PUT', data: json);
  }

  Future<void> deleteRole(int id) async {
    await _client.request('${ApiPaths.roles}/$id', method: 'DELETE');
  }
}
