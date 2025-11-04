import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import '../auth/permission_gate.dart';
import '../../core/constants/app_constants.dart';
import '../../core/network/api_client.dart';
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
  final _keywordCtrl = TextEditingController();
  @override
  void initState() {
    super.initState();
    WidgetsBinding.instance.addPostFrameCallback((_) {
      context.read<RoleProvider>().load();
    });
  }

  @override
  void dispose() {
    _keywordCtrl.dispose();
    super.dispose();
  }

  void _openCreate() async {
    final req = await showDialog<CreateRoleRequest>(
      context: context,
      builder: (_) => const RoleFormDialog(),
    );
    if (req != null) {
      final ok = await context.read<RoleProvider>().create(req);
      if (!mounted) return;
      ScaffoldMessenger.of(context).showSnackBar(SnackBar(
        content: Text(ok ? '创建成功' : '创建失败'),
      ));
    }
  }

  void _openEdit(RoleItem r) async {
    final req = await showDialog<UpdateRoleRequest>(
      context: context,
      builder: (_) => RoleFormDialog(editing: r),
    );
    if (req != null) {
      final ok = await context.read<RoleProvider>().update(r.id, req);
      if (!mounted) return;
      ScaffoldMessenger.of(context).showSnackBar(SnackBar(
        content: Text(ok ? '更新成功' : '更新失败'),
      ));
    }
  }

  void _delete(RoleItem r) async {
    final confirm = await showDialog<bool>(
      context: context,
      builder: (ctx) => AlertDialog(
        title: const Text('确认删除'),
        content: Text('确定删除角色 "${r.name}" ?'),
        actions: [
          TextButton(
              onPressed: () => Navigator.pop(ctx, false),
              child: const Text('取消')),
          ElevatedButton(
              onPressed: () => Navigator.pop(ctx, true),
              child: const Text('删除')),
        ],
      ),
    );
    if (confirm == true) {
      final ok = await context.read<RoleProvider>().remove(r.id);
      if (!mounted) return;
      ScaffoldMessenger.of(context).showSnackBar(SnackBar(
        content: Text(ok ? '删除成功' : '删除失败'),
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
      padding: const EdgeInsets.all(20),
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
          const Text('角色管理',
              style: TextStyle(fontSize: 20, fontWeight: FontWeight.w600)),
          const Spacer(),
          SizedBox(
            width: 200,
            child: TextField(
              controller: _keywordCtrl,
              decoration: InputDecoration(
                hintText: '关键字',
                isDense: true,
                filled: true,
                fillColor: Colors.grey.shade100,
                border: OutlineInputBorder(
                    borderSide: BorderSide.none,
                    borderRadius: BorderRadius.circular(8)),
                prefixIcon: const Icon(Icons.search, size: 18),
              ),
              onSubmitted: (_) => context
                  .read<RoleProvider>()
                  .load(keyword: _keywordCtrl.text.trim()),
            ),
          ),
          const SizedBox(width: 12),
          ElevatedButton(
            onPressed: () => context
                .read<RoleProvider>()
                .load(keyword: _keywordCtrl.text.trim()),
            child: const Text('搜索'),
          ),
          const SizedBox(width: 12),
          PermissionGate(
            required: PermissionCodes.roleAdd,
            child: ElevatedButton.icon(
              onPressed: _openCreate,
              icon: const Icon(Icons.add),
              label: const Text('新增角色'),
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
                onPressed: () =>
                    provider.load(keyword: _keywordCtrl.text.trim()),
                child: const Text('重试'),
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
              const Text('暂无角色',
                  style: TextStyle(fontSize: 16, color: Colors.grey)),
            ],
          ),
        );
      }
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
                  offset: const Offset(0, 2)),
            ],
          ),
          child: Column(
            children: [
              _buildHeaderRow(),
              Divider(height: 1, color: Colors.grey.shade200),
              Expanded(
                child: ListView.separated(
                  itemCount: provider.list.length,
                  separatorBuilder: (_, __) =>
                      Divider(height: 1, color: Colors.grey.shade100),
                  itemBuilder: (ctx, i) {
                    final r = provider.list[i];
                    final active = (r.status ?? 1) == 1;
                    return Container(
                      color: i.isEven ? Colors.white : Colors.grey.shade50,
                      padding: const EdgeInsets.symmetric(
                          horizontal: 12, vertical: 10),
                      child: Row(
                        children: [
                          SizedBox(
                              width: 60,
                              child: Text('${r.id}',
                                  style: const TextStyle(fontSize: 12))),
                          SizedBox(
                              width: 160,
                              child: Text(r.name,
                                  maxLines: 1,
                                  overflow: TextOverflow.ellipsis,
                                  style: const TextStyle(
                                      fontSize: 14,
                                      fontWeight: FontWeight.w500))),
                          SizedBox(
                              width: 160,
                              child: Text(r.code,
                                  maxLines: 1,
                                  overflow: TextOverflow.ellipsis,
                                  style: TextStyle(
                                      fontSize: 12,
                                      color: Colors.blueGrey.shade600))),
                          SizedBox(
                              width: 200,
                              child: Text(r.remark ?? '-',
                                  maxLines: 1,
                                  overflow: TextOverflow.ellipsis,
                                  style: TextStyle(
                                      fontSize: 12,
                                      color: Colors.grey.shade600))),
                          SizedBox(width: 80, child: _buildStatusTag(active)),
                          SizedBox(
                              width: 160,
                              child: Text(r.createdAt ?? '-',
                                  style: TextStyle(
                                      fontSize: 12,
                                      color: Colors.grey.shade500))),
                          Expanded(
                            child: Row(
                              mainAxisAlignment: MainAxisAlignment.end,
                              children: [
                                PermissionGate(
                                  required: PermissionCodes.roleEdit,
                                  child: TextButton(
                                      onPressed: () => _openEdit(r),
                                      child: const Text('修改')),
                                ),
                                PermissionGate(
                                  required: PermissionCodes.roleDelete,
                                  child: TextButton(
                                    onPressed: () => _delete(r),
                                    child: const Text('删除',
                                        style:
                                            TextStyle(color: Colors.redAccent)),
                                  ),
                                ),
                              ],
                            ),
                          )
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
    Text h(String t) => Text(t,
        style: const TextStyle(fontSize: 12, fontWeight: FontWeight.w600));
    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 10),
      child: Row(
        children: [
          SizedBox(width: 60, child: h('ID')),
          SizedBox(width: 160, child: h('角色名称')),
          SizedBox(width: 160, child: h('角色编码')),
          SizedBox(width: 200, child: h('备注')),
          SizedBox(width: 80, child: h('状态')),
          SizedBox(width: 160, child: h('创建时间')),
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

  Widget _buildStatusTag(bool active) {
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
