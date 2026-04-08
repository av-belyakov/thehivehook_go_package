package natsapi

import (
	"errors"

	"github.com/av-belyakov/thehivehook_go_package/internal/interfaces"
)

// New настраивает новый модуль взаимодействия с API NATS
func NewNatsApi(logger interfaces.Logger, opts ...NatsApiOptions_New) (*NatsApi, error) {
	api := &NatsApi{
		logger:   logger,
		cachettl: 10,
		//входящие в модуль данные (обогащённые alerts и case)
		chanInputModule: make(chan InputData),
		//исходящие из модуля данные (команды добавления tags и custom field)
		chanOutputModule: make(chan OutputData),
	}

	// поиск команды для объекта с определенным id поступившей за ближайшее время
	// это своего рода защитный механизм для предотвращения цикличных запросов
	//if _, ok := api.storageCache.GetObject(keyId); ok {
	//подобная команда уже есть в хранилище, исключаем её передачу
	//берём временную паузу, равную времени жизни объекта
	//----- natsapi storage -----
	//sc, err := storage.NewStorageAcceptedCommands(
	//	storage.WithMaxSize(16),
	//	storage.WithMaxTtl(180), //поставим пока время равное 3 минутам
	//	storage.WithTimeTick(2))
	//if err != nil {
	//	return api, err
	//}
	// МОЖЕТ БЫТЬ ВООБЩЕ УБРАТЬ ХРАНИЛИЩ? ТЕМ БОЛЕЕ В НАСТОЯЩИЙ МОМЕНТ ОНО НЕ ИСПОЛЬЗУЕТСЯ!!!
	//api.storageCache = sc

	for _, opt := range opts {
		if err := opt(api); err != nil {
			return api, err
		}
	}

	return api, nil
}

// WithHost метод устанавливает имя или ip адрес хоста API
func WithHost_New(v string) NatsApiOptions_New {
	return func(n *NatsApi) error {
		if v == "" {
			return errors.New("the value of 'host' cannot be empty")
		}

		n.host = v

		return nil
	}
}

// WithPort метод устанавливает порт API
func WithPort_New(v int) NatsApiOptions_New {
	return func(n *NatsApi) error {
		if v <= 0 || v > 65535 {
			return errors.New("an incorrect network port value was received")
		}

		n.port = v

		return nil
	}
}

// WithCacheTTL устанавливает время жизни для кэша хранящего функции-обработчики
// запросов к модулю
func WithCacheTTL_New(v int) NatsApiOptions_New {
	return func(th *NatsApi) error {
		if v <= 10 || v > 86400 {
			return errors.New("the lifetime of a cache entry should be between 10 and 86400 seconds")
		}

		th.cachettl = v

		return nil
	}
}

// WithSubSenderCase устанавливает канал в который будут отправлятся объекты типа 'case'
func WithSubSenderCase_New(v string) NatsApiOptions_New {
	return func(n *NatsApi) error {
		if v == "" {
			return errors.New("the value of 'sender_case' cannot be empty")
		}

		n.subscriptions.senderCase = v

		return nil
	}
}

// WithSubSenderAlert устанавливает канал в который будут отправлятся объекты типа 'alert'
func WithSubSenderAlert_New(v string) NatsApiOptions_New {
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
func WithSubListenerCommand_New(v string) NatsApiOptions_New {
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
func WithNameRegionalObject_New(v string) NatsApiOptions_New {
	return func(n *NatsApi) error {
		n.nameRegionalObject = v

		return nil
	}
}
