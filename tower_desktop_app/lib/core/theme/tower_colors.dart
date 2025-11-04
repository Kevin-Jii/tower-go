import 'package:flutter/material.dart';

class TowerColors extends ThemeExtension<TowerColors> {
  final Color navBarBackground;
  final Color navBarForeground;
  final Color navBarBorder;
  final Color sideBarBackground;
  final Color sideBarBorder;
  final Color contentBackground;
  final Color subtleFill;
  final Color accentGradientStart;
  final Color accentGradientEnd;

  const TowerColors({
    required this.navBarBackground,
    required this.navBarForeground,
    required this.navBarBorder,
    required this.sideBarBackground,
    required this.sideBarBorder,
    required this.contentBackground,
    required this.subtleFill,
    required this.accentGradientStart,
    required this.accentGradientEnd,
  });

  @override
  TowerColors copyWith({
    Color? navBarBackground,
    Color? navBarForeground,
    Color? navBarBorder,
    Color? sideBarBackground,
    Color? sideBarBorder,
    Color? contentBackground,
    Color? subtleFill,
    Color? accentGradientStart,
    Color? accentGradientEnd,
  }) =>
      TowerColors(
        navBarBackground: navBarBackground ?? this.navBarBackground,
        navBarForeground: navBarForeground ?? this.navBarForeground,
        navBarBorder: navBarBorder ?? this.navBarBorder,
        sideBarBackground: sideBarBackground ?? this.sideBarBackground,
        sideBarBorder: sideBarBorder ?? this.sideBarBorder,
        contentBackground: contentBackground ?? this.contentBackground,
        subtleFill: subtleFill ?? this.subtleFill,
        accentGradientStart: accentGradientStart ?? this.accentGradientStart,
        accentGradientEnd: accentGradientEnd ?? this.accentGradientEnd,
      );

  @override
  TowerColors lerp(ThemeExtension<TowerColors>? other, double t) {
    if (other is! TowerColors) return this;
    return TowerColors(
      navBarBackground:
          Color.lerp(navBarBackground, other.navBarBackground, t)!,
      navBarForeground:
          Color.lerp(navBarForeground, other.navBarForeground, t)!,
      navBarBorder: Color.lerp(navBarBorder, other.navBarBorder, t)!,
      sideBarBackground:
          Color.lerp(sideBarBackground, other.sideBarBackground, t)!,
      sideBarBorder: Color.lerp(sideBarBorder, other.sideBarBorder, t)!,
      contentBackground:
          Color.lerp(contentBackground, other.contentBackground, t)!,
      subtleFill: Color.lerp(subtleFill, other.subtleFill, t)!,
      accentGradientStart:
          Color.lerp(accentGradientStart, other.accentGradientStart, t)!,
      accentGradientEnd:
          Color.lerp(accentGradientEnd, other.accentGradientEnd, t)!,
    );
  }
}
