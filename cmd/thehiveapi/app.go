// Модуль для взаимодействия с API TheHive
package thehiveapi

import (
	"errors"
)

// New настраивает модуль взаимодействия с API TheHive
func New(opts ...theHiveApiOptions) (*TheHiveApi, error) {
	api := &TheHiveApi{
		settings: theHiveApiSettings{},
	}

	for _, opt := range opts {
		if err := opt(api); err != nil {
			return api, err
		}
	}

	return api, nil
}

// WithAPIKey идентификатор-ключ для API
func WithAPIKey(v string) theHiveApiOptions {
	return func(th *TheHiveApi) error {
		if v == "" {
			return errors.New("the value of 'apiKey' cannot be empty")
		}

		th.settings.apiKey = v

		return nil
	}
}

// WithHost имя или ip адрес хоста API
func WithHost(v string) theHiveApiOptions {
	return func(th *TheHiveApi) error {
		if v == "" {
			return errors.New("the value of 'host' cannot be empty")
		}

		th.settings.host = v

		return nil
	}
}

// WithPort сетевой порт API
func WithPort(v int) theHiveApiOptions {
	return func(th *TheHiveApi) error {
		if v <= 0 || v > 65535 {
			return errors.New("an incorrect network port value was received")
		}

		th.settings.port = v

		return nil
	}
}

// WithNameRegionalObject наименование регионального объекта
func WithNameRegionalObject(v string) theHiveApiOptions {
	return func(th *TheHiveApi) error {
		if v == "" {
			return errors.New("the value of 'nameRegionalObject' cannot be empty")
		}

		th.settings.nameRegionalObject = v

		return nil
	}
}
