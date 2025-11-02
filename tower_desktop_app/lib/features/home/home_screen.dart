import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import '../auth/login_screen.dart';
import '../auth/session_manager.dart';
import '../menu/menu_provider.dart';
import '../menu/menu_tree.dart';
import 'widgets/menu_content.dart';

class HomeScreen extends StatefulWidget {
  const HomeScreen({super.key});

  @override
  State<HomeScreen> createState() => _HomeScreenState();
}

class _HomeScreenState extends State<HomeScreen> {
  void _handleLogout() {
    SessionManager().clear();
    Navigator.of(context).pushAndRemoveUntil(
        MaterialPageRoute(builder: (_) => const LoginScreen()), (_) => false);
  }

  @override
  Widget build(BuildContext context) {
    final mp = context.watch<MenuProvider>();

    return Scaffold(
      appBar: AppBar(
        title: const Text('Tower 管理桌面端'),
        actions: [
          IconButton(
            onPressed: _handleLogout,
            icon: const Icon(Icons.logout),
            tooltip: '退出登录',
          )
        ],
      ),
      body: Row(
        children: [
          // 左侧菜单栏
          SizedBox(
            width: 260,
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.stretch,
              children: [
                Container(
                  padding: const EdgeInsets.all(12),
                  decoration: BoxDecoration(
                    color:
                        Theme.of(context).colorScheme.primary.withOpacity(.08),
                  ),
                  child: Text(
                    '功能菜单',
                    style: Theme.of(context).textTheme.titleMedium,
                  ),
                ),
                const Expanded(child: MenuTree()),
              ],
            ),
          ),
          const VerticalDivider(width: 1),
          // 右侧内容区
          Expanded(
            child: mp.selected == null
                ? const Center(child: Text('请选择左侧菜单'))
                : MenuContent(menuItem: mp.selected!),
          )
        ],
      ),
    );
  }
}
