package natsapi

import (
	"context"
	"fmt"
	"time"

	"github.com/av-belyakov/thehivehook_go_package/internal/supportingfunctions"
	"github.com/nats-io/nats.go"
)

// Start инициализирует новый модуль взаимодействия с API NATS
// при инициализации возращается канал для взаимодействия с модулем, все
// запросы к модулю выполняются через данный канал
func (api *NatsApi) Start(ctx context.Context) (chan<- InputData, <-chan OutputData, error) {
	if ctx.Err() != nil {
		return api.chanInputModule, api.chanOutputModule, ctx.Err()
	}

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
		return api.chanInputModule, api.chanOutputModule, supportingfunctions.CustomError(err)
	}

	api.natsConnection = nc

	//обработчик подписки
	go api.subscriptionHandler(ctx)

	//обработчик данных изнутри приложения
	go api.receivingChannelHandler(ctx)

	go func(ctx context.Context, nc *nats.Conn) {
		<-ctx.Done()
		nc.Drain()
	}(ctx, nc)

	return api.chanInputModule, api.chanOutputModule, nil
}
