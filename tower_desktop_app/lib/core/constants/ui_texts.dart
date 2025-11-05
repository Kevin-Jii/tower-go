/// UI 文案常量集中管理，统一后续做国际化或批量修改。
/// 命名约定：模块_用途 或通用语义。
class UITexts {
  // 登录页
  static const loginWelcome = '欢迎回来';
  static const loginSubtitle = '请登录您的账户';
  static const loginRemember = '记住账号密码';
  static const loginForgot = '忘记密码?';
  static const loginForgotHint = '请联系管理员重置密码';
  static const loginButton = '登录';
  static const loginVersion = 'Tower 餐饮管理系统 v1.0';

  // 角色管理
  static const roleTitle = '角色管理';
  static const roleAddButton = '新增角色';
  static const roleEditButton = '修改';
  static const roleDeleteButton = '删除';
  static const roleDeleteConfirmTitle = '确认删除';
  static const roleDeleteConfirmContentPrefix = '确定删除角色';
  static const roleDeleteConfirmCancel = '取消';
  static const roleDeleteConfirmOk = '删除';
  static const roleCreateSuccess = '创建成功';
  static const roleCreateFailed = '创建失败';
  static const roleUpdateSuccess = '更新成功';
  static const roleUpdateFailed = '更新失败';
  static const roleDeleteSuccess = '删除成功';
  static const roleDeleteFailed = '删除失败';
  static const roleSearchPlaceholder = '关键字';
  static const roleSearchButton = '搜索';
  static const roleRetryButton = '重试';
  static const roleEmpty = '暂无角色';
  static const roleColumnId = 'ID';
  static const roleColumnName = '角色名称';
  static const roleColumnCode = '角色编码';
  static const roleColumnDesc = '备注';
  static const roleColumnStatus = '状态';
  static const roleColumnCreated = '创建时间';
  static const roleColumnActions = '操作';
  static const roleStatusEnabled = '正常';
  static const roleStatusDisabled = '停用';

  // 菜品管理
  static const dishTitle = '菜品管理';
  static const dishStoreSelectLabel = '选择门店';
  static const dishCategoryTitle = '菜品分类';
  static const dishCategoryAdd = '新增分类';
  static const dishCategoryEdit = '编辑分类';
  static const dishCategoryDelete = '删除分类';
  static const dishCategoryDeleteConfirm = '确认删除该分类？';
  static const dishCategoryFormName = '分类名称';
  static const dishCategoryEmpty = '暂无分类';
  static const dishCategoryLoading = '分类加载中...';
  static const dishCategoryError = '分类加载失败';
  static const dishAdd = '新增菜品';
  static const dishEdit = '编辑菜品';
  static const dishDelete = '删除菜品';
  static const dishDeleteConfirm = '确认删除该菜品？';
  static const dishDeleteConfirmTitle = '确认删除';
  static const dishDeleteConfirmMessage = '确定要删除菜品';
  static const dishEmpty = '该分类下暂无菜品';
  static const dishLoading = '菜品加载中...';
  static const dishError = '菜品加载失败';
  static const dishColumnName = '菜品名称';
  static const dishColumnPrice = '价格';
  static const dishColumnStatus = '状态';
  static const dishColumnActions = '操作';
  static const dishStatusEnabled = '上架';
  static const dishStatusDisabled = '下架';
  static const dishFormName = '名称';
  static const dishFormDesc = '描述(可选)';
  static const dishFormPrice = '价格(元)';
  static const dishFormStatus = '上架状态';
  static const dishCreateSuccess = '菜品创建成功';
  static const dishCreateFailed = '菜品创建失败';
  static const dishUpdateSuccess = '菜品更新成功';
  static const dishUpdateFailed = '菜品更新失败';
  static const dishDeleteSuccess = '菜品删除成功';
  static const dishDeleteFailed = '菜品删除失败';

  // 通用
  static const commonCancel = '取消';
  static const commonOk = '确定';
}
