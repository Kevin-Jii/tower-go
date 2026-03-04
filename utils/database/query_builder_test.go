package database

import (
	"testing"
)

// TestEncodeCursor 测试游标编码
func TestEncodeCursor(t *testing.T) {
	paginator := NewCursorPaginator("id", false)

	// 测试编码
	cursor := paginator.EncodeCursor(123, nil)
	if cursor == "" {
		t.Error("Expected non-empty cursor")
	}

	// 测试解码
	lastID, lastValue, err := paginator.DecodeCursor(cursor)
	if err != nil {
		t.Errorf("DecodeCursor failed: %v", err)
	}

	if lastID != 123 {
		t.Errorf("Expected lastID=123, got %d", lastID)
	}

	if lastValue != nil {
		t.Errorf("Expected lastValue=nil, got %v", lastValue)
	}
}

// TestEncodeCursorWithValue 测试带值的游标编码
func TestEncodeCursorWithValue(t *testing.T) {
	paginator := NewCursorPaginator("created_at", true)

	// 测试编码
	cursor := paginator.EncodeCursor(456, "2024-01-01")
	if cursor == "" {
		t.Error("Expected non-empty cursor")
	}

	// 测试解码
	lastID, lastValue, err := paginator.DecodeCursor(cursor)
	if err != nil {
		t.Errorf("DecodeCursor failed: %v", err)
	}

	if lastID != 456 {
		t.Errorf("Expected lastID=456, got %d", lastID)
	}

	if lastValue == nil {
		t.Error("Expected non-nil lastValue")
	}
}

// TestDecodeCursorEmpty 测试空游标解码
func TestDecodeCursorEmpty(t *testing.T) {
	paginator := NewCursorPaginator("id", false)

	lastID, lastValue, err := paginator.DecodeCursor("")
	if err != nil {
		t.Errorf("DecodeCursor with empty cursor should not error: %v", err)
	}

	if lastID != 0 {
		t.Errorf("Expected lastID=0, got %d", lastID)
	}

	if lastValue != nil {
		t.Errorf("Expected lastValue=nil, got %v", lastValue)
	}
}

// TestDecodeCursorInvalid 测试无效游标解码
func TestDecodeCursorInvalid(t *testing.T) {
	paginator := NewCursorPaginator("id", false)

	_, _, err := paginator.DecodeCursor("invalid-cursor")
	if err == nil {
		t.Error("Expected error for invalid cursor")
	}
}

// TestEncodeCursorZeroID 测试零ID编码
func TestEncodeCursorZeroID(t *testing.T) {
	paginator := NewCursorPaginator("id", false)

	cursor := paginator.EncodeCursor(0, nil)
	if cursor != "" {
		t.Error("Expected empty cursor for zero ID")
	}
}
