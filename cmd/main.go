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
	"github.com/av-belyakov/thehivehook_go_package/cmd/zabbixapi"
	"github.com/av-belyakov/thehivehook_go_package/internal/confighandler"
	"github.com/av-belyakov/thehivehook_go_package/internal/logginghandler"
	"github.com/av-belyakov/thehivehook_go_package/internal/supportingfunctions"
)

const ROOT_DIR = "thehivehook_go_package"

func main() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
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

	//инициализируем модуль чтения конфигурационного файла
	confApp, err := confighandler.NewConfig(rootPath)
	if err != nil {
		log.Fatalf("error module 'confighandler': %v", err)
	}

	//инициализируем модуль логирования
	sl, err := simplelogger.NewSimpleLogger(ROOT_DIR, getLoggerSettings(confApp.GetListLogs()))
	if err != nil {
		log.Fatalf("error module 'simplelogger': %v", err)
	}

	//взаимодействие с Zabbix
	channelZabbix := make(chan zabbixapi.MessageSettings)
	if err := interactionZabbix(ctx, confApp, sl, channelZabbix); err != nil {
		_, f, l, _ := runtime.Caller(0)
		_ = sl.WriteLoggingData(fmt.Sprintf(" '%s' %s:%d", err.Error(), f, l-3), "error")

		log.Fatalf("error module 'zabbixinteraction': %v\n", err)
	}

	// логирование данных
	logging := make(chan logginghandler.MessageLogging)
	go logginghandler.LoggingHandler(channelZabbix, sl, logging)

}
