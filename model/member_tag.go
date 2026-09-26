package model

import "time"

// MemberTag is owned by exactly one store. Store 0 is reserved for legacy/HQ
// campaigns and is not visible to store-bound users.
type MemberTag struct {
	ID          uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	StoreID     uint      `json:"store_id" gorm:"not null;index;uniqueIndex:idx_member_tags_store_name,priority:1"`
	Name        string    `json:"name" gorm:"type:varchar(80);not null;uniqueIndex:idx_member_tags_store_name,priority:2"`
	Color       string    `json:"color" gorm:"type:varchar(32)"`
	Description string    `json:"description" gorm:"type:varchar(255)"`
	MemberCount int64     `json:"member_count" gorm:"-"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (MemberTag) TableName() string { return "member_tags" }

type MemberTagBinding struct {
	ID        uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	StoreID   uint      `json:"store_id" gorm:"not null;index;uniqueIndex:idx_member_tag_bindings_unique,priority:1"`
	MemberID  uint      `json:"member_id" gorm:"not null;index;uniqueIndex:idx_member_tag_bindings_unique,priority:2"`
	TagID     uint      `json:"tag_id" gorm:"not null;index;uniqueIndex:idx_member_tag_bindings_unique,priority:3"`
	CreatedAt time.Time `json:"created_at"`
}

func (MemberTagBinding) TableName() string { return "member_tag_bindings" }

type UpsertMemberTagReq struct {
	StoreID     uint   `json:"store_id"`
	Name        string `json:"name" binding:"required,max=80"`
	Color       string `json:"color" binding:"max=32"`
	Description string `json:"description" binding:"max=255"`
}

type AssignMemberTagsReq struct {
	TagIDs UintList `json:"tag_ids"`
}

type BindMemberTagReq struct {
	MemberID uint `json:"member_id" binding:"required"`
}

// SmsMemberSummary intentionally excludes balances, points and other financial
// fields from marketing-only member/tag APIs.
type SmsMemberSummary struct {
	ID      uint   `json:"id"`
	StoreID uint   `json:"store_id"`
	UID     string `json:"uid"`
	Name    string `json:"name"`
	Phone   string `json:"phone"`
}
