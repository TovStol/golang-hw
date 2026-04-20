package hw04lrucache

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type entry struct {
	key   Key
	value interface{}
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
	if item, exists := c.items[key]; exists {
		c.queue.MoveToFront(item)
		item.Value = entry{key: key, value: value}
		return true
	}

	if c.queue.Len() >= c.capacity {
		oldest := c.queue.Back()
		if oldest != nil {
			oldEntry := oldest.Value.(entry)
			delete(c.items, oldEntry.key)
			c.queue.Remove(oldest)
		}
	}

	newItem := c.queue.PushFront(entry{key: key, value: value})
	c.items[key] = newItem
	return false
}

func (c *lruCache) Get(key Key) (interface{}, bool) {
	if item, exists := c.items[key]; exists {
		c.queue.MoveToFront(item)
		return item.Value.(entry).value, true
	}
	return nil, false
}

func (c *lruCache) Clear() {
	c.queue = NewList()
	c.items = make(map[Key]*ListItem, c.capacity)
}
