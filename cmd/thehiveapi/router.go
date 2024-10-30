package thehiveapi

import (
	"context"
	"fmt"
)

func (api *apiTheHiveSettings) router(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return

		case msg := <-api.receivingChannel:
			switch msg.GetCommand() {
			case "get_observables":
				res, statusCode, err := api.GetObservables(ctx, msg.GetRootId())
				if err != nil {
					api.logger.Send("error", err.Error())

					continue
				}

				newRes := NewChannelRespons()
				newRes.SetRequestId(msg.GetRequestId())
				newRes.SetStatusCode(statusCode)
				newRes.SetData(res)

				msg.GetChanOutput() <- newRes
				close(msg.GetChanOutput())

			case "get_ttp":
				res, statusCode, err := api.GetTTP(ctx, msg.GetRootId())
				if err != nil {
					api.logger.Send("error", err.Error())

					continue
				}

				newRes := NewChannelRespons()
				newRes.SetRequestId(msg.GetRequestId())
				newRes.SetStatusCode(statusCode)
				newRes.SetData(res)

				msg.GetChanOutput() <- newRes
				close(msg.GetChanOutput())

			case "send command":
				// Вот здесь нужно использовать temporaryStorage как кеширующее
				// хранилище команд корторые нужно отправить в TheHive и которые
				// из-за, по какой то причине, недоступности TheHive отправить сразу
				// не получается

				switch msg.GetOrder() {
				case "add case tags":
					_, statusCode, err := api.AddCaseTags(ctx, msg.GetRootId(), msg.GetData())
					if err != nil {
						api.logger.Send("error", err.Error())

						continue
					}

					api.logger.Send("info", fmt.Sprintf("when making a request to add a new tag for the rootId '%s', the following is received status code '%d'", msg.GetRootId(), statusCode))

				case "add case custom fields":

				case "add case task":
				}
			}
		}
	}
}
