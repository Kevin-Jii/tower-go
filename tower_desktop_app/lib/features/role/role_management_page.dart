import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import '../auth/permission_gate.dart';
import '../../core/constants/app_constants.dart';
import '../../core/constants/ui_texts.dart';
import '../../core/network/api_client.dart';
import '../../core/widgets/date_text.dart';
import '../../core/widgets/status_tag.dart';
import '../../core/widgets/admin_table.dart';
import 'role_provider.dart';
import 'role_api.dart';
import 'role_form_dialog.dart';
import 'role_models.dart';

class RoleManagementPage extends StatefulWidget {
  const RoleManagementPage({super.key});
  @override
  State<RoleManagementPage> createState() => _RoleManagementPageState();
}

class _RoleManagementPageState extends State<RoleManagementPage> {
  final TextEditingController _searchCtrl = TextEditingController();
  @override
  void initState() {
    super.initState();
    WidgetsBinding.instance.addPostFrameCallback((_) {
      context.read<RoleProvider>().loadRoles();
    });
  }

  @override
  void dispose() {
    _searchCtrl.dispose();
    super.dispose();
  }

  void _openCreate() async {
    final req = await showDialog<CreateRoleRequest>(
      context: context,
      // 重新在对话框作用域内提供 RoleProvider，避免在某些布局/多 Navigator 场景下
      // showDialog 创建的 route 脱离原有 InheritedElement 链导致的 Provider 未找到异常。
      builder: (_) => ChangeNotifierProvider.value(
        value: context.read<RoleProvider>(),
        child: const RoleFormDialog(),
      ),
    );
    if (req != null) {
      final ok = await context.read<RoleProvider>().create(req);
      if (!mounted) return;
      ScaffoldMessenger.of(context).showSnackBar(SnackBar(
        content:
            Text(ok ? UITexts.roleCreateSuccess : UITexts.roleCreateFailed),
      ));
    }
  }

  void _openEdit(RoleItem r) async {
    final req = await showDialog<UpdateRoleRequest>(
      context: context,
      builder: (_) => ChangeNotifierProvider.value(
        value: context.read<RoleProvider>(),
        child: RoleFormDialog(editing: r),
      ),
    );
    if (req != null) {
      final ok = await context.read<RoleProvider>().update(r.id, req);
      if (!mounted) return;
      ScaffoldMessenger.of(context).showSnackBar(SnackBar(
        content:
            Text(ok ? UITexts.roleUpdateSuccess : UITexts.roleUpdateFailed),
      ));
    }
  }

  void _delete(RoleItem r) async {
    final confirm = await showDialog<bool>(
      context: context,
      builder: (ctx) => AlertDialog(
        title: const Text(UITexts.roleDeleteConfirmTitle),
        content:
            Text('${UITexts.roleDeleteConfirmContentPrefix} "${r.name}" ?'),
        actions: [
          TextButton(
              onPressed: () => Navigator.pop(ctx, false),
              child: const Text(UITexts.roleDeleteConfirmCancel)),
          ElevatedButton(
              onPressed: () => Navigator.pop(ctx, true),
              child: const Text(UITexts.roleDeleteConfirmOk)),
        ],
      ),
    );
    if (confirm == true) {
      final ok = await context.read<RoleProvider>().remove(r.id);
      if (!mounted) return;
      ScaffoldMessenger.of(context).showSnackBar(SnackBar(
        content:
            Text(ok ? UITexts.roleDeleteSuccess : UITexts.roleDeleteFailed),
      ));
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: Colors.grey.shade50,
      body: Column(
        children: [
          _buildToolbar(),
          Expanded(child: _buildTable()),
        ],
      ),
    );
  }

