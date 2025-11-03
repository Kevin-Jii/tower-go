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

Store _$StoreFromJson(Map<String, dynamic> json) {
  return _Store.fromJson(json);
}

/// @nodoc
mixin _$Store {
  int get id => throw _privateConstructorUsedError;
  String get name => throw _privateConstructorUsedError;
  String? get address => throw _privateConstructorUsedError;
  int? get status => throw _privateConstructorUsedError; // 0=禁用 1=启用
  String? get remark => throw _privateConstructorUsedError;
  @JsonKey(name: 'business_hours')
  String? get businessHours => throw _privateConstructorUsedError;
  @JsonKey(name: 'contact_person')
  String? get contactPerson => throw _privateConstructorUsedError; // 负责人（字符串类型）
  String? get phone => throw _privateConstructorUsedError; // 负责人电话
  @JsonKey(name: 'created_at')
  String? get createdAt => throw _privateConstructorUsedError;
  @JsonKey(name: 'updated_at')
  String? get updatedAt => throw _privateConstructorUsedError;

  /// Serializes this Store to a JSON map.
  Map<String, dynamic> toJson() => throw _privateConstructorUsedError;

  /// Create a copy of Store
  /// with the given fields replaced by the non-null parameter values.
  @JsonKey(includeFromJson: false, includeToJson: false)
  $StoreCopyWith<Store> get copyWith => throw _privateConstructorUsedError;
}

/// @nodoc
abstract class $StoreCopyWith<$Res> {
  factory $StoreCopyWith(Store value, $Res Function(Store) then) =
      _$StoreCopyWithImpl<$Res, Store>;
  @useResult
  $Res call(
      {int id,
      String name,
      String? address,
      int? status,
      String? remark,
      @JsonKey(name: 'business_hours') String? businessHours,
      @JsonKey(name: 'contact_person') String? contactPerson,
      String? phone,
      @JsonKey(name: 'created_at') String? createdAt,
      @JsonKey(name: 'updated_at') String? updatedAt});
}

/// @nodoc
class _$StoreCopyWithImpl<$Res, $Val extends Store>
    implements $StoreCopyWith<$Res> {
  _$StoreCopyWithImpl(this._value, this._then);

  // ignore: unused_field
  final $Val _value;
  // ignore: unused_field
  final $Res Function($Val) _then;

  /// Create a copy of Store
  /// with the given fields replaced by the non-null parameter values.
  @pragma('vm:prefer-inline')
  @override
  $Res call({
    Object? id = null,
    Object? name = null,
    Object? address = freezed,
    Object? status = freezed,
    Object? remark = freezed,
    Object? businessHours = freezed,
    Object? contactPerson = freezed,
    Object? phone = freezed,
    Object? createdAt = freezed,
    Object? updatedAt = freezed,
  }) {
    return _then(_value.copyWith(
      id: null == id
          ? _value.id
          : id // ignore: cast_nullable_to_non_nullable
              as int,
      name: null == name
          ? _value.name
          : name // ignore: cast_nullable_to_non_nullable
              as String,
      address: freezed == address
          ? _value.address
          : address // ignore: cast_nullable_to_non_nullable
              as String?,
      status: freezed == status
          ? _value.status
          : status // ignore: cast_nullable_to_non_nullable
              as int?,
      remark: freezed == remark
          ? _value.remark
          : remark // ignore: cast_nullable_to_non_nullable
              as String?,
      businessHours: freezed == businessHours
          ? _value.businessHours
          : businessHours // ignore: cast_nullable_to_non_nullable
              as String?,
      contactPerson: freezed == contactPerson
          ? _value.contactPerson
          : contactPerson // ignore: cast_nullable_to_non_nullable
              as String?,
      phone: freezed == phone
          ? _value.phone
          : phone // ignore: cast_nullable_to_non_nullable
              as String?,
      createdAt: freezed == createdAt
          ? _value.createdAt
          : createdAt // ignore: cast_nullable_to_non_nullable
              as String?,
      updatedAt: freezed == updatedAt
          ? _value.updatedAt
          : updatedAt // ignore: cast_nullable_to_non_nullable
              as String?,
    ) as $Val);
  }
}

