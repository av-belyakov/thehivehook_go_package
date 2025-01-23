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
				crm.DeleteElement(id)

				fmt.Println("func 'automaticExecutionMethods' new tick: before delete id:", id)

				continue
			}

			//удаление записей если функция в настоящее время не выполняется и вернула
			// положительный результат
			if !storage.isFunctionExecution && storage.isCompletedSuccessfully {
				crm.DeleteElement(id)

				fmt.Println("func 'automaticExecutionMethods' new tick: delete id:", id)

				continue
			}

			if storage.isFunctionExecution {
				continue
			}

			//выполнение кешированной функции
			go func(cache *CacheRunningFunctions, id string, f func(int) bool) {
				fmt.Println("func 'automaticExecutionMethods' new tick: cacheFunc, id:", id)

				//устанавливает что функция выполняется
				cache.setIsFunctionExecution(id)
				//увеличивает количество попыток выполения функции на 1
				cache.increaseNumberAttempts(id)

				//при вызове, функция принимает количество попыток обработки
				if f(cache.getNumberAttempts(id)) {
					cache.setIsCompletedSuccessfully(id)
				}

				//отмечает что функция завершила выполнение
				cache.setIsFunctionNotExecution(id)
			}(crm, id, storage.cacheFunc)
		}
		crm.cacheStorage.mutex.RUnlock()
	}
}
