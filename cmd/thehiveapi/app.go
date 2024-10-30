// Модуль для взаимодействия с API TheHive
package thehiveapi

import (
	"context"
	"errors"

	"github.com/av-belyakov/thehivehook_go_package/cmd/commoninterfaces"
	"github.com/av-belyakov/thehivehook_go_package/cmd/thehiveapi/cacherunningmethods"
)

// New настраивает модуль взаимодействия с API TheHive
func New(logger commoninterfaces.Logger, opts ...theHiveAPIOptions) (*apiTheHiveSettings, error) {
	api := &apiTheHiveSettings{
		logger:           logger,
		receivingChannel: make(chan commoninterfaces.ChannelRequester),
	}

	for _, opt := range opts {
		if err := opt(api); err != nil {
			return api, err
		}
	}

	return api, nil
}

// Start инициализирует новый модуль взаимодействия с API TheHive
// при инициализации возращается канал для взаимодействия с модулем, все
// запросы к модулю выполняются через данный канал
func (api *apiTheHiveSettings) Start(ctx context.Context) (chan<- commoninterfaces.ChannelRequester, error) {
	crm, err := cacherunningmethods.New(ctx, 30)
	if err != nil {
		return api.receivingChannel, err
	}

	api.cacheRunningMethods = crm

	go api.router(ctx)

	return api.receivingChannel, nil
}

// WithAPIKey устанавливает идентификатор-ключ для API
func WithAPIKey(v string) theHiveAPIOptions {
	return func(th *apiTheHiveSettings) error {
		if v == "" {
			return errors.New("the value of 'apiKey' cannot be empty")
		}

		th.apiKey = v

		return nil
	}
}

// WithHost устанавливает имя или ip адрес хоста API
func WithHost(v string) theHiveAPIOptions {
	return func(th *apiTheHiveSettings) error {
		if v == "" {
			return errors.New("the value of 'host' cannot be empty")
		}

		th.host = v

		return nil
	}
}

// WithPort устанавливает порт API
func WithPort(v int) theHiveAPIOptions {
	return func(th *apiTheHiveSettings) error {
		if v <= 0 || v > 65535 {
			return errors.New("an incorrect network port value was received")
		}

		th.port = v

		return nil
	}
}
