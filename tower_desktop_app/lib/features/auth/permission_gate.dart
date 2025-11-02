import 'package:flutter/widgets.dart';
import 'package:provider/provider.dart';
import 'permission_provider.dart';

/// 用于包裹需要权限的组件
class PermissionGate extends StatelessWidget {
  final String? required;
  final List<String>? any;
  final List<String>? all;
  final Widget child;
  final Widget? no;

  const PermissionGate({
    super.key,
    this.required,
    this.any,
    this.all,
    required this.child,
    this.no,
  });

  @override
  Widget build(BuildContext context) {
    final p = context.watch<PermissionProvider>();

    bool allow = true;
    if (required != null) {
      allow = p.has(required!);
    }
    if (allow && any != null && any!.isNotEmpty) {
      allow = p.hasAny(any!);
    }
    if (allow && all != null && all!.isNotEmpty) {
      allow = p.hasAll(all!);
    }

    if (allow) return child;
    return no ?? const SizedBox.shrink();
  }
}
