import 'package:flutter/material.dart';

class AdminTableColumn {
  final double? width;
  final String label;
  final Alignment alignment;
  const AdminTableColumn(
      {this.width, required this.label, this.alignment = Alignment.centerLeft});
}

class AdminTable<T> extends StatelessWidget {
  final List<AdminTableColumn> columns;
  final List<T> data;
  final Widget Function(T item, int index) rowBuilder;
  final double headerHeight;
  final double rowHeight;
  final Color? stripeColor;
  final EdgeInsets cellPadding;
  // 分页相关（可选）：如果 total/page/pageSize 都提供则显示分页条
  final int? page;
  final int? pageSize;
  final int? total;
  final ValueChanged<int>? onPageChange; // 传入新的页码（从 1 开始）
  // 搜索钩子（可选）：仅提供基础输入框，不做去抖；交由外部处理
  final bool enableSearch;
  final String searchPlaceholder;
  final ValueChanged<String>? onSearch;
  final TextEditingController? searchController;

  const AdminTable({
    super.key,
    required this.columns,
    required this.data,
    required this.rowBuilder,
    this.headerHeight = 40,
    this.rowHeight = 44,
    this.stripeColor,
    this.cellPadding = const EdgeInsets.symmetric(horizontal: 12, vertical: 10),
    this.page,
    this.pageSize,
    this.total,
    this.onPageChange,
    this.enableSearch = false,
    this.searchPlaceholder = '搜索...',
    this.onSearch,
    this.searchController,
  });

  @override
  Widget build(BuildContext context) {
    return Container(
      decoration: BoxDecoration(
        color: Colors.white,
        borderRadius: BorderRadius.circular(8),
        boxShadow: [
          BoxShadow(
              color: Colors.black.withOpacity(.05),
              blurRadius: 6,
              offset: const Offset(0, 2)),
        ],
      ),
      child: Column(
        children: [
          if (enableSearch) _buildSearchBar(),
          _buildHeader(),
          Divider(height: 1, color: Colors.grey.shade200),
          Expanded(
            child: ListView.separated(
              itemCount: data.length,
              separatorBuilder: (_, __) =>
                  Divider(height: 1, color: Colors.grey.shade100),
              itemBuilder: (_, i) {
                final item = data[i];
                return Container(
                  height: rowHeight,
                  color: i.isEven
                      ? Colors.white
                      : (stripeColor ?? Colors.grey.shade50),
                  child: Padding(
                    padding: cellPadding,
                    child: rowBuilder(item, i),
                  ),
                );
              },
            ),
          ),
          if (_showPagination) Divider(height: 1, color: Colors.grey.shade200),
          if (_showPagination) _buildPaginationBar(),
        ],
      ),
    );
  }

  Widget _buildHeader() {
    return Container(
      height: headerHeight,
      padding: cellPadding,
      child: Row(
        children: [
          for (final c in columns)
            SizedBox(
              width: c.width,
              child: Align(
                alignment: c.alignment,
                child: Text(c.label,
                    style: const TextStyle(
                        fontSize: 12, fontWeight: FontWeight.w600)),
              ),
            ),
          if (columns.isEmpty) const SizedBox.shrink(),
        ],
      ),
    );
  }

  bool get _showPagination =>
      page != null && pageSize != null && total != null && onPageChange != null;

  Widget _buildPaginationBar() {
    final current = page!;
    final size = pageSize!;
    final totalItems = total!;
    final totalPages =
        totalItems == 0 ? 1 : ((totalItems + size - 1) / size).ceil();
    final canPrev = current > 1;
    final canNext = current < totalPages;
    return Container(
      height: 48,
      padding: const EdgeInsets.symmetric(horizontal: 16),
      child: Row(
        children: [
          Text('共 $totalItems 条',
              style: TextStyle(fontSize: 12, color: Colors.grey.shade600)),
          const Spacer(),
          IconButton(
            icon: const Icon(Icons.chevron_left),
            tooltip: '上一页',
            onPressed: canPrev ? () => onPageChange!(current - 1) : null,
          ),
          Text('$current / $totalPages', style: const TextStyle(fontSize: 12)),
          IconButton(
            icon: const Icon(Icons.chevron_right),
            tooltip: '下一页',
            onPressed: canNext ? () => onPageChange!(current + 1) : null,
          ),
        ],
      ),
    );
  }

  Widget _buildSearchBar() {
    final controller = searchController ?? TextEditingController();
    return Container(
      padding: const EdgeInsets.fromLTRB(12, 12, 12, 8),
      child: TextField(
        controller: controller,
        decoration: InputDecoration(
          isDense: true,
          hintText: searchPlaceholder,
          prefixIcon: const Icon(Icons.search, size: 18),
          filled: true,
          fillColor: Colors.grey.shade100,
          border: OutlineInputBorder(
            borderSide: BorderSide.none,
            borderRadius: BorderRadius.circular(8),
          ),
        ),
        onSubmitted: (value) => onSearch?.call(value.trim()),
      ),
    );
  }
}
