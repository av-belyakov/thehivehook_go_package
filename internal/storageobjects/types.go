package storageobjects

import (
	"sync"
	"time"
)

// storageObjects хранилище объектов
type StorageObjects[T any] struct {
	timeTick        time.Duration //интервал, в секундах, с которым будут выполнятся автоматические действия
	timeToLive      time.Duration //время, в секундах, по истечении которого запись в storages будет удалена
	timeDelayToSend time.Duration //время, в секундах, по истечении которого объект из записи будет передан в канал,
	// а сама запись будет удалена
	storage     objects[T]                //хранилище объектов
	chanOutSize int                       //размер канала для передачи данных
	chanOut     chan StorageObjectData[T] //канал для передачи данных
}

type objects[T any] struct {
	mtx     sync.RWMutex
	objects []object[T]
}

type object[T any] struct {
	timeSending time.Time //время отправки объекта
	timeCreated time.Time //время добавления объекта
	timeExpiry  time.Time //общее время истечения жизни, время по истечению которого объект удаляется
	element     T
	id          string
}

type cacheOptions[T any] func(*StorageObjects[T]) error

type StorageObjectData[T any] struct {
	TimeCreated time.Time
	Data        T
	Id          string
}
