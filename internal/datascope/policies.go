package datascope

// 列表 / 行级隔离常用列（须带表名或别名，避免 JOIN 歧义）。D2：集中注册，禁止在业务里散落拼字符串。

// PolicyPurchaseOrders 采购单主表 purchase_orders。
var PolicyPurchaseOrders = TablePolicy{
	StoreColumn:   "purchase_orders.store_id",
	CreatorColumn: "purchase_orders.created_by",
}

// PolicyStoreAccounts 门店记账 store_accounts（本人维度为 operator_id）。
var PolicyStoreAccounts = TablePolicy{
	StoreColumn:   "store_accounts.store_id",
	CreatorColumn: "store_accounts.operator_id",
}

// PolicyInventoriesAsI 库存列表查询使用表别名 i。
var PolicyInventoriesAsI = TablePolicy{
	StoreColumn:   "i.store_id",
	CreatorColumn: "",
}

// PolicyInventoryOrders 出入库单主表 inventory_orders。
var PolicyInventoryOrders = TablePolicy{
	StoreColumn:   "inventory_orders.store_id",
	CreatorColumn: "inventory_orders.operator_id",
}
