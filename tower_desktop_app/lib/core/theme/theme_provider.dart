import 'package:flutter/material.dart';
import 'app_themes.dart';

class ThemeProvider extends ChangeNotifier {
  final List<ThemeData> _themes = AppThemes.all();
  int _index = 0;

  ThemeData get theme => _themes[_index];
  int get index => _index;
  List<String> get names => AppThemes.names;

  void setTheme(int i) {
    if (i < 0 || i >= _themes.length) return;
    _index = i;
    notifyListeners();
  }

  void next() {
    _index = (_index + 1) % _themes.length;
    notifyListeners();
  }
}

extension ThemeProviderExt on BuildContext {
  ThemeProvider themeProvider() =>
      dependOnInheritedWidgetOfExactType<_ThemeInherited>()!.data;
}

/// 简单 Inherited 用于无需层层 Provider.of 访问，可选
class ThemeProviderHost extends StatefulWidget {
  final ThemeProvider provider;
  final Widget child;
  const ThemeProviderHost(
      {super.key, required this.provider, required this.child});

  @override
  State<ThemeProviderHost> createState() => _ThemeProviderHostState();
}

class _ThemeProviderHostState extends State<ThemeProviderHost> {
  @override
  void initState() {
    super.initState();
    widget.provider.addListener(_onChange);
  }

  void _onChange() => setState(() {});

  @override
  void dispose() {
    widget.provider.removeListener(_onChange);
    super.dispose();
  }

  @override
  Widget build(BuildContext context) =>
      _ThemeInherited(data: widget.provider, child: widget.child);
}

class _ThemeInherited extends InheritedWidget {
  final ThemeProvider data;
  const _ThemeInherited({required this.data, required super.child});
  @override
  bool updateShouldNotify(covariant _ThemeInherited oldWidget) =>
      oldWidget.data.index != data.index;
}
