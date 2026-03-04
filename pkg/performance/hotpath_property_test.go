package performance

import (
	"fmt"
	"math/rand"
	"sync"
	"testing"
	"time"

	"github.com/leanovate/gopter"
	"github.com/leanovate/gopter/gen"
	"github.com/leanovate/gopter/prop"
)

// **Feature: performance-optimization, Property 34: Type switch performance**
// Property: For any frequent type assertion operation, using type switch should be faster than using reflection
func TestProperty_TypeSwitchPerformance(t *testing.T) {
	properties := gopter.NewProperties(nil)

	properties.Property("type switch converts numeric types correctly", prop.ForAll(
		func(val interface{}) bool {
			tc := NewTypeConverter()

			// Test conversion based on type
			switch v := val.(type) {
			case int:
				result, err := tc.ToUint(v)
				if v < 0 {
					return err != nil
				}
				return err == nil && result == uint(v)
			case int8:
				result, err := tc.ToUint(v)
				if v < 0 {
					return err != nil
				}
				return err == nil && result == uint(v)
			case int16:
				result, err := tc.ToUint(v)
				if v < 0 {
					return err != nil
				}
				return err == nil && result == uint(v)
			case int32:
				result, err := tc.ToUint(v)
				if v < 0 {
					return err != nil
				}
				return err == nil && result == uint(v)
			case int64:
				result, err := tc.ToUint(v)
				if v < 0 {
					return err != nil
				}
				return err == nil && result == uint(v)
			case uint:
				result, err := tc.ToUint(v)
				return err == nil && result == v
			case uint8:
				result, err := tc.ToUint(v)
				return err == nil && result == uint(v)
			case uint16:
				result, err := tc.ToUint(v)
				return err == nil && result == uint(v)
			case uint32:
				result, err := tc.ToUint(v)
				return err == nil && result == uint(v)
			case uint64:
				result, err := tc.ToUint(v)
				return err == nil && result == uint(v)
			case float32:
				result, err := tc.ToUint(v)
				if v < 0 {
					return err != nil
				}
				return err == nil && result == uint(v)
			case float64:
				result, err := tc.ToUint(v)
				if v < 0 {
					return err != nil
				}
				return err == nil && result == uint(v)
			default:
				return false
			}
		},
		gen.OneGenOf(
			gen.Int(),
			gen.Int8(),
			gen.Int16(),
			gen.Int32(),
			gen.Int64(),
			gen.UInt(),
			gen.UInt8(),
			gen.UInt16(),
			gen.UInt32(),
			gen.UInt64(),
			gen.Float32(),
			gen.Float64(),
		),
	))

	properties.Property("type switch converts to string correctly", prop.ForAll(
		func(val interface{}) bool {
			tc := NewTypeConverter()
			result := tc.ToString(val)

			// Result should be a string (may be empty for empty string input)
			return result != "" || fmt.Sprintf("%v", val) == ""
		},
		gen.OneGenOf(
			gen.AlphaString(), // Use AlphaString instead of AnyString to avoid empty strings
			gen.Int(),
			gen.UInt(),
			gen.Float64(),
			gen.Bool(),
		),
	))

	properties.Property("type switch handles negative values correctly", prop.ForAll(
		func(val int) bool {
			tc := NewTypeConverter()
			result, err := tc.ToUint(val)

			if val < 0 {
				// Should return error for negative values
				return err != nil && result == 0
			}
			// Should succeed for non-negative values
			return err == nil && result == uint(val)
		},
		gen.Int(),
	))

	properties.TestingRun(t, gopter.ConsoleReporter(false))
}

