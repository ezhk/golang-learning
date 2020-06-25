package hw04_lru_cache //nolint:golint,stylecheck

import (
	"sync"
)

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type lruCache struct {
	sync.Mutex
	capacity int
	Queue    List
	Items    map[Key]*listItem
}

type cacheItem struct {
	cKey   Key
	cValue interface{}
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		Queue:    NewList(),
		Items:    make(map[Key]*listItem),
	}
}

func newItem(key Key, value interface{}) *cacheItem {
	return &cacheItem{
		cKey:   key,
		cValue: value,
	}
}

func (c *lruCache) Set(key Key, value interface{}) bool {
	c.Lock()
	defer c.Unlock()

	if item, ok := c.Items[key]; ok {
		item.Value.(*cacheItem).cValue = value
		c.Queue.MoveToFront(item)
		return true
	}

	for c.Queue.Len() >= c.capacity {
		latestItem := c.Queue.Back()

		c.Queue.Remove(latestItem)
		delete(c.Items, latestItem.Value.(*cacheItem).cKey)
	}

	item := newItem(key, value)
	c.Items[key] = c.Queue.PushFront(item)
	return false
}

func (c *lruCache) Get(key Key) (interface{}, bool) {
	c.Lock()
	defer c.Unlock()

	if item, ok := c.Items[key]; ok {
		c.Queue.MoveToFront(item)
		return item.Value.(*cacheItem).cValue, true
	}

	return nil, false
}

func (c *lruCache) Clear() {
	c.Lock()
	defer c.Unlock()

	c.Queue = NewList()
	c.Items = make(map[Key]*listItem)
}
