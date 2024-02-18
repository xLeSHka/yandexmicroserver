package cache

import "sync"

/*простой кеш ключ-значение с потокобезопастными операциями чтения и записи*/
type Cache struct {
	Data  map[string]interface{}
	mutex sync.RWMutex
}

func NewCache() *Cache {
	return &Cache{
		Data: make(map[string]interface{}),
	}
}
func (c *Cache) GetAll() map[string]interface{} {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	return c.Data
}
func (c *Cache) Delete(id string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	delete(c.Data, id)
}
func (c *Cache) Get(id string) (interface{}, bool) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	value, ok := c.Data[id]
	return value, ok
}

func (c *Cache) Set(id string, value interface{}) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.Data[id] = value
}
