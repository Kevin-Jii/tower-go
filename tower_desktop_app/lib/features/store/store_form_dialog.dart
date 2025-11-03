import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import 'package:tdesign_flutter/tdesign_flutter.dart';
import 'models.dart';
import '../user/user_provider.dart';
import '../user/models.dart';

class StoreFormDialog extends StatefulWidget {
  final Store? store; // null 表示新增，非 null 表示编辑

  const StoreFormDialog({super.key, this.store});

  @override
  State<StoreFormDialog> createState() => _StoreFormDialogState();
}

class _StoreFormDialogState extends State<StoreFormDialog> {
  final _formKey = GlobalKey<FormState>();
  late final TextEditingController _nameController;
  late final TextEditingController _addressController;
  late final TextEditingController _remarkController;
  int _status = 1; // 1=营业中 0=停业中
  User? _selectedUser; // 选中的负责人
  List<User> _users = [];
  bool _loadingUsers = false;

  // 营业时间
  TimeOfDay? _startTime;
  TimeOfDay? _endTime;

  @override
  void initState() {
    super.initState();
    _nameController = TextEditingController(text: widget.store?.name ?? '');
    _addressController =
        TextEditingController(text: widget.store?.address ?? '');
    _remarkController = TextEditingController(text: widget.store?.remark ?? '');
    _status = widget.store?.status ?? 1;

    // 解析已有的营业时间
    if (widget.store?.businessHours != null) {
      final parts = widget.store!.businessHours!.split('-');
      if (parts.length == 2) {
        _startTime = _parseTime(parts[0].trim());
        _endTime = _parseTime(parts[1].trim());
      }
    }

    _loadUsers();
  }

  TimeOfDay? _parseTime(String timeStr) {
    try {
      final parts = timeStr.split(':');
      if (parts.length == 2) {
        return TimeOfDay(
          hour: int.parse(parts[0]),
          minute: int.parse(parts[1]),
        );
      }
    } catch (e) {
      // 解析失败
    }
    return null;
  }

  String _formatTime(TimeOfDay time) {
    return '${time.hour.toString().padLeft(2, '0')}:${time.minute.toString().padLeft(2, '0')}';
  }

  Future<void> _loadUsers() async {
    setState(() => _loadingUsers = true);
    try {
      final provider = context.read<UserProvider>();
      await provider.loadUsers(page: 1, pageSize: 100); // 加载所有用户
      setState(() {
        _users = provider.users;
        _loadingUsers = false;
      });
    } catch (e) {
      setState(() => _loadingUsers = false);
    }
  }

  @override
  void dispose() {
    _nameController.dispose();
    _addressController.dispose();
    _remarkController.dispose();
    super.dispose();
  }

  bool get _isEdit => widget.store != null;

  void _handleSubmit() {
    if (!_formKey.currentState!.validate()) return;

    // 生成营业时间字符串
    String? businessHours;
    if (_startTime != null && _endTime != null) {
      businessHours = '${_formatTime(_startTime!)}-${_formatTime(_endTime!)}';
    }

    // 获取用户名（优先使用nickname，否则使用username）
    String? contactPersonName;
    if (_selectedUser != null) {
      contactPersonName = _selectedUser!.nickname ?? _selectedUser!.username;
    }

    if (_isEdit) {
      // 编辑模式
      final req = UpdateStoreRequest(
        name: _nameController.text.trim(),
        address: _addressController.text.trim().isEmpty
            ? null
            : _addressController.text.trim(),
        status: _status,
        remark: _remarkController.text.trim().isEmpty
            ? null
            : _remarkController.text.trim(),
        businessHours: businessHours,
        contactPerson: contactPersonName,
        phone: _selectedUser?.phone,
      );
      Navigator.of(context).pop(req);
    } else {
      // 新增模式
      final req = CreateStoreRequest(
        name: _nameController.text.trim(),
        address: _addressController.text.trim().isEmpty
            ? null
            : _addressController.text.trim(),
        status: _status,
        remark: _remarkController.text.trim().isEmpty
            ? null
            : _remarkController.text.trim(),
        businessHours: businessHours,
        contactPerson: contactPersonName,
        phone: _selectedUser?.phone,
      );
      Navigator.of(context).pop(req);
    }
  }

