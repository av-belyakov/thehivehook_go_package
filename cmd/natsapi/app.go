// Модуль для взаимодействия с API NATS
package natsapi

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"runtime"
	"time"

	"github.com/google/uuid"
	"github.com/nats-io/nats.go"

	cint "github.com/av-belyakov/thehivehook_go_package/cmd/commoninterfaces"
	temporarystoarge "github.com/av-belyakov/thehivehook_go_package/cmd/natsapi/temporarystorage"
)

// New настраивает новый модуль взаимодействия с API NATS
func New(logger cint.Logger, opts ...NatsApiOptions) (*apiNatsModule, error) {
	api := &apiNatsModule{
		cachettl: 10,
		logger:   logger,
		//прием запросов в NATS
		receivingChannel: make(chan cint.ChannelRequester),
		//передача запросов из NATS
		sendingChannel: make(chan cint.ChannelRequester),
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
func (api *apiNatsModule) Start(ctx context.Context) (chan<- cint.ChannelRequester, <-chan cint.ChannelRequester, error) {
	//
	//похоже что временное хранилище для данного модуля
	// вообще не очень то и нужно, надо обдумать
	//этот момент
	//временное хранилище
	ts, err := temporarystoarge.NewTemporaryStorage(ctx, api.cachettl)
	if err != nil {
		return api.receivingChannel, api.sendingChannel, err
	}
	api.temporaryStorage = ts

	nc, err := nats.Connect(
		fmt.Sprintf("%s:%d", api.host, api.port),
		nats.MaxReconnects(-1),
		nats.ReconnectWait(3*time.Second))
	_, f, l, _ := runtime.Caller(0)
	if err != nil {
		return api.receivingChannel, api.sendingChannel, fmt.Errorf("'%w' %s:%d", err, f, l-4)
	}

	//обработка разрыва соединения с NATS
	nc.SetDisconnectErrHandler(func(c *nats.Conn, err error) {
		api.logger.Send("error", fmt.Sprintf("the connection with NATS has been disconnected (%s) %s:%d", err.Error(), f, l-4))
	})

	//обработка переподключения к NATS
	nc.SetReconnectHandler(func(c *nats.Conn) {
		api.logger.Send("info", fmt.Sprintf("the connection to NATS has been re-established (%s) %s:%d", err.Error(), f, l-4))
	})

	api.natsConnection = nc

	//обработчик подписки
	go api.subscriptionHandler(ctx)

	//обработчик данных изнутри приложения
	go api.receivingChannelHandler(ctx)

	go func(ctx context.Context, nc *nats.Conn) {
		<-ctx.Done()
		nc.Close()
	}(ctx, nc)

	return api.receivingChannel, api.sendingChannel, nil
}

// subscriptionHandler обработчик команд
func (api *apiNatsModule) subscriptionHandler(ctx context.Context) {
	api.natsConnection.Subscribe(api.subscriptions.listenerCommand, func(m *nats.Msg) {
		rc := RequestCommand{}
		if err := json.Unmarshal(m.Data, &rc); err != nil {
			_, f, l, _ := runtime.Caller(0)
			api.logger.Send("error", fmt.Sprintf("%s %s:%d", err.Error(), f, l-2))

			return
		}

		go api.handlerIncomingCommands(ctx, rc, m)
	})
}

// handlerIncomingCommands обработчик входящих, через NATS, команд
func (api *apiNatsModule) handlerIncomingCommands(ctx context.Context, rc RequestCommand, m *nats.Msg) {
	id := uuid.New().String()
	chRes := make(chan cint.ChannelResponser)

	ttlTime := (time.Duration(api.cachettl) * time.Second)
	ctxTimeout, ctxTimeoutCancel := context.WithTimeout(ctx, ttlTime)
	defer func(cancel context.CancelFunc, ch chan cint.ChannelResponser) {
		cancel()

		close(ch)
		ch = nil
	}(ctxTimeoutCancel, chRes)

	req := RequestFromNats{
		RequestId:  id,
		Command:    "send_command",
		Order:      rc.Command,
		Data:       m.Data,
		ChanOutput: chRes,
	}
	api.sendingChannel <- &req

	for {
		select {
		case <-ctxTimeout.Done():
			return

		case msg := <-chRes:
			api.logger.Send("info", fmt.Sprintf("the command '%s' from service '%s' (case_id: '%s', root_id: '%s') returned a response '%d'", rc.Command, rc.Service, rc.CaseId, rc.RootId, msg.GetStatusCode()))

			res := []byte(fmt.Sprintf("{id: \"%s\", status_code: \"%d\", data: %v}", msg.GetRequestId(), msg.GetStatusCode(), msg.GetData()))
			if err := api.natsConnection.Publish(m.Reply, res); err != nil {
				_, f, l, _ := runtime.Caller(0)
				api.logger.Send("error", fmt.Sprintf("%s %s:%d", err.Error(), f, l-2))
			}

			return
		}
	}
}

// receivingChannelHandler обработчик данных изнутри приложения
func (api *apiNatsModule) receivingChannelHandler(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return

		case msg := <-api.receivingChannel:
			data, ok := msg.GetData().([]byte)
			if !ok {
				_, f, l, _ := runtime.Caller(0)
				api.logger.Send("error", fmt.Sprintf("it is not possible to convert a value to a %s:%d", f, l-2))

				continue
			}

			var subscription string
			switch msg.GetElementType() {
			case "case":
				subscription = api.subscriptions.senderCase
			case "alert":
				subscription = api.subscriptions.senderAlert

			default:
				_, f, l, _ := runtime.Caller(0)
				api.logger.Send("error", fmt.Sprintf("undefined type '%s' for sending a message to NATS, cannot be processed %s:%d", msg.GetElementType(), f, l-6))

				continue
			}

			if err := api.natsConnection.Publish(subscription, data); err != nil {
				_, f, l, _ := runtime.Caller(0)
				api.logger.Send("error", fmt.Sprintf("%s %s:%d", err.Error(), f, l-1))
			}
		}
	}
}

