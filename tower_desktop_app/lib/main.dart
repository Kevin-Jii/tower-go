import 'package:flutter/material.dart';
import 'features/auth/auth_api.dart';
import 'features/auth/models.dart';
import 'features/auth/session_manager.dart';
import 'features/menu/menu_api.dart';

void main() {
  runApp(const TowerApp());
}

class TowerApp extends StatelessWidget {
  const TowerApp({super.key});

  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      title: 'Tower Desktop',
      theme: ThemeData(useMaterial3: true, colorSchemeSeed: Colors.blue),
      home: const LoginScreen(),
    );
  }
}

class LoginScreen extends StatefulWidget {
  const LoginScreen({super.key});
  @override
  State<LoginScreen> createState() => _LoginScreenState();
}

class _LoginScreenState extends State<LoginScreen> {
  final _userCtrl = TextEditingController();
  final _pwdCtrl = TextEditingController();
  bool _loading = false;
  String? _error;

  Future<void> _doLogin() async {
    setState(() { _loading = true; _error = null; });
    try {
      final authApi = AuthApi();
      final resp = await authApi.login(LoginRequest(username: _userCtrl.text.trim(), password: _pwdCtrl.text));
      final menuApi = MenuApi();
      final perms = await menuApi.getUserPermissions();
      SessionManager().updateSession(token: resp.token, user: resp.user, permissions: perms);
      if (!mounted) return;
      Navigator.of(context).pushReplacement(MaterialPageRoute(builder: (_) => const HomeScreen()));
    } catch (e) {
      setState(() { _error = e.toString(); });
    } finally {
      if (mounted) setState(() { _loading = false; });
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: Center(
        child: ConstrainedBox(
          constraints: const BoxConstraints(maxWidth: 360),
          child: Column(
            mainAxisSize: MainAxisSize.min,
            children: [
              TextField(controller: _userCtrl, decoration: const InputDecoration(labelText: '用户名')),          
              TextField(controller: _pwdCtrl, decoration: const InputDecoration(labelText: '密码'), obscureText: true),
              const SizedBox(height: 16),
              if (_error != null) Text(_error!, style: const TextStyle(color: Colors.red)),
              FilledButton(
                onPressed: _loading ? null : _doLogin,
                child: _loading ? const SizedBox(width:16,height:16,child:CircularProgressIndicator(strokeWidth:2)) : const Text('登录'),
              ),
            ],
          ),
        ),
      ),
    );
  }
}

class HomeScreen extends StatefulWidget {
  const HomeScreen({super.key});
  @override
  State<HomeScreen> createState() => _HomeScreenState();
}

class _HomeScreenState extends State<HomeScreen> {
  int _index = 0;
  @override
  Widget build(BuildContext context) {
    final tabs = <Widget>[
      const Center(child: Text('菜单 / 权限 待实现')), // 可替换为动态树
      const Center(child: Text('菜品列表占位')),
      const Center(child: Text('报菜记录占位')),
      const Center(child: Text('门店管理占位')),
    ];
    return Scaffold(
      appBar: AppBar(title: const Text('Tower 管理桌面端'), actions: [
        IconButton(onPressed: () { SessionManager().clear(); Navigator.of(context).pushAndRemoveUntil(MaterialPageRoute(builder: (_)=>const LoginScreen()), (_) => false); }, icon: const Icon(Icons.logout))
      ]),
      body: Row(
        children: [
          NavigationRail(
            selectedIndex: _index,
            onDestinationSelected: (i)=> setState(()=> _index = i),
            labelType: NavigationRailLabelType.all,
            destinations: const [
              NavigationRailDestination(icon: Icon(Icons.menu), label: Text('权限')),
              NavigationRailDestination(icon: Icon(Icons.restaurant), label: Text('菜品')),
              NavigationRailDestination(icon: Icon(Icons.list_alt), label: Text('报菜')),
              NavigationRailDestination(icon: Icon(Icons.store), label: Text('门店')),
            ],
          ),
          const VerticalDivider(width: 1),
          Expanded(child: tabs[_index]),
        ],
      ),
    );
  }
}
