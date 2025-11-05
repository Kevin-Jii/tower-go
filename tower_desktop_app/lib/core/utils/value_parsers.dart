/// 通用后端字段解析工具，兼容数字/字符串/空值等不规范返回。
/// 后续如果需要支持 bool/int 的更多变体，可在此集中维护，避免各模型重复实现。
library value_parsers;

int parseInt(dynamic v, {int defaultValue = 0}) {
  if (v == null) return defaultValue;
  if (v is int) return v;
  if (v is num) return v.toInt();
  if (v is String) {
    final t = v.trim();
    if (t.isEmpty) return defaultValue;
    return int.tryParse(t) ?? defaultValue;
  }
  return defaultValue;
}

int? parseIntNullable(dynamic v) {
  if (v == null) return null;
  if (v is int) return v;
  if (v is num) return v.toInt();
  if (v is String) {
    final t = v.trim();
    if (t.isEmpty) return null;
    return int.tryParse(t);
  }
  return null; // 保守返回 null 以避免抛异常
}

String? parseStringNullable(dynamic v) {
  if (v == null) return null;
  if (v is String) return v;
  if (v is num || v is bool) return v.toString();
  return v.toString();
}

double parseDouble(dynamic v, {double defaultValue = 0.0}) {
  if (v == null) return defaultValue;
  if (v is double) return v;
  if (v is num) return v.toDouble();
  if (v is String) {
    final t = v.trim();
    if (t.isEmpty) return defaultValue;
    return double.tryParse(t) ?? defaultValue;
  }
  return defaultValue;
}

double? parseDoubleNullable(dynamic v) {
  if (v == null) return null;
  if (v is double) return v;
  if (v is num) return v.toDouble();
  if (v is String) {
    final t = v.trim();
    if (t.isEmpty) return null;
    return double.tryParse(t);
  }
  return null;
}
