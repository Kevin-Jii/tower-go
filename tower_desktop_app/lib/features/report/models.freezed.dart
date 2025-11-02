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

MenuReport _$MenuReportFromJson(Map<String, dynamic> json) {
  return _MenuReport.fromJson(json);
}

/// @nodoc
mixin _$MenuReport {
  int get id => throw _privateConstructorUsedError;
  int get storeId => throw _privateConstructorUsedError;
  int get dishId => throw _privateConstructorUsedError;
  String get date => throw _privateConstructorUsedError;
  int? get quantity => throw _privateConstructorUsedError;
  String? get status => throw _privateConstructorUsedError;

  /// Serializes this MenuReport to a JSON map.
  Map<String, dynamic> toJson() => throw _privateConstructorUsedError;

  /// Create a copy of MenuReport
  /// with the given fields replaced by the non-null parameter values.
  @JsonKey(includeFromJson: false, includeToJson: false)
  $MenuReportCopyWith<MenuReport> get copyWith =>
      throw _privateConstructorUsedError;
}

/// @nodoc
abstract class $MenuReportCopyWith<$Res> {
  factory $MenuReportCopyWith(
          MenuReport value, $Res Function(MenuReport) then) =
      _$MenuReportCopyWithImpl<$Res, MenuReport>;
  @useResult
  $Res call(
      {int id,
      int storeId,
      int dishId,
      String date,
      int? quantity,
      String? status});
}

/// @nodoc
class _$MenuReportCopyWithImpl<$Res, $Val extends MenuReport>
    implements $MenuReportCopyWith<$Res> {
  _$MenuReportCopyWithImpl(this._value, this._then);

  // ignore: unused_field
  final $Val _value;
  // ignore: unused_field
  final $Res Function($Val) _then;

  /// Create a copy of MenuReport
  /// with the given fields replaced by the non-null parameter values.
  @pragma('vm:prefer-inline')
  @override
  $Res call({
    Object? id = null,
    Object? storeId = null,
    Object? dishId = null,
    Object? date = null,
    Object? quantity = freezed,
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
      dishId: null == dishId
          ? _value.dishId
          : dishId // ignore: cast_nullable_to_non_nullable
              as int,
      date: null == date
          ? _value.date
          : date // ignore: cast_nullable_to_non_nullable
              as String,
      quantity: freezed == quantity
          ? _value.quantity
          : quantity // ignore: cast_nullable_to_non_nullable
              as int?,
      status: freezed == status
          ? _value.status
          : status // ignore: cast_nullable_to_non_nullable
              as String?,
    ) as $Val);
  }
}

/// @nodoc
abstract class _$$MenuReportImplCopyWith<$Res>
    implements $MenuReportCopyWith<$Res> {
  factory _$$MenuReportImplCopyWith(
          _$MenuReportImpl value, $Res Function(_$MenuReportImpl) then) =
      __$$MenuReportImplCopyWithImpl<$Res>;
  @override
  @useResult
  $Res call(
      {int id,
      int storeId,
      int dishId,
      String date,
      int? quantity,
      String? status});
}

/// @nodoc
class __$$MenuReportImplCopyWithImpl<$Res>
    extends _$MenuReportCopyWithImpl<$Res, _$MenuReportImpl>
    implements _$$MenuReportImplCopyWith<$Res> {
  __$$MenuReportImplCopyWithImpl(
      _$MenuReportImpl _value, $Res Function(_$MenuReportImpl) _then)
      : super(_value, _then);

  /// Create a copy of MenuReport
  /// with the given fields replaced by the non-null parameter values.
  @pragma('vm:prefer-inline')
  @override
  $Res call({
    Object? id = null,
    Object? storeId = null,
    Object? dishId = null,
    Object? date = null,
    Object? quantity = freezed,
    Object? status = freezed,
  }) {
    return _then(_$MenuReportImpl(
      id: null == id
          ? _value.id
          : id // ignore: cast_nullable_to_non_nullable
              as int,
      storeId: null == storeId
          ? _value.storeId
          : storeId // ignore: cast_nullable_to_non_nullable
              as int,
      dishId: null == dishId
          ? _value.dishId
          : dishId // ignore: cast_nullable_to_non_nullable
              as int,
      date: null == date
          ? _value.date
          : date // ignore: cast_nullable_to_non_nullable
              as String,
      quantity: freezed == quantity
          ? _value.quantity
          : quantity // ignore: cast_nullable_to_non_nullable
              as int?,
      status: freezed == status
          ? _value.status
          : status // ignore: cast_nullable_to_non_nullable
              as String?,
    ));
  }
}