// WithHost метод устанавливает имя или ip адрес хоста API
func WithHost(v string) NatsApiOptions {
	return func(n *apiNatsModule) error {
		if v == "" {
			return errors.New("the value of 'host' cannot be empty")
		}

		n.host = v

		return nil
	}
}

// WithPort метод устанавливает порт API
func WithPort(v int) NatsApiOptions {
	return func(n *apiNatsModule) error {
		if v <= 0 || v > 65535 {
			return errors.New("an incorrect network port value was received")
		}

		n.port = v

		return nil
	}
}

// WithCacheTTL устанавливает время жизни для кэша хранящего функции-обработчики
// запросов к модулю
func WithCacheTTL(v int) NatsApiOptions {
	return func(th *apiNatsModule) error {
		if v <= 10 || v > 86400 {
			return errors.New("the lifetime of a cache entry should be between 10 and 86400 seconds")
		}

		th.cachettl = v

		return nil
	}
}

// WithSubSenderCase устанавливает канал в который будут отправлятся объекты типа 'case'
func WithSubSenderCase(v string) NatsApiOptions {
	return func(n *apiNatsModule) error {
		if v == "" {
			return errors.New("the value of 'sender_case' cannot be empty")
		}

		n.subscriptions.senderCase = v

		return nil
	}
}

// WithSubSenderAlert устанавливает канал в который будут отправлятся объекты типа 'alert'
func WithSubSenderAlert(v string) NatsApiOptions {
	return func(n *apiNatsModule) error {
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
	return func(n *apiNatsModule) error {
		if v == "" {
			return errors.New("the value of 'listener_command' cannot be empty")
		}

		n.subscriptions.listenerCommand = v

		return nil
	}
}

func (mnats *ModuleNATS) GetDataReceptionChannel() <-chan SettingsOutputChan {
	return mnats.chanOutputNATS
}

func (mnats *ModuleNATS) SendingData(data SettingsOutputChan) {
	mnats.chanOutputNATS <- data
}

// WithSubscribers метод добавляет абонентов NATS
//func WithSubscribers(event string, responders []string) NatsApiOptions {
//	return func(n *apiNatsModule) error {
//		if event == "" {
//			return errors.New("the subscriber element 'event' must not be empty")
//		}
//
//		if len(responders) == 0 {
//			return errors.New("the subscriber element 'responders' must not be empty")
//		}
//
//		n.subscribers = append(n.subscribers, SubscriberNATS{
//			Event:      event,
//			Responders: responders,
//		})
//
//		return nil
//	}
//}
