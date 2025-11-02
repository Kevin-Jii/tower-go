import 'package:dio/dio.dart';
import 'package:flutter/foundation.dart';
import 'dart:convert';

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
      return _extractData<T>(resp, converter);
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
      return _extractData<T>(resp, converter);
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
      return _extractData<T>(resp, converter);
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

  // 统一解析 data 节点: {code,message,data}
  T _extractData<T>(Response resp, T Function(dynamic json)? converter) {
    final raw = resp.data;
    dynamic body = raw;
    if (raw is String) {
      try {
        body = json.decode(raw);
      } catch (_) {}
    }
    if (body is Map && body.containsKey('data')) {
      body = body['data'];
    }
    if (converter != null) {
      return converter(body);
    }
    // 如果 T 是 Map<String,dynamic>
    if (T == Map<String, dynamic>) {
      if (body is Map<String, dynamic>) return body as T;
      throw ApiException('响应格式不正确，期望 Map');
    }
    // 如果 T 是 List<dynamic>
    if (T == List) {
      if (body is List) return body as T;
      throw ApiException('响应格式不正确，期望 List');
    }
    return body as T; // 可能抛出 cast 错误，交给调用方
  }
}

class ApiException implements Exception {
  final String message;
  final int? statusCode;
  ApiException(this.message, {this.statusCode});
  @override
  String toString() => 'ApiException($statusCode): $message';
}
