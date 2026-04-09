package natsapi

import (
	"context"
	"fmt"
	"time"

	"github.com/nats-io/nats.go"

	"github.com/av-belyakov/thehivehook_go_package/internal/supportingfunctions"
)

// subscriptionHandler обработчик подписки
func (api *NatsApi) subscriptionHandler(ctx context.Context) {
	api.natsConnection.Subscribe(api.subscriptions.listenerCommand, func(m *nats.Msg) {
		go api.handlerIncomingCommands(ctx, m)
	})
}

// handlerIncomingCommands обработчик входящих команд (добавление tags, custom fields)
func (api *NatsApi) handlerIncomingCommands(ctx context.Context, m *nats.Msg) {
	chRes := make(chan []byte)
	chDone := make(chan struct{})

	ctx, cancel := context.WithTimeout(ctx, (15 * time.Second))
	defer func() {
		cancel()

		close(chDone)
		close(chRes)
	}()

	//keyId := fmt.Sprintf("%s_%s", rc.RootId, rc.Command)
	// !!!!!!
	// Вот это прерывание цикла может серьезно мешать добавлению тегов и custom fields
	// тем более что от placeholder_misp сразу приходит две команды, одна на добавление
	// тега. вторая на добавление custom field. В итоге может быть выполнена только одна команда.
	// !!!!!!

	// поиск команды для объекта с определенным id поступившей за ближайшее время
	// это своего рода защитный механизм для предотвращения цикличных запросов
	//if _, ok := api.storageCache.GetObject(keyId); ok {
	//подобная команда уже есть в хранилище, исключаем её передачу
	//берём временную паузу, равную времени жизни объекта
	//	return
	//}
	//api.storageCache.SetObject(keyId, []byte(rc.Command))

	api.chanOutputModule <- OutputData{
		Data:       m.Data,
		ChanDone:   chDone,
		ChanOutput: chRes,
	}

	select {
	case <-ctx.Done():
		chDone <- struct{}{}

		return

	case msg := <-chRes:
		if err := api.natsConnection.Publish(m.Reply, msg); err != nil {
			api.logger.Send("error", supportingfunctions.CustomError(err).Error())
		}
	}
}

// receivingChannelHandler обработчик данных изнутри приложения (приходящих alerts и case)
func (api *NatsApi) receivingChannelHandler(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return

		case msg := <-api.chanInputModule:
			//--------------------------------------------------------------
			//----------- запись в файл обработанных объектов --------------
			//--------------------------------------------------------------
			go func() {
				if str, err := supportingfunctions.NewReadReflectJSONSprint(msg.Data); err == nil {
					if str != "" {
						api.logger.Send("processed_objects", fmt.Sprintf("\n%s\n", str))
					}
				}
			}()
			//--------------------------------------------------------------

			var subscription, description string
			switch msg.ElementType {
			case "alert":
				subscription = api.subscriptions.senderAlert
				description = fmt.Sprintf("%s with id: '%s' has been successfully transferred", msg.ElementType, msg.RootId)

			case "case":
				subscription = api.subscriptions.senderCase
				description = fmt.Sprintf("%s with id: '%s', rootId:'%s' has been successfully transferred", msg.ElementType, msg.CaseId, msg.RootId)

			default:
				api.logger.Send("error", supportingfunctions.CustomError(fmt.Errorf("undefined type '%s' for sending a message to NATS, cannot be processed", msg.ElementType)).Error())

				return
			}

			if err := api.natsConnection.Publish(subscription, msg.Data); err != nil {
				api.logger.Send("error", supportingfunctions.CustomError(err).Error())
			}

			api.natsConnection.Flush()

			api.logger.Send("info", description)
		}
	}
}
