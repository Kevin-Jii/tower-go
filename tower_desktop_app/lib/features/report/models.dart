import 'package:freezed_annotation/freezed_annotation.dart';
part 'models.freezed.dart';
part 'models.g.dart';

@freezed
class MenuReport with _$MenuReport {
  const factory MenuReport({
    required int id,
    required int storeId,
    required int dishId,
    required String date,
    int? quantity,
    String? status,
  }) = _MenuReport;

  factory MenuReport.fromJson(Map<String, dynamic> json) => _$MenuReportFromJson(json);
}

@freezed
class CreateMenuReportRequest with _$CreateMenuReportRequest {
  const factory CreateMenuReportRequest({
    required int dishId,
    required String date,
    required int quantity,
  }) = _CreateMenuReportRequest;

  factory CreateMenuReportRequest.fromJson(Map<String, dynamic> json) => _$CreateMenuReportRequestFromJson(json);
}

@freezed
class UpdateMenuReportRequest with _$UpdateMenuReportRequest {
  const factory UpdateMenuReportRequest({
    int? quantity,
    String? status,
  }) = _UpdateMenuReportRequest;

  factory UpdateMenuReportRequest.fromJson(Map<String, dynamic> json) => _$UpdateMenuReportRequestFromJson(json);
}
