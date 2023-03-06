package main

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/allegro/bigcache"
	"github.com/go-redis/redis/v8"
	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v8"
	"github.com/willf/bloom"
	"golang.org/x/sync/singleflight"
)

type MultiCache struct {
	memoryCache *bigcache.BigCache
	redisClient *redis.Client
	redisMutex  *redsync.Mutex
	bloomFilter *bloom.BloomFilter
	singleGroup singleflight.Group
	mu          sync.Mutex
}

func NewMultiCache(redisAddr string) (*MultiCache, error) {
	memoryCache, err := bigcache.NewBigCache(bigcache.DefaultConfig(10 * time.Minute))
	if err != nil {
		return nil, err
	}
	redisClient := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: "redispw",
	})
	redisMutex := redsync.New(goredis.NewPool(redisClient))
	bloomFilter := bloom.New(10000000, 5)
	return &MultiCache{
		memoryCache: memoryCache,
		redisClient: redisClient,
		redisMutex:  redisMutex.NewMutex("redsync_mutex"),
		bloomFilter: bloomFilter,
	}, nil
}

// Get gets the value for the given keys from the cache. If the value is not found in the cache,
// it will be loaded using the provided loader function and set to the cache.
func (mc *MultiCache) Get(keys []string, ttl time.Duration, loaderFunc func([]string) (map[string]interface{}, error)) (map[string]interface{}, error) {
	results := make(map[string]interface{})
	notFoundKeys := make([]string, 0)

	// Try to get the values from the memory cache
	for _, key := range keys {
		val, err := mc.memoryCache.Get(key)
		if err == nil {
			results[key] = val
		} else {
			notFoundKeys = append(notFoundKeys, key)
		}
	}

	// Try to get the values from the redis cache
	if len(notFoundKeys) > 0 {
		vals, err := mc.redisClient.MGet(context.Background(), notFoundKeys...).Result()
		if err == nil {
			for i, key := range notFoundKeys {
				val := vals[i]
				if val != nil {
					mc.memoryCache.Set(key, []byte(val.(string)))
					results[key] = val
				}
			}
		}
	}

	// Try to get the values from the loader function
	if len(notFoundKeys) > 0 {
		vals, err, _ := mc.singleGroup.Do("multi:"+notFoundKeys[0], func() (interface{}, error) {
			// Check if the keys are already in the bloom filter to avoid repeated requests
			mc.mu.Lock()
			defer mc.mu.Unlock()
			notFoundKeys2 := make([]string, 0)
			for _, key := range notFoundKeys {
				if !mc.bloomFilter.Test([]byte(key)) {
					notFoundKeys2 = append(notFoundKeys2, key)
				}
			}
			if len(notFoundKeys2) == 0 {
				return nil, errors.New("keys not found")
			}

			// Acquire a distributed lock to avoid cache stampede
			err := mc.redisMutex.Lock()
			if err != nil {
				return nil, err
			}
			defer mc.redisMutex.Unlock()

			// Try to get the values from the memory cache again
			for _, key := range notFoundKeys2 {
				val, err := mc.memoryCache.Get(key)
				if err == nil {
					results[key] = val
				} else {
					notFoundKeys = append(notFoundKeys, key)
				}
			}

			// Try to get the values from the redis cache again
			if len(notFoundKeys2) > 0 {
				vals, err := mc.redisClient.MGet(context.Background(), notFoundKeys2...).Result()
				if err == nil {
					for i, key := range notFoundKeys2 {
						val := vals[i]
						if val != nil {
							mc.memoryCache.Set(key, []byte(val.(string)))
							results[key] = val
						} else {
							notFoundKeys = append(notFoundKeys, key)
						}
					}
				}
			}

			// Get the values from the loader function and set them to the memory and redis cache
			if len(notFoundKeys2) > 0 {
				vals2, err := loaderFunc(notFoundKeys2)
				if err != nil {
					return nil, err
				}
				for key, val := range vals2 {
					if val != nil {
						mc.memoryCache.Set(key, []byte(val.(string)))
						err = mc.redisClient.Set(context.Background(), key, val, ttl).Err()
						if err != nil {
							return nil, err
						}
						results[key] = val
					}
				}
				// Add the keys to the bloom filter to avoid repeated requests
				for _, key := range notFoundKeys2 {
					mc.bloomFilter.Add([]byte(key))
				}
			}

			return results, nil
		})
		if err == nil {
			return vals.(map[string]interface{}), nil
		}
	}

	return results, nil
}

// Set sets the values for the given keys to the cache.
func (mc *MultiCache) Set(keyVals map[string]string, ttl time.Duration) error {
	cmds := make([]redis.Cmder, 0, len(keyVals))
	for key, val := range keyVals {
		err := mc.redisClient.Set(context.Background(), key, val, ttl).Err()
		if err != nil {
			return err
		}
		cmds = append(cmds, mc.redisClient.Set(context.Background(), key, val, ttl))
	}
	_, err := mc.redisClient.Pipelined(context.Background(), func(pipe redis.Pipeliner) error {
		for range cmds {
			pipe.Exec(context.Background())
		}
		return nil
	})
	if err != nil {
		return err
	}
	for key, val := range keyVals {
		mc.memoryCache.Set(key, []byte(val))
	}
	return nil
}

// Delete deletes the value for the given key from the cache.
func (mc *MultiCache) Delete(key string) error {
	err := mc.redisClient.Del(context.Background(), key).Err()
	if err != nil {
		return err
	}
	mc.memoryCache.Delete(key)
	return nil
}