  Widget _buildToolbar() {
    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 20, vertical: 16),
      decoration: BoxDecoration(
        color: Colors.white,
        boxShadow: [
          BoxShadow(
              color: Colors.black.withOpacity(0.05),
              blurRadius: 4,
              offset: const Offset(0, 2)),
        ],
      ),
      child: Row(
        children: [
          const Text(UITexts.roleTitle,
              style: TextStyle(fontSize: 20, fontWeight: FontWeight.w600)),
          const Spacer(),
          PermissionGate(
            required: PermissionCodes.roleAdd,
            child: ElevatedButton.icon(
              onPressed: _openCreate,
              icon: const Icon(Icons.add),
              label: const Text(UITexts.roleAddButton),
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildTable() {
    return Consumer<RoleProvider>(builder: (context, provider, _) {
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
              ElevatedButton(
                onPressed: () => provider.loadRoles(keyword: provider.keyword),
                child: const Text(UITexts.roleRetryButton),
              )
            ],
          ),
        );
      }
      if (provider.list.isEmpty) {
        return Center(
          child: Column(
            mainAxisSize: MainAxisSize.min,
            children: [
              Icon(Icons.inbox_outlined, size: 80, color: Colors.grey.shade300),
              const SizedBox(height: 16),
              const Text(UITexts.roleEmpty,
                  style: TextStyle(fontSize: 16, color: Colors.grey)),
            ],
          ),
        );
      }
      return Padding(
        padding: const EdgeInsets.all(20),
        child: AdminTable<RoleItem>(
          enableSearch: true,
          searchController: _searchCtrl,
          searchPlaceholder: UITexts.roleSearchPlaceholder,
          onSearch: (kw) => context.read<RoleProvider>().loadRoles(keyword: kw),
          page: provider.page,
          pageSize: provider.pageSize,
          total: provider.total,
          onPageChange: (p) => provider.setPage(p),
          columns: const [
            AdminTableColumn(
                width: _colId,
                label: UITexts.roleColumnId,
                alignment: Alignment.center),
            AdminTableColumn(width: _colName, label: UITexts.roleColumnName),
            AdminTableColumn(width: _colCode, label: UITexts.roleColumnCode),
            AdminTableColumn(
                width: _colStatus,
                label: UITexts.roleColumnStatus,
                alignment: Alignment.center),
            AdminTableColumn(
                width: _colCreated, label: UITexts.roleColumnCreated),
            AdminTableColumn(width: _colDesc, label: UITexts.roleColumnDesc),
            AdminTableColumn(
                label: UITexts.roleColumnActions,
                alignment: Alignment.centerRight),
          ],
          data: provider.currentPageItems,
          rowBuilder: (r, i) {
            final active = (r.status ?? 1) == 1;
            return Row(
              children: [
                SizedBox(
                  width: _colId,
                  child: Text('${r.id}',
                      textAlign: TextAlign.center,
                      style: const TextStyle(fontSize: 12)),
                ),
                SizedBox(
                  width: _colName,
                  child: Text(r.name,
                      maxLines: 1,
                      overflow: TextOverflow.ellipsis,
                      style: const TextStyle(
                          fontSize: 14, fontWeight: FontWeight.w500)),
                ),
                SizedBox(
                  width: _colCode,
                  child: Text(r.code,
                      maxLines: 1,
                      overflow: TextOverflow.ellipsis,
                      style: TextStyle(
                          fontSize: 12, color: Colors.blueGrey.shade600)),
                ),
                SizedBox(
                  width: _colStatus,
                  child: StatusTag(
                      active: active,
                      enabledText: UITexts.roleStatusEnabled,
                      disabledText: UITexts.roleStatusDisabled),
                ),
                SizedBox(
                  width: _colCreated,
                  child: DateText(
                      raw: r.createdAt,
                      style:
                          TextStyle(fontSize: 12, color: Colors.grey.shade500)),
                ),
                SizedBox(
                  width: _colDesc,
                  child: Text(r.description ?? '-',
                      maxLines: 1,
                      overflow: TextOverflow.ellipsis,
                      style:
                          TextStyle(fontSize: 12, color: Colors.grey.shade600)),
                ),
                Expanded(
                  child: Row(
                    mainAxisAlignment: MainAxisAlignment.end,
                    children: [
                      PermissionGate(
                        required: PermissionCodes.roleEdit,
                        child: TextButton(
                          onPressed: () => _openEdit(r),
                          child: const Text(UITexts.roleEditButton),
                        ),
                      ),
                      PermissionGate(
                        required: PermissionCodes.roleDelete,
                        child: TextButton(
                          onPressed: () => _delete(r),
                          child: const Text(UITexts.roleDeleteButton,
                              style: TextStyle(color: Colors.redAccent)),
                        ),
                      ),
                    ],
                  ),
                ),
              ],
            );
          },
        ),
      );
    });
  }

  // 统一列宽与对齐配置
  static const double _colId = 70;
  static const double _colName = 160;
  static const double _colCode = 160;
  static const double _colStatus = 130;
  static const double _colCreated = 170;
  static const double _colDesc = 220;
  // 旧表头、行渲染与状态标签已被 AdminTable + StatusTag 替换
}

class RoleManagementScope extends StatelessWidget {
  const RoleManagementScope({super.key});
  @override
  Widget build(BuildContext context) {
    return ChangeNotifierProvider(
      create: (_) => RoleProvider(RoleApi(ApiClient())),
      child: const RoleManagementPage(),
    );
  }
}
