import 'package:freezed_annotation/freezed_annotation.dart';

part 'models.freezed.dart';
part 'models.g.dart';

// 辅助函数来安全地解析各种类型的数值
int? _parseIntFromDynamic(dynamic value) {
  if (value == null) return null;
  if (value is int) return value;
  if (value is String) return int.tryParse(value);
  if (value is num) return value.toInt();
  return null;
}

@freezed
class DingTalkRobot with _$DingTalkRobot {
  const factory DingTalkRobot({
    required int id,
    @Default('') String name,
    @JsonKey(name: 'bot_type') @Default('') String botType,
    @Default('') String webhook,
    @Default('') String secret,
    @JsonKey(name: 'client_id') @Default('') String clientId,
    @JsonKey(name: 'client_secret') @Default('') String clientSecret,
    @JsonKey(name: 'agent_id', fromJson: _parseIntFromDynamic) int? agentId,
    @JsonKey(name: 'store_id') int? storeId,
    @JsonKey(name: 'is_enabled') @Default(true) bool isEnabled,
    @JsonKey(name: 'msg_type') @Default('') String msgType,
    @JsonKey(name: 'robot_code') @Default('') String robotCode,
    String? remark,
    @JsonKey(name: 'created_at') String? createdAt,
    @JsonKey(name: 'updated_at') String? updatedAt,
  }) = _DingTalkRobot;

  factory DingTalkRobot.fromJson(Map<String, dynamic> json) =>
      _$DingTalkRobotFromJson(json);
}

// 扩展方法提供兼容性
extension DingTalkRobotX on DingTalkRobot {
  // 为了兼容旧的 UI 代码，提供 token 属性
  String get token => webhook.isNotEmpty ? webhook : clientId;

  // 提供 status 属性（1:启用 0:禁用）
  int get status => isEnabled ? 1 : 0;
}

@freezed
class CreateDingTalkRobotRequest with _$CreateDingTalkRobotRequest {
  const factory CreateDingTalkRobotRequest({
    required String name,
    required String token,
    @JsonKey(name: 'store_id') int? storeId,
    int? status,
    String? remark,
  }) = _CreateDingTalkRobotRequest;

  factory CreateDingTalkRobotRequest.fromJson(Map<String, dynamic> json) =>
      _$CreateDingTalkRobotRequestFromJson(json);
}

@freezed
class UpdateDingTalkRobotRequest with _$UpdateDingTalkRobotRequest {
  const factory UpdateDingTalkRobotRequest({
    String? name,
    String? token,
    @JsonKey(name: 'store_id') int? storeId,
    int? status,
    String? remark,
  }) = _UpdateDingTalkRobotRequest;

  factory UpdateDingTalkRobotRequest.fromJson(Map<String, dynamic> json) =>
      _$UpdateDingTalkRobotRequestFromJson(json);
}
