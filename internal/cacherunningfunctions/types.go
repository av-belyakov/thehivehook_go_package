// Хранилище выполняющее кэширование некоторых функции для
// их последующего выполнения. Функция будет удалена из кэша
// при ее успешном выполнении или по истечении определенного времени.
package cacherunningfunctions

import (
	"sync"
	"time"
)

// CacheRunningFunctions хранилище функций
type CacheRunningFunctions struct {
	ttl          time.Duration
	cacheStorage cacheStorageParameters
}

type cacheStorageParameters struct {
	mutex    sync.RWMutex
	storages map[string]storageParameters
}

type storageParameters struct {
	cacheFunc               func(int) bool
	timeExpiry              time.Time
	numberAttempts          int
	isFunctionExecution     bool
	isCompletedSuccessfully bool
}
