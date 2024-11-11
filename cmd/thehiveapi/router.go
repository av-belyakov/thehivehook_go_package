package thehiveapi

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"runtime"
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

			case "send_command":
				rc, err := getRequestCommandData(msg.GetData())
				if err != nil {
					_, f, l, _ := runtime.Caller(0)
					api.logger.Send("error", fmt.Sprintf("%s %s:%d", err.Error(), f, l-2))
				}

				chRes := msg.GetChanOutput()

				var isSuccess bool = true
				res := NewChannelRespons()
				res.SetRequestId(msg.GetRequestId())

				switch msg.GetOrder() {
				case "add_case_tags":
					api.cacheRunningFunction.SetMethod(msg.GetRootId(), func() bool {
						_, statusCode, err := api.AddCaseTags(ctx, rc)
						if err != nil {
							api.logger.Send("error", err.Error())
							res.SetError(err)

							isSuccess = false
						}

						if isSuccess {
							api.logger.Send("info", fmt.Sprintf("when making a request to add a new tag for the rootId '%s', caseId '%s', the following is received status code '%d'", rc.RootId, rc.CaseId, statusCode))
						}

						res.SetStatusCode(statusCode)
						res.sendToChan(chRes)

						return isSuccess
					})

				case "add_case_task":
					api.cacheRunningFunction.SetMethod(msg.GetRootId(), func() bool {
						_, statusCode, err := api.AddCaseTask(ctx, rc)
						if err != nil {
							api.logger.Send("error", err.Error())
							res.SetError(err)

							isSuccess = false
						}

						if isSuccess {
							api.logger.Send("info", fmt.Sprintf("when making a request to add a new task for the rootId '%s', caseId '%s', the following is received status code '%d'", rc.RootId, rc.CaseId, statusCode))
						}

						res.SetStatusCode(statusCode)
						res.sendToChan(chRes)

						return isSuccess
					})

				case "set_case_custom_field":
					api.cacheRunningFunction.SetMethod(msg.GetRootId(), func() bool {
						_, statusCode, err := api.AddCaseCustomFields(ctx, rc)
						if err != nil {
							api.logger.Send("error", err.Error())
							res.SetError(err)

							isSuccess = false
						}

						if isSuccess {
							api.logger.Send("info", fmt.Sprintf("when making a request to add a new custom field for the rootId '%s', caseId '%s', the following is received status code '%d'", rc.RootId, rc.CaseId, statusCode))
						}

						res.SetStatusCode(statusCode)
						res.sendToChan(chRes)

						return isSuccess
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
