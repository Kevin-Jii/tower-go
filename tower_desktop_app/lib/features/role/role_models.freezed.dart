// coverage:ignore-file
// GENERATED CODE - DO NOT MODIFY BY HAND
// ignore_for_file: type=lint
// ignore_for_file: unused_element, deprecated_member_use, deprecated_member_use_from_same_package, use_function_type_syntax_for_parameters, unnecessary_const, avoid_init_to_null, invalid_override_different_default_values_named, prefer_expression_function_bodies, annotate_overrides, invalid_annotation_target, unnecessary_question_mark

part of 'role_models.dart';

// **************************************************************************
// FreezedGenerator
// **************************************************************************

T _$identity<T>(T value) => value;

final _privateConstructorUsedError = UnsupportedError(
    'It seems like you constructed your class using `MyClass._()`. This constructor is only meant to be used by freezed and you are not supposed to need it nor use it.\nPlease check the documentation here for more information: https://github.com/rrousselGit/freezed#adding-getters-and-methods-to-our-models');

RoleItem _$RoleItemFromJson(Map<String, dynamic> json) {
  return _RoleItem.fromJson(json);
}

/// @nodoc
mixin _$RoleItem {
  @JsonKey(fromJson: parseInt)
  int get id => throw _privateConstructorUsedError;
  String get name => throw _privateConstructorUsedError;
  String get code => throw _privateConstructorUsedError;
  @JsonKey(name: 'description')
  String? get description => throw _privateConstructorUsedError;
  @JsonKey(fromJson: parseIntNullable)
  int? get status => throw _privateConstructorUsedError; // 1启用 0禁用
  @JsonKey(name: 'created_at', fromJson: parseStringNullable)
  String? get createdAt => throw _privateConstructorUsedError;
  @JsonKey(name: 'updated_at', fromJson: parseStringNullable)
  String? get updatedAt => throw _privateConstructorUsedError;

  /// Serializes this RoleItem to a JSON map.
  Map<String, dynamic> toJson() => throw _privateConstructorUsedError;

  /// Create a copy of RoleItem
  /// with the given fields replaced by the non-null parameter values.
  @JsonKey(includeFromJson: false, includeToJson: false)
  $RoleItemCopyWith<RoleItem> get copyWith =>
      throw _privateConstructorUsedError;
}

/// @nodoc
abstract class $RoleItemCopyWith<$Res> {
  factory $RoleItemCopyWith(RoleItem value, $Res Function(RoleItem) then) =
      _$RoleItemCopyWithImpl<$Res, RoleItem>;
  @useResult
  $Res call(
      {@JsonKey(fromJson: parseInt) int id,
      String name,
      String code,
      @JsonKey(name: 'description') String? description,
      @JsonKey(fromJson: parseIntNullable) int? status,
      @JsonKey(name: 'created_at', fromJson: parseStringNullable)
      String? createdAt,
      @JsonKey(name: 'updated_at', fromJson: parseStringNullable)
      String? updatedAt});
}

