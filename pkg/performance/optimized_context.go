package performance

import (
	"errors"
	"fmt"

	"github.com/gin-gonic/gin"
)

// ContextExtractor provides optimized context value extraction using type switches
type ContextExtractor struct {
	converter *TypeConverter
}

// NewContextExtractor creates a new ContextExtractor instance
func NewContextExtractor() *ContextExtractor {
	return &ContextExtractor{
		converter: GetTypeConverter(),
	}
}

// GetStoreID extracts StoreID from Gin Context with optimized type assertion
func (ce *ContextExtractor) GetStoreID(ctx *gin.Context) (uint, error) {
	val, exists := ctx.Get("StoreID")
	if !exists {
		return 0, errors.New("StoreID not found in context (Middleware missing?)")
	}

	// Use type switch for fast type assertion
	switch v := val.(type) {
	case uint:
		if v == 0 {
			return 0, errors.New("Invalid StoreID in context")
		}
		return v, nil
	case int:
		if v <= 0 {
			return 0, errors.New("Invalid StoreID in context")
		}
		return uint(v), nil
	case int64:
		if v <= 0 {
			return 0, errors.New("Invalid StoreID in context")
		}
		return uint(v), nil
	case float64:
		if v <= 0 {
			return 0, errors.New("Invalid StoreID in context")
		}
		return uint(v), nil
	default:
		// Fallback to converter for other types
		result, err := ce.converter.ToUint(val)
		if err != nil {
			return 0, fmt.Errorf("Invalid StoreID type in context: %T", val)
		}
		if result == 0 {
			return 0, errors.New("Invalid StoreID in context")
		}
		return result, nil
	}
}

// GetUserID extracts UserID from Gin Context with optimized type assertion
func (ce *ContextExtractor) GetUserID(ctx *gin.Context) (uint, error) {
	val, exists := ctx.Get("UserID")
	if !exists {
		return 0, errors.New("UserID not found in context")
	}

	// Use type switch for fast type assertion
	switch v := val.(type) {
	case uint:
		if v == 0 {
			return 0, errors.New("Invalid UserID in context")
		}
		return v, nil
	case int:
		if v <= 0 {
			return 0, errors.New("Invalid UserID in context")
		}
		return uint(v), nil
	case int64:
		if v <= 0 {
			return 0, errors.New("Invalid UserID in context")
		}
		return uint(v), nil
	case float64:
		if v <= 0 {
			return 0, errors.New("Invalid UserID in context")
		}
		return uint(v), nil
	default:
		result, err := ce.converter.ToUint(val)
		if err != nil {
			return 0, fmt.Errorf("Invalid UserID type in context: %T", val)
		}
		if result == 0 {
			return 0, errors.New("Invalid UserID in context")
		}
		return result, nil
	}
}

// GetString extracts a string value from Gin Context with optimized type assertion
func (ce *ContextExtractor) GetString(ctx *gin.Context, key string) (string, error) {
	val, exists := ctx.Get(key)
	if !exists {
		return "", fmt.Errorf("%s not found in context", key)
	}

	// Use type switch for fast type assertion
	switch v := val.(type) {
	case string:
		return v, nil
	case []byte:
		return string(v), nil
	default:
		return ce.converter.ToString(val), nil
	}
}

// GetInt extracts an int value from Gin Context with optimized type assertion
func (ce *ContextExtractor) GetInt(ctx *gin.Context, key string) (int, error) {
	val, exists := ctx.Get(key)
	if !exists {
		return 0, fmt.Errorf("%s not found in context", key)
	}

	// Use type switch for fast type assertion
	switch v := val.(type) {
	case int:
		return v, nil
	case int8:
		return int(v), nil
	case int16:
		return int(v), nil
	case int32:
		return int(v), nil
	case int64:
		return int(v), nil
	case uint:
		return int(v), nil
	case uint8:
		return int(v), nil
	case uint16:
		return int(v), nil
	case uint32:
		return int(v), nil
	case uint64:
		return int(v), nil
	case float32:
		return int(v), nil
	case float64:
		return int(v), nil
	default:
		return ce.converter.ToInt(val)
	}
}

// GetBool extracts a bool value from Gin Context with optimized type assertion
func (ce *ContextExtractor) GetBool(ctx *gin.Context, key string) (bool, error) {
	val, exists := ctx.Get(key)
	if !exists {
		return false, fmt.Errorf("%s not found in context", key)
	}

	// Use type switch for fast type assertion
	switch v := val.(type) {
	case bool:
		return v, nil
	case int:
		return v != 0, nil
	case int64:
		return v != 0, nil
	case string:
		return v == "true" || v == "1" || v == "yes", nil
	default:
		return false, fmt.Errorf("cannot convert %T to bool", val)
	}
}

// Global instance for convenience
var globalContextExtractor = NewContextExtractor()

// GetContextExtractor returns the global ContextExtractor instance
func GetContextExtractor() *ContextExtractor {
	return globalContextExtractor
}
