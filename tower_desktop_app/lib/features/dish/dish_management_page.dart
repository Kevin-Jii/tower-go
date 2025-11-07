import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import '../../core/theme/design_tokens.dart';
import '../../core/constants/ui_texts.dart';
import '../../core/constants/error_texts.dart';
import '../../core/constants/app_constants.dart';
import '../../core/utils/value_parsers.dart';
import '../../core/widgets/status_tag.dart';
import '../../core/widgets/admin_table.dart';
import '../store/store_api.dart';
import '../auth/permission_gate.dart';
import '../auth/permission_provider.dart';
import 'dish_api.dart';
import 'models.dart';
import '../../core/network/api_client.dart';
import 'providers/store_selector_provider.dart';
import 'providers/dish_category_provider.dart';
import 'providers/dish_list_provider.dart';
import 'dish_category_api.dart';

class DishManagementPage extends StatefulWidget {
  const DishManagementPage({super.key});
  @override
  State<DishManagementPage> createState() => _DishManagementPageState();
}

class _DishManagementPageState extends State<DishManagementPage> {
  @override
  void initState() {
    super.initState();
    WidgetsBinding.instance.addPostFrameCallback((_) {
      context.read<StoreSelectorProvider>().initializeStore();
    });
  }

  void _openCreateDish() async {
    final categoryProvider = context.read<DishCategoryProvider>();
    if (categoryProvider.selectedCategory == null) return;
    final req = await showDialog<CreateDishRequest>(
      context: context,
      builder: (_) => const DishFormDialog(),
    );
    if (req != null) {
      final ok = await context.read<DishListProvider>().createDish(req);
      if (!mounted) return;
      ScaffoldMessenger.of(context).showSnackBar(SnackBar(
        content:
            Text(ok ? UITexts.dishCreateSuccess : UITexts.dishCreateFailed),
      ));
    }
  }

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    return Scaffold(
      backgroundColor: theme.scaffoldBackgroundColor,
      body: Row(
        children: [
          _buildCategorySidebar(),
          Expanded(
            child: Container(
              decoration: BoxDecoration(
                border: Border(
                    left: BorderSide(color: theme.dividerColor, width: 1)),
              ),
              child: _buildDishArea(),
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildCategorySidebar() {
    return Consumer2<DishCategoryProvider, DishListProvider>(
        builder: (context, categoryProvider, dishList, __) {
      final theme = Theme.of(context);

      // Show error as SnackBar instead of inline display
      if (categoryProvider.error != null) {
        WidgetsBinding.instance.addPostFrameCallback((_) {
          if (context.mounted) {
            ScaffoldMessenger.of(context).showSnackBar(
              SnackBar(
                content: Text(
                    '${UITexts.dishCategoryError}: ${categoryProvider.error}'),
                backgroundColor: DesignTokens.danger,
                duration: const Duration(seconds: 1),
              ),
            );
          }
        });
      }

      return Container(
        width: DesignTokens.sidebarWidth,
        decoration: BoxDecoration(
          color: theme.cardColor,
          boxShadow: DesignTokens.shadowSm,
        ),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Container(
              padding: const EdgeInsets.fromLTRB(16, 18, 16, 12),
              child: Row(
                children: [
                  Text(
                    UITexts.dishCategoryTitle,
                    style: DesignTokens.heading3
                        .copyWith(color: theme.textTheme.bodyLarge?.color),
                  ),
                  const Spacer(),
                  PermissionGate(
                    required: PermissionCodes.dishAdd,
                    child: IconButton(
                      icon: const Icon(Icons.add_circle_outline,
                          color: DesignTokens.brandPrimary),
                      tooltip: UITexts.dishCategoryAdd,
                      onPressed: () async {
                        final catProvider =
                            context.read<DishCategoryProvider>();
                        final name = await showDialog<String>(
                          context: context,
                          builder: (_) => ChangeNotifierProvider.value(
                            value: catProvider,
                            child: const DishCategoryFormDialog(),
                          ),
                        );
                        if (name != null && name.trim().isNotEmpty) {
                          final ok =
                              await catProvider.createCategory(name.trim());
                          if (!mounted) return;
                          if (!ok) {
                            ScaffoldMessenger.of(context).showSnackBar(
                                const SnackBar(
                                    content: Text(ErrorTexts.createCategory)));
                          }
                        }
                      },
                    ),
                  ),
                ],
              ),
            ),
            Divider(height: 1, color: theme.dividerColor),
            Expanded(
              child: categoryProvider.loading
                  ? Center(
                      child: Text(UITexts.dishCategoryLoading,
                          style: DesignTokens.body.copyWith(
                              color: theme.textTheme.bodySmall?.color)))
                  : categoryProvider.categories.isEmpty
                      ? Center(
                          child: Text(UITexts.dishCategoryEmpty,
                              style: DesignTokens.body.copyWith(
                                  color: theme.textTheme.bodySmall?.color)))
                      : ListView.separated(
                          itemCount: categoryProvider.categories.length,
                          separatorBuilder: (_, __) =>
                              Divider(height: 1, color: theme.dividerColor),
                          itemBuilder: (_, i) {
                            final c = categoryProvider.categories[i];
                            final selected =
                                categoryProvider.selectedCategory?.id == c.id;
                            return InkWell(
                              onTap: () => categoryProvider.selectCategory(c),
                              child: AnimatedContainer(
                                duration: const Duration(milliseconds: 160),
                                decoration: BoxDecoration(
                                  color: selected
                                      ? DesignTokens.brandPrimarySoft
                                      : null,
                                  border: selected
                                      ? const Border(
                                          left: BorderSide(
                                              color: DesignTokens.brandPrimary,
                                              width: 4))
                                      : const Border(
                                          left: BorderSide(
                                              color: Colors.transparent,
                                              width: 4)),
                                ),
                                padding: const EdgeInsets.symmetric(
                                    horizontal: 14, vertical: 10),
                                height: DesignTokens.categoryItemHeight,
                                child: Row(
                                  children: [
                                    Expanded(
                                      child: Text(
                                        c.name,
                                        maxLines: 1,
                                        overflow: TextOverflow.ellipsis,
                                        style: (selected
                                            ? DesignTokens.subtitle.copyWith(
                                                fontWeight: FontWeight.w600,
                                                color: DesignTokens
                                                    .brandPrimaryHover)
                                            : DesignTokens.subtitle.copyWith(
                                                fontWeight: FontWeight.w500,
                                                color:
                                                    DesignTokens.neutral700)),
                                      ),
                                    ),
                                    IconButton(
                                      icon: const Icon(Icons.edit,
                                          size: 18,
                                          color: DesignTokens.neutral500),
                                      tooltip: UITexts.dishCategoryEdit,
                                      padding: EdgeInsets.zero,
                                      onPressed: () async {
                                        final catProvider = context
                                            .read<DishCategoryProvider>();
                                        final newName =
                                            await showDialog<String>(
                                          context: context,
                                          builder: (_) =>
                                              ChangeNotifierProvider.value(
                                            value: catProvider,
                                            child: DishCategoryFormDialog(
                                                editing: c),
                                          ),
                                        );
                                        if (newName != null &&
                                            newName.trim().isNotEmpty &&
                                            newName.trim() != c.name) {
                                          final ok =
                                              await catProvider.updateCategory(
                                                  c, newName.trim());
                                          if (!mounted) return;
                                          if (!ok) {
                                            ScaffoldMessenger.of(context)
                                                .showSnackBar(const SnackBar(
                                                    content: Text(ErrorTexts
                                                        .updateCategory)));
                                          }
                                        }
                                      },
                                    ),
                                    IconButton(
                                      icon: const Icon(Icons.delete_outline,
                                          size: 18, color: DesignTokens.danger),
                                      tooltip: UITexts.dishCategoryDelete,
                                      padding: EdgeInsets.zero,
                                      onPressed: () async {
                                        final confirm = await showDialog<bool>(
                                          context: context,
                                          builder: (_) => AlertDialog(
                                            title: const Text(
                                                UITexts.dishCategoryDelete),
                                            content: const Text(UITexts
                                                .dishCategoryDeleteConfirm),
                                            actions: [
                                              TextButton(
                                                  onPressed: () =>
                                                      Navigator.pop(
                                                          context, false),
                                                  child: const Text(
                                                      UITexts.commonCancel)),
                                              ElevatedButton(
                                                  onPressed: () =>
                                                      Navigator.pop(
                                                          context, true),
                                                  child: const Text(
                                                      UITexts.commonOk)),
                                            ],
                                          ),
                                        );
                                        if (confirm == true) {
                                          final dishProvider =
                                              context.read<DishListProvider>();
                                          final ok = await context
                                              .read<DishCategoryProvider>()
                                              .deleteCategory(
                                            c,
                                            hasDishes: (category) {
                                              // Check if this category has dishes
                                              final selectedCategory = context
                                                  .read<DishCategoryProvider>()
                                                  .selectedCategory;
                                              if (selectedCategory?.id ==
                                                  category.id) {
                                                return dishProvider
                                                    .dishes.isNotEmpty;
                                              }
                                              // If not currently selected, assume safe to delete
                                              // Backend will validate
                                              return false;
                                            },
                                          );
                                          if (!mounted) return;
                                          if (!ok) {
                                            final errorMsg = context
                                                    .read<
                                                        DishCategoryProvider>()
                                                    .error ??
                                                ErrorTexts.deleteCategory;
                                            ScaffoldMessenger.of(context)
                                                .showSnackBar(SnackBar(
                                                    content: Text(errorMsg)));
                                          }
                                        }
                                      },
                                    ),
                                  ],
                                ),
                              ),
                            );
                          },
                        ),
            ),
          ],
        ),
      );
    });
  }

  Widget _buildDishArea() {
    return Builder(builder: (context) {
      final theme = Theme.of(context);
      final storeSelector = context.watch<StoreSelectorProvider>();
      final categoryProvider = context.watch<DishCategoryProvider>();
      final dishList = context.watch<DishListProvider>();

      // Show dish list error as SnackBar
      if (dishList.error != null) {
        WidgetsBinding.instance.addPostFrameCallback((_) {
          if (context.mounted) {
            ScaffoldMessenger.of(context).showSnackBar(
              SnackBar(
                content: Text('${UITexts.dishError}: ${dishList.error}'),
                backgroundColor: DesignTokens.danger,
                duration: const Duration(seconds: 3),
              ),
            );
          }
        });
      }

      return Column(
        children: [
          // 顶部工具栏：门店选择 + 当前分类 + 新增菜品
          Container(
            padding: const EdgeInsets.symmetric(
                horizontal: DesignTokens.spaceLg,
                vertical: DesignTokens.spaceMd),
            decoration: BoxDecoration(
              color: theme.cardColor,
              boxShadow: DesignTokens.shadowSm,
            ),
            child: Row(
              children: [
                Text(UITexts.dishTitle,
                    style: DesignTokens.heading2
                        .copyWith(color: theme.textTheme.bodyLarge?.color)),
                const SizedBox(width: DesignTokens.spaceLg),
                // 只有 admin 角色才显示门店选择器
                if (storeSelector.isAdmin)
                  SizedBox(
                    width: 220,
                    child: storeSelector.loading
                        ? const Align(
                            alignment: Alignment.centerLeft,
                            child: SizedBox(
                                height: 18,
                                width: 18,
                                child:
                                    CircularProgressIndicator(strokeWidth: 2)))
                        : storeSelector.error != null
                            ? Row(
                                children: [
                                  Expanded(
                                    child: Text(
                                      ErrorTexts.loadStores,
                                      style: TextStyle(
                                          color: Colors.red.shade600,
                                          fontSize: 12),
                                    ),
                                  ),
                                  IconButton(
                                    icon: const Icon(Icons.refresh, size: 18),
                                    tooltip: '重试加载门店', // 可后续常量化
                                    onPressed: () => storeSelector.loadStores(),
                                  )
                                ],
                              )
                            : DropdownButton<int>(
                                isExpanded: true,
                                value: storeSelector.selectedStoreId,
                                hint: const Text(UITexts.dishStoreSelectLabel),
                                items: storeSelector.stores
                                    .map((s) => DropdownMenuItem<int>(
                                        value: s.id, child: Text(s.name)))
                                    .toList(),
                                onChanged: (v) {
                                  if (v != null) storeSelector.selectStore(v);
                                },
                              ),
                  ),
                const Spacer(),
                // 临时调试:显示当前权限
                Consumer<PermissionProvider>(
                  builder: (context, permProvider, _) {
                    final hasAdd = permProvider.has(PermissionCodes.dishAdd);
                    final hasEdit = permProvider.has(PermissionCodes.dishEdit);
                    final hasDelete =
                        permProvider.has(PermissionCodes.dishDelete);
                    return Tooltip(
                      message:
                          'Add:$hasAdd Edit:$hasEdit Delete:$hasDelete\nAll permissions: ${permProvider.all.join(", ")}',
                      child: Icon(
                        hasAdd && hasEdit && hasDelete
                            ? Icons.check_circle
                            : Icons.warning,
                        color: hasAdd && hasEdit && hasDelete
                            ? Colors.green
                            : Colors.orange,
                        size: 16,
                      ),
                    );
                  },
                ),
                const SizedBox(width: DesignTokens.spaceMd),
                PermissionGate(
                  required: PermissionCodes.dishAdd,
                  child: ElevatedButton.icon(
                    onPressed: categoryProvider.selectedCategory == null
                        ? null
                        : _openCreateDish,
                    icon: const Icon(Icons.add),
                    style: ElevatedButton.styleFrom(
                        backgroundColor: DesignTokens.brandPrimary,
                        foregroundColor: Colors.white,
                        padding: const EdgeInsets.symmetric(
                            horizontal: 20, vertical: 12),
                        textStyle: DesignTokens.subtitle
                            .copyWith(fontWeight: FontWeight.w600)),
                    label: const Text(UITexts.dishAdd),
                  ),
                )
              ],
            ),
          ),
          Expanded(
            child: categoryProvider.selectedCategory == null
                ? Center(
                    child: Text('请选择左侧分类',
                        style: DesignTokens.subtitle
                            .copyWith(color: DesignTokens.neutral500)))
                : dishList.loading
                    ? Center(
                        child: Text(UITexts.dishLoading,
                            style: DesignTokens.body.copyWith(
                                color: theme.textTheme.bodySmall?.color)))
                    : dishList.dishes.isEmpty
                        ? Center(
                            child: Text(UITexts.dishEmpty,
                                style: DesignTokens.body.copyWith(
                                    color: theme.textTheme.bodySmall?.color)))
                        : _buildDishTable(dishList.dishes),
          )
        ],
      );
    });
  }

  Widget _buildDishTable(List<Dish> list) {
    return Padding(
      padding: const EdgeInsets.all(16),
      child: AdminTable<Dish>(
        columns: const [
          AdminTableColumn(width: 220, label: UITexts.dishColumnName),
          AdminTableColumn(width: 120, label: UITexts.dishColumnPrice),
          AdminTableColumn(width: 100, label: UITexts.dishColumnStatus),
          AdminTableColumn(
              label: UITexts.dishColumnActions,
              alignment: Alignment.centerRight),
        ],
        data: list,
        rowBuilder: (d, i) {
          final active = (d.status ?? 1) == 1;
          return Row(
            children: [
              SizedBox(
                width: 220,
                child: Text(d.name,
                    maxLines: 1,
                    overflow: TextOverflow.ellipsis,
                    style: DesignTokens.body.copyWith(
                        fontWeight: FontWeight.w500,
                        color: DesignTokens.neutral700)),
              ),
              SizedBox(
                width: 120,
                child: Text(d.price?.toStringAsFixed(2) ?? '-',
                    style: DesignTokens.bodySm
                        .copyWith(color: DesignTokens.neutral600)),
              ),
              SizedBox(
                  width: 100,
                  child: StatusTag(
                      active: active,
                      enabledText: UITexts.dishStatusEnabled,
                      disabledText: UITexts.dishStatusDisabled)),
              Expanded(
                child: Row(
                  mainAxisAlignment: MainAxisAlignment.end,
                  children: [
                    PermissionGate(
                      required: PermissionCodes.dishEdit,
                      child: TextButton(
                          onPressed: () async {
                            final dishProvider =
                                context.read<DishListProvider>();
                            final updated = await showDialog<UpdateDishRequest>(
                              context: context,
                              builder: (_) => DishFormDialog(editing: d),
                            );
                            if (updated != null) {
                              final ok =
                                  await dishProvider.updateDish(d, updated);
                              if (!context.mounted) return;
                              ScaffoldMessenger.of(context).showSnackBar(
                                SnackBar(
                                  content: Text(ok
                                      ? UITexts.dishUpdateSuccess
                                      : UITexts.dishUpdateFailed),
                                ),
                              );
                            }
                          },
                          style: TextButton.styleFrom(
                              foregroundColor: DesignTokens.brandPrimary,
                              textStyle: DesignTokens.body
                                  .copyWith(fontWeight: FontWeight.w500)),
                          child: const Text(UITexts.dishEdit)),
                    ),
                    PermissionGate(
                      required: PermissionCodes.dishDelete,
                      child: TextButton(
                        onPressed: () async {
                          final dishProvider = context.read<DishListProvider>();
                          final confirm = await showDialog<bool>(
                            context: context,
                            builder: (ctx) => AlertDialog(
                              title: Text(UITexts.dishDeleteConfirmTitle),
                              content: Text(
                                  '${UITexts.dishDeleteConfirmMessage} "${d.name}"?'),
                              actions: [
                                TextButton(
                                    onPressed: () => Navigator.pop(ctx, false),
                                    child: const Text(UITexts.commonCancel)),
                                TextButton(
                                    onPressed: () => Navigator.pop(ctx, true),
                                    style: TextButton.styleFrom(
                                        foregroundColor: DesignTokens.danger),
                                    child: const Text(UITexts.commonOk)),
                              ],
                            ),
                          );
                          if (confirm == true) {
                            final ok = await dishProvider.deleteDish(d);
                            if (!context.mounted) return;
                            ScaffoldMessenger.of(context).showSnackBar(
                              SnackBar(
                                content: Text(ok
                                    ? UITexts.dishDeleteSuccess
                                    : UITexts.dishDeleteFailed),
                              ),
                            );
                          }
                        },
                        style: TextButton.styleFrom(
                            foregroundColor: DesignTokens.danger,
                            textStyle: DesignTokens.body
                                .copyWith(fontWeight: FontWeight.w500)),
                        child: const Text(UITexts.dishDelete),
                      ),
                    ),
                  ],
                ),
              )
            ],
          );
        },
      ),
    );
  }
}

/// 菜品表单对话框（简化版）
class DishFormDialog extends StatefulWidget {
  final Dish? editing;
  const DishFormDialog({super.key, this.editing});
  @override
  State<DishFormDialog> createState() => _DishFormDialogState();
}

class _DishFormDialogState extends State<DishFormDialog> {
  final _formKey = GlobalKey<FormState>();
  final _nameCtrl = TextEditingController();
  final _descCtrl = TextEditingController();
  final _priceCtrl = TextEditingController();
  bool _status = true;

  @override
  void initState() {
    super.initState();
    if (widget.editing != null) {
      final d = widget.editing!;
      _nameCtrl.text = d.name;
      _descCtrl.text = d.description ?? '';
      _priceCtrl.text = d.price?.toString() ?? '';
      _status = (d.status ?? 1) == 1;
    }
  }

  @override
  void dispose() {
    _nameCtrl.dispose();
    _descCtrl.dispose();
    _priceCtrl.dispose();
    super.dispose();
  }

  void _submit() {
    if (!_formKey.currentState!.validate()) return;
    if (widget.editing != null) {
      // Edit mode: return UpdateDishRequest
      final req = UpdateDishRequest(
        name: _nameCtrl.text.trim(),
        description:
            _descCtrl.text.trim().isEmpty ? null : _descCtrl.text.trim(),
        price: parseDoubleNullable(_priceCtrl.text.trim()),
        status: _status ? 1 : 0,
      );
      Navigator.pop(context, req);
    } else {
      // Create mode: return CreateDishRequest
      final req = CreateDishRequest(
        name: _nameCtrl.text.trim(),
        description:
            _descCtrl.text.trim().isEmpty ? null : _descCtrl.text.trim(),
        price: parseDoubleNullable(_priceCtrl.text.trim()),
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
                Text(
                    widget.editing == null ? UITexts.dishAdd : UITexts.dishEdit,
                    style: const TextStyle(
                        fontSize: 20, fontWeight: FontWeight.w600)),
                const SizedBox(height: 16),
                TextFormField(
                  controller: _nameCtrl,
                  decoration:
                      const InputDecoration(labelText: UITexts.dishFormName),
                  maxLength: 64,
                  validator: (v) {
                    final value = v?.trim() ?? '';
                    if (value.isEmpty) return '必填';
                    if (value.length < 2) return '至少 2 个字符';
                    return null;
                  },
                ),
                const SizedBox(height: 12),
                TextFormField(
                  controller: _descCtrl,
                  decoration:
                      const InputDecoration(labelText: UITexts.dishFormDesc),
                  maxLines: 3,
                  maxLength: 200,
                ),
                const SizedBox(height: 12),
                TextFormField(
                  controller: _priceCtrl,
                  keyboardType:
                      const TextInputType.numberWithOptions(decimal: true),
                  decoration:
                      const InputDecoration(labelText: UITexts.dishFormPrice),
                  validator: (v) {
                    final value = v?.trim() ?? '';
                    if (value.isEmpty) return null; // 可选
                    return double.tryParse(value) == null ? '请输入合法数字或留空' : null;
                  },
                ),
                const SizedBox(height: 12),
                Row(
                  children: [
                    Switch(
                        value: _status,
                        onChanged: (v) => setState(() => _status = v)),
                    const Text(UITexts.dishFormStatus),
                  ],
                ),
                const SizedBox(height: 20),
                Row(
                  children: [
                    Expanded(
                      child: OutlinedButton(
                        onPressed: () => Navigator.pop(context),
                        child: const Text(UITexts.commonCancel),
                      ),
                    ),
                    const SizedBox(width: 16),
                    Expanded(
                      child: ElevatedButton(
                        onPressed: _submit,
                        child: const Text(UITexts.commonOk),
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

class DishManagementScope extends StatelessWidget {
  const DishManagementScope({super.key});
  @override
  Widget build(BuildContext context) {
    final client = ApiClient();
    return MultiProvider(
      providers: [
        // 门店 Provider
        ChangeNotifierProvider<StoreSelectorProvider>(
          create: (_) => StoreSelectorProvider(StoreApi(client)),
        ),
        // 分类 Provider 依赖门店：直接读取已存在的 StoreSelectorProvider 实例，避免重复实例导致状态不同步
        ChangeNotifierProvider<DishCategoryProvider>(
          create: (ctx) => DishCategoryProvider(
              DishCategoryApi(client), ctx.read<StoreSelectorProvider>()),
        ),
        // 菜品列表 Provider 依赖门店与分类：同理复用实例
        ChangeNotifierProvider<DishListProvider>(
          create: (ctx) => DishListProvider(
              DishApi(client),
              ctx.read<StoreSelectorProvider>(),
              ctx.read<DishCategoryProvider>()),
        ),
      ],
      child: const DishManagementPage(),
    );
  }
}

/// 分类新增/编辑表单对话框
class DishCategoryFormDialog extends StatefulWidget {
  final DishCategory? editing;
  const DishCategoryFormDialog({super.key, this.editing});
  @override
  State<DishCategoryFormDialog> createState() => _DishCategoryFormDialogState();
}

class _DishCategoryFormDialogState extends State<DishCategoryFormDialog> {
  final _formKey = GlobalKey<FormState>();
  final _nameCtrl = TextEditingController();

  @override
  void initState() {
    super.initState();
    if (widget.editing != null) {
      _nameCtrl.text = widget.editing!.name;
    }
  }

  @override
  void dispose() {
    _nameCtrl.dispose();
    super.dispose();
  }

  void _submit() {
    if (!_formKey.currentState!.validate()) return;
    Navigator.pop(context, _nameCtrl.text.trim());
  }

  @override
  Widget build(BuildContext context) {
    return Dialog(
      child: ConstrainedBox(
        constraints: const BoxConstraints(maxWidth: 420),
        child: Padding(
          padding: const EdgeInsets.all(24),
          child: Form(
            key: _formKey,
            child: Column(
              mainAxisSize: MainAxisSize.min,
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text(
                  widget.editing == null
                      ? UITexts.dishCategoryAdd
                      : UITexts.dishCategoryEdit,
                  style: const TextStyle(
                      fontSize: 18, fontWeight: FontWeight.w600),
                ),
                const SizedBox(height: 16),
                TextFormField(
                  controller: _nameCtrl,
                  decoration: const InputDecoration(
                      labelText: UITexts.dishCategoryFormName),
                  maxLength: 40,
                  validator: (v) {
                    final value = v?.trim() ?? '';
                    if (value.isEmpty) return '必填';
                    if (value.length < 2) return '至少 2 个字符';
                    // 唯一性校验（同门店下名称不重复）
                    final provider = context.read<DishCategoryProvider>();
                    final editingId = widget.editing?.id;
                    final exists = provider.categories
                        .any((c) => c.name == value && c.id != editingId);
                    if (exists) return '该名称已存在';
                    return null;
                  },
                ),
                const SizedBox(height: 20),
                Row(
                  children: [
                    Expanded(
                      child: OutlinedButton(
                        onPressed: () => Navigator.pop(context),
                        child: const Text(UITexts.commonCancel),
                      ),
                    ),
                    const SizedBox(width: 16),
                    Expanded(
                      child: ElevatedButton(
                        onPressed: _submit,
                        child: const Text(UITexts.commonOk),
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
