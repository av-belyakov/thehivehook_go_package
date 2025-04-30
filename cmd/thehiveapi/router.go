package thehiveapi

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/av-belyakov/thehivehook_go_package/internal/datamodels"
	"github.com/av-belyakov/thehivehook_go_package/internal/supportingfunctions"
)

func (api *apiTheHiveModule) router(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return

		case msg := <-api.receivingChannel:
			//switch msg.GetCommand() {
			switch msg.Command {
			case "get_observables":
				go func(command, id string) {
					so := NewSpecialObjectForCache[any]()
					so.SetID(command + id)

					//для того что бы выполнить сравнение объектов нужно передать
					//этот объект so.SetObject
					//хотя для thehivehook_go может это не надо, надо обдумать!!!
					//so.SetObject(msg.GetData())

					so.SetFunc(func(_ int) bool {
						api.logger.Send("info", fmt.Sprintf("request to TheHive, command:'%s', root id:'%d' (case:'%s')", command, id, msg.CaseId))

						res, statusCode, err := api.GetObservables(ctx, id)
						if err != nil {
							api.logger.Send("error", supportingfunctions.CustomError(err).Error())

							return false
						}

						//newRes := NewChannelRespons()
						//newRes.SetRequestId(id)
						//newRes.SetStatusCode(statusCode)
						//newRes.SetData(res)

						api.logger.Send("info", fmt.Sprintf("successful response to TheHive request, command:'%s', root id:'%s', status code:'%d'", command, id, statusCode))

						select {
						//case <-msg.GetContext().Done():
						case <-msg.Context.Done():
							return false

						default:
							//msg.GetChanOutput() <- newRes
							msg.ChOutput <- datamodels.ResponseChan{
								RequestId:  id,
								StatusCode: statusCode,
								Data:       res,
							}

						}

						return true
					})

					//добавляем объект в очередь для обработки
					api.cache.PushObjectToQueue(so)
					//}(msg.GetCommand(), msg.GetRootId())
				}(msg.Command, msg.RootId)

			case "get_ttp":
				go func(command, id string) {
					so := NewSpecialObjectForCache[any]()
					so.SetID(command + id)
					so.SetFunc(func(_ int) bool {
						//api.logger.Send("info", fmt.Sprintf("request to TheHive, command:'%s', root id:'%s' (case:'%s')", command, id, msg.GetCaseId()))
						api.logger.Send("info", fmt.Sprintf("request to TheHive, command:'%s', root id:'%s' (case:'%d')", command, id, msg.CaseId))

						res, statusCode, err := api.GetTTP(ctx, id)
						if err != nil {
							api.logger.Send("error", supportingfunctions.CustomError(err).Error())

							return false
						}

						//newRes := NewChannelRespons()
						//newRes.SetRequestId(id)
						//newRes.SetStatusCode(statusCode)
						//newRes.SetData(res)

						api.logger.Send("info", fmt.Sprintf("successful response to TheHive request, command:'%s', root id:'%s', status code:'%d'", command, id, statusCode))

						select {
						//case <-msg.GetContext().Done():
						case <-msg.Context.Done():
							return false

						default:
							//msg.GetChanOutput() <- newRes
							msg.ChOutput <- datamodels.ResponseChan{
								RequestId:  id,
								StatusCode: statusCode,
								Data:       res,
							}

						}

						return true
					})

					//добавляем объект в очередь для обработки
					api.cache.PushObjectToQueue(so)
					//}(msg.GetCommand(), msg.GetRootId())
				}(msg.Command, msg.RootId)

			case "send_command":
				//rc, err := getRequestCommandData(msg.GetData())
				rc, err := getRequestCommandData(msg.Data)
				if err != nil {
					api.logger.Send("error", supportingfunctions.CustomError(err).Error())

					continue
				}

				//api.logger.Send("info", fmt.Sprintf("the command '%s' has been received, order:'%s', rootId:'%s'", rc.Command, msg.GetOrder(), msg.GetRootId()))
				api.logger.Send("info", fmt.Sprintf("the command '%s' has been received, order:'%s', rootId:'%s'", rc.Command, msg.Order, msg.RootId))

				//switch msg.GetOrder() {
				switch msg.Order {
				case "add_case_tag":
					go func(id string) {
						//res := NewChannelRespons()
						//res.SetRequestId(id)

						res := datamodels.ResponseChan{RequestId: id}

						_, statusCode, err := api.AddCaseTags(ctx, rc)
						//res.SetStatusCode(statusCode)
						if err != nil {
							res.StatusCode = statusCode
							res.Error = err
							//res.SetError(err)
							api.logger.Send("error", supportingfunctions.CustomError(err).Error())
						} else {
							api.logger.Send("info", fmt.Sprintf("when making a request to add a new tag for the rootId '%s', caseId '%s', the following is received status code '%d'", rc.RootId, rc.CaseId, statusCode))
						}

						res.StatusCode = statusCode

						//res.SendToChan(msg.GetChanOutput())
						msg.ChOutput <- res
						//}(msg.GetRequestId())
					}(msg.RequestId)

				case "add_case_task":
					go func(id string) {
						//res := NewChannelRespons()
						//res.SetRequestId(id)

						res := datamodels.ResponseChan{RequestId: id}
						_, statusCode, err := api.AddCaseTask(ctx, rc)
						//res.SetStatusCode(statusCode)
						if err != nil {
							res.StatusCode = statusCode
							res.Error = err
							//res.SetError(err)
							api.logger.Send("error", supportingfunctions.CustomError(err).Error())
						} else {
							api.logger.Send("info", fmt.Sprintf("when making a request to add a new tag for the rootId '%s', caseId '%s', the following is received status code '%d'", rc.RootId, rc.CaseId, statusCode))
						}

						res.StatusCode = statusCode
						//res.SendToChan(msg.ChanOutput())
						msg.ChOutput <- res
						//}(msg.GetRequestId())
					}(msg.RequestId)

				case "set_case_custom_field":

					go func(id string) {
						//res := NewChannelRespons()
						//res.SetRequestId(id)

						res := datamodels.ResponseChan{RequestId: id}

						_, statusCode, err := api.AddCaseCustomFields(ctx, rc)
						//res.SetStatusCode(statusCode)
						if err != nil {
							res.StatusCode = statusCode
							res.Error = err
							//res.SetError(err)
							api.logger.Send("error", supportingfunctions.CustomError(err).Error())
						} else {
							api.logger.Send("info", fmt.Sprintf("when making a request to add a new tag for the rootId '%s', caseId '%s', the following is received status code '%d'", rc.RootId, rc.CaseId, statusCode))
						}

						res.StatusCode = statusCode
						//res.SendToChan(msg.ChanOutput())
						msg.ChOutput <- res
						//}(msg.GetRequestId())
					}(msg.RequestId)
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
