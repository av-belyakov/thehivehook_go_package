package webhookserver

import (
	"context"
	"net/http"

	"github.com/av-belyakov/thehivehook_go_package/cmd/commoninterfaces"
	temporarystorage "github.com/av-belyakov/thehivehook_go_package/cmd/webhookserver/temporarystorage"
	"github.com/av-belyakov/thehivehook_go_package/internal/logginghandler"
)

// WebHookServer непосредственно сам сервер
type WebHookServer struct {
	ttl       int
	port      int
	host      string
	name      string //gcm, rcmmsk и т.д.
	version   string
	ctx       context.Context
	server    *http.Server
	storage   *temporarystorage.WebHookTemporaryStorage
	logger    *logginghandler.LoggingChan
	chanInput chan<- ChanFormWebHookServer
}

// webHookServerOptions функциональные параметры
type webHookServerOptions func(*WebHookServer)

// WebHookServerOptions основные опции
type WebHookServerOptions struct {
	TTL     int
	Port    int
	Host    string
	Name    string
	Version string
}

// ChanFormWebHookServer структура канала для взаимодействия сторонних сервисов с webhookserver
type ChanFormWebHookServer struct {
	ForSomebody string
	Data        commoninterfaces.ChannelRequester
}

// EventElement типовой элемент описывающий события приходящие из TheHive
type EventElement struct {
	Operation  string `json:"operation"`
	ObjectType string `json:"objectType"`
	RootId     string `json:"rootId"`
}

// ReadyMadeEventCase готовый, сформированный объект содержащий информацию по кейсу
type ReadyMadeEventCase struct {
	Source      string                 `json:"source"`
	Case        map[string]interface{} `json:"event"`
	Observables []interface{}          `json:"observables"`
	TTPs        []interface{}          `json:"ttp"`
}

// ReadyMadeEventAlert готовый, сформированный объект содержащий информацию по алерту
type ReadyMadeEventAlert struct {
	Source string                 `json:"source"`
	Event  map[string]interface{} `json:"event"`
	Alert  map[string]interface{} `json:"alert"`
}
