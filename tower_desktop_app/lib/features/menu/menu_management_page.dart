import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import 'package:tdesign_flutter/tdesign_flutter.dart';
import '../auth/permission_gate.dart';
import '../../core/constants/menu_types.dart';
import '../../core/constants/app_constants.dart';
import 'menu_manage_provider.dart';
import 'menu_api.dart';
import 'menu_form_dialog.dart';
import 'models.dart';

/// 菜单管理页面（与用户侧边栏使用的菜单树分开）
class MenuManagementPage extends StatefulWidget {
  const MenuManagementPage({super.key});

  @override
  State<MenuManagementPage> createState() => _MenuManagementPageState();
}

class _MenuManagementPageState extends State<MenuManagementPage> {
  final _searchCtrl = TextEditingController();

  @override
  void initState() {
    super.initState();
    WidgetsBinding.instance.addPostFrameCallback((_) {
      context.read<MenuManageProvider>().load();
    });
  }

  @override
  void dispose() {
    _searchCtrl.dispose();
    super.dispose();
  }

  void _openCreate() async {
    final provider = context.read<MenuManageProvider>();
    final parents = provider
        .flattened()
        .map((n) => n.item)
        .where((m) => m.type == MenuType.directory || m.type == MenuType.page)
        .toList();
    final req = await showDialog<CreateMenuRequest>(
      context: context,
      builder: (_) => MenuFormDialog(parents: parents),
    );
    if (req != null) {
      final ok = await provider.create(req);
      if (ok) {
        TDToast.showSuccess('创建成功', context: context);
      } else {
        TDToast.showFail(provider.error ?? '创建失败', context: context);
      }
    }
  }

  void _openEdit(MenuItem item) async {
    final provider = context.read<MenuManageProvider>();
    final parents = provider
        .flattened()
        .map((n) => n.item)
        .where((m) =>
            m.id != item.id &&
            (m.type == MenuType.directory || m.type == MenuType.page))
        .toList();
    final req = await showDialog<UpdateMenuRequest>(
      context: context,
      builder: (_) => MenuFormDialog(editing: item, parents: parents),
    );
    if (req != null) {
      final ok = await provider.update(item.id, req);
      if (ok) {
        TDToast.showSuccess('更新成功', context: context);
      } else {
        TDToast.showFail(provider.error ?? '更新失败', context: context);
      }
    }
  }

  void _delete(MenuItem item) async {
    final provider = context.read<MenuManageProvider>();
    final confirm = await showDialog<bool>(
      context: context,
      builder: (ctx) => AlertDialog(
        title: const Text('确认删除'),
        content: Text('确定要删除菜单 "${item.title}" 吗?'),
        actions: [
          TextButton(
              onPressed: () => Navigator.pop(ctx, false),
              child: const Text('取消')),
          TDButton(
              text: '删除',
              theme: TDButtonTheme.danger,
              onTap: () => Navigator.pop(ctx, true)),
        ],
      ),
    );
    if (confirm == true) {
      final ok = await provider.remove(item.id);
      if (ok) {
        TDMessage.showMessage(
          context: context,
          content: '删除成功',
          theme: MessageTheme.success,
        );
      } else {
        TDToast.showFail(provider.error ?? '删除失败', context: context);
      }
    }
  }

