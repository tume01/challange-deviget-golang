package sample1

import (
	"fmt"
	"sync"
	"time"
)

// PriceService is a service that we can use to get prices for the items
// Calls to this service are expensive (they take time)
type PriceService interface {
	GetPriceFor(itemCode string) (float64, error)
}

// TransparentCache is a cache that wraps the actual service
// The cache will remember prices we ask for, so that we don't have to wait on every call
// Cache should only return a price if it is not older than "maxAge", so that we don't get stale prices
type TransparentCache struct {
	actualPriceService PriceService
	maxAge             time.Duration
	prices             map[string]float64
	m                  sync.RWMutex
}

func NewTransparentCache(actualPriceService PriceService, maxAge time.Duration) *TransparentCache {
	return &TransparentCache{
		actualPriceService: actualPriceService,
		maxAge:             maxAge,
		prices:             map[string]float64{},
	}
}

// GetPriceFor gets the price for the item, either from the cache or the actual service if it was not cached or too old
func (c *TransparentCache) GetPriceFor(itemCode string) (float64, error) {
	c.m.RLock()
	price, ok := c.prices[itemCode]
	c.m.RUnlock()
	if ok {
		// TODO: check that the price was retrieved less than "maxAge" ago!
		return price, nil
	}
	price, err := c.actualPriceService.GetPriceFor(itemCode)
	if err != nil {
		return 0, fmt.Errorf("getting price from service : %v", err.Error())
	}
	c.m.Lock()
	c.prices[itemCode] = price
	c.m.Unlock()
	time.AfterFunc(c.maxAge, func() {
		c.m.Lock()
		delete(c.prices, itemCode)
		c.m.Unlock()
	})
	return price, nil
}

// GetPricesFor gets the prices for several items at once, some might be found in the cache, others might not
// If any of the operations returns an error, it should return an error as well
func (c *TransparentCache) GetPricesFor(itemCodes ...string) ([]float64, error) {
	m := sync.Mutex{}
	results := []float64{}
	wg := &sync.WaitGroup{}
	var err error
	for _, itemCode := range itemCodes {
		wg.Add(1)
		go func(err error, itemCode string) {
			defer wg.Done()
			price, err := c.GetPriceFor(itemCode)
			if err == nil {
				m.Lock()
				results = append(results, price)
				m.Unlock()
			}
		}(err, itemCode)
	}
	wg.Wait()
	return results, nil
}
