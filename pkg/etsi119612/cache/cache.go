package cache

import(
	"context"
	"sync"
	"time"
	"fmt"
)

type CacheInterface interface {
	Get(ctxt context.Context, key string)  (any, error)
	Set(ctxt context.Context,key string, value any, ttl time.Duration) error
	Delete(ctxt context.Context, key string)  error 
}

type CacheEntry  struct {
	Value  any
	ExpiresAt time.Time
}

type CustomCache struct {
  Data map[string]*CacheEntry 
  mu sync.RWMutex
}

func  NewCustomCache () *CustomCache {
return &CustomCache{Data: make(map[string]*CacheEntry)}
}

func (c *CustomCache) Get(ctxt context.Context, key string) (*CacheEntry, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	value, exists := c.Data[key]
	return value, exists
}


func (c *CustomCache) Set(ctxt context.Context, key string, value any, ttl time.Duration) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	entry := &CacheEntry{Value: value, ExpiresAt: time.Now().Add(ttl)}
	c.Data[key]=entry
	return nil
}

func (c *CustomCache) Delete(ctxt context.Context, key string) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	_, exists :=c.Data[key]
	if !exists{
		return fmt.Errorf("key %v not found", key)
	}
	delete(c.Data, key)
	return nil
}


