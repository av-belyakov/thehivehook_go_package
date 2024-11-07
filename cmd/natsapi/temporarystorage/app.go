// - модуль реализующий хранилище временной информации
package temporarystoarge

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"
)

var (
	once sync.Once
	ts   TemporaryStorage
)

// NewTemporaryStorage конструктор временного хранилища
// ttl - time-to-live время жизни хранящейся информации в секундах,
// минимальное значение 5 секунд, максимальное не должно превышать 86400 секунд
// что соответствует 1-им суткам.
// Внимание! Чрезмерно большое время жизни временной информации может повлечь за
// собой утечку памяти.
func NewTemporaryStorage(ctx context.Context, ttl int) (*TemporaryStorage, error) {
	ts = TemporaryStorage{}

	if ttl < 5 || ttl > 86400 {
		return &ts, errors.New("the lifetime of the temporary information should not be less than 10 seconds and more than 86400 seconds")
	}

	go func(ctx context.Context, ts *TemporaryStorage) {
		<-ctx.Done()
		ts.ttlStorage.storage = make(map[string]repository)
	}(ctx, &ts)

	var err error
	once.Do(func() {
		timeToLive, newErr := time.ParseDuration(fmt.Sprintf("%ds", ttl))
		if newErr != nil {
			err = errors.Join(err, newErr)

			return
		}

		ts.ttl = timeToLive
		ts.ttlStorage = ttlStorage{
			storage: make(map[string]repository),
		}

		go checkLiveTime(ctx, &ts)
	})

	return &ts, err
}

// checkLiveTime удаляет устаревшую временную информацию
func checkLiveTime(ctx context.Context, ts *TemporaryStorage) {
	tick := time.NewTicker(5 * time.Second)
	go func(ctx context.Context, tick *time.Ticker) {
		<-ctx.Done()
		tick.Stop()
	}(ctx, tick)

	for range tick.C {
		go func() {
			ts.ttlStorage.mutex.Lock()
			defer ts.ttlStorage.mutex.Unlock()

			for k, v := range ts.ttlStorage.storage {
				if v.timeExpiry.Before(time.Now()) {
					delete(ts.ttlStorage.storage, k)
				}
			}
		}()
	}
}
