/// 全局常量集中维护，避免散落在各业务文件中。
/// 按领域分组：应用文本、权限码、状态值、分页、环境等。

class AppTexts {
  static const String appName = 'Tower 管理桌面端';
  static const String selectMenuPlaceholder = '请选择左侧菜单';
  static const String defaultStoreName = '默认总部门店';
  static const String envProd = 'PROD';
}

class PermissionCodes {
  // 用户相关
  static const String userAdd = 'system:user:add';
  static const String userDelete = 'system:user:delete';
  static const String userEdit = 'system:user:edit';
  static const String userList = 'system:user:list';

  // 角色相关
  static const String roleList = 'system:role:list';
  static const String roleAdd = 'system:role:add';
  static const String roleEdit = 'system:role:edit';
  static const String roleDelete = 'system:role:delete';
  static const String roleMenu = 'system:role:menu'; // 分配菜单

  // 菜单相关
  static const String menuList = 'system:menu:list';
  static const String menuAdd = 'system:menu:add';
  static const String menuEdit = 'system:menu:edit';
  static const String menuDelete = 'system:menu:delete';

  // 门店相关（注意：后端给的是 store:* 而非 system:store:*）
  static const String storeList = 'store:list';
  static const String storeAdd = 'store:add';
  static const String storeEdit = 'store:edit';
  static const String storeDelete = 'store:delete';
  static const String storeMenu = 'store:menu'; // 配置权限

  // 菜品相关（dish）
  static const String dishList = 'dish:list';
  static const String dishAdd = 'dish:add';
  static const String dishEdit = 'dish:edit';
  static const String dishDelete = 'dish:delete';
  static const String dishStatus = 'dish:status'; // 上下架

  // 报菜 / 报表相关（report）
  static const String reportList = 'report:list';
  static const String reportAdd = 'report:add';
  static const String reportEdit = 'report:edit';
  static const String reportDelete = 'report:delete';
  static const String reportStatistics = 'report:statistics';

  // 用于后续快速判断是否为按钮级(操作级)权限（Type=3）时可维护一个集合
  // static const Set<String> actionPermissions = {
  //   userAdd, userEdit, userDelete, roleAdd, roleEdit, roleDelete, roleMenu,
  //   menuAdd, menuEdit, menuDelete, storeAdd, storeEdit, storeDelete, storeMenu,
  //   dishAdd, dishEdit, dishDelete, dishStatus, reportAdd, reportEdit, reportDelete,
  // };
}

class StatusValues {
  static const int enabled = 1;
  static const int disabled = 0;
}

class PaginationDefaults {
  static const int pageSize = 10;
  static const int firstPage = 1;
}

class EnvironmentTags {
  static const String prod = 'PROD';
  static const String dev = 'DEV';
  static const String test = 'TEST';
}

/// 统一地方管理路由名称（如有需要）
class RouteNames {
  static const String login = '/login';
  static const String home = '/home';
  static const String dishes = '/dishes';
  // static const String userList = '/user/list';
}

/// 统一接口路径常量（根据需要逐步迁移）
class ApiPaths {
  static const String users = '/users';
  static const String stores = '/stores';
  static const String menus = '/menus';
  static const String dishes = '/dishes';
  static const String roles = '/roles';
  // 如后端有统计/报表模块可继续补充：static const String reports = '/reports';
}
