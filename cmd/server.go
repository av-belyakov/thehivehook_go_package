package main

import (
	"context"
	"fmt"
	"log"
	"runtime"

	"github.com/av-belyakov/simplelogger"
	"github.com/av-belyakov/thehivehook_go_package/cmd/thehiveapi"
	"github.com/av-belyakov/thehivehook_go_package/cmd/webhookserver"
	"github.com/av-belyakov/thehivehook_go_package/cmd/zabbixapi"
	"github.com/av-belyakov/thehivehook_go_package/internal/confighandler"
	"github.com/av-belyakov/thehivehook_go_package/internal/logginghandler"
	"github.com/av-belyakov/thehivehook_go_package/internal/supportingfunctions"
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

	//********** инициализация модуля логирования **********
	simpleLogger, err := simplelogger.NewSimpleLogger(ROOT_DIR, getLoggerSettings(confApp.GetListLogs()))
	if err != nil {
		log.Fatalf("error module 'simplelogger': %v", err)
	}

	//********** инициализация модуля взаимодействия с Zabbix **********
	channelZabbix := make(chan zabbixapi.MessageSettings)
	if err := interactionZabbix(ctx, confApp, simpleLogger, channelZabbix); err != nil {
		_, f, l, _ := runtime.Caller(0)
		_ = simpleLogger.WriteLoggingData(fmt.Sprintf(" '%s' %s:%d", err.Error(), f, l-3), "error")

		//
		// здесь надо сделать возможность отключение передачи данных в zabbix
		// для этого в конфиге надо добавить дополнительную опцию или сделать
		// не критичным ошибки модуля interactionZabbix
		//

		log.Fatalf("error module 'zabbixinteraction': %v\n", err)
	}

	//********** инициализация обработчика логирования данных **********
	logging := logginghandler.New()
	go logginghandler.LoggingHandler(ctx, channelZabbix, simpleLogger, logging.GetChan())

	//********** инициализация модуля взаимодействия с TheHive **********
	confTheHiveAPI := confApp.GetApplicationTheHive()
	chanRequestTheHiveAPI, err := thehiveapi.New(ctx, confTheHiveAPI.ApiKey, confTheHiveAPI.Host, confTheHiveAPI.Port, logging)
	if err != nil {
		_, f, l, _ := runtime.Caller(0)
		_ = simpleLogger.WriteLoggingData(fmt.Sprintf(" '%s' %s:%d", err.Error(), f, l-3), "error")

		log.Fatalf("error module 'thehiveapi': %v\n", err)
	}

	//********** инициализация модуля взаимодействия с NATS **********

	//********** инициализация WEBHOOKSERVER модуля **********
	confWebHook := confApp.GetApplicationWebHookServer()
	webHook, chanForSomebody, err := webhookserver.New(
		ctx,
		logging,
		webhookserver.WithTTL(confApp.TTLTmpInfo),
		webhookserver.WithPort(confWebHook.Port),
		webhookserver.WithHost(confWebHook.Host),
		webhookserver.WithName(confWebHook.Name),
		webhookserver.WithVersion("0.1.1"))
	if err != nil {
		_, f, l, _ := runtime.Caller(0)
		_ = simpleLogger.WriteLoggingData(fmt.Sprintf(" '%s' %s:%d", err.Error(), f, l-3), "error")

		log.Fatalf("error module 'webhookserver': %v\n", err)
	}

	go router(chanForSomebody, chanRequestTheHiveAPI)

	webHook.Start()
	webHook.Shutdown(ctx)
}
