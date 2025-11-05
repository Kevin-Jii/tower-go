/// 时间解析与格式化工具
/// 兼容以下输入：
/// - ISO8601 字符串 (2025-11-05T12:34:56Z / 2025-11-05 12:34:56)
/// - 纯时间戳（秒或毫秒）数字或字符串
/// - 其它无法识别时返回 null
library time_utils;

DateTime? parseFlexibleDateTime(dynamic input) {
  if (input == null) return null;
  if (input is DateTime) return input.toLocal();
  if (input is int) {
    return _fromIntTimestamp(input);
  }
  if (input is num) {
    return _fromIntTimestamp(input.toInt());
  }
  if (input is String) {
    final v = input.trim();
    if (v.isEmpty) return null;
    // 尝试直接 DateTime.parse
    try {
      final dt = DateTime.parse(v.replaceAll(' ', 'T'));
      return dt.toLocal();
    } catch (_) {
      // 尝试数值解析
      final n = int.tryParse(v);
      if (n != null) {
        return _fromIntTimestamp(n);
      }
    }
  }
  return null;
}

DateTime? _fromIntTimestamp(int ts) {
  // 根据位数或大小判断：> 10^12 基本是毫秒；否则可能是秒。
  if (ts > 1000000000000) {
    return DateTime.fromMillisecondsSinceEpoch(ts).toLocal();
  }
  // 如果是秒级，转换为毫秒
  if (ts < 1000000000000 && ts > 1000000000) {
    return DateTime.fromMillisecondsSinceEpoch(ts * 1000).toLocal();
  }
  // 其它过小值不可信，返回 null
  return null;
}

String formatDateTime(DateTime dt, {bool withSeconds = true}) {
  final y = dt.year.toString().padLeft(4, '0');
  final m = dt.month.toString().padLeft(2, '0');
  final d = dt.day.toString().padLeft(2, '0');
  final hh = dt.hour.toString().padLeft(2, '0');
  final mm = dt.minute.toString().padLeft(2, '0');
  final ss = dt.second.toString().padLeft(2, '0');
  return withSeconds ? '$y-$m-$d $hh:$mm:$ss' : '$y-$m-$d $hh:$mm';
}

/// 入口方法：传入多种格式原始值，输出人类可读格式；失败返回 '-'
String formatDateTimeFlexible(dynamic raw, {bool withSeconds = true}) {
  final dt = parseFlexibleDateTime(raw);
  if (dt == null) return '-';
  return formatDateTime(dt, withSeconds: withSeconds);
}
