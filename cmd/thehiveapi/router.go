package thehiveapi

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/av-belyakov/thehivehook_go_package/internal/supportingfunctions"
)

func (api *apiTheHiveModule) router(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return

		case msg := <-api.receivingChannel:
			switch msg.GetCommand() {
			case "get_observables":
				go func(command, id string) {
					so := NewSpecialObjectForCache[interface{}]()
					so.SetID(command + id)
					so.SetFunc(func(_ int) bool {

						fmt.Printf("=== func 'router', command:'%s', root id:'%s' (%s)\n", command, id, command+id)

						res, statusCode, err := api.GetObservables(ctx, id)
						if err != nil {
							api.logger.Send("error", supportingfunctions.CustomError(err).Error())

							return false
						}

						newRes := NewChannelRespons()
						newRes.SetRequestId(id)
						newRes.SetStatusCode(statusCode)
						newRes.SetData(res)

						select {
						case <-msg.GetContext().Done():
							return false

						default:
							msg.GetChanOutput() <- newRes

						}

						return true
					})

					//добавляем объект в очередь для обработки
					api.cache.PushObjectToQueue(so)
				}(msg.GetCommand(), msg.GetRootId())

			case "get_ttp":
				go func(command, id string) {
					so := NewSpecialObjectForCache[interface{}]()
					so.SetID(command + id)
					so.SetFunc(func(_ int) bool {

						fmt.Printf("=== func 'router', command:'%s', root id:'%s' (%s)\n", command, id, command+id)

						res, statusCode, err := api.GetTTP(ctx, id)
						if err != nil {
							api.logger.Send("error", supportingfunctions.CustomError(err).Error())

							return false
						}

						newRes := NewChannelRespons()
						newRes.SetRequestId(id)
						newRes.SetStatusCode(statusCode)
						newRes.SetData(res)

						select {
						case <-msg.GetContext().Done():
							return false

						default:
							msg.GetChanOutput() <- newRes

						}

						return true
					})

					//добавляем объект в очередь для обработки
					api.cache.PushObjectToQueue(so)
				}(msg.GetCommand(), msg.GetRootId())

			case "send_command":
				rc, err := getRequestCommandData(msg.GetData())
				if err != nil {
					api.logger.Send("error", supportingfunctions.CustomError(err).Error())

					continue
				}

				switch msg.GetOrder() {
				case "add_case_tag":
					go func(id string) {
						res := NewChannelRespons()
						res.SetRequestId(id)

						_, statusCode, err := api.AddCaseTags(ctx, rc)
						res.SetStatusCode(statusCode)
						if err != nil {
							res.SetError(err)
							api.logger.Send("error", supportingfunctions.CustomError(err).Error())
						} else {
							api.logger.Send("info", fmt.Sprintf("when making a request to add a new tag for the rootId '%s', caseId '%s', the following is received status code '%d'", rc.RootId, rc.CaseId, statusCode))
						}

						res.SendToChan(msg.GetChanOutput())
					}(msg.GetRequestId())

				case "add_case_task":
					go func(id string) {
						res := NewChannelRespons()
						res.SetRequestId(id)

						_, statusCode, err := api.AddCaseTask(ctx, rc)
						res.SetStatusCode(statusCode)
						if err != nil {
							res.SetError(err)
							api.logger.Send("error", supportingfunctions.CustomError(err).Error())
						} else {
							api.logger.Send("info", fmt.Sprintf("when making a request to add a new tag for the rootId '%s', caseId '%s', the following is received status code '%d'", rc.RootId, rc.CaseId, statusCode))
						}

						res.SendToChan(msg.GetChanOutput())
					}(msg.GetRequestId())

				case "set_case_custom_field":
					go func(id string) {
						res := NewChannelRespons()
						res.SetRequestId(id)

						_, statusCode, err := api.AddCaseCustomFields(ctx, rc)
						res.SetStatusCode(statusCode)
						if err != nil {
							res.SetError(err)
							api.logger.Send("error", supportingfunctions.CustomError(err).Error())
						} else {
							api.logger.Send("info", fmt.Sprintf("when making a request to add a new tag for the rootId '%s', caseId '%s', the following is received status code '%d'", rc.RootId, rc.CaseId, statusCode))
						}

						res.SendToChan(msg.GetChanOutput())
					}(msg.GetRequestId())
				}
			}
		}
	}
}

func getRequestCommandData(i interface{}) (RequestCommand, error) {
	rc := RequestCommand{}

	data, ok := i.([]byte)
	if !ok {
		return rc, errors.New("'it is not possible to convert a value msg.GetData() to a []byte'")
	}

	if err := json.Unmarshal(data, &rc); err != nil {
		return rc, err
	}

	return rc, nil
}
