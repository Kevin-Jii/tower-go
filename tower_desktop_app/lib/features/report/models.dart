import 'package:freezed_annotation/freezed_annotation.dart';
part 'models.freezed.dart';
part 'models.g.dart';

@freezed
class MenuReport with _$MenuReport {
  const factory MenuReport({
    required int id,
    @JsonKey(name: 'store_id') required int storeId,
    @JsonKey(name: 'store_name') String? storeName,
    @JsonKey(name: 'dish_id') required int dishId,
    @JsonKey(name: 'dish_name') required String dishName,
    required int quantity,
    String? remark,
    @JsonKey(name: 'created_at') String? createdAt,
    @JsonKey(name: 'updated_at') String? updatedAt,
  }) = _MenuReport;

  factory MenuReport.fromJson(Map<String, dynamic> json) =>
      _$MenuReportFromJson(json);
}

@freezed
class CreateMenuReportRequest with _$CreateMenuReportRequest {
  const factory CreateMenuReportRequest({
    @JsonKey(name: 'dish_id') required int dishId,
    required int quantity,
    String? remark,
  }) = _CreateMenuReportRequest;

  factory CreateMenuReportRequest.fromJson(Map<String, dynamic> json) =>
      _$CreateMenuReportRequestFromJson(json);
}

@freezed
class UpdateMenuReportRequest with _$UpdateMenuReportRequest {
  const factory UpdateMenuReportRequest({
    int? quantity,
    String? remark,
  }) = _UpdateMenuReportRequest;

  factory UpdateMenuReportRequest.fromJson(Map<String, dynamic> json) =>
      _$UpdateMenuReportRequestFromJson(json);
}
