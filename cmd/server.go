package main

import (
	"context"
	"log"
	"os"

	"net/http"
	_ "net/http/pprof"

	"github.com/av-belyakov/simplelogger"

	"github.com/av-belyakov/thehivehook_go_package/cmd/commoninterfaces"
	"github.com/av-belyakov/thehivehook_go_package/cmd/constants"
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
	version, err := appversion.GetAppVersion()
	if err != nil {
		log.Println(err)
	}

	rootPath, err := supportingfunctions.GetRootPath(constants.Root_Dir)
	if err != nil {
		log.Fatalf("error, it is impossible to form root path (%s)", err.Error())
	}

	//чтение конфигурационного файла
	confApp, err := confighandler.NewConfig(rootPath)
	if err != nil {
		log.Fatalf("error module 'confighandler': %s", err.Error())
	}

	confWebHook := confApp.GetApplicationWebHookServer()

	// ****************************************************************************
	// ********************* инициализация модуля логирования *********************
	var listLog []simplelogger.OptionsManager
	for _, v := range confApp.GetListLogs() {
		listLog = append(listLog, v)
	}
	opts := simplelogger.CreateOptions(listLog...)
	simpleLogger, err := simplelogger.NewSimpleLogger(ctx, constants.Root_Dir, opts)
	if err != nil {
		log.Fatalf("error module 'simplelogger': %v", err)
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
		_ = simpleLogger.Write("error", supportingfunctions.CustomError(err).Error())
	} else {
		//подключение логирования в БД
		simpleLogger.SetDataBaseInteraction(esc)
	}

	//******************************************************************
	//********** инициализация модуля взаимодействия с Zabbix **********
	chZabbix := make(chan commoninterfaces.Messager)
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
	wzSettings := wrappers.WrappersZabbixInteractionSettings{
		NetworkPort: confApp.Zabbix.NetworkPort,
		NetworkHost: confApp.Zabbix.NetworkHost,
		ZabbixHost:  confApp.Zabbix.ZabbixHost,
		EventTypes:  eventTypes,
	}
	wrappers.WrappersZabbixInteraction(ctx, wzSettings, simpleLogger, chZabbix)

	//******************************************************************
	//********** инициализация обработчика логирования данных **********
	logging := logginghandler.New(simpleLogger, chZabbix)
	logging.Start(ctx)

	//******************************************************************
	//************** инициализация TheHive API модуля ******************
	confTheHiveAPI := confApp.GetApplicationTheHive()
	apiTheHive, err := thehiveapi.New(
		logging,
		thehiveapi.WithAPIKey(confTheHiveAPI.ApiKey),
		thehiveapi.WithHost(confTheHiveAPI.Host),
		thehiveapi.WithPort(confTheHiveAPI.Port),
		thehiveapi.WithCacheTTL(confTheHiveAPI.CacheTTL))
	if err != nil {
		_ = simpleLogger.Write("error", supportingfunctions.CustomError(err).Error())
		log.Fatalf("error module 'thehiveapi': %s\n", err.Error())
	}
	//запуск модуля
	chReqTheHiveAPI, err := apiTheHive.Start(ctx)
	if err != nil {
		_ = simpleLogger.Write("error", supportingfunctions.CustomError(err).Error())
		log.Fatalf("error module 'thehiveapi': %s\n", err.Error())
	}

	//***************************************************
	//********** инициализация NATS API модуля **********
	confNatsSAPI := confApp.GetApplicationNATS()
	natsOptsAPI := []natsapi.NatsApiOptions{
		natsapi.WithHost(confNatsSAPI.Host),
		natsapi.WithPort(confNatsSAPI.Port),
		natsapi.WithCacheTTL(confNatsSAPI.CacheTTL),
		natsapi.WithNameRegionalObject(confWebHook.Name),
		natsapi.WithSubSenderCase(confNatsSAPI.Subscriptions.SenderCase),
		natsapi.WithSubSenderAlert(confNatsSAPI.Subscriptions.SenderAlert),
		natsapi.WithSubListenerCommand(confNatsSAPI.Subscriptions.ListenerCommand)}
	apiNats, err := natsapi.New(logging, natsOptsAPI...)
	if err != nil {
		_ = simpleLogger.Write("error", supportingfunctions.CustomError(err).Error())
		log.Fatalf("error module 'natsapi': %s\n", err.Error())
	}
	//запуск модуля
	chReqNatsAPI, chNatsAPIReq, err := apiNats.Start(ctx)
	if err != nil {
		_ = simpleLogger.Write("error", supportingfunctions.CustomError(err).Error())
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
		webhookserver.WithVersion(version))
	if err != nil {
		_ = simpleLogger.Write("error", supportingfunctions.CustomError(err).Error())
		log.Fatalf("error module 'webhookserver': %s\n", err.Error())
	}

	//мост между каналами различных модулей
	router(ctx, chForSomebody, chNatsAPIReq, chReqTheHiveAPI, chReqNatsAPI)

	//для отладки через pprof
	//http://localhost:6060/debug/pprof/
	//go tool pprof http://localhost:6060/debug/pprof/heap
	//go tool pprof http://localhost:6060/debug/pprof/goroutine
	//go tool pprof http://localhost:6060/debug/pprof/allocs
	if os.Getenv("GO_HIVEHOOK_MAIN") == "development" {
		go func() {
			log.Println(http.ListenAndServe("localhost:6060", nil))
		}()
	}
	//------------------------------------------

	//запуск модуля
	if err = webHook.Start(ctx); err != nil {
		_ = simpleLogger.Write("error", supportingfunctions.CustomError(err).Error())
	}
}