/// @nodoc
abstract class _$$StoreImplCopyWith<$Res> implements $StoreCopyWith<$Res> {
  factory _$$StoreImplCopyWith(
          _$StoreImpl value, $Res Function(_$StoreImpl) then) =
      __$$StoreImplCopyWithImpl<$Res>;
  @override
  @useResult
  $Res call(
      {int id,
      String name,
      String? address,
      int? status,
      String? remark,
      @JsonKey(name: 'business_hours') String? businessHours,
      @JsonKey(name: 'contact_person') String? contactPerson,
      String? phone,
      @JsonKey(name: 'created_at') String? createdAt,
      @JsonKey(name: 'updated_at') String? updatedAt});
}

/// @nodoc
class __$$StoreImplCopyWithImpl<$Res>
    extends _$StoreCopyWithImpl<$Res, _$StoreImpl>
    implements _$$StoreImplCopyWith<$Res> {
  __$$StoreImplCopyWithImpl(
      _$StoreImpl _value, $Res Function(_$StoreImpl) _then)
      : super(_value, _then);

  /// Create a copy of Store
  /// with the given fields replaced by the non-null parameter values.
  @pragma('vm:prefer-inline')
  @override
  $Res call({
    Object? id = null,
    Object? name = null,
    Object? address = freezed,
    Object? status = freezed,
    Object? remark = freezed,
    Object? businessHours = freezed,
    Object? contactPerson = freezed,
    Object? phone = freezed,
    Object? createdAt = freezed,
    Object? updatedAt = freezed,
  }) {
    return _then(_$StoreImpl(
      id: null == id
          ? _value.id
          : id // ignore: cast_nullable_to_non_nullable
              as int,
      name: null == name
          ? _value.name
          : name // ignore: cast_nullable_to_non_nullable
              as String,
      address: freezed == address
          ? _value.address
          : address // ignore: cast_nullable_to_non_nullable
              as String?,
      status: freezed == status
          ? _value.status
          : status // ignore: cast_nullable_to_non_nullable
              as int?,
      remark: freezed == remark
          ? _value.remark
          : remark // ignore: cast_nullable_to_non_nullable
              as String?,
      businessHours: freezed == businessHours
          ? _value.businessHours
          : businessHours // ignore: cast_nullable_to_non_nullable
              as String?,
      contactPerson: freezed == contactPerson
          ? _value.contactPerson
          : contactPerson // ignore: cast_nullable_to_non_nullable
              as String?,
      phone: freezed == phone
          ? _value.phone
          : phone // ignore: cast_nullable_to_non_nullable
              as String?,
      createdAt: freezed == createdAt
          ? _value.createdAt
          : createdAt // ignore: cast_nullable_to_non_nullable
              as String?,
      updatedAt: freezed == updatedAt
          ? _value.updatedAt
          : updatedAt // ignore: cast_nullable_to_non_nullable
              as String?,
    ));
  }
}

/// @nodoc
@JsonSerializable()
class _$StoreImpl implements _Store {
  const _$StoreImpl(
      {required this.id,
      required this.name,
      this.address,
      this.status,
      this.remark,
      @JsonKey(name: 'business_hours') this.businessHours,
      @JsonKey(name: 'contact_person') this.contactPerson,
      this.phone,
      @JsonKey(name: 'created_at') this.createdAt,
      @JsonKey(name: 'updated_at') this.updatedAt});

  factory _$StoreImpl.fromJson(Map<String, dynamic> json) =>
      _$$StoreImplFromJson(json);

  @override
  final int id;
  @override
  final String name;
  @override
  final String? address;
  @override
  final int? status;
// 0=禁用 1=启用
  @override
  final String? remark;
  @override
  @JsonKey(name: 'business_hours')
  final String? businessHours;
  @override
  @JsonKey(name: 'contact_person')
  final String? contactPerson;
// 负责人（字符串类型）
  @override
  final String? phone;
// 负责人电话
  @override
  @JsonKey(name: 'created_at')
  final String? createdAt;
  @override
  @JsonKey(name: 'updated_at')
  final String? updatedAt;

  @override
  String toString() {
    return 'Store(id: $id, name: $name, address: $address, status: $status, remark: $remark, businessHours: $businessHours, contactPerson: $contactPerson, phone: $phone, createdAt: $createdAt, updatedAt: $updatedAt)';
  }

