package cacherunningmethods

import (
	"sync"
	"time"
)

type storageParameters struct {
	timeExpiry  time.Time
	cacheMethod func() bool
}

type cacheStorageParameters struct {
	mutex    sync.RWMutex
	storages map[string]storageParameters
}

type CacheRunningMethods struct {
	ttl          time.Duration
	cacheStorage cacheStorageParameters
}
