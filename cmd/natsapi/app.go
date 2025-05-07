// Модуль для взаимодействия с API NATS
package natsapi

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/nats-io/nats.go"

	cint "github.com/av-belyakov/thehivehook_go_package/cmd/commoninterfaces"
	"github.com/av-belyakov/thehivehook_go_package/cmd/constants"
	"github.com/av-belyakov/thehivehook_go_package/cmd/natsapi/storage"
	"github.com/av-belyakov/thehivehook_go_package/internal/supportingfunctions"
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

	//----- natsapi storage -----
	sc, err := storage.NewStorageAcceptedCommands(
		storage.WithMaxSize(16),
		storage.WithMaxTtl(180), //поставим пока время равное 3 минутам
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

// Start инициализирует новый модуль взаимодействия с API NATS
// при инициализации возращается канал для взаимодействия с модулем, все
// запросы к модулю выполняются через данный канал
func (api *apiNatsModule) Start(ctx context.Context) (chan<- cint.ChannelRequester, <-chan cint.ChannelRequester, error) {
	if ctx.Err() != nil {
		return api.receivingChannel, api.sendingChannel, ctx.Err()
	}

	//инициализация автоматической очистки хранилища используемого для хранения
	//принимаемых, через NATS, команд
	//сделал для того что бы избежать повторной отправки одной и той же команды
	//предназначенной для одного и того же объекта передаваемой за короткий
	//промежуток времени
	api.storageCache.Start(ctx)

	nc, err := nats.Connect(
		fmt.Sprintf("%s:%d", api.host, api.port),
		//имя клиента
		nats.Name(fmt.Sprintf("thehivehook.%s", api.nameRegionalObject)),
		//неограниченное количество попыток переподключения
		nats.MaxReconnects(-1),
		//время ожидания до следующей попытки переподключения (по умолчанию 2 сек.)
		nats.ReconnectWait(3*time.Second),
		//обработка разрыва соединения с NATS
		nats.DisconnectErrHandler(func(c *nats.Conn, err error) {
			api.logger.Send("error", supportingfunctions.CustomError(fmt.Errorf("the connection with NATS has been disconnected (%w)", err)).Error())
		}),
		//обработка переподключения к NATS
		nats.ReconnectHandler(func(c *nats.Conn) {
			api.logger.Send("info", "the connection to NATS has been re-established")
		}),
		//поиск медленных получателей (не обязательный для данного приложения параметр)
		nats.ErrorHandler(func(c *nats.Conn, s *nats.Subscription, err error) {
			if err == nats.ErrSlowConsumer {
				pendingMsgs, _, err := s.Pending()
				if err != nil {
					api.logger.Send("warning", fmt.Sprintf("couldn't get pending messages: %v", err))

					return
				}

				api.logger.Send("warning", fmt.Sprintf("Falling behind with %d pending messages on subject %q.\n", pendingMsgs, s.Subject))
			}
		}))
	if err != nil {
		return api.receivingChannel, api.sendingChannel, supportingfunctions.CustomError(err)
	}

	log.Printf("%vconnect to NATS with address %v%s:%d%v\n", constants.Ansi_Bright_Green, constants.Ansi_Dark_Gray, api.host, api.port, constants.Ansi_Reset)

	api.natsConnection = nc

	//обработчик подписки
	go api.subscriptionHandler(ctx)

	//обработчик данных изнутри приложения
	go api.receivingChannelHandler(ctx)

	go func(ctx context.Context, nc *nats.Conn) {
		<-ctx.Done()
		nc.Drain()
	}(ctx, nc)

	return api.receivingChannel, api.sendingChannel, nil
}
