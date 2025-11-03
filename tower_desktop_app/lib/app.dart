import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import 'features/auth/login_screen.dart';
import 'features/auth/permission_provider.dart';
import 'features/auth/session_manager.dart';
import 'features/home/home_screen.dart' as home;
import 'features/menu/menu_api.dart';
import 'features/menu/menu_provider.dart';
import 'features/user/user_api.dart';
import 'features/user/user_provider.dart';
import 'features/store/store_api.dart';
import 'features/store/store_provider.dart';
import 'core/network/api_client.dart';

class TowerApp extends StatelessWidget {
  const TowerApp({super.key});

  @override
  Widget build(BuildContext context) {
    final apiClient = ApiClient();
    final storeApi = StoreApi(apiClient);
    final userApi = UserApi(apiClient);

    return MultiProvider(
      providers: [
        ChangeNotifierProvider(create: (_) => PermissionProvider()),
        ChangeNotifierProvider(create: (_) => MenuProvider(MenuApi())),
        ChangeNotifierProvider(create: (_) => UserProvider(userApi, storeApi)),
        ChangeNotifierProvider(create: (_) => StoreProvider(storeApi)),
      ],
      child: FutureBuilder(
        future: _bootstrap(context),
        builder: (context, snap) {
          if (snap.connectionState != ConnectionState.done) {
            return const MaterialApp(
              home: Scaffold(
                body: Center(child: CircularProgressIndicator()),
              ),
            );
          }
          return MaterialApp(
            title: 'Tower Desktop',
            theme: ThemeData(useMaterial3: true, colorSchemeSeed: Colors.blue),
            home:
                snap.data == true ? const _DeferredHome() : const LoginScreen(),
          );
        },
      ),
    );
  }
}

/// 延迟加载 Home：需要 MenuProvider 先加载菜单
class _DeferredHome extends StatefulWidget {
  const _DeferredHome();
  @override
  State<_DeferredHome> createState() => _DeferredHomeState();
}

class _DeferredHomeState extends State<_DeferredHome> {
  bool _ready = false;
  @override
  void initState() {
    super.initState();
    WidgetsBinding.instance.addPostFrameCallback((_) {
      _prepare();
    });
  }

  Future<void> _prepare() async {
    final sm = SessionManager();
    if (!sm.isLoggedIn || sm.isExpired) {
      if (mounted)
        Navigator.of(context).pushReplacement(
            MaterialPageRoute(builder: (_) => const LoginScreen()));
      return;
    }
    // 恢复权限到 provider
    final permProvider = context.read<PermissionProvider>();
    permProvider.setPermissions(sm.permissions);
    // 加载菜单
    final menuProvider = context.read<MenuProvider>();
    await menuProvider.load(permissionProvider: permProvider);
    if (mounted) {
      setState(() => _ready = true);
    }
  }

  @override
  Widget build(BuildContext context) {
    if (!_ready) {
      return const Scaffold(body: Center(child: CircularProgressIndicator()));
    }
    // 直接复用 HomeScreen，不在此文件重复导入依赖以避免循环；延迟 import
    return const home.HomeScreen();
  }
}

Future<bool> _bootstrap(BuildContext context) async {
  final sm = SessionManager();
  if (sm.isLoggedIn && !sm.isExpired) {
    return true; // 走自动登录流程
  }
  return false;
}