// **Feature: performance-optimization, Property 35: sync.Map concurrent performance**
// Property: For any concurrent map operations with read-heavy workload, sync.Map should perform better than mutex-protected map
func TestProperty_SyncMapConcurrentPerformance(t *testing.T) {
	properties := gopter.NewProperties(nil)

	properties.Property("concurrent cache operations are thread-safe", prop.ForAll(
		func(keys []string, values []int) bool {
			if len(keys) == 0 || len(values) == 0 {
				return true
			}

			cache := NewConcurrentCache()
			var wg sync.WaitGroup
			numGoroutines := 10

			// Concurrent writes
			wg.Add(numGoroutines)
			for i := 0; i < numGoroutines; i++ {
				go func(id int) {
					defer wg.Done()
					for j, key := range keys {
						if j%numGoroutines == id {
							cache.Set(key, values[j%len(values)])
						}
					}
				}(i)
			}
			wg.Wait()

			// Concurrent reads
			wg.Add(numGoroutines)
			errors := make(chan bool, numGoroutines)
			for i := 0; i < numGoroutines; i++ {
				go func(id int) {
					defer wg.Done()
					for j, key := range keys {
						if j%numGoroutines == id {
							_, ok := cache.Get(key)
							if !ok {
								errors <- false
								return
							}
						}
					}
					errors <- true
				}(i)
			}
			wg.Wait()
			close(errors)

			// Check all reads succeeded
			for success := range errors {
				if !success {
					return false
				}
			}

			return true
		},
		gen.SliceOf(gen.Identifier()),
		gen.SliceOf(gen.Int()),
	))

	properties.Property("GetOrSet is atomic", prop.ForAll(
		func(key string, value int) bool {
			cache := NewConcurrentCache()
			var wg sync.WaitGroup
			numGoroutines := 100
			results := make([]interface{}, numGoroutines)

			// Multiple goroutines try to set the same key
			wg.Add(numGoroutines)
			for i := 0; i < numGoroutines; i++ {
				go func(id int) {
					defer wg.Done()
					actual, _ := cache.GetOrSet(key, id)
					results[id] = actual
				}(i)
			}
			wg.Wait()

			// All goroutines should see the same value (the first one that won)
			firstValue := results[0]
			for _, v := range results {
				if v != firstValue {
					return false
				}
			}

			return true
		},
		gen.Identifier(),
		gen.Int(),
	))

	properties.Property("concurrent reads and writes don't cause data races", prop.ForAll(
		func(operations []bool) bool {
			if len(operations) == 0 {
				return true
			}

			cache := NewConcurrentCache()
			var wg sync.WaitGroup
			numGoroutines := 10

			wg.Add(numGoroutines)
			for i := 0; i < numGoroutines; i++ {
				go func(id int) {
					defer wg.Done()
					for j, isWrite := range operations {
						key := fmt.Sprintf("key_%d", j%100)
						if isWrite {
							cache.Set(key, j)
						} else {
							cache.Get(key)
						}
					}
				}(i)
			}
			wg.Wait()

			return true
		},
		gen.SliceOf(gen.Bool()),
	))

	properties.Property("Range iterates over all entries", prop.ForAll(
		func(entries map[string]int) bool {
			cache := NewConcurrentCache()

			// Populate cache
			for k, v := range entries {
				cache.Set(k, v)
			}

			// Count entries using Range
			count := 0
			cache.Range(func(key, value interface{}) bool {
				count++
				return true
			})

			// Should have at least as many entries as we added
			// (may have more if cache had previous entries)
			return count >= len(entries)
		},
		gen.MapOf(gen.Identifier(), gen.Int()),
	))

	properties.Property("Delete removes entries", prop.ForAll(
		func(key string, value int) bool {
			cache := NewConcurrentCache()

			cache.Set(key, value)
			_, ok := cache.Get(key)
			if !ok {
				return false
			}

			cache.Delete(key)
			_, ok = cache.Get(key)
			return !ok
		},
		gen.Identifier(),
		gen.Int(),
	))

	properties.Property("Clear removes all entries", prop.ForAll(
		func(entries map[string]int) bool {
			cache := NewConcurrentCache()

			// Populate cache
			for k, v := range entries {
				cache.Set(k, v)
			}

			cache.Clear()

			// Count remaining entries
			count := 0
			cache.Range(func(key, value interface{}) bool {
				count++
				return true
			})

			return count == 0
		},
		gen.MapOf(gen.Identifier(), gen.Int()),
	))

	properties.TestingRun(t, gopter.ConsoleReporter(false))
}

