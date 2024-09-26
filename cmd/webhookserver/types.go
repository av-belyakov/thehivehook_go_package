package webhookserver

import (
	"context"
	"net/http"
	"sync"

	"github.com/av-belyakov/thehivehook_go_package/internal/logginghandler"
)

type WebHookServer struct {
	port    int
	host    string
	name    string //gcm, rcmmsk и т.д.
	version string
	ctx     context.Context
	server  *http.Server
	storage *WebHookTemporaryStorage
	logger  *logginghandler.LoggingChan
}

// WebHookTemporaryStorage временное хранилище для WebHookServer
// ttl - количество секунд после истечении котрых объект будет считатся
// устаревшим и подлежащим автоматическому удалению
// ttlStorage - хранилище данных со сроком жизни
type WebHookTemporaryStorage struct {
	ttl        int
	ttlStorage ttlStorage
}

type ttlStorage struct {
	mutex   sync.Mutex
	storage map[string]messageDescriptors
}

type messageDescriptors struct {
	timeCreate int64
	eventId    string
}

type EventElement struct {
	Operation  string `json:"operation"`
	ObjectType string `json:"objectType"`
	RootId     string `json:"rootId"`
}

type wrappers struct {
}
