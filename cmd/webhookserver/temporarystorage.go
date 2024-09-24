package webhookserver

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
)

var (
	once sync.Once
	whts *WebHookTemporaryStorage
)

// NewWebHookTemporaryStorage конструктор временного хранилища сервера WebHook
// ttl - time-to-live время жизни хранящейся информации в секундах,
// минимальное значение 5 секунд, максимальное не должно превышать 86400 секунд
// что соответствует 1-им суткам.
// Внимание! Чрезмерно большое время жизни временной информации может повлечь за
// собой утечку памяти.
func NewWebHookTemporaryStorage(ttl int) (*WebHookTemporaryStorage, error) {
	whts := WebHookTemporaryStorage{}

	if ttl < 5 || ttl > 86400 {
		return &whts, errors.New("the lifetime of the temporary information should not be less than 10 seconds and more than 86400 seconds")
	}

	once.Do(func() {
		whts.ttl = ttl
		whts.ttlStorage = ttlStorage{
			storage: make(map[string]messageDescriptors),
		}

		go checkLiveTime(&whts)
	})

	return &whts, nil
}

func checkLiveTime(whts *WebHookTemporaryStorage) {
	for range time.Tick(5 * time.Second) {
		go func() {
			whts.ttlStorage.mutex.Lock()
			defer whts.ttlStorage.mutex.Unlock()

			for k, v := range whts.ttlStorage.storage {
				fmt.Println("int64(whts.ttl)=", int64(whts.ttl))
				fmt.Println("time.Now().Unix():", time.Now().Unix(), ">", v.timeCreate+int64(whts.ttl), "timeCreate")

				if time.Now().Unix() > (v.timeCreate + int64(whts.ttl)) {
					fmt.Println("DELETE")

					whts.DeleteElement(k)
				}
			}
		}()
	}
}

func (whts *WebHookTemporaryStorage) SetElementId(eventId string) string {
	id := uuid.New().String()

	whts.ttlStorage.mutex.Lock()
	defer whts.ttlStorage.mutex.Unlock()

	whts.ttlStorage.storage[id] = messageDescriptors{
		timeCreate: time.Now().Unix(),
		eventId:    eventId,
	}

	return id
}

func (whts *WebHookTemporaryStorage) GetElementId(id string) (string, bool) {
	if data, ok := whts.ttlStorage.storage[id]; ok {
		return data.eventId, ok
	}

	return "", false
}

func (whts *WebHookTemporaryStorage) DeleteElement(id string) {
	delete(whts.ttlStorage.storage, id)
}
