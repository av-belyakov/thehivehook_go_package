// Модуль для взаимодействия с API NATS
package natsapi

import (
	"context"
	"errors"

	"github.com/av-belyakov/thehivehook_go_package/cmd/commoninterfaces"
	"github.com/av-belyakov/thehivehook_go_package/internal/logginghandler"
)

// New инициализирует новый модуль взаимодействия с API NATS
// при инициализации возращается канал для взаимодействия с модулем, все
// запросы к модулю выполняются через данный канал
func New(ctx context.Context, logging *logginghandler.LoggingChan, opts ...NatsAPIOptions) (chan<- commoninterfaces.ChannelRequester, error) {
	receivingChannel := make(chan commoninterfaces.ChannelRequester)

	api := &apiNatsSettings{
		subscribers: []SubscriberNATS(nil),
	}

	for _, opt := range opts {
		if err := opt(api); err != nil {
			return receivingChannel, err
		}
	}

	go func() {

	}()

	return receivingChannel, nil
}

// WithHost метод устанавливает имя или ip адрес хоста API
func WithHost(v string) NatsAPIOptions {
	return func(n *apiNatsSettings) error {
		if v == "" {
			return errors.New("the value of 'host' cannot be empty")
		}

		n.host = v

		return nil
	}
}

// WithPort метод устанавливает порт API
func WithPort(v int) NatsAPIOptions {
	return func(n *apiNatsSettings) error {
		if v <= 0 || v > 65535 {
			return errors.New("an incorrect network port value was received")
		}

		n.port = v

		return nil
	}
}

// WithSubscribers метод добавляет абонентов NATS
func WithSubscribers(event string, responders []string) NatsAPIOptions {
	return func(n *apiNatsSettings) error {
		if event == "" {
			return errors.New("the subscriber element 'event' must not be empty")
		}

		if len(responders) == 0 {
			return errors.New("the subscriber element 'responders' must not be empty")
		}

		n.subscribers = append(n.subscribers, SubscriberNATS{
			Event:      event,
			Responders: responders,
		})

		return nil
	}
}

func (mnats *ModuleNATS) GetDataReceptionChannel() <-chan SettingsOutputChan {
	return mnats.chanOutputNATS
}

func (mnats *ModuleNATS) SendingData(data SettingsOutputChan) {
	mnats.chanOutputNATS <- data
}
