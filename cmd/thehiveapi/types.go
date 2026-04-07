package thehiveapi

import (
	"github.com/av-belyakov/cachingstoragewithqueue"
	"github.com/av-belyakov/thehivehook_go_package/cmd/thehiveapi/storage"
	"github.com/av-belyakov/thehivehook_go_package/internal/interfaces"
)

type TheHiveApi_New struct {
	settings theHiveApiSettings
}

// apiTheHiveModule модуль для взаимодействия с API TheHive
type apiTheHiveModule struct {
	cache            *cachingstoragewithqueue.CacheStorageWithQueue[any]
	storageCache     *storage.StorageFoundObjects
	logger           interfaces.Logger
	receivingChannel chan interfaces.ChannelRequester
	settings         theHiveApiSettings
}

// theHiveApiSettings настройки модуля
type theHiveApiSettings struct {
	apiKey             string
	host               string
	nameRegionalObject string
	cachettl           int
	port               int
}

// theHiveAPIOptions функциональные опции
type theHiveApiOptions func(*apiTheHiveModule) error

// theHiveAPIOptions_New функциональные опции
type theHiveApiOptions_New func(*TheHiveApi_New) error

type InputReguest struct {
	Data    []byte
	Command string
	RootId  string
}

// Querys перечень запросов к TheHive
type Querys struct {
	Query []Query `json:"query"`
}

// Query структура запроса к TheHive
type Query struct {
	ExtraData []string `json:"extraData,omitempty"`
	Name      string   `json:"_name,omitempty"`
	IDOrName  string   `json:"idOrName,omitempty"`
	From      int64    `json:"from"`
	To        int      `json:"to,omitempty"`
}

// ErrorAnswer структура описания ошибок получаемых от TheHive
type ErrorAnswer struct {
	Err     string `json:"type"`
	Message string `json:"message"`
}
