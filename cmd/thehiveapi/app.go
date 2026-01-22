// Модуль для взаимодействия с API TheHive
package thehiveapi

import (
	"errors"

	"github.com/av-belyakov/cachingstoragewithqueue"
	"github.com/av-belyakov/thehivehook_go_package/cmd/commoninterfaces"
	"github.com/av-belyakov/thehivehook_go_package/cmd/thehiveapi/storage"
)

// New настраивает модуль взаимодействия с API TheHive
func New(logger commoninterfaces.Logger, opts ...theHiveApiOptions) (*apiTheHiveModule, error) {
	api := &apiTheHiveModule{
		settings: theHiveApiSettings{
			cachettl: 10,
		},
		logger:           logger,
		receivingChannel: make(chan commoninterfaces.ChannelRequester),
	}

	//---- пока уберем для тестирования использования своего собственого хранилища ----
	l := NewLogWrite(logger)
	cache, err := cachingstoragewithqueue.NewCacheStorage(
		cachingstoragewithqueue.WithMaxTtl[any](60),
		cachingstoragewithqueue.WithTimeTick[any](2),
		cachingstoragewithqueue.WithMaxSize[any](360),
		cachingstoragewithqueue.WithEnableAsyncProcessing[any](3),
		cachingstoragewithqueue.WithLogging[any](l))
	if err != nil {
		return api, err
	}
	api.cache = cache

	//----- thehiveapi storage -----
	sc, err := storage.NewStorageFoundObjects(
		storage.WithMaxSize(360),
		storage.WithMaxTtl(60),
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

// WithAPIKey идентификатор-ключ для API
func WithAPIKey(v string) theHiveApiOptions {
	return func(th *apiTheHiveModule) error {
		if v == "" {
			return errors.New("the value of 'apiKey' cannot be empty")
		}

		th.settings.apiKey = v

		return nil
	}
}

// WithHost имя или ip адрес хоста API
func WithHost(v string) theHiveApiOptions {
	return func(th *apiTheHiveModule) error {
		if v == "" {
			return errors.New("the value of 'host' cannot be empty")
		}

		th.settings.host = v

		return nil
	}
}

// WithPort сетевой порт API
func WithPort(v int) theHiveApiOptions {
	return func(th *apiTheHiveModule) error {
		if v <= 0 || v > 65535 {
			return errors.New("an incorrect network port value was received")
		}

		th.settings.port = v

		return nil
	}
}

// WithCacheTTL время жизни для кэша хранящего функции-обработчики
// запросов к модулю
func WithCacheTTL(v int) theHiveApiOptions {
	return func(th *apiTheHiveModule) error {
		if v <= 10 || v > 86400 {
			return errors.New("the lifetime of a cache entry should be between 10 and 86400 seconds")
		}

		th.settings.cachettl = v

		return nil
	}
}

// WithNameRegionalObject наименование регионального объекта
func WithNameRegionalObject(v string) theHiveApiOptions {
	return func(th *apiTheHiveModule) error {
		if v == "" {
			return errors.New("the value of 'nameRegionalObject' cannot be empty")
		}

		th.settings.nameRegionalObject = v

		return nil
	}
}
