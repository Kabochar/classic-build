package v2

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/allegro/bigcache"
	"github.com/go-redis/redis/v8"
	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v8"
	"golang.org/x/sync/singleflight"
)

// MultiCache 是一个多级缓存结构体
type MultiCache struct {
	memoryCache *bigcache.BigCache
	redisClient *redis.Client
	redisMutex  *redsync.Mutex
	singleGroup singleflight.Group
	mu          sync.Mutex
}

// NewMultiCache 创建一个新的多级缓存实例
func NewMultiCache(redisAddr string) (*MultiCache, error) {
	// 创建内存缓存
	memoryCache, err := bigcache.NewBigCache(bigcache.DefaultConfig(10 * time.Minute))
	if err != nil {
		return nil, err
	}

	// 创建 Redis 客户端
	redisClient := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: "redispw",
	})

	// 创建分布式锁
	redisMutex := redsync.New(goredis.NewPool(redisClient))

	return &MultiCache{
		memoryCache: memoryCache,
		redisClient: redisClient,
		redisMutex:  redisMutex.NewMutex(""),
	}, nil
}

// Get 从多级缓存中获取数据，如果缓存中不存在，则从数据库中加载数据
func (mc *MultiCache) Get(ctx context.Context, keys []string, ttl time.Duration, loaderFunc func([]string) (map[string]interface{}, error)) (map[string]interface{}, error) {
	results := make(map[string]interface{})
	notFoundKeys := make([]string, 0)

	// 尝试从内存缓存中获取数据
	for _, key := range keys {
		val, err := mc.memoryCache.Get(key)
		if err == nil {
			results[key] = val
		} else {
			notFoundKeys = append(notFoundKeys, key)
		}
	}

	// 尝试从 Redis 缓存中获取数据
	if len(notFoundKeys) > 0 {
		vals, err := mc.redisClient.MGet(ctx, notFoundKeys...).Result()
		if err == nil {
			for i, key := range notFoundKeys {
				val := vals[i]
				if val != nil {
					mc.memoryCache.Set(key, []byte(fmt.Sprintf("%v", val)))
					results[key] = val
				}
			}
		}
	}

	// 尝试从数据库中加载数据
	if len(notFoundKeys) > 0 {
		vals, err, _ := mc.singleGroup.Do("", func() (interface{}, error) {
			// 尝试从内存缓存中获取数据
			mc.mu.Lock()
			defer mc.mu.Unlock()
			notFoundKeys2 := make([]string, 0)
			for _, key := range notFoundKeys {
				if _, ok := results[key]; !ok {
					notFoundKeys2 = append(notFoundKeys2, key)
				}
			}
			if len(notFoundKeys2) == 0 {
				return results, nil
			}

			// 获取分布式锁
			err := mc.redisMutex.Lock()
			if err != nil {
				return nil, err
			}
			defer mc.redisMutex.Unlock()

			// 再次尝试从内存缓存中获取数据
			for _, key := range notFoundKeys2 {
				val, err := mc.memoryCache.Get(key)
				if err == nil {
					results[key] = val
				} else {
					notFoundKeys = append(notFoundKeys, key)
				}
			}

			// 再次尝试从 Redis 缓存中获取数据
			if len(notFoundKeys2) > 0 {
				vals, err := mc.redisClient.MGet(ctx, notFoundKeys2...).Result()
				if err == nil {
					for i, key := range notFoundKeys2 {
						val := vals[i]
						if val != nil {
							mc.memoryCache.Set(key, []byte(fmt.Sprintf("%v", val)))
							results[key] = val
						} else {
							notFoundKeys = append(notFoundKeys, key)
						}
					}
				}
			}

			// 从数据库中加载数据
			if len(notFoundKeys2) > 0 {
				vals2, err := loaderFunc(notFoundKeys2)
				if err != nil {
					return nil, err
				}
				for key, val := range vals2 {
					if val != nil {
						mc.memoryCache.Set(key, []byte(fmt.Sprintf("%v", val)))
						err = mc.redisClient.Set(ctx, key, val, ttl).Err()
						if err != nil {
							return nil, err
						}
						results[key] = val
					}
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

// Set 将数据设置到多级缓存中
func (mc *MultiCache) Set(ctx context.Context, keyVals map[string]string, ttl time.Duration) error {
	pipe := mc.redisClient.TxPipeline()
	for key, val := range keyVals {
		pipe.Set(ctx, key, val, ttl)
	}
	_, err := pipe.Exec(ctx)
	if err != nil {
		return err
	}

	// 如果 Redis 写入失败，则不允许写入内存缓存
	for key, val := range keyVals {
		// err = mc.redisClient.Get(ctx, key).Err()
		// if err != nil {
		// 	return err
		// }
		mc.memoryCache.Set(key, []byte(val))
	}

	return nil
}

// Delete 从多级缓存中删除数据
func (mc *MultiCache) Delete(ctx context.Context, key string) error {
	err := mc.redisClient.Del(ctx, key).Err()
	if err != nil {
		return err
	}
	mc.memoryCache.Delete(key)
	return nil
}
