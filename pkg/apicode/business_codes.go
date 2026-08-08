package apicode

// 业务错误码集中维护。
// 新增错误码时按 HTTP 状态码 + 两位序号分配，禁止在 Controller/Service 中直接写裸数字。
var (
	// 认证与权限 401xx / 403xx
	InvalidCredentials = Code{40104, "账号或密码错误"}
	AccountDisabled    = Code{40310, "账号已被禁用"}
	StoreRequired      = Code{40308, "当前操作需要有效门店"}
	OperationDenied    = Code{40309, "当前账号无权执行此操作"}

	// 请求与校验 400xx / 422xx
	InvalidParameter       = Code{40003, "请求参数无效"}
	MissingParameter       = Code{40004, "缺少必要参数"}
	InvalidDate            = Code{40005, "日期格式无效"}
	ValidationFailed       = Code{42201, "数据校验失败"}
	ExpenseCategoryInvalid = Code{42202, "支出分类无效"}
	ReturnDateInvalid      = Code{42203, "返厂日期格式无效"}
	PaymentStatusInvalid   = Code{42204, "支付状态无效"}

	// 资源 404xx
	UserNotFound              = Code{40410, "用户不存在"}
	StoreNotFound             = Code{40411, "门店不存在"}
	SupplierNotFound          = Code{40412, "供应商不存在"}
	ProductNotFound           = Code{40413, "商品不存在"}
	InventoryNotFound         = Code{40414, "库存记录不存在"}
	DictTypeNotFound          = Code{40415, "字典类型不存在"}
	DictDataNotFound          = Code{40416, "字典数据不存在"}
	MemberNotFound            = Code{40424, "会员不存在"}
	OrderNotFound             = Code{40417, "单据不存在"}
	CustomerNotFound          = Code{40418, "客户不存在"}
	UnitSpecNotFound          = Code{40419, "商品规格不存在"}
	PriceListNotFound         = Code{40420, "价目单不存在"}
	CategoryNotFound          = Code{40421, "分类不存在"}
	PrinterNotFound           = Code{40422, "打印机不存在"}
	ItemNotFound              = Code{40423, "明细不存在"}
	ExpenseNotFound           = Code{40425, "支出记录不存在"}
	ReturnNotFound            = Code{40426, "返厂记录不存在"}
	WineStorageNotFound       = Code{40427, "存酒记录不存在"}
	ThirdPartyAccountNotFound = Code{40428, "第三方账号不存在"}
	WechatNotBound            = Code{40429, "微信未绑定账号"}

	// 冲突与业务条件 409xx
	Conflict                  = Code{40901, "数据冲突"}
	DictTypeAlreadyExists     = Code{40910, "字典类型编码已存在"}
	SupplierAlreadyBound      = Code{40911, "供应商已绑定"}
	InventoryInsufficient     = Code{40912, "库存不足"}
	OrderStateConflict        = Code{40913, "单据状态不允许当前操作"}
	MemberPhoneExists         = Code{40914, "手机号已注册"}
	ResourceInUse             = Code{40915, "资源正在使用中"}
	DuplicateOperation        = Code{40916, "重复操作"}
	UnitSpecMismatch          = Code{40917, "商品规格与商品不匹配"}
	ProductNotBound           = Code{40918, "商品未绑定到当前门店"}
	OrderDeletionDenied       = Code{40919, "当前订单状态不允许删除"}
	CustomerDisabled          = Code{40920, "客户已停用"}
	OptimisticLockConflict    = Code{40921, "数据已被其他操作修改"}
	BalanceInsufficient       = Code{40922, "余额不足"}
	WineInsufficient          = Code{40923, "存酒数量不足"}
	WechatAlreadyBound        = Code{40924, "微信已绑定其他账号"}
	PhoneAlreadyExists        = Code{40925, "手机号已注册"}
	SupplierAlreadyUsed       = Code{40926, "供应商正在使用中"}
	StoreAccountOrderNoExists = Code{40927, "当前渠道下的外卖订单号已存在"}

	// 服务与外部依赖 500xx / 502xx
	ConfigMissing = Code{50002, "服务配置缺失"}
	// 外部服务 502xx
	ExternalServiceFailed = Code{50201, "外部服务调用失败"}
)
