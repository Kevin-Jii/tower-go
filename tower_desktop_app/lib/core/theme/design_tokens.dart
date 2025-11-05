import 'package:flutter/material.dart';

/// 设计令牌：集中管理颜色、间距、圆角、阴影、文本层级等。
/// 后续可与品牌色/暗色模式切换结合。
class DesignTokens {
  // 基础色板（可根据品牌再调）
  static const Color brandPrimary = Color(0xFF2563EB); // 主色 (蓝)
  static const Color brandPrimaryHover = Color(0xFF1D4ED8);
  static const Color brandPrimarySoft = Color(0xFFE0F2FE);
  static const Color danger = Color(0xFFDC2626);
  static const Color dangerSoft = Color(0xFFFEE2E2);
  static const Color neutral0 = Colors.white;
  static const Color neutral50 = Color(0xFFF8FAFC);
  static const Color neutral100 = Color(0xFFF1F5F9);
  static const Color neutral200 = Color(0xFFE2E8F0);
  static const Color neutral300 = Color(0xFFCBD5E1);
  static const Color neutral400 = Color(0xFF94A3B8);
  static const Color neutral500 = Color(0xFF64748B);
  static const Color neutral600 = Color(0xFF475569);
  static const Color neutral700 = Color(0xFF334155);
  static const Color neutral800 = Color(0xFF1E293B);

  // 间距 (统一使用 4 的倍数节奏)
  static const double spaceXs = 4;
  static const double spaceSm = 8;
  static const double spaceMd = 16;
  static const double spaceLg = 24;
  static const double spaceXl = 32;

  // 圆角
  static const BorderRadius radiusSm = BorderRadius.all(Radius.circular(4));
  static const BorderRadius radiusMd = BorderRadius.all(Radius.circular(8));
  static const BorderRadius radiusLg = BorderRadius.all(Radius.circular(12));

  // 阴影（可后续抽离暗色模式）
  static List<BoxShadow> shadowSm = [
    BoxShadow(
      color: Colors.black.withOpacity(.05),
      blurRadius: 6,
      offset: const Offset(0, 2),
    )
  ];
  static List<BoxShadow> shadowMd = [
    BoxShadow(
      color: Colors.black.withOpacity(.06),
      blurRadius: 12,
      spreadRadius: 1,
      offset: const Offset(0, 4),
    )
  ];

  // 文本层级
  static const TextStyle heading1 =
      TextStyle(fontSize: 22, fontWeight: FontWeight.w600, letterSpacing: .2);
  static const TextStyle heading2 =
      TextStyle(fontSize: 18, fontWeight: FontWeight.w600, letterSpacing: .15);
  static const TextStyle heading3 =
      TextStyle(fontSize: 16, fontWeight: FontWeight.w600, letterSpacing: .15);
  static const TextStyle subtitle =
      TextStyle(fontSize: 14, fontWeight: FontWeight.w500, letterSpacing: .1);
  static const TextStyle body =
      TextStyle(fontSize: 14, fontWeight: FontWeight.w400, letterSpacing: .1);
  static const TextStyle bodySm =
      TextStyle(fontSize: 12, fontWeight: FontWeight.w400, letterSpacing: .1);
  static const TextStyle caption =
      TextStyle(fontSize: 11, fontWeight: FontWeight.w400, letterSpacing: .2);

  // 侧边栏宽度
  static const double sidebarWidth = 264;

  // 分类项高度
  static const double categoryItemHeight = 44;
}
