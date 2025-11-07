import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import 'dart:convert';
import '../../core/theme/design_tokens.dart';
import '../../core/widgets/status_tag.dart';
import 'dingtalk_provider.dart';
import 'dingtalk_api.dart';
import '../../core/network/api_client.dart';

class DingTalkManagementPage extends StatefulWidget {
  const DingTalkManagementPage({super.key});

  @override
  State<DingTalkManagementPage> createState() => _DingTalkManagementPageState();
}

class _DingTalkManagementPageState extends State<DingTalkManagementPage> {
  @override
  void initState() {
    super.initState();
    WidgetsBinding.instance.addPostFrameCallback((_) {
      context.read<DingTalkProvider>().loadRobots();
    });
  }

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    return Scaffold(
      backgroundColor: theme.scaffoldBackgroundColor,
      body: Column(
        children: [
          _buildHeader(),
          Expanded(child: _buildContent()),
        ],
      ),
    );
  }

  Widget _buildHeader() {
    final theme = Theme.of(context);
    return Container(
      padding: const EdgeInsets.all(DesignTokens.spaceLg),
      decoration: BoxDecoration(
        color: theme.cardColor,
        boxShadow: DesignTokens.shadowSm,
      ),
      child: Row(
        children: [
          Text(
            '钉钉机器人管理',
            style: DesignTokens.heading2
                .copyWith(color: theme.textTheme.bodyLarge?.color),
          ),
          const Spacer(),
          ElevatedButton.icon(
            onPressed: () {},
            icon: const Icon(Icons.add),
            label: const Text('新增机器人'),
            style: ElevatedButton.styleFrom(
              backgroundColor: DesignTokens.brandPrimary,
              foregroundColor: Colors.white,
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildContent() {
    return Consumer<DingTalkProvider>(
      builder: (context, provider, _) {
        if (provider.loading) {
          return const Center(child: CircularProgressIndicator());
        }
        if (provider.error != null) {
          return Center(
            child: Column(
              mainAxisAlignment: MainAxisAlignment.center,
              children: [
                Text(provider.error!,
                    style: TextStyle(color: Colors.red.shade600)),
                const SizedBox(height: 16),
                ElevatedButton(
                  onPressed: () => provider.loadRobots(),
                  child: const Text('重试'),
                ),
              ],
            ),
          );
        }
        if (provider.robots.isEmpty) {
          return const Center(child: Text('暂无钉钉机器人配置'));
        }
        return _buildCardGrid(provider.robots);
      },
    );
  }

  Widget _buildCardGrid(List<dynamic> robots) {
    return Padding(
      padding: const EdgeInsets.all(DesignTokens.spaceLg),
      child: GridView.builder(
        gridDelegate: const SliverGridDelegateWithFixedCrossAxisCount(
          crossAxisCount: 3,
          childAspectRatio: 1.3, // 调小比例，卡片更贴合内容
          crossAxisSpacing: DesignTokens.spaceMd,
          mainAxisSpacing: DesignTokens.spaceMd,
        ),
        itemCount: robots.length,
        itemBuilder: (context, index) {
          final robot = robots[index];
          return Align(
            alignment: Alignment.topCenter,
            child: IntrinsicHeight(
              child: _buildRobotCard(robot),
            ),
          );
        },
      ),
    );
  }

  Widget _buildRobotCard(dynamic robot) {
    return Card(
      elevation: 8,
      shape: RoundedRectangleBorder(
        borderRadius: BorderRadius.circular(18),
      ),
      shadowColor: Colors.blue.shade100,
      child: Stack(
        children: [
          Container(
            padding:
                const EdgeInsets.symmetric(horizontal: 16, vertical: 12), // 更紧凑
            decoration: BoxDecoration(
              borderRadius: BorderRadius.circular(18),
              gradient: LinearGradient(
                begin: Alignment.topLeft,
                end: Alignment.bottomRight,
                // 修复：添加缺失的 colors 参数
                colors: [Colors.white, Colors.blue.shade50],
              ),
            ),
            child: Column(
              mainAxisSize: MainAxisSize.min,
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                // 图标+名称+状态
                Row(
                  crossAxisAlignment: CrossAxisAlignment.center,
                  children: [
                    Container(
                      width: 40,
                      height: 40,
                      decoration: BoxDecoration(
                        color: robot.isEnabled
                            ? Colors.green.shade100
                            : Colors.red.shade100,
                        borderRadius: BorderRadius.circular(20),
                      ),
                      child: Icon(
                        Icons.smart_toy,
                        color: robot.isEnabled
                            ? Colors.green.shade600
                            : Colors.red.shade600,
                        size: 24,
                      ),
                    ),
                    const SizedBox(width: 14),
                    Expanded(
                      child: Column(
                        crossAxisAlignment: CrossAxisAlignment.start,
                        children: [
                          Text(robot.name,
                              style: const TextStyle(
                                  fontSize: 16,
                                  fontWeight: FontWeight.bold,
                                  color: Color(0xFF222B45)),
                              maxLines: 2,
                              overflow: TextOverflow.ellipsis),
                        ],
                      ),
                    ),
                  ],
                ),
                const SizedBox(height: 10),
                // 类型标签
                Container(
                  padding:
                      const EdgeInsets.symmetric(horizontal: 10, vertical: 4),
                  decoration: BoxDecoration(
                      color: Colors.blue.shade100,
                      borderRadius: BorderRadius.circular(10)),
                  child: Text(robot.botType.toUpperCase(),
                      style: TextStyle(
                          fontSize: 12,
                          color: Colors.blue.shade700,
                          fontWeight: FontWeight.w600)),
                ),
                const SizedBox(height: 10),
                // 信息区
                Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    _buildInfoRow(
                        Icons.info_outline, 'ID', robot.id.toString()),
                    if (robot.webhook.isNotEmpty)
                      _buildInfoRow(Icons.link, 'Webhook', robot.webhook),
                    _buildInfoRow(
                        Icons.fingerprint, 'Client ID', robot.clientId),
                    if (robot.agentId != null)
                      _buildInfoRow(
                          Icons.settings, 'Agent ID', robot.agentId.toString()),
                  ],
                ),
                const SizedBox(height: 18),
                // 按钮区
                Row(
                  mainAxisAlignment: MainAxisAlignment.center,
                  children: [
                    OutlinedButton.icon(
                      onPressed: () => _testPush(robot),
                      icon: const Icon(Icons.send, size: 16),
                      label: const Text('测试', style: TextStyle(fontSize: 13)),
                      style: OutlinedButton.styleFrom(
                        padding: const EdgeInsets.symmetric(
                            vertical: 10, horizontal: 18),
                        side: BorderSide(color: Colors.green.shade300),
                        foregroundColor: Colors.green.shade600,
                        shape: RoundedRectangleBorder(
                            borderRadius: BorderRadius.circular(24)),
                      ),
                    ),
                    const SizedBox(width: 16),
                    ElevatedButton.icon(
                      onPressed: () => _openEditDialog(robot),
                      icon: const Icon(Icons.edit, size: 16),
                      label: const Text('编辑', style: TextStyle(fontSize: 13)),
                      style: ElevatedButton.styleFrom(
                        padding: const EdgeInsets.symmetric(
                            vertical: 10, horizontal: 18),
                        backgroundColor: DesignTokens.brandPrimary,
                        foregroundColor: Colors.white,
                        shape: RoundedRectangleBorder(
                            borderRadius: BorderRadius.circular(24)),
                      ),
                    ),
                  ],
                ),
              ],
            ),
          ),
          // 右上角启用/禁用按钮
          Positioned(
            top: 12,
            right: 12,
            child: ElevatedButton(
              onPressed: () {
                // TODO: 启用/禁用逻辑
              },
              style: ElevatedButton.styleFrom(
                backgroundColor: robot.isEnabled
                    ? Colors.red.shade400
                    : Colors.green.shade400,
                foregroundColor: Colors.white,
                padding:
                    const EdgeInsets.symmetric(horizontal: 14, vertical: 6),
                shape: RoundedRectangleBorder(
                    borderRadius: BorderRadius.circular(18)),
                elevation: 0,
              ),
              child: Text(robot.isEnabled ? '禁用' : '启用',
                  style: const TextStyle(fontSize: 12)),
            ),
          ),
        ],
      ),
    );
  }

  // 修复：将 _buildInfoRow 定义为 _DingTalkManagementPageState 的独立方法
  Widget _buildInfoRow(IconData icon, String label, String value) {
    return Padding(
      padding: const EdgeInsets.only(bottom: 4.0),
      child: Row(
        children: [
          Icon(icon, size: 12, color: Colors.grey.shade600),
          const SizedBox(width: 4),
          Text('$label: ',
              style: TextStyle(
                  fontSize: 10,
                  color: Colors.grey.shade600,
                  fontWeight: FontWeight.w500)),
          Expanded(
            child: Text(value,
                style:
                    const TextStyle(fontSize: 10, fontWeight: FontWeight.w600),
                maxLines: 1,
                overflow: TextOverflow.ellipsis),
          ),
        ],
      ),
    );
  }

  void _testPush(dynamic robot) async {
    final robotId = robot.id;
    showDialog(
      context: context,
      barrierDismissible: false,
      builder: (ctx) => const Center(child: CircularProgressIndicator()),
    );
    try {
      // 用统一API工具调用测试接口，和列表加载一样
      final api = DingTalkApi(ApiClient());
      final result = await api.testRobot(robot);
      Navigator.of(context).pop();
      if (result.error == null) {
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(content: Text('测试成功')),
        );
      } else {
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(
              content: Text('测试失败: ${result.error}'),
              backgroundColor: Colors.red),
        );
      }
    } catch (e) {
      Navigator.of(context).pop();
      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(content: Text('测试失败: $e'), backgroundColor: Colors.red),
      );
    }
  }

  void _openEditDialog(dynamic robot) {
    ScaffoldMessenger.of(context).showSnackBar(
      SnackBar(content: Text('编辑机器人: ${robot.name}')),
    );
  }
}
