import 'package:flutter/material.dart';
import '../../core/theme/tower_colors.dart';
import 'package:provider/provider.dart';
import 'menu_provider.dart';
import 'models.dart';
import '../../core/icons/td_icon_mapper.dart';

class MenuTree extends StatelessWidget {
  final void Function(MenuItem item)? onSelect;
  const MenuTree({super.key, this.onSelect});

  @override
  Widget build(BuildContext context) {
    return Consumer<MenuProvider>(
      builder: (context, mp, _) {
        if (mp.loading) {
          return const Center(child: CircularProgressIndicator());
        }
        if (mp.error != null) {
          return Center(
            child: Column(
              mainAxisSize: MainAxisSize.min,
              children: [
                Text('加载失败: ${mp.error}',
                    style: const TextStyle(color: Colors.red)),
                const SizedBox(height: 8),
                FilledButton(onPressed: mp.load, child: const Text('重试')),
              ],
            ),
          );
        }
        return ListView(
          padding: const EdgeInsets.symmetric(vertical: 4),
          children: mp.tree
              .map((e) => _MenuNode(
                  item: e,
                  depth: 0,
                  onTap: (m) {
                    mp.select(m);
                    onSelect?.call(m);
                  }))
              .toList(),
        );
      },
    );
  }
}

class _MenuNode extends StatefulWidget {
  final MenuItem item;
  final int depth;
  final void Function(MenuItem) onTap;
  const _MenuNode(
      {required this.item, required this.depth, required this.onTap});

  @override
  State<_MenuNode> createState() => _MenuNodeState();
}

class _MenuNodeState extends State<_MenuNode> {
  bool _expanded = true;

  @override
  Widget build(BuildContext context) {
    final hasChildren = widget.item.children.isNotEmpty;
    final mp = context.read<MenuProvider>();
    final selected = mp.selected?.id == widget.item.id;
    final isDirectory = hasChildren; // 目录节点
    final canSelect = !isDirectory; // 只有叶子(页面)才触发 select

    final tc = Theme.of(context).extension<TowerColors>();
    final accent = Theme.of(context).colorScheme.primary;
    final tile = InkWell(
      borderRadius: BorderRadius.circular(8),
      onTap: () {
        if (isDirectory) {
          setState(() => _expanded = !_expanded);
        } else {
          widget.onTap(widget.item);
        }
      },
      child: Padding(
        padding: EdgeInsets.only(
          left: 6 + widget.depth * 14,
          right: 8,
          top: 4,
          bottom: 4,
        ),
        child: Row(
          children: [
            // 图标 / 展开指示
            if (isDirectory)
              AnimatedRotation(
                turns: _expanded ? 0.0 : -0.25,
                duration: const Duration(milliseconds: 200),
                child: Icon(
                  Icons.expand_more,
                  size: 18,
                  color: tc?.navBarForeground.withOpacity(.6) ??
                      Colors.blueGrey.shade600,
                ),
              )
            else
              TdIconMapper.build(
                widget.item.icon,
                size: 18,
                color: selected
                    ? (tc?.accentGradientStart ?? accent)
                    : tc?.navBarForeground.withOpacity(.65) ??
                        Colors.blueGrey.shade600,
              ),
            const SizedBox(width: 8),
            Expanded(
              child: Text(
                widget.item.title,
                style: TextStyle(
                  fontSize: 13.5,
                  fontWeight: selected ? FontWeight.w600 : FontWeight.w400,
                  color: selected
                      ? (tc?.accentGradientStart ?? accent)
                      : tc?.navBarForeground.withOpacity(.75) ??
                          Colors.blueGrey.shade700,
                ),
                overflow: TextOverflow.ellipsis,
              ),
            ),
            if (canSelect && selected)
              Container(
                width: 6,
                height: 6,
                decoration: BoxDecoration(
                  color: tc?.accentGradientStart ?? accent,
                  borderRadius: BorderRadius.circular(3),
                ),
              ),
          ],
        ),
      ),
    );
    Widget styledTile = Theme(
      data: Theme.of(context).copyWith(
        listTileTheme: ListTileThemeData(
          selectedColor: (tc?.navBarForeground ?? Colors.white),
          selectedTileColor:
              (tc?.accentGradientStart ?? accent).withOpacity(0.18),
        ),
      ),
      child: tile,
    );

    if (!hasChildren) return styledTile;
    return Column(
      crossAxisAlignment: CrossAxisAlignment.stretch,
      children: [
        styledTile,
        AnimatedCrossFade(
          firstChild: const SizedBox.shrink(),
          secondChild: Column(
            children: widget.item.children
                .map((c) => _MenuNode(
                    item: c, depth: widget.depth + 1, onTap: widget.onTap))
                .toList(),
          ),
          crossFadeState:
              _expanded ? CrossFadeState.showSecond : CrossFadeState.showFirst,
          duration: const Duration(milliseconds: 180),
        ),
      ],
    );
  }
}
// 类型标签组件已移除，目录/菜单不再显示额外徽章
