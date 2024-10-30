package thehiveapi

import (
	"context"
	"fmt"
	"time"
)

func (api *apiTheHiveSettings) router(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return

		case msg := <-api.receivingChannel:
			switch msg.GetCommand() {
			case "get_observables":
				api.cacheRunningMethods.SetMethod(msg.GetRootId(), func() bool {
					ctxTimeout, ctxClose := context.WithTimeout(ctx, 5*time.Second)
					defer ctxClose()

					res, statusCode, err := api.GetObservables(ctxTimeout, msg.GetRootId())
					if err != nil {
						api.logger.Send("error", err.Error())

						return false
					}

					newRes := NewChannelRespons()
					newRes.SetRequestId(msg.GetRequestId())
					newRes.SetStatusCode(statusCode)
					newRes.SetData(res)

					msg.GetChanOutput() <- newRes
					close(msg.GetChanOutput())

					return true
				})

			case "get_ttp":
				api.cacheRunningMethods.SetMethod(msg.GetRootId(), func() bool {
					ctxTimeout, ctxClose := context.WithTimeout(ctx, 5*time.Second)
					defer ctxClose()

					res, statusCode, err := api.GetTTP(ctxTimeout, msg.GetRootId())
					if err != nil {
						api.logger.Send("error", err.Error())

						return false
					}

					newRes := NewChannelRespons()
					newRes.SetRequestId(msg.GetRequestId())
					newRes.SetStatusCode(statusCode)
					newRes.SetData(res)

					msg.GetChanOutput() <- newRes
					close(msg.GetChanOutput())

					return true
				})

			case "send command":
				switch msg.GetOrder() {
				case "add case tags":
					api.cacheRunningMethods.SetMethod(msg.GetRootId(), func() bool {
						ctxTimeout, ctxClose := context.WithTimeout(ctx, 5*time.Second)
						defer ctxClose()

						_, statusCode, err := api.AddCaseTags(ctxTimeout, msg.GetRootId(), msg.GetData())
						if err != nil {
							api.logger.Send("error", err.Error())

							return false
						}

						api.logger.Send("info", fmt.Sprintf("when making a request to add a new tag for the rootId '%s', the following is received status code '%d'", msg.GetRootId(), statusCode))

						return true
					})
				case "add case custom fields":

				case "add case task":
				}
			}
		}
	}
}
