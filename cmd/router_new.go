package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"golang.org/x/sync/errgroup"

	"github.com/av-belyakov/thehivehook_go_package/cmd/natsapi"
	"github.com/av-belyakov/thehivehook_go_package/cmd/thehiveapi"
	"github.com/av-belyakov/thehivehook_go_package/cmd/webhookserver"
	"github.com/av-belyakov/thehivehook_go_package/internal/confighandler"
	"github.com/av-belyakov/thehivehook_go_package/internal/datamodels"
	"github.com/av-belyakov/thehivehook_go_package/internal/interfaces"
	"github.com/av-belyakov/thehivehook_go_package/internal/storageobjects"
	"github.com/av-belyakov/thehivehook_go_package/internal/supportingfunctions"
)

func NewMajorRouter(cfg confighandler.ConfigApp, logger interfaces.Logger) *majorRouter {
	return &majorRouter{
		cfg:    cfg,
		logger: logger,
	}
}

func (r *majorRouter) start(
	ctx context.Context,
	apiTheHive *thehiveapi.TheHiveApi_New,
	toNatsApi chan<- natsapi.InputData,
	fromWebHook <-chan webhookserver.OutputData,
	fromNatsApi <-chan natsapi.OutputData,
) error {
	var timeToSend int = 10

	// очередь для хранения объектов с отложенной отправкой в канал
	storage, err := storageobjects.New(
		storageobjects.WithChannelSize[map[string]any](4),
		storageobjects.WithTimeTick[map[string]any](1),
		storageobjects.WithTimeToLive[map[string]any](r.cfg.GetApplicationTemporaryStorage().StorageObjectTTL),
	)
	if err != nil {
		return supportingfunctions.CustomError(err)
	}
	storage.Start(ctx)

	go func() {
		for {
			select {
			case <-ctx.Done():
				return

			case msg := <-fromWebHook:
				// -------------------------------------
				//все данные полученные из webhookserver
				if msg.ObjectType == "alert" {
					timeToSend = r.cfg.StorageDelayToSendAlert
				}
				if msg.ObjectType == "case" {
					timeToSend = r.cfg.StorageDelayToSendCase
				}

				//сохраняем объект для отложеного выполнения
				// ожидание истечения времени из параметра storageObject[n].timeSending
				storage.AddObject(
					timeToSend,
					storageobjects.StorageObjectDataSettings[map[string]any]{
						Id:         msg.RootId,
						ObjectType: msg.ObjectType,
						Data:       msg.Data,
					})

			case data := <-storage.GetObjects():
				// -------------------------------------
				//все данные полученные из очереди хранилища (для временной задержки)
				switch data.ObjectType {
				case "alert":
					go func() {
						ctx, cancel := context.WithTimeout(ctx, 15*time.Second)
						defer cancel()

						res, statusCode, err := apiTheHive.GetAlert(ctx, data.Id)
						if err != nil {
							if _, ok := errors.AsType[thehiveapi.ErrorInformation](err); ok {
								r.logger.Send("info", err.Error())
							}

							r.logger.Send("error", supportingfunctions.CustomError(fmt.Errorf("%w, root id:'%s'", err, data.Id)).Error())

							return
						}

						r.logger.Send("info", fmt.Sprintf("the request for additional information about 'alert' was completed successfully, successful response to TheHive request, root id:'%s', status code:'%d'", data.Id, statusCode))

						element := map[string]any{}
						if err := json.Unmarshal(res, &element); err != nil {
							r.logger.Send("error", supportingfunctions.CustomError(err).Error())

							return
						}

						b, err := json.Marshal(datamodels.VerifiedObjectEventAlert{
							Source: r.cfg.AppConfigWebHookServer.Name,
							Event:  data.Data,
							Alert:  element,
						})
						if err != nil {
							r.logger.Send("error", supportingfunctions.CustomError(err).Error())

							return
						}

						toNatsApi <- natsapi.InputData{
							ElementType: data.ObjectType,
							RootId:      data.Id,
							Data:        b,
						}
					}()

				case "case":
					go func() {
						ctx, cancel := context.WithTimeout(ctx, 15*time.Second)
						defer cancel()

						verifiedObject := datamodels.VerifiedObjectEventCase{
							Source: r.cfg.AppConfigWebHookServer.Name,
							Case:   data.Data,
						}

						var g errgroup.Group
						g.Go(func() error {
							// запрос объекта типа Observables
							res, statusCode, err := apiTheHive.GetObservables(ctx, data.Id)
							if err != nil {
								return err
							}

							r.logger.Send("info", fmt.Sprintf("the request for additional information about 'observables' was completed successfully, successful response to TheHive request, root id:'%s', status code:'%d'", data.Id, statusCode))

							element := []any{}
							if err := json.Unmarshal(res, &element); err != nil {
								r.logger.Send("error", supportingfunctions.CustomError(err).Error())

								return err
							}

							verifiedObject.Observables = element

							return nil
						})
						g.Go(func() error {
							// запрос объекта типа TTP
							res, statusCode, err := apiTheHive.GetTTP(ctx, data.Id)
							if err != nil {
								return err
							}

							r.logger.Send("info", fmt.Sprintf("the request for additional information about 'ttp' was completed successfully, successful response to TheHive request, root id:'%s', status code:'%d'", data.Id, statusCode))

							element := []any{}
							if err := json.Unmarshal(res, &element); err != nil {
								r.logger.Send("error", supportingfunctions.CustomError(err).Error())

								return err
							}

							verifiedObject.TTPs = element

							return nil
						})
						if err := g.Wait(); err != nil {
							if _, ok := errors.AsType[thehiveapi.ErrorInformation](err); ok {
								r.logger.Send("info", err.Error())
							} else {
								r.logger.Send("error", supportingfunctions.CustomError(fmt.Errorf("%w, root id:'%s'", err, data.Id)).Error())
							}

							return
						}

						var caseId string
						caseIdInt, err := supportingfunctions.GetCaseIdFromEventTheHive(data.Data)
						if err == nil {
							caseId = strconv.Itoa(caseIdInt)
						}

						b, err := json.Marshal(verifiedObject)
						if err != nil {
							r.logger.Send("error", supportingfunctions.CustomError(fmt.Errorf("%w, root id:'%s', case id:'%s'", err, data.Id, caseId)).Error())

							return
						}

						toNatsApi <- natsapi.InputData{
							ElementType: data.ObjectType,
							RootId:      data.Id,
							CaseId:      caseId,
							Data:        b,
						}
					}()
				}

			case msg := <-fromNatsApi:
				// -------------------------------------
				//обработка входящих команд (добавление tags, custom fields и т.д.)
				verifiedResponse := datamodels.VerifiedResponseAcceptedCommand{
					Source:     r.cfg.GetApplicationWebHookServer().Name,
					StatusCode: http.StatusBadRequest,
				}

				// парсим входящие команды
				rc := thehiveapi.RequestCommand{}
				if err := json.Unmarshal(msg.Data, &rc); err != nil {
					errMsg := fmt.Errorf("the request contains an invalid json object (%s)", err.Error())
					verifiedResponse.Error = errMsg.Error()
					r.logger.Send("error", supportingfunctions.CustomError(errMsg).Error())

					sendData(verifiedResponse, msg.ChanDone, msg.ChanOutput)

					continue
				}

				verifiedResponse.Id = rc.RootId
				verifiedResponse.Command = rc.Command

				//проверяем какому из thehivehook_go_package было предназначена команда
				// так как данный модуль распространяется по регионам, необходимо отличать
				// запросы на изменения тегов, добавления задач или специальных полей,
				// предназначенных для определённого модуля
				if r.cfg.GetApplicationWebHookServer().Name != rc.RegionalObject {
					errMsg := fmt.Errorf(
						"the command '%s' cannot be executed because the name of the regional object '%s' does not match with name '%s'",
						rc.Command,
						rc.RegionalObject,
						r.cfg.GetApplicationWebHookServer().Name,
					)
					verifiedResponse.Error = errMsg.Error()
					verifiedResponse.StatusCode = http.StatusPreconditionFailed // условие ложно

					r.logger.Send("error", supportingfunctions.CustomError(errMsg).Error())
					sendData(verifiedResponse, msg.ChanDone, msg.ChanOutput)

					continue
				}

				/*
					!!! ЭТО БУДУЩИЙ ответ на команду !!!

						//ответ на команду
						res, err := json.Marshal(struct {
							Id         string `json:"id"`
							Source     string `json:"source"`
							Command    string `json:"command"`
							StatusCode int    `json:"status_code"`
							Data       any    `json:"data"`
							Error      string `json:"error"`
						}{
							Id:         msg.GetRequestId(),
							Source:     msg.GetSource(),
							Command:    rc.Command,
							StatusCode: msg.GetStatusCode(),
							Data:       msg.GetData(),
							Error:      errMsg,
						})

				*/

				r.logger.Send(
					"info",
					fmt.Sprintf("the command '%s' has been received, rootId '%s', for regional object '%s', currnet regional object '%s'",
						rc.Command,
						rc.RootId,
						rc.RegionalObject,
						r.cfg.GetApplicationWebHookServer().Name,
					))

				switch rc.Command {
				case "add_case_tag":
					go func() {
						_, statusCode, err := apiTheHive.AddCaseTags(ctx, rc)
						if err != nil {
							verifiedResponse.Error = err.Error()
							r.logger.Send("error", supportingfunctions.CustomError(err).Error())
						} else {
							r.logger.Send(
								"info",
								fmt.Sprintf("when making a request to add a new tag '%s' for the service '%s' rootId '%s', caseId '%s', the following is received status code '%d'",
									rc.Value,
									rc.Service,
									rc.RootId,
									rc.CaseId,
									statusCode,
								))
						}

						verifiedResponse.StatusCode = statusCode
					}()

				case "add_case_task":
					go func() {
						_, statusCode, err := apiTheHive.AddCaseTask(ctx, rc)
						if err != nil {
							verifiedResponse.Error = err.Error()
							r.logger.Send("error", supportingfunctions.CustomError(err).Error())
						} else {
							r.logger.Send(
								"info",
								fmt.Sprintf("when making a request to add a new case task '%s' for the service '%s' rootId '%s', caseId '%s', the following is received status code '%d'",
									rc.Value,
									rc.Service,
									rc.RootId,
									rc.CaseId,
									statusCode,
								))
						}

						verifiedResponse.StatusCode = statusCode
					}()

				case "set_case_custom_field":
					go func() {
						_, statusCode, err := apiTheHive.AddCaseCustomFields(ctx, rc)
						if err != nil {
							verifiedResponse.Error = err.Error()
							if _, ok := errors.AsType[thehiveapi.ErrorInformation](err); ok {
								r.logger.Send("info", err.Error())
							} else {
								r.logger.Send("error", supportingfunctions.CustomError(err).Error())
							}
						} else {
							r.logger.Send(
								"info",
								fmt.Sprintf("when making a request to add a new custom field '%s' for the service '%s' rootId '%s', caseId '%s', the following is received status code '%d'",
									rc.Value,
									rc.Service,
									rc.RootId,
									rc.CaseId,
									statusCode,
								))
						}

						verifiedResponse.StatusCode = statusCode
					}()
				}

				sendData(verifiedResponse, msg.ChanDone, msg.ChanOutput)
			}
		}
	}()

	return nil
}
