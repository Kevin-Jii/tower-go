package performance

import (
	"fmt"
	"sync"
	"testing"
)

// TestTypeConverter_ToUint tests type conversion to uint
func TestTypeConverter_ToUint(t *testing.T) {
	tc := NewTypeConverter()

	tests := []struct {
		name    string
		input   interface{}
		want    uint
		wantErr bool
	}{
		{"uint", uint(42), 42, false},
		{"uint8", uint8(42), 42, false},
		{"uint16", uint16(42), 42, false},
		{"uint32", uint32(42), 42, false},
		{"uint64", uint64(42), 42, false},
		{"int", int(42), 42, false},
		{"int8", int8(42), 42, false},
		{"int16", int16(42), 42, false},
		{"int32", int32(42), 42, false},
		{"int64", int64(42), 42, false},
		{"float32", float32(42.5), 42, false},
		{"float64", float64(42.5), 42, false},
		{"negative int", int(-1), 0, true},
		{"negative float", float64(-1.5), 0, true},
		{"string", "42", 0, true},
		{"unsupported", struct{}{}, 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tc.ToUint(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ToUint() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ToUint() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestTypeConverter_ToString tests type conversion to string
func TestTypeConverter_ToString(t *testing.T) {
	tc := NewTypeConverter()

	tests := []struct {
		name  string
		input interface{}
		want  string
	}{
		{"string", "hello", "hello"},
		{"bytes", []byte("hello"), "hello"},
		{"int", int(42), "42"},
		{"uint", uint(42), "42"},
		{"float64", float64(42.5), "42.500000"},
		{"bool true", true, "true"},
		{"bool false", false, "false"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tc.ToString(tt.input)
			if got != tt.want {
				t.Errorf("ToString() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestConcurrentCache tests concurrent cache operations
func TestConcurrentCache(t *testing.T) {
	cache := NewConcurrentCache()

	// Test basic operations
	cache.Set("key1", "value1")
	val, ok := cache.Get("key1")
	if !ok || val != "value1" {
		t.Errorf("Get() = %v, %v, want value1, true", val, ok)
	}

	// Test GetOrSet
	actual, loaded := cache.GetOrSet("key2", "value2")
	if loaded || actual != "value2" {
		t.Errorf("GetOrSet() = %v, %v, want value2, false", actual, loaded)
	}

	actual, loaded = cache.GetOrSet("key2", "value3")
	if !loaded || actual != "value2" {
		t.Errorf("GetOrSet() = %v, %v, want value2, true", actual, loaded)
	}

	// Test Delete
	cache.Delete("key1")
	_, ok = cache.Get("key1")
	if ok {
		t.Error("Delete() failed, key still exists")
	}

	// Test Range
	cache.Set("a", 1)
	cache.Set("b", 2)
	cache.Set("c", 3)
	count := 0
	cache.Range(func(key, value interface{}) bool {
		count++
		return true
	})
	if count < 3 {
		t.Errorf("Range() counted %d items, want at least 3", count)
	}

	// Test Clear
	cache.Clear()
	count = 0
	cache.Range(func(key, value interface{}) bool {
		count++
		return true
	})
	if count != 0 {
		t.Errorf("Clear() failed, %d items remain", count)
	}
}

// TestConcurrentCache_Concurrent tests concurrent access to cache
func TestConcurrentCache_Concurrent(t *testing.T) {
	cache := NewConcurrentCache()
	var wg sync.WaitGroup
	numGoroutines := 100
	numOperations := 1000

	// Concurrent writes
	wg.Add(numGoroutines)
	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			defer wg.Done()
			for j := 0; j < numOperations; j++ {
				key := fmt.Sprintf("key_%d_%d", id, j)
				cache.Set(key, j)
			}
		}(i)
	}
	wg.Wait()

	// Concurrent reads
	wg.Add(numGoroutines)
	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			defer wg.Done()
			for j := 0; j < numOperations; j++ {
				key := fmt.Sprintf("key_%d_%d", id, j)
				cache.Get(key)
			}
		}(i)
	}
	wg.Wait()

	// Concurrent mixed operations
	wg.Add(numGoroutines)
	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			defer wg.Done()
			for j := 0; j < numOperations; j++ {
				key := fmt.Sprintf("key_%d_%d", id, j)
				if j%2 == 0 {
					cache.Set(key, j)
				} else {
					cache.Get(key)
				}
			}
		}(i)
	}
	wg.Wait()
}

// TestRegexCache tests regex caching
func TestRegexCache(t *testing.T) {
	cache := NewRegexCache()

	// Test Get
	pattern := `^\d{3}-\d{4}$`
	regex1, err := cache.Get(pattern)
	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}

	// Test that second Get returns cached version
	regex2, err := cache.Get(pattern)
	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}

	// They should be the same instance (cached)
	if regex1 != regex2 {
		t.Error("Get() did not return cached regex")
	}

	// Test MustGet
	regex3 := cache.MustGet(pattern)
	if regex3 != regex1 {
		t.Error("MustGet() did not return cached regex")
	}

	// Test invalid pattern
	_, err = cache.Get(`[invalid`)
	if err == nil {
		t.Error("Get() should return error for invalid pattern")
	}

	// Test Precompile
	patterns := []string{
		`^\d+$`,
		`^[a-z]+$`,
		`^[A-Z]+$`,
	}
	err = cache.Precompile(patterns)
	if err != nil {
		t.Fatalf("Precompile() error = %v", err)
	}

	// Verify all patterns are cached
	for _, p := range patterns {
		regex, err := cache.Get(p)
		if err != nil {
			t.Errorf("Precompile() failed to cache pattern %s", p)
		}
		if !regex.MatchString("123") && p == `^\d+$` {
			t.Error("Cached regex does not work correctly")
		}
	}
}

// BenchmarkTypeSwitch benchmarks type switch vs reflection
func BenchmarkTypeSwitch(b *testing.B) {
	tc := NewTypeConverter()
	values := []interface{}{
		uint(42),
		int(42),
		int64(42),
		float64(42.5),
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		val := values[i%len(values)]
		_, _ = tc.ToUint(val)
	}
}

// BenchmarkConcurrentCache_Read benchmarks concurrent cache reads
func BenchmarkConcurrentCache_Read(b *testing.B) {
	cache := NewConcurrentCache()
	for i := 0; i < 1000; i++ {
		cache.Set(fmt.Sprintf("key%d", i), i)
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			key := fmt.Sprintf("key%d", i%1000)
			cache.Get(key)
			i++
		}
	})
}

// BenchmarkConcurrentCache_Write benchmarks concurrent cache writes
func BenchmarkConcurrentCache_Write(b *testing.B) {
	cache := NewConcurrentCache()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			key := fmt.Sprintf("key%d", i)
			cache.Set(key, i)
			i++
		}
	})
}

// BenchmarkRegexCache benchmarks regex caching
func BenchmarkRegexCache(b *testing.B) {
	cache := NewRegexCache()
	pattern := `^\d{3}-\d{4}$`
	testString := "123-4567"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		regex, _ := cache.Get(pattern)
		regex.MatchString(testString)
	}
}
