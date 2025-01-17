package main

import (
	"context"
	"fmt"
	"log"
	"runtime"

	"github.com/av-belyakov/simplelogger"
	"github.com/av-belyakov/thehivehook_go_package/cmd/commoninterfaces"
	"github.com/av-belyakov/thehivehook_go_package/cmd/elasticsearchapi"
	"github.com/av-belyakov/thehivehook_go_package/cmd/natsapi"
	"github.com/av-belyakov/thehivehook_go_package/cmd/thehiveapi"
	"github.com/av-belyakov/thehivehook_go_package/cmd/webhookserver"
	"github.com/av-belyakov/thehivehook_go_package/cmd/wrappers"
	"github.com/av-belyakov/thehivehook_go_package/internal/appversion"
	"github.com/av-belyakov/thehivehook_go_package/internal/confighandler"
	"github.com/av-belyakov/thehivehook_go_package/internal/logginghandler"
	"github.com/av-belyakov/thehivehook_go_package/internal/supportingfunctions"
)

// server здесь реализована вся логика запуска thehivehook_go_package
func server(ctx context.Context) {
	rootPath, err := supportingfunctions.GetRootPath(Root_Dir)
	if err != nil {
		log.Fatalf("error, it is impossible to form root path (%s)", err.Error())
	}

	//чтение конфигурационного файла
	confApp, err := confighandler.NewConfig(rootPath)
	if err != nil {
		log.Fatalf("error module 'confighandler': %s", err.Error())
	}

	confWebHook := confApp.GetApplicationWebHookServer()

	//******************************************************
	//********** инициализация модуля логирования **********
	simpleLogger, err := simplelogger.NewSimpleLogger(ctx, Root_Dir, getLoggerSettings(confApp.GetListLogs()))
	if err != nil {
		log.Fatalf("error module 'simplelogger': %s", err.Error())
	}

	//*********************************************************************************
	//********** инициализация модуля взаимодействия с БД для передачи логов **********
	confDB := confApp.GetApplicationWriteLogDB()
	if esc, err := elasticsearchapi.NewElasticsearchConnect(elasticsearchapi.Settings{
		Port:               confDB.Port,
		Host:               confDB.Host,
		User:               confDB.User,
		Passwd:             confDB.Passwd,
		IndexDB:            confDB.StorageNameDB,
		NameRegionalObject: confWebHook.Name,
	}); err != nil {
		_, f, l, _ := runtime.Caller(0)
		_ = simpleLogger.Write("error", fmt.Sprintf(" '%s' %s:%d", err.Error(), f, l-7))

		log.Println(err.Error())
	} else {
		simpleLogger.SetDataBaseInteraction(esc)
	}

	//******************************************************************
	//********** инициализация модуля взаимодействия с Zabbix **********
	channelZabbix := make(chan commoninterfaces.Messager)
	wzis := wrappers.WrappersZabbixInteractionSettings{
		NetworkPort: confApp.Zabbix.NetworkPort,
		NetworkHost: confApp.Zabbix.NetworkHost,
		ZabbixHost:  confApp.Zabbix.ZabbixHost,
	}

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
	wzis.EventTypes = eventTypes
	wrappers.WrappersZabbixInteraction(ctx, wzis, simpleLogger, channelZabbix)

	//******************************************************************
	//********** инициализация обработчика логирования данных **********
	logging := logginghandler.New()
	go logginghandler.LoggingHandler(ctx, simpleLogger, channelZabbix, logging.GetChan())

	//********************************************************
	//********** инициализация TheHive API модуля ************
	confTheHiveAPI := confApp.GetApplicationTheHive()
	apiTheHive, err := thehiveapi.New(
		logging,
		thehiveapi.WithAPIKey(confTheHiveAPI.ApiKey),
		thehiveapi.WithHost(confTheHiveAPI.Host),
		thehiveapi.WithPort(confTheHiveAPI.Port),
		thehiveapi.WithCacheTTL(confTheHiveAPI.CacheTTL))
	if err != nil {
		_, f, l, _ := runtime.Caller(0)
		_ = simpleLogger.Write("error", fmt.Sprintf(" '%s' %s:%d", err.Error(), f, l-3))

		log.Fatalf("error module 'thehiveapi': %s\n", err.Error())
	}
	chReqTheHiveAPI, err := apiTheHive.Start(ctx)
	if err != nil {
		_, f, l, _ := runtime.Caller(0)
		_ = simpleLogger.Write("error", fmt.Sprintf(" '%s' %s:%d", err.Error(), f, l-3))

		log.Fatalf("error module 'thehiveapi': %s\n", err.Error())
	}

	//***************************************************
	//********** инициализация NATS API модуля **********
	confNatsSAPI := confApp.GetApplicationNATS()
	natsOptsAPI := []natsapi.NatsApiOptions{
		natsapi.WithHost(confNatsSAPI.Host),
		natsapi.WithPort(confNatsSAPI.Port),
		natsapi.WithCacheTTL(confNatsSAPI.CacheTTL),
		natsapi.WithSubSenderCase(confNatsSAPI.Subscriptions.SenderCase),
		natsapi.WithSubSenderAlert(confNatsSAPI.Subscriptions.SenderAlert),
		natsapi.WithSubListenerCommand(confNatsSAPI.Subscriptions.ListenerCommand)}
	apiNats, err := natsapi.New(logging, natsOptsAPI...)
	if err != nil {
		_, f, l, _ := runtime.Caller(0)
		_ = simpleLogger.Write("error", fmt.Sprintf(" '%s' %s:%d", err.Error(), f, l-3))

		log.Fatalf("error module 'natsapi': %s\n", err.Error())
	}
	chReqNatsAPI, chNatsAPIReq, err := apiNats.Start(ctx)
	if err != nil {
		_, f, l, _ := runtime.Caller(0)
		_ = simpleLogger.Write("error", fmt.Sprintf(" '%s' %s:%d", err.Error(), f, l-3))

		log.Fatalf("error module 'natsapi': %s\n", err.Error())
	}

	//***********************************************************
	//********** инициализация WEBHOOKSERVER сервера ************
	webHook, chForSomebody, err := webhookserver.New(
		logging,
		webhookserver.WithTTL(confApp.TTLTmpInfo),
		webhookserver.WithPort(confWebHook.Port),
		webhookserver.WithHost(confWebHook.Host),
		webhookserver.WithName(confWebHook.Name),
		webhookserver.WithVersion(appversion.GetVersion()))
	if err != nil {
		_, f, l, _ := runtime.Caller(0)
		_ = simpleLogger.Write("error", fmt.Sprintf(" '%s' %s:%d", err.Error(), f, l-3))

		log.Fatalf("error module 'webhookserver': %s\n", err.Error())
	}

	//мост между каналами различных модулей
	go router(ctx, chForSomebody, chNatsAPIReq, chReqTheHiveAPI, chReqNatsAPI)

	/*
		1. устранить гонку данных
		2. может сделать какой нибудь модуль для контролировании расхода памяти модулей
	*/

	printMemStats()

	/*
			==================
		WARNING: DATA RACE
		Write at 0x00c0002e4d20 by goroutine 829:
		  runtime.mapdelete_faststr()
		      /usr/local/go/src/runtime/map_faststr.go:321 +0x0
		  github.com/av-belyakov/thehivehook_go_package/internal/cacherunningfunctions.(*CacheRunningFunctions).DeleteElement()
		      /home/artemij/go/src/thehivehook_go_package/internal/cacherunningfunctions/methods.go:38 +0xc9
		  github.com/av-belyakov/thehivehook_go_package/internal/cacherunningfunctions.(*CacheRunningFunctions).automaticExecutionMethods.gowrap3()
		      /home/artemij/go/src/thehivehook_go_package/internal/cacherunningfunctions/app.go:62 +0x4f

		Previous read at 0x00c0002e4d20 by goroutine 833:
		  runtime.mapaccess2_faststr()
		      /usr/local/go/src/runtime/map_faststr.go:117 +0x0
		  github.com/av-belyakov/thehivehook_go_package/internal/cacherunningfunctions.(*CacheRunningFunctions).getNumberAttempts()
		      /home/artemij/go/src/thehivehook_go_package/internal/cacherunningfunctions/methods.go:43 +0x184
		  github.com/av-belyakov/thehivehook_go_package/internal/cacherunningfunctions.(*CacheRunningFunctions).automaticExecutionMethods.func2()
		      /home/artemij/go/src/thehivehook_go_package/internal/cacherunningfunctions/app.go:83 +0x1f5
		  github.com/av-belyakov/thehivehook_go_package/internal/cacherunningfunctions.(*CacheRunningFunctions).automaticExecutionMethods.gowrap4()
		      /home/artemij/go/src/thehivehook_go_package/internal/cacherunningfunctions/app.go:89 +0x6e

		Goroutine 829 (running) created at:
		  github.com/av-belyakov/thehivehook_go_package/internal/cacherunningfunctions.(*CacheRunningFunctions).automaticExecutionMethods()
		      /home/artemij/go/src/thehivehook_go_package/internal/cacherunningfunctions/app.go:62 +0x5e8
		  github.com/av-belyakov/thehivehook_go_package/internal/cacherunningfunctions.CreateCache.gowrap1()
		      /home/artemij/go/src/thehivehook_go_package/internal/cacherunningfunctions/app.go:32 +0x4f

		Goroutine 833 (running) created at:
		  github.com/av-belyakov/thehivehook_go_package/internal/cacherunningfunctions.(*CacheRunningFunctions).automaticExecutionMethods()
		      /home/artemij/go/src/thehivehook_go_package/internal/cacherunningfunctions/app.go:74 +0x52a
		  github.com/av-belyakov/thehivehook_go_package/internal/cacherunningfunctions.CreateCache.gowrap1()
		      /home/artemij/go/src/thehivehook_go_package/internal/cacherunningfunctions/app.go:32 +0x4f
		==================
	*/

	if err = webHook.Start(ctx); err != nil {
		_, f, l, _ := runtime.Caller(0)
		_ = simpleLogger.Write("error", fmt.Sprintf("'%s' %s:%d", err.Error(), f, l-1))
		log.Fatalln(err)
	}
}
