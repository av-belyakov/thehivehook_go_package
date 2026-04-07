package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/av-belyakov/thehivehook_go_package/cmd/thehiveapi"
	"github.com/av-belyakov/thehivehook_go_package/cmd/webhookserver"
	"github.com/av-belyakov/thehivehook_go_package/internal/confighandler"
	"github.com/av-belyakov/thehivehook_go_package/internal/datamodels"
	"github.com/av-belyakov/thehivehook_go_package/internal/interfaces"
	"github.com/av-belyakov/thehivehook_go_package/internal/storageobjects"
	"github.com/av-belyakov/thehivehook_go_package/internal/supportingfunctions"
	"golang.org/x/sync/errgroup"
)

func router_new(
	ctx context.Context,
	cfg confighandler.ConfigApp,
	apiTheHive *thehiveapi.TheHiveApi_New,
	fromWebHook <-chan webhookserver.OutputData,
	logger interfaces.Logger,
) error {
	var timeToSend int = 10

	// очередь для хранения объектов с отложенной отправкой в канал
	storage, err := storageobjects.New[map[string]any](
		storageobjects.WithChannelSize[map[string]any](4),
		storageobjects.WithTimeTick[map[string]any](1),
		storageobjects.WithTimeToLive[map[string]any](cfg.GetApplicationTemporaryStorage().StorageObjectTTL),
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

			case data := <-storage.GetObjects():
				//просто что бы понимать какие объекты передаются
				//rootId := data.Id
				//objectType := data.ObjectType
				//object := data.Data

				switch data.ObjectType {
				case "alert":
					go func() {
						res, statusCode, err := apiTheHive.GetAlert(ctx, data.Id)
						if err != nil {
							if _, ok := errors.AsType[thehiveapi.ErrorInformation](err); ok {
								logger.Send("info", err.Error())
							}

							logger.Send("error", supportingfunctions.CustomError(fmt.Errorf("%w, root id:'%s'", err, data.Id)).Error())

							return
						}

						logger.Send("info", fmt.Sprintf("the request for additional information about 'alert' was completed successfully, successful response to TheHive request, root id:'%s', status code:'%d'", data.Id, statusCode))

						element := map[string]any{}
						if err := json.Unmarshal(res, &element); err != nil {
							logger.Send("error", supportingfunctions.CustomError(err).Error())

							return
						}

						// !!!!!!! доделать как появится новая структура канала примёма NatsApi
						//отправка события в nats
						verifiedObject := datamodels.VerifiedObjectEventAlert{
							Source: cfg.AppConfigWebHookServer.Name,
							Event:  data.Data,
							Alert:  element,
						}
					}()

				case "case":
					verifiedObject := datamodels.VerifiedObjectEventCase{
						Source: cfg.AppConfigWebHookServer.Name,
						Case:   data.Data,
					}

					var g errgroup.Group
					g.Go(func() error {
						// запрос объекта типа Observables
						res, statusCode, err := apiTheHive.GetObservables(ctx, data.Id)
						if err != nil {
							return err
						}

						logger.Send("info", fmt.Sprintf("the request for additional information about 'observables' was completed successfully, successful response to TheHive request, root id:'%s', status code:'%d'", data.Id, statusCode))

						element := []any{}
						if err := json.Unmarshal(res, &element); err != nil {
							logger.Send("error", supportingfunctions.CustomError(err).Error())

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

						logger.Send("info", fmt.Sprintf("the request for additional information about 'ttp' was completed successfully, successful response to TheHive request, root id:'%s', status code:'%d'", data.Id, statusCode))

						element := []any{}
						if err := json.Unmarshal(res, &element); err != nil {
							logger.Send("error", supportingfunctions.CustomError(err).Error())

							return err
						}

						verifiedObject.TTPs = element

						return nil
					})
					if err := g.Wait(); err != nil {
						if err != nil {
							if _, ok := errors.AsType[thehiveapi.ErrorInformation](err); ok {
								logger.Send("info", err.Error())
							}

							logger.Send("error", supportingfunctions.CustomError(fmt.Errorf("%w, root id:'%s'", err, data.Id)).Error())

							return
						}
					}

					// !!!!!!! доделать как появится новая структура канала примёма NatsApi
					//отправка события в nats
					verifiedObject
				}

			case msg := <-fromWebHook:
				if msg.ObjectType == "alert" {
					timeToSend = cfg.StorageDelayToSendAlert
				}
				if msg.ObjectType == "case" {
					timeToSend = cfg.StorageDelayToSendCase
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

				/*
					switch msg.ForSomebody {
					case "to thehive":
						toTheHiveAPI <- msg.Data

					case "to nats":
						toNatsAPI <- msg.Data
					}
				*/

			case msg := <-fromNatsAPI:
				//обработка входящих команд
				toTheHiveAPI <- msg

			}
		}
	}()

	return nil
}
