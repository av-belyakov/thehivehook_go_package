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
				fmt.Printf("=== thehiveapi router, command:'%s' rootId:'%s'\n", msg.GetCommand(), msg.GetRequestId())

				go api.cacheRunningFunction.SetMethod("observables"+msg.GetRootId(), func(_ int) bool {
					defer close(msg.GetChanOutput())

					fmt.Println("zzz Method GET OBSERVABLE")

					res, statusCode, err := api.GetObservables(ctx, msg.GetRootId())
					if err != nil {
						api.logger.Send("error", supportingfunctions.CustomError(err).Error())

						return false
					}

					newRes := NewChannelRespons()
					newRes.SetRequestId(msg.GetRootId())
					newRes.SetStatusCode(statusCode)
					newRes.SetData(res)

					msg.GetChanOutput() <- newRes

					//что бы данную гроутину не держала ссылка на объект
					newRes = NewChannelRespons()

					return true
				})

			case "get_ttp":
				fmt.Printf("=== thehiveapi router, command:'%s' rootId:'%s'\n", msg.GetCommand(), msg.GetRequestId())

				go api.cacheRunningFunction.SetMethod("ttp"+msg.GetRootId(), func(_ int) bool {
					defer close(msg.GetChanOutput())

					fmt.Println("zzz Method GET TTP")

					res, statusCode, err := api.GetTTP(ctx, msg.GetRootId())
					if err != nil {
						api.logger.Send("error", supportingfunctions.CustomError(err).Error())

						return false
					}

					newRes := NewChannelRespons()
					newRes.SetRequestId(msg.GetRootId())
					newRes.SetStatusCode(statusCode)
					newRes.SetData(res)

					msg.GetChanOutput() <- newRes

					//что бы данную гроутину не держала ссылка на объект
					newRes = NewChannelRespons()

					return true
				})

			case "send_command":
				rc, err := getRequestCommandData(msg.GetData())
				if err != nil {
					api.logger.Send("error", supportingfunctions.CustomError(err).Error())
				}

				chRes := msg.GetChanOutput()

				res := NewChannelRespons()
				res.SetRequestId(msg.GetRequestId())

				switch msg.GetOrder() {
				case "add_case_tag":
					api.cacheRunningFunction.SetMethod(msg.GetRequestId(), func(countAttempts int) bool {
						_, statusCode, err := api.AddCaseTags(ctx, rc)

						if err != nil {
							api.logger.Send("error", supportingfunctions.CustomError(err).Error())

							if countAttempts >= 10 {
								res.SetStatusCode(statusCode)
								res.SetError(err)
								res.sendToChan(chRes)

								return true
							}

							return false
						}

						api.logger.Send("info", fmt.Sprintf("when making a request to add a new tag for the rootId '%s', caseId '%s', the following is received status code '%d'", rc.RootId, rc.CaseId, statusCode))

						res.SetStatusCode(statusCode)
						res.sendToChan(chRes)

						return true
					})

				case "add_case_task":
					api.cacheRunningFunction.SetMethod(msg.GetRequestId(), func(countAttempts int) bool {
						_, statusCode, err := api.AddCaseTask(ctx, rc)

						if err != nil {
							api.logger.Send("error", supportingfunctions.CustomError(err).Error())

							if countAttempts >= 10 {
								res.SetStatusCode(statusCode)
								res.SetError(err)
								res.sendToChan(chRes)

								return true
							}

							return false
						}

						api.logger.Send("info", fmt.Sprintf("when making a request to add a new task for the rootId '%s', caseId '%s', the following is received status code '%d'", rc.RootId, rc.CaseId, statusCode))

						res.SetStatusCode(statusCode)
						res.sendToChan(chRes)

						return true
					})

				case "set_case_custom_field":
					api.cacheRunningFunction.SetMethod(msg.GetRequestId(), func(countAttempts int) bool {
						_, statusCode, err := api.AddCaseCustomFields(ctx, rc)

						if err != nil {
							api.logger.Send("error", supportingfunctions.CustomError(err).Error())

							if countAttempts >= 10 {
								res.SetStatusCode(statusCode)
								res.SetError(err)
								res.sendToChan(chRes)

								return true
							}

							return false
						}

						api.logger.Send("info", fmt.Sprintf("when making a request to add a new custom field for the rootId '%s', caseId '%s', the following is received status code '%d'", rc.RootId, rc.CaseId, statusCode))

						res.SetStatusCode(statusCode)
						res.sendToChan(chRes)

						return true
					})
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
