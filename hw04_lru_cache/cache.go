package hw04lrucache

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type CacheElement struct {
	Key   Key
	Value interface{}
}

type lruCache struct {
	capacity int
	queue    List
	items    map[Key]*ListItem
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}

func (c *lruCache) Set(key Key, value interface{}) bool {
	item, ok := c.items[key]
	if ok {
		item.Value = &CacheElement{
			Key:   key,
			Value: value,
		}
		c.queue.MoveToFront(item)

		return true
	}

	if c.queue.Len() == c.capacity {
		lastItem := c.queue.Back()
		c.queue.Remove(lastItem)
		cacheElement := lastItem.Value.(*CacheElement)
		delete(c.items, cacheElement.Key)
	}

	c.items[key] = c.queue.PushFront(
		&CacheElement{
			Key:   key,
			Value: value,
		},
	)

	return false
}

func (c *lruCache) Get(key Key) (interface{}, bool) {
	item, ok := c.items[key]
	if ok {
		c.queue.MoveToFront(item)

		cacheElement := item.Value.(*CacheElement)

		return cacheElement.Value, true
	}

	return nil, false
}

func (c *lruCache) Clear() {
	c.queue = NewList()
	c.items = make(map[Key]*ListItem, c.capacity)
}
