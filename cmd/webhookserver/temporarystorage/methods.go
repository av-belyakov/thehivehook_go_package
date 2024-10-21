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
	whts WebHookTemporaryStorage
)

// NewWebHookTemporaryStorage конструктор временного хранилища сервера WebHook
// ttl - time-to-live время жизни хранящейся информации в секундах,
// минимальное значение 5 секунд, максимальное не должно превышать 86400 секунд
// что соответствует 1-им суткам.
// Внимание! Чрезмерно большое время жизни временной информации может повлечь за
// собой утечку памяти.
func NewWebHookTemporaryStorage(ttl int) (*WebHookTemporaryStorage, error) {
	whts = WebHookTemporaryStorage{}

	if ttl < 5 || ttl > 86400 {
		return &whts, errors.New("the lifetime of the temporary information should not be less than 10 seconds and more than 86400 seconds")
	}

	var err error
	once.Do(func() {
		timeToLive, newErr := time.ParseDuration(fmt.Sprintf("%ds", ttl))
		if newErr != nil {
			err = newErr

			return
		}

		whts.ttl = timeToLive
		whts.ttlStorage = ttlStorage{
			storage: make(map[string]messageDescriptors),
		}

		go checkLiveTime(&whts)
	})

	return &whts, err
}

// checkLiveTime удаляет устаревшую временную информацию
func checkLiveTime(whts *WebHookTemporaryStorage) {
	for range time.Tick(5 * time.Second) {
		go func() {
			whts.ttlStorage.mutex.Lock()
			defer whts.ttlStorage.mutex.Unlock()

			for k, v := range whts.ttlStorage.storage {
				if v.timeExpiry.Before(time.Now()) {
					delete(whts.ttlStorage.storage, k)
				}
			}
		}()
	}
}

// SetValue создает новую запись, принимает значение которое нужно сохранить и
// id по которому данное значение можно будет найти
func (whts *WebHookTemporaryStorage) SetValue(id, value string) string {
	whts.ttlStorage.mutex.Lock()
	defer whts.ttlStorage.mutex.Unlock()

	whts.ttlStorage.storage[id] = messageDescriptors{
		timeExpiry: time.Now().Add(whts.ttl),
		value:      value,
	}

	return id
}

// GetValue возвращает данные по полученому id
func (whts *WebHookTemporaryStorage) GetValue(id string) (string, bool) {
	if data, ok := whts.ttlStorage.storage[id]; ok {
		return data.value, ok
	}

	return "", false
}

// DeleteElement удаляет заданный элемент по его uuid
func (whts *WebHookTemporaryStorage) DeleteElement(id string) {
	whts.ttlStorage.mutex.Lock()
	defer whts.ttlStorage.mutex.Unlock()

	delete(whts.ttlStorage.storage, id)
}
