// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'models.dart';

// **************************************************************************
// JsonSerializableGenerator
// **************************************************************************

_$MenuReportImpl _$$MenuReportImplFromJson(Map<String, dynamic> json) =>
    _$MenuReportImpl(
      id: (json['id'] as num).toInt(),
      storeId: (json['store_id'] as num).toInt(),
      storeName: json['store_name'] as String?,
      dishId: (json['dish_id'] as num).toInt(),
      dishName: json['dish_name'] as String,
      quantity: (json['quantity'] as num).toInt(),
      remark: json['remark'] as String?,
      createdAt: json['created_at'] as String?,
      updatedAt: json['updated_at'] as String?,
    );

Map<String, dynamic> _$$MenuReportImplToJson(_$MenuReportImpl instance) =>
    <String, dynamic>{
      'id': instance.id,
      'store_id': instance.storeId,
      'store_name': instance.storeName,
      'dish_id': instance.dishId,
      'dish_name': instance.dishName,
      'quantity': instance.quantity,
      'remark': instance.remark,
      'created_at': instance.createdAt,
      'updated_at': instance.updatedAt,
    };

_$CreateMenuReportRequestImpl _$$CreateMenuReportRequestImplFromJson(
        Map<String, dynamic> json) =>
    _$CreateMenuReportRequestImpl(
      dishId: (json['dish_id'] as num).toInt(),
      quantity: (json['quantity'] as num).toInt(),
      remark: json['remark'] as String?,
    );

Map<String, dynamic> _$$CreateMenuReportRequestImplToJson(
        _$CreateMenuReportRequestImpl instance) =>
    <String, dynamic>{
      'dish_id': instance.dishId,
      'quantity': instance.quantity,
      'remark': instance.remark,
    };

_$UpdateMenuReportRequestImpl _$$UpdateMenuReportRequestImplFromJson(
        Map<String, dynamic> json) =>
    _$UpdateMenuReportRequestImpl(
      quantity: (json['quantity'] as num?)?.toInt(),
      remark: json['remark'] as String?,
    );

Map<String, dynamic> _$$UpdateMenuReportRequestImplToJson(
        _$UpdateMenuReportRequestImpl instance) =>
    <String, dynamic>{
      'quantity': instance.quantity,
      'remark': instance.remark,
    };
