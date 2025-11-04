import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import '../auth/login_screen.dart';
import '../auth/session_manager.dart';
import '../auth/models.dart';
import '../menu/menu_provider.dart';
import '../menu/menu_tree.dart';
import 'widgets/menu_content.dart';
import '../../core/theme/theme_provider.dart';
import '../../core/theme/tower_colors.dart';
import '../../core/constants/app_constants.dart';

class HomeScreen extends StatefulWidget {
  const HomeScreen({super.key});

  @override
  State<HomeScreen> createState() => _HomeScreenState();
}

class _HomeScreenState extends State<HomeScreen> {
  void _handleLogout() {
    SessionManager().clear(persist: true);
    Navigator.of(context).pushAndRemoveUntil(
        MaterialPageRoute(builder: (_) => const LoginScreen()), (_) => false);
  }

  @override
  Widget build(BuildContext context) {
    final mp = context.watch<MenuProvider>();

    final user = SessionManager().userInfo;

    return Scaffold(
      // 自定义顶部导航条替换默认 AppBar
      body: Column(
        children: [
          _buildTopBar(context, user),
          Expanded(
            child: Row(
              children: [
                // 左侧菜单栏
                SizedBox(
                  width: 260,
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.stretch,
                    children: [
                      // 移除“功能菜单”标题，直接展示菜单树
                      const Expanded(child: MenuTree()),
                    ],
                  ),
                ),
                const VerticalDivider(width: 1),
                // 右侧内容区
                Expanded(
                  child: mp.selected == null
                      ? const Center(
                          child: Text(AppTexts.selectMenuPlaceholder))
                      : MenuContent(menuItem: mp.selected!),
                )
              ],
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildTopBar(BuildContext context, UserInfo? user) {
    final tc = Theme.of(context).extension<TowerColors>()!;
    return Container(
      height: 60,
      padding: const EdgeInsets.symmetric(horizontal: 20),
      decoration: BoxDecoration(
        color: tc.navBarBackground,
        boxShadow: [
          BoxShadow(
            color: tc.navBarBorder.withOpacity(.35),
            blurRadius: 8,
            offset: const Offset(0, 2),
          ),
        ],
      ),
      child: Row(
        children: [
          // 左侧项目名称与可选环境标签
          Row(
            children: [
              Icon(Icons.dashboard_customize,
                  color: Theme.of(context).colorScheme.primary, size: 26),
              const SizedBox(width: 10),
              ShaderMask(
                shaderCallback: (r) => LinearGradient(
                  colors: [tc.accentGradientStart, tc.accentGradientEnd],
                ).createShader(r),
                child: const Text(
                  AppTexts.appName,
                  style: TextStyle(
                    fontSize: 18,
                    fontWeight: FontWeight.w700,
                    color: Colors.white,
                    letterSpacing: .6,
                  ),
                ),
              ),
              const SizedBox(width: 12),
              _envTag(AppTexts.envProd),
            ],
          ),
          const Spacer(),
          // 搜索框占位
          // 主题切换按钮组（替换原搜索框）
          _themeSwitcher(context),
          const SizedBox(width: 12),
          // 通知
          _iconButton(context, Icons.notifications_none, '通知', () {}),
          _iconButton(context, Icons.settings_outlined, '设置', () {}),
          const SizedBox(width: 8),
          // 用户头像 + 下拉菜单
          _userMenu(context, user),
        ],
      ),
    );
  }

  Widget _themeSwitcher(BuildContext context) {
    final tp = Provider.of<ThemeProvider>(context);
    final tc = Theme.of(context).extension<TowerColors>()!;
    return PopupMenuButton<int>(
      tooltip: '切换主题',
      onSelected: (i) => tp.setTheme(i),
      itemBuilder: (ctx) => List.generate(tp.names.length, (i) {
        return PopupMenuItem<int>(
          value: i,
          child: Row(
            children: [
              Icon(
                i == tp.index ? Icons.check_circle : Icons.circle_outlined,
                size: 18,
                color: i == tp.index
                    ? tc.accentGradientStart
                    : Colors.grey.shade500,
              ),
              const SizedBox(width: 8),
              Text(tp.names[i]),
            ],
          ),
        );
      }),
      child: Container(
        height: 36,
        padding: const EdgeInsets.symmetric(horizontal: 12),
        decoration: BoxDecoration(
          color: tc.subtleFill,
          borderRadius: BorderRadius.circular(8),
        ),
        child: Row(
          children: [
            Icon(Icons.palette_outlined,
                size: 18, color: tc.navBarForeground.withOpacity(.7)),
            const SizedBox(width: 6),
            Text(tp.names[tp.index],
                style: TextStyle(fontSize: 13, color: tc.navBarForeground)),
            const Icon(Icons.keyboard_arrow_down, size: 18),
          ],
        ),
      ),
    );
  }

  Widget _envTag(String text) {
    final tc = Theme.of(context).extension<TowerColors>()!;
    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 2),
      decoration: BoxDecoration(
        color: tc.subtleFill,
        borderRadius: BorderRadius.circular(6),
        border: Border.all(color: tc.navBarBorder),
      ),
      child: Text(
        text,
        style: TextStyle(
          fontSize: 11,
          fontWeight: FontWeight.w600,
          color: tc.accentGradientStart,
          letterSpacing: .5,
        ),
      ),
    );
  }

  Widget _iconButton(
      BuildContext context, IconData icon, String tooltip, VoidCallback onTap) {
    final tc = Theme.of(context).extension<TowerColors>()!;
    return Padding(
      padding: const EdgeInsets.symmetric(horizontal: 4),
      child: InkWell(
        onTap: onTap,
        borderRadius: BorderRadius.circular(30),
        child: Tooltip(
          message: tooltip,
          child: Container(
            width: 36,
            height: 36,
            decoration: BoxDecoration(
              color: tc.subtleFill,
              borderRadius: BorderRadius.circular(18),
            ),
            child: Icon(icon,
                size: 20, color: tc.navBarForeground.withOpacity(.75)),
          ),
        ),
      ),
    );
  }

  Widget _userMenu(BuildContext context, UserInfo? user) {
    final displayName = (user != null && user.nickname.isNotEmpty)
        ? user.nickname
        : (user?.username ?? '未登录');
    final tc = Theme.of(context).extension<TowerColors>()!;
    return PopupMenuButton<String>(
      offset: const Offset(0, 50),
      onSelected: (value) {
        if (value == 'logout') {
          _handleLogout();
        } else if (value == 'profile') {
          // 预留：跳转个人设置页面
          ScaffoldMessenger.of(context).showSnackBar(
            const SnackBar(content: Text('个人设置功能待接入')),
          );
        }
      },
      itemBuilder: (ctx) => [
        PopupMenuItem(
          value: 'profile',
          child: Row(
            children: const [
              Icon(Icons.person_outline, size: 18),
              SizedBox(width: 8),
              Text('个人设置'),
            ],
          ),
        ),
        const PopupMenuDivider(),
        PopupMenuItem(
          value: 'logout',
          child: Row(
            children: const [
              Icon(Icons.logout, size: 18),
              SizedBox(width: 8),
              Text('退出登录'),
            ],
          ),
        ),
      ],
      child: Row(
        children: [
          CircleAvatar(
            radius: 16,
            backgroundColor: tc.accentGradientStart,
            child: Text(
              displayName.isNotEmpty ? displayName.characters.first : '?',
              style: const TextStyle(color: Colors.white),
            ),
          ),
          const SizedBox(width: 8),
          Text(
            displayName,
            style: TextStyle(
              fontSize: 14,
              fontWeight: FontWeight.w500,
              color: tc.navBarForeground,
            ),
          ),
          const Icon(Icons.keyboard_arrow_down, size: 18),
        ],
      ),
    );
  }
}
