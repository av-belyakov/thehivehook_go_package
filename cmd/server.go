package main

import (
	"context"
	"log"
	"strings"

	_ "net/http/pprof"

	"github.com/av-belyakov/simplelogger"

	"github.com/av-belyakov/thehivehook_go_package/cmd/elasticsearchapi"
	"github.com/av-belyakov/thehivehook_go_package/cmd/natsapi"
	"github.com/av-belyakov/thehivehook_go_package/cmd/webhookserver"
	"github.com/av-belyakov/thehivehook_go_package/cmd/wrappers"
	"github.com/av-belyakov/thehivehook_go_package/internal/appversion"
	"github.com/av-belyakov/thehivehook_go_package/internal/confighandler"
	"github.com/av-belyakov/thehivehook_go_package/internal/constants"
	"github.com/av-belyakov/thehivehook_go_package/internal/interfaces"
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
	cfg, err := confighandler.NewConfig(rootPath)
	if err != nil {
		log.Fatalf("error module 'confighandler': %s", err.Error())
	}

	// ****************************************************************************
	// ********************* инициализация модуля логирования *********************
	listLog := make([]simplelogger.OptionsManager, 0, len(cfg.GetListLogs()))
	for _, v := range cfg.GetListLogs() {
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
		Port:               cfg.GetApplicationWriteLogDB().Port,
		Host:               cfg.GetApplicationWriteLogDB().Host,
		User:               cfg.GetApplicationWriteLogDB().User,
		Passwd:             cfg.GetApplicationWriteLogDB().Passwd,
		IndexDB:            cfg.GetApplicationWriteLogDB().StorageNameDB,
		NameRegionalObject: cfg.GetApplicationWebHookServer().Name,
	}); err != nil {
		_ = simpleLogger.Write("error", supportingfunctions.CustomError(err).Error())
	} else {
		//подключение логирования в БД
		simpleLogger.SetDataBaseInteraction(esc)
	}

	//******************************************************************
	//********** инициализация модуля взаимодействия с Zabbix **********
	chZabbix := make(chan interfaces.Messager)
	zabbixSettings := wrappers.WrappersZabbixInteractionSettings{
		NetworkPort: cfg.Zabbix.NetworkPort,
		NetworkHost: cfg.Zabbix.NetworkHost,
		ZabbixHost:  cfg.Zabbix.ZabbixHost,
		EventTypes:  make([]wrappers.EventType, len(cfg.Zabbix.EventTypes)),
	}
	for _, v := range cfg.Zabbix.EventTypes {
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

	//***************************************************
	//********** инициализация модуля для взаимодействия с API NATS **********
	natsOptsAPI := []natsapi.NatsApiOptions{
		natsapi.WithHost(cfg.GetApplicationNATS().Host),
		natsapi.WithPort(cfg.GetApplicationNATS().Port),
		natsapi.WithNameRegionalObject(cfg.GetApplicationWebHookServer().Name),
		natsapi.WithSubSenderCase(cfg.GetApplicationNATS().Subscriptions.SenderCase),
		natsapi.WithSubSenderAlert(cfg.GetApplicationNATS().Subscriptions.SenderAlert),
		natsapi.WithSubListenerCommand(cfg.GetApplicationNATS().Subscriptions.ListenerCommand)}
	natsApi, err := natsapi.New(logging, natsOptsAPI...)
	if err != nil {
		_ = simpleLogger.Write("error", supportingfunctions.CustomError(err).Error())
		log.Fatalf("error module 'natsapi': %s\n", err.Error())
	}
	chanInputApiNats, chanOutputApiNats, err := natsApi.Start(ctx)
	if err != nil {
		_ = simpleLogger.Write("error", supportingfunctions.CustomError(err).Error())
		log.Fatalf("error module 'natsapi': %s\n", err.Error())
	}

	//********************************************************************
	//*************** инициализация WEBHOOKSERVER сервера ****************
	webHookServer, chanOutputWebHookServer, err := webhookserver.New(
		logging,
		webhookserver.WithVersion(version),
		webhookserver.WithName(cfg.GetApplicationWebHookServer().Name),
		webhookserver.WithHost(cfg.GetApplicationWebHookServer().Host),
		webhookserver.WithPort(cfg.GetApplicationWebHookServer().Port),
	)
	if err != nil {
		_ = simpleLogger.Write("error", supportingfunctions.CustomError(err).Error())
		log.Fatalf("error module 'webhookserver': %s\n", err.Error())
	}

	//*********************************************************************************
	//********* инициализация маршрутизатора между каналами различных модулей *********
	router := NewMajorRouter(*cfg, logging)
	if err := router.start(
		ctx,
		chanInputApiNats,
		chanOutputWebHookServer,
		chanOutputApiNats,
	); err != nil {
		_ = simpleLogger.Write("error", supportingfunctions.CustomError(err).Error())
	}

	// вывод информационного сообщения при старте приложения
	infoMsg := getInformationMessage(cfg)
	_ = simpleLogger.Write("info", strings.ToLower(infoMsg))

	//запуск модуля
	if err = webHookServer.Start(ctx); err != nil {
		_ = simpleLogger.Write("error", supportingfunctions.CustomError(err).Error())
	}
}
