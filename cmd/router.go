package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"golang.org/x/sync/errgroup"

	"github.com/av-belyakov/thehivehook_go_package/v2/cmd/natsapi"
	"github.com/av-belyakov/thehivehook_go_package/v2/cmd/thehiveapi"
	"github.com/av-belyakov/thehivehook_go_package/v2/cmd/webhookserver"
	"github.com/av-belyakov/thehivehook_go_package/v2/internal/confighandler"
	"github.com/av-belyakov/thehivehook_go_package/v2/internal/counterelements"
	"github.com/av-belyakov/thehivehook_go_package/v2/internal/datamodels"
	"github.com/av-belyakov/thehivehook_go_package/v2/internal/interfaces"
	"github.com/av-belyakov/thehivehook_go_package/v2/internal/storageobjects"
	"github.com/av-belyakov/thehivehook_go_package/v2/internal/supportingfunctions"
)

const Time_After_Delete = 30

// NewMajorRouter основной маршрутизатор
func NewMajorRouter(cfg confighandler.ConfigApp, logger interfaces.Logger) *majorRouter {
	return &majorRouter{
		cfg:    cfg,
		logger: logger,
	}
}

func (r *majorRouter) start(
	ctx context.Context,
	toNatsApi chan<- natsapi.InputData,
	fromWebHook <-chan webhookserver.OutputData,
	fromNatsApi <-chan natsapi.OutputData,
) error {
	var timeToSend int = 10

	//очередь для хранения объектов с отложенной отправкой в канал
	storage, err := storageobjects.New(
		storageobjects.WithChannelSize[map[string]any](4),
		storageobjects.WithTimeTick[map[string]any](1),
		storageobjects.WithTimeToLive[map[string]any](r.cfg.GetApplicationTemporaryStorage().StorageObjectTTL),
	)
	if err != nil {
		return supportingfunctions.CustomError(err)
	}
	storage.Start(ctx)

	//счетчик комад получаемых через Nats
	// нужен для отслеживания количества команд, на каждую команду TheHive отправляет ответ в виде кейса,
	// в результате получается замкнутый цикл когда изменение в кейсах порождает команды на добавления
	// тегов и т.д., а команды, в свою очередь порождают изменения в кейсах
	counter := counterelements.New(Time_After_Delete)
	counter.Start(ctx)

	//модуль для взаимодействия с API TheHive
	apiTheHive, err := thehiveapi.New(
		thehiveapi.WithHost(r.cfg.GetApplicationTheHive().Host),
		thehiveapi.WithPort(r.cfg.GetApplicationTheHive().Port),
		thehiveapi.WithAPIKey(r.cfg.GetApplicationTheHive().ApiKey),
		thehiveapi.WithNameRegionalObject(r.cfg.GetApplicationWebHookServer().Name),
	)
	if err != nil {
		return supportingfunctions.CustomError(err)
	}

	go func() {
		for {
			select {
			case <-ctx.Done():
				return

			case msg := <-fromWebHook:
				r.logger.Send("info", fmt.Sprintf("'majorRouter' received an object of the type:'%s'", msg.ObjectType))

				// -------------------------------------
				//все данные полученные из webhookserver
				if msg.ObjectType == "alert" {
					timeToSend = r.cfg.StorageDelayToSendAlert
				}
				if msg.ObjectType == "case" {
					//уменьшаем счётчик на единицу и проверяем достиг ли счётчик нуля, если счётчик
					// равен нулю, то запись об элементе удаляется из счётчика, а сам объект игнорируется
					// если счётчик больше нуля, то игнорируем объект так как этот объект мог быть порождён
					// результатом выполнения команды
					counter.Done(msg.RootId)
					countElements := counter.Get(msg.RootId)
					if countElements >= 0 {
						if countElements == 0 {
							counter.DeleteKey(msg.RootId)
						}

						r.logger.Send("info", fmt.Sprintf("'majorRouter' the object type '%s', may be a consequence of executing some command, skip adding to the cache", msg.ObjectType))

						continue
					}

					timeToSend = r.cfg.StorageDelayToSendCase
				}

				r.logger.Send("info", fmt.Sprintf("'majorRouter' saving a type '%s' object for deferred execution", msg.ObjectType))

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
				r.logger.Send("info", fmt.Sprintf("'majorRouter' data was received from the deferred cache, object type '%s'", data.ObjectType))

				// -------------------------------------
				//все данные полученные из очереди хранилища (для временной задержки)
				switch data.ObjectType {
				case "alert":
					go func() {
						r.logger.Send("info", "'majorRouter' created goroutine for to enrich an object of the type 'alert'")

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

						r.logger.Send("info", "'majorRouter' transfer of an enriched verified object of type 'alert' to the module natsapi")

						toNatsApi <- natsapi.InputData{
							ElementType: data.ObjectType,
							RootId:      data.Id,
							Data:        b,
						}
					}()

				case "case":
					go func() {
						verifiedObject := datamodels.VerifiedObjectEventCase{
							Source: r.cfg.AppConfigWebHookServer.Name,
							Case:   data.Data,
						}

						r.logger.Send("info", "'majorRouter' created goroutine for to enrich an object of the type 'case'")

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
							//запрос объекта типа TTP
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

						r.logger.Send("info", "'majorRouter' transfer of an enriched verified object of type 'case' to the module natsapi")

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

				//парсим входящие команды
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
						//увеличиваем счётчик для объекта на 1
						counter.Add(rc.RootId)

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
						//увеличиваем счётчик для объекта на 1
						counter.Add(rc.RootId)

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
						//увеличиваем счётчик для объекта на 1
						counter.Add(rc.RootId)

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
