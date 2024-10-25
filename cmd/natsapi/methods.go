// Модуль для взаимодействия с API NATS
package natsapi

import (
	"context"
	"errors"

	"github.com/av-belyakov/thehivehook_go_package/cmd/commoninterfaces"
	temporarystoarge "github.com/av-belyakov/thehivehook_go_package/cmd/natsapi/temporarystorage"
	"github.com/av-belyakov/thehivehook_go_package/internal/logginghandler"
)

// New настраивает новый модуль взаимодействия с API NATS
func New(logger *logginghandler.LoggingChan, opts ...NatsAPIOptions) (*apiNatsSettings, error) {
	ts, err := temporarystoarge.NewTemporaryStorage(30)
	if err != nil {
		return &apiNatsSettings{}, err
	}

	api := &apiNatsSettings{
		subscribers:      []SubscriberNATS(nil),
		logger:           logger,
		receivingChannel: make(chan commoninterfaces.ChannelRequester),
		temporaryStorage: ts,
	}

	for _, opt := range opts {
		if err := opt(api); err != nil {
			return api, err
		}
	}

	return api, nil
}

// Start инициализирует новый модуль взаимодействия с API NATS
// при инициализации возращается канал для взаимодействия с модулем, все
// запросы к модулю выполняются через данный канал
func (api *apiNatsSettings) Start(ctx context.Context) chan<- commoninterfaces.ChannelRequester {
	go func() {
		//здесь temporarystorage будет использоватся для хранения двух
		// основных типов данных:
		// 1. хранение дескрипторов соединения с NATS
		// 2. выполнение функции кеширования case или alert которые отправляются
		// в NATS. Если NATS по какой то причине не будет доступен, то хранить
		// вышеуказанные виды объектов и пытатся их отправить до тех пор
		// пока они не будут отправлены или не истечет заданный срок после которых
		// их можно будет удалить
	}()

	return api.receivingChannel
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
