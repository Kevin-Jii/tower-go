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

User _$UserFromJson(Map<String, dynamic> json) {
  return _User.fromJson(json);
}

/// @nodoc
mixin _$User {
  int get id => throw _privateConstructorUsedError;
  String get username => throw _privateConstructorUsedError;
  String get phone => throw _privateConstructorUsedError;
  String? get email => throw _privateConstructorUsedError;
  String? get nickname => throw _privateConstructorUsedError;
  String? get avatar => throw _privateConstructorUsedError;
  int? get gender => throw _privateConstructorUsedError; // 1=男 2=女
  @JsonKey(name: 'role_id')
  int? get roleId => throw _privateConstructorUsedError;
  @JsonKey(name: 'role_name')
  String? get roleName => throw _privateConstructorUsedError;
  @JsonKey(name: 'store_id')
  int? get storeId => throw _privateConstructorUsedError;
  @JsonKey(name: 'store_name')
  String? get storeName => throw _privateConstructorUsedError;
  int? get status => throw _privateConstructorUsedError; // 0=禁用 1=启用
  String? get remark => throw _privateConstructorUsedError;
  @JsonKey(name: 'created_at')
  String? get createdAt => throw _privateConstructorUsedError;
  @JsonKey(name: 'updated_at')
  String? get updatedAt => throw _privateConstructorUsedError;

  /// Serializes this User to a JSON map.
  Map<String, dynamic> toJson() => throw _privateConstructorUsedError;

  /// Create a copy of User
  /// with the given fields replaced by the non-null parameter values.
  @JsonKey(includeFromJson: false, includeToJson: false)
  $UserCopyWith<User> get copyWith => throw _privateConstructorUsedError;
}

/// @nodoc
abstract class $UserCopyWith<$Res> {
  factory $UserCopyWith(User value, $Res Function(User) then) =
      _$UserCopyWithImpl<$Res, User>;
  @useResult
  $Res call(
      {int id,
      String username,
      String phone,
      String? email,
      String? nickname,
      String? avatar,
      int? gender,
      @JsonKey(name: 'role_id') int? roleId,
      @JsonKey(name: 'role_name') String? roleName,
      @JsonKey(name: 'store_id') int? storeId,
      @JsonKey(name: 'store_name') String? storeName,
      int? status,
      String? remark,
      @JsonKey(name: 'created_at') String? createdAt,
      @JsonKey(name: 'updated_at') String? updatedAt});
}

