import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import 'package:tdesign_flutter/tdesign_flutter.dart';
import 'user_provider.dart';
import 'models.dart';
import 'user_form_dialog.dart';
import '../../features/auth/permission_gate.dart';

class UserManagementPage extends StatefulWidget {
  const UserManagementPage({super.key});

  @override
  State<UserManagementPage> createState() => _UserManagementPageState();
}

class _UserManagementPageState extends State<UserManagementPage> {
  final TextEditingController _searchController = TextEditingController();

  @override
  void initState() {
    super.initState();
    WidgetsBinding.instance.addPostFrameCallback((_) {
      context.read<UserProvider>().loadUsers();
    });
  }

  @override
  void dispose() {
    _searchController.dispose();
    super.dispose();
  }

  void _handleSearch() {
    context.read<UserProvider>().loadUsers(
          page: 1,
          keyword: _searchController.text.trim(),
        );
  }

  void _handleCreate() async {
    final result = await showDialog<CreateUserRequest>(
      context: context,
      builder: (ctx) => const UserFormDialog(),
    );
    if (result != null && mounted) {
      final success = await context.read<UserProvider>().createUser(result);
      if (success && mounted) {
        TDToast.showSuccess('创建成功', context: context);
      } else if (mounted) {
        final err = context.read<UserProvider>().error ?? '创建失败';
        TDToast.showFail(err, context: context);
      }
    }
  }

  void _handleDelete(User user) async {
    final confirm = await showDialog<bool>(
      context: context,
      builder: (ctx) => AlertDialog(
        title: const Text('确认删除'),
        content: Text('确定要删除用户 "${user.username}" 吗?'),
        actions: [
          TextButton(
            onPressed: () => Navigator.pop(ctx, false),
            child: const Text('取消'),
          ),
          TDButton(
            text: '删除',
            theme: TDButtonTheme.danger,
            onTap: () => Navigator.pop(ctx, true),
          ),
        ],
      ),
    );
    if (confirm == true && mounted) {
      final success = await context.read<UserProvider>().deleteUser(user.id);
      if (success && mounted) {
        TDMessage.showMessage(
          context: context,
          content: '删除成功',
          theme: MessageTheme.success,
          duration: 3000,
        );
      } else if (mounted) {
        final err = context.read<UserProvider>().error ?? '删除失败';
        TDToast.showFail(err, context: context);
      }
    }
  }

  void _handleStatusChange(User user, bool value) async {
    final newStatus = value ? 1 : 0;
    final req = UpdateUserRequest(
      status: newStatus,
    );
    final success = await context.read<UserProvider>().updateUser(user.id, req);
    if (success && mounted) {
      TDMessage.showMessage(
        context: context,
        content: '状态已更新',
        theme: MessageTheme.success,
        duration: 2000,
      );
    } else if (mounted) {
      final err = context.read<UserProvider>().error ?? '状态更新失败';
      TDToast.showFail(err, context: context);
    }
  }

