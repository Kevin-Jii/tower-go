import 'package:flutter/material.dart';
import '../../../core/routes/route_manager.dart';
import '../../menu/models.dart';

/// 菜单内容展示区域
class MenuContent extends StatelessWidget {
  final MenuItem menuItem;

  const MenuContent({super.key, required this.menuItem});

  @override
  Widget build(BuildContext context) {
    // 使用路由管理器获取对应页面
    final page = RouteManager().getPageForMenuItem(menuItem);

    if (page != null) {
      return page;
    }

    // 如果没有对应页面，显示默认提示
    return Center(
      child: Column(
        mainAxisSize: MainAxisSize.min,
        children: [
          Icon(Icons.dashboard, size: 48, color: Colors.grey[400]),
          const SizedBox(height: 16),
          Text(
            '当前菜单: ${menuItem.title}',
            style: Theme.of(context).textTheme.headlineSmall,
          ),
          const SizedBox(height: 8),
          Text(
            '路由: ${menuItem.path ?? "无"}',
            style: Theme.of(context).textTheme.bodyMedium?.copyWith(
                  color: Colors.grey[600],
                ),
          ),
        ],
      ),
    );
  }
}
