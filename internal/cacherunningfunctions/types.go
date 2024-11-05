// Хранилище выполняющее кэширование некоторых функции для
// их последующего выполнения. Функция будет удалена из кэша
// при ее успешном выполнении или по истечении определенного времени.
package cacherunningfunctions

import (
	"sync"
	"time"
)

type storageParameters struct {
	timeExpiry time.Time
	cacheFunc  func() bool
}

type cacheStorageParameters struct {
	mutex    sync.RWMutex
	storages map[string]storageParameters
}

// CacheRunningFunctions хранилище функций
type CacheRunningFunctions struct {
	ttl          time.Duration
	cacheStorage cacheStorageParameters
}
