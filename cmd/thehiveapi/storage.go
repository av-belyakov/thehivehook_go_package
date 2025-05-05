package thehiveapi

import (
	"sync"
	"time"
)

type CachingStorage struct {
}

type handlerStorage struct {
	mutex    sync.RWMutex
	handlers []func(int) bool
	//количество попыток выполнения функции
	numberExecutionAttempts int
	//общее время истечения жизни, время по истечению которого объект удаляется в любом
	//случае в независимости от того, был ли он выполнен или нет, формируется time.Now().Add(c.maxTTL)
	timeExpiry time.Time
	//основное время, по нему можно найти самый старый объект в кэше
	timeMain time.Time
	//результат выполнения
	isCompletedSuccessfully bool
	//статус выполнения
	isExecution bool
}
