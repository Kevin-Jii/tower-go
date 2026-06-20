package model

import (
	"time"

	"gorm.io/gorm"
)

const StoreExpenseCategoryDictCode = "EXPENDITURECLASS"

// StoreExpense 门店支出记录。
type StoreExpense struct {
	ID           uint           `json:"id" gorm:"primaryKey;autoIncrement"`
	ExpenseNo    string         `json:"expense_no" gorm:"type:varchar(50);uniqueIndex;not null;comment:支出单号"`
	StoreID      uint           `json:"store_id" gorm:"not null;index;comment:门店ID"`
	ExpenseDate  time.Time      `json:"expense_date" gorm:"type:date;index;comment:支出日期"`
	CategoryCode string         `json:"category_code" gorm:"type:varchar(100);not null;index;comment:支出分类编码(字典:EXPENDITURECLASS)"`
	CategoryName string         `json:"category_name" gorm:"type:varchar(100);not null;comment:支出分类名称"`
	Amount       float64        `json:"amount" gorm:"type:decimal(10,2);not null;default:0;comment:支出金额"`
	Remark       string         `json:"remark" gorm:"type:varchar(500);comment:备注说明"`
	OperatorID   uint           `json:"operator_id" gorm:"not null;comment:操作人ID"`
	OperatorName string         `json:"operator_name" gorm:"type:varchar(100);comment:操作人名称"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `json:"-" gorm:"index"`
	Store        *Store         `json:"store,omitempty" gorm:"foreignKey:StoreID"`
	Operator     *User          `json:"operator,omitempty" gorm:"foreignKey:OperatorID"`
}

func (StoreExpense) TableName() string {
	return "store_expenses"
}

type CreateStoreExpenseReq struct {
	StoreID      uint    `json:"store_id"`
	CategoryCode string  `json:"category_code" binding:"required,max=100"`
	Amount       float64 `json:"amount" binding:"required,gt=0"`
	Remark       string  `json:"remark" binding:"max=500"`
}

type UpdateStoreExpenseReq struct {
	CategoryCode string   `json:"category_code" binding:"omitempty,max=100"`
	Amount       *float64 `json:"amount" binding:"omitempty,gt=0"`
	Remark       string   `json:"remark" binding:"max=500"`
}

type ListStoreExpenseReq struct {
	StoreID      uint   `form:"store_id"`
	CategoryCode string `form:"category_code"`
	Keyword      string `form:"keyword"`
	StartDate    string `form:"start_date"`
	EndDate      string `form:"end_date"`
	Page         int    `form:"page,default=1" binding:"min=1"`
	PageSize     int    `form:"page_size,default=20" binding:"min=1,max=100"`
}

type StoreExpenseStats struct {
	TotalAmount float64 `json:"total_amount"`
	Count       int64   `json:"count"`
}
