package cache

import (
	"fmt"
	"time"
)

// Example: Batch loading user profiles
// This example demonstrates how to efficiently load multiple user profiles using BatchGetTyped

type UserProfile struct {
	ID      uint   `json:"id"`
	Name    string `json:"name"`
	Email   string `json:"email"`
	StoreID uint   `json:"store_id"`
	RoleID  uint   `json:"role_id"`
}

// LoadUserProfiles loads multiple user profiles from cache using Pipeline
// Falls back to database if cache misses occur
func LoadUserProfiles(userIDs []uint, dbLoader func([]uint) (map[uint]UserProfile, error)) (map[uint]UserProfile, error) {
	if len(userIDs) == 0 {
		return make(map[uint]UserProfile), nil
	}

	// Step 1: Generate cache keys
	keys := make([]string, len(userIDs))
	keyToID := make(map[string]uint, len(userIDs))
	for i, id := range userIDs {
		key := fmt.Sprintf("%suser:profile:%d", CachePrefix, id)
		keys[i] = key
		keyToID[key] = id
	}

	// Step 2: Batch get from cache using Pipeline
	var cachedProfiles map[string]UserProfile
	err := BatchGetTyped(keys, &cachedProfiles)
	if err != nil {
		// If batch get fails, fall back to database for all
		return dbLoader(userIDs)
	}

	// Step 3: Identify cache misses
	result := make(map[uint]UserProfile, len(userIDs))
	missingIDs := []uint{}

	for key, id := range keyToID {
		if profile, exists := cachedProfiles[key]; exists {
			result[id] = profile
		} else {
			missingIDs = append(missingIDs, id)
		}
	}

	// Step 4: Load missing profiles from database
	if len(missingIDs) > 0 {
		dbProfiles, err := dbLoader(missingIDs)
		if err != nil {
			return result, err // Return partial results
		}

		// Step 5: Cache the loaded profiles using Pipeline
		cacheItems := make(map[string]interface{}, len(dbProfiles))
		for id, profile := range dbProfiles {
			key := fmt.Sprintf("%suser:profile:%d", CachePrefix, id)
			cacheItems[key] = profile
			result[id] = profile
		}

		// Batch set to cache (fire and forget)
		go BatchSet(cacheItems, 10*time.Minute)
	}

	return result, nil
}

// Example: Batch caching store statistics
// This example shows how to cache multiple store statistics efficiently

type StoreStats struct {
	StoreID        uint    `json:"store_id"`
	TotalOrders    int     `json:"total_orders"`
	TotalRevenue   float64 `json:"total_revenue"`
	ActiveProducts int     `json:"active_products"`
	LastUpdated    string  `json:"last_updated"`
}

// CacheStoreStatistics caches statistics for multiple stores using Pipeline
func CacheStoreStatistics(stats map[uint]StoreStats, ttl time.Duration) error {
	if len(stats) == 0 {
		return nil
	}

	// Prepare cache items
	cacheItems := make(map[string]interface{}, len(stats))
	for storeID, stat := range stats {
		key := fmt.Sprintf("%sstore:stats:%d", CachePrefix, storeID)
		cacheItems[key] = stat
	}

	// Batch set using Pipeline
	return BatchSet(cacheItems, ttl)
}

// GetStoreStatistics retrieves statistics for multiple stores from cache
func GetStoreStatistics(storeIDs []uint) (map[uint]StoreStats, error) {
	if len(storeIDs) == 0 {
		return make(map[uint]StoreStats), nil
	}

	// Generate cache keys
	keys := make([]string, len(storeIDs))
	keyToID := make(map[string]uint, len(storeIDs))
	for i, id := range storeIDs {
		key := fmt.Sprintf("%sstore:stats:%d", CachePrefix, id)
		keys[i] = key
		keyToID[key] = id
	}

	// Batch get using Pipeline
	var cachedStats map[string]StoreStats
	err := BatchGetTyped(keys, &cachedStats)
	if err != nil {
		return nil, err
	}

	// Convert to result map
	result := make(map[uint]StoreStats, len(cachedStats))
	for key, stats := range cachedStats {
		if id, exists := keyToID[key]; exists {
			result[id] = stats
		}
	}

	return result, nil
}

// Example: Batch cache warming
// This example demonstrates how to pre-load frequently accessed data

// WarmupProductCache pre-loads popular products into cache using Pipeline
func WarmupProductCache(products []interface{}, ttl time.Duration) error {
	if len(products) == 0 {
		return nil
	}

	// Prepare cache items
	cacheItems := make(map[string]interface{}, len(products))
	for i, product := range products {
		// Assuming product has an ID field
		key := fmt.Sprintf("%sproduct:%d", CachePrefix, i+1)
		cacheItems[key] = product
	}

	// Batch set using Pipeline
	return BatchSet(cacheItems, ttl)
}

// Example: Batch invalidation pattern
// This example shows how to invalidate multiple related cache keys

// InvalidateUserRelatedCache invalidates all cache entries related to a user
func InvalidateUserRelatedCache(userID uint) error {
	// Collect all related cache keys
	keys := []string{
		fmt.Sprintf("%suser:profile:%d", CachePrefix, userID),
		fmt.Sprintf("%suser:permissions:%d", CachePrefix, userID),
		fmt.Sprintf("%suser:stores:%d", CachePrefix, userID),
		fmt.Sprintf("%suser:roles:%d", CachePrefix, userID),
	}

	// Delete all keys at once
	return CacheDelete(keys...)
}

// Example: Chunked batch processing
// This example demonstrates how to process large batches in chunks

// BatchGetInChunks retrieves a large number of keys in manageable chunks
func BatchGetInChunks(keys []string, chunkSize int) (map[string]interface{}, error) {
	if chunkSize <= 0 {
		chunkSize = 100 // Default chunk size
	}

	result := make(map[string]interface{})

	// Process in chunks
	for i := 0; i < len(keys); i += chunkSize {
		end := i + chunkSize
		if end > len(keys) {
			end = len(keys)
		}

		chunk := keys[i:end]
		chunkResult, err := BatchGet(chunk)
		if err != nil {
			return result, err
		}

		// Merge chunk results
		for k, v := range chunkResult {
			result[k] = v
		}
	}

	return result, nil
}

// BatchSetInChunks sets a large number of items in manageable chunks
func BatchSetInChunks(items map[string]interface{}, ttl time.Duration, chunkSize int) error {
	if chunkSize <= 0 {
		chunkSize = 100 // Default chunk size
	}

	// Convert map to slice for chunking
	keys := make([]string, 0, len(items))
	for k := range items {
		keys = append(keys, k)
	}

	// Process in chunks
	for i := 0; i < len(keys); i += chunkSize {
		end := i + chunkSize
		if end > len(keys) {
			end = len(keys)
		}

		// Create chunk map
		chunk := make(map[string]interface{}, end-i)
		for j := i; j < end; j++ {
			key := keys[j]
			chunk[key] = items[key]
		}

		// Batch set chunk
		if err := BatchSet(chunk, ttl); err != nil {
			return err
		}
	}

	return nil
}