/// @nodoc
class _$RoleItemCopyWithImpl<$Res, $Val extends RoleItem>
    implements $RoleItemCopyWith<$Res> {
  _$RoleItemCopyWithImpl(this._value, this._then);

  // ignore: unused_field
  final $Val _value;
  // ignore: unused_field
  final $Res Function($Val) _then;

  /// Create a copy of RoleItem
  /// with the given fields replaced by the non-null parameter values.
  @pragma('vm:prefer-inline')
  @override
  $Res call({
    Object? id = null,
    Object? name = null,
    Object? code = null,
    Object? description = freezed,
    Object? status = freezed,
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
      code: null == code
          ? _value.code
          : code // ignore: cast_nullable_to_non_nullable
              as String,
      description: freezed == description
          ? _value.description
          : description // ignore: cast_nullable_to_non_nullable
              as String?,
      status: freezed == status
          ? _value.status
          : status // ignore: cast_nullable_to_non_nullable
              as int?,
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
abstract class _$$RoleItemImplCopyWith<$Res>
    implements $RoleItemCopyWith<$Res> {
  factory _$$RoleItemImplCopyWith(
          _$RoleItemImpl value, $Res Function(_$RoleItemImpl) then) =
      __$$RoleItemImplCopyWithImpl<$Res>;
  @override
  @useResult
  $Res call(
      {@JsonKey(fromJson: parseInt) int id,
      String name,
      String code,
      @JsonKey(name: 'description') String? description,
      @JsonKey(fromJson: parseIntNullable) int? status,
      @JsonKey(name: 'created_at', fromJson: parseStringNullable)
      String? createdAt,
      @JsonKey(name: 'updated_at', fromJson: parseStringNullable)
      String? updatedAt});
}

/// @nodoc
class __$$RoleItemImplCopyWithImpl<$Res>
    extends _$RoleItemCopyWithImpl<$Res, _$RoleItemImpl>
    implements _$$RoleItemImplCopyWith<$Res> {
  __$$RoleItemImplCopyWithImpl(
      _$RoleItemImpl _value, $Res Function(_$RoleItemImpl) _then)
      : super(_value, _then);

  /// Create a copy of RoleItem
  /// with the given fields replaced by the non-null parameter values.
  @pragma('vm:prefer-inline')
  @override
  $Res call({
    Object? id = null,
    Object? name = null,
    Object? code = null,
    Object? description = freezed,
    Object? status = freezed,
    Object? createdAt = freezed,
    Object? updatedAt = freezed,
  }) {
    return _then(_$RoleItemImpl(
      id: null == id
          ? _value.id
          : id // ignore: cast_nullable_to_non_nullable
              as int,
      name: null == name
          ? _value.name
          : name // ignore: cast_nullable_to_non_nullable
              as String,
      code: null == code
          ? _value.code
          : code // ignore: cast_nullable_to_non_nullable
              as String,
      description: freezed == description
          ? _value.description
          : description // ignore: cast_nullable_to_non_nullable
              as String?,
      status: freezed == status
          ? _value.status
          : status // ignore: cast_nullable_to_non_nullable
              as int?,
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
class _$RoleItemImpl implements _RoleItem {
  const _$RoleItemImpl(
      {@JsonKey(fromJson: parseInt) required this.id,
      required this.name,
      required this.code,
      @JsonKey(name: 'description') this.description,
      @JsonKey(fromJson: parseIntNullable) this.status,
      @JsonKey(name: 'created_at', fromJson: parseStringNullable)
      this.createdAt,
      @JsonKey(name: 'updated_at', fromJson: parseStringNullable)
      this.updatedAt});

  factory _$RoleItemImpl.fromJson(Map<String, dynamic> json) =>
      _$$RoleItemImplFromJson(json);

  @override
  @JsonKey(fromJson: parseInt)
  final int id;
  @override
  final String name;
  @override
  final String code;
  @override
  @JsonKey(name: 'description')
  final String? description;
  @override
  @JsonKey(fromJson: parseIntNullable)
  final int? status;
// 1启用 0禁用
  @override
  @JsonKey(name: 'created_at', fromJson: parseStringNullable)
  final String? createdAt;
  @override
  @JsonKey(name: 'updated_at', fromJson: parseStringNullable)
  final String? updatedAt;

  @override
  String toString() {
    return 'RoleItem(id: $id, name: $name, code: $code, description: $description, status: $status, createdAt: $createdAt, updatedAt: $updatedAt)';
  }

  @override
  bool operator ==(Object other) {
    return identical(this, other) ||
        (other.runtimeType == runtimeType &&
            other is _$RoleItemImpl &&
            (identical(other.id, id) || other.id == id) &&
            (identical(other.name, name) || other.name == name) &&
            (identical(other.code, code) || other.code == code) &&
            (identical(other.description, description) ||
                other.description == description) &&
            (identical(other.status, status) || other.status == status) &&
            (identical(other.createdAt, createdAt) ||
                other.createdAt == createdAt) &&
            (identical(other.updatedAt, updatedAt) ||
                other.updatedAt == updatedAt));
  }

  @JsonKey(includeFromJson: false, includeToJson: false)
  @override
  int get hashCode => Object.hash(
      runtimeType, id, name, code, description, status, createdAt, updatedAt);

  /// Create a copy of RoleItem
  /// with the given fields replaced by the non-null parameter values.
  @JsonKey(includeFromJson: false, includeToJson: false)
  @override
  @pragma('vm:prefer-inline')
  _$$RoleItemImplCopyWith<_$RoleItemImpl> get copyWith =>
      __$$RoleItemImplCopyWithImpl<_$RoleItemImpl>(this, _$identity);

  @override
  Map<String, dynamic> toJson() {
    return _$$RoleItemImplToJson(
      this,
    );
  }
}

abstract class _RoleItem implements RoleItem {
  const factory _RoleItem(
      {@JsonKey(fromJson: parseInt) required final int id,
      required final String name,
      required final String code,
      @JsonKey(name: 'description') final String? description,
      @JsonKey(fromJson: parseIntNullable) final int? status,
      @JsonKey(name: 'created_at', fromJson: parseStringNullable)
      final String? createdAt,
      @JsonKey(name: 'updated_at', fromJson: parseStringNullable)
      final String? updatedAt}) = _$RoleItemImpl;

  factory _RoleItem.fromJson(Map<String, dynamic> json) =
      _$RoleItemImpl.fromJson;

  @override
  @JsonKey(fromJson: parseInt)
  int get id;
  @override
  String get name;
  @override
  String get code;
  @override
  @JsonKey(name: 'description')
  String? get description;
  @override
  @JsonKey(fromJson: parseIntNullable)
  int? get status; // 1启用 0禁用
  @override
  @JsonKey(name: 'created_at', fromJson: parseStringNullable)
  String? get createdAt;
  @override
  @JsonKey(name: 'updated_at', fromJson: parseStringNullable)
  String? get updatedAt;

  /// Create a copy of RoleItem
  /// with the given fields replaced by the non-null parameter values.
  @override
  @JsonKey(includeFromJson: false, includeToJson: false)
  _$$RoleItemImplCopyWith<_$RoleItemImpl> get copyWith =>
      throw _privateConstructorUsedError;
}