  @override
  bool operator ==(Object other) {
    return identical(this, other) ||
        (other.runtimeType == runtimeType &&
            other is _$StoreImpl &&
            (identical(other.id, id) || other.id == id) &&
            (identical(other.name, name) || other.name == name) &&
            (identical(other.address, address) || other.address == address) &&
            (identical(other.status, status) || other.status == status) &&
            (identical(other.remark, remark) || other.remark == remark) &&
            (identical(other.businessHours, businessHours) ||
                other.businessHours == businessHours) &&
            (identical(other.contactPerson, contactPerson) ||
                other.contactPerson == contactPerson) &&
            (identical(other.phone, phone) || other.phone == phone) &&
            (identical(other.createdAt, createdAt) ||
                other.createdAt == createdAt) &&
            (identical(other.updatedAt, updatedAt) ||
                other.updatedAt == updatedAt));
  }

  @JsonKey(includeFromJson: false, includeToJson: false)
  @override
  int get hashCode => Object.hash(runtimeType, id, name, address, status,
      remark, businessHours, contactPerson, phone, createdAt, updatedAt);

  /// Create a copy of Store
  /// with the given fields replaced by the non-null parameter values.
  @JsonKey(includeFromJson: false, includeToJson: false)
  @override
  @pragma('vm:prefer-inline')
  _$$StoreImplCopyWith<_$StoreImpl> get copyWith =>
      __$$StoreImplCopyWithImpl<_$StoreImpl>(this, _$identity);

  @override
  Map<String, dynamic> toJson() {
    return _$$StoreImplToJson(
      this,
    );
  }
}

abstract class _Store implements Store {
  const factory _Store(
      {required final int id,
      required final String name,
      final String? address,
      final int? status,
      final String? remark,
      @JsonKey(name: 'business_hours') final String? businessHours,
      @JsonKey(name: 'contact_person') final String? contactPerson,
      final String? phone,
      @JsonKey(name: 'created_at') final String? createdAt,
      @JsonKey(name: 'updated_at') final String? updatedAt}) = _$StoreImpl;

  factory _Store.fromJson(Map<String, dynamic> json) = _$StoreImpl.fromJson;

  @override
  int get id;
  @override
  String get name;
  @override
  String? get address;
  @override
  int? get status; // 0=禁用 1=启用
  @override
  String? get remark;
  @override
  @JsonKey(name: 'business_hours')
  String? get businessHours;
  @override
  @JsonKey(name: 'contact_person')
  String? get contactPerson; // 负责人（字符串类型）
  @override
  String? get phone; // 负责人电话
  @override
  @JsonKey(name: 'created_at')
  String? get createdAt;
  @override
  @JsonKey(name: 'updated_at')
  String? get updatedAt;

  /// Create a copy of Store
  /// with the given fields replaced by the non-null parameter values.
  @override
  @JsonKey(includeFromJson: false, includeToJson: false)
  _$$StoreImplCopyWith<_$StoreImpl> get copyWith =>
      throw _privateConstructorUsedError;
}

CreateStoreRequest _$CreateStoreRequestFromJson(Map<String, dynamic> json) {
  return _CreateStoreRequest.fromJson(json);
}

/// @nodoc
mixin _$CreateStoreRequest {
  String get name => throw _privateConstructorUsedError;
  String? get address => throw _privateConstructorUsedError;
  int get status => throw _privateConstructorUsedError;
  String? get remark => throw _privateConstructorUsedError;
  @JsonKey(name: 'business_hours')
  String? get businessHours => throw _privateConstructorUsedError;
  @JsonKey(name: 'contact_person')
  String? get contactPerson => throw _privateConstructorUsedError;
  String? get phone => throw _privateConstructorUsedError;

  /// Serializes this CreateStoreRequest to a JSON map.
  Map<String, dynamic> toJson() => throw _privateConstructorUsedError;

  /// Create a copy of CreateStoreRequest
  /// with the given fields replaced by the non-null parameter values.
  @JsonKey(includeFromJson: false, includeToJson: false)
  $CreateStoreRequestCopyWith<CreateStoreRequest> get copyWith =>
      throw _privateConstructorUsedError;
}

/// @nodoc
abstract class $CreateStoreRequestCopyWith<$Res> {
  factory $CreateStoreRequestCopyWith(
          CreateStoreRequest value, $Res Function(CreateStoreRequest) then) =
      _$CreateStoreRequestCopyWithImpl<$Res, CreateStoreRequest>;
  @useResult
  $Res call(
      {String name,
      String? address,
      int status,
      String? remark,
      @JsonKey(name: 'business_hours') String? businessHours,
      @JsonKey(name: 'contact_person') String? contactPerson,
      String? phone});
}