/// @nodoc
@JsonSerializable()
class _$MenuReportImpl implements _MenuReport {
  const _$MenuReportImpl(
      {required this.id,
      required this.storeId,
      required this.dishId,
      required this.date,
      this.quantity,
      this.status});

  factory _$MenuReportImpl.fromJson(Map<String, dynamic> json) =>
      _$$MenuReportImplFromJson(json);

  @override
  final int id;
  @override
  final int storeId;
  @override
  final int dishId;
  @override
  final String date;
  @override
  final int? quantity;
  @override
  final String? status;

  @override
  String toString() {
    return 'MenuReport(id: $id, storeId: $storeId, dishId: $dishId, date: $date, quantity: $quantity, status: $status)';
  }

  @override
  bool operator ==(Object other) {
    return identical(this, other) ||
        (other.runtimeType == runtimeType &&
            other is _$MenuReportImpl &&
            (identical(other.id, id) || other.id == id) &&
            (identical(other.storeId, storeId) || other.storeId == storeId) &&
            (identical(other.dishId, dishId) || other.dishId == dishId) &&
            (identical(other.date, date) || other.date == date) &&
            (identical(other.quantity, quantity) ||
                other.quantity == quantity) &&
            (identical(other.status, status) || other.status == status));
  }

  @JsonKey(includeFromJson: false, includeToJson: false)
  @override
  int get hashCode =>
      Object.hash(runtimeType, id, storeId, dishId, date, quantity, status);

  /// Create a copy of MenuReport
  /// with the given fields replaced by the non-null parameter values.
  @JsonKey(includeFromJson: false, includeToJson: false)
  @override
  @pragma('vm:prefer-inline')
  _$$MenuReportImplCopyWith<_$MenuReportImpl> get copyWith =>
      __$$MenuReportImplCopyWithImpl<_$MenuReportImpl>(this, _$identity);

  @override
  Map<String, dynamic> toJson() {
    return _$$MenuReportImplToJson(
      this,
    );
  }
}

abstract class _MenuReport implements MenuReport {
  const factory _MenuReport(
      {required final int id,
      required final int storeId,
      required final int dishId,
      required final String date,
      final int? quantity,
      final String? status}) = _$MenuReportImpl;

  factory _MenuReport.fromJson(Map<String, dynamic> json) =
      _$MenuReportImpl.fromJson;

  @override
  int get id;
  @override
  int get storeId;
  @override
  int get dishId;
  @override
  String get date;
  @override
  int? get quantity;
  @override
  String? get status;

  /// Create a copy of MenuReport
  /// with the given fields replaced by the non-null parameter values.
  @override
  @JsonKey(includeFromJson: false, includeToJson: false)
  _$$MenuReportImplCopyWith<_$MenuReportImpl> get copyWith =>
      throw _privateConstructorUsedError;
}

CreateMenuReportRequest _$CreateMenuReportRequestFromJson(
    Map<String, dynamic> json) {
  return _CreateMenuReportRequest.fromJson(json);
}

/// @nodoc
mixin _$CreateMenuReportRequest {
  int get dishId => throw _privateConstructorUsedError;
  String get date => throw _privateConstructorUsedError;
  int get quantity => throw _privateConstructorUsedError;

  /// Serializes this CreateMenuReportRequest to a JSON map.
  Map<String, dynamic> toJson() => throw _privateConstructorUsedError;

  /// Create a copy of CreateMenuReportRequest
  /// with the given fields replaced by the non-null parameter values.
  @JsonKey(includeFromJson: false, includeToJson: false)
  $CreateMenuReportRequestCopyWith<CreateMenuReportRequest> get copyWith =>
      throw _privateConstructorUsedError;
}

