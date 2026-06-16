package model

import "time"

const (
	AuditStatusSuccess = "success"
	AuditStatusFail    = "fail"

	AuditActionLogin  = "login"
	AuditActionCreate = "create"
	AuditActionUpdate = "update"
	AuditActionDelete = "delete"
	AuditActionQuery  = "query"
	AuditActionOther  = "other"
)

// AuditLog 操作日志/审计日志。
type AuditLog struct {
	ID uint `json:"id" gorm:"primarykey"`

	TraceID  string `json:"trace_id" gorm:"type:varchar(64);index;comment:请求链路ID"`
	UserID   uint   `json:"user_id" gorm:"index:idx_audit_user_time,priority:1;comment:操作人ID"`
	Username string `json:"username" gorm:"type:varchar(191);comment:操作人账号快照"`
	Nickname string `json:"nickname" gorm:"type:varchar(100);comment:操作人昵称快照"`
	Phone    string `json:"phone" gorm:"type:varchar(20);comment:操作人手机号快照"`
	RoleName string `json:"role_name" gorm:"type:varchar(100);comment:角色名称快照"`
	RoleCode string `json:"role_code" gorm:"type:varchar(50);comment:角色编码快照"`

	StoreID   uint   `json:"store_id" gorm:"index:idx_audit_store_time,priority:1;comment:门店ID"`
	StoreName string `json:"store_name" gorm:"type:varchar(100);comment:门店名称快照"`

	Module     string `json:"module" gorm:"type:varchar(64);index:idx_audit_module_time,priority:1;comment:模块编码"`
	ModuleName string `json:"module_name" gorm:"type:varchar(100);comment:模块名称"`
	Action     string `json:"action" gorm:"type:varchar(64);index:idx_audit_action_time,priority:1;comment:动作编码"`
	ActionName string `json:"action_name" gorm:"type:varchar(100);comment:动作名称"`

	ResourceType string `json:"resource_type" gorm:"type:varchar(64);comment:资源类型"`
	ResourceID   string `json:"resource_id" gorm:"type:varchar(64);index:idx_audit_resource,priority:2;comment:资源ID"`
	ResourceNo   string `json:"resource_no" gorm:"type:varchar(100);index:idx_audit_resource,priority:3;comment:业务编号"`
	ResourceName string `json:"resource_name" gorm:"type:varchar(191);comment:资源名称"`

	Method      string `json:"method" gorm:"type:varchar(16);comment:HTTP方法"`
	Path        string `json:"path" gorm:"type:varchar(255);comment:接口路径"`
	Query       string `json:"query" gorm:"type:text;comment:查询参数"`
	RequestBody string `json:"request_body" gorm:"type:longtext;comment:请求体摘要"`
	BeforeData  string `json:"before_data" gorm:"type:longtext;comment:修改前JSON"`
	AfterData   string `json:"after_data" gorm:"type:longtext;comment:修改后JSON"`
	DiffData    string `json:"diff_data" gorm:"type:longtext;comment:字段差异JSON"`

	Status       string `json:"status" gorm:"type:varchar(16);index;comment:success/fail"`
	StatusCode   int    `json:"status_code" gorm:"comment:HTTP状态码"`
	ErrorMessage string `json:"error_message" gorm:"type:text;comment:错误信息"`
	LatencyMs    int64  `json:"latency_ms" gorm:"comment:请求耗时毫秒"`

	ClientIP     string `json:"client_ip" gorm:"type:varchar(64);comment:客户端IP"`
	ClientSource string `json:"client_source" gorm:"type:varchar(64);comment:客户端来源"`
	DeviceType   string `json:"device_type" gorm:"type:varchar(32);comment:设备类型 desktop/mobile/tablet/bot/unknown"`
	OS           string `json:"os" gorm:"type:varchar(64);comment:操作系统"`
	Browser      string `json:"browser" gorm:"type:varchar(64);comment:浏览器/客户端"`
	UserAgent    string `json:"user_agent" gorm:"type:varchar(512);comment:User-Agent"`

	CreatedAt time.Time `json:"created_at" gorm:"index:idx_audit_created_at;index:idx_audit_user_time,priority:2;index:idx_audit_store_time,priority:2;index:idx_audit_module_time,priority:2;index:idx_audit_action_time,priority:2"`
}

func (AuditLog) TableName() string {
	return "audit_logs"
}

type AuditLogListReq struct {
	Page      int    `json:"page" form:"page"`
	PageSize  int    `json:"page_size" form:"page_size"`
	StartTime string `json:"start_time" form:"start_time"`
	EndTime   string `json:"end_time" form:"end_time"`
	UserID    uint   `json:"user_id" form:"user_id"`
	StoreID   uint   `json:"store_id" form:"store_id"`
	Module    string `json:"module" form:"module"`
	Action    string `json:"action" form:"action"`
	Status    string `json:"status" form:"status"`
	Keyword   string `json:"keyword" form:"keyword"`
}
