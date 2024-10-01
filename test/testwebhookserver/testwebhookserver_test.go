package testwebhookserver_test

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/av-belyakov/simplelogger"
	"github.com/av-belyakov/thehivehook_go_package/cmd/commoninterfaces"
	"github.com/av-belyakov/thehivehook_go_package/cmd/thehiveapi"
	"github.com/av-belyakov/thehivehook_go_package/cmd/webhookserver"
	"github.com/av-belyakov/thehivehook_go_package/cmd/zabbixapi"
	"github.com/av-belyakov/thehivehook_go_package/internal/confighandler"
	"github.com/av-belyakov/thehivehook_go_package/internal/logginghandler"
)

const ROOT_DIR = "thehivehook_go_package"

var _ = Describe("Testwebhookserver", Ordered, func() {
	var (
		rootDir string = "thehivehook_go_package"

		webHookServer *webhookserver.WebHookServer

		conf              *confighandler.ConfigApp
		confTheHiveAPI    *confighandler.AppConfigTheHive
		confWebHookServer *confighandler.AppConfigWebHookServer

		errConf, errServer, errTheHiveAPI error
	)

	BeforeAll(func() {
		//
		// ВАЖНО!!!
		//
		//перед запуском теста установите переменную окружения GO_HIVEHOOK_THAPIKEY
		//с ключем-идентификатором, необходимым для авторизации в API TheHive,
		//командой export GO_HIVEHOOK_THAPIKEY=<api_key>

		conf, errConf = confighandler.NewConfig(rootDir)
		confTheHiveAPI = conf.GetApplicationTheHive()
		confWebHookServer = conf.GetApplicationWebHookServer()
	})

	Context("Тест 1. Чтение конфигурационного файла config_prod.yaml (если не задано GO_HIVEHOOK_MAIN=development)", func() {
		It("При чтение конфигурационного файла не должно быть ошибок", func() {
			Expect(errConf).ShouldNot(HaveOccurred())
		})
	})

	Context("Тест 2. Проверка инициализации модуля TheHiveAPI", func() {
		It("При инициализации модуля не должно быть ошибок", func() {
			Expect(errTheHiveAPI).ShouldNot(HaveOccurred())
		})
	})

	Context("Тест 3. Проверка работы WebHookServer", func() {
		var (
			ctx    context.Context
			cancel context.CancelFunc

			chanForSomebody       <-chan webhookserver.ChanFormWebHookServer
			chanRequestTheHiveAPI chan<- commoninterfaces.ChannelRequester
		)

		BeforeAll(func() {
			ctx, cancel = signal.NotifyContext(context.Background(),
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

			simpleLogger, err := simplelogger.NewSimpleLogger(ROOT_DIR, getLoggerSettings(conf.GetListLogs()))
			if err != nil {
				log.Fatalf("error module 'simplelogger': %v", err)
			}

			channelZabbix := make(chan zabbixapi.MessageSettings)
			go func() {
				for msg := range channelZabbix {
					fmt.Println("INFO for Zabbix:", msg)
				}
			}()

			logging := logginghandler.New()
			go logginghandler.LoggingHandler(ctx, channelZabbix, simpleLogger, logging.GetChan())

			//инициализация модуля взаимодействия с TheHive
			chanRequestTheHiveAPI, errTheHiveAPI = thehiveapi.New(ctx, confTheHiveAPI.ApiKey, confTheHiveAPI.Host, confTheHiveAPI.Port, logging)

			//инициализация webhookserver
			webHookServer, chanForSomebody, errServer = webhookserver.New(ctx, webhookserver.WebHookServerOptions{
				TTL:     confWebHookServer.TTLTmpInfo,
				Port:    confWebHookServer.Port,
				Host:    confWebHookServer.Host,
				Name:    confWebHookServer.Name,
				Version: "1.1.0",
			}, logging)

			go func() {
				for msg := range chanForSomebody {
					switch msg.ForSomebody {
					case "for thehive":
						chanRequestTheHiveAPI <- msg.Data

					case "for nats":
					}
				}
			}()
		})

		It("Ошибок при инициализации сервера быть не должно", func() {
			Expect(errServer).ShouldNot(HaveOccurred())
		})

		It("Работоспособность сервера", func() {
			webHookServer.Start()
			webHookServer.Shutdown(ctx)

			Expect(true).ShouldNot(BeTrue())
		})
	})
})

func getLoggerSettings(cls []confighandler.LogSet) []simplelogger.MessageTypeSettings {
	loggerConf := make([]simplelogger.MessageTypeSettings, 0, len(cls))

	for _, v := range cls {
		loggerConf = append(loggerConf, simplelogger.MessageTypeSettings{
			MsgTypeName:   v.MsgTypeName,
			WritingFile:   v.WritingFile,
			PathDirectory: v.PathDirectory,
			WritingStdout: v.WritingStdout,
			MaxFileSize:   v.MaxFileSize,
		})
	}

	return loggerConf
}
