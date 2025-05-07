package storage

import (
	"sync"
	"time"
)

// StorageAcceptedCommands хранилище принятых, через NATS, команд
type StorageAcceptedCommands struct {
	storages storage
	maxTtl   time.Duration //максимальное время, в секундах, по истечении которого запись в storages будет удалена
	timeTick time.Duration //интервал, в секундах, с которым будут выполнятся автоматические действия
	maxSize  int
}

type storage struct {
	mutex   sync.RWMutex
	objects map[string]Object
}

type Object struct {
	//ранее найденный объект
	Data []byte
	//общее время истечения жизни, время по истечению которого объект удаляется в любом
	//случае в независимости от того, был ли он выполнен или нет, формируется time.Now().Add(c.maxTTL)
	timeExpiry time.Time
}

type cacheOptions func(*StorageAcceptedCommands) error
