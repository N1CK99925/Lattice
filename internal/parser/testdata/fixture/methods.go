package fixture

type Cache struct {
	size int
}

func NewCache() *Cache { return &Cache{} }

func (c *Cache) Put(k string, v int) { c.size = v }

func (c *Cache) Get(k string) int { return c.size }
