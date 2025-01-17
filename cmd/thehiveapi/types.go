package thehiveapi

import (
	"github.com/av-belyakov/thehivehook_go_package/cmd/commoninterfaces"
)

// apiTheHiveModule модуль для взаимодействия с API TheHive
type apiTheHiveModule struct {
	cachettl             int
	port                 int
	host                 string
	apiKey               string
	logger               commoninterfaces.Logger
	receivingChannel     chan commoninterfaces.ChannelRequester
	cacheRunningFunction commoninterfaces.CacheFuncRunner
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
