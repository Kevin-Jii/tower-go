package model

import "time"

// PrinterType 打印机类型
type PrinterType int

const (
	PrinterTypeReceipt PrinterType = 1 // 小票打印机
	PrinterTypeLabel   PrinterType = 2 // 标签打印机
)

// Printer 打印机表
type Printer struct {
	ID            uint        `json:"id" gorm:"primarykey"`
	StoreID       uint        `json:"store_id" gorm:"index;not null"`                  // 关联门店ID
	Sn            string      `json:"sn" gorm:"uniqueIndex;type:varchar(32)"`         // 打印机SN号
	Name          string      `json:"name" gorm:"type:varchar(100)"`                  // 打印机名称
	Type          PrinterType `json:"type" gorm:"default:1"`                          // 打印机类型：1=小票，2=标签
	Status        int         `json:"status" gorm:"default:1"`                        // 状态：1=正常，2=停用
	IsDefault     int         `json:"is_default" gorm:"default:0"`                    // 是否为默认打印机：0=否，1=是
	Online        int         `json:"online" gorm:"default:0"`                        // 在线状态：0=离线，1=在线，2=异常
	LastHeartbeat *time.Time  `json:"last_heartbeat,omitempty" gorm:"type:datetime"`  // 最后心跳时间（可为空）
	Remark        string      `json:"remark" gorm:"type:text"`                        // 备注
	CreatedAt     time.Time   `json:"created_at"`
	UpdatedAt     time.Time   `json:"updated_at"`
}

// TableName 指定表名
func (Printer) TableName() string {
	return "printers"
}

// CreatePrinterReq 创建打印机请求
type CreatePrinterReq struct {
	StoreID   uint        `json:"store_id" binding:"required"`
	Sn        string      `json:"sn" binding:"required"`
	Name      string      `json:"name"`
	Type      PrinterType `json:"type"`
	IsDefault int         `json:"is_default"`
	Remark    string      `json:"remark"`
}

// UpdatePrinterReq 更新打印机请求
type UpdatePrinterReq struct {
	Name      *string      `json:"name,omitempty"`
	Type      *PrinterType `json:"type,omitempty"`
	Status    *int         `json:"status,omitempty"`
	IsDefault *int         `json:"is_default,omitempty"`
	Remark    *string      `json:"remark,omitempty"`
}

// BindPrinterReq 绑定打印机到门店请求
type BindPrinterReq struct {
	StoreID   uint   `json:"store_id" binding:"required"`
	Sn        string `json:"sn" binding:"required"`
	Name      string `json:"name"`
	Type      int    `json:"type"`
	IsDefault int    `json:"is_default"`
	Remark    string `json:"remark"`
}

// PrinterResp 打印机响应（包含状态）
type PrinterResp struct {
	ID            uint       `json:"id"`
	StoreID       uint       `json:"store_id"`
	StoreName     string     `json:"store_name,omitempty"`
	Sn            string     `json:"sn"`
	Name          string     `json:"name"`
	Type          int        `json:"type"`
	TypeName      string     `json:"type_name"`
	Status        int        `json:"status"`
	StatusName    string     `json:"status_name"`
	IsDefault     int        `json:"is_default"`
	Online        int        `json:"online"`               // 在线状态：0=离线，1=在线，2=异常
	LastHeartbeat *time.Time `json:"last_heartbeat,omitempty"` // 最后心跳时间
	Remark        string     `json:"remark"`
	CreatedAt     time.Time  `json:"created_at"`
}

// PrinterStatus 打印机在线状态
type PrinterStatus struct {
	Sn     string `json:"sn"`
	Online int    `json:"online"` // 0-离线 1-在线正常 2-在线异常
}