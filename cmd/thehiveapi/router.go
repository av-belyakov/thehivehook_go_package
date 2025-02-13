package thehiveapi

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/av-belyakov/thehivehook_go_package/internal/supportingfunctions"
)

/*
___ func 'RouteWebHook', BEFORE func 'CreateEvenCase'
!!! func 'CreateEvenCase', root id:'~88356888792' START
___ func 'RouteWebHook', BEFORE func 'CreateEvenCase'
!!! func 'CreateEvenCase', root id:'~88356888792' START
=== func 'router', command:'get_ttp', root id:'~88356888792'
func 'CreateEventCase', goroutine 'observable' received data
__________________ goroutine g.Go TTP STOP
___ func 'RouteWebHook', BEFORE func 'CreateEvenCase'
!!! func 'CreateEvenCase', root id:'~88356888792' START
=== func 'router', command:'get_observables', root id:'~88356888792'
func 'CreateEventCase', goroutine 'observable' received data
__________________ goroutine g.Go Observable STOP
func 'CreateEvenCase', STOP...
___ func 'RouteWebHook', AFTER func 'CreateEvenCase'
___ func 'RouteWebHook', BEFORE func 'CreateEvenCase'
!!! func 'CreateEvenCase', root id:'~88356888792' START
=== func 'router', command:'get_ttp', root id:'~88356888792'
func 'CreateEventCase', goroutine 'observable' received data
__________________ goroutine g.Go TTP STOP
=== func 'router', command:'get_observables', root id:'~88356888792'
func 'CreateEventCase', goroutine 'observable' received data
__________________ goroutine g.Go Observable STOP
func 'CreateEvenCase', STOP...
___ func 'RouteWebHook', AFTER func 'CreateEvenCase'
=== func 'router', command:'get_ttp', root id:'~88356888792'
func 'CreateEventCase', goroutine 'observable' received data
__________________ goroutine g.Go TTP STOP
=== func 'router', command:'get_observables', root id:'~88356888792'
func 'CreateEventCase', goroutine 'observable' received data
__________________ goroutine g.Go Observable STOP
func 'CreateEvenCase', STOP...
___ func 'RouteWebHook', AFTER func 'CreateEvenCase'
=== func 'router', command:'get_ttp', root id:'~88356888792'
func 'CreateEventCase', goroutine 'observable' received data
__________________ goroutine g.Go TTP STOP
func 'CreateEvenCase', STOP...
2025-02-13 17:43:56 ERR - thehivehook_go_package - context deadline exceeded /home/artemij/go/src/thehivehook_go_package/cmd/webhookserver/routes.go:53
=== func 'router', command:'get_observables', root id:'~88356888792'
=== func 'router', command:'get_observables', root id:'~88356888792'
=== func 'router', command:'get_observables', root id:'~88356888792'
___ func 'RouteWebHook', BEFORE func 'CreateEvenCase'
!!! func 'CreateEvenCase', root id:'~88356888792' START
___ func 'RouteWebHook', BEFORE func 'CreateEvenCase'
!!! func 'CreateEvenCase', root id:'~88356888792' START
=== func 'router', command:'get_observables', root id:'~88356888792'
func 'CreateEventCase', goroutine 'observable' received data
__________________ goroutine g.Go Observable STOP
___ func 'RouteWebHook', BEFORE func 'CreateEvenCase'
!!! func 'CreateEvenCase', root id:'~88356888792' START
=== func 'router', command:'get_ttp', root id:'~88356888792'
func 'CreateEventCase', goroutine 'observable' received data
__________________ goroutine g.Go TTP STOP
func 'CreateEvenCase', STOP...
___ func 'RouteWebHook', AFTER func 'CreateEvenCase'
___ func 'RouteWebHook', BEFORE func 'CreateEvenCase'
!!! func 'CreateEvenCase', root id:'~88356888792' START
=== func 'router', command:'get_ttp', root id:'~88356888792'
func 'CreateEventCase', goroutine 'observable' received data
__________________ goroutine g.Go TTP STOP
___ func 'RouteWebHook', BEFORE func 'CreateEvenCase'
!!! func 'CreateEvenCase', root id:'~88356888792' START
=== func 'router', command:'get_observables', root id:'~88356888792'
func 'CreateEventCase', goroutine 'observable' received data
__________________ goroutine g.Go Observable STOP
func 'CreateEvenCase', STOP...
___ func 'RouteWebHook', AFTER func 'CreateEvenCase'
=== func 'router', command:'get_observables', root id:'~88356888792'
func 'CreateEventCase', goroutine 'observable' received data
__________________ goroutine g.Go Observable STOP
=== func 'router', command:'get_ttp', root id:'~88356888792'
func 'CreateEventCase', goroutine 'observable' received data
__________________ goroutine g.Go TTP STOP
func 'CreateEvenCase', STOP...
___ func 'RouteWebHook', AFTER func 'CreateEvenCase'
=== func 'router', command:'get_ttp', root id:'~88356888792'
func 'CreateEventCase', goroutine 'observable' received data
__________________ goroutine g.Go TTP STOP
___ func 'RouteWebHook', BEFORE func 'CreateEvenCase'
!!! func 'CreateEvenCase', root id:'~88356888792' START
___ func 'RouteWebHook', BEFORE func 'CreateEvenCase'
!!! func 'CreateEvenCase', root id:'~88356888792' START
func 'CreateEvenCase', STOP...
2025-02-13 17:44:25 ERR - thehivehook_go_package - context deadline exceeded /home/artemij/go/src/thehivehook_go_package/cmd/webhookserver/routes.go:53
=== func 'router', command:'get_observables', root id:'~88356888792'
func 'CreateEvenCase', STOP...
2025-02-13 17:44:26 ERR - thehivehook_go_package - context deadline exceeded /home/artemij/go/src/thehivehook_go_package/cmd/webhookserver/routes.go:53
=== func 'router', command:'get_ttp', root id:'~88356888792'
=== func 'router', command:'get_observables', root id:'~88356888792'
=== func 'router', command:'get_ttp', root id:'~88356888792'
func 'CreateEventCase', goroutine 'observable' received data
__________________ goroutine g.Go TTP STOP
=== func 'router', command:'get_observables', root id:'~88356888792'
func 'CreateEventCase', goroutine 'observable' received data
__________________ goroutine g.Go Observable STOP
func 'CreateEvenCase', STOP...
___ func 'RouteWebHook', AFTER func 'CreateEvenCase'
func 'CreateEvenCase', STOP...
2025-02-13 17:44:40 ERR - thehivehook_go_package - context deadline exceeded /home/artemij/go/src/thehivehook_go_package/cmd/webhookserver/routes.go:53
=== func 'router', command:'get_ttp', root id:'~88356888792'
=== func 'router', command:'get_observables', root id:'~88356888792'
=== func 'router', command:'get_observables', root id:'~88356888792'
=== func 'router', command:'get_observables', root id:'~88356888792'
=== func 'router', command:'get_observables', root id:'~88356888792'
=== func 'router', command:'get_observables', root id:'~88356888792'
=== func 'router', command:'get_observables', root id:'~88356888792'
=== func 'router', command:'get_observables', root id:'~88356888792'
=== func 'router', command:'get_observables', root id:'~88356888792'
___ func 'RouteWebHook', BEFORE func 'CreateEvenCase'
!!! func 'CreateEvenCase', root id:'~90893435000' START
=== func 'router', command:'get_observables', root id:'~88356888792'
=== func 'router', command:'get_observables', root id:'~88356888792'
___ func 'RouteWebHook', BEFORE func 'CreateEvenCase'
!!! func 'CreateEvenCase', root id:'~90893435000' START
=== func 'router', command:'get_observables', root id:'~88356888792'
=== func 'router', command:'get_observables', root id:'~88356888792'
=== func 'router', command:'get_observables', root id:'~88356888792'
func 'CreateEvenCase', STOP...
2025-02-13 17:45:23 ERR - thehivehook_go_package - context deadline exceeded /home/artemij/go/src/thehivehook_go_package/cmd/webhookserver/routes.go:53
=== func 'router', command:'get_observables', root id:'~88356888792'
=== func 'router', command:'get_observables', root id:'~88356888792'
func 'CreateEvenCase', STOP...
2025-02-13 17:45:26 ERR - thehivehook_go_package - context deadline exceeded /home/artemij/go/src/thehivehook_go_package/cmd/webhookserver/routes.go:53
___ func 'RouteWebHook', BEFORE func 'CreateEvenCase'
!!! func 'CreateEvenCase', root id:'~90893435000' START
___ func 'RouteWebHook', BEFORE func 'CreateEvenCase'
!!! func 'CreateEvenCase', root id:'~90893435000' START
=== func 'router', command:'get_observables', root id:'~88356888792'
___ func 'RouteWebHook', BEFORE func 'CreateEvenCase'
!!! func 'CreateEvenCase', root id:'~90893435000' START
=== func 'router', command:'get_observables', root id:'~88356888792'
=== func 'router', command:'get_observables', root id:'~88356888792'
=== func 'router', command:'get_observables', root id:'~88356888792'
=== func 'router', command:'get_observables', root id:'~88356888792'
func 'CreateEvenCase', STOP...
2025-02-13 17:45:42 ERR - thehivehook_go_package - context deadline exceeded /home/artemij/go/src/thehivehook_go_package/cmd/webhookserver/routes.go:53
func 'CreateEvenCase', STOP...
2025-02-13 17:45:43 ERR - thehivehook_go_package - context deadline exceeded /home/artemij/go/src/thehivehook_go_package/cmd/webhookserver/routes.go:53
=== func 'router', command:'get_observables', root id:'~88356888792'
func 'CreateEvenCase', STOP...
2025-02-13 17:45:44 ERR - thehivehook_go_package - context deadline exceeded /home/artemij/go/src/thehivehook_go_package/cmd/webhookserver/routes.go:53
=== func 'router', command:'get_observables', root id:'~88356888792'
=== func 'router', command:'get_observables', root id:'~88356888792'
=== func 'router', command:'get_observables', root id:'~88356888792'
=== func 'router', command:'get_observables', root id:'~88356888792'
=== func 'router', command:'get_observables', root id:'~88356888792'
=== func 'router', command:'get_observables', root id:'~88356888792'
=== func 'router', command:'get_observables', root id:'~88356888792'
=== func 'router', command:'get_observables', root id:'~88356888792'
___ func 'RouteWebHook', BEFORE func 'CreateEvenCase'
!!! func 'CreateEvenCase', root id:'~90893029496' START
=== func 'router', command:'get_observables', root id:'~88356888792'
___ func 'RouteWebHook', BEFORE func 'CreateEvenCase'
!!! func 'CreateEvenCase', root id:'~90893029496' START
=== func 'router', command:'get_observables', root id:'~88356888792'
=== func 'router', command:'get_observables', root id:'~88356888792'
=== func 'router', command:'get_observables', root id:'~88356888792'
=== func 'router', command:'get_observables', root id:'~88356888792'
func 'CreateEvenCase', STOP...
2025-02-13 17:46:26 ERR - thehivehook_go_package - context deadline exceeded /home/artemij/go/src/thehivehook_go_package/cmd/webhookserver/routes.go:53
=== func 'router', command:'get_observables', root id:'~88356888792'
func 'CreateEvenCase', STOP...
2025-02-13 17:46:27 ERR - thehivehook_go_package - context deadline exceeded /home/artemij/go/src/thehivehook_go_package/cmd/webhookserver/routes.go:53
=== func 'router', command:'get_observables', root id:'~88356888792'
=== func 'router', command:'get_observables', root id:'~88356888792'
=== func 'router', command:'get_observables', root id:'~88356888792'
=== func 'router', command:'get_observables', root id:'~88356888792'
___ func 'RouteWebHook', BEFORE func 'CreateEvenCase'
!!! func 'CreateEvenCase', root id:'~90893029496' START
=== func 'router', command:'get_observables', root id:'~88356888792'
=== func 'router', command:'get_observables', root id:'~88356888792'
___ func 'RouteWebHook', BEFORE func 'CreateEvenCase'
!!! func 'CreateEvenCase', root id:'~90893029496' START
=== func 'router', command:'get_observables', root id:'~88356888792'
=== func 'router', command:'get_observables', root id:'~88356888792'
=== func 'router', command:'get_observables', root id:'~88356888792'
func 'CreateEvenCase', STOP...
2025-02-13 17:46:54 ERR - thehivehook_go_package - context deadline exceeded /home/artemij/go/src/thehivehook_go_package/cmd/webhookserver/routes.go:53
=== func 'router', command:'get_observables', root id:'~88356888792'
=== func 'router', command:'get_observables', root id:'~88356888792'
func 'CreateEvenCase', STOP...
2025-02-13 17:47:01 ERR - thehivehook_go_package - context deadline exceeded /home/artemij/go/src/thehivehook_go_package/cmd/webhookserver/routes.go:53
=== func 'router', command:'get_observables', root id:'~88356888792'
=== func 'router', command:'get_observables', root id:'~88356888792'
___ func 'RouteWebHook', BEFORE func 'CreateEvenCase'
!!! func 'CreateEvenCase', root id:'~90320322744' START
=== func 'router', command:'get_observables', root id:'~88356888792'
=== func 'router', command:'get_observables', root id:'~88356888792'
=== func 'router', command:'get_observables', root id:'~88356888792'
=== func 'router', command:'get_observables', root id:'~88356888792'
___ func 'RouteWebHook', BEFORE func 'CreateEvenCase'
!!! func 'CreateEvenCase', root id:'~90320322744' START
=== func 'router', command:'get_observables', root id:'~88356888792'
___ func 'RouteWebHook', BEFORE func 'CreateEvenCase'
!!! func 'CreateEvenCase', root id:'~90320322744' START
func 'CreateEvenCase', STOP...
2025-02-13 17:47:22 ERR - thehivehook_go_package - context deadline exceeded /home/artemij/go/src/thehivehook_go_package/cmd/webhookserver/routes.go:53
=== func 'router', command:'get_observables', root id:'~88356888792'
=== func 'router', command:'get_observables', root id:'~88356888792'
___ func 'RouteWebHook', BEFORE func 'CreateEvenCase'
!!! func 'CreateEvenCase', root id:'~90320322744' START
=== func 'router', command:'get_observables', root id:'~88356888792'
___ func 'RouteWebHook', BEFORE func 'CreateEvenCase'
!!! func 'CreateEvenCase', root id:'~90320322744' START
=== func 'router', command:'get_observables', root id:'~88356888792'
func 'CreateEvenCase', STOP...
2025-02-13 17:47:33 ERR - thehivehook_go_package - context deadline exceeded /home/artemij/go/src/thehivehook_go_package/cmd/webhookserver/routes.go:53
=== func 'router', command:'get_observables', root id:'~88356888792'
func 'CreateEvenCase', STOP...
2025-02-13 17:47:35 ERR - thehivehook_go_package - context deadline exceeded /home/artemij/go/src/thehivehook_go_package/cmd/webhookserver/routes.go:53
=== func 'router', command:'get_observables', root id:'~88356888792'
=== func 'router', command:'get_observables', root id:'~88356888792'
func 'CreateEvenCase', STOP...
2025-02-13 17:47:42 ERR - thehivehook_go_package - context deadline exceeded /home/artemij/go/src/thehivehook_go_package/cmd/webhookserver/routes.go:53
=== func 'router', command:'get_observables', root id:'~88356888792'
func 'CreateEvenCase', STOP...
2025-02-13 17:47:44 ERR - thehivehook_go_package - context deadline exceeded /home/artemij/go/src/thehivehook_go_package/cmd/webhookserver/routes.go:53
=== func 'router', command:'get_observables', root id:'~88356888792'
=== func 'router', command:'get_observables', root id:'~88356888792'
=== func 'router', command:'get_observables', root id:'~88356888792'
=== func 'router', command:'get_observables', root id:'~88356888792'
=== func 'router', command:'get_observables', root id:'~88356888792'
=== func 'router', command:'get_observables', root id:'~88356888792'
=== func 'router', command:'get_observables', root id:'~88356888792'
=== func 'router', command:'get_observables', root id:'~88356888792'
=== func 'router', command:'get_observables', root id:'~88356888792'
=== func 'router', command:'get_observables', root id:'~88356888792'
=== func 'router', command:'get_observables', root id:'~88356888792'
=== func 'router', command:'get_observables', root id:'~88356888792'
=== func 'router', command:'get_observables', root id:'~88356888792'
=== func 'router', command:'get_observables', root id:'~88356888792'
=== func 'router', command:'get_observables', root id:'~88356888792'
=== func 'router', command:'get_observables', root id:'~88356888792'
=== func 'router', command:'get_observables', root id:'~88356888792'
=== func 'router', command:'get_observables', root id:'~88356888792'
*/

