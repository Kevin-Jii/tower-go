import 'package:flutter/material.dart';
import 'app.dart';
import 'features/auth/session_manager.dart';

Future<void> main() async {
  WidgetsFlutterBinding.ensureInitialized();
  await SessionManager().init();
  runApp(const TowerApp());
}
