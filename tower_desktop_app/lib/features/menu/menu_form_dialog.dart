import 'package:flutter/material.dart';
import 'package:tdesign_flutter/tdesign_flutter.dart';
import '../../core/constants/menu_types.dart';
import 'models.dart';

/// 菜单新增/编辑弹窗
class MenuFormDialog extends StatefulWidget {
  final MenuItem? editing;
  final List<MenuItem> parents; // 用于选择父级
  final int? initialParentId; // 新增子菜单时预选父级
  const MenuFormDialog(
      {super.key, this.editing, required this.parents, this.initialParentId});

  @override
  State<MenuFormDialog> createState() => _MenuFormDialogState();
}

class _MenuFormDialogState extends State<MenuFormDialog> {
  final _formKey = GlobalKey<FormState>();
  late int _parentId;
  late int _type;
  final _nameCtrl = TextEditingController();
  final _titleCtrl = TextEditingController();
  final _iconCtrl = TextEditingController();
  final _pathCtrl = TextEditingController();
  final _componentCtrl = TextEditingController();
  final _permissionCtrl = TextEditingController();
  final _sortCtrl = TextEditingController();
  bool _visible = true;
  bool _status = true;
  bool _isExternal = false;
  bool _cache = false;
  final _routeParamsCtrl = TextEditingController();

  @override
  void initState() {
    super.initState();
    if (widget.editing != null) {
      final m = widget.editing!;
      _parentId = m.parentId;
      _type = m.type ?? MenuType.page;
      _nameCtrl.text = m.name;
      _titleCtrl.text = m.title;
      _iconCtrl.text = m.icon ?? '';
      _pathCtrl.text = m.path ?? '';
      _componentCtrl.text = m.component ?? '';
      _permissionCtrl.text = m.permission ?? '';
      _sortCtrl.text = (m.sort ?? '').toString();
      _visible = (m.visible ?? 1) == 1;
      _status = (m.status ?? 1) == 1;
      _isExternal = (m.isExternal ?? 0) == 1;
      _cache = (m.cache ?? 0) == 1;
      _routeParamsCtrl.text = m.routeParams ?? '';
    } else {
      _parentId = widget.initialParentId ?? 0;
      _type = MenuType.page;
    }
  }

  @override
  void dispose() {
    _nameCtrl.dispose();
    _titleCtrl.dispose();
    _iconCtrl.dispose();
    _pathCtrl.dispose();
    _componentCtrl.dispose();
    _permissionCtrl.dispose();
    _sortCtrl.dispose();
    _routeParamsCtrl.dispose();
    super.dispose();
  }

  void _submit() {
    if (!_formKey.currentState!.validate()) return;
    if (widget.editing == null) {
      final req = CreateMenuRequest(
        parentId: _parentId,
        name: _nameCtrl.text.trim(),
        title: _titleCtrl.text.trim(),
        type: _type,
        sort: _sortCtrl.text.trim().isEmpty
            ? null
            : int.tryParse(_sortCtrl.text.trim()),
        icon: _iconCtrl.text.trim().isEmpty ? null : _iconCtrl.text.trim(),
        path: _pathCtrl.text.trim().isEmpty ? null : _pathCtrl.text.trim(),
        component: _componentCtrl.text.trim().isEmpty
            ? null
            : _componentCtrl.text.trim(),
        permission: _permissionCtrl.text.trim().isEmpty
            ? null
            : _permissionCtrl.text.trim(),
        visible: _visible ? 1 : 0,
        status: _status ? 1 : 0,
        isExternal: _isExternal ? 1 : 0,
        cache: _cache ? 1 : 0,
        routeParams: _routeParamsCtrl.text.trim().isEmpty
            ? null
            : _routeParamsCtrl.text.trim(),
      );
      Navigator.of(context).pop(req);
    } else {
      final req = UpdateMenuRequest(
        parentId: _parentId == widget.editing!.parentId ? null : _parentId,
        name: _nameCtrl.text.trim() == widget.editing!.name
            ? null
            : _nameCtrl.text.trim(),
        title: _titleCtrl.text.trim() == widget.editing!.title
            ? null
            : _titleCtrl.text.trim(),
        type: _type == (widget.editing!.type ?? MenuType.page) ? null : _type,
        sort: _sortCtrl.text.trim().isEmpty
            ? null
            : (int.tryParse(_sortCtrl.text.trim()) == widget.editing!.sort
                ? null
                : int.tryParse(_sortCtrl.text.trim())),
        icon: _iconCtrl.text.trim() == (widget.editing!.icon ?? '')
            ? null
            : _iconCtrl.text.trim(),
        path: _pathCtrl.text.trim() == (widget.editing!.path ?? '')
            ? null
            : _pathCtrl.text.trim(),
        component:
            _componentCtrl.text.trim() == (widget.editing!.component ?? '')
                ? null
                : _componentCtrl.text.trim(),
        permission:
            _permissionCtrl.text.trim() == (widget.editing!.permission ?? '')
                ? null
                : _permissionCtrl.text.trim(),
        visible: (_visible ? 1 : 0) == (widget.editing!.visible ?? 1)
            ? null
            : (_visible ? 1 : 0),
        status: (_status ? 1 : 0) == (widget.editing!.status ?? 1)
            ? null
            : (_status ? 1 : 0),
        isExternal: (_isExternal ? 1 : 0) == (widget.editing!.isExternal ?? 0)
            ? null
            : (_isExternal ? 1 : 0),
        cache: (_cache ? 1 : 0) == (widget.editing!.cache ?? 0)
            ? null
            : (_cache ? 1 : 0),
        routeParams:
            _routeParamsCtrl.text.trim() == (widget.editing!.routeParams ?? '')
                ? null
                : _routeParamsCtrl.text.trim(),
      );
      Navigator.of(context).pop(req);
    }
  }

