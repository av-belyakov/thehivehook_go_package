// Основной модуль приложения
package webhookserver

import (
	"time"

	"github.com/av-belyakov/thehivehook_go_package/internal/interfaces"
)

// New конструктор webhookserver принимает функциональные опции для настройки модуля перед запуском
func New(logger interfaces.Logger, opts ...webHookServerOptions) (*WebHookServer, <-chan []byte, error) {
	whs := &WebHookServer{
		name:      "nobody",
		version:   "0.1.1",
		timeStart: time.Now(),
		host:      "127.0.0.1",
		port:      7575,
		logger:    logger,
		chanOut:   make(chan []byte),
	}

	for _, opt := range opts {
		if err := opt(whs); err != nil {
			return whs, whs.chanOut, err
		}
	}

	return whs, whs.chanOut, nil
}

/*// New конструктор webhookserver принимает функциональные опции для настройки модуля перед запуском
func New(logger interfaces.Logger, opts ...webHookServerOptions) (*WebHookServer, <-chan ChanFromWebHookServer, error) {
	chanOutput := make(chan ChanFromWebHookServer)

	whs := &WebHookServer{
		name:      "nobody",
		version:   "0.1.1",
		timeStart: time.Now(),
		host:      "127.0.0.1",
		port:      7575,
		logger:    logger,
		chanInput: chanOutput,
	}

	for _, opt := range opts {
		if err := opt(whs); err != nil {
			return whs, chanOutput, err
		}
	}

	return whs, chanOutput, nil
}*/

// WithPort устанавливает порт для взаимодействия с модулем
func WithPort(v int) webHookServerOptions {
	return func(whs *WebHookServer) error {
		whs.port = v

		return nil
	}
}

// WithHost устанавливает хост для взаимодействия с модулем
func WithHost(v string) webHookServerOptions {
	return func(whs *WebHookServer) error {
		whs.host = v

		return nil
	}
}

// WithName устанавливает наименование модуля (обязательно). Наименование основывается
// на имени организации или подразделения эксплуатирующем модуль. Например, gcm, rcmslx и т.д.
func WithName(v string) webHookServerOptions {
	return func(whs *WebHookServer) error {
		whs.name = v

		return nil
	}
}

// WithVersion устанавливает версию модуля (опционально)
func WithVersion(v string) webHookServerOptions {
	return func(whs *WebHookServer) error {
		whs.version = v

		return nil
	}
}
