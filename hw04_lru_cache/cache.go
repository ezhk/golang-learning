package hw04_lru_cache //nolint:golint,stylecheck

type Key string

type Cache interface {
	Set(key string, value interface{}) bool
	Get(key string) (interface{}, bool)
	Clear()
}

type lruCache struct {
	capacity int
	Queue    List
	Items    map[string]*listItem
}

type cacheItem struct {
	cKey   string
	cValue interface{}
}

func NewCache(capacity int) Cache {
	list := NewList()
	return &lruCache{
		capacity: capacity,
		Queue:    list,
		Items:    make(map[string]*listItem),
	}
}

func NewItem(key string, value interface{}) *cacheItem {
	return &cacheItem{
		cKey:   key,
		cValue: value,
	}
}

func (c lruCache) Set(key string, value interface{}) bool {
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

	item := NewItem(key, value)
	c.Items[key] = c.Queue.PushFront(item)
	return false
}

func (c lruCache) Get(key string) (interface{}, bool) {
	if item, ok := c.Items[key]; ok {
		c.Queue.MoveToFront(item)
		return item.Value.(*cacheItem).cValue, true
	}

	return nil, false
}

func (c lruCache) Clear() {
	for c.Queue.Len() > 0 {
		c.Queue.Remove(c.Queue.Back())
	}
	c.Items = make(map[string]*listItem)
}
