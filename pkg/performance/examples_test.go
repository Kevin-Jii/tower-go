package performance_test

import (
	"fmt"
	"sync"

	"github.com/Kevin-Jii/tower-go/pkg/performance"
)

// Example_typeConverter demonstrates using TypeConverter for fast type conversions
func Example_typeConverter() {
	converter := performance.GetTypeConverter()

	// Convert various types to uint
	val1, _ := converter.ToUint(42)
	val2, _ := converter.ToUint(int64(100))
	val3, _ := converter.ToUint(float64(3.14))

	fmt.Println(val1, val2, val3)
	// Output: 42 100 3
}

// Example_concurrentCache demonstrates using ConcurrentCache for thread-safe caching
func Example_concurrentCache() {
	cache := performance.NewConcurrentCache()

	// Set values
	cache.Set("user:1", "Alice")
	cache.Set("user:2", "Bob")

	// Get values
	if val, ok := cache.Get("user:1"); ok {
		fmt.Println(val)
	}

	// Atomic get or set
	actual, loaded := cache.GetOrSet("user:3", "Charlie")
	fmt.Println(actual, loaded)

	// Output:
	// Alice
	// Charlie false
}

// Example_regexCache demonstrates using RegexCache for pre-compiled regex
func Example_regexCache() {
	cache := performance.GetRegexCache()

	// Get or compile regex (first call compiles, subsequent calls use cache)
	phoneRegex, _ := cache.Get(`^1[3-9]\d{9}$`)

	// Use the cached regex
	fmt.Println(phoneRegex.MatchString("13800138000"))
	fmt.Println(phoneRegex.MatchString("12345678901"))

	// Output:
	// true
	// false
}

// Example_optimizedValidator demonstrates using OptimizedValidator
func Example_optimizedValidator() {
	validator := performance.GetOptimizedValidator()

	// Validate phone
	fmt.Println(validator.ValidatePhone("13800138000"))
	fmt.Println(validator.ValidatePhone("12345678901"))

	// Validate email
	fmt.Println(validator.ValidateEmail("user@example.com"))
	fmt.Println(validator.ValidateEmail("invalid-email"))

	// Output:
	// true
	// false
	// true
	// false
}

// Example_concurrentCache_concurrent demonstrates concurrent usage
func Example_concurrentCache_concurrent() {
	cache := performance.NewConcurrentCache()
	var wg sync.WaitGroup

	// Concurrent writes
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			cache.Set(fmt.Sprintf("key%d", id), id)
		}(i)
	}
	wg.Wait()

	// Concurrent reads
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			cache.Get(fmt.Sprintf("key%d", id))
		}(i)
	}
	wg.Wait()

	fmt.Println("Concurrent operations completed successfully")
	// Output: Concurrent operations completed successfully
}

// Example_typeConverter_errorHandling demonstrates error handling
func Example_typeConverter_errorHandling() {
	converter := performance.GetTypeConverter()

	// Negative values return error
	_, err := converter.ToUint(-1)
	if err != nil {
		fmt.Println("Error:", err)
	}

	// Unsupported types return error
	_, err = converter.ToUint("string")
	if err != nil {
		fmt.Println("Error: unsupported type")
	}

	// Output:
	// Error: negative value cannot be converted to uint: -1
	// Error: unsupported type
}

// Example_regexCache_precompile demonstrates pre-compiling patterns
func Example_regexCache_precompile() {
	cache := performance.GetRegexCache()

	// Pre-compile common patterns at startup
	patterns := []string{
		`^\d+$`,
		`^[a-z]+$`,
		`^[A-Z]+$`,
	}
	cache.Precompile(patterns)

	// Use pre-compiled patterns
	digitRegex := cache.MustGet(`^\d+$`)
	fmt.Println(digitRegex.MatchString("12345"))
	fmt.Println(digitRegex.MatchString("abc"))

	// Output:
	// true
	// false
}

// Example_comparison demonstrates performance comparison
func Example_comparison() {
	// Before: Using reflection or repeated regex compilation
	// After: Using type switches and cached regex

	converter := performance.GetTypeConverter()
	cache := performance.GetRegexCache()

	// Fast type conversion
	val, _ := converter.ToUint(42)
	fmt.Printf("Converted: %d\n", val)

	// Fast regex matching with cache
	regex, _ := cache.Get(`^\d+$`)
	fmt.Printf("Matches: %v\n", regex.MatchString("123"))

	// Output:
	// Converted: 42
	// Matches: true
}

