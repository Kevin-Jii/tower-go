// coverage:ignore-file
// GENERATED CODE - DO NOT MODIFY BY HAND
// ignore_for_file: type=lint
// ignore_for_file: unused_element, deprecated_member_use, deprecated_member_use_from_same_package, use_function_type_syntax_for_parameters, unnecessary_const, avoid_init_to_null, invalid_override_different_default_values_named, prefer_expression_function_bodies, annotate_overrides, invalid_annotation_target, unnecessary_question_mark

part of 'models.dart';

// **************************************************************************
// FreezedGenerator
// **************************************************************************

T _$identity<T>(T value) => value;

final _privateConstructorUsedError = UnsupportedError(
    'It seems like you constructed your class using `MyClass._()`. This constructor is only meant to be used by freezed and you are not supposed to need it nor use it.\nPlease check the documentation here for more information: https://github.com/rrousselGit/freezed#adding-getters-and-methods-to-our-models');

Dish _$DishFromJson(Map<String, dynamic> json) {
  return _Dish.fromJson(json);
}

/// @nodoc
mixin _$Dish {
  int get id => throw _privateConstructorUsedError;
  int get storeId => throw _privateConstructorUsedError;
  String get name => throw _privateConstructorUsedError;
  String? get description => throw _privateConstructorUsedError;
  double? get price => throw _privateConstructorUsedError;
  int? get status => throw _privateConstructorUsedError;

  /// Serializes this Dish to a JSON map.
  Map<String, dynamic> toJson() => throw _privateConstructorUsedError;

  /// Create a copy of Dish
  /// with the given fields replaced by the non-null parameter values.
  @JsonKey(includeFromJson: false, includeToJson: false)
  $DishCopyWith<Dish> get copyWith => throw _privateConstructorUsedError;
}

/// @nodoc
abstract class $DishCopyWith<$Res> {
  factory $DishCopyWith(Dish value, $Res Function(Dish) then) =
      _$DishCopyWithImpl<$Res, Dish>;
  @useResult
  $Res call(
      {int id,
      int storeId,
      String name,
      String? description,
      double? price,
      int? status});
}

/// @nodoc
class _$DishCopyWithImpl<$Res, $Val extends Dish>
    implements $DishCopyWith<$Res> {
  _$DishCopyWithImpl(this._value, this._then);

  // ignore: unused_field
  final $Val _value;
  // ignore: unused_field
  final $Res Function($Val) _then;

  /// Create a copy of Dish
  /// with the given fields replaced by the non-null parameter values.
  @pragma('vm:prefer-inline')
  @override
  $Res call({
    Object? id = null,
    Object? storeId = null,
    Object? name = null,
    Object? description = freezed,
    Object? price = freezed,
    Object? status = freezed,
  }) {
    return _then(_value.copyWith(
      id: null == id
          ? _value.id
          : id // ignore: cast_nullable_to_non_nullable
              as int,
      storeId: null == storeId
          ? _value.storeId
          : storeId // ignore: cast_nullable_to_non_nullable
              as int,
      name: null == name
          ? _value.name
          : name // ignore: cast_nullable_to_non_nullable
              as String,
      description: freezed == description
          ? _value.description
          : description // ignore: cast_nullable_to_non_nullable
              as String?,
      price: freezed == price
          ? _value.price
          : price // ignore: cast_nullable_to_non_nullable
              as double?,
      status: freezed == status
          ? _value.status
          : status // ignore: cast_nullable_to_non_nullable
              as int?,
    ) as $Val);
  }
}

/// @nodoc
abstract class _$$DishImplCopyWith<$Res> implements $DishCopyWith<$Res> {
  factory _$$DishImplCopyWith(
          _$DishImpl value, $Res Function(_$DishImpl) then) =
      __$$DishImplCopyWithImpl<$Res>;
  @override
  @useResult
  $Res call(
      {int id,
      int storeId,
      String name,
      String? description,
      double? price,
      int? status});
}

