package main

import (
	"fmt"
	"math/rand"
	"sync"
	"testing"
	"time"
)

func TestMultiCache(t *testing.T) {
	redisAddr := "localhost:32768"
	mc, err := NewMultiCache(redisAddr)
	if err != nil {
		t.Fatalf("Failed to create MultiCache: %v", err)
	}

	// Generate random keys
	numKeys := 10000
	keys := make([]string, numKeys)
	for i := 0; i < numKeys; i++ {
		keys[i] = fmt.Sprintf("key%d", rand.Intn(numKeys))
	}

	// Test Set function
	start := time.Now()
	keyVals := make(map[string]string)
	for i := 0; i < numKeys; i++ {
		keyVals[keys[i]] = fmt.Sprintf("%d", rand.Intn(100))
	}
	err = mc.Set(keyVals, time.Minute)
	if err != nil {
		t.Fatalf("Failed to set values to MultiCache: %v", err)
	}
	elapsed := time.Since(start)
	fmt.Printf("Set %d keys in MultiCache in %s\n", numKeys, elapsed)

	// Test Get function
	start = time.Now()
	wg := sync.WaitGroup{}
	wg.Add(numKeys)
	for i := 0; i < numKeys; i++ {
		go func(i int) {
			defer wg.Done()
			_, err := mc.Get([]string{keys[i]}, time.Minute, func(keys []string) (map[string]interface{}, error) {
				return nil, nil
			})
			if err != nil {
				t.Fatalf("Failed to get value from MultiCache: %v", err)
			}
		}(i)
	}
	wg.Wait()
	elapsed = time.Since(start)
	fmt.Printf("Get %d keys from MultiCache in %s\n", numKeys, elapsed)
}
