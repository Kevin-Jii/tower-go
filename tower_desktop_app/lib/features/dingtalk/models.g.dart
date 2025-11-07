// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'models.dart';

// **************************************************************************
// JsonSerializableGenerator
// **************************************************************************

_$DingTalkRobotImpl _$$DingTalkRobotImplFromJson(Map<String, dynamic> json) =>
    _$DingTalkRobotImpl(
      id: (json['id'] as num).toInt(),
      name: json['name'] as String? ?? '',
      botType: json['bot_type'] as String? ?? '',
      webhook: json['webhook'] as String? ?? '',
      secret: json['secret'] as String? ?? '',
      clientId: json['client_id'] as String? ?? '',
      clientSecret: json['client_secret'] as String? ?? '',
      agentId: _parseIntFromDynamic(json['agent_id']),
      storeId: (json['store_id'] as num?)?.toInt(),
      isEnabled: json['is_enabled'] as bool? ?? true,
      msgType: json['msg_type'] as String? ?? '',
      robotCode: json['robot_code'] as String? ?? '',
      remark: json['remark'] as String?,
      createdAt: json['created_at'] as String?,
      updatedAt: json['updated_at'] as String?,
    );

Map<String, dynamic> _$$DingTalkRobotImplToJson(_$DingTalkRobotImpl instance) =>
    <String, dynamic>{
      'id': instance.id,
      'name': instance.name,
      'bot_type': instance.botType,
      'webhook': instance.webhook,
      'secret': instance.secret,
      'client_id': instance.clientId,
      'client_secret': instance.clientSecret,
      'agent_id': instance.agentId,
      'store_id': instance.storeId,
      'is_enabled': instance.isEnabled,
      'msg_type': instance.msgType,
      'robot_code': instance.robotCode,
      'remark': instance.remark,
      'created_at': instance.createdAt,
      'updated_at': instance.updatedAt,
    };

_$CreateDingTalkRobotRequestImpl _$$CreateDingTalkRobotRequestImplFromJson(
        Map<String, dynamic> json) =>
    _$CreateDingTalkRobotRequestImpl(
      name: json['name'] as String,
      token: json['token'] as String,
      storeId: (json['store_id'] as num?)?.toInt(),
      status: (json['status'] as num?)?.toInt(),
      remark: json['remark'] as String?,
    );

Map<String, dynamic> _$$CreateDingTalkRobotRequestImplToJson(
        _$CreateDingTalkRobotRequestImpl instance) =>
    <String, dynamic>{
      'name': instance.name,
      'token': instance.token,
      'store_id': instance.storeId,
      'status': instance.status,
      'remark': instance.remark,
    };

_$UpdateDingTalkRobotRequestImpl _$$UpdateDingTalkRobotRequestImplFromJson(
        Map<String, dynamic> json) =>
    _$UpdateDingTalkRobotRequestImpl(
      name: json['name'] as String?,
      token: json['token'] as String?,
      storeId: (json['store_id'] as num?)?.toInt(),
      status: (json['status'] as num?)?.toInt(),
      remark: json['remark'] as String?,
    );

Map<String, dynamic> _$$UpdateDingTalkRobotRequestImplToJson(
        _$UpdateDingTalkRobotRequestImpl instance) =>
    <String, dynamic>{
      'name': instance.name,
      'token': instance.token,
      'store_id': instance.storeId,
      'status': instance.status,
      'remark': instance.remark,
    };
