import 'package:flutter/widgets.dart';
import 'package:tdesign_flutter/tdesign_flutter.dart';

/// 后端传来的 icon 字符串 与 TDesign Icon 的映射
/// 参考: https://tdesign.tencent.com/vue/components/icon
class TdIconMapper {
  TdIconMapper._();

  static Widget build(String? name, {double size = 18, Color? color}) {
    final icon = _iconData(name);
    if (icon != null) {
      return Icon(icon, size: size, color: color);
    }
    // fallback 使用一个默认的方块或常规图标
    return Icon(TDIcons.app, size: size, color: color);
  }

  static IconData? _iconData(String? name) {
    if (name == null || name.isEmpty) return null;
    switch (name) {
      // 系统/设置类
      case 'setting':
      case 'settings':
        return TDIcons.setting;
      case 'user':
        return TDIcons.user;
      case 'team':
      case 'users':
        return TDIcons.usergroup;
      case 'menu':
        return TDIcons.menu;
      case 'shop':
      case 'store':
        return TDIcons.shop;
      case 'unordered-list':
      case 'list':
        return TDIcons.list;
      case 'coffee':
      case 'dish':
        return TDIcons.tea;
      case 'file-text':
      case 'document':
        return TDIcons.file;
      case 'bar-chart':
      case 'chart':
        return TDIcons.chart_bar;
      case 'role':
        return TDIcons.secured;
      default:
        return null;
    }
  }
}
