// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'models.dart';

// **************************************************************************
// JsonSerializableGenerator
// **************************************************************************

_$StoreImpl _$$StoreImplFromJson(Map<String, dynamic> json) => _$StoreImpl(
      id: (json['id'] as num).toInt(),
      name: json['name'] as String,
      address: json['address'] as String?,
      status: (json['status'] as num?)?.toInt(),
      remark: json['remark'] as String?,
      businessHours: json['business_hours'] as String?,
      contactPerson: json['contact_person'] as String?,
      phone: json['phone'] as String?,
      createdAt: json['created_at'] as String?,
      updatedAt: json['updated_at'] as String?,
    );

Map<String, dynamic> _$$StoreImplToJson(_$StoreImpl instance) =>
    <String, dynamic>{
      'id': instance.id,
      'name': instance.name,
      'address': instance.address,
      'status': instance.status,
      'remark': instance.remark,
      'business_hours': instance.businessHours,
      'contact_person': instance.contactPerson,
      'phone': instance.phone,
      'created_at': instance.createdAt,
      'updated_at': instance.updatedAt,
    };

_$CreateStoreRequestImpl _$$CreateStoreRequestImplFromJson(
        Map<String, dynamic> json) =>
    _$CreateStoreRequestImpl(
      name: json['name'] as String,
      address: json['address'] as String?,
      status: (json['status'] as num?)?.toInt() ?? 1,
      remark: json['remark'] as String?,
      businessHours: json['business_hours'] as String?,
      contactPerson: json['contact_person'] as String?,
      phone: json['phone'] as String?,
    );

Map<String, dynamic> _$$CreateStoreRequestImplToJson(
        _$CreateStoreRequestImpl instance) =>
    <String, dynamic>{
      'name': instance.name,
      'address': instance.address,
      'status': instance.status,
      'remark': instance.remark,
      'business_hours': instance.businessHours,
      'contact_person': instance.contactPerson,
      'phone': instance.phone,
    };

_$UpdateStoreRequestImpl _$$UpdateStoreRequestImplFromJson(
        Map<String, dynamic> json) =>
    _$UpdateStoreRequestImpl(
      name: json['name'] as String?,
      address: json['address'] as String?,
      status: (json['status'] as num?)?.toInt(),
      remark: json['remark'] as String?,
      businessHours: json['business_hours'] as String?,
      contactPerson: json['contact_person'] as String?,
      phone: json['phone'] as String?,
    );

Map<String, dynamic> _$$UpdateStoreRequestImplToJson(
        _$UpdateStoreRequestImpl instance) =>
    <String, dynamic>{
      'name': instance.name,
      'address': instance.address,
      'status': instance.status,
      'remark': instance.remark,
      'business_hours': instance.businessHours,
      'contact_person': instance.contactPerson,
      'phone': instance.phone,
    };
