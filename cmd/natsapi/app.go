package natsapi

import (
	"errors"

	"github.com/av-belyakov/thehivehook_go_package/internal/interfaces"
)

// New настраивает новый модуль взаимодействия с API NATS
func New(logger interfaces.Logger, opts ...NatsApiOptions) (*NatsApi, error) {
	api := &NatsApi{
		logger:   logger,
		cachettl: 10,
		//входящие в модуль данные (обогащённые alerts и case)
		chanInputModule: make(chan InputData),
		//исходящие из модуля данные (команды добавления tags и custom field)
		chanOutputModule: make(chan OutputData),
	}

	for _, opt := range opts {
		if err := opt(api); err != nil {
			return api, err
		}
	}

	return api, nil
}

// WithHost метод устанавливает имя или ip адрес хоста API
func WithHost(v string) NatsApiOptions {
	return func(n *NatsApi) error {
		if v == "" {
			return errors.New("the value of 'host' cannot be empty")
		}

		n.host = v

		return nil
	}
}

// WithPort метод устанавливает порт API
func WithPort(v int) NatsApiOptions {
	return func(n *NatsApi) error {
		if v <= 0 || v > 65535 {
			return errors.New("an incorrect network port value was received")
		}

		n.port = v

		return nil
	}
}

// WithSubSenderCase устанавливает канал в который будут отправлятся объекты типа 'case'
func WithSubSenderCase(v string) NatsApiOptions {
	return func(n *NatsApi) error {
		if v == "" {
			return errors.New("the value of 'sender_case' cannot be empty")
		}

		n.subscriptions.senderCase = v

		return nil
	}
}

// WithSubSenderAlert устанавливает канал в который будут отправлятся объекты типа 'alert'
func WithSubSenderAlert(v string) NatsApiOptions {
	return func(n *NatsApi) error {
		if v == "" {
			return errors.New("the value of 'sender_alert' cannot be empty")
		}

		n.subscriptions.senderAlert = v

		return nil
	}
}

// WithSubListenerCommand устанавливает канал через которые будут приходить команды для
// выполнения определенных действий в TheHive
func WithSubListenerCommand(v string) NatsApiOptions {
	return func(n *NatsApi) error {
		if v == "" {
			return errors.New("the value of 'listener_command' cannot be empty")
		}

		n.subscriptions.listenerCommand = v

		return nil
	}
}

// WithNameRegionalObject устанавливает наименование которое будет отображатся в
// статистике подключенных клиентов NATS
func WithNameRegionalObject(v string) NatsApiOptions {
	return func(n *NatsApi) error {
		n.nameRegionalObject = v

		return nil
	}
}
