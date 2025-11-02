import 'package:flutter/material.dart';
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
    final tile = ListTile(
      dense: true,
      selected: selected,
      leading: hasChildren
          ? GestureDetector(
              onTap: () => setState(() => _expanded = !_expanded),
              child: Icon(
                _expanded ? Icons.folder_open : Icons.folder,
                size: 18,
                color: selected ? Theme.of(context).colorScheme.primary : null,
              ),
            )
          : TdIconMapper.build(
              widget.item.icon,
              size: 18,
              color: selected ? Theme.of(context).colorScheme.primary : null,
            ),
      title: Text(
        widget.item.title,
        style: TextStyle(
          fontWeight: selected ? FontWeight.w600 : FontWeight.normal,
          color: selected
              ? Theme.of(context).colorScheme.onPrimary
              : Theme.of(context).colorScheme.onSurface,
        ),
        overflow: TextOverflow.ellipsis,
      ),
      onTap: () {
        widget.onTap(widget.item);
      },
      contentPadding: EdgeInsets.only(left: 8.0 + widget.depth * 12, right: 8),
      trailing: hasChildren
          ? IconButton(
              icon: Icon(_expanded ? Icons.expand_less : Icons.expand_more),
              onPressed: () => setState(() => _expanded = !_expanded),
            )
          : null,
    );
    Widget styledTile = Theme(
      data: Theme.of(context).copyWith(
        listTileTheme: ListTileThemeData(
          selectedColor: Theme.of(context).colorScheme.onPrimary,
          selectedTileColor:
              Theme.of(context).colorScheme.primary.withOpacity(0.85),
        ),
      ),
      child: tile,
    );

    if (!hasChildren) return styledTile;
    return Column(
      crossAxisAlignment: CrossAxisAlignment.stretch,
      children: [
        styledTile,
        if (_expanded)
          ...widget.item.children.map((c) =>
              _MenuNode(item: c, depth: widget.depth + 1, onTap: widget.onTap)),
      ],
    );
  }
}
// 类型标签组件已移除，目录/菜单不再显示额外徽章
