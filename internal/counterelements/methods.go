package counterelements

import (
	"context"
	"time"
)

// Start автоматическая обработка удаления объектов с истёкшим временеем жизни
func (c *CaseElementCount) Start(ctx context.Context) {
	go func() {
		tick := time.NewTicker(1 * time.Second)
		defer tick.Stop()

		for {
			select {
			case <-ctx.Done():
				return

			case <-tick.C:
				//удаление всех объектов у которых истекло время жизни
				c.deleteObjectsAfterLifetimeExpired()

			}
		}
	}()
}

// Get возвращает количество элементов по ключу, если ключ не найден то -1
func (c *CaseElementCount) Get(key string) int {
	c.mtx.RLock()
	defer c.mtx.RUnlock()

	if _, ok := c.counters[key]; !ok {
		return -1
	}

	return c.counters[key].count
}

// Add добавляет новый ключ
func (c *CaseElementCount) Add(key string) {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	currentTime := time.Now()
	if _, ok := c.counters[key]; !ok {
		c.counters[key] = struct {
			timeExpiry time.Time
			createTime time.Time
			count      int
		}{
			timeExpiry: currentTime.Add(time.Duration(c.deletingAfterTime) * time.Second),
			createTime: currentTime,
			count:      1,
		}

		return
	}

	element := c.counters[key]
	element.count++

	c.counters[key] = element
}

// Size количество элементов
func (c *CaseElementCount) Size() int {
	c.mtx.RLock()
	defer c.mtx.RUnlock()

	return len(c.counters)
}

// Done уменьшает количество элементов по ключу на единицу
func (c *CaseElementCount) Done(key string) {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	if element, ok := c.counters[key]; ok && element.count > 0 {
		element := c.counters[key]
		element.count--

		c.counters[key] = element
	}
}

// DeleteKey удаляет ключ
func (c *CaseElementCount) DeleteKey(key string) {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	delete(c.counters, key)
}

// deleteObjectAfterLifetimeExpired удаляет все объекты у которых истекло время жизни
func (c *CaseElementCount) deleteObjectsAfterLifetimeExpired() {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	for k, v := range c.counters {
		if v.timeExpiry.Before(time.Now()) {
			delete(c.counters, k)
		}
	}
}
