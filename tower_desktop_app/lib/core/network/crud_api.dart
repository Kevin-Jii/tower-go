/// 通用 CRUD API 基类（简单场景）。
/// 适用于基础的 list/create/update/delete，不含分页与复杂查询。
import 'api_client.dart';
import 'response_utils.dart';

abstract class CrudApi<T> {
  final ApiClient client;
  CrudApi(this.client);

  /// 子类需提供资源基础路径，如 ApiPaths.roles
  String get basePath;
  T fromJson(Map<String, dynamic> json);

  Future<List<T>> list({Map<String, dynamic>? query, String? listKey}) async {
    final resp = await client.dio.get(basePath, queryParameters: query);
    final maps = ResponseUtils.payloadToList(resp.data, listKey: listKey);
    return maps.map(fromJson).toList();
  }

  Future<T?> create(Map<String, dynamic> data) async {
    final resp = await client.dio.post(basePath, data: data);
    final payload = ResponseUtils.extractData(resp.data);
    if (payload is Map) return fromJson(Map<String, dynamic>.from(payload));
    return null;
  }

  Future<T?> update(int id, Map<String, dynamic> data) async {
    final resp = await client.dio.put('$basePath/$id', data: data);
    final payload = ResponseUtils.extractData(resp.data);
    if (payload is Map) return fromJson(Map<String, dynamic>.from(payload));
    return null;
  }

  Future<void> delete(int id) async {
    await client.dio.delete('$basePath/$id');
  }
}
