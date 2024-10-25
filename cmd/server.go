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
	"github.com/av-belyakov/thehivehook_go_package/internal/confighandler"
	"github.com/av-belyakov/thehivehook_go_package/internal/logginghandler"
	"github.com/av-belyakov/thehivehook_go_package/internal/supportingfunctions"
	"github.com/av-belyakov/thehivehook_go_package/internal/versionandname"
)

func server(ctx context.Context) {
	rootPath, err := supportingfunctions.GetRootPath(ROOT_DIR)
	if err != nil {
		log.Fatalf("error, it is impossible to form root path (%v)", err)
	}

	//чтение конфигурационного файла
	confApp, err := confighandler.NewConfig(rootPath)
	if err != nil {
		log.Fatalf("error module 'confighandler': %v", err)
	}

	//******************************************************
	//********** инициализация модуля логирования **********
	simpleLogger, err := simplelogger.NewSimpleLogger(ROOT_DIR, getLoggerSettings(confApp.GetListLogs()))
	if err != nil {
		log.Fatalf("error module 'simplelogger': %v", err)
	}

	//*****************************************************************
	//********** инициализация модуля взаимодействия с Zabbix **********
	channelZabbix := make(chan commoninterfaces.Messager)
	if err := WrappersZabbixInteraction(ctx, *confApp, simpleLogger, channelZabbix); err != nil {
		_, f, l, _ := runtime.Caller(0)
		_ = simpleLogger.WriteLoggingData(fmt.Sprintf(" '%s' %s:%d", err.Error(), f, l-3), "error")

		log.Fatalf("error module 'zabbixinteraction': %v\n", err)
	}

	//******************************************************
	//********** инициализация обработчика логирования данных **********
	logging := logginghandler.New()
	go logginghandler.LoggingHandler(ctx, channelZabbix, simpleLogger, logging.GetChan())

	//******************************************************
	//********** инициализация модуля взаимодействия с TheHive **********
	confTheHiveAPI := confApp.GetApplicationTheHive()
	apiTheHive, err := thehiveapi.New(
		logging,
		thehiveapi.WithAPIKey(confTheHiveAPI.ApiKey),
		thehiveapi.WithHost(confTheHiveAPI.Host),
		thehiveapi.WithPort(confTheHiveAPI.Port))
	if err != nil {
		_, f, l, _ := runtime.Caller(0)
		_ = simpleLogger.WriteLoggingData(fmt.Sprintf(" '%s' %s:%d", err.Error(), f, l-3), "error")

		log.Fatalf("error module 'thehiveapi': %v\n", err)
	}
	chanRequestTheHiveAPI := apiTheHive.Start(ctx)

	//******************************************************
	//********** инициализация модуля взаимодействия с NATS **********
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

	//******************************************************
	//********** инициализация WEBHOOKSERVER модуля **********
	confWebHook := confApp.GetApplicationWebHookServer()
	webHook, chanForSomebody, err := webhookserver.New(
		logging,
		webhookserver.WithTTL(confApp.TTLTmpInfo),
		webhookserver.WithPort(confWebHook.Port),
		webhookserver.WithHost(confWebHook.Host),
		webhookserver.WithName(confWebHook.Name),
		webhookserver.WithVersion(versionandname.GetVersion()))
	if err != nil {
		_, f, l, _ := runtime.Caller(0)
		_ = simpleLogger.WriteLoggingData(fmt.Sprintf(" '%s' %s:%d", err.Error(), f, l-3), "error")

		log.Fatalf("error module 'webhookserver': %v\n", err)
	}

	go router(chanForSomebody, chanRequestTheHiveAPI, chanRequestNatsAPI)

	webHook.Start(ctx)
	webHook.Shutdown(ctx)
}