  void _openCreateChild(MenuItem parent) async {
    final provider = context.read<MenuManageProvider>();
    final parents = provider
        .flattened()
        .map((n) => n.item)
        .where((m) => m.type == MenuType.directory || m.type == MenuType.page)
        .toList();
    final req = await showDialog<CreateMenuRequest>(
      context: context,
      builder: (_) => MenuFormDialog(
        parents: parents,
        initialParentId: parent.id,
      ),
    );
    if (req != null) {
      final ok = await provider.create(req);
      if (ok) {
        TDToast.showSuccess('创建成功', context: context);
        context.read<MenuManageProvider>().toggleExpand(parent.id); // 展开新父级
      } else {
        TDToast.showFail(provider.error ?? '创建失败', context: context);
      }
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: Colors.grey.shade50,
      body: Column(
        children: [
          _buildToolbar(),
          Expanded(child: _buildTreeTable()),
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
          const Text('菜单管理',
              style: TextStyle(fontSize: 20, fontWeight: FontWeight.w600)),
          const Spacer(),
          SizedBox(
            width: 260,
            child: TextField(
              controller: _searchCtrl,
              decoration: InputDecoration(
                hintText: '搜索标题/名称/权限码',
                prefixIcon: const Icon(Icons.search, size: 20),
                filled: true,
                fillColor: Colors.grey.shade50,
                border: OutlineInputBorder(
                  borderRadius: BorderRadius.circular(8),
                  borderSide: BorderSide.none,
                ),
                contentPadding:
                    const EdgeInsets.symmetric(horizontal: 16, vertical: 12),
              ),
              onSubmitted: (_) => context.read<MenuManageProvider>().load(),
            ),
          ),
          const SizedBox(width: 12),
          TDButton(
            text: '搜索',
            theme: TDButtonTheme.primary,
            size: TDButtonSize.medium,
            onTap: () => context.read<MenuManageProvider>().load(),
          ),
          const SizedBox(width: 12),
          PermissionGate(
            required: PermissionCodes.menuAdd,
            child: TDButton(
              text: '新增菜单',
              theme: TDButtonTheme.primary,
              size: TDButtonSize.medium,
              icon: TDIcons.add,
              onTap: _openCreate,
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildTreeTable() {
    return Consumer<MenuManageProvider>(builder: (context, provider, _) {
      if (provider.loading) {
        return const Center(child: CircularProgressIndicator());
      }
      if (provider.error != null) {
        return Center(
          child: Column(
            mainAxisSize: MainAxisSize.min,
            children: [
              Icon(Icons.error_outline, size: 64, color: Colors.red.shade300),
              const SizedBox(height: 12),
              Text(provider.error!,
                  style: TextStyle(color: Colors.red.shade700)),
              const SizedBox(height: 12),
              TDButton(
                text: '重试',
                theme: TDButtonTheme.primary,
                onTap: () => provider.load(),
              )
            ],
          ),
        );
      }
      if (provider.tree.isEmpty) {
        return Center(
          child: Column(
            mainAxisSize: MainAxisSize.min,
            children: [
              Icon(Icons.inbox_outlined, size: 80, color: Colors.grey.shade300),
              const SizedBox(height: 16),
              const Text('暂无菜单',
                  style: TextStyle(fontSize: 16, color: Colors.grey)),
            ],
          ),
        );
      }

      final flat = provider.flattened();
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
          child: Column(
            children: [
              _buildHeaderRow(),
              Divider(height: 1, color: Colors.grey.shade200),
              Expanded(
                child: ListView.separated(
                  itemCount: flat.length,
                  separatorBuilder: (_, __) =>
                      Divider(height: 1, color: Colors.grey.shade100),
                  itemBuilder: (ctx, i) {
                    final node = flat[i];
                    final m = node.item;
                    final hasChildren = node.hasChildren;
                    final expanded = provider.isExpanded(m.id);
                    final isDir = m.type == MenuType.directory;
                    final isPage = m.type == MenuType.page;
                    // 按钮行当前未做单独逻辑，如需后续单独样式可使用 isButton 判断
                    // final isButton = m.type == MenuType.button;
                    final rowColor = isDir
                        ? Colors.grey.shade100
                        : (i.isEven ? Colors.white : Colors.grey.shade50);
                    final titleStyle = TextStyle(
                      fontSize: 14,
                      fontWeight: isDir
                          ? FontWeight.w600
                          : isPage
                              ? FontWeight.w500
                              : FontWeight.normal,
                    );
                    return Container(
                      color: rowColor,
                      padding: const EdgeInsets.symmetric(
                          horizontal: 12, vertical: 8),
                      child: Row(
                        children: [
                          // 菜单名称列 (含展开箭头与缩进)
                          SizedBox(
                            width: 260,
                            child: Row(
                              children: [
                                SizedBox(width: node.depth * 16.0),
                                if (hasChildren)
                                  GestureDetector(
                                    onTap: () => provider.toggleExpand(m.id),
                                    child: Icon(
                                      expanded
                                          ? Icons.keyboard_arrow_down
                                          : Icons.chevron_right,
                                      size: 18,
                                      color: Colors.blueGrey.shade600,
                                    ),
                                  )
                                else
                                  const SizedBox(width: 18),
                                const SizedBox(width: 4),
                                Icon(
                                  m.type == MenuType.directory
                                      ? Icons.folder_outlined
                                      : m.type == MenuType.page
                                          ? Icons.insert_drive_file_outlined
                                          : Icons.extension_outlined,
                                  size: 16,
                                  color: Colors.blueGrey.shade500,
                                ),
                                const SizedBox(width: 6),
                                Expanded(
                                  child: Text(
                                    m.title,
                                    maxLines: 1,
                                    overflow: TextOverflow.ellipsis,
                                    style: titleStyle,
                                  ),
                                ),
                              ],
                            ),
                          ),
                          // 图标列
                          SizedBox(
                            width: 80,
                            child: Text(
                              m.icon ?? '-',
                              maxLines: 1,
                              overflow: TextOverflow.ellipsis,
                              style: TextStyle(
                                  fontSize: 12, color: Colors.grey.shade700),
                            ),
                          ),
                          // 排序列
                          SizedBox(
                            width: 60,
                            child: Text('${m.sort ?? '-'}',
                                textAlign: TextAlign.center,
                                style: TextStyle(
                                    fontSize: 12, color: Colors.grey.shade800)),
                          ),
                          // 权限标识列
                          SizedBox(
                            width: 160,
                            child: Text(
                              m.permission ?? '-',
                              maxLines: 1,
                              overflow: TextOverflow.ellipsis,
                              style: TextStyle(
                                  fontSize: 12, color: Colors.teal.shade700),
                            ),
                          ),
                          // 组件路径列
                          SizedBox(
                            width: 180,
                            child: Text(
                              m.component ?? '-',
                              maxLines: 1,
                              overflow: TextOverflow.ellipsis,
                              style: TextStyle(
                                  fontSize: 12, color: Colors.indigo.shade600),
                            ),
                          ),
                          // 状态列
                          SizedBox(
                            width: 80,
                            child: _buildStatusTag(m.status),
                          ),
                          // 创建时间列
                          SizedBox(
                            width: 160,
                            child: Text(
                              m.createdAt ?? '-',
                              maxLines: 1,
                              overflow: TextOverflow.ellipsis,
                              style: TextStyle(
                                  fontSize: 12, color: Colors.grey.shade600),
                            ),
                          ),
                          // 操作列
                          Expanded(
                            child: Row(
                              mainAxisAlignment: MainAxisAlignment.end,
                              children: [
                                PermissionGate(
                                  required: PermissionCodes.menuEdit,
                                  child: TextButton(
                                    onPressed: () => _openEdit(m),
                                    child: const Text('修改'),
                                  ),
                                ),
                                PermissionGate(
                                  required: PermissionCodes.menuAdd,
                                  child: TextButton(
                                    onPressed: () => _openCreateChild(m),
                                    child: const Text('+新增'),
                                  ),
                                ),
                                PermissionGate(
                                  required: PermissionCodes.menuDelete,
                                  child: TextButton(
                                    onPressed: () => _delete(m),
                                    child: const Text('删除',
                                        style:
                                            TextStyle(color: Colors.redAccent)),
                                  ),
                                ),
                              ],
                            ),
                          ),
                        ],
                      ),
                    );
                  },
                ),
              ),
            ],
          ),
        ),
      );
    });
  }

  Widget _buildHeaderRow() {
    Text buildHeader(String text) => Text(
          text,
          style: const TextStyle(
              fontSize: 12, fontWeight: FontWeight.w600, color: Colors.black87),
        );
    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 10),
      child: Row(
        children: [
          SizedBox(width: 260, child: buildHeader('菜单名称')),
          SizedBox(width: 80, child: buildHeader('图标')),
          SizedBox(width: 60, child: buildHeader('排序')),
          SizedBox(width: 160, child: buildHeader('权限标识')),
          SizedBox(width: 180, child: buildHeader('组件路径')),
          SizedBox(width: 80, child: buildHeader('状态')),
          SizedBox(width: 160, child: buildHeader('创建时间')),
          const Expanded(
              child: Align(
                  alignment: Alignment.centerRight,
                  child: Text('操作',
                      style: TextStyle(
                          fontSize: 12, fontWeight: FontWeight.w600))))
        ],
      ),
    );
  }

  Widget _buildStatusTag(int? status) {
    final active = (status ?? 1) == 1;
    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 10, vertical: 4),
      decoration: BoxDecoration(
        color: active ? Colors.blue.shade50 : Colors.grey.shade200,
        borderRadius: BorderRadius.circular(4),
        border: Border.all(
            color: active ? Colors.blue.shade200 : Colors.grey.shade400),
      ),
      child: Text(
        active ? '正常' : '停用',
        textAlign: TextAlign.center,
        style: TextStyle(
          fontSize: 12,
          fontWeight: FontWeight.w500,
          color: active ? Colors.blue.shade700 : Colors.grey.shade600,
        ),
      ),
    );
  }

  // 分页组件已移除（树结构不分页）
}

/// 将页面放入 MultiProvider 时的便捷构造
class MenuManagementScope extends StatelessWidget {
  const MenuManagementScope({super.key});

  @override
  Widget build(BuildContext context) {
    return ChangeNotifierProvider(
      create: (_) => MenuManageProvider(MenuApi()),
      child: const MenuManagementPage(),
    );
  }
}
