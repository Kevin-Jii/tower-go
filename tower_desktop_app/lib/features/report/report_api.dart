import 'package:dio/dio.dart';
import '../../core/network/api_client.dart';
import 'models.dart';

class ReportApi {
  final Dio _dio = ApiClient().dio;

  Future<List<MenuReport>> listReports() async {
    try {
      final resp = await _dio.get('/menu-reports');
      final data = resp.data['data'] ?? resp.data;
      final List list = data is List ? data : (data['list'] ?? []);
      return list.map((e) => MenuReport.fromJson(Map<String, dynamic>.from(e))).toList();
    } on DioException catch (e) {
      throw ApiException(e.message ?? '加载报菜记录失败', statusCode: e.response?.statusCode);
    }
  }

  Future<MenuReport> createReport(CreateMenuReportRequest req) async {
    try {
      final resp = await _dio.post('/menu-reports', data: req.toJson());
      final data = resp.data['data'] ?? resp.data;
      return MenuReport.fromJson(Map<String, dynamic>.from(data));
    } on DioException catch (e) {
      throw ApiException(e.message ?? '创建报菜失败', statusCode: e.response?.statusCode);
    }
  }
}
