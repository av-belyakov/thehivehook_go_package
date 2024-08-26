package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"runtime/debug"
	"syscall"

	"github.com/lmittmann/tint"

	"github.com/av-belyakov/simplelogger"
	"github.com/av-belyakov/thehivehook_go_package/internal/confighandler"
	"github.com/av-belyakov/thehivehook_go_package/internal/datamodels"
	"github.com/av-belyakov/thehivehook_go_package/internal/supportingfunctions"
	"github.com/av-belyakov/thehivehook_go_package/internal/versionandname"
)

const ROOT_DIR = "thehivehook_go_package"

func main() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	rootPath, err := supportingfunctions.GetRootPath(ROOT_DIR)
	if err != nil {
		log.Fatalf("error, it is impossible to form root path (%w)", err)
	}

	//инициализируем модуль чтения конфигурационного файла
	confApp, err := confighandler.NewConfig(rootPath)
	if err != nil {
		log.Fatalf("error module 'confighandler': %w", err)
	}

	//инициализируем модуль логирования
	sl, err := simplelogger.NewSimpleLogger(ROOT_DIR, getLoggerSettings(confApp.GetListLogs()))
	if err != nil {
		log.Fatalf("error module 'simplelogger': %v", err)
	}

	ctxCore, ctxCancelCore := context.WithCancel(context.Background())

	go func() {
		osCall := <-sigChan
		msg := fmt.Sprintf("stop 'main' function, %s", osCall.String())
		_ = sl.WriteLoggingData(msg, "info")

		ctxCancelCore()
	}()

	loggerColor := slog.New(tint.NewHandler(os.Stdout, &tint.Options{Level: slog.LevelDebug}))

	err := run(loggerColor)
	if err != nil {
		trace := string(debug.Stack())
		loggerColor.Error(err.Error(), "trace", trace)
		os.Exit(1)
	}
}

func run(logger *slog.Logger) error {
	var cfg datamodels.Config

	flag.StringVar(&cfg.baseURL, "base-url", "http://localhost:4444", "base URL for the application")
	flag.IntVar(&cfg.httpPort, "http-port", 4444, "port to listen on for HTTP requests")

	showVersion := flag.Bool("version", false, "display version and exit")

	flag.Parse()

	if *showVersion {
		fmt.Printf("version: %s\n", versionandname.GetVersion())
		return nil
	}

	app := &datamodels.Application{
		config: cfg,
		logger: logger,
	}

	return app.serveHTTP()
}