  @override
  Widget build(BuildContext context) {
    return Dialog(
      shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(12)),
      child: Container(
        width: 500,
        constraints: const BoxConstraints(maxHeight: 650),
        padding: const EdgeInsets.all(24),
        child: Column(
          mainAxisSize: MainAxisSize.min,
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Row(
              mainAxisAlignment: MainAxisAlignment.spaceBetween,
              children: [
                Text(
                  _isEdit ? '编辑门店' : '新增门店',
                  style: const TextStyle(
                      fontSize: 20, fontWeight: FontWeight.w600),
                ),
                IconButton(
                  icon: const Icon(Icons.close),
                  onPressed: () => Navigator.of(context).pop(),
                ),
              ],
            ),
            const SizedBox(height: 24),
            Flexible(
              child: SingleChildScrollView(
                child: Form(
                  key: _formKey,
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      _buildTextField(
                        label: '门店名称',
                        controller: _nameController,
                        required: true,
                        validator: (v) {
                          if (v == null || v.trim().isEmpty) {
                            return '请输入门店名称';
                          }
                          return null;
                        },
                      ),
                      const SizedBox(height: 20),
                      _buildTextField(
                        label: '门店地址',
                        controller: _addressController,
                        required: false,
                        maxLines: 2,
                      ),
                      const SizedBox(height: 20),
                      _buildBusinessHoursSelector(),
                      const SizedBox(height: 20),
                      _buildUserSelector(),
                      const SizedBox(height: 20),
                      const Text(
                        '状态',
                        style: TextStyle(
                          fontSize: 14,
                          fontWeight: FontWeight.w500,
                          color: Colors.black87,
                        ),
                      ),
                      const SizedBox(height: 8),
                      Row(
                        children: [
                          Expanded(
                            child: TDButton(
                              text: '营业中',
                              type: _status == 1
                                  ? TDButtonType.fill
                                  : TDButtonType.outline,
                              theme: _status == 1
                                  ? TDButtonTheme.primary
                                  : TDButtonTheme.defaultTheme,
                              onTap: () => setState(() => _status = 1),
                            ),
                          ),
                          const SizedBox(width: 12),
                          Expanded(
                            child: TDButton(
                              text: '停业中',
                              type: _status == 0
                                  ? TDButtonType.fill
                                  : TDButtonType.outline,
                              theme: _status == 0
                                  ? TDButtonTheme.danger
                                  : TDButtonTheme.defaultTheme,
                              onTap: () => setState(() => _status = 0),
                            ),
                          ),
                        ],
                      ),
                      const SizedBox(height: 20),
                      _buildTextField(
                        label: '备注',
                        controller: _remarkController,
                        required: false,
                        maxLines: 3,
                      ),
                    ],
                  ),
                ),
              ),
            ),
            const SizedBox(height: 24),
            Row(
              mainAxisAlignment: MainAxisAlignment.end,
              children: [
                TDButton(
                  text: '取消',
                  type: TDButtonType.outline,
                  theme: TDButtonTheme.defaultTheme,
                  onTap: () => Navigator.of(context).pop(),
                ),
                const SizedBox(width: 12),
                TDButton(
                  text: _isEdit ? '保存' : '创建',
                  theme: TDButtonTheme.primary,
                  onTap: _handleSubmit,
                ),
              ],
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildTextField({
    required String label,
    required TextEditingController controller,
    required bool required,
    int maxLines = 1,
    String? hintText,
    String? Function(String?)? validator,
  }) {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Row(
          children: [
            if (required)
              const Text(
                '* ',
                style: TextStyle(color: Colors.red, fontSize: 14),
              ),
            Text(
              label,
              style: const TextStyle(
                fontSize: 14,
                fontWeight: FontWeight.w500,
                color: Colors.black87,
              ),
            ),
          ],
        ),
        const SizedBox(height: 8),
        TextFormField(
          controller: controller,
          maxLines: maxLines,
          decoration: InputDecoration(
            hintText: hintText ?? '请输入$label',
            border: OutlineInputBorder(
              borderRadius: BorderRadius.circular(6),
            ),
            contentPadding:
                const EdgeInsets.symmetric(horizontal: 12, vertical: 10),
          ),
          validator: validator,
        ),
      ],
    );
  }

