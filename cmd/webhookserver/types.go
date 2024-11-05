package webhookserver

import (
	"net/http"

	"github.com/av-belyakov/thehivehook_go_package/cmd/commoninterfaces"
)

// WebHookServer непосредственно сам сервер
type WebHookServer struct {
	ttl        int
	port       int
	host       string
	name       string //gcm, rcmmsk и т.д.
	version    string
	pathSqlite string
	server     *http.Server
	logger     commoninterfaces.Logger
	chanInput  chan<- ChanFromWebHookServer
}

// webHookServerOptions функциональные параметры
type webHookServerOptions func(*WebHookServer) error

// WebHookServerOptions основные опции
type WebHookServerOptions struct {
	TTL     int    //Time to live в секундах, по умолчанию 10 сек. (не обязательный параметр)
	Port    int    //сетевой порт, по умолчанию 7575 (не обязательный параметр)
	Host    string //сетевой хост, доменное имя или ip адрес, по умолчанию 127.0.0.1
	Name    string //наименование Webhook сервера (не обязательный параметр)
	Version string //версия сервера (не обязательный параметр)
}

// ChanFromWebHookServer структура канала для взаимодействия сторонних сервисов с webhookserver
type ChanFromWebHookServer struct {
	ForSomebody string //от кого данные
	Data        commoninterfaces.ChannelRequester
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