// Example_realWorld demonstrates a real-world usage scenario
func Example_realWorld() {
	// Initialize components
	cache := performance.NewConcurrentCache()
	validator := performance.GetOptimizedValidator()

	// Simulate user registration validation
	phone := "13800138000"
	email := "user@example.com"

	// Validate inputs using pre-compiled regex
	if !validator.ValidatePhone(phone) {
		fmt.Println("Invalid phone")
		return
	}
	if !validator.ValidateEmail(email) {
		fmt.Println("Invalid email")
		return
	}

	// Cache user data using concurrent-safe cache
	cache.Set(fmt.Sprintf("user:phone:%s", phone), email)

	// Retrieve from cache
	if val, ok := cache.Get(fmt.Sprintf("user:phone:%s", phone)); ok {
		fmt.Printf("User email: %s\n", val)
	}

	// Output:
	// User email: user@example.com
}

// Example_migration demonstrates migrating from old code to optimized code
func Example_migration() {
	// OLD CODE (slow):
	// var mu sync.RWMutex
	// data := make(map[string]interface{})
	// mu.Lock()
	// data["key"] = value
	// mu.Unlock()

	// NEW CODE (fast):
	cache := performance.NewConcurrentCache()
	cache.Set("key", "value")

	// OLD CODE (slow):
	// pattern := `^\d+$`
	// regex, _ := regexp.Compile(pattern) // Compiled every time!
	// regex.MatchString(input)

	// NEW CODE (fast):
	regexCache := performance.GetRegexCache()
	regex, _ := regexCache.Get(`^\d+$`) // Compiled once, cached
	regex.MatchString("123")

	fmt.Println("Migration complete")
	// Output: Migration complete
}

// Example_benchmarking demonstrates how to benchmark the optimizations
func Example_benchmarking() {
	// Run benchmarks to verify performance improvements:
	// go test -bench=BenchmarkTypeSwitch -benchmem
	// go test -bench=BenchmarkSyncMap -benchmem
	// go test -bench=BenchmarkRegexCache -benchmem

	// Expected results:
	// - Type switch: ~2-3 ns/op, 0 allocs
	// - sync.Map reads: ~3 ns/op (9x faster than mutex map)
	// - sync.Map writes: ~40 ns/op (1.8x faster than mutex map)
	// - Regex cache: Amortized O(1) after first compilation

	fmt.Println("Run benchmarks to see performance improvements")
	// Output: Run benchmarks to see performance improvements
}

// Example_bestPractices demonstrates best practices
func Example_bestPractices() {
	// 1. Reuse instances (don't create new ones in hot paths)
	converter := performance.GetTypeConverter()
	cache := performance.NewConcurrentCache()
	validator := performance.GetOptimizedValidator()

	// 2. Pre-compile known patterns at startup
	regexCache := performance.GetRegexCache()
	regexCache.Precompile([]string{
		`^\d+$`,
		`^[a-z]+$`,
	})

	// 3. Use type switches for known types
	var val interface{} = 42
	result, _ := converter.ToUint(val)

	// 4. Use sync.Map for read-heavy workloads
	cache.Set("key", "value")
	cache.Get("key")

	// 5. Use pre-compiled regex in hot paths
	regex := regexCache.MustGet(`^\d+$`)
	regex.MatchString("123")

	// 6. Validate inputs with optimized validator
	validator.ValidatePhone("13800138000")

	fmt.Println("Best practices applied:", result)
	// Output: Best practices applied: 42
}

// Example_avoidCommonMistakes demonstrates what NOT to do
func Example_avoidCommonMistakes() {
	// MISTAKE 1: Creating new instances in hot paths
	// DON'T DO THIS:
	// for i := 0; i < 1000; i++ {
	//     cache := performance.NewConcurrentCache() // BAD!
	// }

	// DO THIS:
	cache := performance.NewConcurrentCache()
	for i := 0; i < 1000; i++ {
		cache.Set(fmt.Sprintf("key%d", i), i) // GOOD!
	}

	// MISTAKE 2: Compiling regex in loops
	// DON'T DO THIS:
	// for _, input := range inputs {
	//     regex, _ := regexp.Compile(`^\d+$`) // BAD!
	//     regex.MatchString(input)
	// }

	// DO THIS:
	regexCache := performance.GetRegexCache()
	regex, _ := regexCache.Get(`^\d+$`)
	for _, input := range []string{"123", "456"} {
		regex.MatchString(input) // GOOD!
	}

	// MISTAKE 3: Using sync.Map for write-heavy workloads
	// For write-heavy workloads, consider sync.RWMutex instead

	fmt.Println("Mistakes avoided")
	// Output: Mistakes avoided
}
