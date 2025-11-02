/// 菜单类型常量，避免魔法数字
class MenuType {
  static const int directory = 1; // 目录，仅用于分组，不直接作为页面
  static const int page = 2; // 页面，可路由和显示内容
  static const int button = 3; // 按钮/权限点，不参与侧边栏 & 路由
}
