import 'package:flutter/material.dart';
import '../utils/time_utils.dart';

/// 全局统一的日期展示组件。
/// - raw: 可为字符串 / 数字时间戳(秒或毫秒) / DateTime / null
/// - placeholder: 无法解析时显示占位符（默认 '-'）
/// - withSeconds: 是否包含秒
/// - style: 文本样式（未提供则使用主题 bodySmall）
/// - tooltip: 是否显示 tooltip 原始值
class DateText extends StatelessWidget {
  final dynamic raw;
  final String placeholder;
  final bool withSeconds;
  final TextStyle? style;
  final bool tooltip;
  const DateText({
    super.key,
    required this.raw,
    this.placeholder = '-',
    this.withSeconds = true,
    this.style,
    this.tooltip = true,
  });

  @override
  Widget build(BuildContext context) {
    final dt = parseFlexibleDateTime(raw);
    final display =
        dt == null ? placeholder : formatDateTime(dt, withSeconds: withSeconds);
    final textWidget =
        Text(display, style: style ?? Theme.of(context).textTheme.bodySmall);
    if (!tooltip || dt == null) return textWidget;
    // Tooltip 展示原始 & ISO8601
    final iso = dt.toIso8601String();
    return Tooltip(
      message: '原始: ${raw ?? 'null'}\nISO: $iso',
      child: textWidget,
    );
  }
}
