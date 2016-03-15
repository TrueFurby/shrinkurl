package main

import (
	"sync"
	"time"
)

type CacheItem struct {
	key string
	val string
	age int64
}

type Cache struct {
	items     map[string]CacheItem
	lock      *sync.RWMutex
	maxAgeSec int64
}

func NewCache(maxAgeSeconds int64) *Cache {
	return &Cache{
		items:     make(map[string]CacheItem, 1024),
		lock:      new(sync.RWMutex),
		maxAgeSec: maxAgeSeconds,
	}
}

func (c *Cache) Get(key string) string {
	c.lock.RLock()
	defer c.lock.RUnlock()

	now := time.Now()
	currTime := now.Unix()
	if currTime-c.items[key].age > c.maxAgeSec {
		return ""
	}
	return c.items[key].val
}

func (c *Cache) Add(key string, val string) {
	c.lock.Lock()
	defer c.lock.Unlock()

	now := time.Now()
	age := now.Unix()
	c.items[key] = CacheItem{key, val, age}
}

func (c *Cache) Remove(id string) {
	c.lock.Lock()
	defer c.lock.Unlock()
	delete(c.items, id)
}
