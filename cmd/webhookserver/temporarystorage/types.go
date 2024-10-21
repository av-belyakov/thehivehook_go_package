package temporarystoarge

import (
	"sync"
	"time"
)

// WebHookTemporaryStorage временное хранилище для WebHookServer
type WebHookTemporaryStorage struct {
	ttl        time.Duration //количество секунд после истечении котрых объект будет считатся устаревшим и подлежащим автоматическому удалению
	ttlStorage ttlStorage    //хранилище данных со сроком жизни
}

// ttlStorage хранилище данных со сроком жизни
type ttlStorage struct {
	mutex   sync.RWMutex
	storage map[string]messageDescriptors
}

// messageDescriptors структура с описанием хранящихся значений
type messageDescriptors struct {
	timeExpiry time.Time
	value      string
}
