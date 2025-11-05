/// 通用后端响应解析工具，减少各 API 重复样板代码。
/// 约定：后端常见返回 { code, message, data } 或直接 data(list/map)。
library response_utils;

import 'dart:convert';

class ResponseUtils {
  /// 抽取最内层数据：如果是 Map 且含有 data 键，则返回其值，否则返回原始对象。
  /// 如果响应是标准格式 {code, message, meta} 且没有 data 字段，返回 null。
  static dynamic extractData(dynamic raw) {
    // 原始为字符串：尝试 JSON 解析（忽略异常，保留原始值）
    if (raw is String) {
      final trimmed = raw.trim();
      if (trimmed.isNotEmpty &&
          (trimmed.startsWith('{') || trimmed.startsWith('['))) {
        try {
          raw = jsonDecode(trimmed);
        } catch (_) {
          // 保留原样，调用方可自行处理
        }
      }
    }
    if (raw is Map) {
      // 如果有 data 键，返回其值
      if (raw.containsKey('data')) return raw['data'];
      // 如果是标准响应格式但没有 data 字段（如 {code, message, meta}），返回 null
      if (raw.containsKey('code') && raw.containsKey('message')) {
        return null;
      }
    }
    return raw;
  }

  /// 确保转换为 Map 列表。支持：
  /// - 直接 List
  /// - 包裹 Map 再取 listKey 对应的 List
  /// - 其他情况返回空列表
  static List<Map<String, dynamic>> toMapList(dynamic raw, {String? listKey}) {
    // 先确保已尝试抽取内部 data 与字符串解析
    raw = extractData(raw);
    dynamic data = raw;
    if (listKey != null && raw is Map) {
      final inner = raw[listKey];
      if (inner is List) {
        data = inner;
      } else {
        return const [];
      }
    }
    if (data is List) {
      return data
          .where((e) => e is Map)
          .map((e) => Map<String, dynamic>.from(e as Map))
          .toList();
    }
    return const [];
  }

  /// 直接组合：先 extractData 再 toMapList。
  static List<Map<String, dynamic>> payloadToList(dynamic raw,
      {String? listKey}) {
    return toMapList(raw, listKey: listKey);
  }
}
