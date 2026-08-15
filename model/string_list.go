package model

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

// StringList stores a string slice as JSON while keeping the API representation as an array.
type StringList []string

func (s *StringList) Scan(value interface{}) error {
	if value == nil {
		*s = StringList{}
		return nil
	}

	var data []byte
	switch typed := value.(type) {
	case []byte:
		data = typed
	case string:
		data = []byte(typed)
	default:
		return fmt.Errorf("scan StringList from %T", value)
	}
	if len(data) == 0 {
		*s = StringList{}
		return nil
	}
	if err := json.Unmarshal(data, s); err != nil {
		return fmt.Errorf("decode StringList: %w", err)
	}
	if *s == nil {
		*s = StringList{}
	}
	return nil
}

func (s StringList) Value() (driver.Value, error) {
	if s == nil {
		s = StringList{}
	}
	data, err := json.Marshal(s)
	if err != nil {
		return nil, fmt.Errorf("encode StringList: %w", err)
	}
	return string(data), nil
}
