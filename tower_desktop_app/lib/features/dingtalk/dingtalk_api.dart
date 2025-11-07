import 'package:dio/dio.dart';
import '../../core/network/api_client.dart';
import '../../core/network/response_utils.dart';
import 'models.dart';

class DingTalkApi {
  final ApiClient _client;
  DingTalkApi(this._client);
  Dio get _dio => _client.dio;

  String get _basePath => '/dingtalk/robots';

  /// 获取钉钉机器人列表
  Future<ApiResult<List<DingTalkRobot>>> list({int? storeId}) async {
    try {
      final params = <String, dynamic>{};
      if (storeId != null) params['store_id'] = storeId;

      final resp = await _dio.get(_basePath, queryParameters: params);
      print('钉钉机器人列表响应: ${resp.data}');

      final list = ResponseUtils.payloadToList(resp.data);
      print('解析后的列表: $list');

      final items = list.map((e) {
        print('解析机器人数据: $e');
        final robot = DingTalkRobot.fromJson(Map<String, dynamic>.from(e));
        print(
            '解析后的机器人: id=${robot.id}, name=${robot.name}, webhook="${robot.webhook}", clientId="${robot.clientId}", token="${robot.token}"');
        return robot;
      }).toList();

      print('机器人数量: ${items.length}');
      return ApiResult.success(items);
    } on DioException catch (e) {
      print('DioException: ${e.message}');
      return ApiResult.error('加载钉钉机器人失败', e.response?.statusCode);
    } catch (e) {
      print('解析错误: $e');
      return ApiResult.error('加载钉钉机器人失败: $e', null);
    }
  }

  /// 创建钉钉机器人
  Future<ApiResult<DingTalkRobot?>> create(
      CreateDingTalkRobotRequest payload) async {
    try {
      final resp = await _dio.post(_basePath, data: payload.toJson());
      final data = ResponseUtils.extractData(resp.data);

      if (data == null) {
        return ApiResult.success(null);
      }

      final item = DingTalkRobot.fromJson(Map<String, dynamic>.from(data));
      return ApiResult.success(item);
    } on DioException catch (e) {
      return ApiResult.error('创建钉钉机器人失败', e.response?.statusCode);
    } catch (e) {
      return ApiResult.error('创建钉钉机器人失败: $e', null);
    }
  }

  /// 更新钉钉机器人
  Future<ApiResult<DingTalkRobot?>> update(
      int id, UpdateDingTalkRobotRequest payload) async {
    try {
      final resp = await _dio.put('$_basePath/$id', data: payload.toJson());
      final data = ResponseUtils.extractData(resp.data);

      if (data == null) {
        return ApiResult.success(null);
      }

      final item = DingTalkRobot.fromJson(Map<String, dynamic>.from(data));
      return ApiResult.success(item);
    } on DioException catch (e) {
      return ApiResult.error('更新钉钉机器人失败', e.response?.statusCode);
    } catch (e) {
      return ApiResult.error('更新钉钉机器人失败: $e', null);
    }
  }

  /// 删除钉钉机器人
  Future<ApiResult<void>> delete(int id) async {
    try {
      final resp = await _dio.delete('$_basePath/$id');
      ResponseUtils.extractData(resp.data);
      return ApiResult.success(null);
    } on DioException catch (e) {
      return ApiResult.error('删除钉钉机器人失败', e.response?.statusCode);
    } catch (e) {
      return ApiResult.error('删除钉钉机器人失败: $e', null);
    }
  }

  /// 测试钉钉机器人推送
  Future<ApiResult<void>> testRobot(DingTalkRobot robot) async {
    try {
      final resp =
          await _dio.post('$_basePath/${robot.id}/test', data: robot.toJson());
      ResponseUtils.extractData(resp.data);
      return ApiResult.success(null);
    } on DioException catch (e) {
      return ApiResult.error('测试推送失败', e.response?.statusCode);
    } catch (e) {
      return ApiResult.error('测试推送失败: $e', null);
    }
  }
}
