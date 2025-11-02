import 'package:flutter/material.dart';
import 'package:tdesign_flutter/tdesign_flutter.dart';
import 'models.dart';

class UserFormDialog extends StatefulWidget {
  const UserFormDialog({super.key});

  @override
  State<UserFormDialog> createState() => _UserFormDialogState();
}

class _UserFormDialogState extends State<UserFormDialog> {
  final _formKey = GlobalKey<FormState>();
  final _usernameController = TextEditingController();
  final _phoneController = TextEditingController();
  final _emailController = TextEditingController();
  final _nicknameController = TextEditingController();
  final _passwordController = TextEditingController();
  final _remarkController = TextEditingController();

  int? _gender; // 1=男 2=女

  @override
  void dispose() {
    _usernameController.dispose();
    _phoneController.dispose();
    _emailController.dispose();
    _nicknameController.dispose();
    _passwordController.dispose();
    _remarkController.dispose();
    super.dispose();
  }

  void _handleSubmit() {
    // 验证表单
    if (!_formKey.currentState!.validate()) {
      TDToast.showText('请检查表单输入', context: context);
      return;
    }

    final req = CreateUserRequest(
      username: _usernameController.text.trim(),
      phone: _phoneController.text.trim(),
      password: _passwordController.text.trim(),
      email: _emailController.text.trim().isEmpty
          ? null
          : _emailController.text.trim(),
      nickname: _nicknameController.text.trim().isEmpty
          ? null
          : _nicknameController.text.trim(),
      gender: _gender,
      status: 1, // 默认启用
      remark: _remarkController.text.trim().isEmpty
          ? null
          : _remarkController.text.trim(),
    );
    Navigator.pop(context, req);
  }

  Widget _buildFormItem({
    required String label,
    required Widget child,
    bool required = false,
  }) {
    return Padding(
      padding: const EdgeInsets.only(bottom: 16),
      child: Row(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          SizedBox(
            width: 80,
            child: Padding(
              padding: const EdgeInsets.only(top: 12),
              child: Row(
                children: [
                  if (required)
                    const Text(
                      '*',
                      style: TextStyle(color: Colors.red),
                    ),
                  Text(
                    label,
                    style: const TextStyle(fontSize: 14),
                  ),
                ],
              ),
            ),
          ),
          Expanded(child: child),
        ],
      ),
    );
  }

  @override
  Widget build(BuildContext context) {
    return AlertDialog(
      title: const Text('新增用户'),
      content: SizedBox(
        width: 560,
        child: Form(
          key: _formKey,
          child: SingleChildScrollView(
            child: Column(
              mainAxisSize: MainAxisSize.min,
              children: [
                // 用户名
                _buildFormItem(
                  label: '用户名',
                  required: true,
                  child: TDInput(
                    controller: _usernameController,
                    hintText: '请输入用户名',
                    additionInfo: '用户名创建后不可修改',
                    needClear: true,
                    onChanged: (value) {
                      _formKey.currentState?.validate();
                    },
                  ),
                ),

                // 手机号
                _buildFormItem(
                  label: '手机号',
                  required: true,
                  child: TDInput(
                    controller: _phoneController,
                    hintText: '请输入11位手机号',
                    additionInfo: '用于登录账号，创建后不可修改',
                    needClear: true,
                    onChanged: (value) {
                      _formKey.currentState?.validate();
                    },
                  ),
                ),

                // 邮箱
                _buildFormItem(
                  label: '邮箱',
                  child: TDInput(
                    controller: _emailController,
                    hintText: '请输入邮箱（可选）',
                    needClear: true,
                  ),
                ),

                // 昵称
                _buildFormItem(
                  label: '昵称',
                  child: TDInput(
                    controller: _nicknameController,
                    hintText: '请输入昵称（可选）',
                    needClear: true,
                  ),
                ),

                // 密码
                _buildFormItem(
                  label: '密码',
                  required: true,
                  child: TDInput(
                    controller: _passwordController,
                    hintText: '请输入密码',
                    additionInfo: '密码至少6位',
                    obscureText: true,
                    needClear: true,
                    onChanged: (value) {
                      _formKey.currentState?.validate();
                    },
                  ),
                ),

                // 性别
                _buildFormItem(
                  label: '性别',
                  child: Row(
                    children: [
                      TDButton(
                        text: '男',
                        size: TDButtonSize.medium,
                        type: _gender == 1
                            ? TDButtonType.fill
                            : TDButtonType.outline,
                        theme: _gender == 1
                            ? TDButtonTheme.primary
                            : TDButtonTheme.defaultTheme,
                        onTap: () {
                          setState(() {
                            _gender = 1;
                          });
                        },
                      ),
                      const SizedBox(width: 12),
                      TDButton(
                        text: '女',
                        size: TDButtonSize.medium,
                        type: _gender == 2
                            ? TDButtonType.fill
                            : TDButtonType.outline,
                        theme: _gender == 2
                            ? TDButtonTheme.primary
                            : TDButtonTheme.defaultTheme,
                        onTap: () {
                          setState(() {
                            _gender = 2;
                          });
                        },
                      ),
                    ],
                  ),
                ),

                // 备注
                _buildFormItem(
                  label: '备注',
                  child: TDInput(
                    controller: _remarkController,
                    hintText: '请输入备注（可选）',
                    maxLines: 3,
                    needClear: true,
                  ),
                ),
              ],
            ),
          ),
        ),
      ),
      actions: [
        TDButton(
          text: '取消',
          theme: TDButtonTheme.defaultTheme,
          onTap: () => Navigator.pop(context),
        ),
        const SizedBox(width: 12),
        TDButton(
          text: '确定',
          theme: TDButtonTheme.primary,
          onTap: _handleSubmit,
        ),
      ],
    );
  }
}
