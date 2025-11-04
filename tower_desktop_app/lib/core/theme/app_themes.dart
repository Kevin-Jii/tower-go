import 'package:flutter/material.dart';
import 'tower_colors.dart';

/// 定义多套可选主题。可继续扩展。
class AppThemes {
  static final ThemeData defaultLight = ThemeData(
    useMaterial3: true,
    colorScheme: ColorScheme.fromSeed(
        seedColor: Colors.blue, brightness: Brightness.light),
    scaffoldBackgroundColor: const Color(0xFFF5F7FA),
    cardTheme: CardThemeData(
      elevation: 3,
      shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(14)),
    ),
    inputDecorationTheme: InputDecorationTheme(
      border: OutlineInputBorder(borderRadius: BorderRadius.circular(10)),
    ),
    extensions: const [
      TowerColors(
        navBarBackground: Colors.white,
        navBarForeground: Color(0xFF1A1F29),
        navBarBorder: Color(0xFFE3E8EF),
        sideBarBackground: Color(0xFFF8FAFC),
        sideBarBorder: Color(0xFFE2E8F0),
        contentBackground: Color(0xFFF5F7FA),
        subtleFill: Color(0xFFF1F5F9),
        accentGradientStart: Color(0xFF3B82F6),
        accentGradientEnd: Color(0xFF6366F1),
      ),
    ],
  );

  static final ThemeData dark = ThemeData(
    useMaterial3: true,
    colorScheme: ColorScheme.fromSeed(
        seedColor: Colors.indigo, brightness: Brightness.dark),
    scaffoldBackgroundColor: const Color(0xFF12161C),
    cardColor: const Color(0xFF1E252E),
    inputDecorationTheme: InputDecorationTheme(
      border: OutlineInputBorder(borderRadius: BorderRadius.circular(10)),
    ),
    extensions: const [
      TowerColors(
        navBarBackground: Color(0xFF1E252E),
        navBarForeground: Color(0xFFE6EDF5),
        navBarBorder: Color(0xFF2C3742),
        sideBarBackground: Color(0xFF192027),
        sideBarBorder: Color(0xFF27313C),
        contentBackground: Color(0xFF12161C),
        subtleFill: Color(0xFF222B33),
        accentGradientStart: Color(0xFF6366F1),
        accentGradientEnd: Color(0xFF8B5CF6),
      ),
    ],
  );

  static final ThemeData ocean = ThemeData(
    useMaterial3: true,
    colorScheme: ColorScheme.fromSeed(
        seedColor: Colors.teal, brightness: Brightness.light),
    scaffoldBackgroundColor: const Color(0xFFE6F4F3),
    appBarTheme: const AppBarTheme(
        backgroundColor: Color(0xFF0B7A75), foregroundColor: Colors.white),
    chipTheme: ChipThemeData(
      backgroundColor: Colors.teal.shade50,
      selectedColor: Colors.teal.shade200,
      labelStyle: const TextStyle(color: Colors.teal),
      shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(8)),
    ),
    extensions: const [
      TowerColors(
        navBarBackground: Color(0xFF0B7A75),
        navBarForeground: Colors.white,
        navBarBorder: Color(0xFF0B6E69),
        sideBarBackground: Color(0xFFF0FBFA),
        sideBarBorder: Color(0xFFCCE7E5),
        contentBackground: Color(0xFFE6F4F3),
        subtleFill: Color(0xFFD9EFEE),
        accentGradientStart: Color(0xFF14B8A6),
        accentGradientEnd: Color(0xFF0D9488),
      ),
    ],
  );

  static final ThemeData forest = ThemeData(
    useMaterial3: true,
    colorScheme: ColorScheme.fromSeed(
        seedColor: Colors.green, brightness: Brightness.light),
    scaffoldBackgroundColor: const Color(0xFFF4F8F4),
    appBarTheme: const AppBarTheme(
        backgroundColor: Color(0xFF2E7D32), foregroundColor: Colors.white),
    floatingActionButtonTheme:
        const FloatingActionButtonThemeData(backgroundColor: Color(0xFF43A047)),
    extensions: const [
      TowerColors(
        navBarBackground: Color(0xFF2E7D32),
        navBarForeground: Colors.white,
        navBarBorder: Color(0xFF256428),
        sideBarBackground: Color(0xFFF2F9F2),
        sideBarBorder: Color(0xFFD7E7D6),
        contentBackground: Color(0xFFF4F8F4),
        subtleFill: Color(0xFFE3EEE2),
        accentGradientStart: Color(0xFF43A047),
        accentGradientEnd: Color(0xFF2E7D32),
      ),
    ],
  );

  static List<ThemeData> all() => [defaultLight, dark, ocean, forest];
  static const List<String> names = [
    '默认',
    '暗黑',
  ];
}
