package natsapi

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
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

	rc := RequestCommand{}
	if err := json.Unmarshal(m.Data, &rc); err != nil {
		api.logger.Send("error", supportingfunctions.CustomError(err).Error())

		return
	}

	api.chanOutputModule <- OutputData{
		RequestId:  uuid.New().String(),
		Data:       m.Data,
		ChanDone:   chDone,
		ChanOutput: chRes,
	}

	/*
		api.chanOutputModule <- &RequestFromNats{
			RequestId:  uuid.New().String(),
			RootId:     rc.RootId,
			Service:    rc.Service,
			Command:    "send_command",
			Order:      rc.Command,
			Data:       m.Data,
			ChanOutput: chRes,
		}
	*/

	select {
	case <-ctx.Done():
		chDone <- struct{}{}

		return

	case msg := <-chRes:
		/*
			errMsg := "no error"
			if err := msg.GetError(); err == nil {
				api.logger.Send("info", fmt.Sprintf("the command '%s' from service '%s' (case_id: '%s', root_id: '%s') returned a response '%d'", rc.Command, rc.Service, rc.CaseId, rc.RootId, msg.GetStatusCode()))
			} else {
				errMsg = err.Error()
			}

			//ответ на команду
			res, err := json.Marshal(struct {
				Id         string `json:"id"`
				Source     string `json:"source"`
				Command    string `json:"command"`
				StatusCode int    `json:"status_code"`
				Data       any    `json:"data"`
				Error      string `json:"error"`
			}{
				Id:         msg.GetRequestId(),
				Source:     msg.GetSource(),
				Command:    rc.Command,
				StatusCode: msg.GetStatusCode(),
				Data:       msg.GetData(),
				Error:      errMsg,
			})
			if err != nil {
				api.logger.Send("error", supportingfunctions.CustomError(err).Error())

				return
			}
		*/
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
			if str, err := supportingfunctions.NewReadReflectJSONSprint(msg.Data); err == nil {
				if str != "" {
					api.logger.Send("processed_objects", fmt.Sprintf("\n%s\n", str))
				}
			}
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
