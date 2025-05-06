package storage

import (
	"sync"
	"time"
)

// StorageFoundObjects хранилище найденных объектов
type StorageFoundObjects struct {
	storages storage
	maxTtl   time.Duration //максимальное время, в секундах, по истечении которого запись в cacheStorages будет удалена
	timeTick time.Duration //интервал, в секундах, с которым будут выполнятся автоматические действия
	maxSize  int
}

type storage struct {
	mutex        sync.RWMutex
	foundObjects map[string]foundObject
}

type foundObject struct {
	//ранее найденный объект
	object []byte
	//общее время истечения жизни, время по истечению которого объект удаляется в любом
	//случае в независимости от того, был ли он выполнен или нет, формируется time.Now().Add(c.maxTTL)
	timeExpiry time.Time
}

type cacheOptions func(*StorageFoundObjects) error