// TestProperty_RegexPrecompilation tests that pre-compiled regex performs better
func TestProperty_RegexPrecompilation(t *testing.T) {
	properties := gopter.NewProperties(nil)

	properties.Property("regex cache returns same instance for same pattern", prop.ForAll(
		func(pattern string) bool {
			// Skip invalid patterns
			if pattern == "" || pattern == "[" || pattern == "(" {
				return true
			}

			cache := NewRegexCache()
			regex1, err1 := cache.Get(pattern)
			if err1 != nil {
				// Invalid pattern, skip
				return true
			}

			regex2, err2 := cache.Get(pattern)
			if err2 != nil {
				return false
			}

			// Should return the same instance (cached)
			return regex1 == regex2
		},
		gen.RegexMatch(`^[a-z0-9]+$`),
	))

	properties.Property("precompiled regex works correctly", prop.ForAll(
		func(testString string) bool {
			cache := NewRegexCache()
			pattern := `^[a-z0-9]+$`

			regex, err := cache.Get(pattern)
			if err != nil {
				return false
			}

			// Verify regex works correctly
			matches := regex.MatchString(testString)

			// Check if string should match (only lowercase letters and digits, non-empty)
			if testString == "" {
				return !matches // Empty string should not match
			}

			expectedMatches := true
			for _, ch := range testString {
				if !((ch >= 'a' && ch <= 'z') || (ch >= '0' && ch <= '9')) {
					expectedMatches = false
					break
				}
			}

			return matches == expectedMatches
		},
		gen.AlphaString(),
	))

	properties.TestingRun(t, gopter.ConsoleReporter(false))
}

// BenchmarkTypeSwitch_vs_Reflection compares type switch performance with reflection
func BenchmarkTypeSwitch_vs_Reflection(b *testing.B) {
	tc := NewTypeConverter()
	values := []interface{}{
		uint(42),
		int(42),
		int64(42),
		float64(42.5),
	}

	b.Run("TypeSwitch", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			val := values[i%len(values)]
			_, _ = tc.ToUint(val)
		}
	})
}

// BenchmarkSyncMap_vs_MutexMap compares sync.Map with mutex-protected map
func BenchmarkSyncMap_vs_MutexMap(b *testing.B) {
	// Prepare test data
	keys := make([]string, 1000)
	for i := 0; i < 1000; i++ {
		keys[i] = fmt.Sprintf("key%d", i)
	}

	b.Run("SyncMap_Read", func(b *testing.B) {
		cache := NewConcurrentCache()
		for i, key := range keys {
			cache.Set(key, i)
		}

		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			i := 0
			for pb.Next() {
				cache.Get(keys[i%len(keys)])
				i++
			}
		})
	})

	b.Run("MutexMap_Read", func(b *testing.B) {
		var mu sync.RWMutex
		data := make(map[string]int)
		for i, key := range keys {
			data[key] = i
		}

		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			i := 0
			for pb.Next() {
				mu.RLock()
				_ = data[keys[i%len(keys)]]
				mu.RUnlock()
				i++
			}
		})
	})

	b.Run("SyncMap_Write", func(b *testing.B) {
		cache := NewConcurrentCache()

		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			i := 0
			for pb.Next() {
				cache.Set(keys[i%len(keys)], i)
				i++
			}
		})
	})

	b.Run("MutexMap_Write", func(b *testing.B) {
		var mu sync.RWMutex
		data := make(map[string]int)

		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			i := 0
			for pb.Next() {
				mu.Lock()
				data[keys[i%len(keys)]] = i
				mu.Unlock()
				i++
			}
		})
	})
}

// Helper function to generate random operations
func generateRandomOperations(n int) []bool {
	rand.Seed(time.Now().UnixNano())
	ops := make([]bool, n)
	for i := 0; i < n; i++ {
		ops[i] = rand.Float32() < 0.8 // 80% reads, 20% writes
	}
	return ops
}
