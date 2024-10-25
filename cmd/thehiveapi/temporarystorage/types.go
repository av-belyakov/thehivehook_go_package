package temporarystoarge

import (
	"sync"
	"time"
)

// TemporaryStorage временное хранилище
type TemporaryStorage struct {
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
	timeExpiry time.Time //время добавления значения
	value      []byte    //значение, в данном случае это команда к TheHive
}
