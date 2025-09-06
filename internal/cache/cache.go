package cache

import (
	"context"
	"log"
	"sync"

	"github.com/venexene/wbl0-orders-service/internal/models"
	"github.com/venexene/wbl0-orders-service/internal/db"
)

// Структура кэша
type Cache struct {
	elems map[string] *models.Order
	mu sync.RWMutex
}

// Конструктор кэша
func NewCache() *Cache {
	elems := make(map[string] *models.Order)

	cache := Cache {
		elems: elems,
	}

	return &cache
}


// Добавление данных в кэш
func (c *Cache) Set(order *models.Order) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.elems[order.OrderUID] = order
	log.Printf("Added order %s to cache", order.OrderUID)
}

// Получение данных из кэша
func (c *Cache) Get(key string) (*models.Order, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	order, exists := c.elems[key]
	return order, exists
}

// Удаление из кэша
func (c *Cache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.elems, key)
	log.Printf("Deleted order %s from cache", key)
}

// Получение UID всех заказов в кэше
func (c *Cache) GetAllUIDs() []string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	uids := make([]string, 0, len(c.elems))
	for uid := range c.elems {
		uids = append(uids, uid)
	}
	return uids
}

// Получение размера кэша
func (c *Cache) Size() int { 
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.elems)
}

// Заполнение кэша
func (c *Cache) Populate(ctx context.Context, storage *database.Storage) error {
	uids, err := storage.GetAllOrdersUID(ctx)
	if err != nil {
		log.Printf("Failed to get all orders from db: %v", err)
		return err
	}

	var loadCount int
	for _, uid := range uids {
		order, err := storage.GetOrderByUID(ctx, uid)
		if err != nil {
			log.Printf("Failed to load order %s into cache: %v", uid, err)
			continue
		}
		c.Set(order)
		loadCount++
	}

	log.Printf("Populated cache with %d orders", loadCount)
	return nil
}