  void _handleResetPassword(User user) async {
    showGeneralDialog(
      context: context,
      pageBuilder: (BuildContext buildContext, Animation<double> animation,
          Animation<double> secondaryAnimation) {
        return TDAlertDialog(
          content: '确定要重置 ${user.username} 的密码吗？',
          buttonWidget: Row(
            children: [
              Expanded(
                child: TDButton(
                  text: '取消',
                  size: TDButtonSize.large,
                  type: TDButtonType.outline,
                  theme: TDButtonTheme.defaultTheme,
                  onTap: () => Navigator.of(buildContext).pop(),
                ),
              ),
              const SizedBox(width: 16),
              Expanded(
                child: TDButton(
                  text: '确定',
                  size: TDButtonSize.large,
                  theme: TDButtonTheme.primary,
                  onTap: () async {
                    Navigator.of(buildContext).pop();
                    // 生成默认密码（可以根据需求调整）
                    const defaultPassword = '123456';
                    if (mounted) {
                      final provider = context.read<UserProvider>();
                      final success = await provider.resetPassword(
                          user.id, defaultPassword);
                      if (success && mounted) {
                        TDMessage.showMessage(
                          context: context,
                          content: '密码已重置为：$defaultPassword',
                          theme: MessageTheme.success,
                          duration: 3000,
                        );
                      } else if (mounted) {
                        TDMessage.showMessage(
                          context: context,
                          content: provider.error ?? '重置失败',
                          theme: MessageTheme.error,
                          duration: 3000,
                        );
                      }
                    }
                  },
                ),
              ),
            ],
          ),
        );
      },
    );
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: Colors.grey.shade50,
      body: Column(
        children: [
          _buildToolbar(),
          Expanded(child: _buildTable()),
          _buildPagination(),
        ],
      ),
    );
  }

  Widget _buildToolbar() {
    return Container(
      padding: const EdgeInsets.all(20),
      decoration: BoxDecoration(
        color: Colors.white,
        boxShadow: [
          BoxShadow(
            color: Colors.black.withOpacity(0.05),
            blurRadius: 4,
            offset: const Offset(0, 2),
          ),
        ],
      ),
      child: Row(
        children: [
          const Text(
            '用户管理',
            style: TextStyle(fontSize: 20, fontWeight: FontWeight.w600),
          ),
          const Spacer(),
          SizedBox(
            width: 280,
            child: TextField(
              controller: _searchController,
              decoration: InputDecoration(
                hintText: '搜索用户名/手机号后4位',
                prefixIcon: const Icon(Icons.search, size: 20),
                filled: true,
                fillColor: Colors.grey.shade50,
                border: OutlineInputBorder(
                  borderRadius: BorderRadius.circular(8),
                  borderSide: BorderSide.none,
                ),
                contentPadding: const EdgeInsets.symmetric(
                  horizontal: 16,
                  vertical: 12,
                ),
              ),
              onSubmitted: (_) => _handleSearch(),
            ),
          ),
          const SizedBox(width: 12),
          TDButton(
            text: '搜索',
            theme: TDButtonTheme.primary,
            size: TDButtonSize.medium,
            onTap: _handleSearch,
          ),
          const SizedBox(width: 12),
          PermissionGate(
            required: 'system:user:add',
            child: TDButton(
              text: '新增用户',
              theme: TDButtonTheme.primary,
              size: TDButtonSize.medium,
              icon: TDIcons.add,
              onTap: _handleCreate,
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildTable() {
    return Consumer<UserProvider>(
      builder: (context, provider, _) {
        if (provider.loading) {
          return _buildTableSkeleton(rows: provider.pageSize);
        }
        if (provider.error != null) {
          return Center(
            child: Column(
              mainAxisSize: MainAxisSize.max,
              children: [
                Icon(Icons.error_outline, size: 64, color: Colors.red.shade300),
                const SizedBox(height: 16),
                Text('加载失败: ${provider.error}',
                    style: TextStyle(color: Colors.red.shade700)),
                const SizedBox(height: 16),
                TDButton(
                  text: '重试',
                  theme: TDButtonTheme.primary,
                  onTap: () => provider.loadUsers(),
                ),
              ],
            ),
          );
        }

        final columns = <TDTableCol>[
          TDTableCol(title: '用户名', colKey: 'username', ellipsis: true),
          TDTableCol(title: '手机号', colKey: 'phone', ellipsis: true),
          TDTableCol(
            title: '性别',
            colKey: 'gender',
            width: 80,
            cellBuilder: (ctx, index) {
              final u = provider.users[index];
              if (u.gender == 1) {
                // 男
                return Row(
                  mainAxisSize: MainAxisSize.min,
                  children: [
                    Icon(Icons.male, size: 18, color: Colors.blue.shade600),
                    const SizedBox(width: 4),
                    const Text('男', style: TextStyle(fontSize: 14)),
                  ],
                );
              } else if (u.gender == 2) {
                // 女
                return Row(
                  mainAxisSize: MainAxisSize.min,
                  children: [
                    Icon(Icons.female, size: 18, color: Colors.pink.shade400),
                    const SizedBox(width: 4),
                    const Text('女', style: TextStyle(fontSize: 14)),
                  ],
                );
              } else {
                return const Text('—', style: TextStyle(fontSize: 14));
              }
            },
          ),
          TDTableCol(title: '角色', colKey: 'role'),
          TDTableCol(title: '门店', colKey: 'store'),
          TDTableCol(
            title: '状态',
            colKey: 'status',
            width: 90,
            cellBuilder: (ctx, index) {
              final u = provider.users[index];
              final isActive = u.status == 1;
              return TDSwitch(
                isOn: isActive,
                type: TDSwitchType.icon,
                onChanged: (value) {
                  _handleStatusChange(u, value);
                  return value;
                },
              );
            },
          ),
          TDTableCol(
            title: '操作',
            colKey: 'actions',
            width: 150,
            cellBuilder: (ctx, index) {
              final u = provider.users[index];
              return Row(
                mainAxisAlignment: MainAxisAlignment.start,
                children: [
                  PermissionGate(
                    required: 'system:user:delete',
                    child: GestureDetector(
                      onTap: () => _handleDelete(u),
                      child: Padding(
                        padding: const EdgeInsets.symmetric(
                            horizontal: 4, vertical: 4),
                        child: Text(
                          '删除',
                          style: const TextStyle(
                            color: Colors.redAccent,
                            fontSize: 14,
                          ),
                        ),
                      ),
                    ),
                  ),
                  PermissionGate(
                    required: 'system:user:edit',
                    child: GestureDetector(
                      onTap: () => _handleResetPassword(u),
                      child: Padding(
                        padding: const EdgeInsets.symmetric(
                            horizontal: 4, vertical: 4),
                        child: Text(
                          '重置密码',
                          style: TextStyle(
                            color: TDTheme.of(context).warningColor5,
                            fontSize: 14,
                          ),
                        ),
                      ),
                    ),
                  ),
                ],
              );
            },
          ),
        ];

        final data = provider.users
            .map((u) => {
                  'username': u.username,
                  'phone': u.phone,
                  'gender': '', // 占位，实际通过 cellBuilder
                  'role': u.roleName ?? '—',
                  'store': u.storeName ?? '—',
                  'status': u.status == 1 ? '启用' : '禁用',
                  'actions': '', // 占位，实际通过 cellBuilder
                })
            .toList();

        return Padding(
          padding: const EdgeInsets.all(20),
          child: Container(
            decoration: BoxDecoration(
              color: Colors.white,
              borderRadius: BorderRadius.circular(8),
              boxShadow: [
                BoxShadow(
                  color: Colors.black.withOpacity(0.05),
                  blurRadius: 6,
                  offset: const Offset(0, 2),
                ),
              ],
            ),
            child: data.isEmpty
                ? Center(
                    child: Column(
                      mainAxisSize: MainAxisSize.min,
                      children: [
                        Icon(Icons.inbox_outlined,
                            size: 80, color: Colors.grey.shade300),
                        const SizedBox(height: 16),
                        const Text(
                          '暂无用户',
                          style: TextStyle(fontSize: 16, color: Colors.grey),
                        ),
                      ],
                    ),
                  )
                : ClipRRect(
                    borderRadius: BorderRadius.circular(8),
                    child: TDTable(
                      columns: columns,
                      data: data,
                      rowHeight: 50,
                    ),
                  ),
          ),
        );
      },
    );
  }

  Widget _buildTableSkeleton({int rows = 8}) {
    if (rows <= 0) rows = 1;
    if (rows > 50) rows = 50; // 安全上限
    // 简单骨架：模拟 8 行，列布局与真实表头对应
    return Padding(
      padding: const EdgeInsets.all(20),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          // 表头骨架
          Container(
            padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 12),
            decoration: BoxDecoration(
              color: Colors.white,
              borderRadius: BorderRadius.circular(8),
              boxShadow: [
                BoxShadow(
                  color: Colors.black.withOpacity(0.04),
                  blurRadius: 4,
                  offset: const Offset(0, 2),
                ),
              ],
            ),
            child: Row(
              children: [
                _skeletonBlock(width: 80, height: 14),
                const SizedBox(width: 40),
                _skeletonBlock(width: 80, height: 14),
                const SizedBox(width: 40),
                _skeletonBlock(width: 60, height: 14),
                const SizedBox(width: 40),
                _skeletonBlock(width: 60, height: 14),
                const Spacer(),
                _skeletonBlock(width: 50, height: 14),
              ],
            ),
          ),
          const SizedBox(height: 8),
          Expanded(
            child: Container(
              decoration: BoxDecoration(
                color: Colors.white,
                borderRadius: BorderRadius.circular(8),
                boxShadow: [
                  BoxShadow(
                    color: Colors.black.withOpacity(0.04),
                    blurRadius: 4,
                    offset: const Offset(0, 2),
                  ),
                ],
              ),
              child: ListView.separated(
                itemCount: rows,
                separatorBuilder: (_, __) => Divider(
                  height: 1,
                  color: Colors.grey.shade100,
                ),
                itemBuilder: (ctx, index) {
                  return AnimatedContainer(
                    duration: const Duration(milliseconds: 800),
                    curve: Curves.easeInOut,
                    color: index.isEven ? Colors.white : Colors.grey.shade50,
                    padding: const EdgeInsets.symmetric(
                        horizontal: 20, vertical: 18), // 增大骨架行高度
                    child: Row(
                      children: [
                        _skeletonBlock(width: 150, height: 16),
                        const SizedBox(width: 48),
                        _skeletonBlock(width: 120, height: 16),
                        const SizedBox(width: 48),
                        _skeletonBlock(width: 80, height: 16),
                        const SizedBox(width: 48),
                        _skeletonBlock(width: 80, height: 16),
                        const Spacer(),
                        Row(
                          children: [
                            _skeletonBlock(width: 40, height: 16),
                            const SizedBox(width: 20),
                            _skeletonBlock(width: 40, height: 16),
                          ],
                        )
                      ],
                    ),
                  );
                },
              ),
            ),
          ),
        ],
      ),
    );
  }

  Widget _skeletonBlock({required double width, required double height}) {
    return TweenAnimationBuilder<double>(
      tween: Tween(begin: 0, end: 1),
      duration: const Duration(milliseconds: 1200),
      curve: Curves.easeInOut,
      builder: (context, value, child) {
        final base = Colors.grey.shade200;
        final highlight = Colors.grey.shade300;
        Color blended(Color a, Color b, double t) => Color.lerp(a, b, t)!;
        return Container(
          width: width,
          height: height,
          decoration: BoxDecoration(
            color: blended(base, highlight, (value * 2 % 1).clamp(0, 1)),
            borderRadius: BorderRadius.circular(4),
          ),
        );
      },
    );
  }

  Widget _buildPagination() {
    return Consumer<UserProvider>(
      builder: (context, provider, _) {
        final totalPages = (provider.total / provider.pageSize).ceil();
        return Container(
          padding: const EdgeInsets.symmetric(horizontal: 20, vertical: 16),
          decoration: BoxDecoration(
            color: Colors.white,
            boxShadow: [
              BoxShadow(
                color: Colors.black.withOpacity(0.05),
                blurRadius: 4,
                offset: const Offset(0, -2),
              ),
            ],
          ),
          child: Row(
            mainAxisAlignment: MainAxisAlignment.spaceBetween,
            children: [
              Text(
                '共 ${provider.total} 条记录',
                style: TextStyle(fontSize: 14, color: Colors.grey.shade700),
              ),
              Row(
                children: [
                  IconButton(
                    icon: const Icon(Icons.chevron_left),
                    onPressed: provider.page > 1
                        ? () => provider.loadUsers(page: provider.page - 1)
                        : null,
                    color: provider.page > 1
                        ? Colors.blue.shade700
                        : Colors.grey.shade400,
                  ),
                  Container(
                    padding:
                        const EdgeInsets.symmetric(horizontal: 16, vertical: 8),
                    decoration: BoxDecoration(
                      color: Colors.grey.shade100,
                      borderRadius: BorderRadius.circular(6),
                    ),
                    child: Text(
                      '${provider.page} / ${totalPages > 0 ? totalPages : 1}',
                      style: const TextStyle(
                          fontSize: 14, fontWeight: FontWeight.w500),
                    ),
                  ),
                  IconButton(
                    icon: const Icon(Icons.chevron_right),
                    onPressed: provider.page < totalPages
                        ? () => provider.loadUsers(page: provider.page + 1)
                        : null,
                    color: provider.page < totalPages
                        ? Colors.blue.shade700
                        : Colors.grey.shade400,
                  ),
                ],
              ),
            ],
          ),
        );
      },
    );
  }
}