/// @nodoc
abstract class $CreateMenuReportRequestCopyWith<$Res> {
  factory $CreateMenuReportRequestCopyWith(CreateMenuReportRequest value,
          $Res Function(CreateMenuReportRequest) then) =
      _$CreateMenuReportRequestCopyWithImpl<$Res, CreateMenuReportRequest>;
  @useResult
  $Res call({int dishId, String date, int quantity});
}

/// @nodoc
class _$CreateMenuReportRequestCopyWithImpl<$Res,
        $Val extends CreateMenuReportRequest>
    implements $CreateMenuReportRequestCopyWith<$Res> {
  _$CreateMenuReportRequestCopyWithImpl(this._value, this._then);

  // ignore: unused_field
  final $Val _value;
  // ignore: unused_field
  final $Res Function($Val) _then;

  /// Create a copy of CreateMenuReportRequest
  /// with the given fields replaced by the non-null parameter values.
  @pragma('vm:prefer-inline')
  @override
  $Res call({
    Object? dishId = null,
    Object? date = null,
    Object? quantity = null,
  }) {
    return _then(_value.copyWith(
      dishId: null == dishId
          ? _value.dishId
          : dishId // ignore: cast_nullable_to_non_nullable
              as int,
      date: null == date
          ? _value.date
          : date // ignore: cast_nullable_to_non_nullable
              as String,
      quantity: null == quantity
          ? _value.quantity
          : quantity // ignore: cast_nullable_to_non_nullable
              as int,
    ) as $Val);
  }
}

/// @nodoc
abstract class _$$CreateMenuReportRequestImplCopyWith<$Res>
    implements $CreateMenuReportRequestCopyWith<$Res> {
  factory _$$CreateMenuReportRequestImplCopyWith(
          _$CreateMenuReportRequestImpl value,
          $Res Function(_$CreateMenuReportRequestImpl) then) =
      __$$CreateMenuReportRequestImplCopyWithImpl<$Res>;
  @override
  @useResult
  $Res call({int dishId, String date, int quantity});
}

/// @nodoc
class __$$CreateMenuReportRequestImplCopyWithImpl<$Res>
    extends _$CreateMenuReportRequestCopyWithImpl<$Res,
        _$CreateMenuReportRequestImpl>
    implements _$$CreateMenuReportRequestImplCopyWith<$Res> {
  __$$CreateMenuReportRequestImplCopyWithImpl(
      _$CreateMenuReportRequestImpl _value,
      $Res Function(_$CreateMenuReportRequestImpl) _then)
      : super(_value, _then);

  /// Create a copy of CreateMenuReportRequest
  /// with the given fields replaced by the non-null parameter values.
  @pragma('vm:prefer-inline')
  @override
  $Res call({
    Object? dishId = null,
    Object? date = null,
    Object? quantity = null,
  }) {
    return _then(_$CreateMenuReportRequestImpl(
      dishId: null == dishId
          ? _value.dishId
          : dishId // ignore: cast_nullable_to_non_nullable
              as int,
      date: null == date
          ? _value.date
          : date // ignore: cast_nullable_to_non_nullable
              as String,
      quantity: null == quantity
          ? _value.quantity
          : quantity // ignore: cast_nullable_to_non_nullable
              as int,
    ));
  }
}

/// @nodoc
@JsonSerializable()
class _$CreateMenuReportRequestImpl implements _CreateMenuReportRequest {
  const _$CreateMenuReportRequestImpl(
      {required this.dishId, required this.date, required this.quantity});

  factory _$CreateMenuReportRequestImpl.fromJson(Map<String, dynamic> json) =>
      _$$CreateMenuReportRequestImplFromJson(json);

  @override
  final int dishId;
  @override
  final String date;
  @override
  final int quantity;

  @override
  String toString() {
    return 'CreateMenuReportRequest(dishId: $dishId, date: $date, quantity: $quantity)';
  }

