package thehiveapi

import (
	"context"
	"fmt"
)

func (api *apiTheHiveModule) router(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return

		case msg := <-api.receivingChannel:
			switch msg.GetCommand() {
			case "get_observables":
				api.cacheRunningFunction.SetMethod(msg.GetRootId(), func() bool {
					res, statusCode, err := api.GetObservables(ctx, msg.GetRootId())
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
				api.cacheRunningFunction.SetMethod(msg.GetRootId(), func() bool {
					res, statusCode, err := api.GetTTP(ctx, msg.GetRootId())
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
					api.cacheRunningFunction.SetMethod(msg.GetRootId(), func() bool {
						_, statusCode, err := api.AddCaseTags(ctx, msg.GetRootId(), msg.GetData())
						if err != nil {
							api.logger.Send("error", err.Error())

							return false
						}

						api.logger.Send("info", fmt.Sprintf("when making a request to add a new tag for the rootId '%s', caseId '%s', the following is received status code '%d'", msg.GetRootId(), msg.GetCaseId(), statusCode))

						return true
					})
				case "add case custom fields":
					api.cacheRunningFunction.SetMethod(msg.GetRootId(), func() bool {
						_, statusCode, err := api.AddCaseCustomFields(ctx, msg.GetRootId(), msg.GetData())
						if err != nil {
							api.logger.Send("error", err.Error())

							return false
						}

						api.logger.Send("info", fmt.Sprintf("when making a request to add a new tag for the rootId '%s', caseId '%s', the following is received status code '%d'", msg.GetRootId(), msg.GetCaseId(), statusCode))

						return true
					})

				case "add case task":
					api.cacheRunningFunction.SetMethod(msg.GetRootId(), func() bool {
						_, statusCode, err := api.AddCaseTask(ctx, msg.GetRootId(), msg.GetData())
						if err != nil {
							api.logger.Send("error", err.Error())

							return false
						}

						api.logger.Send("info", fmt.Sprintf("when making a request to add a new tag for the rootId '%s', caseId '%s', the following is received status code '%d'", msg.GetRootId(), msg.GetCaseId(), statusCode))

						return true
					})

				case "set severity":

				}
			}
		}
	}
}
