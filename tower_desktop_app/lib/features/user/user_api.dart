import '../../../core/network/api_client.dart';
import '../../../core/constants/app_constants.dart';
import 'models.dart';

class UserApi {
  final ApiClient _client;
  UserApi(this._client);

  /// 获取用户列表
  Future<UserListResponse> getUsers({
    int page = 1,
    int pageSize = 10,
    String? keyword,
    int? roleId,
    int? storeId,
    int? status,
  }) async {
    final pageResp = await _client.getPage<User>(
      ApiPaths.users,
      queryParameters: {
        'page': page,
        'page_size': pageSize,
        if (keyword != null && keyword.isNotEmpty) 'keyword': keyword,
        if (roleId != null) 'role_id': roleId,
        if (storeId != null) 'store_id': storeId,
        if (status != null) 'status': status,
      },
      itemParser: (json) => User.fromJson(json),
    );
    return UserListResponse(
      list: pageResp.list,
      total: pageResp.total,
      page: pageResp.page,
      pageSize: pageResp.pageSize,
    );
  }

  /// 创建用户
  Future<void> createUser(CreateUserRequest req) async {
    await _client.post<void>(
      ApiPaths.users,
      data: req.toJson(),
    );
  }

  /// 更新用户
  Future<void> updateUser(int id, UpdateUserRequest req) async {
    // 过滤掉 null 值,仅发送有值的字段
    final json = req.toJson();
    json.removeWhere((key, value) => value == null);

    await _client.request<void>(
      '${ApiPaths.users}/$id',
      method: 'PUT',
      data: json,
    );
  }

  /// 删除用户
  Future<void> deleteUser(int id) async {
    await _client.request('${ApiPaths.users}/$id', method: 'DELETE');
  }

  /// 批量删除用户
  Future<void> batchDelete(List<int> ids) async {
    await _client.post('${ApiPaths.users}/batch-delete', data: {'ids': ids});
  }

  /// 重置密码
  Future<void> resetPassword(int id, String newPassword) async {
    await _client.post('${ApiPaths.users}/$id/reset-password',
        data: {'password': newPassword});
  }
}