  Widget _buildBusinessHoursSelector() {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        const Text(
          '营业时间',
          style: TextStyle(
            fontSize: 14,
            fontWeight: FontWeight.w500,
            color: Colors.black87,
          ),
        ),
        const SizedBox(height: 8),
        Row(
          children: [
            Expanded(
              child: InkWell(
                onTap: () async {
                  final time = await showTimePicker(
                    context: context,
                    initialTime:
                        _startTime ?? const TimeOfDay(hour: 9, minute: 0),
                    builder: (context, child) {
                      return MediaQuery(
                        data: MediaQuery.of(context).copyWith(
                          alwaysUse24HourFormat: true,
                        ),
                        child: child!,
                      );
                    },
                  );
                  if (time != null) {
                    setState(() => _startTime = time);
                  }
                },
                child: Container(
                  padding:
                      const EdgeInsets.symmetric(horizontal: 12, vertical: 10),
                  decoration: BoxDecoration(
                    border: Border.all(color: Colors.grey.shade300),
                    borderRadius: BorderRadius.circular(6),
                  ),
                  child: Row(
                    mainAxisAlignment: MainAxisAlignment.spaceBetween,
                    children: [
                      Text(
                        _startTime != null ? _formatTime(_startTime!) : '开始时间',
                        style: TextStyle(
                          fontSize: 14,
                          color: _startTime != null
                              ? Colors.black87
                              : Colors.grey.shade500,
                        ),
                      ),
                      Icon(Icons.access_time,
                          size: 20, color: Colors.grey.shade600),
                    ],
                  ),
                ),
              ),
            ),
            const Padding(
              padding: EdgeInsets.symmetric(horizontal: 12),
              child: Text('至', style: TextStyle(fontSize: 14)),
            ),
            Expanded(
              child: InkWell(
                onTap: () async {
                  final time = await showTimePicker(
                    context: context,
                    initialTime:
                        _endTime ?? const TimeOfDay(hour: 21, minute: 0),
                    builder: (context, child) {
                      return MediaQuery(
                        data: MediaQuery.of(context).copyWith(
                          alwaysUse24HourFormat: true,
                        ),
                        child: child!,
                      );
                    },
                  );
                  if (time != null) {
                    setState(() => _endTime = time);
                  }
                },
                child: Container(
                  padding:
                      const EdgeInsets.symmetric(horizontal: 12, vertical: 10),
                  decoration: BoxDecoration(
                    border: Border.all(color: Colors.grey.shade300),
                    borderRadius: BorderRadius.circular(6),
                  ),
                  child: Row(
                    mainAxisAlignment: MainAxisAlignment.spaceBetween,
                    children: [
                      Text(
                        _endTime != null ? _formatTime(_endTime!) : '结束时间',
                        style: TextStyle(
                          fontSize: 14,
                          color: _endTime != null
                              ? Colors.black87
                              : Colors.grey.shade500,
                        ),
                      ),
                      Icon(Icons.access_time,
                          size: 20, color: Colors.grey.shade600),
                    ],
                  ),
                ),
              ),
            ),
          ],
        ),
      ],
    );
  }

  Widget _buildUserSelector() {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        const Text(
          '负责人',
          style: TextStyle(
            fontSize: 14,
            fontWeight: FontWeight.w500,
            color: Colors.black87,
          ),
        ),
        const SizedBox(height: 8),
        _loadingUsers
            ? const Center(
                child: Padding(
                  padding: EdgeInsets.all(16.0),
                  child: CircularProgressIndicator(),
                ),
              )
            : Container(
                decoration: BoxDecoration(
                  border: Border.all(color: Colors.grey.shade300),
                  borderRadius: BorderRadius.circular(6),
                ),
                child: DropdownButtonHideUnderline(
                  child: DropdownButton<User>(
                    value: _selectedUser,
                    hint: const Padding(
                      padding: EdgeInsets.symmetric(horizontal: 12),
                      child: Text('请选择负责人'),
                    ),
                    isExpanded: true,
                    items: _users.map((user) {
                      return DropdownMenuItem<User>(
                        value: user,
                        child: Padding(
                          padding: const EdgeInsets.symmetric(horizontal: 12),
                          child: Text(
                            '${user.nickname ?? user.username} (${user.phone})',
                          ),
                        ),
                      );
                    }).toList(),
                    onChanged: (User? user) {
                      setState(() {
                        _selectedUser = user;
                      });
                    },
                  ),
                ),
              ),
        if (_selectedUser != null) ...[
          const SizedBox(height: 8),
          Text(
            '电话: ${_selectedUser!.phone}',
            style: TextStyle(
              fontSize: 13,
              color: Colors.grey.shade600,
            ),
          ),
        ],
      ],
    );
  }
}
