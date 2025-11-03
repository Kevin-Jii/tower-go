import '../../core/network/api_client.dart';
import '../../core/network/api_response.dart';
import 'models.dart';

class StoreApi {
  final ApiClient _client;

  StoreApi(this._client);

  /// 获取门店列表（分页）
  Future<PageResponse<Store>> getStores({
    int page = 1,
    int pageSize = 10,
    String? keyword,
  }) async {
    return _client.getPage<Store>(
      '/stores',
      itemParser: Store.fromJson,
      queryParameters: {
        'page': page,
        'page_size': pageSize,
        if (keyword != null && keyword.isNotEmpty) 'keyword': keyword,
      },
    );
  }

  /// 获取所有门店列表（不分页，用于下拉选择）
  Future<List<Store>> listStores() async {
    try {
      return await _client.getListSmart<Store>(
        path: '/stores/all',
        fromJson: Store.fromJson,
      );
    } catch (e) {
      // 如果接口不存在，返回空列表
      return [];
    }
  }

  /// 创建门店
  /// 返回创建的门店对象（如果后端返回），否则返回 null
  Future<Store?> createStore(CreateStoreRequest request) async {
    return await _client.postSmart<Store>(
      path: '/stores',
      data: request.toJson(),
      fromJson: Store.fromJson,
      filterNulls: true, // 自动过滤 null 值
    );
  }

  /// 更新门店
  /// 返回更新后的门店对象（如果后端返回），否则返回 null
  Future<Store?> updateStore(int id, UpdateStoreRequest request) async {
    return await _client.putSmart<Store>(
      path: '/stores/$id',
      data: request.toJson(),
      fromJson: Store.fromJson,
      filterNulls: true, // 自动过滤 null 值
    );
  }

  /// 删除门店
  Future<void> deleteStore(int id) async {
    await _client.deleteSmart(path: '/stores/$id');
  }
}
