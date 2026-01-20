package natsapi

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/nats-io/nats.go"

	cint "github.com/av-belyakov/thehivehook_go_package/cmd/commoninterfaces"
	"github.com/av-belyakov/thehivehook_go_package/internal/supportingfunctions"
)

// subscriptionHandler обработчик подписки
func (api *apiNatsModule) subscriptionHandler(ctx context.Context) {
	api.natsConnection.Subscribe(api.subscriptions.listenerCommand, func(m *nats.Msg) {
		rc := RequestCommand{}
		if err := json.Unmarshal(m.Data, &rc); err != nil {
			api.logger.Send("error", supportingfunctions.CustomError(err).Error())

			return
		}

		go api.handlerIncomingCommands(ctx, rc, m)
	})
}

// handlerIncomingCommands обработчик входящих команд
func (api *apiNatsModule) handlerIncomingCommands(ctx context.Context, rc RequestCommand, m *nats.Msg) {
	id := uuid.New().String()
	chRes := make(chan cint.ChannelResponser)

	ttlTime := (time.Duration(api.cachettl) * time.Second)
	ctxTimeout, ctxTimeoutCancel := context.WithTimeout(ctx, ttlTime)
	defer func(cancel context.CancelFunc, ch chan cint.ChannelResponser) {
		cancel()
		close(ch)
	}(ctxTimeoutCancel, chRes)

	keyId := fmt.Sprintf("%s_%s", rc.RootId, rc.Command)

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

	api.storageCache.SetObject(keyId, []byte(rc.Command))

	api.sendingChannel <- &RequestFromNats{
		RequestId:  id,
		RootId:     rc.RootId,
		Service:    rc.Service,
		Command:    "send_command",
		Order:      rc.Command,
		Data:       m.Data,
		ChanOutput: chRes,
	}

	for {
		select {
		case <-ctxTimeout.Done():
			return

		case msg := <-chRes:
			errMsg := "no error"
			if err := msg.GetError(); err == nil {
				api.logger.Send("info", fmt.Sprintf("the command '%s' from service '%s' (case_id: '%s', root_id: '%s') returned a response '%d'", rc.Command, rc.Service, rc.CaseId, rc.RootId, msg.GetStatusCode()))
			} else {
				errMsg = err.Error()
			}

			//ответ на команду
			res := fmt.Appendf(nil, `{
					"id": \"%s\",
					"source": \"%s\",
					"command": \"%s\",
					"status_code": \"%d\",
					"data": %v
					"error": \"%s\",
					}`,
				msg.GetRequestId(),
				msg.GetSource(),
				rc.Command,
				msg.GetStatusCode(),
				msg.GetData(),
				errMsg,
			)
			if err := api.natsConnection.Publish(m.Reply, res); err != nil {
				api.logger.Send("error", supportingfunctions.CustomError(err).Error())
			}
			api.natsConnection.Flush()

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
			isSendCase := msg.GetCommand() != "send case"
			isSendAlert := msg.GetCommand() != "send alert"

			data, ok := msg.GetData().([]byte)
			if !ok {
				api.logger.Send("error", supportingfunctions.CustomError(errors.New("it is not possible to convert a value")).Error())

				return
			}

			//--------------------------------------------------------------
			//----------- запись в файл обработанных объектов --------------
			//--------------------------------------------------------------
			if str, err := supportingfunctions.NewReadReflectJSONSprint(data); err == nil {
				if str != "" {
					api.logger.Send("processed_objects", fmt.Sprintf("\n%s\n", str))
				}
			}
			//--------------------------------------------------------------

			if !isSendCase && !isSendAlert {
				return
			}

			var subscription, description string
			switch msg.GetElementType() {
			case "case":
				subscription = api.subscriptions.senderCase
				description = fmt.Sprintf("%s with id: '%s', rootId:'%s' has been successfully transferred", msg.GetElementType(), msg.GetCaseId(), msg.GetRootId())

			case "alert":
				subscription = api.subscriptions.senderAlert
				description = fmt.Sprintf("%s with id: '%s' has been successfully transferred", msg.GetElementType(), msg.GetRootId())

			default:
				api.logger.Send("error", supportingfunctions.CustomError(fmt.Errorf("undefined type '%s' for sending a message to NATS, cannot be processed", msg.GetElementType())).Error())

				return
			}

			if err := api.natsConnection.Publish(subscription, data); err != nil {
				api.logger.Send("error", supportingfunctions.CustomError(err).Error())
			}

			api.natsConnection.Flush()

			api.logger.Send("info", description)
		}
	}
}
