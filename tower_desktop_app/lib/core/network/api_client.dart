import 'package:dio/dio.dart';
import 'package:flutter/foundation.dart';
import 'dart:convert';
import 'api_response.dart';

class ApiClient {
  static final ApiClient _instance = ApiClient._internal();
  factory ApiClient() => _instance;

  late final Dio dio;

  // 根据你的后端地址调整 (开发环境)
  static const String baseUrl = 'http://127.0.0.1:10024/api/v1';

  String? _token;
  void setToken(String? token) {
    _token = token;
  }

  ApiClient._internal() {
    dio = Dio(BaseOptions(
      baseUrl: baseUrl,
      connectTimeout: const Duration(seconds: 10),
      receiveTimeout: const Duration(seconds: 15),
      headers: {
        'Content-Type': 'application/json',
      },
    ));

    dio.interceptors.add(InterceptorsWrapper(
      onRequest: (options, handler) {
        if (_token != null && _token!.isNotEmpty) {
          options.headers['Authorization'] = 'Bearer $_token';
        }
        return handler.next(options);
      },
      onResponse: (response, handler) {
        return handler.next(response);
      },
      onError: (e, handler) {
        // 统一错误转换
        if (kDebugMode) {
          debugPrint(
              'API Error: ${e.type} ${e.response?.statusCode} => ${e.message}');
        }
        return handler.next(e);
      },
    ));
  }

  // 通用 GET
  Future<T> get<T>(String path,
      {Map<String, dynamic>? queryParameters,
      T Function(dynamic json)? converter}) async {
    try {
      final resp = await dio.get(path, queryParameters: queryParameters);
      return _extractData<T>(resp, converter, method: 'GET');
    } on DioException catch (e) {
      throw _toApiException(e);
    }
  }

  // 通用 POST
  Future<T> post<T>(String path,
      {Object? data,
      Map<String, dynamic>? queryParameters,
      T Function(dynamic json)? converter}) async {
    try {
      final resp =
          await dio.post(path, data: data, queryParameters: queryParameters);
      return _extractData<T>(resp, converter, method: 'POST');
    } on DioException catch (e) {
      throw _toApiException(e);
    }
  }

  // 通用 DELETE
  Future<T> delete<T>(String path,
      {Object? data,
      Map<String, dynamic>? queryParameters,
      T Function(dynamic json)? converter}) async {
    try {
      final resp =
          await dio.delete(path, data: data, queryParameters: queryParameters);
      return _extractData<T>(resp, converter, method: 'DELETE');
    } on DioException catch (e) {
      throw _toApiException(e);
    }
  }

  // 支持其它 method
  Future<T> request<T>(String path,
      {required String method,
      Object? data,
      Map<String, dynamic>? queryParameters,
      T Function(dynamic json)? converter}) async {
    try {
      final resp = await dio.request(path,
          data: data,
          queryParameters: queryParameters,
          options: Options(method: method));
      return _extractData<T>(resp, converter, method: method.toUpperCase());
    } on DioException catch (e) {
      throw _toApiException(e);
    }
  }

  /// 分页 GET 请求：自动解析 meta 并构建 PageResponse
  Future<PageResponse<T>> getPage<T>(
    String path, {
    Map<String, dynamic>? queryParameters,
    required T Function(Map<String, dynamic>) itemParser,
  }) async {
    try {
      final resp = await dio.get(path, queryParameters: queryParameters);
      final body = _validateResponse(resp);
      return PageResponse.fromEnvelope(body, itemParser);
    } on DioException catch (e) {
      throw _toApiException(e);
    }
  }

  /// 分页 POST 请求：自动解析 meta 并构建 PageResponse
  Future<PageResponse<T>> postPage<T>(
    String path, {
    Object? data,
    Map<String, dynamic>? queryParameters,
    required T Function(Map<String, dynamic>) itemParser,
  }) async {
    try {
      final resp =
          await dio.post(path, data: data, queryParameters: queryParameters);
      final body = _validateResponse(resp);
      return PageResponse.fromEnvelope(body, itemParser);
    } on DioException catch (e) {
      throw _toApiException(e);
    }
  }

  ApiException _toApiException(DioException e) {
    String msg;
    switch (e.type) {
      case DioExceptionType.connectionTimeout:
        msg = '连接服务器超时';
        break;
      case DioExceptionType.sendTimeout:
        msg = '请求发送超时';
        break;
      case DioExceptionType.receiveTimeout:
        msg = '服务器响应超时';
        break;
      case DioExceptionType.badResponse:
        msg = e.response?.data is Map && (e.response?.data['message'] != null)
            ? e.response?.data['message']
            : '服务器返回错误(${e.response?.statusCode})';
        break;
      case DioExceptionType.cancel:
        msg = '请求已取消';
        break;
      default:
        msg = e.message ?? '网络错误';
    }
    return ApiException(msg, statusCode: e.response?.statusCode);
  }

  // 统一解析响应: {code, message, data, [total, page, page_size]}
  T _extractData<T>(Response resp, T Function(dynamic json)? converter,
      {required String method}) {
    final body = _validateResponse(resp);

    // 如果使用了 converter，直接传整个响应体（包含 data、total 等）
    if (converter != null) {
      return converter(body);
    }

    // 否则只提取 data 字段
    final data = body['data'];

    // 类型安全检查
    // 当 T 为 void 时（增删改操作），data 可以为 null
    if (T == dynamic) {
      return data as T;
    }

    // data 为空：GET 认为错误，非 GET 允许（增删改成功通常无 data）
    if (data == null) {
      if (method == 'GET') {
        throw ApiException('data 字段为空 (GET ${resp.requestOptions.path})');
      }
      return null as T; // POST/PUT/DELETE 允许无 data
    }

    if (T == Map<String, dynamic>) {
      if (data is Map) {
        return Map<String, dynamic>.from(data) as T;
      }
      throw ApiException('data 字段格式不正确，期望 Map，实际为 ${data.runtimeType}');
    }

    if (T.toString().contains('List')) {
      if (data is List) {
        return data as T;
      }
      throw ApiException('data 字段格式不正确，期望 List，实际为 ${data.runtimeType}');
    }

    return data as T;
  }

  /// 校验响应格式和业务状态码，返回响应体 Map
  Map<String, dynamic> _validateResponse(Response resp) {
    final raw = resp.data;
    dynamic body = raw;

    // 如果响应是字符串，尝试解码为 JSON
    if (raw is String) {
      try {
        body = json.decode(raw);
      } catch (_) {
        throw ApiException('响应数据格式错误，无法解析 JSON');
      }
    }

    // 检查后端统一响应格式 {code, message, data}
    if (body is! Map) {
      throw ApiException('响应格式错误，期望 JSON 对象');
    }

    final code = body['code'];
    final message = body['message'];

    // 检查业务状态码
    if (code != null && code != 200) {
      throw ApiException(message?.toString() ?? '请求失败', statusCode: code);
    }

    return body as Map<String, dynamic>;
  }
}

class ApiException implements Exception {
  final String message;
  final int? statusCode;
  ApiException(this.message, {this.statusCode});
  @override
  String toString() => 'ApiException($statusCode): $message';
}
