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
import '../../core/constants/app_constants.dart';

class LoginScreen extends StatefulWidget {
  const LoginScreen({super.key});

  @override
  State<LoginScreen> createState() => _LoginScreenState();
}

class _LoginScreenState extends State<LoginScreen>
    with SingleTickerProviderStateMixin {
  final _formKey = GlobalKey<FormState>();
  final _phoneCtrl = TextEditingController();
  final _pwdCtrl = TextEditingController();
  bool _loading = false;
  bool _obscurePassword = true;
  late AnimationController _animationController;
  late Animation<double> _fadeAnimation;

  @override
  void initState() {
    super.initState();
    _animationController = AnimationController(
      vsync: this,
      duration: const Duration(milliseconds: 800),
    );
    _fadeAnimation = Tween<double>(begin: 0.0, end: 1.0).animate(
      CurvedAnimation(parent: _animationController, curve: Curves.easeIn),
    );
    _animationController.forward();
  }

  @override
  void dispose() {
    _phoneCtrl.dispose();
    _pwdCtrl.dispose();
    _animationController.dispose();
    super.dispose();
  }

  Future<void> _doLogin() async {
    if (!_formKey.currentState!.validate()) {
      return;
    }

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
    final screenWidth = MediaQuery.of(context).size.width;
    final isMobile = screenWidth < 800;

    return Scaffold(
      body: Container(
        decoration: BoxDecoration(
          gradient: LinearGradient(
            begin: Alignment.topLeft,
            end: Alignment.bottomRight,
            colors: [
              Colors.blue.shade50,
              Colors.purple.shade50,
              Colors.pink.shade50,
            ],
          ),
        ),
        child: Center(
          child: FadeTransition(
            opacity: _fadeAnimation,
            child: ConstrainedBox(
              constraints: BoxConstraints(
                maxWidth: isMobile ? screenWidth * 0.9 : 1000,
                maxHeight: isMobile ? double.infinity : 700,
              ),
              child: Card(
                elevation: 12,
                shadowColor: Colors.black.withOpacity(0.2),
                shape: RoundedRectangleBorder(
                  borderRadius: BorderRadius.circular(24),
                ),
                child: isMobile ? _buildMobileLayout() : _buildDesktopLayout(),
              ),
            ),
          ),
        ),
      ),
    );
  }

  Widget _buildMobileLayout() {
    return SingleChildScrollView(
      padding: const EdgeInsets.all(32),
      child: Column(
        mainAxisSize: MainAxisSize.min,
        children: [
          _buildBrandingSection(compact: true),
          const SizedBox(height: 32),
          _buildLoginForm(),
        ],
      ),
    );
  }

  Widget _buildDesktopLayout() {
    return Row(
      children: [
        // 左侧品牌展示区域
        Expanded(
          flex: 1,
          child: Container(
            decoration: BoxDecoration(
              gradient: LinearGradient(
                begin: Alignment.topCenter,
                end: Alignment.bottomCenter,
                colors: [
                  Colors.blue.shade600,
                  Colors.indigo.shade600,
                ],
              ),
              borderRadius: const BorderRadius.only(
                topLeft: Radius.circular(24),
                bottomLeft: Radius.circular(24),
              ),
            ),
            child: Padding(
              padding: const EdgeInsets.fromLTRB(40, 48, 40, 40),
              child: _buildBrandingSection(),
            ),
          ),
        ),

        // 右侧登录表单区域
        Expanded(
          flex: 1,
          child: Container(
            padding: const EdgeInsets.symmetric(horizontal: 56, vertical: 48),
            child: SingleChildScrollView(
              child: _buildLoginForm(),
            ),
          ),
        ),
      ],
    );
  }

  Widget _buildBrandingSection({bool compact = false}) {
    // 左侧品牌区域新的视觉：顶部 Logo + 系统名 + 门店/环境标签 + 功能列表 + 底部版权
    final brandTitle = Row(
      crossAxisAlignment: CrossAxisAlignment.center,
      children: [
        Container(
          width: compact ? 46 : 54,
          height: compact ? 46 : 54,
          decoration: BoxDecoration(
            gradient: LinearGradient(
              colors: [Colors.white, Colors.white.withOpacity(.85)],
              begin: Alignment.topLeft,
              end: Alignment.bottomRight,
            ),
            borderRadius: BorderRadius.circular(16),
            boxShadow: [
              BoxShadow(
                color: Colors.black.withOpacity(.15),
                blurRadius: 18,
                offset: const Offset(0, 8),
              )
            ],
          ),
          child: Icon(Icons.storefront_rounded,
              size: compact ? 30 : 34, color: Colors.indigo.shade600),
        ),
        const SizedBox(width: 14),
        Expanded(
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              Text(
                AppTexts.appName,
                style: TextStyle(
                  fontSize: compact ? 18 : 22,
                  fontWeight: FontWeight.w700,
                  letterSpacing: .5,
                  color: Colors.white,
                ),
              ),
              const SizedBox(height: 4),
              Row(
                children: [
                  _tag('API'),
                  const SizedBox(width: 6),
                  // 门店名称占位（登录前可展示“未登录”或主品牌）
                  Flexible(
                    child: Text(
                      AppTexts.defaultStoreName,
                      style: TextStyle(
                        fontSize: compact ? 11 : 12,
                        color: Colors.white.withOpacity(.9),
                        fontWeight: FontWeight.w500,
                        overflow: TextOverflow.ellipsis,
                      ),
                    ),
                  ),
                ],
              )
            ],
          ),
        )
      ],
    );

    final featureList = !compact
        ? Padding(
            padding: const EdgeInsets.only(top: 36),
            child: _buildFeaturesList(),
          )
        : const SizedBox.shrink();

    final footer = !compact
        ? Padding(
            padding: const EdgeInsets.only(top: 48),
            child: Opacity(
              opacity: .85,
              child: Text(
                '© 2025 Tower Suite\nAll Rights Reserved',
                style: TextStyle(
                  fontSize: 11,
                  height: 1.4,
                  color: Colors.white.withOpacity(.85),
                  letterSpacing: .5,
                ),
              ),
            ),
          )
        : const SizedBox.shrink();

    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        brandTitle,
        const SizedBox(height: 28),
        Text(
          '智能化餐饮管理平台',
          style: TextStyle(
            fontSize: compact ? 16 : 18,
            fontWeight: FontWeight.w600,
            color: Colors.white,
            letterSpacing: .8,
          ),
        ),
        const SizedBox(height: 10),
        Text(
          '聚合门店 · 人员 · 菜品 · 库存 · 经营数据\n助力高效决策与增长',
          style: TextStyle(
            fontSize: compact ? 12 : 13.5,
            height: 1.5,
            color: Colors.white.withOpacity(.90),
          ),
        ),
        featureList,
        const Spacer(),
        footer,
      ],
    );
  }

  Widget _tag(String text) {
    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 2),
      decoration: BoxDecoration(
        color: Colors.white.withOpacity(.15),
        borderRadius: BorderRadius.circular(6),
        border: Border.all(color: Colors.white.withOpacity(.4), width: .8),
      ),
      child: Text(
        text,
        style: const TextStyle(
          fontSize: 11,
          fontWeight: FontWeight.w600,
          letterSpacing: .5,
          color: Colors.white,
        ),
      ),
    );
  }

  Widget _buildFeaturesList() {
    final features = [
      {'icon': Icons.inventory_2_outlined, 'text': '库存管理'},
      {'icon': Icons.people_outline, 'text': '员工管理'},
      {'icon': Icons.assessment_outlined, 'text': '数据分析'},
      {'icon': Icons.store_outlined, 'text': '多门店支持'},
    ];

    return Column(
      children: features.map((feature) {
        return Padding(
          padding: const EdgeInsets.symmetric(vertical: 8),
          child: Row(
            children: [
              Icon(
                feature['icon'] as IconData,
                color: Colors.white,
                size: 24,
              ),
              const SizedBox(width: 16),
              Text(
                feature['text'] as String,
                style: const TextStyle(
                  color: Colors.white,
                  fontSize: 16,
                ),
              ),
            ],
          ),
        );
      }).toList(),
    );
  }

  Widget _buildLoginForm() {
    return Form(
      key: _formKey,
      child: Column(
        mainAxisSize: MainAxisSize.min,
        crossAxisAlignment: CrossAxisAlignment.stretch,
        children: [
          // 标题
          Text(
            '欢迎回来',
            style: TextStyle(
              fontSize: 28,
              fontWeight: FontWeight.bold,
              color: Colors.grey.shade800,
            ),
          ),
          const SizedBox(height: 8),
          Text(
            '请登录您的账户',
            style: TextStyle(
              fontSize: 16,
              color: Colors.grey.shade600,
            ),
          ),
          const SizedBox(height: 40),

          // 手机号输入框
          TextFormField(
            controller: _phoneCtrl,
            keyboardType: TextInputType.phone,
            decoration: InputDecoration(
              labelText: '手机号',
              hintText: '请输入11位手机号',
              prefixIcon:
                  Icon(Icons.phone_android, color: Colors.blue.shade600),
              border: OutlineInputBorder(
                borderRadius: BorderRadius.circular(12),
              ),
              enabledBorder: OutlineInputBorder(
                borderRadius: BorderRadius.circular(12),
                borderSide: BorderSide(color: Colors.grey.shade300),
              ),
              focusedBorder: OutlineInputBorder(
                borderRadius: BorderRadius.circular(12),
                borderSide: BorderSide(color: Colors.blue.shade600, width: 2),
              ),
              filled: true,
              fillColor: Colors.grey.shade50,
            ),
            validator: (value) {
              if (value == null || value.trim().isEmpty) {
                return '请输入手机号';
              }
              if (value.trim().length != 11) {
                return '请输入11位手机号';
              }
              return null;
            },
          ),
          const SizedBox(height: 20),

          // 密码输入框
          TextFormField(
            controller: _pwdCtrl,
            obscureText: _obscurePassword,
            decoration: InputDecoration(
              labelText: '密码',
              hintText: '请输入密码',
              prefixIcon: Icon(Icons.lock_outline, color: Colors.blue.shade600),
              suffixIcon: IconButton(
                icon: Icon(
                  _obscurePassword ? Icons.visibility_off : Icons.visibility,
                  color: Colors.grey.shade600,
                ),
                onPressed: () {
                  setState(() {
                    _obscurePassword = !_obscurePassword;
                  });
                },
              ),
              border: OutlineInputBorder(
                borderRadius: BorderRadius.circular(12),
              ),
              enabledBorder: OutlineInputBorder(
                borderRadius: BorderRadius.circular(12),
                borderSide: BorderSide(color: Colors.grey.shade300),
              ),
              focusedBorder: OutlineInputBorder(
                borderRadius: BorderRadius.circular(12),
                borderSide: BorderSide(color: Colors.blue.shade600, width: 2),
              ),
              filled: true,
              fillColor: Colors.grey.shade50,
            ),
            validator: (value) {
              if (value == null || value.isEmpty) {
                return '请输入密码';
              }
              return null;
            },
            onFieldSubmitted: (_) => _doLogin(),
          ),
          const SizedBox(height: 32),

          // 登录按钮
          SizedBox(
            height: 56,
            child: ElevatedButton(
              onPressed: _loading ? null : _doLogin,
              style: ElevatedButton.styleFrom(
                backgroundColor: Colors.blue.shade600,
                foregroundColor: Colors.white,
                disabledBackgroundColor: Colors.grey.shade300,
                elevation: 2,
                shadowColor: Colors.blue.shade200,
                shape: RoundedRectangleBorder(
                  borderRadius: BorderRadius.circular(12),
                ),
              ),
              child: _loading
                  ? const SizedBox(
                      width: 24,
                      height: 24,
                      child: CircularProgressIndicator(
                        strokeWidth: 2.5,
                        color: Colors.white,
                      ),
                    )
                  : const Text(
                      '登录',
                      style: TextStyle(
                        fontSize: 18,
                        fontWeight: FontWeight.w600,
                        letterSpacing: 1,
                      ),
                    ),
            ),
          ),
          const SizedBox(height: 24),

          // 底部提示
          Center(
            child: Text(
              'Tower 餐饮管理系统 v1.0',
              style: TextStyle(
                color: Colors.grey.shade500,
                fontSize: 12,
              ),
            ),
          ),
        ],
      ),
    );
  }
}
