package main

import (
	"context"
	"fmt"
	"log"
	"runtime"

	"github.com/av-belyakov/simplelogger"
	"github.com/av-belyakov/thehivehook_go_package/cmd/commoninterfaces"
	"github.com/av-belyakov/thehivehook_go_package/cmd/natsapi"
	"github.com/av-belyakov/thehivehook_go_package/cmd/thehiveapi"
	"github.com/av-belyakov/thehivehook_go_package/cmd/webhookserver"
	"github.com/av-belyakov/thehivehook_go_package/cmd/wrappers"
	"github.com/av-belyakov/thehivehook_go_package/internal/confighandler"
	"github.com/av-belyakov/thehivehook_go_package/internal/logginghandler"
	"github.com/av-belyakov/thehivehook_go_package/internal/supportingfunctions"
	"github.com/av-belyakov/thehivehook_go_package/internal/versionandname"
)

// server здесь реализована вся логика запуска thehivehook_go_package
func server(ctx context.Context) {
	rootPath, err := supportingfunctions.GetRootPath(ROOT_DIR)
	if err != nil {
		log.Fatalf("error, it is impossible to form root path (%w)", err)
	}

	//чтение конфигурационного файла
	confApp, err := confighandler.NewConfig(rootPath)
	if err != nil {
		log.Fatalf("error module 'confighandler': %w", err)
	}

	//******************************************************
	//********** инициализация модуля логирования **********
	simpleLogger, err := simplelogger.NewSimpleLogger(ctx, ROOT_DIR, getLoggerSettings(confApp.GetListLogs()))
	if err != nil {
		log.Fatalf("error module 'simplelogger': %w", err)
	}

	//*****************************************************************
	//********** инициализация модуля взаимодействия с Zabbix **********
	channelZabbix := make(chan commoninterfaces.Messager)
	wzis := wrappers.WrappersZabbixInteractionSettings{
		NetworkPort: confApp.Zabbix.NetworkPort,
		NetworkHost: confApp.Zabbix.NetworkHost,
		ZabbixHost:  confApp.Zabbix.ZabbixHost}
	eventTypes := []wrappers.EventType(nil)
	for _, v := range confApp.Zabbix.EventTypes {
		eventTypes = append(eventTypes, wrappers.EventType{
			IsTransmit: v.IsTransmit,
			EventType:  v.EventType,
			ZabbixKey:  v.ZabbixKey,
			Handshake: wrappers.Handshake{
				TimeInterval: v.Handshake.TimeInterval,
				Message:      v.Handshake.Message,
			},
		})
	}
	if err := wrappers.WrappersZabbixInteraction(ctx, simpleLogger, wzis, channelZabbix); err != nil {
		_, f, l, _ := runtime.Caller(0)
		_ = simpleLogger.WriteLoggingData(fmt.Sprintf(" '%w' %s:%d", err, f, l-3), "error")

		log.Fatalf("error module 'zabbixinteraction': %w\n", err)
	}

	//******************************************************************
	//********** инициализация обработчика логирования данных **********
	logging := logginghandler.New()
	go logginghandler.LoggingHandler(ctx, simpleLogger, channelZabbix, logging.GetChan())

	//******************************************************
	//********** инициализация TheHive API модуля **********
	confTheHiveAPI := confApp.GetApplicationTheHive()
	apiTheHive, err := thehiveapi.New(
		logging,
		thehiveapi.WithAPIKey(confTheHiveAPI.ApiKey),
		thehiveapi.WithHost(confTheHiveAPI.Host),
		thehiveapi.WithPort(confTheHiveAPI.Port))
	if err != nil {
		_, f, l, _ := runtime.Caller(0)
		_ = simpleLogger.WriteLoggingData(fmt.Sprintf(" '%w' %s:%d", err, f, l-3), "error")

		log.Fatalf("error module 'thehiveapi': %w\n", err)
	}
	chanRequestTheHiveAPI, err := apiTheHive.Start(ctx)
	if err != nil {
		if err != nil {
			_, f, l, _ := runtime.Caller(0)
			_ = simpleLogger.WriteLoggingData(fmt.Sprintf(" '%w' %s:%d", err, f, l-3), "error")

			log.Fatalf("error module 'thehiveapi': %w\n", err)
		}
	}

	//***************************************************
	//********** инициализация NATS API модуля **********
	confNatsSAPI := confApp.GetApplicationNATS()
	natsOptsAPI := []natsapi.NatsAPIOptions{
		natsapi.WithHost(confNatsSAPI.Host),
		natsapi.WithPort(confNatsSAPI.Port),
	}
	for _, v := range confNatsSAPI.Subscribers {
		natsOptsAPI = append(natsOptsAPI, natsapi.WithSubscribers(v.Event, v.Responders))
	}
	apiNats, err := natsapi.New(logging, natsOptsAPI...)
	chanRequestNatsAPI := apiNats.Start(ctx)

	//*********************************************************
	//********** инициализация WEBHOOKSERVER сервера **********
	confWebHook := confApp.GetApplicationWebHookServer()
	webHook, chanForSomebody, err := webhookserver.New(
		logging,
		webhookserver.WithTTL(confApp.TTLTmpInfo),
		webhookserver.WithPort(confWebHook.Port),
		webhookserver.WithHost(confWebHook.Host),
		webhookserver.WithName(confWebHook.Name),
		webhookserver.WithVersion(versionandname.GetVersion()),
		webhookserver.WithPathSqlite(confApp.GetApplicationSqlite().PathDatabase))
	if err != nil {
		_, f, l, _ := runtime.Caller(0)
		_ = simpleLogger.WriteLoggingData(fmt.Sprintf(" '%w' %s:%d", err, f, l-3), "error")

		log.Fatalf("error module 'webhookserver': %w\n", err)
	}

	//мост между каналами различных модулей, где любой канал модуля должен
	//удовлетворять интерфейсу commoninterfaces.ChannelRequester и каналом
	//для взаимодействия с webhookserver
	go router(chanForSomebody, chanRequestTheHiveAPI, chanRequestNatsAPI)

	webHook.Start(ctx)
	webHook.Shutdown(ctx)
}