/// @nodoc
class _$CreateStoreRequestCopyWithImpl<$Res, $Val extends CreateStoreRequest>
    implements $CreateStoreRequestCopyWith<$Res> {
  _$CreateStoreRequestCopyWithImpl(this._value, this._then);

  // ignore: unused_field
  final $Val _value;
  // ignore: unused_field
  final $Res Function($Val) _then;

  /// Create a copy of CreateStoreRequest
  /// with the given fields replaced by the non-null parameter values.
  @pragma('vm:prefer-inline')
  @override
  $Res call({
    Object? name = null,
    Object? address = freezed,
    Object? status = null,
    Object? remark = freezed,
    Object? businessHours = freezed,
    Object? contactPerson = freezed,
    Object? phone = freezed,
  }) {
    return _then(_value.copyWith(
      name: null == name
          ? _value.name
          : name // ignore: cast_nullable_to_non_nullable
              as String,
      address: freezed == address
          ? _value.address
          : address // ignore: cast_nullable_to_non_nullable
              as String?,
      status: null == status
          ? _value.status
          : status // ignore: cast_nullable_to_non_nullable
              as int,
      remark: freezed == remark
          ? _value.remark
          : remark // ignore: cast_nullable_to_non_nullable
              as String?,
      businessHours: freezed == businessHours
          ? _value.businessHours
          : businessHours // ignore: cast_nullable_to_non_nullable
              as String?,
      contactPerson: freezed == contactPerson
          ? _value.contactPerson
          : contactPerson // ignore: cast_nullable_to_non_nullable
              as String?,
      phone: freezed == phone
          ? _value.phone
          : phone // ignore: cast_nullable_to_non_nullable
              as String?,
    ) as $Val);
  }
}

/// @nodoc
abstract class _$$CreateStoreRequestImplCopyWith<$Res>
    implements $CreateStoreRequestCopyWith<$Res> {
  factory _$$CreateStoreRequestImplCopyWith(_$CreateStoreRequestImpl value,
          $Res Function(_$CreateStoreRequestImpl) then) =
      __$$CreateStoreRequestImplCopyWithImpl<$Res>;
  @override
  @useResult
  $Res call(
      {String name,
      String? address,
      int status,
      String? remark,
      @JsonKey(name: 'business_hours') String? businessHours,
      @JsonKey(name: 'contact_person') String? contactPerson,
      String? phone});
}

/// @nodoc
class __$$CreateStoreRequestImplCopyWithImpl<$Res>
    extends _$CreateStoreRequestCopyWithImpl<$Res, _$CreateStoreRequestImpl>
    implements _$$CreateStoreRequestImplCopyWith<$Res> {
  __$$CreateStoreRequestImplCopyWithImpl(_$CreateStoreRequestImpl _value,
      $Res Function(_$CreateStoreRequestImpl) _then)
      : super(_value, _then);

  /// Create a copy of CreateStoreRequest
  /// with the given fields replaced by the non-null parameter values.
  @pragma('vm:prefer-inline')
  @override
  $Res call({
    Object? name = null,
    Object? address = freezed,
    Object? status = null,
    Object? remark = freezed,
    Object? businessHours = freezed,
    Object? contactPerson = freezed,
    Object? phone = freezed,
  }) {
    return _then(_$CreateStoreRequestImpl(
      name: null == name
          ? _value.name
          : name // ignore: cast_nullable_to_non_nullable
              as String,
      address: freezed == address
          ? _value.address
          : address // ignore: cast_nullable_to_non_nullable
              as String?,
      status: null == status
          ? _value.status
          : status // ignore: cast_nullable_to_non_nullable
              as int,
      remark: freezed == remark
          ? _value.remark
          : remark // ignore: cast_nullable_to_non_nullable
              as String?,
      businessHours: freezed == businessHours
          ? _value.businessHours
          : businessHours // ignore: cast_nullable_to_non_nullable
              as String?,
      contactPerson: freezed == contactPerson
          ? _value.contactPerson
          : contactPerson // ignore: cast_nullable_to_non_nullable
              as String?,
      phone: freezed == phone
          ? _value.phone
          : phone // ignore: cast_nullable_to_non_nullable
              as String?,
    ));
  }
}

/// @nodoc
@JsonSerializable()
class _$CreateStoreRequestImpl implements _CreateStoreRequest {
  const _$CreateStoreRequestImpl(
      {required this.name,
      this.address,
      this.status = 1,
      this.remark,
      @JsonKey(name: 'business_hours') this.businessHours,
      @JsonKey(name: 'contact_person') this.contactPerson,
      this.phone});

