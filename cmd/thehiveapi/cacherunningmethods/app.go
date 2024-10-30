package cacherunningmethods

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"
)

var (
	once        sync.Once
	cacheRunCom CacheRunningMethods
)

func New(ctx context.Context, ttl int) (*CacheRunningMethods, error) {
	cacheRunCom = CacheRunningMethods{}

	if ttl < 5 || ttl > 86400 {
		return &cacheRunCom, errors.New("the lifetime of the temporary information should not be less than 10 seconds and more than 86400 seconds")
	}

	var err error
	once.Do(func() {
		timeToLive, newErr := time.ParseDuration(fmt.Sprintf("%ds", ttl))
		if newErr != nil {
			err = newErr

			return
		}

		cacheRunCom.ttl = timeToLive
		cacheRunCom.cacheStorage = cacheStorageParameters{
			storages: make(map[string]storageParameters),
		}

		go cacheRunCom.automaticExecutionMethods(ctx)
	})

	return &cacheRunCom, err
}

func (crm *CacheRunningMethods) automaticExecutionMethods(ctx context.Context) {
	tick := time.NewTicker(5 * time.Second)

	go func(ctx context.Context, tick *time.Ticker) {
		<-ctx.Done()
		tick.Stop()
	}(ctx, tick)

	for range tick.C {
		crm.cacheStorage.mutex.Lock()
		for k, v := range crm.cacheStorage.storages {
			//удаляем если записи слишком старые
			if v.timeExpiry.Before(time.Now()) {
				delete(crm.cacheStorage.storages, k)
			}

			if v.cacheMethod() {
				delete(crm.cacheStorage.storages, k)
			}
		}
		crm.cacheStorage.mutex.Unlock()
	}
}

// SetMethod создает новую запись, принимает значение которое нужно сохранить
// и id по которому данное значение можно будет найти
func (crm *CacheRunningMethods) SetMethod(id string, f func() bool) string {
	crm.cacheStorage.mutex.Lock()
	defer crm.cacheStorage.mutex.Unlock()

	crm.cacheStorage.storages[id] = storageParameters{
		timeExpiry:  time.Now().Add(crm.ttl),
		cacheMethod: f,
	}

	return id
}

// GetMethod возвращает данные по полученому id
func (crm *CacheRunningMethods) GetMethod(id string) (func() bool, bool) {
	if stoarge, ok := crm.cacheStorage.storages[id]; ok {
		return stoarge.cacheMethod, ok
	}

	return nil, false
}

// DeleteElement удаляет заданный элемент по его id
func (crm *CacheRunningMethods) DeleteElement(id string) {
	crm.cacheStorage.mutex.Lock()
	defer crm.cacheStorage.mutex.Unlock()

	delete(crm.cacheStorage.storages, id)
}
