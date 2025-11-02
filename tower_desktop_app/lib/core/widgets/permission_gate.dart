import 'package:flutter/widgets.dart';
import '../../features/auth/session_manager.dart';

class PermissionGate extends StatelessWidget {
  final String permission;
  final Widget child;
  final Widget? fallback;
  const PermissionGate({super.key, required this.permission, required this.child, this.fallback});

  @override
  Widget build(BuildContext context) {
    final has = SessionManager().hasPermission(permission);
    if (has) return child;
    return fallback ?? const SizedBox.shrink();
  }
}