/// @nodoc
class __$$DishImplCopyWithImpl<$Res>
    extends _$DishCopyWithImpl<$Res, _$DishImpl>
    implements _$$DishImplCopyWith<$Res> {
  __$$DishImplCopyWithImpl(_$DishImpl _value, $Res Function(_$DishImpl) _then)
      : super(_value, _then);

  /// Create a copy of Dish
  /// with the given fields replaced by the non-null parameter values.
  @pragma('vm:prefer-inline')
  @override
  $Res call({
    Object? id = null,
    Object? storeId = null,
    Object? name = null,
    Object? description = freezed,
    Object? price = freezed,
    Object? status = freezed,
  }) {
    return _then(_$DishImpl(
      id: null == id
          ? _value.id
          : id // ignore: cast_nullable_to_non_nullable
              as int,
      storeId: null == storeId
          ? _value.storeId
          : storeId // ignore: cast_nullable_to_non_nullable
              as int,
      name: null == name
          ? _value.name
          : name // ignore: cast_nullable_to_non_nullable
              as String,
      description: freezed == description
          ? _value.description
          : description // ignore: cast_nullable_to_non_nullable
              as String?,
      price: freezed == price
          ? _value.price
          : price // ignore: cast_nullable_to_non_nullable
              as double?,
      status: freezed == status
          ? _value.status
          : status // ignore: cast_nullable_to_non_nullable
              as int?,
    ));
  }
}

/// @nodoc
@JsonSerializable()
class _$DishImpl implements _Dish {
  const _$DishImpl(
      {required this.id,
      required this.storeId,
      required this.name,
      this.description,
      this.price,
      this.status});

  factory _$DishImpl.fromJson(Map<String, dynamic> json) =>
      _$$DishImplFromJson(json);

  @override
  final int id;
  @override
  final int storeId;
  @override
  final String name;
  @override
  final String? description;
  @override
  final double? price;
  @override
  final int? status;

  @override
  String toString() {
    return 'Dish(id: $id, storeId: $storeId, name: $name, description: $description, price: $price, status: $status)';
  }

  @override
  bool operator ==(Object other) {
    return identical(this, other) ||
        (other.runtimeType == runtimeType &&
            other is _$DishImpl &&
            (identical(other.id, id) || other.id == id) &&
            (identical(other.storeId, storeId) || other.storeId == storeId) &&
            (identical(other.name, name) || other.name == name) &&
            (identical(other.description, description) ||
                other.description == description) &&
            (identical(other.price, price) || other.price == price) &&
            (identical(other.status, status) || other.status == status));
  }

  @JsonKey(includeFromJson: false, includeToJson: false)
  @override
  int get hashCode =>
      Object.hash(runtimeType, id, storeId, name, description, price, status);

  /// Create a copy of Dish
  /// with the given fields replaced by the non-null parameter values.
  @JsonKey(includeFromJson: false, includeToJson: false)
  @override
  @pragma('vm:prefer-inline')
  _$$DishImplCopyWith<_$DishImpl> get copyWith =>
      __$$DishImplCopyWithImpl<_$DishImpl>(this, _$identity);

  @override
  Map<String, dynamic> toJson() {
    return _$$DishImplToJson(
      this,
    );
  }
}

abstract class _Dish implements Dish {
  const factory _Dish(
      {required final int id,
      required final int storeId,
      required final String name,
      final String? description,
      final double? price,
      final int? status}) = _$DishImpl;

  factory _Dish.fromJson(Map<String, dynamic> json) = _$DishImpl.fromJson;

  @override
  int get id;
  @override
  int get storeId;
  @override
  String get name;
  @override
  String? get description;
  @override
  double? get price;
  @override
  int? get status;

  /// Create a copy of Dish
  /// with the given fields replaced by the non-null parameter values.
  @override
  @JsonKey(includeFromJson: false, includeToJson: false)
  _$$DishImplCopyWith<_$DishImpl> get copyWith =>
      throw _privateConstructorUsedError;
}

CreateDishRequest _$CreateDishRequestFromJson(Map<String, dynamic> json) {
  return _CreateDishRequest.fromJson(json);
}