  @override
  bool operator ==(Object other) {
    return identical(this, other) ||
        (other.runtimeType == runtimeType &&
            other is _$CreateMenuReportRequestImpl &&
            (identical(other.dishId, dishId) || other.dishId == dishId) &&
            (identical(other.date, date) || other.date == date) &&
            (identical(other.quantity, quantity) ||
                other.quantity == quantity));
  }

  @JsonKey(includeFromJson: false, includeToJson: false)
  @override
  int get hashCode => Object.hash(runtimeType, dishId, date, quantity);

  /// Create a copy of CreateMenuReportRequest
  /// with the given fields replaced by the non-null parameter values.
  @JsonKey(includeFromJson: false, includeToJson: false)
  @override
  @pragma('vm:prefer-inline')
  _$$CreateMenuReportRequestImplCopyWith<_$CreateMenuReportRequestImpl>
      get copyWith => __$$CreateMenuReportRequestImplCopyWithImpl<
          _$CreateMenuReportRequestImpl>(this, _$identity);

  @override
  Map<String, dynamic> toJson() {
    return _$$CreateMenuReportRequestImplToJson(
      this,
    );
  }
}

abstract class _CreateMenuReportRequest implements CreateMenuReportRequest {
  const factory _CreateMenuReportRequest(
      {required final int dishId,
      required final String date,
      required final int quantity}) = _$CreateMenuReportRequestImpl;

  factory _CreateMenuReportRequest.fromJson(Map<String, dynamic> json) =
      _$CreateMenuReportRequestImpl.fromJson;

  @override
  int get dishId;
  @override
  String get date;
  @override
  int get quantity;

  /// Create a copy of CreateMenuReportRequest
  /// with the given fields replaced by the non-null parameter values.
  @override
  @JsonKey(includeFromJson: false, includeToJson: false)
  _$$CreateMenuReportRequestImplCopyWith<_$CreateMenuReportRequestImpl>
      get copyWith => throw _privateConstructorUsedError;
}

UpdateMenuReportRequest _$UpdateMenuReportRequestFromJson(
    Map<String, dynamic> json) {
  return _UpdateMenuReportRequest.fromJson(json);
}

/// @nodoc
mixin _$UpdateMenuReportRequest {
  int? get quantity => throw _privateConstructorUsedError;
  String? get status => throw _privateConstructorUsedError;

  /// Serializes this UpdateMenuReportRequest to a JSON map.
  Map<String, dynamic> toJson() => throw _privateConstructorUsedError;

  /// Create a copy of UpdateMenuReportRequest
  /// with the given fields replaced by the non-null parameter values.
  @JsonKey(includeFromJson: false, includeToJson: false)
  $UpdateMenuReportRequestCopyWith<UpdateMenuReportRequest> get copyWith =>
      throw _privateConstructorUsedError;
}

/// @nodoc
abstract class $UpdateMenuReportRequestCopyWith<$Res> {
  factory $UpdateMenuReportRequestCopyWith(UpdateMenuReportRequest value,
          $Res Function(UpdateMenuReportRequest) then) =
      _$UpdateMenuReportRequestCopyWithImpl<$Res, UpdateMenuReportRequest>;
  @useResult
  $Res call({int? quantity, String? status});
}

/// @nodoc
class _$UpdateMenuReportRequestCopyWithImpl<$Res,
        $Val extends UpdateMenuReportRequest>
    implements $UpdateMenuReportRequestCopyWith<$Res> {
  _$UpdateMenuReportRequestCopyWithImpl(this._value, this._then);

  // ignore: unused_field
  final $Val _value;
  // ignore: unused_field
  final $Res Function($Val) _then;

  /// Create a copy of UpdateMenuReportRequest
  /// with the given fields replaced by the non-null parameter values.
  @pragma('vm:prefer-inline')
  @override
  $Res call({
    Object? quantity = freezed,
    Object? status = freezed,
  }) {
    return _then(_value.copyWith(
      quantity: freezed == quantity
          ? _value.quantity
          : quantity // ignore: cast_nullable_to_non_nullable
              as int?,
      status: freezed == status
          ? _value.status
          : status // ignore: cast_nullable_to_non_nullable
              as String?,
    ) as $Val);
  }
}

