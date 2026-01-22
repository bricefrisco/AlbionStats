package util

type IDCache struct {
	cache    map[int64]struct{}
	order    []int64
	capacity int
}

func NewIDCache(capacity int) *IDCache {
	return &IDCache{
		cache:    make(map[int64]struct{}),
		order:    make([]int64, 0, capacity),
		capacity: capacity,
	}
}

func (c *IDCache) Exists(id int64) bool {
	_, exists := c.cache[id]
	return exists
}

func (c *IDCache) Add(id int64) bool {
	if _, exists := c.cache[id]; exists {
		return false
	}

	c.cache[id] = struct{}{}
	c.order = append(c.order, id)

	if len(c.order) > c.capacity {
		toRemove := c.order[0]
		delete(c.cache, toRemove)
		c.order = c.order[1:]
	}

	return true
}
