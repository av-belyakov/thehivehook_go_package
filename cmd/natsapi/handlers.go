package natsapi

import (
	"context"
	"encoding/json"
	"fmt"
	"runtime"
	"time"

	"github.com/google/uuid"
	"github.com/nats-io/nats.go"

	cint "github.com/av-belyakov/thehivehook_go_package/cmd/commoninterfaces"
)

// subscriptionHandler обработчик подписки
func (api *apiNatsModule) subscriptionHandler(ctx context.Context) {
	api.natsConnection.Subscribe(api.subscriptions.listenerCommand, func(m *nats.Msg) {

		fmt.Println("((( subscriptionHandler ))) reseived ", string(m.Data))

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
					id: \"%s\", 
					error: \"%s\",
					command: \"%s\", 
					status_code: \"%d\", 
					data: %v}`,
					msg.GetRequestId(),
					msg.GetError().Error(),
					rc.Command,
					msg.GetStatusCode(),
					msg.GetData()))
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
			isSendCase := msg.GetCommand() != "send case"
			isSendAlert := msg.GetCommand() != "send alert"

			if !isSendCase && !isSendAlert {
				continue
			}

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
