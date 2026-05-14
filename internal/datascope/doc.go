// Package datascope 提供行级数据范围 GORM Scope 与 TablePolicy 注册表（D2/D3）。
//
// 约定（D3）：
//   - 列表类数据范围优先通过本包 listDataScopeScope + Policy* 组合，避免 module 内手写 SQL 片段。
//   - module 仍通过 pkg/datascope.Apply* 入口调用，便于渐进迁移；后续可抽到 repository 层唯一组装 Scopes。
//
// 与 RBAC 分离：接口「能否访问」由 middleware.Permission 等负责；本包只处理「能看哪些行」。

package datascope
