// Основной модуль приложения
package webhookserver

import (
	"time"

	"github.com/av-belyakov/thehivehook_go_package/internal/interfaces"
	"github.com/av-belyakov/thehivehook_go_package/internal/storageobjects"
)

// New конструктор webhookserver принимает функциональные опции для настройки модуля перед запуском
func New[T any](logger interfaces.Logger, opts ...webHookServerOptions[T]) (*WebHookServer[T], <-chan ChanFromWebHookServer, error) {
	chanOutput := make(chan ChanFromWebHookServer)

	whs := &WebHookServer[T]{
		name:         "nobody",
		version:      "0.1.1",
		timeStart:    time.Now(),
		host:         "127.0.0.1",
		port:         7575,
		ttl:          30,
		delaySending: 10,
		logger:       logger,
		chanInput:    chanOutput,
	}

	for _, opt := range opts {
		if err := opt(whs); err != nil {
			return whs, chanOutput, err
		}
	}

	storage, err := storageobjects.New(
		storageobjects.WithTimeTick[T](1),
		storageobjects.WithChannelSize[T](10),
		storageobjects.WithTimeToLive[T](whs.ttl),
		storageobjects.WithTimeDelayToSend[T](whs.delaySending),
	)
	if err != nil {
		return whs, chanOutput, err
	}

	whs.storage = storage

	return whs, chanOutput, nil
}

// WithStorageDelayToSend устанавливает время задержки отправки данных из хранилища
func WithStorageDelayToSend[T any](v int) webHookServerOptions[T] {
	return func(whs *WebHookServer[T]) error {
		whs.delaySending = v

		return nil
	}
}

// WithTTL устанавливает время TimeToLive для временного хранилища информации в модуле
func WithTTL[T any](v int) webHookServerOptions[T] {
	return func(whs *WebHookServer[T]) error {
		whs.ttl = v

		return nil
	}
}

// WithPort устанавливает порт для взаимодействия с модулем
func WithPort[T any](v int) webHookServerOptions[T] {
	return func(whs *WebHookServer[T]) error {
		whs.port = v

		return nil
	}
}

// WithHost устанавливает хост для взаимодействия с модулем
func WithHost[T any](v string) webHookServerOptions[T] {
	return func(whs *WebHookServer[T]) error {
		whs.host = v

		return nil
	}
}

// WithName устанавливает наименование модуля (обязательно). Наименование основывается
// на имени организации или подразделения эксплуатирующем модуль. Например, gcm, rcmslx и т.д.
func WithName[T any](v string) webHookServerOptions[T] {
	return func(whs *WebHookServer[T]) error {
		whs.name = v

		return nil
	}
}

// WithVersion устанавливает версию модуля (опционально)
func WithVersion[T any](v string) webHookServerOptions[T] {
	return func(whs *WebHookServer[T]) error {
		whs.version = v

		return nil
	}
}
