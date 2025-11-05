import 'package:freezed_annotation/freezed_annotation.dart';
import '../../core/utils/value_parsers.dart';

part 'models.freezed.dart';
part 'models.g.dart';

@freezed
class Dish with _$Dish {
  const factory Dish({
    required int id,
    @JsonKey(name: 'store_id') required int storeId,
    required String name,
    String? description,
    @JsonKey(fromJson: parseDoubleNullable) double? price,
    @JsonKey(name: 'category_id') int? categoryId,
    int? status,
    String? image,
    String? remark,
    @JsonKey(name: 'created_at') String? createdAt,
    @JsonKey(name: 'updated_at') String? updatedAt,
  }) = _Dish;

  factory Dish.fromJson(Map<String, dynamic> json) => _$DishFromJson(json);
}

@freezed
class CreateDishRequest with _$CreateDishRequest {
  const factory CreateDishRequest({
    required String name,
    String? description,
    double? price,
  }) = _CreateDishRequest;

  factory CreateDishRequest.fromJson(Map<String, dynamic> json) =>
      _$CreateDishRequestFromJson(json);
}

@freezed
class UpdateDishRequest with _$UpdateDishRequest {
  const factory UpdateDishRequest({
    String? name,
    String? description,
    double? price,
    int? status,
  }) = _UpdateDishRequest;

  factory UpdateDishRequest.fromJson(Map<String, dynamic> json) =>
      _$UpdateDishRequestFromJson(json);
}
