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
			case "get_alert":
				keyId := msg.GetCommand() + msg.GetRootId()

				api.logger.Send("info", fmt.Sprintf("--- Search alert request accepted, command:'%s', root id:'%s' keyId:'%s'", msg.GetCommand(), msg.GetRootId(), keyId))

				so := NewSpecialObjectForCache[any]()
				so.SetID(keyId)

				so.SetFunc(func(_ int) bool {
					api.logger.Send("info", fmt.Sprintf("start search object for alert, command:'%s', root id:'%s'", msg.GetCommand(), msg.GetRootId()))

					newRes := NewChannelRespons()
					newRes.SetRequestId(msg.GetRequestId())
					newRes.SetStatusCode(200)

					//ищем в хранилище объект, который возможно уже запрашивали ранее
					if obj, ok := api.storageCache.GetObject(keyId); ok {
						api.logger.Send("info", fmt.Sprintf("the object with id:'%s' was found in the cache", keyId))

						newRes.SetData(obj)
					} else {
						api.logger.Send("info", fmt.Sprintf("request object for alert to TheHive, command:'%s', root id:'%s'", msg.GetCommand(), msg.GetRootId()))

						//делаем запрос к TheHive для получения дополнительной информации по объекту
						res, statusCode, err := api.GetAlert(ctx, msg.GetRootId())
						if err != nil {
							api.logger.Send("error", supportingfunctions.CustomError(err).Error())
						} else {
							api.logger.Send("info", fmt.Sprintf("successful response to TheHive request, command:'%s', root id:'%s', status code:'%d'", msg.GetCommand(), msg.GetRootId(), statusCode))

							// добавляем найденный объект в кеш
							api.storageCache.SetObject(keyId, res)
							newRes.SetData(res)
						}

						newRes.SetStatusCode(statusCode)
					}

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

			case "get_observables":
				keyId := msg.GetCommand() + msg.GetRootId()

				//api.logger.Send("info", fmt.Sprintf("--- Search request accepted, command:'%s', root id:'%s' (case:'%s') keyId:'%s'", msg.GetCommand(), msg.GetRootId(), msg.GetCaseId(), keyId))

				so := NewSpecialObjectForCache[any]()
				so.SetID(keyId)

				//для того что бы выполнить сравнение объектов нужно передать
				//этот объект so.SetObject
				//хотя для thehivehook_go может это не надо, надо обдумать!!!
				//so.SetObject(msg.GetData())

				so.SetFunc(func(_ int) bool {
					api.logger.Send("info", fmt.Sprintf("start search object, command:'%s', root id:'%s' (case:'%s')", msg.GetCommand(), msg.GetRootId(), msg.GetCaseId()))

					newRes := NewChannelRespons()
					newRes.SetRequestId(msg.GetRequestId())
					newRes.SetStatusCode(200)

					//ищем в хранилище объект, который возможно уже запрашивали ранее
					if obj, ok := api.storageCache.GetObject(keyId); ok {
						api.logger.Send("info", fmt.Sprintf("the object with id:'%s' was found in the cache", keyId))

						newRes.SetData(obj)
					} else {
						api.logger.Send("info", fmt.Sprintf("request to TheHive, command:'%s', root id:'%s' (case:'%s')", msg.GetCommand(), msg.GetRootId(), msg.GetCaseId()))

						//делаем запрос к TheHive для получения дополнительной информации по объекту
						res, statusCode, err := api.GetObservables(ctx, msg.GetRootId())
						if err != nil {
							api.logger.Send("error", supportingfunctions.CustomError(err).Error())
						} else {
							api.logger.Send("info", fmt.Sprintf("successful response to TheHive request, command:'%s', root id:'%s', status code:'%d'", msg.GetCommand(), msg.GetRootId(), statusCode))

							// добавляем найденный объект в кеш
							api.storageCache.SetObject(keyId, res)
							newRes.SetData(res)
						}

						newRes.SetStatusCode(statusCode)
					}

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

			case "get_ttp":
				keyId := msg.GetCommand() + msg.GetRootId()

				//api.logger.Send("info", fmt.Sprintf("--- Search request accepted, command:'%s', root id:'%s' (case:'%s') keyId:'%s'", msg.GetCommand(), msg.GetRootId(), msg.GetCaseId(), keyId))

				so := NewSpecialObjectForCache[any]()
				so.SetID(keyId)

				//для того что бы выполнить сравнение объектов нужно передать
				//этот объект so.SetObject
				//хотя для thehivehook_go может это не надо, надо обдумать!!!
				//so.SetObject(msg.GetData())

				so.SetFunc(func(_ int) bool {
					api.logger.Send("info", fmt.Sprintf("start search object, command:'%s', root id:'%s' (case:'%s')", msg.GetCommand(), msg.GetRootId(), msg.GetCaseId()))

					newRes := NewChannelRespons()
					newRes.SetRequestId(msg.GetRequestId())
					newRes.SetStatusCode(200)

					//ищем в хранилище объект, который возможно уже запрашивали ранее
					if obj, ok := api.storageCache.GetObject(keyId); ok {
						api.logger.Send("info", fmt.Sprintf("the object with id:'%s' was found in the cache", keyId))

						newRes.SetData(obj)
					} else {
						api.logger.Send("info", fmt.Sprintf("request to TheHive, command:'%s', root id:'%s' (case:'%s')", msg.GetCommand(), msg.GetRootId(), msg.GetCaseId()))

						//делаем запрос к TheHive для получения дополнительной информации по объекту
						res, statusCode, err := api.GetTTP(ctx, msg.GetRootId())
						if err != nil {
							api.logger.Send("error", supportingfunctions.CustomError(err).Error())
						} else {
							api.logger.Send("info", fmt.Sprintf("successful response to TheHive request, command:'%s', root id:'%s', status code:'%d'", msg.GetCommand(), msg.GetRootId(), statusCode))

							// добавляем найденный объект в кеш
							api.storageCache.SetObject(keyId, res)
							newRes.SetData(res)
						}

						newRes.SetStatusCode(statusCode)
					}

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

			case "send_command":
				rc, err := getRequestCommandData(msg.GetData())
				if err != nil {
					api.logger.Send("error", supportingfunctions.CustomError(err).Error())

					continue
				}

				api.logger.Send("info", fmt.Sprintf("the command '%s' has been received, order:'%s', rootId:'%s'", rc.Command, msg.GetOrder(), msg.GetRootId()))

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

func getRequestCommandData(i any) (RequestCommand, error) {
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
