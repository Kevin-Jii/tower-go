import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import 'package:tdesign_flutter/tdesign_flutter.dart';
import 'store_provider.dart';
import 'models.dart';
import 'store_form_dialog.dart';

class StoreManagementPage extends StatefulWidget {
  const StoreManagementPage({super.key});

  @override
  State<StoreManagementPage> createState() => _StoreManagementPageState();
}

class _StoreManagementPageState extends State<StoreManagementPage> {
  final Set<int> _updatingStores = {}; // 记录正在更新的门店ID

  @override
  void initState() {
    super.initState();
    WidgetsBinding.instance.addPostFrameCallback((_) {
      context.read<StoreProvider>().loadStores();
    });
  }

  void _handleCreate() async {
    final result = await showDialog<CreateStoreRequest>(
      context: context,
      builder: (ctx) => const StoreFormDialog(),
    );
    if (result != null && mounted) {
      final success = await context.read<StoreProvider>().createStore(result);
      if (success && mounted) {
        TDMessage.showMessage(
          context: context,
          content: '创建成功',
          theme: MessageTheme.success,
          duration: 2000,
        );
      } else if (mounted) {
        final err = context.read<StoreProvider>().error ?? '创建失败';
        TDMessage.showMessage(
          context: context,
          content: err,
          theme: MessageTheme.error,
          duration: 3000,
        );
      }
    }
  }

  void _handleStatusChange(Store store, bool value) async {
    // 标记为正在更新
    setState(() => _updatingStores.add(store.id));

    final newStatus = value ? 1 : 0;
    final req = UpdateStoreRequest(
      status: newStatus,
    );

    final success =
        await context.read<StoreProvider>().updateStore(store.id, req);

    // 移除更新标记
    if (mounted) {
      setState(() => _updatingStores.remove(store.id));
    }

    if (success && mounted) {
      TDMessage.showMessage(
        context: context,
        content: '状态已更新',
        theme: MessageTheme.success,
        duration: 2000,
      );
    } else if (mounted) {
      final err = context.read<StoreProvider>().error ?? '状态更新失败';
      TDMessage.showMessage(
        context: context,
        content: err,
        theme: MessageTheme.error,
        duration: 3000,
      );
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
            '门店管理',
            style: TextStyle(fontSize: 20, fontWeight: FontWeight.w600),
          ),
          const Spacer(),
          TDButton(
            text: '新增门店',
            theme: TDButtonTheme.primary,
            size: TDButtonSize.medium,
            icon: TDIcons.add,
            onTap: _handleCreate,
          ),
        ],
      ),
    );
  }

  Widget _buildTable() {
    return Consumer<StoreProvider>(
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
                  onTap: () => provider.loadStores(),
                ),
              ],
            ),
          );
        }

        final columns = <TDTableCol>[
          TDTableCol(title: 'ID', colKey: 'id', width: 80, ellipsis: true),
          TDTableCol(title: '门店名称', colKey: 'name', ellipsis: true),
          TDTableCol(title: '地址', colKey: 'address', ellipsis: true),
          TDTableCol(title: '营业时间', colKey: 'businessHours', ellipsis: true),
          TDTableCol(title: '门店负责人', colKey: 'contactPerson', ellipsis: true),
          TDTableCol(title: '负责人电话', colKey: 'iphone', ellipsis: true),
          TDTableCol(
            title: '门店状态',
            colKey: 'status',
            width: 140,
            cellBuilder: (ctx, index) {
              final store = provider.stores[index];
              final isActive = store.status == 1;
              final isUpdating = _updatingStores.contains(store.id);

              if (isUpdating) {
                return Row(
                  mainAxisSize: MainAxisSize.min,
                  children: [
                    const SizedBox(
                      width: 16,
                      height: 16,
                      child: CircularProgressIndicator(
                        strokeWidth: 2,
                      ),
                    ),
                    const SizedBox(width: 8),
                    Text(
                      isActive ? '营业中' : '停业中',
                      style: TextStyle(
                        fontSize: 14,
                        color: Colors.grey.shade600,
                      ),
                    ),
                  ],
                );
              }

              return Row(
                mainAxisSize: MainAxisSize.min,
                children: [
                  TDSwitch(
                    isOn: isActive,
                    type: TDSwitchType.icon,
                    size: TDSwitchSize.small,
                    onChanged: (value) {
                      _handleStatusChange(store, value);
                      return value;
                    },
                  ),
                  const SizedBox(width: 8),
                  Text(
                    isActive ? '营业中' : '停业中',
                    style: TextStyle(
                      fontSize: 14,
                      color: isActive
                          ? Colors.green.shade600
                          : Colors.grey.shade600,
                      fontWeight: FontWeight.w500,
                    ),
                  ),
                ],
              );
            },
          ),
          TDTableCol(title: '备注', colKey: 'remark', ellipsis: true),
        ];

        final data = provider.stores
            .map((s) => {
                  'id': s.id.toString(),
                  'name': s.name,
                  'address': s.address ?? '—',
                  'businessHours': s.businessHours ?? '—',
                  'contactPerson': s.contactPerson ?? '—',
                  'iphone': s.phone ?? '—',
                  'status': s.status == 1 ? '营业中' : '停业中',
                  'remark': s.remark ?? '—',
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
                          '暂无门店',
                          style: TextStyle(fontSize: 16, color: Colors.grey),
                        ),
                      ],
                    ),
                  )
                : ClipRRect(
                    borderRadius: BorderRadius.circular(8),
                    child: SingleChildScrollView(
                      child: TDTable(
                        columns: columns,
                        data: data,
                        rowHeight: 50,
                      ),
                    ),
                  ),
          ),
        );
      },
    );
  }

  Widget _buildTableSkeleton({int rows = 8}) {
    if (rows <= 0) rows = 1;
    if (rows > 50) rows = 50;
    return Padding(
      padding: const EdgeInsets.all(20),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
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
                _skeletonBlock(width: 120, height: 14),
                const SizedBox(width: 40),
                _skeletonBlock(width: 100, height: 14),
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
                        horizontal: 20, vertical: 18),
                    child: Row(
                      children: [
                        _skeletonBlock(width: 60, height: 16),
                        const SizedBox(width: 48),
                        _skeletonBlock(width: 150, height: 16),
                        const SizedBox(width: 48),
                        _skeletonBlock(width: 200, height: 16),
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
    return Consumer<StoreProvider>(
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
                        ? () => provider.loadStores(page: provider.page - 1)
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
                        ? () => provider.loadStores(page: provider.page + 1)
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