/// @nodoc
mixin _$CreateDishRequest {
  String get name => throw _privateConstructorUsedError;
  String? get description => throw _privateConstructorUsedError;
  double? get price => throw _privateConstructorUsedError;

  /// Serializes this CreateDishRequest to a JSON map.
  Map<String, dynamic> toJson() => throw _privateConstructorUsedError;

  /// Create a copy of CreateDishRequest
  /// with the given fields replaced by the non-null parameter values.
  @JsonKey(includeFromJson: false, includeToJson: false)
  $CreateDishRequestCopyWith<CreateDishRequest> get copyWith =>
      throw _privateConstructorUsedError;
}

/// @nodoc
abstract class $CreateDishRequestCopyWith<$Res> {
  factory $CreateDishRequestCopyWith(
          CreateDishRequest value, $Res Function(CreateDishRequest) then) =
      _$CreateDishRequestCopyWithImpl<$Res, CreateDishRequest>;
  @useResult
  $Res call({String name, String? description, double? price});
}

/// @nodoc
class _$CreateDishRequestCopyWithImpl<$Res, $Val extends CreateDishRequest>
    implements $CreateDishRequestCopyWith<$Res> {
  _$CreateDishRequestCopyWithImpl(this._value, this._then);

  // ignore: unused_field
  final $Val _value;
  // ignore: unused_field
  final $Res Function($Val) _then;

  /// Create a copy of CreateDishRequest
  /// with the given fields replaced by the non-null parameter values.
  @pragma('vm:prefer-inline')
  @override
  $Res call({
    Object? name = null,
    Object? description = freezed,
    Object? price = freezed,
  }) {
    return _then(_value.copyWith(
      name: null == name
          ? _value.name
          : name // ignore: cast_nullable_to_non_nullable
              as String,
      description: freezed == description
          ? _value.description
          : description // ignore: cast_nullable_to_non_nullable
              as String?,
      price: freezed == price
          ? _value.price
          : price // ignore: cast_nullable_to_non_nullable
              as double?,
    ) as $Val);
  }
}

/// @nodoc
abstract class _$$CreateDishRequestImplCopyWith<$Res>
    implements $CreateDishRequestCopyWith<$Res> {
  factory _$$CreateDishRequestImplCopyWith(_$CreateDishRequestImpl value,
          $Res Function(_$CreateDishRequestImpl) then) =
      __$$CreateDishRequestImplCopyWithImpl<$Res>;
  @override
  @useResult
  $Res call({String name, String? description, double? price});
}

/// @nodoc
class __$$CreateDishRequestImplCopyWithImpl<$Res>
    extends _$CreateDishRequestCopyWithImpl<$Res, _$CreateDishRequestImpl>
    implements _$$CreateDishRequestImplCopyWith<$Res> {
  __$$CreateDishRequestImplCopyWithImpl(_$CreateDishRequestImpl _value,
      $Res Function(_$CreateDishRequestImpl) _then)
      : super(_value, _then);

  /// Create a copy of CreateDishRequest
  /// with the given fields replaced by the non-null parameter values.
  @pragma('vm:prefer-inline')
  @override
  $Res call({
    Object? name = null,
    Object? description = freezed,
    Object? price = freezed,
  }) {
    return _then(_$CreateDishRequestImpl(
      name: null == name
          ? _value.name
          : name // ignore: cast_nullable_to_non_nullable
              as String,
      description: freezed == description
          ? _value.description
          : description // ignore: cast_nullable_to_non_nullable
              as String?,
      price: freezed == price
          ? _value.price
          : price // ignore: cast_nullable_to_non_nullable
              as double?,
    ));
  }
}

/// @nodoc
@JsonSerializable()
class _$CreateDishRequestImpl implements _CreateDishRequest {
  const _$CreateDishRequestImpl(
      {required this.name, this.description, this.price});

  factory _$CreateDishRequestImpl.fromJson(Map<String, dynamic> json) =>
      _$$CreateDishRequestImplFromJson(json);

  @override
  final String name;
  @override
  final String? description;
  @override
  final double? price;

  @override
  String toString() {
    return 'CreateDishRequest(name: $name, description: $description, price: $price)';
  }

  @override
  bool operator ==(Object other) {
    return identical(this, other) ||
        (other.runtimeType == runtimeType &&
            other is _$CreateDishRequestImpl &&
            (identical(other.name, name) || other.name == name) &&
            (identical(other.description, description) ||
                other.description == description) &&
            (identical(other.price, price) || other.price == price));
  }

