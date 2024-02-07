package cache

import "sync"

/*простой кеш ключ-значение с потокобезопастными операциями чтения и записи*/
type Cache struct {
	data  map[string]interface{}
	mutex sync.RWMutex
}

func NewCache() *Cache {
	return &Cache{
		data: make(map[string]interface{}),
	}
}

func (c *Cache) Get(id string) (interface{}, bool) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	value, ok := c.data[id]
	return value, ok
}

func (c *Cache) Set(id string, value interface{}) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.data[id] = value
}
