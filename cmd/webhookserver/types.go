package webhookserver

import (
	"net/http"
	"time"

	"github.com/av-belyakov/thehivehook_go_package/internal/interfaces"
)

// WebHookServer непосредственно сам сервер
type WebHookServer struct {
	server    *http.Server
	logger    interfaces.Logger
	timeStart time.Time
	version   string
	name      string //gcm, rcmmsk и т.д.
	host      string
	port      int
	chanOut   chan []byte
}

// webHookServerOptions функциональные параметры
type webHookServerOptions func(*WebHookServer) error

// WebHookServerOptions основные опции
type WebHookServerOptions struct {
	Host    string //сетевой хост, доменное имя или ip адрес, по умолчанию 127.0.0.1
	Name    string //наименование Webhook сервера (не обязательный параметр)
	Version string //версия сервера (не обязательный параметр)
	TTL     int    //Time to live в секундах, по умолчанию 10 сек. (не обязательный параметр)
	Port    int    //сетевой порт, по умолчанию 7575 (не обязательный параметр)
}

// ChanFromWebHookServer структура канала для взаимодействия сторонних сервисов с webhookserver
type ChanFromWebHookServer struct {
	Data        interfaces.ChannelRequester
	ForSomebody string //для кого данные
}

// ReadyMadeEventCase готовый, сформированный объект содержащий информацию по кейсу
type ReadyMadeEventCase struct {
	Case        map[string]any `json:"event"`
	Observables []any          `json:"observables"`
	TTPs        []any          `json:"ttp"`
	Source      string         `json:"source"`
}

// ReadyMadeEventAlert готовый, сформированный объект содержащий информацию по алерту
type ReadyMadeEventAlert struct {
	Event  map[string]any `json:"event"`
	Alert  map[string]any `json:"alert"`
	Source string         `json:"source"`
}

type EventElements map[string]any
