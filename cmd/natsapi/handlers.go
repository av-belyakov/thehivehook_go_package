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

	api.logger.Send("info", fmt.Sprintf("the command '%s' has been received, data '%v'", rc.Command, m.Data))

	api.sendingChannel <- &RequestFromNats{
		RequestId:  id,
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
			api.logger.Send("info", fmt.Sprintf("the command '%s' from service '%s' (case_id: '%s', root_id: '%s') returned a response '%d'", rc.Command, rc.Service, rc.CaseId, rc.RootId, msg.GetStatusCode()))

			res := []byte(
				fmt.Sprintf(`{
					"id": \"%s\", 
					"error": \"%s\",
					"command": \"%s\", 
					"status_code": \"%d\", 
					"data": %v
					}`,
					msg.GetRequestId(),
					msg.GetError().Error(),
					rc.Command,
					msg.GetStatusCode(),
					msg.GetData()))
			if err := api.natsConnection.Publish(m.Reply, res); err != nil {
				api.logger.Send("error", supportingfunctions.CustomError(err).Error())
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
			isSendCase := msg.GetCommand() != "send case"
			isSendAlert := msg.GetCommand() != "send alert"

			if !isSendCase && !isSendAlert {
				continue
			}

			data, ok := msg.GetData().([]byte)
			if !ok {
				api.logger.Send("error", supportingfunctions.CustomError(errors.New("it is not possible to convert a value")).Error())

				continue
			}

			var subscription, description string
			switch msg.GetElementType() {
			case "case":
				subscription = api.subscriptions.senderCase
				description = fmt.Sprintf("%s with id: '%s' has been successfully transferred", msg.GetElementType(), msg.GetCaseId())

			case "alert":
				subscription = api.subscriptions.senderAlert
				description = fmt.Sprintf("%s with id: '%s' has been successfully transferred", msg.GetElementType(), msg.GetRootId())

			default:
				api.logger.Send("error", supportingfunctions.CustomError(fmt.Errorf("undefined type '%s' for sending a message to NATS, cannot be processed", msg.GetElementType())).Error())

				continue
			}

			if err := api.natsConnection.Publish(subscription, data); err != nil {
				api.logger.Send("error", supportingfunctions.CustomError(err).Error())
			}

			api.logger.Send("log_to_db", description)
		}
	}
}
