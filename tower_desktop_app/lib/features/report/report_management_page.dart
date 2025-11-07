import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import '../../core/widgets/admin_table.dart';
import '../auth/permission_gate.dart';
import 'report_provider.dart';
import 'models.dart';

class ReportManagementPage extends StatefulWidget {
  const ReportManagementPage({super.key});

  @override
  State<ReportManagementPage> createState() => _ReportManagementPageState();
}

class _ReportManagementPageState extends State<ReportManagementPage> {
  @override
  void initState() {
    super.initState();
    WidgetsBinding.instance.addPostFrameCallback((_) {
      context.read<ReportProvider>().loadReports();
    });
  }

  void _openCreateDialog() async {
    final req = await showDialog<CreateMenuReportRequest>(
      context: context,
      builder: (_) => const ReportFormDialog(),
    );

    if (req != null) {
      final ok = await context.read<ReportProvider>().createReport(req);
      if (!mounted) return;

      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(
          content: Text(ok ? '创建成功' : '创建失败'),
          backgroundColor: ok ? Colors.green : Colors.red,
        ),
      );
    }
  }

  void _openEditDialog(MenuReport report) async {
    final req = await showDialog<UpdateMenuReportRequest>(
      context: context,
      builder: (_) => ReportFormDialog(editing: report),
    );

    if (req != null) {
      final ok =
          await context.read<ReportProvider>().updateReport(report.id, req);
      if (!mounted) return;

      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(
          content: Text(ok ? '更新成功' : '更新失败'),
          backgroundColor: ok ? Colors.green : Colors.red,
        ),
      );
    }
  }

  void _deleteReport(MenuReport report) async {
    final confirmed = await showDialog<bool>(
      context: context,
      builder: (_) => AlertDialog(
        title: const Text('确认删除'),
        content: Text('确定要删除报菜记录 "${report.dishName}" 吗？'),
        actions: [
          TextButton(
            onPressed: () => Navigator.pop(context, false),
            child: const Text('取消'),
          ),
          TextButton(
            onPressed: () => Navigator.pop(context, true),
            style: TextButton.styleFrom(foregroundColor: Colors.red),
            child: const Text('删除'),
          ),
        ],
      ),
    );

    if (confirmed == true) {
      final ok = await context.read<ReportProvider>().deleteReport(report.id);
      if (!mounted) return;

      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(
          content: Text(ok ? '删除成功' : '删除失败'),
          backgroundColor: ok ? Colors.green : Colors.red,
        ),
      );
    }
  }

  Widget _buildCell(String text, double? width,
      {TextOverflow overflow = TextOverflow.clip}) {
    return Container(
      width: width,
      padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 10),
      child: Text(text, overflow: overflow, maxLines: 1),
    );
  }

  Widget _buildRow(MenuReport report, int index) {
    return Row(
      children: [
        _buildCell(report.id.toString(), 80),
        _buildCell(report.storeName ?? '-', 150),
        _buildCell(report.dishName, 200),
        _buildCell(report.quantity.toString(), 100),
        _buildCell(report.remark ?? '-', 200),
        _buildCell(report.createdAt ?? '-', 180),
        SizedBox(
          width: 150,
          child: Row(
            mainAxisSize: MainAxisSize.min,
            children: [
              TextButton(
                onPressed: () => _openEditDialog(report),
                child: const Text('编辑'),
              ),
              TextButton(
                onPressed: () => _deleteReport(report),
                style: TextButton.styleFrom(foregroundColor: Colors.red),
                child: const Text('删除'),
              ),
            ],
          ),
        ),
      ],
    );
  }

  @override
  Widget build(BuildContext context) {
    return PermissionGate(
      required: 'report:list',
      child: Scaffold(
        appBar: AppBar(
          title: const Text('报菜管理'),
          actions: [
            Padding(
              padding: const EdgeInsets.only(right: 16),
              child: ElevatedButton.icon(
                onPressed: _openCreateDialog,
                icon: const Icon(Icons.add),
                label: const Text('新增报菜'),
              ),
            ),
          ],
        ),
        body: Consumer<ReportProvider>(
          builder: (context, provider, child) {
            if (provider.loading && provider.reports.isEmpty) {
              return const Center(child: CircularProgressIndicator());
            }

            if (provider.error != null) {
              return Center(
                child: Column(
                  mainAxisAlignment: MainAxisAlignment.center,
                  children: [
                    Text(
                      '加载失败: ${provider.error}',
                      style: const TextStyle(color: Colors.red),
                    ),
                    const SizedBox(height: 16),
                    ElevatedButton(
                      onPressed: () => provider.loadReports(),
                      child: const Text('重试'),
                    ),
                  ],
                ),
              );
            }

            return AdminTable<MenuReport>(
              columns: const [
                AdminTableColumn(label: 'ID', width: 80),
                AdminTableColumn(label: '门店', width: 150),
                AdminTableColumn(label: '菜品', width: 200),
                AdminTableColumn(label: '数量', width: 100),
                AdminTableColumn(label: '备注', width: 200),
                AdminTableColumn(label: '创建时间', width: 180),
                AdminTableColumn(label: '操作', width: 150),
              ],
              data: provider.reports,
              rowBuilder: _buildRow,
            );
          },
        ),
      ),
    );
  }
}

