package thehiveapi

import (
	"github.com/av-belyakov/cachingstoragewithqueue"
	"github.com/av-belyakov/thehivehook_go_package/cmd/commoninterfaces"
)

// apiTheHiveModule модуль для взаимодействия с API TheHive
type apiTheHiveModule struct {
	cache            *cachingstoragewithqueue.CacheStorageWithQueue[interface{}]
	logger           commoninterfaces.Logger
	apiKey           string
	host             string
	receivingChannel chan commoninterfaces.ChannelRequester
	cachettl         int
	port             int
}

// theHiveAPIOptions функциональные опции
type theHiveApiOptions func(*apiTheHiveModule) error

// Querys перечень запросов к TheHive
type Querys struct {
	Query []Query `json:"query"`
}

// Query структура запроса к TheHive
type Query struct {
	Name      string   `json:"_name,omitempty"`
	IDOrName  string   `json:"idOrName,omitempty"`
	From      int64    `json:"from"`
	To        int      `json:"to,omitempty"`
	ExtraData []string `json:"extraData,omitempty"`
}

// ErrorAnswer структура описания ошибок получаемых от TheHive
type ErrorAnswer struct {
	Err     string `json:"type"`
	Message string `json:"message"`
}

// SpecialObjectForCache является вспомогательным типом который реализует интерфейс
// CacheStorageFuncHandler[T any] где в методе Comparison(objFromCache T) bool необходимо
// реализовать подробное сравнение объекта типа T.
// Нужен для пакета cachingstoragewithqueue
type SpecialObjectForCache[T any] struct {
	object      T
	handlerFunc func(int) bool
	id          string
}
