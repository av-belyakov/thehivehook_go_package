// Модуль для взаимодействия с API TheHive
package thehiveapi

import (
	"errors"
)

// New настраивает модуль взаимодействия с API TheHive
func NewTheHiveApi(opts ...theHiveApiOptions_New) (*TheHiveApi_New, error) {
	api := &TheHiveApi_New{
		settings: theHiveApiSettings{
			cachettl: 10,
		},
	}

	for _, opt := range opts {
		if err := opt(api); err != nil {
			return api, err
		}
	}

	return api, nil
}

// WithAPIKey_New идентификатор-ключ для API
func WithAPIKey_New(v string) theHiveApiOptions_New {
	return func(th *TheHiveApi_New) error {
		if v == "" {
			return errors.New("the value of 'apiKey' cannot be empty")
		}

		th.settings.apiKey = v

		return nil
	}
}

// WithHost_New имя или ip адрес хоста API
func WithHost_New(v string) theHiveApiOptions_New {
	return func(th *TheHiveApi_New) error {
		if v == "" {
			return errors.New("the value of 'host' cannot be empty")
		}

		th.settings.host = v

		return nil
	}
}

// WithPort_New сетевой порт API
func WithPort_New(v int) theHiveApiOptions_New {
	return func(th *TheHiveApi_New) error {
		if v <= 0 || v > 65535 {
			return errors.New("an incorrect network port value was received")
		}

		th.settings.port = v

		return nil
	}
}

// WithNameRegionalObject_New наименование регионального объекта
func WithNameRegionalObject_New(v string) theHiveApiOptions_New {
	return func(th *TheHiveApi_New) error {
		if v == "" {
			return errors.New("the value of 'nameRegionalObject' cannot be empty")
		}

		th.settings.nameRegionalObject = v

		return nil
	}
}
