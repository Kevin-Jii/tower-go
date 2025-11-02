// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'models.dart';

// **************************************************************************
// JsonSerializableGenerator
// **************************************************************************

_$MenuReportImpl _$$MenuReportImplFromJson(Map<String, dynamic> json) =>
    _$MenuReportImpl(
      id: (json['id'] as num).toInt(),
      storeId: (json['storeId'] as num).toInt(),
      dishId: (json['dishId'] as num).toInt(),
      date: json['date'] as String,
      quantity: (json['quantity'] as num?)?.toInt(),
      status: json['status'] as String?,
    );

Map<String, dynamic> _$$MenuReportImplToJson(_$MenuReportImpl instance) =>
    <String, dynamic>{
      'id': instance.id,
      'storeId': instance.storeId,
      'dishId': instance.dishId,
      'date': instance.date,
      'quantity': instance.quantity,
      'status': instance.status,
    };

_$CreateMenuReportRequestImpl _$$CreateMenuReportRequestImplFromJson(
        Map<String, dynamic> json) =>
    _$CreateMenuReportRequestImpl(
      dishId: (json['dishId'] as num).toInt(),
      date: json['date'] as String,
      quantity: (json['quantity'] as num).toInt(),
    );

Map<String, dynamic> _$$CreateMenuReportRequestImplToJson(
        _$CreateMenuReportRequestImpl instance) =>
    <String, dynamic>{
      'dishId': instance.dishId,
      'date': instance.date,
      'quantity': instance.quantity,
    };

_$UpdateMenuReportRequestImpl _$$UpdateMenuReportRequestImplFromJson(
        Map<String, dynamic> json) =>
    _$UpdateMenuReportRequestImpl(
      quantity: (json['quantity'] as num?)?.toInt(),
      status: json['status'] as String?,
    );

Map<String, dynamic> _$$UpdateMenuReportRequestImplToJson(
        _$UpdateMenuReportRequestImpl instance) =>
    <String, dynamic>{
      'quantity': instance.quantity,
      'status': instance.status,
    };
