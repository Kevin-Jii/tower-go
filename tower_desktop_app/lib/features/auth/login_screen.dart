import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import 'package:tdesign_flutter/tdesign_flutter.dart';
import 'auth_api.dart';
import 'models.dart';
import 'session_manager.dart';
import '../menu/menu_api.dart';
import '../menu/menu_provider.dart';
import 'permission_provider.dart';
import '../home/home_screen.dart';

class LoginScreen extends StatefulWidget {
  const LoginScreen({super.key});

  @override
  State<LoginScreen> createState() => _LoginScreenState();
}

class _LoginScreenState extends State<LoginScreen> {
  final _phoneCtrl = TextEditingController();
  final _pwdCtrl = TextEditingController();
  bool _loading = false;

  @override
  void dispose() {
    _phoneCtrl.dispose();
    _pwdCtrl.dispose();
    super.dispose();
  }

  Future<void> _doLogin() async {
    setState(() {
      _loading = true;
    });
    try {
      final authApi = AuthApi();
      final resp = await authApi.login(
          LoginRequest(phone: _phoneCtrl.text.trim(), password: _pwdCtrl.text));
      final menuApi = MenuApi();
      final perms = await menuApi.getUserPermissions();
      // 保存权限到 PermissionProvider 以便 UI 立即生效
      final permProvider = context.read<PermissionProvider>();
      permProvider.setPermissions(perms);
      SessionManager().updateSession(
          token: resp.token,
          userInfo: resp.userInfo,
          permissions: perms,
          expiresIn: resp.expiresIn == 0 ? null : resp.expiresIn);
      if (!mounted) return;

      // 登录成功后先加载菜单
      final menuProvider = context.read<MenuProvider>();
      await menuProvider.load(permissionProvider: permProvider);
      if (!mounted) return;

      Navigator.of(context).pushReplacement(
          MaterialPageRoute(builder: (_) => const HomeScreen()));
    } catch (e) {
      if (!mounted) return;
      TDToast.showText(e.toString(), context: context);
    } finally {
      if (mounted) {
        setState(() {
          _loading = false;
        });
      }
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
              TextField(
                controller: _phoneCtrl,
                decoration: const InputDecoration(labelText: '手机号'),
                keyboardType: TextInputType.phone,
              ),
              const SizedBox(height: 12),
              TextField(
                controller: _pwdCtrl,
                decoration: const InputDecoration(labelText: '密码'),
                obscureText: true,
                onSubmitted: (_) => _doLogin(),
              ),
              const SizedBox(height: 24),
              TDButton(
                text: _loading ? '登录中...' : '登录',
                onTap: _loading ? null : _doLogin,
                size: TDButtonSize.large,
                theme: TDButtonTheme.primary,
              )
            ],
          ),
        ),
      ),
    );
  }
}
