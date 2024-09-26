package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"github.com/av-belyakov/simplelogger"
	"github.com/av-belyakov/thehivehook_go_package/cmd/webhookserver"
	"github.com/av-belyakov/thehivehook_go_package/cmd/zabbixapi"
	"github.com/av-belyakov/thehivehook_go_package/internal/confighandler"
	"github.com/av-belyakov/thehivehook_go_package/internal/logginghandler"
	"github.com/av-belyakov/thehivehook_go_package/internal/supportingfunctions"
)

const ROOT_DIR = "thehivehook_go_package"

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	go func() {
		sigChan := make(chan os.Signal, 1)
		osCall := <-sigChan
		log.Printf("system call:%+v", osCall)

		cancel()
	}()

	server(ctx)
}

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

		log.Fatalf("error module 'zabbixinteraction': %v\n", err)
	}

	//********** инициализация обработчика логирования данных **********
	logging := logginghandler.New()
	go logginghandler.LoggingHandler(ctx, channelZabbix, simpleLogger, logging.GetChan())

	//********** инициализация модуля взаимодействия с TheHive **********

	//********** инициализация модуля взаимодействия с NATS **********

	//********** инициализация WEBHOOKSERVER модуля **********
	confWebHook := confApp.GetApplicationWebHookServer()
	webHook, err := webhookserver.New(ctx, confWebHook.Name, confWebHook.Host, confWebHook.Port, confWebHook.TTLTmpInfo, logging)
	if err != nil {
		_, f, l, _ := runtime.Caller(0)
		_ = simpleLogger.WriteLoggingData(fmt.Sprintf(" '%s' %s:%d", err.Error(), f, l-3), "error")

		log.Fatalf("error module 'webhookserver': %v\n", err)
	}
	webHook.Start()
	webHook.Shutdown(ctx)
}
