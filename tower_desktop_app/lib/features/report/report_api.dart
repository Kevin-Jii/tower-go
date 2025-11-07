import 'package:dio/dio.dart';
import '../../core/network/api_client.dart';
import '../../core/network/response_utils.dart';
import 'models.dart';

class ReportApi {
  final ApiClient _client;
  ReportApi(this._client);
  Dio get _dio => _client.dio;

  /// 获取报菜记录列表
  Future<ApiResult<List<MenuReport>>> list({
    int? storeId,
    int? dishId,
    String? startDate,
    String? endDate,
    int? page,
    int? pageSize,
  }) async {
    try {
      final queryParams = <String, dynamic>{};
      if (storeId != null) queryParams['store_id'] = storeId;
      if (dishId != null) queryParams['dish_id'] = dishId;
      if (startDate != null) queryParams['start_date'] = startDate;
      if (endDate != null) queryParams['end_date'] = endDate;
      if (page != null) queryParams['page'] = page;
      if (pageSize != null) queryParams['page_size'] = pageSize;

      final resp = await _dio.get(
        '/admin/menu-reports',
        queryParameters: queryParams,
      );

      final list = ResponseUtils.payloadToList(resp.data);
      final items = list
          .map((e) => MenuReport.fromJson(Map<String, dynamic>.from(e)))
          .toList();
      return ApiResult.success(items);
    } on DioException catch (e) {
      if (e.response?.statusCode == 404) {
        return ApiResult.success([]);
      }
      return ApiResult.error('加载报菜记录失败', e.response?.statusCode);
    } catch (e) {
      return ApiResult.error('加载报菜记录失败: $e', null);
    }
  }

  /// 创建报菜记录
  Future<ApiResult<MenuReport?>> create(CreateMenuReportRequest payload) async {
    try {
      final resp = await _dio.post(
        '/admin/menu-reports',
        data: payload.toJson(),
      );

      final data = ResponseUtils.extractData(resp.data);
      if (data == null) {
        return ApiResult.success(null);
      }
      final item = MenuReport.fromJson(Map<String, dynamic>.from(data));
      return ApiResult.success(item);
    } on DioException catch (e) {
      return ApiResult.error('创建报菜失败', e.response?.statusCode);
    } catch (e) {
      return ApiResult.error('创建报菜失败: $e', null);
    }
  }

  /// 更新报菜记录
  Future<ApiResult<MenuReport?>> update(
    int id,
    UpdateMenuReportRequest payload,
  ) async {
    try {
      final resp = await _dio.put(
        '/admin/menu-reports/$id',
        data: payload.toJson(),
      );

      final data = ResponseUtils.extractData(resp.data);
      if (data == null) {
        return ApiResult.success(null);
      }
      final item = MenuReport.fromJson(Map<String, dynamic>.from(data));
      return ApiResult.success(item);
    } on DioException catch (e) {
      return ApiResult.error('更新报菜失败', e.response?.statusCode);
    } catch (e) {
      return ApiResult.error('更新报菜失败: $e', null);
    }
  }

  /// 删除报菜记录
  Future<ApiResult<void>> delete(int id) async {
    try {
      final resp = await _dio.delete('/admin/menu-reports/$id');
      ResponseUtils.extractData(resp.data);
      return ApiResult.success(null);
    } on DioException catch (e) {
      return ApiResult.error('删除报菜失败', e.response?.statusCode);
    } catch (e) {
      return ApiResult.error('删除报菜失败: $e', null);
    }
  }

  /// 批量删除报菜记录
  Future<ApiResult<void>> batchDelete(List<int> ids) async {
    try {
      final resp = await _dio.post(
        '/admin/menu-reports/batch-delete',
        data: {'ids': ids},
      );
      ResponseUtils.extractData(resp.data);
      return ApiResult.success(null);
    } on DioException catch (e) {
      return ApiResult.error('批量删除失败', e.response?.statusCode);
    } catch (e) {
      return ApiResult.error('批量删除失败: $e', null);
    }
  }
}
