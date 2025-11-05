import 'package:dio/dio.dart';
import '../../core/network/api_client.dart';
import '../../core/constants/app_constants.dart';
import '../../core/network/response_utils.dart';
import '../../core/constants/error_texts.dart';
import '../../core/utils/map_utils.dart';
import 'models.dart';

/// 菜单管理 & 用户菜单 API
/// 假设后端 REST 设计：
/// GET    /menus            分页查询菜单
/// POST   /menus            创建菜单
/// PUT    /menus/{id}       更新菜单
/// DELETE /menus/{id}       删除菜单
/// GET    /menus/user-menus 当前用户可见树
/// GET    /menus/user-permissions 当前用户权限
/// 如与后端不符，可在此集中适配。

class MenuApi {
  final Dio _dio = ApiClient().dio;
  final ApiClient _client = ApiClient();

  Future<List<MenuItem>> getUserMenus() async {
    try {
      final resp = await _dio.get('${ApiPaths.menus}/user-menus');
      final list = ResponseUtils.payloadToList(resp.data, listKey: 'menus');
      return list.map(MenuItem.fromJson).toList();
    } on DioException catch (e) {
      throw ApiException(e.message ?? ErrorTexts.loadMenus,
          statusCode: e.response?.statusCode);
    }
  }

  Future<List<String>> getUserPermissions() async {
    try {
      final resp = await _dio.get('${ApiPaths.menus}/user-permissions');
      final list =
          ResponseUtils.payloadToList(resp.data, listKey: 'permissions');
      return list
          .map((e) => e['permission']?.toString() ?? e.values.first.toString())
          .toList();
    } on DioException catch (e) {
      throw ApiException(e.message ?? ErrorTexts.loadPermissions,
          statusCode: e.response?.statusCode);
    }
  }

  /// 获取完整菜单树（管理界面使用，不带分页 meta）
  Future<List<MenuItem>> getMenuTree() async {
    try {
      final resp = await _dio.get(ApiPaths.menus); // 后端直接返回树或包裹
      final list = ResponseUtils.payloadToList(resp.data, listKey: 'menus');
      return list.map(MenuItem.fromJson).toList();
    } on DioException catch (e) {
      throw ApiException(e.message ?? ErrorTexts.loadMenuTree,
          statusCode: e.response?.statusCode);
    }
  }

  /// 创建菜单
  Future<void> createMenu(CreateMenuRequest req) async {
    await _client.post<void>(ApiPaths.menus, data: req.toJson());
  }

  /// 更新菜单
  Future<void> updateMenu(int id, UpdateMenuRequest req) async {
    final data = compact(req.toJson());
    await _client.request<void>('${ApiPaths.menus}/$id',
        method: 'PUT', data: data);
  }

  /// 删除菜单
  Future<void> deleteMenu(int id) async {
    await _client.request<void>('${ApiPaths.menus}/$id', method: 'DELETE');
  }
}