  factory _$CreateStoreRequestImpl.fromJson(Map<String, dynamic> json) =>
      _$$CreateStoreRequestImplFromJson(json);

  @override
  final String name;
  @override
  final String? address;
  @override
  @JsonKey()
  final int status;
  @override
  final String? remark;
  @override
  @JsonKey(name: 'business_hours')
  final String? businessHours;
  @override
  @JsonKey(name: 'contact_person')
  final String? contactPerson;
  @override
  final String? phone;

  @override
  String toString() {
    return 'CreateStoreRequest(name: $name, address: $address, status: $status, remark: $remark, businessHours: $businessHours, contactPerson: $contactPerson, phone: $phone)';
  }

  @override
  bool operator ==(Object other) {
    return identical(this, other) ||
        (other.runtimeType == runtimeType &&
            other is _$CreateStoreRequestImpl &&
            (identical(other.name, name) || other.name == name) &&
            (identical(other.address, address) || other.address == address) &&
            (identical(other.status, status) || other.status == status) &&
            (identical(other.remark, remark) || other.remark == remark) &&
            (identical(other.businessHours, businessHours) ||
                other.businessHours == businessHours) &&
            (identical(other.contactPerson, contactPerson) ||
                other.contactPerson == contactPerson) &&
            (identical(other.phone, phone) || other.phone == phone));
  }

  @JsonKey(includeFromJson: false, includeToJson: false)
  @override
  int get hashCode => Object.hash(runtimeType, name, address, status, remark,
      businessHours, contactPerson, phone);

  /// Create a copy of CreateStoreRequest
  /// with the given fields replaced by the non-null parameter values.
  @JsonKey(includeFromJson: false, includeToJson: false)
  @override
  @pragma('vm:prefer-inline')
  _$$CreateStoreRequestImplCopyWith<_$CreateStoreRequestImpl> get copyWith =>
      __$$CreateStoreRequestImplCopyWithImpl<_$CreateStoreRequestImpl>(
          this, _$identity);

  @override
  Map<String, dynamic> toJson() {
    return _$$CreateStoreRequestImplToJson(
      this,
    );
  }
}

abstract class _CreateStoreRequest implements CreateStoreRequest {
  const factory _CreateStoreRequest(
      {required final String name,
      final String? address,
      final int status,
      final String? remark,
      @JsonKey(name: 'business_hours') final String? businessHours,
      @JsonKey(name: 'contact_person') final String? contactPerson,
      final String? phone}) = _$CreateStoreRequestImpl;

  factory _CreateStoreRequest.fromJson(Map<String, dynamic> json) =
      _$CreateStoreRequestImpl.fromJson;

  @override
  String get name;
  @override
  String? get address;
  @override
  int get status;
  @override
  String? get remark;
  @override
  @JsonKey(name: 'business_hours')
  String? get businessHours;
  @override
  @JsonKey(name: 'contact_person')
  String? get contactPerson;
  @override
  String? get phone;

  /// Create a copy of CreateStoreRequest
  /// with the given fields replaced by the non-null parameter values.
  @override
  @JsonKey(includeFromJson: false, includeToJson: false)
  _$$CreateStoreRequestImplCopyWith<_$CreateStoreRequestImpl> get copyWith =>
      throw _privateConstructorUsedError;
}

UpdateStoreRequest _$UpdateStoreRequestFromJson(Map<String, dynamic> json) {
  return _UpdateStoreRequest.fromJson(json);
}

/// @nodoc
mixin _$UpdateStoreRequest {
  String? get name => throw _privateConstructorUsedError;
  String? get address => throw _privateConstructorUsedError;
  int? get status => throw _privateConstructorUsedError;
  String? get remark => throw _privateConstructorUsedError;
  @JsonKey(name: 'business_hours')
  String? get businessHours => throw _privateConstructorUsedError;
  @JsonKey(name: 'contact_person')
  String? get contactPerson => throw _privateConstructorUsedError;
  String? get phone => throw _privateConstructorUsedError;

  /// Serializes this UpdateStoreRequest to a JSON map.
  Map<String, dynamic> toJson() => throw _privateConstructorUsedError;

  /// Create a copy of UpdateStoreRequest
  /// with the given fields replaced by the non-null parameter values.
  @JsonKey(includeFromJson: false, includeToJson: false)
  $UpdateStoreRequestCopyWith<UpdateStoreRequest> get copyWith =>
      throw _privateConstructorUsedError;
}