class ReportFormDialog extends StatefulWidget {
  final MenuReport? editing;

  const ReportFormDialog({super.key, this.editing});

  @override
  State<ReportFormDialog> createState() => _ReportFormDialogState();
}

class _ReportFormDialogState extends State<ReportFormDialog> {
  final _formKey = GlobalKey<FormState>();
  late final TextEditingController _dishIdController;
  late final TextEditingController _quantityController;
  late final TextEditingController _remarkController;

  bool get isEditing => widget.editing != null;

  @override
  void initState() {
    super.initState();
    _dishIdController = TextEditingController(
      text: widget.editing?.dishId.toString() ?? '',
    );
    _quantityController = TextEditingController(
      text: widget.editing?.quantity.toString() ?? '',
    );
    _remarkController = TextEditingController(
      text: widget.editing?.remark ?? '',
    );
  }

  @override
  void dispose() {
    _dishIdController.dispose();
    _quantityController.dispose();
    _remarkController.dispose();
    super.dispose();
  }

  void _submit() {
    if (!_formKey.currentState!.validate()) return;

    if (isEditing) {
      final req = UpdateMenuReportRequest(
        quantity: int.tryParse(_quantityController.text),
        remark: _remarkController.text.trim().isEmpty
            ? null
            : _remarkController.text.trim(),
      );
      Navigator.pop(context, req);
    } else {
      final req = CreateMenuReportRequest(
        dishId: int.parse(_dishIdController.text),
        quantity: int.parse(_quantityController.text),
        remark: _remarkController.text.trim().isEmpty
            ? null
            : _remarkController.text.trim(),
      );
      Navigator.pop(context, req);
    }
  }

  @override
  Widget build(BuildContext context) {
    return AlertDialog(
      title: Text(isEditing ? '编辑报菜' : '新增报菜'),
      content: Form(
        key: _formKey,
        child: SizedBox(
          width: 400,
          child: Column(
            mainAxisSize: MainAxisSize.min,
            children: [
              if (!isEditing)
                TextFormField(
                  controller: _dishIdController,
                  decoration: const InputDecoration(
                    labelText: '菜品ID *',
                    hintText: '请输入菜品ID',
                  ),
                  keyboardType: TextInputType.number,
                  validator: (v) {
                    if (v == null || v.trim().isEmpty) return '请输入菜品ID';
                    if (int.tryParse(v) == null) return '请输入有效的数字';
                    return null;
                  },
                ),
              const SizedBox(height: 16),
              TextFormField(
                controller: _quantityController,
                decoration: const InputDecoration(
                  labelText: '数量 *',
                  hintText: '请输入数量',
                ),
                keyboardType: TextInputType.number,
                validator: (v) {
                  if (v == null || v.trim().isEmpty) return '请输入数量';
                  if (int.tryParse(v) == null) return '请输入有效的数字';
                  final num = int.parse(v);
                  if (num <= 0) return '数量必须大于0';
                  return null;
                },
              ),
              const SizedBox(height: 16),
              TextFormField(
                controller: _remarkController,
                decoration: const InputDecoration(
                  labelText: '备注',
                  hintText: '选填',
                ),
                maxLines: 3,
              ),
            ],
          ),
        ),
      ),
      actions: [
        TextButton(
          onPressed: () => Navigator.pop(context),
          child: const Text('取消'),
        ),
        ElevatedButton(
          onPressed: _submit,
          child: Text(isEditing ? '更新' : '创建'),
        ),
      ],
    );
  }
}
