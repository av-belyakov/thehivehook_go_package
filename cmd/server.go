package main

import (
	"context"
	"log"
	"strings"

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
	conf, err := confighandler.NewConfig(rootPath)
	if err != nil {
		log.Fatalf("error module 'confighandler': %s", err.Error())
	}

	// ****************************************************************************
	// ********************* инициализация модуля логирования *********************
	listLog := make([]simplelogger.OptionsManager, 0, len(conf.GetListLogs()))
	for _, v := range conf.GetListLogs() {
		listLog = append(listLog, v)
	}
	opts := simplelogger.CreateOptions(listLog...)
	simpleLogger, err := simplelogger.NewSimpleLogger(ctx, constants.Root_Dir, opts)
	if err != nil {
		log.Fatalf("error module 'simplelogger': %v", err)
	}

	//*********************************************************************************
	//********** инициализация модуля взаимодействия с БД для передачи логов **********
	if esc, err := elasticsearchapi.NewElasticsearchConnect(elasticsearchapi.Settings{
		Port:               conf.GetApplicationWriteLogDB().Port,
		Host:               conf.GetApplicationWriteLogDB().Host,
		User:               conf.GetApplicationWriteLogDB().User,
		Passwd:             conf.GetApplicationWriteLogDB().Passwd,
		IndexDB:            conf.GetApplicationWriteLogDB().StorageNameDB,
		NameRegionalObject: conf.GetApplicationWebHookServer().Name,
	}); err != nil {
		_ = simpleLogger.Write("error", supportingfunctions.CustomError(err).Error())
	} else {
		//подключение логирования в БД
		simpleLogger.SetDataBaseInteraction(esc)
	}

	//******************************************************************
	//********** инициализация модуля взаимодействия с Zabbix **********
	chZabbix := make(chan commoninterfaces.Messager)
	zabbixSettings := wrappers.WrappersZabbixInteractionSettings{
		NetworkPort: conf.Zabbix.NetworkPort,
		NetworkHost: conf.Zabbix.NetworkHost,
		ZabbixHost:  conf.Zabbix.ZabbixHost,
		EventTypes:  make([]wrappers.EventType, len(conf.Zabbix.EventTypes)),
	}
	for _, v := range conf.Zabbix.EventTypes {
		zabbixSettings.EventTypes = append(zabbixSettings.EventTypes, wrappers.EventType{
			IsTransmit: v.IsTransmit,
			EventType:  v.EventType,
			ZabbixKey:  v.ZabbixKey,
			Handshake: wrappers.Handshake{
				TimeInterval: v.Handshake.TimeInterval,
				Message:      v.Handshake.Message,
			},
		})
	}
	wrappers.WrappersZabbixInteraction(ctx, zabbixSettings, simpleLogger, chZabbix)

	//******************************************************************
	//********** инициализация обработчика логирования данных **********
	logging := logginghandler.New(simpleLogger, chZabbix)
	logging.Start(ctx)

	//******************************************************************
	//************** инициализация TheHive API модуля ******************
	apiTheHive, err := thehiveapi.New(
		logging,
		thehiveapi.WithAPIKey(conf.GetApplicationTheHive().ApiKey),
		thehiveapi.WithHost(conf.GetApplicationTheHive().Host),
		thehiveapi.WithPort(conf.GetApplicationTheHive().Port),
		thehiveapi.WithCacheTTL(conf.GetApplicationTheHive().CacheTTL),
		thehiveapi.WithNameRegionalObject(conf.GetApplicationWebHookServer().Name))
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
	natsOptsAPI := []natsapi.NatsApiOptions{
		natsapi.WithHost(conf.GetApplicationNATS().Host),
		natsapi.WithPort(conf.GetApplicationNATS().Port),
		natsapi.WithCacheTTL(conf.GetApplicationNATS().CacheTTL),
		natsapi.WithNameRegionalObject(conf.GetApplicationWebHookServer().Name),
		natsapi.WithSubSenderCase(conf.GetApplicationNATS().Subscriptions.SenderCase),
		natsapi.WithSubSenderAlert(conf.GetApplicationNATS().Subscriptions.SenderAlert),
		natsapi.WithSubListenerCommand(conf.GetApplicationNATS().Subscriptions.ListenerCommand)}
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
		webhookserver.WithTTL(conf.GetApplicationWebHookServer().TTLTmpInfo),
		webhookserver.WithPort(conf.GetApplicationWebHookServer().Port),
		webhookserver.WithHost(conf.GetApplicationWebHookServer().Host),
		webhookserver.WithName(conf.GetApplicationWebHookServer().Name),
		webhookserver.WithVersion(version))
	if err != nil {
		_ = simpleLogger.Write("error", supportingfunctions.CustomError(err).Error())
		log.Fatalf("error module 'webhookserver': %s\n", err.Error())
	}

	//мост между каналами различных модулей
	router(ctx, chForSomebody, chNatsAPIReq, chReqTheHiveAPI, chReqNatsAPI)

	// вывод информационного сообщения при старте приложения
	infoMsg := getInformationMessage(conf)
	_ = simpleLogger.Write("info", strings.ToLower(infoMsg))

	//запуск модуля
	if err = webHook.Start(ctx); err != nil {
		_ = simpleLogger.Write("error", supportingfunctions.CustomError(err).Error())
	}
}
