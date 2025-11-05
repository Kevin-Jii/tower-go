// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'models.dart';

// **************************************************************************
// JsonSerializableGenerator
// **************************************************************************

_$DishImpl _$$DishImplFromJson(Map<String, dynamic> json) => _$DishImpl(
      id: (json['id'] as num).toInt(),
      storeId: (json['store_id'] as num).toInt(),
      name: json['name'] as String,
      description: json['description'] as String?,
      price: parseDoubleNullable(json['price']),
      categoryId: (json['category_id'] as num?)?.toInt(),
      status: (json['status'] as num?)?.toInt(),
      image: json['image'] as String?,
      remark: json['remark'] as String?,
      createdAt: json['created_at'] as String?,
      updatedAt: json['updated_at'] as String?,
    );

Map<String, dynamic> _$$DishImplToJson(_$DishImpl instance) =>
    <String, dynamic>{
      'id': instance.id,
      'store_id': instance.storeId,
      'name': instance.name,
      'description': instance.description,
      'price': instance.price,
      'category_id': instance.categoryId,
      'status': instance.status,
      'image': instance.image,
      'remark': instance.remark,
      'created_at': instance.createdAt,
      'updated_at': instance.updatedAt,
    };

_$CreateDishRequestImpl _$$CreateDishRequestImplFromJson(
        Map<String, dynamic> json) =>
    _$CreateDishRequestImpl(
      name: json['name'] as String,
      description: json['description'] as String?,
      price: (json['price'] as num?)?.toDouble(),
    );

Map<String, dynamic> _$$CreateDishRequestImplToJson(
        _$CreateDishRequestImpl instance) =>
    <String, dynamic>{
      'name': instance.name,
      'description': instance.description,
      'price': instance.price,
    };

_$UpdateDishRequestImpl _$$UpdateDishRequestImplFromJson(
        Map<String, dynamic> json) =>
    _$UpdateDishRequestImpl(
      name: json['name'] as String?,
      description: json['description'] as String?,
      price: (json['price'] as num?)?.toDouble(),
      status: (json['status'] as num?)?.toInt(),
    );

Map<String, dynamic> _$$UpdateDishRequestImplToJson(
        _$UpdateDishRequestImpl instance) =>
    <String, dynamic>{
      'name': instance.name,
      'description': instance.description,
      'price': instance.price,
      'status': instance.status,
    };
