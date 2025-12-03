package model

import "time"

// DictType 字典类型
type DictType struct {
	ID        uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	Code      string    `json:"code" gorm:"type:varchar(100);uniqueIndex;not null;comment:字典类型编码"`
	Name      string    `json:"name" gorm:"type:varchar(100);not null;comment:字典类型名称"`
	Remark    string    `json:"remark" gorm:"type:varchar(500);comment:备注"`
	Status    int8      `json:"status" gorm:"default:1;comment:状态 1=启用 0=禁用"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (DictType) TableName() string {
	return "dict_types"
}

// DictData 字典数据
type DictData struct {
	ID         uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	TypeID     uint      `json:"type_id" gorm:"not null;index;comment:字典类型ID"`
	TypeCode   string    `json:"type_code" gorm:"type:varchar(100);index;not null;comment:字典类型编码"`
	Label      string    `json:"label" gorm:"type:varchar(100);not null;comment:字典标签"`
	Value      string    `json:"value" gorm:"type:varchar(100);not null;comment:字典值"`
	Sort       int       `json:"sort" gorm:"default:0;comment:排序"`
	CssClass   string    `json:"css_class" gorm:"type:varchar(100);comment:样式类名"`
	ListClass  string    `json:"list_class" gorm:"type:varchar(100);comment:列表样式(success/info/warning/danger)"`
	IsDefault  bool      `json:"is_default" gorm:"default:false;comment:是否默认"`
	Remark     string    `json:"remark" gorm:"type:varchar(500);comment:备注"`
	Status     int8      `json:"status" gorm:"default:1;comment:状态 1=启用 0=禁用"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func (DictData) TableName() string {
	return "dict_data"
}

// CreateDictTypeReq 创建字典类型请求
type CreateDictTypeReq struct {
	Code   string `json:"code" binding:"required,max=100"`
	Name   string `json:"name" binding:"required,max=100"`
	Remark string `json:"remark" binding:"max=500"`
	Status int8   `json:"status" binding:"oneof=0 1"`
}

// UpdateDictTypeReq 更新字典类型请求
type UpdateDictTypeReq struct {
	Code   string `json:"code" binding:"max=100"`
	Name   string `json:"name" binding:"max=100"`
	Remark string `json:"remark" binding:"max=500"`
	Status *int8  `json:"status" binding:"omitempty,oneof=0 1"`
}

// CreateDictDataReq 创建字典数据请求
type CreateDictDataReq struct {
	TypeCode  string `json:"type_code" binding:"required,max=100"`
	Label     string `json:"label" binding:"required,max=100"`
	Value     string `json:"value" binding:"required,max=100"`
	Sort      int    `json:"sort"`
	CssClass  string `json:"css_class" binding:"max=100"`
	ListClass string `json:"list_class" binding:"max=100"`
	IsDefault bool   `json:"is_default"`
	Remark    string `json:"remark" binding:"max=500"`
	Status    int8   `json:"status" binding:"oneof=0 1"`
}

// UpdateDictDataReq 更新字典数据请求
type UpdateDictDataReq struct {
	Label     string `json:"label" binding:"max=100"`
	Value     string `json:"value" binding:"max=100"`
	Sort      *int   `json:"sort"`
	CssClass  string `json:"css_class" binding:"max=100"`
	ListClass string `json:"list_class" binding:"max=100"`
	IsDefault *bool  `json:"is_default"`
	Remark    string `json:"remark" binding:"max=500"`
	Status    *int8  `json:"status" binding:"omitempty,oneof=0 1"`
}