/// @nodoc
abstract class $UpdateStoreRequestCopyWith<$Res> {
  factory $UpdateStoreRequestCopyWith(
          UpdateStoreRequest value, $Res Function(UpdateStoreRequest) then) =
      _$UpdateStoreRequestCopyWithImpl<$Res, UpdateStoreRequest>;
  @useResult
  $Res call(
      {String? name,
      String? address,
      int? status,
      String? remark,
      @JsonKey(name: 'business_hours') String? businessHours,
      @JsonKey(name: 'contact_person') String? contactPerson,
      String? phone});
}

/// @nodoc
class _$UpdateStoreRequestCopyWithImpl<$Res, $Val extends UpdateStoreRequest>
    implements $UpdateStoreRequestCopyWith<$Res> {
  _$UpdateStoreRequestCopyWithImpl(this._value, this._then);

  // ignore: unused_field
  final $Val _value;
  // ignore: unused_field
  final $Res Function($Val) _then;

  /// Create a copy of UpdateStoreRequest
  /// with the given fields replaced by the non-null parameter values.
  @pragma('vm:prefer-inline')
  @override
  $Res call({
    Object? name = freezed,
    Object? address = freezed,
    Object? status = freezed,
    Object? remark = freezed,
    Object? businessHours = freezed,
    Object? contactPerson = freezed,
    Object? phone = freezed,
  }) {
    return _then(_value.copyWith(
      name: freezed == name
          ? _value.name
          : name // ignore: cast_nullable_to_non_nullable
              as String?,
      address: freezed == address
          ? _value.address
          : address // ignore: cast_nullable_to_non_nullable
              as String?,
      status: freezed == status
          ? _value.status
          : status // ignore: cast_nullable_to_non_nullable
              as int?,
      remark: freezed == remark
          ? _value.remark
          : remark // ignore: cast_nullable_to_non_nullable
              as String?,
      businessHours: freezed == businessHours
          ? _value.businessHours
          : businessHours // ignore: cast_nullable_to_non_nullable
              as String?,
      contactPerson: freezed == contactPerson
          ? _value.contactPerson
          : contactPerson // ignore: cast_nullable_to_non_nullable
              as String?,
      phone: freezed == phone
          ? _value.phone
          : phone // ignore: cast_nullable_to_non_nullable
              as String?,
    ) as $Val);
  }
}

/// @nodoc
abstract class _$$UpdateStoreRequestImplCopyWith<$Res>
    implements $UpdateStoreRequestCopyWith<$Res> {
  factory _$$UpdateStoreRequestImplCopyWith(_$UpdateStoreRequestImpl value,
          $Res Function(_$UpdateStoreRequestImpl) then) =
      __$$UpdateStoreRequestImplCopyWithImpl<$Res>;
  @override
  @useResult
  $Res call(
      {String? name,
      String? address,
      int? status,
      String? remark,
      @JsonKey(name: 'business_hours') String? businessHours,
      @JsonKey(name: 'contact_person') String? contactPerson,
      String? phone});
}

/// @nodoc
class __$$UpdateStoreRequestImplCopyWithImpl<$Res>
    extends _$UpdateStoreRequestCopyWithImpl<$Res, _$UpdateStoreRequestImpl>
    implements _$$UpdateStoreRequestImplCopyWith<$Res> {
  __$$UpdateStoreRequestImplCopyWithImpl(_$UpdateStoreRequestImpl _value,
      $Res Function(_$UpdateStoreRequestImpl) _then)
      : super(_value, _then);

  /// Create a copy of UpdateStoreRequest
  /// with the given fields replaced by the non-null parameter values.
  @pragma('vm:prefer-inline')
  @override
  $Res call({
    Object? name = freezed,
    Object? address = freezed,
    Object? status = freezed,
    Object? remark = freezed,
    Object? businessHours = freezed,
    Object? contactPerson = freezed,
    Object? phone = freezed,
  }) {
    return _then(_$UpdateStoreRequestImpl(
      name: freezed == name
          ? _value.name
          : name // ignore: cast_nullable_to_non_nullable
              as String?,
      address: freezed == address
          ? _value.address
          : address // ignore: cast_nullable_to_non_nullable
              as String?,
      status: freezed == status
          ? _value.status
          : status // ignore: cast_nullable_to_non_nullable
              as int?,
      remark: freezed == remark
          ? _value.remark
          : remark // ignore: cast_nullable_to_non_nullable
              as String?,
      businessHours: freezed == businessHours
          ? _value.businessHours
          : businessHours // ignore: cast_nullable_to_non_nullable
              as String?,
      contactPerson: freezed == contactPerson
          ? _value.contactPerson
          : contactPerson // ignore: cast_nullable_to_non_nullable
              as String?,
      phone: freezed == phone
          ? _value.phone
          : phone // ignore: cast_nullable_to_non_nullable
              as String?,
    ));
  }
}