  @override
  Widget build(BuildContext context) {
    return Dialog(
      insetPadding: const EdgeInsets.all(40),
      child: ConstrainedBox(
        constraints: const BoxConstraints(maxWidth: 600),
        child: Padding(
          padding: const EdgeInsets.all(24),
          child: SingleChildScrollView(
            child: Form(
              key: _formKey,
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Text(
                    widget.editing == null ? '新增菜单' : '编辑菜单',
                    style: const TextStyle(
                      fontSize: 20,
                      fontWeight: FontWeight.w600,
                    ),
                  ),
                  const SizedBox(height: 20),
                  _buildParentField(),
                  _buildTypeField(),
                  _buildTextField(_nameCtrl,
                      label: '路由名称(name)', validator: _notEmpty),
                  _buildTextField(_titleCtrl,
                      label: '标题(title)', validator: _notEmpty),
                  _buildTextField(_iconCtrl, label: '图标(icon)', hint: '可选'),
                  if (_type != MenuType.directory)
                    _buildTextField(_pathCtrl,
                        label: '路由地址(path)', hint: '页面/按钮可选'),
                  if (_type == MenuType.page)
                    _buildTextField(_componentCtrl,
                        label: '组件(component)', hint: '前端组件路径'),
                  if (_type != MenuType.directory)
                    _buildTextField(_permissionCtrl,
                        label: '权限码(permission)', hint: '如 system:menu:add'),
                  if (_type == MenuType.page)
                    _buildTextField(_routeParamsCtrl,
                        label: '路由参数(route_params)', hint: '如 id=1&mode=edit'),
                  _buildTextField(
                    _sortCtrl,
                    label: '排序(sort)',
                    hint: '数字升序',
                    keyboardType: TextInputType.number,
                  ),
                  const SizedBox(height: 8),
                  Row(
                    children: [
                      _buildSwitch(
                          '显示', _visible, (v) => setState(() => _visible = v)),
                      const SizedBox(width: 24),
                      _buildSwitch(
                          '启用', _status, (v) => setState(() => _status = v)),
                      const SizedBox(width: 24),
                      if (_type == MenuType.page)
                        _buildSwitch('外链', _isExternal,
                            (v) => setState(() => _isExternal = v)),
                      const SizedBox(width: 24),
                      if (_type == MenuType.page)
                        _buildSwitch(
                            '缓存', _cache, (v) => setState(() => _cache = v)),
                    ],
                  ),
                  const SizedBox(height: 24),
                  Row(
                    children: [
                      Expanded(
                        child: TDButton(
                          text: '取消',
                          type: TDButtonType.outline,
                          size: TDButtonSize.large,
                          onTap: () => Navigator.of(context).pop(),
                        ),
                      ),
                      const SizedBox(width: 16),
                      Expanded(
                        child: TDButton(
                          text: '确定',
                          theme: TDButtonTheme.primary,
                          size: TDButtonSize.large,
                          onTap: _submit,
                        ),
                      ),
                    ],
                  ),
                ],
              ),
            ),
          ),
        ),
      ),
    );
  }

  Widget _buildParentField() {
    // 去重并确保当前 parentId 存在; 若不存在回退到 0
    final uniqueParents = <int, MenuItem>{
      for (final p in widget.parents) p.id: p,
    }.values.toList();
    if (_parentId != 0 && uniqueParents.every((e) => e.id != _parentId)) {
      _parentId = 0; // 回退根节点避免 Dropdown 断言
    }
    return Padding(
      padding: const EdgeInsets.only(bottom: 16),
      child: DropdownButtonFormField<int>(
        value: _parentId == 0 ? 0 : _parentId,
        items: [
          const DropdownMenuItem(value: 0, child: Text('根节点')),
          ...uniqueParents.map((m) => DropdownMenuItem<int>(
                value: m.id,
                child: Text('[${m.id}] ${m.title}'),
              ))
        ],
        onChanged: (v) => setState(() => _parentId = v ?? 0),
        decoration: const InputDecoration(labelText: '父级菜单'),
      ),
    );
  }

  Widget _buildTypeField() {
    return Padding(
      padding: const EdgeInsets.only(bottom: 16),
      child: DropdownButtonFormField<int>(
        value: _type,
        items: const [
          DropdownMenuItem(value: MenuType.directory, child: Text('目录(1)')),
          DropdownMenuItem(value: MenuType.page, child: Text('页面(2)')),
          DropdownMenuItem(value: MenuType.button, child: Text('按钮(3)')),
        ],
        onChanged: (v) => setState(() => _type = v ?? MenuType.page),
        decoration: const InputDecoration(labelText: '类型'),
      ),
    );
  }

  Widget _buildTextField(TextEditingController ctrl,
      {required String label,
      String? hint,
      String? Function(String?)? validator,
      TextInputType? keyboardType}) {
    return Padding(
      padding: const EdgeInsets.only(bottom: 16),
      child: TextFormField(
        controller: ctrl,
        decoration: InputDecoration(labelText: label, hintText: hint),
        validator: validator,
        keyboardType: keyboardType,
      ),
    );
  }

  String? _notEmpty(String? v) {
    if (v == null || v.trim().isEmpty) return '必填';
    return null;
  }

  Widget _buildSwitch(String label, bool value, ValueChanged<bool> onChanged) {
    return Row(
      mainAxisSize: MainAxisSize.min,
      children: [
        Switch(value: value, onChanged: onChanged),
        Text(label),
      ],
    );
  }
}
