package cacherunningfunctions

import (
	"context"
	"errors"
	"fmt"
	"time"
)

// CreateCache создает новое хранилище кэширующее исполняемые функции. Время по
// истечение которого кэшированная функция будет удалена, задается в секундах и
// варьируется в диапазоне от 10 до 86400 секунд, что эквивалентно одним суткам.
func CreateCache(ctx context.Context, ttl int) (*CacheRunningFunctions, error) {
	cacheRunCom := CacheRunningFunctions{
		ttl: time.Duration(30 * time.Second),
	}

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
		crm.cacheStorage.mutex.RLock()
		for id, storage := range crm.cacheStorage.storages {
			fmt.Println("func 'automaticExecutionMethods' new tick:")

			//удаление слишком старых записей
			if storage.timeExpiry.Before(time.Now()) {
				//delete(crm.cacheStorage.storages, k)
				go crm.DeleteElement(id)

				fmt.Println("func 'automaticExecutionMethods' new tick: before delete id:", id)

				continue
			}

			//удаление записей если функция в настоящее время не выполняется и вернула
			// положительный результат
			if storage.isCompletedSuccessfully {
				//delete(crm.cacheStorage.storages, k)
				go crm.DeleteElement(id)

				fmt.Println("func 'automaticExecutionMethods' new tick: delete id:", id)

				continue
			}

			//выполнение кешированной функции
			go func(cache *CacheRunningFunctions, id string, f func() bool) {
				fmt.Println("func 'automaticExecutionMethods' new tick: cacheFunc, id:", id)

				cache.setIsFunctionRunning(id)
				if f() {
					cache.setIsCompletedSuccessfully(id)
				}
				cache.setIsFunctionNotRunning(id)
			}(crm, id, storage.cacheFunc)
		}
		crm.cacheStorage.mutex.RUnlock()
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
	crm.cacheStorage.mutex.RLock()
	defer crm.cacheStorage.mutex.Unlock()

	if storage, ok := crm.cacheStorage.storages[id]; ok {
		return storage.cacheFunc, ok
	}

	return nil, false
}

// DeleteElement удаляет заданный элемент по его id
func (crm *CacheRunningFunctions) DeleteElement(id string) {
	crm.cacheStorage.mutex.Lock()
	defer crm.cacheStorage.mutex.Unlock()

	delete(crm.cacheStorage.storages, id)
}

// setIsCompletedSuccessfully выполняемая функция завершилась успехом
func (crm *CacheRunningFunctions) setIsCompletedSuccessfully(id string) {
	crm.cacheStorage.mutex.Lock()
	defer crm.cacheStorage.mutex.Unlock()

	storage, ok := crm.cacheStorage.storages[id]
	if !ok {
		return
	}

	storage.isCompletedSuccessfully = true
	crm.cacheStorage.storages[id] = storage
}

// setIsFunctionRunning функция находится в процессе выполнения
func (crm *CacheRunningFunctions) setIsFunctionRunning(id string) {
	crm.cacheStorage.mutex.Lock()
	defer crm.cacheStorage.mutex.Unlock()

	storage, ok := crm.cacheStorage.storages[id]
	if !ok {
		return
	}

	storage.isFunctionRunning = true
	crm.cacheStorage.storages[id] = storage
}

// setIsFunctionNotRunning функция не выполняется
func (crm *CacheRunningFunctions) setIsFunctionNotRunning(id string) {
	crm.cacheStorage.mutex.Lock()
	defer crm.cacheStorage.mutex.Unlock()

	storage, ok := crm.cacheStorage.storages[id]
	if !ok {
		return
	}

	storage.isFunctionRunning = false
	crm.cacheStorage.storages[id] = storage
}
