package performance

import (
	"fmt"
	"regexp"
	"sync"
)

// TypeConverter provides optimized type conversion using type switches instead of reflection
type TypeConverter struct{}

// NewTypeConverter creates a new TypeConverter instance
func NewTypeConverter() *TypeConverter {
	return &TypeConverter{}
}

// ToUint converts interface{} to uint using type switch (faster than reflection)
func (tc *TypeConverter) ToUint(val interface{}) (uint, error) {
	switch v := val.(type) {
	case uint:
		return v, nil
	case uint8:
		return uint(v), nil
	case uint16:
		return uint(v), nil
	case uint32:
		return uint(v), nil
	case uint64:
		return uint(v), nil
	case int:
		if v < 0 {
			return 0, fmt.Errorf("negative value cannot be converted to uint: %d", v)
		}
		return uint(v), nil
	case int8:
		if v < 0 {
			return 0, fmt.Errorf("negative value cannot be converted to uint: %d", v)
		}
		return uint(v), nil
	case int16:
		if v < 0 {
			return 0, fmt.Errorf("negative value cannot be converted to uint: %d", v)
		}
		return uint(v), nil
	case int32:
		if v < 0 {
			return 0, fmt.Errorf("negative value cannot be converted to uint: %d", v)
		}
		return uint(v), nil
	case int64:
		if v < 0 {
			return 0, fmt.Errorf("negative value cannot be converted to uint: %d", v)
		}
		return uint(v), nil
	case float32:
		if v < 0 {
			return 0, fmt.Errorf("negative value cannot be converted to uint: %f", v)
		}
		return uint(v), nil
	case float64:
		if v < 0 {
			return 0, fmt.Errorf("negative value cannot be converted to uint: %f", v)
		}
		return uint(v), nil
	case string:
		// For string conversion, we'd need strconv, but that's not a hot path optimization
		return 0, fmt.Errorf("string to uint conversion not supported in hot path")
	default:
		return 0, fmt.Errorf("unsupported type: %T", val)
	}
}

// ToString converts interface{} to string using type switch
func (tc *TypeConverter) ToString(val interface{}) string {
	switch v := val.(type) {
	case string:
		return v
	case []byte:
		return string(v)
	case int:
		return fmt.Sprintf("%d", v)
	case int8:
		return fmt.Sprintf("%d", v)
	case int16:
		return fmt.Sprintf("%d", v)
	case int32:
		return fmt.Sprintf("%d", v)
	case int64:
		return fmt.Sprintf("%d", v)
	case uint:
		return fmt.Sprintf("%d", v)
	case uint8:
		return fmt.Sprintf("%d", v)
	case uint16:
		return fmt.Sprintf("%d", v)
	case uint32:
		return fmt.Sprintf("%d", v)
	case uint64:
		return fmt.Sprintf("%d", v)
	case float32:
		return fmt.Sprintf("%f", v)
	case float64:
		return fmt.Sprintf("%f", v)
	case bool:
		if v {
			return "true"
		}
		return "false"
	default:
		return fmt.Sprintf("%v", v)
	}
}

// ToInt converts interface{} to int using type switch
func (tc *TypeConverter) ToInt(val interface{}) (int, error) {
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
		return 0, fmt.Errorf("unsupported type: %T", val)
	}
}

// ConcurrentCache provides a thread-safe cache using sync.Map for read-heavy workloads
// sync.Map is optimized for scenarios where entries are written once and read many times
type ConcurrentCache struct {
	data sync.Map
}

// NewConcurrentCache creates a new ConcurrentCache instance
func NewConcurrentCache() *ConcurrentCache {
	return &ConcurrentCache{}
}

// Get retrieves a value from the cache
func (cc *ConcurrentCache) Get(key string) (interface{}, bool) {
	return cc.data.Load(key)
}

// Set stores a value in the cache
func (cc *ConcurrentCache) Set(key string, value interface{}) {
	cc.data.Store(key, value)
}

// Delete removes a value from the cache
func (cc *ConcurrentCache) Delete(key string) {
	cc.data.Delete(key)
}

// GetOrSet atomically gets or sets a value
func (cc *ConcurrentCache) GetOrSet(key string, value interface{}) (actual interface{}, loaded bool) {
	return cc.data.LoadOrStore(key, value)
}

// Range iterates over all entries in the cache
func (cc *ConcurrentCache) Range(f func(key, value interface{}) bool) {
	cc.data.Range(f)
}

// Clear removes all entries from the cache
func (cc *ConcurrentCache) Clear() {
	cc.data.Range(func(key, value interface{}) bool {
		cc.data.Delete(key)
		return true
	})
}

// RegexCache provides a cache for pre-compiled regular expressions
type RegexCache struct {
	cache sync.Map
}

// NewRegexCache creates a new RegexCache instance
func NewRegexCache() *RegexCache {
	return &RegexCache{}
}

// Get retrieves a compiled regex from cache or compiles and caches it
func (rc *RegexCache) Get(pattern string) (*regexp.Regexp, error) {
	// Try to load from cache
	if cached, ok := rc.cache.Load(pattern); ok {
		return cached.(*regexp.Regexp), nil
	}

	// Compile the regex
	compiled, err := regexp.Compile(pattern)
	if err != nil {
		return nil, err
	}

	// Store in cache
	rc.cache.Store(pattern, compiled)
	return compiled, nil
}

// MustGet retrieves a compiled regex from cache or panics if compilation fails
func (rc *RegexCache) MustGet(pattern string) *regexp.Regexp {
	regex, err := rc.Get(pattern)
	if err != nil {
		panic(fmt.Sprintf("failed to compile regex pattern %s: %v", pattern, err))
	}
	return regex
}

// Precompile pre-compiles a set of regex patterns
func (rc *RegexCache) Precompile(patterns []string) error {
	for _, pattern := range patterns {
		if _, err := rc.Get(pattern); err != nil {
			return fmt.Errorf("failed to precompile pattern %s: %w", pattern, err)
		}
	}
	return nil
}

// Global instances for convenience
var (
	globalTypeConverter = NewTypeConverter()
	globalRegexCache    = NewRegexCache()
)

// GetTypeConverter returns the global TypeConverter instance
func GetTypeConverter() *TypeConverter {
	return globalTypeConverter
}

// GetRegexCache returns the global RegexCache instance
func GetRegexCache() *RegexCache {
	return globalRegexCache
}
