// Модуль для взаимодействия с API TheHive
package thehiveapi

import (
	"context"
	"errors"

	"github.com/av-belyakov/thehivehook_go_package/cmd/commoninterfaces"
	"github.com/av-belyakov/thehivehook_go_package/internal/cacherunningfunctions"
)

// New настраивает модуль взаимодействия с API TheHive
func New(logger commoninterfaces.Logger, opts ...theHiveApiOptions) (*apiTheHiveModule, error) {
	api := &apiTheHiveModule{
		cachettl:         10,
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
func (api *apiTheHiveModule) Start(ctx context.Context) (chan<- commoninterfaces.ChannelRequester, error) {
	crf, err := cacherunningfunctions.CreateCache(ctx, api.cachettl)
	if err != nil {
		return api.receivingChannel, err
	}

	api.cacheRunningFunction = crf

	go api.router(ctx)

	return api.receivingChannel, nil
}

// WithAPIKey устанавливает идентификатор-ключ для API
func WithAPIKey(v string) theHiveApiOptions {
	return func(th *apiTheHiveModule) error {
		if v == "" {
			return errors.New("the value of 'apiKey' cannot be empty")
		}

		th.apiKey = v

		return nil
	}
}

// WithHost устанавливает имя или ip адрес хоста API
func WithHost(v string) theHiveApiOptions {
	return func(th *apiTheHiveModule) error {
		if v == "" {
			return errors.New("the value of 'host' cannot be empty")
		}

		th.host = v

		return nil
	}
}

// WithPort устанавливает порт API
func WithPort(v int) theHiveApiOptions {
	return func(th *apiTheHiveModule) error {
		if v <= 0 || v > 65535 {
			return errors.New("an incorrect network port value was received")
		}

		th.port = v

		return nil
	}
}

// WithCacheTTL устанавливает время жизни для кэша хранящего функции-обработчики
// запросов к модулю
func WithCacheTTL(v int) theHiveApiOptions {
	return func(th *apiTheHiveModule) error {
		if v <= 10 || v > 86400 {
			return errors.New("the lifetime of a cache entry should be between 10 and 86400 seconds")
		}

		th.cachettl = v

		return nil
	}
}