/// @nodoc
class _$UserCopyWithImpl<$Res, $Val extends User>
    implements $UserCopyWith<$Res> {
  _$UserCopyWithImpl(this._value, this._then);

  // ignore: unused_field
  final $Val _value;
  // ignore: unused_field
  final $Res Function($Val) _then;

  /// Create a copy of User
  /// with the given fields replaced by the non-null parameter values.
  @pragma('vm:prefer-inline')
  @override
  $Res call({
    Object? id = null,
    Object? username = null,
    Object? phone = null,
    Object? email = freezed,
    Object? nickname = freezed,
    Object? avatar = freezed,
    Object? gender = freezed,
    Object? roleId = freezed,
    Object? roleName = freezed,
    Object? storeId = freezed,
    Object? storeName = freezed,
    Object? status = freezed,
    Object? remark = freezed,
    Object? createdAt = freezed,
    Object? updatedAt = freezed,
  }) {
    return _then(_value.copyWith(
      id: null == id
          ? _value.id
          : id // ignore: cast_nullable_to_non_nullable
              as int,
      username: null == username
          ? _value.username
          : username // ignore: cast_nullable_to_non_nullable
              as String,
      phone: null == phone
          ? _value.phone
          : phone // ignore: cast_nullable_to_non_nullable
              as String,
      email: freezed == email
          ? _value.email
          : email // ignore: cast_nullable_to_non_nullable
              as String?,
      nickname: freezed == nickname
          ? _value.nickname
          : nickname // ignore: cast_nullable_to_non_nullable
              as String?,
      avatar: freezed == avatar
          ? _value.avatar
          : avatar // ignore: cast_nullable_to_non_nullable
              as String?,
      gender: freezed == gender
          ? _value.gender
          : gender // ignore: cast_nullable_to_non_nullable
              as int?,
      roleId: freezed == roleId
          ? _value.roleId
          : roleId // ignore: cast_nullable_to_non_nullable
              as int?,
      roleName: freezed == roleName
          ? _value.roleName
          : roleName // ignore: cast_nullable_to_non_nullable
              as String?,
      storeId: freezed == storeId
          ? _value.storeId
          : storeId // ignore: cast_nullable_to_non_nullable
              as int?,
      storeName: freezed == storeName
          ? _value.storeName
          : storeName // ignore: cast_nullable_to_non_nullable
              as String?,
      status: freezed == status
          ? _value.status
          : status // ignore: cast_nullable_to_non_nullable
              as int?,
      remark: freezed == remark
          ? _value.remark
          : remark // ignore: cast_nullable_to_non_nullable
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
abstract class _$$UserImplCopyWith<$Res> implements $UserCopyWith<$Res> {
  factory _$$UserImplCopyWith(
          _$UserImpl value, $Res Function(_$UserImpl) then) =
      __$$UserImplCopyWithImpl<$Res>;
  @override
  @useResult
  $Res call(
      {int id,
      String username,
      String phone,
      String? email,
      String? nickname,
      String? avatar,
      int? gender,
      @JsonKey(name: 'role_id') int? roleId,
      @JsonKey(name: 'role_name') String? roleName,
      @JsonKey(name: 'store_id') int? storeId,
      @JsonKey(name: 'store_name') String? storeName,
      int? status,
      String? remark,
      @JsonKey(name: 'created_at') String? createdAt,
      @JsonKey(name: 'updated_at') String? updatedAt});
}

/// @nodoc
class __$$UserImplCopyWithImpl<$Res>
    extends _$UserCopyWithImpl<$Res, _$UserImpl>
    implements _$$UserImplCopyWith<$Res> {
  __$$UserImplCopyWithImpl(_$UserImpl _value, $Res Function(_$UserImpl) _then)
      : super(_value, _then);

  /// Create a copy of User
  /// with the given fields replaced by the non-null parameter values.
  @pragma('vm:prefer-inline')
  @override
  $Res call({
    Object? id = null,
    Object? username = null,
    Object? phone = null,
    Object? email = freezed,
    Object? nickname = freezed,
    Object? avatar = freezed,
    Object? gender = freezed,
    Object? roleId = freezed,
    Object? roleName = freezed,
    Object? storeId = freezed,
    Object? storeName = freezed,
    Object? status = freezed,
    Object? remark = freezed,
    Object? createdAt = freezed,
    Object? updatedAt = freezed,
  }) {
    return _then(_$UserImpl(
      id: null == id
          ? _value.id
          : id // ignore: cast_nullable_to_non_nullable
              as int,
      username: null == username
          ? _value.username
          : username // ignore: cast_nullable_to_non_nullable
              as String,
      phone: null == phone
          ? _value.phone
          : phone // ignore: cast_nullable_to_non_nullable
              as String,
      email: freezed == email
          ? _value.email
          : email // ignore: cast_nullable_to_non_nullable
              as String?,
      nickname: freezed == nickname
          ? _value.nickname
          : nickname // ignore: cast_nullable_to_non_nullable
              as String?,
      avatar: freezed == avatar
          ? _value.avatar
          : avatar // ignore: cast_nullable_to_non_nullable
              as String?,
      gender: freezed == gender
          ? _value.gender
          : gender // ignore: cast_nullable_to_non_nullable
              as int?,
      roleId: freezed == roleId
          ? _value.roleId
          : roleId // ignore: cast_nullable_to_non_nullable
              as int?,
      roleName: freezed == roleName
          ? _value.roleName
          : roleName // ignore: cast_nullable_to_non_nullable
              as String?,
      storeId: freezed == storeId
          ? _value.storeId
          : storeId // ignore: cast_nullable_to_non_nullable
              as int?,
      storeName: freezed == storeName
          ? _value.storeName
          : storeName // ignore: cast_nullable_to_non_nullable
              as String?,
      status: freezed == status
          ? _value.status
          : status // ignore: cast_nullable_to_non_nullable
              as int?,
      remark: freezed == remark
          ? _value.remark
          : remark // ignore: cast_nullable_to_non_nullable
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
class _$UserImpl implements _User {
  const _$UserImpl(
      {required this.id,
      required this.username,
      required this.phone,
      this.email,
      this.nickname,
      this.avatar,
      this.gender,
      @JsonKey(name: 'role_id') this.roleId,
      @JsonKey(name: 'role_name') this.roleName,
      @JsonKey(name: 'store_id') this.storeId,
      @JsonKey(name: 'store_name') this.storeName,
      this.status,
      this.remark,
      @JsonKey(name: 'created_at') this.createdAt,
      @JsonKey(name: 'updated_at') this.updatedAt});

  factory _$UserImpl.fromJson(Map<String, dynamic> json) =>
      _$$UserImplFromJson(json);

  @override
  final int id;
  @override
  final String username;
  @override
  final String phone;
  @override
  final String? email;
  @override
  final String? nickname;
  @override
  final String? avatar;
  @override
  final int? gender;
// 1=男 2=女
  @override
  @JsonKey(name: 'role_id')
  final int? roleId;
  @override
  @JsonKey(name: 'role_name')
  final String? roleName;
  @override
  @JsonKey(name: 'store_id')
  final int? storeId;
  @override
  @JsonKey(name: 'store_name')
  final String? storeName;
  @override
  final int? status;
// 0=禁用 1=启用
  @override
  final String? remark;
  @override
  @JsonKey(name: 'created_at')
  final String? createdAt;
  @override
  @JsonKey(name: 'updated_at')
  final String? updatedAt;

  @override
  String toString() {
    return 'User(id: $id, username: $username, phone: $phone, email: $email, nickname: $nickname, avatar: $avatar, gender: $gender, roleId: $roleId, roleName: $roleName, storeId: $storeId, storeName: $storeName, status: $status, remark: $remark, createdAt: $createdAt, updatedAt: $updatedAt)';
  }

  @override
  bool operator ==(Object other) {
    return identical(this, other) ||
        (other.runtimeType == runtimeType &&
            other is _$UserImpl &&
            (identical(other.id, id) || other.id == id) &&
            (identical(other.username, username) ||
                other.username == username) &&
            (identical(other.phone, phone) || other.phone == phone) &&
            (identical(other.email, email) || other.email == email) &&
            (identical(other.nickname, nickname) ||
                other.nickname == nickname) &&
            (identical(other.avatar, avatar) || other.avatar == avatar) &&
            (identical(other.gender, gender) || other.gender == gender) &&
            (identical(other.roleId, roleId) || other.roleId == roleId) &&
            (identical(other.roleName, roleName) ||
                other.roleName == roleName) &&
            (identical(other.storeId, storeId) || other.storeId == storeId) &&
            (identical(other.storeName, storeName) ||
                other.storeName == storeName) &&
            (identical(other.status, status) || other.status == status) &&
            (identical(other.remark, remark) || other.remark == remark) &&
            (identical(other.createdAt, createdAt) ||
                other.createdAt == createdAt) &&
            (identical(other.updatedAt, updatedAt) ||
                other.updatedAt == updatedAt));
  }

  @JsonKey(includeFromJson: false, includeToJson: false)
  @override
  int get hashCode => Object.hash(
      runtimeType,
      id,
      username,
      phone,
      email,
      nickname,
      avatar,
      gender,
      roleId,
      roleName,
      storeId,
      storeName,
      status,
      remark,
      createdAt,
      updatedAt);

  /// Create a copy of User
  /// with the given fields replaced by the non-null parameter values.
  @JsonKey(includeFromJson: false, includeToJson: false)
  @override
  @pragma('vm:prefer-inline')
  _$$UserImplCopyWith<_$UserImpl> get copyWith =>
      __$$UserImplCopyWithImpl<_$UserImpl>(this, _$identity);

  @override
  Map<String, dynamic> toJson() {
    return _$$UserImplToJson(
      this,
    );
  }
}

abstract class _User implements User {
  const factory _User(
      {required final int id,
      required final String username,
      required final String phone,
      final String? email,
      final String? nickname,
      final String? avatar,
      final int? gender,
      @JsonKey(name: 'role_id') final int? roleId,
      @JsonKey(name: 'role_name') final String? roleName,
      @JsonKey(name: 'store_id') final int? storeId,
      @JsonKey(name: 'store_name') final String? storeName,
      final int? status,
      final String? remark,
      @JsonKey(name: 'created_at') final String? createdAt,
      @JsonKey(name: 'updated_at') final String? updatedAt}) = _$UserImpl;

  factory _User.fromJson(Map<String, dynamic> json) = _$UserImpl.fromJson;

  @override
  int get id;
  @override
  String get username;
  @override
  String get phone;
  @override
  String? get email;
  @override
  String? get nickname;
  @override
  String? get avatar;
  @override
  int? get gender; // 1=男 2=女
  @override
  @JsonKey(name: 'role_id')
  int? get roleId;
  @override
  @JsonKey(name: 'role_name')
  String? get roleName;
  @override
  @JsonKey(name: 'store_id')
  int? get storeId;
  @override
  @JsonKey(name: 'store_name')
  String? get storeName;
  @override
  int? get status; // 0=禁用 1=启用
  @override
  String? get remark;
  @override
  @JsonKey(name: 'created_at')
  String? get createdAt;
  @override
  @JsonKey(name: 'updated_at')
  String? get updatedAt;

  /// Create a copy of User
  /// with the given fields replaced by the non-null parameter values.
  @override
  @JsonKey(includeFromJson: false, includeToJson: false)
  _$$UserImplCopyWith<_$UserImpl> get copyWith =>
      throw _privateConstructorUsedError;
}

UserListResponse _$UserListResponseFromJson(Map<String, dynamic> json) {
  return _UserListResponse.fromJson(json);
}

/// @nodoc
mixin _$UserListResponse {
  List<User> get list => throw _privateConstructorUsedError;
  int get total => throw _privateConstructorUsedError;
  int get page => throw _privateConstructorUsedError;
  @JsonKey(name: 'page_size')
  int get pageSize => throw _privateConstructorUsedError;

  /// Serializes this UserListResponse to a JSON map.
  Map<String, dynamic> toJson() => throw _privateConstructorUsedError;

  /// Create a copy of UserListResponse
  /// with the given fields replaced by the non-null parameter values.
  @JsonKey(includeFromJson: false, includeToJson: false)
  $UserListResponseCopyWith<UserListResponse> get copyWith =>
      throw _privateConstructorUsedError;
}

/// @nodoc
abstract class $UserListResponseCopyWith<$Res> {
  factory $UserListResponseCopyWith(
          UserListResponse value, $Res Function(UserListResponse) then) =
      _$UserListResponseCopyWithImpl<$Res, UserListResponse>;
  @useResult
  $Res call(
      {List<User> list,
      int total,
      int page,
      @JsonKey(name: 'page_size') int pageSize});
}

/// @nodoc
class _$UserListResponseCopyWithImpl<$Res, $Val extends UserListResponse>
    implements $UserListResponseCopyWith<$Res> {
  _$UserListResponseCopyWithImpl(this._value, this._then);

  // ignore: unused_field
  final $Val _value;
  // ignore: unused_field
  final $Res Function($Val) _then;

  /// Create a copy of UserListResponse
  /// with the given fields replaced by the non-null parameter values.
  @pragma('vm:prefer-inline')
  @override
  $Res call({
    Object? list = null,
    Object? total = null,
    Object? page = null,
    Object? pageSize = null,
  }) {
    return _then(_value.copyWith(
      list: null == list
          ? _value.list
          : list // ignore: cast_nullable_to_non_nullable
              as List<User>,
      total: null == total
          ? _value.total
          : total // ignore: cast_nullable_to_non_nullable
              as int,
      page: null == page
          ? _value.page
          : page // ignore: cast_nullable_to_non_nullable
              as int,
      pageSize: null == pageSize
          ? _value.pageSize
          : pageSize // ignore: cast_nullable_to_non_nullable
              as int,
    ) as $Val);
  }
}

/// @nodoc
abstract class _$$UserListResponseImplCopyWith<$Res>
    implements $UserListResponseCopyWith<$Res> {
  factory _$$UserListResponseImplCopyWith(_$UserListResponseImpl value,
          $Res Function(_$UserListResponseImpl) then) =
      __$$UserListResponseImplCopyWithImpl<$Res>;
  @override
  @useResult
  $Res call(
      {List<User> list,
      int total,
      int page,
      @JsonKey(name: 'page_size') int pageSize});
}

/// @nodoc
class __$$UserListResponseImplCopyWithImpl<$Res>
    extends _$UserListResponseCopyWithImpl<$Res, _$UserListResponseImpl>
    implements _$$UserListResponseImplCopyWith<$Res> {
  __$$UserListResponseImplCopyWithImpl(_$UserListResponseImpl _value,
      $Res Function(_$UserListResponseImpl) _then)
      : super(_value, _then);

  /// Create a copy of UserListResponse
  /// with the given fields replaced by the non-null parameter values.
  @pragma('vm:prefer-inline')
  @override
  $Res call({
    Object? list = null,
    Object? total = null,
    Object? page = null,
    Object? pageSize = null,
  }) {
    return _then(_$UserListResponseImpl(
      list: null == list
          ? _value._list
          : list // ignore: cast_nullable_to_non_nullable
              as List<User>,
      total: null == total
          ? _value.total
          : total // ignore: cast_nullable_to_non_nullable
              as int,
      page: null == page
          ? _value.page
          : page // ignore: cast_nullable_to_non_nullable
              as int,
      pageSize: null == pageSize
          ? _value.pageSize
          : pageSize // ignore: cast_nullable_to_non_nullable
              as int,
    ));
  }
}

/// @nodoc
@JsonSerializable()
class _$UserListResponseImpl implements _UserListResponse {
  const _$UserListResponseImpl(
      {required final List<User> list,
      required this.total,
      this.page = 1,
      @JsonKey(name: 'page_size') this.pageSize = 10})
      : _list = list;

  factory _$UserListResponseImpl.fromJson(Map<String, dynamic> json) =>
      _$$UserListResponseImplFromJson(json);

  final List<User> _list;
  @override
  List<User> get list {
    if (_list is EqualUnmodifiableListView) return _list;
    // ignore: implicit_dynamic_type
    return EqualUnmodifiableListView(_list);
  }

  @override
  final int total;
  @override
  @JsonKey()
  final int page;
  @override
  @JsonKey(name: 'page_size')
  final int pageSize;

  @override
  String toString() {
    return 'UserListResponse(list: $list, total: $total, page: $page, pageSize: $pageSize)';
  }

  @override
  bool operator ==(Object other) {
    return identical(this, other) ||
        (other.runtimeType == runtimeType &&
            other is _$UserListResponseImpl &&
            const DeepCollectionEquality().equals(other._list, _list) &&
            (identical(other.total, total) || other.total == total) &&
            (identical(other.page, page) || other.page == page) &&
            (identical(other.pageSize, pageSize) ||
                other.pageSize == pageSize));
  }

  @JsonKey(includeFromJson: false, includeToJson: false)
  @override
  int get hashCode => Object.hash(runtimeType,
      const DeepCollectionEquality().hash(_list), total, page, pageSize);

  /// Create a copy of UserListResponse
  /// with the given fields replaced by the non-null parameter values.
  @JsonKey(includeFromJson: false, includeToJson: false)
  @override
  @pragma('vm:prefer-inline')
  _$$UserListResponseImplCopyWith<_$UserListResponseImpl> get copyWith =>
      __$$UserListResponseImplCopyWithImpl<_$UserListResponseImpl>(
          this, _$identity);

  @override
  Map<String, dynamic> toJson() {
    return _$$UserListResponseImplToJson(
      this,
    );
  }
}

abstract class _UserListResponse implements UserListResponse {
  const factory _UserListResponse(
      {required final List<User> list,
      required final int total,
      final int page,
      @JsonKey(name: 'page_size') final int pageSize}) = _$UserListResponseImpl;

  factory _UserListResponse.fromJson(Map<String, dynamic> json) =
      _$UserListResponseImpl.fromJson;

  @override
  List<User> get list;
  @override
  int get total;
  @override
  int get page;
  @override
  @JsonKey(name: 'page_size')
  int get pageSize;

  /// Create a copy of UserListResponse
  /// with the given fields replaced by the non-null parameter values.
  @override
  @JsonKey(includeFromJson: false, includeToJson: false)
  _$$UserListResponseImplCopyWith<_$UserListResponseImpl> get copyWith =>
      throw _privateConstructorUsedError;
}

CreateUserRequest _$CreateUserRequestFromJson(Map<String, dynamic> json) {
  return _CreateUserRequest.fromJson(json);
}

/// @nodoc
mixin _$CreateUserRequest {
  String get username => throw _privateConstructorUsedError;
  String get phone => throw _privateConstructorUsedError;
  String get password => throw _privateConstructorUsedError;
  String? get email => throw _privateConstructorUsedError;
  String? get nickname => throw _privateConstructorUsedError;
  int? get gender => throw _privateConstructorUsedError; // 1=男 2=女
  @JsonKey(name: 'role_id')
  int? get roleId => throw _privateConstructorUsedError;
  @JsonKey(name: 'store_id')
  int? get storeId => throw _privateConstructorUsedError;
  int? get status => throw _privateConstructorUsedError;
  String? get remark => throw _privateConstructorUsedError;

  /// Serializes this CreateUserRequest to a JSON map.
  Map<String, dynamic> toJson() => throw _privateConstructorUsedError;

  /// Create a copy of CreateUserRequest
  /// with the given fields replaced by the non-null parameter values.
  @JsonKey(includeFromJson: false, includeToJson: false)
  $CreateUserRequestCopyWith<CreateUserRequest> get copyWith =>
      throw _privateConstructorUsedError;
}

/// @nodoc
abstract class $CreateUserRequestCopyWith<$Res> {
  factory $CreateUserRequestCopyWith(
          CreateUserRequest value, $Res Function(CreateUserRequest) then) =
      _$CreateUserRequestCopyWithImpl<$Res, CreateUserRequest>;
  @useResult
  $Res call(
      {String username,
      String phone,
      String password,
      String? email,
      String? nickname,
      int? gender,
      @JsonKey(name: 'role_id') int? roleId,
      @JsonKey(name: 'store_id') int? storeId,
      int? status,
      String? remark});
}

/// @nodoc
class _$CreateUserRequestCopyWithImpl<$Res, $Val extends CreateUserRequest>
    implements $CreateUserRequestCopyWith<$Res> {
  _$CreateUserRequestCopyWithImpl(this._value, this._then);

  // ignore: unused_field
  final $Val _value;
  // ignore: unused_field
  final $Res Function($Val) _then;

  /// Create a copy of CreateUserRequest
  /// with the given fields replaced by the non-null parameter values.
  @pragma('vm:prefer-inline')
  @override
  $Res call({
    Object? username = null,
    Object? phone = null,
    Object? password = null,
    Object? email = freezed,
    Object? nickname = freezed,
    Object? gender = freezed,
    Object? roleId = freezed,
    Object? storeId = freezed,
    Object? status = freezed,
    Object? remark = freezed,
  }) {
    return _then(_value.copyWith(
      username: null == username
          ? _value.username
          : username // ignore: cast_nullable_to_non_nullable
              as String,
      phone: null == phone
          ? _value.phone
          : phone // ignore: cast_nullable_to_non_nullable
              as String,
      password: null == password
          ? _value.password
          : password // ignore: cast_nullable_to_non_nullable
              as String,
      email: freezed == email
          ? _value.email
          : email // ignore: cast_nullable_to_non_nullable
              as String?,
      nickname: freezed == nickname
          ? _value.nickname
          : nickname // ignore: cast_nullable_to_non_nullable
              as String?,
      gender: freezed == gender
          ? _value.gender
          : gender // ignore: cast_nullable_to_non_nullable
              as int?,
      roleId: freezed == roleId
          ? _value.roleId
          : roleId // ignore: cast_nullable_to_non_nullable
              as int?,
      storeId: freezed == storeId
          ? _value.storeId
          : storeId // ignore: cast_nullable_to_non_nullable
              as int?,
      status: freezed == status
          ? _value.status
          : status // ignore: cast_nullable_to_non_nullable
              as int?,
      remark: freezed == remark
          ? _value.remark
          : remark // ignore: cast_nullable_to_non_nullable
              as String?,
    ) as $Val);
  }
}

/// @nodoc
abstract class _$$CreateUserRequestImplCopyWith<$Res>
    implements $CreateUserRequestCopyWith<$Res> {
  factory _$$CreateUserRequestImplCopyWith(_$CreateUserRequestImpl value,
          $Res Function(_$CreateUserRequestImpl) then) =
      __$$CreateUserRequestImplCopyWithImpl<$Res>;
  @override
  @useResult
  $Res call(
      {String username,
      String phone,
      String password,
      String? email,
      String? nickname,
      int? gender,
      @JsonKey(name: 'role_id') int? roleId,
      @JsonKey(name: 'store_id') int? storeId,
      int? status,
      String? remark});
}

/// @nodoc
class __$$CreateUserRequestImplCopyWithImpl<$Res>
    extends _$CreateUserRequestCopyWithImpl<$Res, _$CreateUserRequestImpl>
    implements _$$CreateUserRequestImplCopyWith<$Res> {
  __$$CreateUserRequestImplCopyWithImpl(_$CreateUserRequestImpl _value,
      $Res Function(_$CreateUserRequestImpl) _then)
      : super(_value, _then);

  /// Create a copy of CreateUserRequest
  /// with the given fields replaced by the non-null parameter values.
  @pragma('vm:prefer-inline')
  @override
  $Res call({
    Object? username = null,
    Object? phone = null,
    Object? password = null,
    Object? email = freezed,
    Object? nickname = freezed,
    Object? gender = freezed,
    Object? roleId = freezed,
    Object? storeId = freezed,
    Object? status = freezed,
    Object? remark = freezed,
  }) {
    return _then(_$CreateUserRequestImpl(
      username: null == username
          ? _value.username
          : username // ignore: cast_nullable_to_non_nullable
              as String,
      phone: null == phone
          ? _value.phone
          : phone // ignore: cast_nullable_to_non_nullable
              as String,
      password: null == password
          ? _value.password
          : password // ignore: cast_nullable_to_non_nullable
              as String,
      email: freezed == email
          ? _value.email
          : email // ignore: cast_nullable_to_non_nullable
              as String?,
      nickname: freezed == nickname
          ? _value.nickname
          : nickname // ignore: cast_nullable_to_non_nullable
              as String?,
      gender: freezed == gender
          ? _value.gender
          : gender // ignore: cast_nullable_to_non_nullable
              as int?,
      roleId: freezed == roleId
          ? _value.roleId
          : roleId // ignore: cast_nullable_to_non_nullable
              as int?,
      storeId: freezed == storeId
          ? _value.storeId
          : storeId // ignore: cast_nullable_to_non_nullable
              as int?,
      status: freezed == status
          ? _value.status
          : status // ignore: cast_nullable_to_non_nullable
              as int?,
      remark: freezed == remark
          ? _value.remark
          : remark // ignore: cast_nullable_to_non_nullable
              as String?,
    ));
  }
}

/// @nodoc
@JsonSerializable()
class _$CreateUserRequestImpl implements _CreateUserRequest {
  const _$CreateUserRequestImpl(
      {required this.username,
      required this.phone,
      required this.password,
      this.email,
      this.nickname,
      this.gender,
      @JsonKey(name: 'role_id') this.roleId,
      @JsonKey(name: 'store_id') this.storeId,
      this.status,
      this.remark});

  factory _$CreateUserRequestImpl.fromJson(Map<String, dynamic> json) =>
      _$$CreateUserRequestImplFromJson(json);

  @override
  final String username;
  @override
  final String phone;
  @override
  final String password;
  @override
  final String? email;
  @override
  final String? nickname;
  @override
  final int? gender;
// 1=男 2=女
  @override
  @JsonKey(name: 'role_id')
  final int? roleId;
  @override
  @JsonKey(name: 'store_id')
  final int? storeId;
  @override
  final int? status;
  @override
  final String? remark;

  @override
  String toString() {
    return 'CreateUserRequest(username: $username, phone: $phone, password: $password, email: $email, nickname: $nickname, gender: $gender, roleId: $roleId, storeId: $storeId, status: $status, remark: $remark)';
  }

  @override
  bool operator ==(Object other) {
    return identical(this, other) ||
        (other.runtimeType == runtimeType &&
            other is _$CreateUserRequestImpl &&
            (identical(other.username, username) ||
                other.username == username) &&
            (identical(other.phone, phone) || other.phone == phone) &&
            (identical(other.password, password) ||
                other.password == password) &&
            (identical(other.email, email) || other.email == email) &&
            (identical(other.nickname, nickname) ||
                other.nickname == nickname) &&
            (identical(other.gender, gender) || other.gender == gender) &&
            (identical(other.roleId, roleId) || other.roleId == roleId) &&
            (identical(other.storeId, storeId) || other.storeId == storeId) &&
            (identical(other.status, status) || other.status == status) &&
            (identical(other.remark, remark) || other.remark == remark));
  }

  @JsonKey(includeFromJson: false, includeToJson: false)
  @override
  int get hashCode => Object.hash(runtimeType, username, phone, password, email,
      nickname, gender, roleId, storeId, status, remark);

  /// Create a copy of CreateUserRequest
  /// with the given fields replaced by the non-null parameter values.
  @JsonKey(includeFromJson: false, includeToJson: false)
  @override
  @pragma('vm:prefer-inline')
  _$$CreateUserRequestImplCopyWith<_$CreateUserRequestImpl> get copyWith =>
      __$$CreateUserRequestImplCopyWithImpl<_$CreateUserRequestImpl>(
          this, _$identity);

  @override
  Map<String, dynamic> toJson() {
    return _$$CreateUserRequestImplToJson(
      this,
    );
  }
}

abstract class _CreateUserRequest implements CreateUserRequest {
  const factory _CreateUserRequest(
      {required final String username,
      required final String phone,
      required final String password,
      final String? email,
      final String? nickname,
      final int? gender,
      @JsonKey(name: 'role_id') final int? roleId,
      @JsonKey(name: 'store_id') final int? storeId,
      final int? status,
      final String? remark}) = _$CreateUserRequestImpl;

  factory _CreateUserRequest.fromJson(Map<String, dynamic> json) =
      _$CreateUserRequestImpl.fromJson;

  @override
  String get username;
  @override
  String get phone;
  @override
  String get password;
  @override
  String? get email;
  @override
  String? get nickname;
  @override
  int? get gender; // 1=男 2=女
  @override
  @JsonKey(name: 'role_id')
  int? get roleId;
  @override
  @JsonKey(name: 'store_id')
  int? get storeId;
  @override
  int? get status;
  @override
  String? get remark;

  /// Create a copy of CreateUserRequest
  /// with the given fields replaced by the non-null parameter values.
  @override
  @JsonKey(includeFromJson: false, includeToJson: false)
  _$$CreateUserRequestImplCopyWith<_$CreateUserRequestImpl> get copyWith =>
      throw _privateConstructorUsedError;
}

UpdateUserRequest _$UpdateUserRequestFromJson(Map<String, dynamic> json) {
  return _UpdateUserRequest.fromJson(json);
}

/// @nodoc
mixin _$UpdateUserRequest {
  String? get phone => throw _privateConstructorUsedError;
  String? get email => throw _privateConstructorUsedError;
  String? get nickname => throw _privateConstructorUsedError;
  String? get password => throw _privateConstructorUsedError; // 可选,为空则不修改
  @JsonKey(name: 'role_id')
  int? get roleId => throw _privateConstructorUsedError;
  @JsonKey(name: 'store_id')
  int? get storeId => throw _privateConstructorUsedError;
  int? get status => throw _privateConstructorUsedError;
  String? get remark => throw _privateConstructorUsedError;

  /// Serializes this UpdateUserRequest to a JSON map.
  Map<String, dynamic> toJson() => throw _privateConstructorUsedError;

  /// Create a copy of UpdateUserRequest
  /// with the given fields replaced by the non-null parameter values.
  @JsonKey(includeFromJson: false, includeToJson: false)
  $UpdateUserRequestCopyWith<UpdateUserRequest> get copyWith =>
      throw _privateConstructorUsedError;
}

/// @nodoc
abstract class $UpdateUserRequestCopyWith<$Res> {
  factory $UpdateUserRequestCopyWith(
          UpdateUserRequest value, $Res Function(UpdateUserRequest) then) =
      _$UpdateUserRequestCopyWithImpl<$Res, UpdateUserRequest>;
  @useResult
  $Res call(
      {String? phone,
      String? email,
      String? nickname,
      String? password,
      @JsonKey(name: 'role_id') int? roleId,
      @JsonKey(name: 'store_id') int? storeId,
      int? status,
      String? remark});
}

/// @nodoc
class _$UpdateUserRequestCopyWithImpl<$Res, $Val extends UpdateUserRequest>
    implements $UpdateUserRequestCopyWith<$Res> {
  _$UpdateUserRequestCopyWithImpl(this._value, this._then);

  // ignore: unused_field
  final $Val _value;
  // ignore: unused_field
  final $Res Function($Val) _then;

  /// Create a copy of UpdateUserRequest
  /// with the given fields replaced by the non-null parameter values.
  @pragma('vm:prefer-inline')
  @override
  $Res call({
    Object? phone = freezed,
    Object? email = freezed,
    Object? nickname = freezed,
    Object? password = freezed,
    Object? roleId = freezed,
    Object? storeId = freezed,
    Object? status = freezed,
    Object? remark = freezed,
  }) {
    return _then(_value.copyWith(
      phone: freezed == phone
          ? _value.phone
          : phone // ignore: cast_nullable_to_non_nullable
              as String?,
      email: freezed == email
          ? _value.email
          : email // ignore: cast_nullable_to_non_nullable
              as String?,
      nickname: freezed == nickname
          ? _value.nickname
          : nickname // ignore: cast_nullable_to_non_nullable
              as String?,
      password: freezed == password
          ? _value.password
          : password // ignore: cast_nullable_to_non_nullable
              as String?,
      roleId: freezed == roleId
          ? _value.roleId
          : roleId // ignore: cast_nullable_to_non_nullable
              as int?,
      storeId: freezed == storeId
          ? _value.storeId
          : storeId // ignore: cast_nullable_to_non_nullable
              as int?,
      status: freezed == status
          ? _value.status
          : status // ignore: cast_nullable_to_non_nullable
              as int?,
      remark: freezed == remark
          ? _value.remark
          : remark // ignore: cast_nullable_to_non_nullable
              as String?,
    ) as $Val);
  }
}

/// @nodoc
abstract class _$$UpdateUserRequestImplCopyWith<$Res>
    implements $UpdateUserRequestCopyWith<$Res> {
  factory _$$UpdateUserRequestImplCopyWith(_$UpdateUserRequestImpl value,
          $Res Function(_$UpdateUserRequestImpl) then) =
      __$$UpdateUserRequestImplCopyWithImpl<$Res>;
  @override
  @useResult
  $Res call(
      {String? phone,
      String? email,
      String? nickname,
      String? password,
      @JsonKey(name: 'role_id') int? roleId,
      @JsonKey(name: 'store_id') int? storeId,
      int? status,
      String? remark});
}

/// @nodoc
class __$$UpdateUserRequestImplCopyWithImpl<$Res>
    extends _$UpdateUserRequestCopyWithImpl<$Res, _$UpdateUserRequestImpl>
    implements _$$UpdateUserRequestImplCopyWith<$Res> {
  __$$UpdateUserRequestImplCopyWithImpl(_$UpdateUserRequestImpl _value,
      $Res Function(_$UpdateUserRequestImpl) _then)
      : super(_value, _then);

  /// Create a copy of UpdateUserRequest
  /// with the given fields replaced by the non-null parameter values.
  @pragma('vm:prefer-inline')
  @override
  $Res call({
    Object? phone = freezed,
    Object? email = freezed,
    Object? nickname = freezed,
    Object? password = freezed,
    Object? roleId = freezed,
    Object? storeId = freezed,
    Object? status = freezed,
    Object? remark = freezed,
  }) {
    return _then(_$UpdateUserRequestImpl(
      phone: freezed == phone
          ? _value.phone
          : phone // ignore: cast_nullable_to_non_nullable
              as String?,
      email: freezed == email
          ? _value.email
          : email // ignore: cast_nullable_to_non_nullable
              as String?,
      nickname: freezed == nickname
          ? _value.nickname
          : nickname // ignore: cast_nullable_to_non_nullable
              as String?,
      password: freezed == password
          ? _value.password
          : password // ignore: cast_nullable_to_non_nullable
              as String?,
      roleId: freezed == roleId
          ? _value.roleId
          : roleId // ignore: cast_nullable_to_non_nullable
              as int?,
      storeId: freezed == storeId
          ? _value.storeId
          : storeId // ignore: cast_nullable_to_non_nullable
              as int?,
      status: freezed == status
          ? _value.status
          : status // ignore: cast_nullable_to_non_nullable
              as int?,
      remark: freezed == remark
          ? _value.remark
          : remark // ignore: cast_nullable_to_non_nullable
              as String?,
    ));
  }
}

/// @nodoc
@JsonSerializable()
class _$UpdateUserRequestImpl implements _UpdateUserRequest {
  const _$UpdateUserRequestImpl(
      {this.phone,
      this.email,
      this.nickname,
      this.password,
      @JsonKey(name: 'role_id') this.roleId,
      @JsonKey(name: 'store_id') this.storeId,
      this.status,
      this.remark});

  factory _$UpdateUserRequestImpl.fromJson(Map<String, dynamic> json) =>
      _$$UpdateUserRequestImplFromJson(json);

  @override
  final String? phone;
  @override
  final String? email;
  @override
  final String? nickname;
  @override
  final String? password;
// 可选,为空则不修改
  @override
  @JsonKey(name: 'role_id')
  final int? roleId;
  @override
  @JsonKey(name: 'store_id')
  final int? storeId;
  @override
  final int? status;
  @override
  final String? remark;

  @override
  String toString() {
    return 'UpdateUserRequest(phone: $phone, email: $email, nickname: $nickname, password: $password, roleId: $roleId, storeId: $storeId, status: $status, remark: $remark)';
  }

  @override
  bool operator ==(Object other) {
    return identical(this, other) ||
        (other.runtimeType == runtimeType &&
            other is _$UpdateUserRequestImpl &&
            (identical(other.phone, phone) || other.phone == phone) &&
            (identical(other.email, email) || other.email == email) &&
            (identical(other.nickname, nickname) ||
                other.nickname == nickname) &&
            (identical(other.password, password) ||
                other.password == password) &&
            (identical(other.roleId, roleId) || other.roleId == roleId) &&
            (identical(other.storeId, storeId) || other.storeId == storeId) &&
            (identical(other.status, status) || other.status == status) &&
            (identical(other.remark, remark) || other.remark == remark));
  }

  @JsonKey(includeFromJson: false, includeToJson: false)
  @override
  int get hashCode => Object.hash(runtimeType, phone, email, nickname, password,
      roleId, storeId, status, remark);

  /// Create a copy of UpdateUserRequest
  /// with the given fields replaced by the non-null parameter values.
  @JsonKey(includeFromJson: false, includeToJson: false)
  @override
  @pragma('vm:prefer-inline')
  _$$UpdateUserRequestImplCopyWith<_$UpdateUserRequestImpl> get copyWith =>
      __$$UpdateUserRequestImplCopyWithImpl<_$UpdateUserRequestImpl>(
          this, _$identity);

  @override
  Map<String, dynamic> toJson() {
    return _$$UpdateUserRequestImplToJson(
      this,
    );
  }
}

abstract class _UpdateUserRequest implements UpdateUserRequest {
  const factory _UpdateUserRequest(
      {final String? phone,
      final String? email,
      final String? nickname,
      final String? password,
      @JsonKey(name: 'role_id') final int? roleId,
      @JsonKey(name: 'store_id') final int? storeId,
      final int? status,
      final String? remark}) = _$UpdateUserRequestImpl;

  factory _UpdateUserRequest.fromJson(Map<String, dynamic> json) =
      _$UpdateUserRequestImpl.fromJson;

  @override
  String? get phone;
  @override
  String? get email;
  @override
  String? get nickname;
  @override
  String? get password; // 可选,为空则不修改
  @override
  @JsonKey(name: 'role_id')
  int? get roleId;
  @override
  @JsonKey(name: 'store_id')
  int? get storeId;
  @override
  int? get status;
  @override
  String? get remark;

  /// Create a copy of UpdateUserRequest
  /// with the given fields replaced by the non-null parameter values.
  @override
  @JsonKey(includeFromJson: false, includeToJson: false)
  _$$UpdateUserRequestImplCopyWith<_$UpdateUserRequestImpl> get copyWith =>
      throw _privateConstructorUsedError;
}
