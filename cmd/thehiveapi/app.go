// Модуль для взаимодействия с API TheHive
package thehiveapi

import (
	"context"
	"errors"

	"github.com/av-belyakov/cachingstoragewithqueue"
	"github.com/av-belyakov/thehivehook_go_package/cmd/commoninterfaces"
	"github.com/av-belyakov/thehivehook_go_package/cmd/thehiveapi/storage"
)

// New настраивает модуль взаимодействия с API TheHive
func New(logger commoninterfaces.Logger, opts ...theHiveApiOptions) (*apiTheHiveModule, error) {
	api := &apiTheHiveModule{
		cachettl:         10,
		logger:           logger,
		receivingChannel: make(chan commoninterfaces.ChannelRequester),
	}

	//---- пока уберем для тестирования использования своего собственого хранилища ----
	l := NewLogWrite(logger)
	cache, err := cachingstoragewithqueue.NewCacheStorage(
		cachingstoragewithqueue.WithMaxTtl[any](180),
		cachingstoragewithqueue.WithTimeTick[any](1),
		cachingstoragewithqueue.WithMaxSize[any](15),
		cachingstoragewithqueue.WithEnableAsyncProcessing[any](1),
		cachingstoragewithqueue.WithLogging[any](l))
	if err != nil {
		return api, err
	}
	api.cache = cache

	//----- thehiveapi storage -----
	sc, err := storage.NewStorageFoundObjects(
		storage.WithMaxSize(16),
		storage.WithMaxTtl(180),
		storage.WithTimeTick(2))
	if err != nil {
		return api, err
	}

	api.storageCache = sc

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
	//обработка кэша
	api.cache.StartAutomaticExecution(ctx)

	//инициализация автоматической очистки хранилища
	api.storageCache.Start(ctx)

	//обработка маршрутов
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