  @JsonKey(includeFromJson: false, includeToJson: false)
  @override
  int get hashCode => Object.hash(runtimeType, name, description, price);

  /// Create a copy of CreateDishRequest
  /// with the given fields replaced by the non-null parameter values.
  @JsonKey(includeFromJson: false, includeToJson: false)
  @override
  @pragma('vm:prefer-inline')
  _$$CreateDishRequestImplCopyWith<_$CreateDishRequestImpl> get copyWith =>
      __$$CreateDishRequestImplCopyWithImpl<_$CreateDishRequestImpl>(
          this, _$identity);

  @override
  Map<String, dynamic> toJson() {
    return _$$CreateDishRequestImplToJson(
      this,
    );
  }
}

abstract class _CreateDishRequest implements CreateDishRequest {
  const factory _CreateDishRequest(
      {required final String name,
      final String? description,
      final double? price}) = _$CreateDishRequestImpl;

  factory _CreateDishRequest.fromJson(Map<String, dynamic> json) =
      _$CreateDishRequestImpl.fromJson;

  @override
  String get name;
  @override
  String? get description;
  @override
  double? get price;

  /// Create a copy of CreateDishRequest
  /// with the given fields replaced by the non-null parameter values.
  @override
  @JsonKey(includeFromJson: false, includeToJson: false)
  _$$CreateDishRequestImplCopyWith<_$CreateDishRequestImpl> get copyWith =>
      throw _privateConstructorUsedError;
}

UpdateDishRequest _$UpdateDishRequestFromJson(Map<String, dynamic> json) {
  return _UpdateDishRequest.fromJson(json);
}

/// @nodoc
mixin _$UpdateDishRequest {
  String? get name => throw _privateConstructorUsedError;
  String? get description => throw _privateConstructorUsedError;
  double? get price => throw _privateConstructorUsedError;
  int? get status => throw _privateConstructorUsedError;

  /// Serializes this UpdateDishRequest to a JSON map.
  Map<String, dynamic> toJson() => throw _privateConstructorUsedError;

  /// Create a copy of UpdateDishRequest
  /// with the given fields replaced by the non-null parameter values.
  @JsonKey(includeFromJson: false, includeToJson: false)
  $UpdateDishRequestCopyWith<UpdateDishRequest> get copyWith =>
      throw _privateConstructorUsedError;
}

/// @nodoc
abstract class $UpdateDishRequestCopyWith<$Res> {
  factory $UpdateDishRequestCopyWith(
          UpdateDishRequest value, $Res Function(UpdateDishRequest) then) =
      _$UpdateDishRequestCopyWithImpl<$Res, UpdateDishRequest>;
  @useResult
  $Res call({String? name, String? description, double? price, int? status});
}

/// @nodoc
class _$UpdateDishRequestCopyWithImpl<$Res, $Val extends UpdateDishRequest>
    implements $UpdateDishRequestCopyWith<$Res> {
  _$UpdateDishRequestCopyWithImpl(this._value, this._then);

  // ignore: unused_field
  final $Val _value;
  // ignore: unused_field
  final $Res Function($Val) _then;

  /// Create a copy of UpdateDishRequest
  /// with the given fields replaced by the non-null parameter values.
  @pragma('vm:prefer-inline')
  @override
  $Res call({
    Object? name = freezed,
    Object? description = freezed,
    Object? price = freezed,
    Object? status = freezed,
  }) {
    return _then(_value.copyWith(
      name: freezed == name
          ? _value.name
          : name // ignore: cast_nullable_to_non_nullable
              as String?,
      description: freezed == description
          ? _value.description
          : description // ignore: cast_nullable_to_non_nullable
              as String?,
      price: freezed == price
          ? _value.price
          : price // ignore: cast_nullable_to_non_nullable
              as double?,
      status: freezed == status
          ? _value.status
          : status // ignore: cast_nullable_to_non_nullable
              as int?,
    ) as $Val);
  }
}

