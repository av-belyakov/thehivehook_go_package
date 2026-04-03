package storageobjects

import (
	"sync"
	"time"
)

// storageObjects хранилище объектов
type StorageObjects struct {
	timeTick        time.Duration //интервал, в секундах, с которым будут выполнятся автоматические действия
	timeToLive      time.Duration //время, в секундах, по истечении которого запись в storages будет удалена
	timeDelayToSend time.Duration //время, в секундах, по истечении которого объект из записи будет передан в канал,
	// а сама запись будет удалена
	storage     objects                //хранилище объектов
	chanOutSize int                    //размер канала для передачи данных
	chanOut     chan StorageObjectData //канал для передачи данных
}

type objects struct {
	mtx     sync.RWMutex
	objects []object
}

type object struct {
	timeSending time.Time //время отправки объекта
	timeCreated time.Time //время добавления объекта
	timeExpiry  time.Time //общее время истечения жизни, время по истечению которого объект удаляется
	data        []byte
	objectType  string
	id          string
}

type cacheOptions func(*StorageObjects) error

type StorageObjectData struct {
	TimeCreated time.Time
	Data        []byte
	ObjectType  string
	Id          string
}
