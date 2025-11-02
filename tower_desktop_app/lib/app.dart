import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import 'features/auth/login_screen.dart';
import 'features/menu/menu_api.dart';
import 'features/menu/menu_provider.dart';

class TowerApp extends StatelessWidget {
  const TowerApp({super.key});

  @override
  Widget build(BuildContext context) {
    return MultiProvider(
      providers: [
        ChangeNotifierProvider(create: (_) => MenuProvider(MenuApi())),
      ],
      child: MaterialApp(
        title: 'Tower Desktop',
        theme: ThemeData(useMaterial3: true, colorSchemeSeed: Colors.blue),
        home: const LoginScreen(),
      ),
    );
  }
}