/// @nodoc
@JsonSerializable()
class _$UpdateStoreRequestImpl implements _UpdateStoreRequest {
  const _$UpdateStoreRequestImpl(
      {this.name,
      this.address,
      this.status,
      this.remark,
      @JsonKey(name: 'business_hours') this.businessHours,
      @JsonKey(name: 'contact_person') this.contactPerson,
      this.phone});

  factory _$UpdateStoreRequestImpl.fromJson(Map<String, dynamic> json) =>
      _$$UpdateStoreRequestImplFromJson(json);

  @override
  final String? name;
  @override
  final String? address;
  @override
  final int? status;
  @override
  final String? remark;
  @override
  @JsonKey(name: 'business_hours')
  final String? businessHours;
  @override
  @JsonKey(name: 'contact_person')
  final String? contactPerson;
  @override
  final String? phone;

  @override
  String toString() {
    return 'UpdateStoreRequest(name: $name, address: $address, status: $status, remark: $remark, businessHours: $businessHours, contactPerson: $contactPerson, phone: $phone)';
  }

  @override
  bool operator ==(Object other) {
    return identical(this, other) ||
        (other.runtimeType == runtimeType &&
            other is _$UpdateStoreRequestImpl &&
            (identical(other.name, name) || other.name == name) &&
            (identical(other.address, address) || other.address == address) &&
            (identical(other.status, status) || other.status == status) &&
            (identical(other.remark, remark) || other.remark == remark) &&
            (identical(other.businessHours, businessHours) ||
                other.businessHours == businessHours) &&
            (identical(other.contactPerson, contactPerson) ||
                other.contactPerson == contactPerson) &&
            (identical(other.phone, phone) || other.phone == phone));
  }

  @JsonKey(includeFromJson: false, includeToJson: false)
  @override
  int get hashCode => Object.hash(runtimeType, name, address, status, remark,
      businessHours, contactPerson, phone);

  /// Create a copy of UpdateStoreRequest
  /// with the given fields replaced by the non-null parameter values.
  @JsonKey(includeFromJson: false, includeToJson: false)
  @override
  @pragma('vm:prefer-inline')
  _$$UpdateStoreRequestImplCopyWith<_$UpdateStoreRequestImpl> get copyWith =>
      __$$UpdateStoreRequestImplCopyWithImpl<_$UpdateStoreRequestImpl>(
          this, _$identity);

  @override
  Map<String, dynamic> toJson() {
    return _$$UpdateStoreRequestImplToJson(
      this,
    );
  }
}

abstract class _UpdateStoreRequest implements UpdateStoreRequest {
  const factory _UpdateStoreRequest(
      {final String? name,
      final String? address,
      final int? status,
      final String? remark,
      @JsonKey(name: 'business_hours') final String? businessHours,
      @JsonKey(name: 'contact_person') final String? contactPerson,
      final String? phone}) = _$UpdateStoreRequestImpl;

  factory _UpdateStoreRequest.fromJson(Map<String, dynamic> json) =
      _$UpdateStoreRequestImpl.fromJson;

  @override
  String? get name;
  @override
  String? get address;
  @override
  int? get status;
  @override
  String? get remark;
  @override
  @JsonKey(name: 'business_hours')
  String? get businessHours;
  @override
  @JsonKey(name: 'contact_person')
  String? get contactPerson;
  @override
  String? get phone;

  /// Create a copy of UpdateStoreRequest
  /// with the given fields replaced by the non-null parameter values.
  @override
  @JsonKey(includeFromJson: false, includeToJson: false)
  _$$UpdateStoreRequestImplCopyWith<_$UpdateStoreRequestImpl> get copyWith =>
      throw _privateConstructorUsedError;
}
