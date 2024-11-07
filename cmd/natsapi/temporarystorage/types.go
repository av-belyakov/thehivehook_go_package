package temporarystoarge

import (
	"sync"
	"time"

	"github.com/nats-io/nats.go"
)

// TemporaryStorage временное хранилище
type TemporaryStorage struct {
	ttl        time.Duration //количество секунд после истечении котрых объект будет считатся устаревшим и подлежащим автоматическому удалению
	ttlStorage ttlStorage    //хранилище данных со сроком жизни
}

// ttlStorage хранилище данных со сроком жизни
type ttlStorage struct {
	mutex   sync.RWMutex
	storage map[string]repository
}

type repository struct {
	timeExpiry time.Time
	service    string
	command    string
	rootId     string
	caseId     string
	nsMsg      *nats.Msg
}
