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

DingTalkRobot _$DingTalkRobotFromJson(Map<String, dynamic> json) {
  return _DingTalkRobot.fromJson(json);
}

/// @nodoc
mixin _$DingTalkRobot {
  int get id => throw _privateConstructorUsedError;
  String get name => throw _privateConstructorUsedError;
  @JsonKey(name: 'bot_type')
  String get botType => throw _privateConstructorUsedError;
  String get webhook => throw _privateConstructorUsedError;
  String get secret => throw _privateConstructorUsedError;
  @JsonKey(name: 'client_id')
  String get clientId => throw _privateConstructorUsedError;
  @JsonKey(name: 'client_secret')
  String get clientSecret => throw _privateConstructorUsedError;
  @JsonKey(name: 'agent_id', fromJson: _parseIntFromDynamic)
  int? get agentId => throw _privateConstructorUsedError;
  @JsonKey(name: 'store_id')
  int? get storeId => throw _privateConstructorUsedError;
  @JsonKey(name: 'is_enabled')
  bool get isEnabled => throw _privateConstructorUsedError;
  @JsonKey(name: 'msg_type')
  String get msgType => throw _privateConstructorUsedError;
  @JsonKey(name: 'robot_code')
  String get robotCode => throw _privateConstructorUsedError;
  String? get remark => throw _privateConstructorUsedError;
  @JsonKey(name: 'created_at')
  String? get createdAt => throw _privateConstructorUsedError;
  @JsonKey(name: 'updated_at')
  String? get updatedAt => throw _privateConstructorUsedError;

  /// Serializes this DingTalkRobot to a JSON map.
  Map<String, dynamic> toJson() => throw _privateConstructorUsedError;

  /// Create a copy of DingTalkRobot
  /// with the given fields replaced by the non-null parameter values.
  @JsonKey(includeFromJson: false, includeToJson: false)
  $DingTalkRobotCopyWith<DingTalkRobot> get copyWith =>
      throw _privateConstructorUsedError;
}

/// @nodoc
abstract class $DingTalkRobotCopyWith<$Res> {
  factory $DingTalkRobotCopyWith(
          DingTalkRobot value, $Res Function(DingTalkRobot) then) =
      _$DingTalkRobotCopyWithImpl<$Res, DingTalkRobot>;
  @useResult
  $Res call(
      {int id,
      String name,
      @JsonKey(name: 'bot_type') String botType,
      String webhook,
      String secret,
      @JsonKey(name: 'client_id') String clientId,
      @JsonKey(name: 'client_secret') String clientSecret,
      @JsonKey(name: 'agent_id', fromJson: _parseIntFromDynamic) int? agentId,
      @JsonKey(name: 'store_id') int? storeId,
      @JsonKey(name: 'is_enabled') bool isEnabled,
      @JsonKey(name: 'msg_type') String msgType,
      @JsonKey(name: 'robot_code') String robotCode,
      String? remark,
      @JsonKey(name: 'created_at') String? createdAt,
      @JsonKey(name: 'updated_at') String? updatedAt});
}