/// @nodoc
abstract class _$$UpdateMenuReportRequestImplCopyWith<$Res>
    implements $UpdateMenuReportRequestCopyWith<$Res> {
  factory _$$UpdateMenuReportRequestImplCopyWith(
          _$UpdateMenuReportRequestImpl value,
          $Res Function(_$UpdateMenuReportRequestImpl) then) =
      __$$UpdateMenuReportRequestImplCopyWithImpl<$Res>;
  @override
  @useResult
  $Res call({int? quantity, String? status});
}

/// @nodoc
class __$$UpdateMenuReportRequestImplCopyWithImpl<$Res>
    extends _$UpdateMenuReportRequestCopyWithImpl<$Res,
        _$UpdateMenuReportRequestImpl>
    implements _$$UpdateMenuReportRequestImplCopyWith<$Res> {
  __$$UpdateMenuReportRequestImplCopyWithImpl(
      _$UpdateMenuReportRequestImpl _value,
      $Res Function(_$UpdateMenuReportRequestImpl) _then)
      : super(_value, _then);

  /// Create a copy of UpdateMenuReportRequest
  /// with the given fields replaced by the non-null parameter values.
  @pragma('vm:prefer-inline')
  @override
  $Res call({
    Object? quantity = freezed,
    Object? status = freezed,
  }) {
    return _then(_$UpdateMenuReportRequestImpl(
      quantity: freezed == quantity
          ? _value.quantity
          : quantity // ignore: cast_nullable_to_non_nullable
              as int?,
      status: freezed == status
          ? _value.status
          : status // ignore: cast_nullable_to_non_nullable
              as String?,
    ));
  }
}

/// @nodoc
@JsonSerializable()
class _$UpdateMenuReportRequestImpl implements _UpdateMenuReportRequest {
  const _$UpdateMenuReportRequestImpl({this.quantity, this.status});

  factory _$UpdateMenuReportRequestImpl.fromJson(Map<String, dynamic> json) =>
      _$$UpdateMenuReportRequestImplFromJson(json);

  @override
  final int? quantity;
  @override
  final String? status;

  @override
  String toString() {
    return 'UpdateMenuReportRequest(quantity: $quantity, status: $status)';
  }

  @override
  bool operator ==(Object other) {
    return identical(this, other) ||
        (other.runtimeType == runtimeType &&
            other is _$UpdateMenuReportRequestImpl &&
            (identical(other.quantity, quantity) ||
                other.quantity == quantity) &&
            (identical(other.status, status) || other.status == status));
  }

  @JsonKey(includeFromJson: false, includeToJson: false)
  @override
  int get hashCode => Object.hash(runtimeType, quantity, status);

  /// Create a copy of UpdateMenuReportRequest
  /// with the given fields replaced by the non-null parameter values.
  @JsonKey(includeFromJson: false, includeToJson: false)
  @override
  @pragma('vm:prefer-inline')
  _$$UpdateMenuReportRequestImplCopyWith<_$UpdateMenuReportRequestImpl>
      get copyWith => __$$UpdateMenuReportRequestImplCopyWithImpl<
          _$UpdateMenuReportRequestImpl>(this, _$identity);

  @override
  Map<String, dynamic> toJson() {
    return _$$UpdateMenuReportRequestImplToJson(
      this,
    );
  }
}

abstract class _UpdateMenuReportRequest implements UpdateMenuReportRequest {
  const factory _UpdateMenuReportRequest(
      {final int? quantity,
      final String? status}) = _$UpdateMenuReportRequestImpl;

  factory _UpdateMenuReportRequest.fromJson(Map<String, dynamic> json) =
      _$UpdateMenuReportRequestImpl.fromJson;

  @override
  int? get quantity;
  @override
  String? get status;

  /// Create a copy of UpdateMenuReportRequest
  /// with the given fields replaced by the non-null parameter values.
  @override
  @JsonKey(includeFromJson: false, includeToJson: false)
  _$$UpdateMenuReportRequestImplCopyWith<_$UpdateMenuReportRequestImpl>
      get copyWith => throw _privateConstructorUsedError;
}
