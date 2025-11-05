import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import 'role_models.dart';
import 'role_provider.dart';

class RoleFormDialog extends StatefulWidget {
  final RoleItem? editing;
  const RoleFormDialog({super.key, this.editing});
  @override
  State<RoleFormDialog> createState() => _RoleFormDialogState();
}

class _RoleFormDialogState extends State<RoleFormDialog> {
  final _formKey = GlobalKey<FormState>();
  final _nameCtrl = TextEditingController();
  final _codeCtrl = TextEditingController();
  final _remarkCtrl = TextEditingController();
  bool _status = true;

  @override
  void initState() {
    super.initState();
    if (widget.editing != null) {
      final r = widget.editing!;
      _nameCtrl.text = r.name;
      _codeCtrl.text = r.code;
      _remarkCtrl.text = r.description ?? '';
      _status = (r.status ?? 1) == 1;
    }
  }

  @override
  void dispose() {
    _nameCtrl.dispose();
    _codeCtrl.dispose();
    _remarkCtrl.dispose();
    super.dispose();
  }

  void _submit() {
    if (!_formKey.currentState!.validate()) return;
    if (widget.editing == null) {
      final req = CreateRoleRequest(
        name: _nameCtrl.text.trim(),
        code: _codeCtrl.text.trim(),
        status: _status ? 1 : 0,
        description:
            _remarkCtrl.text.trim().isEmpty ? null : _remarkCtrl.text.trim(),
      );
      Navigator.pop(context, req);
    } else {
      final e = widget.editing!;
      final req = UpdateRoleRequest(
        name: _nameCtrl.text.trim() == e.name ? null : _nameCtrl.text.trim(),
        code: _codeCtrl.text.trim() == e.code ? null : _codeCtrl.text.trim(),
        status: (_status ? 1 : 0) == (e.status ?? 1) ? null : (_status ? 1 : 0),
        description: _remarkCtrl.text.trim() == (e.description ?? '')
            ? null
            : _remarkCtrl.text.trim(),
      );
      Navigator.pop(context, req);
    }
  }

  @override
  Widget build(BuildContext context) {
    return Dialog(
      child: ConstrainedBox(
        constraints: const BoxConstraints(maxWidth: 480),
        child: Padding(
          padding: const EdgeInsets.all(24),
          child: Form(
            key: _formKey,
            child: Column(
              mainAxisSize: MainAxisSize.min,
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text(widget.editing == null ? '新增角色' : '编辑角色',
                    style: const TextStyle(
                        fontSize: 20, fontWeight: FontWeight.w600)),
                const SizedBox(height: 16),
                TextFormField(
                  controller: _nameCtrl,
                  decoration: const InputDecoration(labelText: '角色名称'),
                  maxLength: 32,
                  validator: (v) {
                    final value = v?.trim() ?? '';
                    if (value.isEmpty) return '必填';
                    if (value.length < 2) return '至少 2 个字符';
                    return null;
                  },
                ),
                const SizedBox(height: 12),
                TextFormField(
                  controller: _codeCtrl,
                  decoration: const InputDecoration(labelText: '角色编码'),
                  maxLength: 40,
                  validator: (v) {
                    final value = v?.trim() ?? '';
                    if (value.isEmpty) return '必填';
                    final reg = RegExp(r'^[a-zA-Z0-9_:-]+$');
                    if (!reg.hasMatch(value)) return '仅允许字母、数字、_、:、-';
                    // 唯一性校验
                    final provider = context.read<RoleProvider>();
                    final editingId = widget.editing?.id;
                    final exists = provider.list
                        .any((r) => r.code == value && r.id != editingId);
                    if (exists) return '编码已存在';
                    return null;
                  },
                ),
                const SizedBox(height: 12),
                TextFormField(
                  controller: _remarkCtrl,
                  decoration: const InputDecoration(labelText: '备注(可选)'),
                  maxLines: 3,
                  maxLength: 200,
                  validator: (v) {
                    final value = v?.trim() ?? '';
                    if (value.isEmpty) return null;
                    if (value.length < 2) return '备注至少 2 个字符或留空';
                    return null;
                  },
                ),
                const SizedBox(height: 12),
                Row(
                  children: [
                    Switch(
                        value: _status,
                        onChanged: (v) => setState(() => _status = v)),
                    const Text('启用'),
                  ],
                ),
                const SizedBox(height: 20),
                Row(
                  children: [
                    Expanded(
                      child: OutlinedButton(
                        onPressed: () => Navigator.pop(context),
                        child: const Text('取消'),
                      ),
                    ),
                    const SizedBox(width: 16),
                    Expanded(
                      child: ElevatedButton(
                        onPressed: _submit,
                        child: const Text('确定'),
                      ),
                    ),
                  ],
                )
              ],
            ),
          ),
        ),
      ),
    );
  }
}
