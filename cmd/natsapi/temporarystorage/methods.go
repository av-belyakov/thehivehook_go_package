// - модуль реализующий хранилище временной информации
package temporarystoarge

import (
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
func NewTemporaryStorage(ttl int) (*TemporaryStorage, error) {
	ts = TemporaryStorage{}

	if ttl < 5 || ttl > 86400 {
		return &ts, errors.New("the lifetime of the temporary information should not be less than 10 seconds and more than 86400 seconds")
	}

	var err error
	once.Do(func() {
		timeToLive, newErr := time.ParseDuration(fmt.Sprintf("%ds", ttl))
		if newErr != nil {
			err = newErr

			return
		}

		ts.ttl = timeToLive
		ts.ttlStorage = ttlStorage{
			storage: make(map[string]messageDescriptors),
		}

		go checkLiveTime(&ts)
	})

	return &ts, err
}

// checkLiveTime удаляет устаревшую временную информацию
func checkLiveTime(ts *TemporaryStorage) {
	for range time.Tick(5 * time.Second) {
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

// SetValue создает новую запись, принимает значение которое нужно сохранить и
// id по которому данное значение можно будет найти
func (ts *TemporaryStorage) SetValue(id, value string) string {
	ts.ttlStorage.mutex.Lock()
	defer ts.ttlStorage.mutex.Unlock()

	ts.ttlStorage.storage[id] = messageDescriptors{
		timeExpiry: time.Now().Add(ts.ttl),
		value:      value,
	}

	return id
}

// GetValue возвращает данные по полученому id
func (ts *TemporaryStorage) GetValue(id string) (string, bool) {
	if data, ok := ts.ttlStorage.storage[id]; ok {
		return data.value, ok
	}

	return "", false
}

// DeleteElement удаляет заданный элемент по его uuid
func (ts *TemporaryStorage) DeleteElement(id string) {
	ts.ttlStorage.mutex.Lock()
	defer ts.ttlStorage.mutex.Unlock()

	delete(ts.ttlStorage.storage, id)
}
