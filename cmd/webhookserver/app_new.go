package webhookserver

import (
	"time"

	"github.com/av-belyakov/thehivehook_go_package/internal/interfaces"
)

// NewWebHookServer конструктор webhookserver принимает функциональные опции для настройки модуля перед запуском
func NewWebHookServer(logger interfaces.Logger, opts ...webHookServerOptions_New) (*WebHookServer_New, <-chan OutputData, error) {
	chanOutput := make(chan OutputData)

	whs := &WebHookServer_New{
		name:       "nobody",
		version:    "0.1.1",
		timeStart:  time.Now(),
		host:       "127.0.0.1",
		port:       7575,
		logger:     logger,
		chanOutput: chanOutput,
	}

	for _, opt := range opts {
		if err := opt(whs); err != nil {
			return whs, chanOutput, err
		}
	}

	return whs, chanOutput, nil
}

// WithPort_New устанавливает порт для взаимодействия с модулем
func WithPort_New(v int) webHookServerOptions_New {
	return func(whs *WebHookServer_New) error {
		whs.port = v

		return nil
	}
}

// WithHost_New устанавливает хост для взаимодействия с модулем
func WithHost_New(v string) webHookServerOptions_New {
	return func(whs *WebHookServer_New) error {
		whs.host = v

		return nil
	}
}

// WithName_New устанавливает наименование модуля (обязательно). Наименование основывается
// на имени организации или подразделения эксплуатирующем модуль. Например, gcm, rcmslx и т.д.
func WithName_New(v string) webHookServerOptions_New {
	return func(whs *WebHookServer_New) error {
		whs.name = v

		return nil
	}
}

// WithVersion_New устанавливает версию модуля (опционально)
func WithVersion_New(v string) webHookServerOptions_New {
	return func(whs *WebHookServer_New) error {
		whs.version = v

		return nil
	}
}
