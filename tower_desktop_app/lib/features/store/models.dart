import 'package:freezed_annotation/freezed_annotation.dart';
part 'models.freezed.dart';
part 'models.g.dart';

@freezed
class Store with _$Store {
  const factory Store({
    required int id,
    required String name,
    String? address,
    int? status, // 0=禁用 1=启用
    String? remark,
    @JsonKey(name: 'business_hours') String? businessHours,
    @JsonKey(name: 'contact_person') String? contactPerson, // 负责人（字符串类型）
    String? phone, // 负责人电话
    @JsonKey(name: 'created_at') String? createdAt,
    @JsonKey(name: 'updated_at') String? updatedAt,
  }) = _Store;

  factory Store.fromJson(Map<String, dynamic> json) => _$StoreFromJson(json);
}

@freezed
class CreateStoreRequest with _$CreateStoreRequest {
  const factory CreateStoreRequest({
    required String name,
    String? address,
    @Default(1) int status,
    String? remark,
    @JsonKey(name: 'business_hours') String? businessHours,
    @JsonKey(name: 'contact_person') String? contactPerson,
    String? phone,
  }) = _CreateStoreRequest;

  factory CreateStoreRequest.fromJson(Map<String, dynamic> json) =>
      _$CreateStoreRequestFromJson(json);
}

@freezed
class UpdateStoreRequest with _$UpdateStoreRequest {
  const factory UpdateStoreRequest({
    String? name,
    String? address,
    int? status,
    String? remark,
    @JsonKey(name: 'business_hours') String? businessHours,
    @JsonKey(name: 'contact_person') String? contactPerson,
    String? phone,
  }) = _UpdateStoreRequest;

  factory UpdateStoreRequest.fromJson(Map<String, dynamic> json) =>
      _$UpdateStoreRequestFromJson(json);
}
