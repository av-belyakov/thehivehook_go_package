package cacherunningfunctions

import (
	"context"
	"errors"
	"fmt"
	"time"
)

// CreateCach создает новое хранилище кэширующее исполняемые функции. Время по
// истечение которого кэшированная функция будет удалена, задается в секундах и
// варьируется в диапазоне от 10 до 86400 секунд, что эквивалентно одним суткам.
func CreateCach(ctx context.Context, ttl int) (*CacheRunningFunctions, error) {
	cacheRunCom := CacheRunningFunctions{}

	if ttl < 10 || ttl > 86400 {
		return &cacheRunCom, errors.New("the lifetime of the temporary information should not be less than 10 seconds and more than 86400 seconds")
	}

	timeToLive, err := time.ParseDuration(fmt.Sprintf("%ds", ttl))
	if err != nil {
		return &cacheRunCom, err
	}

	cacheRunCom.ttl = timeToLive
	cacheRunCom.cacheStorage = cacheStorageParameters{
		storages: make(map[string]storageParameters),
	}

	go cacheRunCom.automaticExecutionMethods(ctx)

	return &cacheRunCom, err
}

func (crm *CacheRunningFunctions) automaticExecutionMethods(ctx context.Context) {
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

			//выполнение кешированной функции
			if v.cacheFunc() {
				delete(crm.cacheStorage.storages, k)
			}
		}
		crm.cacheStorage.mutex.Unlock()
	}
}

// SetMethod создает новую запись, принимает значение которое нужно сохранить
// и id по которому данное значение можно будет найти
func (crm *CacheRunningFunctions) SetMethod(id string, f func() bool) string {
	crm.cacheStorage.mutex.Lock()
	defer crm.cacheStorage.mutex.Unlock()

	crm.cacheStorage.storages[id] = storageParameters{
		timeExpiry: time.Now().Add(crm.ttl),
		cacheFunc:  f,
	}

	return id
}

// GetMethod возвращает данные по полученому id
func (crm *CacheRunningFunctions) GetMethod(id string) (func() bool, bool) {
	if stoarge, ok := crm.cacheStorage.storages[id]; ok {
		return stoarge.cacheFunc, ok
	}

	return nil, false
}

// DeleteElement удаляет заданный элемент по его id
func (crm *CacheRunningFunctions) DeleteElement(id string) {
	crm.cacheStorage.mutex.Lock()
	defer crm.cacheStorage.mutex.Unlock()

	delete(crm.cacheStorage.storages, id)
}