/// @nodoc
abstract class _$$UpdateDishRequestImplCopyWith<$Res>
    implements $UpdateDishRequestCopyWith<$Res> {
  factory _$$UpdateDishRequestImplCopyWith(_$UpdateDishRequestImpl value,
          $Res Function(_$UpdateDishRequestImpl) then) =
      __$$UpdateDishRequestImplCopyWithImpl<$Res>;
  @override
  @useResult
  $Res call({String? name, String? description, double? price, int? status});
}

/// @nodoc
class __$$UpdateDishRequestImplCopyWithImpl<$Res>
    extends _$UpdateDishRequestCopyWithImpl<$Res, _$UpdateDishRequestImpl>
    implements _$$UpdateDishRequestImplCopyWith<$Res> {
  __$$UpdateDishRequestImplCopyWithImpl(_$UpdateDishRequestImpl _value,
      $Res Function(_$UpdateDishRequestImpl) _then)
      : super(_value, _then);

  /// Create a copy of UpdateDishRequest
  /// with the given fields replaced by the non-null parameter values.
  @pragma('vm:prefer-inline')
  @override
  $Res call({
    Object? name = freezed,
    Object? description = freezed,
    Object? price = freezed,
    Object? status = freezed,
  }) {
    return _then(_$UpdateDishRequestImpl(
      name: freezed == name
          ? _value.name
          : name // ignore: cast_nullable_to_non_nullable
              as String?,
      description: freezed == description
          ? _value.description
          : description // ignore: cast_nullable_to_non_nullable
              as String?,
      price: freezed == price
          ? _value.price
          : price // ignore: cast_nullable_to_non_nullable
              as double?,
      status: freezed == status
          ? _value.status
          : status // ignore: cast_nullable_to_non_nullable
              as int?,
    ));
  }
}

/// @nodoc
@JsonSerializable()
class _$UpdateDishRequestImpl implements _UpdateDishRequest {
  const _$UpdateDishRequestImpl(
      {this.name, this.description, this.price, this.status});

  factory _$UpdateDishRequestImpl.fromJson(Map<String, dynamic> json) =>
      _$$UpdateDishRequestImplFromJson(json);

  @override
  final String? name;
  @override
  final String? description;
  @override
  final double? price;
  @override
  final int? status;

  @override
  String toString() {
    return 'UpdateDishRequest(name: $name, description: $description, price: $price, status: $status)';
  }

  @override
  bool operator ==(Object other) {
    return identical(this, other) ||
        (other.runtimeType == runtimeType &&
            other is _$UpdateDishRequestImpl &&
            (identical(other.name, name) || other.name == name) &&
            (identical(other.description, description) ||
                other.description == description) &&
            (identical(other.price, price) || other.price == price) &&
            (identical(other.status, status) || other.status == status));
  }

  @JsonKey(includeFromJson: false, includeToJson: false)
  @override
  int get hashCode =>
      Object.hash(runtimeType, name, description, price, status);

  /// Create a copy of UpdateDishRequest
  /// with the given fields replaced by the non-null parameter values.
  @JsonKey(includeFromJson: false, includeToJson: false)
  @override
  @pragma('vm:prefer-inline')
  _$$UpdateDishRequestImplCopyWith<_$UpdateDishRequestImpl> get copyWith =>
      __$$UpdateDishRequestImplCopyWithImpl<_$UpdateDishRequestImpl>(
          this, _$identity);

  @override
  Map<String, dynamic> toJson() {
    return _$$UpdateDishRequestImplToJson(
      this,
    );
  }
}

abstract class _UpdateDishRequest implements UpdateDishRequest {
  const factory _UpdateDishRequest(
      {final String? name,
      final String? description,
      final double? price,
      final int? status}) = _$UpdateDishRequestImpl;

  factory _UpdateDishRequest.fromJson(Map<String, dynamic> json) =
      _$UpdateDishRequestImpl.fromJson;

  @override
  String? get name;
  @override
  String? get description;
  @override
  double? get price;
  @override
  int? get status;

  /// Create a copy of UpdateDishRequest
  /// with the given fields replaced by the non-null parameter values.
  @override
  @JsonKey(includeFromJson: false, includeToJson: false)
  _$$UpdateDishRequestImplCopyWith<_$UpdateDishRequestImpl> get copyWith =>
      throw _privateConstructorUsedError;
}