/// @nodoc
class _$DingTalkRobotCopyWithImpl<$Res, $Val extends DingTalkRobot>
    implements $DingTalkRobotCopyWith<$Res> {
  _$DingTalkRobotCopyWithImpl(this._value, this._then);

  // ignore: unused_field
  final $Val _value;
  // ignore: unused_field
  final $Res Function($Val) _then;

  /// Create a copy of DingTalkRobot
  /// with the given fields replaced by the non-null parameter values.
  @pragma('vm:prefer-inline')
  @override
  $Res call({
    Object? id = null,
    Object? name = null,
    Object? botType = null,
    Object? webhook = null,
    Object? secret = null,
    Object? clientId = null,
    Object? clientSecret = null,
    Object? agentId = freezed,
    Object? storeId = freezed,
    Object? isEnabled = null,
    Object? msgType = null,
    Object? robotCode = null,
    Object? remark = freezed,
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
      botType: null == botType
          ? _value.botType
          : botType // ignore: cast_nullable_to_non_nullable
              as String,
      webhook: null == webhook
          ? _value.webhook
          : webhook // ignore: cast_nullable_to_non_nullable
              as String,
      secret: null == secret
          ? _value.secret
          : secret // ignore: cast_nullable_to_non_nullable
              as String,
      clientId: null == clientId
          ? _value.clientId
          : clientId // ignore: cast_nullable_to_non_nullable
              as String,
      clientSecret: null == clientSecret
          ? _value.clientSecret
          : clientSecret // ignore: cast_nullable_to_non_nullable
              as String,
      agentId: freezed == agentId
          ? _value.agentId
          : agentId // ignore: cast_nullable_to_non_nullable
              as int?,
      storeId: freezed == storeId
          ? _value.storeId
          : storeId // ignore: cast_nullable_to_non_nullable
              as int?,
      isEnabled: null == isEnabled
          ? _value.isEnabled
          : isEnabled // ignore: cast_nullable_to_non_nullable
              as bool,
      msgType: null == msgType
          ? _value.msgType
          : msgType // ignore: cast_nullable_to_non_nullable
              as String,
      robotCode: null == robotCode
          ? _value.robotCode
          : robotCode // ignore: cast_nullable_to_non_nullable
              as String,
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
abstract class _$$DingTalkRobotImplCopyWith<$Res>
    implements $DingTalkRobotCopyWith<$Res> {
  factory _$$DingTalkRobotImplCopyWith(
          _$DingTalkRobotImpl value, $Res Function(_$DingTalkRobotImpl) then) =
      __$$DingTalkRobotImplCopyWithImpl<$Res>;
  @override
  @useResult
  $Res call(
      {int id,
      String name,
      @JsonKey(name: 'bot_type') String botType,
      String webhook,
      String secret,
      @JsonKey(name: 'client_id') String clientId,
      @JsonKey(name: 'client_secret') String clientSecret,
      @JsonKey(name: 'agent_id', fromJson: _parseIntFromDynamic) int? agentId,
      @JsonKey(name: 'store_id') int? storeId,
      @JsonKey(name: 'is_enabled') bool isEnabled,
      @JsonKey(name: 'msg_type') String msgType,
      @JsonKey(name: 'robot_code') String robotCode,
      String? remark,
      @JsonKey(name: 'created_at') String? createdAt,
      @JsonKey(name: 'updated_at') String? updatedAt});
}

/// @nodoc
class __$$DingTalkRobotImplCopyWithImpl<$Res>
    extends _$DingTalkRobotCopyWithImpl<$Res, _$DingTalkRobotImpl>
    implements _$$DingTalkRobotImplCopyWith<$Res> {
  __$$DingTalkRobotImplCopyWithImpl(
      _$DingTalkRobotImpl _value, $Res Function(_$DingTalkRobotImpl) _then)
      : super(_value, _then);

  /// Create a copy of DingTalkRobot
  /// with the given fields replaced by the non-null parameter values.
  @pragma('vm:prefer-inline')
  @override
  $Res call({
    Object? id = null,
    Object? name = null,
    Object? botType = null,
    Object? webhook = null,
    Object? secret = null,
    Object? clientId = null,
    Object? clientSecret = null,
    Object? agentId = freezed,
    Object? storeId = freezed,
    Object? isEnabled = null,
    Object? msgType = null,
    Object? robotCode = null,
    Object? remark = freezed,
    Object? createdAt = freezed,
    Object? updatedAt = freezed,
  }) {
    return _then(_$DingTalkRobotImpl(
      id: null == id
          ? _value.id
          : id // ignore: cast_nullable_to_non_nullable
              as int,
      name: null == name
          ? _value.name
          : name // ignore: cast_nullable_to_non_nullable
              as String,
      botType: null == botType
          ? _value.botType
          : botType // ignore: cast_nullable_to_non_nullable
              as String,
      webhook: null == webhook
          ? _value.webhook
          : webhook // ignore: cast_nullable_to_non_nullable
              as String,
      secret: null == secret
          ? _value.secret
          : secret // ignore: cast_nullable_to_non_nullable
              as String,
      clientId: null == clientId
          ? _value.clientId
          : clientId // ignore: cast_nullable_to_non_nullable
              as String,
      clientSecret: null == clientSecret
          ? _value.clientSecret
          : clientSecret // ignore: cast_nullable_to_non_nullable
              as String,
      agentId: freezed == agentId
          ? _value.agentId
          : agentId // ignore: cast_nullable_to_non_nullable
              as int?,
      storeId: freezed == storeId
          ? _value.storeId
          : storeId // ignore: cast_nullable_to_non_nullable
              as int?,
      isEnabled: null == isEnabled
          ? _value.isEnabled
          : isEnabled // ignore: cast_nullable_to_non_nullable
              as bool,
      msgType: null == msgType
          ? _value.msgType
          : msgType // ignore: cast_nullable_to_non_nullable
              as String,
      robotCode: null == robotCode
          ? _value.robotCode
          : robotCode // ignore: cast_nullable_to_non_nullable
              as String,
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
class _$DingTalkRobotImpl implements _DingTalkRobot {
  const _$DingTalkRobotImpl(
      {required this.id,
      this.name = '',
      @JsonKey(name: 'bot_type') this.botType = '',
      this.webhook = '',
      this.secret = '',
      @JsonKey(name: 'client_id') this.clientId = '',
      @JsonKey(name: 'client_secret') this.clientSecret = '',
      @JsonKey(name: 'agent_id', fromJson: _parseIntFromDynamic) this.agentId,
      @JsonKey(name: 'store_id') this.storeId,
      @JsonKey(name: 'is_enabled') this.isEnabled = true,
      @JsonKey(name: 'msg_type') this.msgType = '',
      @JsonKey(name: 'robot_code') this.robotCode = '',
      this.remark,
      @JsonKey(name: 'created_at') this.createdAt,
      @JsonKey(name: 'updated_at') this.updatedAt});

  factory _$DingTalkRobotImpl.fromJson(Map<String, dynamic> json) =>
      _$$DingTalkRobotImplFromJson(json);

  @override
  final int id;
  @override
  @JsonKey()
  final String name;
  @override
  @JsonKey(name: 'bot_type')
  final String botType;
  @override
  @JsonKey()
  final String webhook;
  @override
  @JsonKey()
  final String secret;
  @override
  @JsonKey(name: 'client_id')
  final String clientId;
  @override
  @JsonKey(name: 'client_secret')
  final String clientSecret;
  @override
  @JsonKey(name: 'agent_id', fromJson: _parseIntFromDynamic)
  final int? agentId;
  @override
  @JsonKey(name: 'store_id')
  final int? storeId;
  @override
  @JsonKey(name: 'is_enabled')
  final bool isEnabled;
  @override
  @JsonKey(name: 'msg_type')
  final String msgType;
  @override
  @JsonKey(name: 'robot_code')
  final String robotCode;
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
    return 'DingTalkRobot(id: $id, name: $name, botType: $botType, webhook: $webhook, secret: $secret, clientId: $clientId, clientSecret: $clientSecret, agentId: $agentId, storeId: $storeId, isEnabled: $isEnabled, msgType: $msgType, robotCode: $robotCode, remark: $remark, createdAt: $createdAt, updatedAt: $updatedAt)';
  }

  @override
  bool operator ==(Object other) {
    return identical(this, other) ||
        (other.runtimeType == runtimeType &&
            other is _$DingTalkRobotImpl &&
            (identical(other.id, id) || other.id == id) &&
            (identical(other.name, name) || other.name == name) &&
            (identical(other.botType, botType) || other.botType == botType) &&
            (identical(other.webhook, webhook) || other.webhook == webhook) &&
            (identical(other.secret, secret) || other.secret == secret) &&
            (identical(other.clientId, clientId) ||
                other.clientId == clientId) &&
            (identical(other.clientSecret, clientSecret) ||
                other.clientSecret == clientSecret) &&
            (identical(other.agentId, agentId) || other.agentId == agentId) &&
            (identical(other.storeId, storeId) || other.storeId == storeId) &&
            (identical(other.isEnabled, isEnabled) ||
                other.isEnabled == isEnabled) &&
            (identical(other.msgType, msgType) || other.msgType == msgType) &&
            (identical(other.robotCode, robotCode) ||
                other.robotCode == robotCode) &&
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
      name,
      botType,
      webhook,
      secret,
      clientId,
      clientSecret,
      agentId,
      storeId,
      isEnabled,
      msgType,
      robotCode,
      remark,
      createdAt,
      updatedAt);

  /// Create a copy of DingTalkRobot
  /// with the given fields replaced by the non-null parameter values.
  @JsonKey(includeFromJson: false, includeToJson: false)
  @override
  @pragma('vm:prefer-inline')
  _$$DingTalkRobotImplCopyWith<_$DingTalkRobotImpl> get copyWith =>
      __$$DingTalkRobotImplCopyWithImpl<_$DingTalkRobotImpl>(this, _$identity);

  @override
  Map<String, dynamic> toJson() {
    return _$$DingTalkRobotImplToJson(
      this,
    );
  }
}

abstract class _DingTalkRobot implements DingTalkRobot {
  const factory _DingTalkRobot(
          {required final int id,
          final String name,
          @JsonKey(name: 'bot_type') final String botType,
          final String webhook,
          final String secret,
          @JsonKey(name: 'client_id') final String clientId,
          @JsonKey(name: 'client_secret') final String clientSecret,
          @JsonKey(name: 'agent_id', fromJson: _parseIntFromDynamic)
          final int? agentId,
          @JsonKey(name: 'store_id') final int? storeId,
          @JsonKey(name: 'is_enabled') final bool isEnabled,
          @JsonKey(name: 'msg_type') final String msgType,
          @JsonKey(name: 'robot_code') final String robotCode,
          final String? remark,
          @JsonKey(name: 'created_at') final String? createdAt,
          @JsonKey(name: 'updated_at') final String? updatedAt}) =
      _$DingTalkRobotImpl;

  factory _DingTalkRobot.fromJson(Map<String, dynamic> json) =
      _$DingTalkRobotImpl.fromJson;

  @override
  int get id;
  @override
  String get name;
  @override
  @JsonKey(name: 'bot_type')
  String get botType;
  @override
  String get webhook;
  @override
  String get secret;
  @override
  @JsonKey(name: 'client_id')
  String get clientId;
  @override
  @JsonKey(name: 'client_secret')
  String get clientSecret;
  @override
  @JsonKey(name: 'agent_id', fromJson: _parseIntFromDynamic)
  int? get agentId;
  @override
  @JsonKey(name: 'store_id')
  int? get storeId;
  @override
  @JsonKey(name: 'is_enabled')
  bool get isEnabled;
  @override
  @JsonKey(name: 'msg_type')
  String get msgType;
  @override
  @JsonKey(name: 'robot_code')
  String get robotCode;
  @override
  String? get remark;
  @override
  @JsonKey(name: 'created_at')
  String? get createdAt;
  @override
  @JsonKey(name: 'updated_at')
  String? get updatedAt;

  /// Create a copy of DingTalkRobot
  /// with the given fields replaced by the non-null parameter values.
  @override
  @JsonKey(includeFromJson: false, includeToJson: false)
  _$$DingTalkRobotImplCopyWith<_$DingTalkRobotImpl> get copyWith =>
      throw _privateConstructorUsedError;
}

CreateDingTalkRobotRequest _$CreateDingTalkRobotRequestFromJson(
    Map<String, dynamic> json) {
  return _CreateDingTalkRobotRequest.fromJson(json);
}

/// @nodoc
mixin _$CreateDingTalkRobotRequest {
  String get name => throw _privateConstructorUsedError;
  String get token => throw _privateConstructorUsedError;
  @JsonKey(name: 'store_id')
  int? get storeId => throw _privateConstructorUsedError;
  int? get status => throw _privateConstructorUsedError;
  String? get remark => throw _privateConstructorUsedError;

  /// Serializes this CreateDingTalkRobotRequest to a JSON map.
  Map<String, dynamic> toJson() => throw _privateConstructorUsedError;

  /// Create a copy of CreateDingTalkRobotRequest
  /// with the given fields replaced by the non-null parameter values.
  @JsonKey(includeFromJson: false, includeToJson: false)
  $CreateDingTalkRobotRequestCopyWith<CreateDingTalkRobotRequest>
      get copyWith => throw _privateConstructorUsedError;
}

/// @nodoc
abstract class $CreateDingTalkRobotRequestCopyWith<$Res> {
  factory $CreateDingTalkRobotRequestCopyWith(CreateDingTalkRobotRequest value,
          $Res Function(CreateDingTalkRobotRequest) then) =
      _$CreateDingTalkRobotRequestCopyWithImpl<$Res,
          CreateDingTalkRobotRequest>;
  @useResult
  $Res call(
      {String name,
      String token,
      @JsonKey(name: 'store_id') int? storeId,
      int? status,
      String? remark});
}

/// @nodoc
class _$CreateDingTalkRobotRequestCopyWithImpl<$Res,
        $Val extends CreateDingTalkRobotRequest>
    implements $CreateDingTalkRobotRequestCopyWith<$Res> {
  _$CreateDingTalkRobotRequestCopyWithImpl(this._value, this._then);

  // ignore: unused_field
  final $Val _value;
  // ignore: unused_field
  final $Res Function($Val) _then;

  /// Create a copy of CreateDingTalkRobotRequest
  /// with the given fields replaced by the non-null parameter values.
  @pragma('vm:prefer-inline')
  @override
  $Res call({
    Object? name = null,
    Object? token = null,
    Object? storeId = freezed,
    Object? status = freezed,
    Object? remark = freezed,
  }) {
    return _then(_value.copyWith(
      name: null == name
          ? _value.name
          : name // ignore: cast_nullable_to_non_nullable
              as String,
      token: null == token
          ? _value.token
          : token // ignore: cast_nullable_to_non_nullable
              as String,
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
abstract class _$$CreateDingTalkRobotRequestImplCopyWith<$Res>
    implements $CreateDingTalkRobotRequestCopyWith<$Res> {
  factory _$$CreateDingTalkRobotRequestImplCopyWith(
          _$CreateDingTalkRobotRequestImpl value,
          $Res Function(_$CreateDingTalkRobotRequestImpl) then) =
      __$$CreateDingTalkRobotRequestImplCopyWithImpl<$Res>;
  @override
  @useResult
  $Res call(
      {String name,
      String token,
      @JsonKey(name: 'store_id') int? storeId,
      int? status,
      String? remark});
}

/// @nodoc
class __$$CreateDingTalkRobotRequestImplCopyWithImpl<$Res>
    extends _$CreateDingTalkRobotRequestCopyWithImpl<$Res,
        _$CreateDingTalkRobotRequestImpl>
    implements _$$CreateDingTalkRobotRequestImplCopyWith<$Res> {
  __$$CreateDingTalkRobotRequestImplCopyWithImpl(
      _$CreateDingTalkRobotRequestImpl _value,
      $Res Function(_$CreateDingTalkRobotRequestImpl) _then)
      : super(_value, _then);

  /// Create a copy of CreateDingTalkRobotRequest
  /// with the given fields replaced by the non-null parameter values.
  @pragma('vm:prefer-inline')
  @override
  $Res call({
    Object? name = null,
    Object? token = null,
    Object? storeId = freezed,
    Object? status = freezed,
    Object? remark = freezed,
  }) {
    return _then(_$CreateDingTalkRobotRequestImpl(
      name: null == name
          ? _value.name
          : name // ignore: cast_nullable_to_non_nullable
              as String,
      token: null == token
          ? _value.token
          : token // ignore: cast_nullable_to_non_nullable
              as String,
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
class _$CreateDingTalkRobotRequestImpl implements _CreateDingTalkRobotRequest {
  const _$CreateDingTalkRobotRequestImpl(
      {required this.name,
      required this.token,
      @JsonKey(name: 'store_id') this.storeId,
      this.status,
      this.remark});

  factory _$CreateDingTalkRobotRequestImpl.fromJson(
          Map<String, dynamic> json) =>
      _$$CreateDingTalkRobotRequestImplFromJson(json);

  @override
  final String name;
  @override
  final String token;
  @override
  @JsonKey(name: 'store_id')
  final int? storeId;
  @override
  final int? status;
  @override
  final String? remark;

  @override
  String toString() {
    return 'CreateDingTalkRobotRequest(name: $name, token: $token, storeId: $storeId, status: $status, remark: $remark)';
  }

  @override
  bool operator ==(Object other) {
    return identical(this, other) ||
        (other.runtimeType == runtimeType &&
            other is _$CreateDingTalkRobotRequestImpl &&
            (identical(other.name, name) || other.name == name) &&
            (identical(other.token, token) || other.token == token) &&
            (identical(other.storeId, storeId) || other.storeId == storeId) &&
            (identical(other.status, status) || other.status == status) &&
            (identical(other.remark, remark) || other.remark == remark));
  }

  @JsonKey(includeFromJson: false, includeToJson: false)
  @override
  int get hashCode =>
      Object.hash(runtimeType, name, token, storeId, status, remark);

  /// Create a copy of CreateDingTalkRobotRequest
  /// with the given fields replaced by the non-null parameter values.
  @JsonKey(includeFromJson: false, includeToJson: false)
  @override
  @pragma('vm:prefer-inline')
  _$$CreateDingTalkRobotRequestImplCopyWith<_$CreateDingTalkRobotRequestImpl>
      get copyWith => __$$CreateDingTalkRobotRequestImplCopyWithImpl<
          _$CreateDingTalkRobotRequestImpl>(this, _$identity);

  @override
  Map<String, dynamic> toJson() {
    return _$$CreateDingTalkRobotRequestImplToJson(
      this,
    );
  }
}

abstract class _CreateDingTalkRobotRequest
    implements CreateDingTalkRobotRequest {
  const factory _CreateDingTalkRobotRequest(
      {required final String name,
      required final String token,
      @JsonKey(name: 'store_id') final int? storeId,
      final int? status,
      final String? remark}) = _$CreateDingTalkRobotRequestImpl;

  factory _CreateDingTalkRobotRequest.fromJson(Map<String, dynamic> json) =
      _$CreateDingTalkRobotRequestImpl.fromJson;

  @override
  String get name;
  @override
  String get token;
  @override
  @JsonKey(name: 'store_id')
  int? get storeId;
  @override
  int? get status;
  @override
  String? get remark;

  /// Create a copy of CreateDingTalkRobotRequest
  /// with the given fields replaced by the non-null parameter values.
  @override
  @JsonKey(includeFromJson: false, includeToJson: false)
  _$$CreateDingTalkRobotRequestImplCopyWith<_$CreateDingTalkRobotRequestImpl>
      get copyWith => throw _privateConstructorUsedError;
}

UpdateDingTalkRobotRequest _$UpdateDingTalkRobotRequestFromJson(
    Map<String, dynamic> json) {
  return _UpdateDingTalkRobotRequest.fromJson(json);
}

/// @nodoc
mixin _$UpdateDingTalkRobotRequest {
  String? get name => throw _privateConstructorUsedError;
  String? get token => throw _privateConstructorUsedError;
  @JsonKey(name: 'store_id')
  int? get storeId => throw _privateConstructorUsedError;
  int? get status => throw _privateConstructorUsedError;
  String? get remark => throw _privateConstructorUsedError;

  /// Serializes this UpdateDingTalkRobotRequest to a JSON map.
  Map<String, dynamic> toJson() => throw _privateConstructorUsedError;

  /// Create a copy of UpdateDingTalkRobotRequest
  /// with the given fields replaced by the non-null parameter values.
  @JsonKey(includeFromJson: false, includeToJson: false)
  $UpdateDingTalkRobotRequestCopyWith<UpdateDingTalkRobotRequest>
      get copyWith => throw _privateConstructorUsedError;
}

/// @nodoc
abstract class $UpdateDingTalkRobotRequestCopyWith<$Res> {
  factory $UpdateDingTalkRobotRequestCopyWith(UpdateDingTalkRobotRequest value,
          $Res Function(UpdateDingTalkRobotRequest) then) =
      _$UpdateDingTalkRobotRequestCopyWithImpl<$Res,
          UpdateDingTalkRobotRequest>;
  @useResult
  $Res call(
      {String? name,
      String? token,
      @JsonKey(name: 'store_id') int? storeId,
      int? status,
      String? remark});
}

/// @nodoc
class _$UpdateDingTalkRobotRequestCopyWithImpl<$Res,
        $Val extends UpdateDingTalkRobotRequest>
    implements $UpdateDingTalkRobotRequestCopyWith<$Res> {
  _$UpdateDingTalkRobotRequestCopyWithImpl(this._value, this._then);

  // ignore: unused_field
  final $Val _value;
  // ignore: unused_field
  final $Res Function($Val) _then;

  /// Create a copy of UpdateDingTalkRobotRequest
  /// with the given fields replaced by the non-null parameter values.
  @pragma('vm:prefer-inline')
  @override
  $Res call({
    Object? name = freezed,
    Object? token = freezed,
    Object? storeId = freezed,
    Object? status = freezed,
    Object? remark = freezed,
  }) {
    return _then(_value.copyWith(
      name: freezed == name
          ? _value.name
          : name // ignore: cast_nullable_to_non_nullable
              as String?,
      token: freezed == token
          ? _value.token
          : token // ignore: cast_nullable_to_non_nullable
              as String?,
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
abstract class _$$UpdateDingTalkRobotRequestImplCopyWith<$Res>
    implements $UpdateDingTalkRobotRequestCopyWith<$Res> {
  factory _$$UpdateDingTalkRobotRequestImplCopyWith(
          _$UpdateDingTalkRobotRequestImpl value,
          $Res Function(_$UpdateDingTalkRobotRequestImpl) then) =
      __$$UpdateDingTalkRobotRequestImplCopyWithImpl<$Res>;
  @override
  @useResult
  $Res call(
      {String? name,
      String? token,
      @JsonKey(name: 'store_id') int? storeId,
      int? status,
      String? remark});
}

/// @nodoc
class __$$UpdateDingTalkRobotRequestImplCopyWithImpl<$Res>
    extends _$UpdateDingTalkRobotRequestCopyWithImpl<$Res,
        _$UpdateDingTalkRobotRequestImpl>
    implements _$$UpdateDingTalkRobotRequestImplCopyWith<$Res> {
  __$$UpdateDingTalkRobotRequestImplCopyWithImpl(
      _$UpdateDingTalkRobotRequestImpl _value,
      $Res Function(_$UpdateDingTalkRobotRequestImpl) _then)
      : super(_value, _then);

  /// Create a copy of UpdateDingTalkRobotRequest
  /// with the given fields replaced by the non-null parameter values.
  @pragma('vm:prefer-inline')
  @override
  $Res call({
    Object? name = freezed,
    Object? token = freezed,
    Object? storeId = freezed,
    Object? status = freezed,
    Object? remark = freezed,
  }) {
    return _then(_$UpdateDingTalkRobotRequestImpl(
      name: freezed == name
          ? _value.name
          : name // ignore: cast_nullable_to_non_nullable
              as String?,
      token: freezed == token
          ? _value.token
          : token // ignore: cast_nullable_to_non_nullable
              as String?,
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
class _$UpdateDingTalkRobotRequestImpl implements _UpdateDingTalkRobotRequest {
  const _$UpdateDingTalkRobotRequestImpl(
      {this.name,
      this.token,
      @JsonKey(name: 'store_id') this.storeId,
      this.status,
      this.remark});

  factory _$UpdateDingTalkRobotRequestImpl.fromJson(
          Map<String, dynamic> json) =>
      _$$UpdateDingTalkRobotRequestImplFromJson(json);

  @override
  final String? name;
  @override
  final String? token;
  @override
  @JsonKey(name: 'store_id')
  final int? storeId;
  @override
  final int? status;
  @override
  final String? remark;

  @override
  String toString() {
    return 'UpdateDingTalkRobotRequest(name: $name, token: $token, storeId: $storeId, status: $status, remark: $remark)';
  }

  @override
  bool operator ==(Object other) {
    return identical(this, other) ||
        (other.runtimeType == runtimeType &&
            other is _$UpdateDingTalkRobotRequestImpl &&
            (identical(other.name, name) || other.name == name) &&
            (identical(other.token, token) || other.token == token) &&
            (identical(other.storeId, storeId) || other.storeId == storeId) &&
            (identical(other.status, status) || other.status == status) &&
            (identical(other.remark, remark) || other.remark == remark));
  }

  @JsonKey(includeFromJson: false, includeToJson: false)
  @override
  int get hashCode =>
      Object.hash(runtimeType, name, token, storeId, status, remark);

  /// Create a copy of UpdateDingTalkRobotRequest
  /// with the given fields replaced by the non-null parameter values.
  @JsonKey(includeFromJson: false, includeToJson: false)
  @override
  @pragma('vm:prefer-inline')
  _$$UpdateDingTalkRobotRequestImplCopyWith<_$UpdateDingTalkRobotRequestImpl>
      get copyWith => __$$UpdateDingTalkRobotRequestImplCopyWithImpl<
          _$UpdateDingTalkRobotRequestImpl>(this, _$identity);

  @override
  Map<String, dynamic> toJson() {
    return _$$UpdateDingTalkRobotRequestImplToJson(
      this,
    );
  }
}

abstract class _UpdateDingTalkRobotRequest
    implements UpdateDingTalkRobotRequest {
  const factory _UpdateDingTalkRobotRequest(
      {final String? name,
      final String? token,
      @JsonKey(name: 'store_id') final int? storeId,
      final int? status,
      final String? remark}) = _$UpdateDingTalkRobotRequestImpl;

  factory _UpdateDingTalkRobotRequest.fromJson(Map<String, dynamic> json) =
      _$UpdateDingTalkRobotRequestImpl.fromJson;

  @override
  String? get name;
  @override
  String? get token;
  @override
  @JsonKey(name: 'store_id')
  int? get storeId;
  @override
  int? get status;
  @override
  String? get remark;

  /// Create a copy of UpdateDingTalkRobotRequest
  /// with the given fields replaced by the non-null parameter values.
  @override
  @JsonKey(includeFromJson: false, includeToJson: false)
  _$$UpdateDingTalkRobotRequestImplCopyWith<_$UpdateDingTalkRobotRequestImpl>
      get copyWith => throw _privateConstructorUsedError;
}