func (api *apiTheHiveModule) router(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return

		case msg := <-api.receivingChannel:
			switch msg.GetCommand() {
			case "get_observables":
				go func(id string) {
					so := NewSpecialObjectForCache[interface{}]()
					so.SetID(id)
					so.SetFunc(func(_ int) bool {

						fmt.Printf("=== func 'router', command:'%s', root id:'%s'\n", msg.GetCommand(), msg.GetRootId())

						res, statusCode, err := api.GetObservables(ctx, msg.GetRootId())
						if err != nil {
							api.logger.Send("error", supportingfunctions.CustomError(err).Error())

							return false
						}

						newRes := NewChannelRespons()
						newRes.SetRequestId(msg.GetRootId())
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
				}(msg.GetRootId())

			case "get_ttp":
				go func(id string) {
					so := NewSpecialObjectForCache[interface{}]()
					so.SetID(id)
					so.SetFunc(func(_ int) bool {

						fmt.Printf("=== func 'router', command:'%s', root id:'%s'\n", msg.GetCommand(), msg.GetRootId())

						res, statusCode, err := api.GetTTP(ctx, msg.GetRootId())
						if err != nil {
							api.logger.Send("error", supportingfunctions.CustomError(err).Error())

							return false
						}

						newRes := NewChannelRespons()
						newRes.SetRequestId(msg.GetRootId())
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
				}(msg.GetRootId())

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
					go func(id string) {
						so := NewSpecialObjectForCache[interface{}]()
						so.SetID(id)
						so.SetFunc(func(countAttempts int) bool {
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

						//добавляем объект в очередь для обработки
						api.cache.PushObjectToQueue(so)
					}(msg.GetRootId())

				case "add_case_task":
					go func(id string) {
						so := NewSpecialObjectForCache[interface{}]()
						so.SetID(id)
						so.SetFunc(func(countAttempts int) bool {
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

						//добавляем объект в очередь для обработки
						api.cache.PushObjectToQueue(so)
					}(msg.GetRequestId())

				case "set_case_custom_field":
					go func(id string) {
						so := NewSpecialObjectForCache[interface{}]()
						so.SetID(id)
						so.SetFunc(func(countAttempts int) bool {
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

						//добавляем объект в очередь для обработки
						api.cache.PushObjectToQueue(so)
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
