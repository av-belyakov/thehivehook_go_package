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

	"github.com/av-belyakov/thehivehook_go_package/cmd/webhookserver"
	"github.com/av-belyakov/thehivehook_go_package/internal/confighandler"
)

var _ = Describe("Testwebhookserver", Ordered, func() {
	var (
		rootDir             string = "thehivehook_go_package"
		theHiveApiKey       string = "70e97faa558d188822c55ec9e00744fd"
		elasticsearchPasswd string = "yD7T27#e28"

		webHookServer *webhookserver.WebHookServer

		conf              *confighandler.ConfigApp
		confWebHookServer *confighandler.AppConfigWebHookServer

		errConf, errServer error
	)

	BeforeAll(func() {
		// это для того что бы тест на чтение конфигурационног файла проходил успешно
		// так как такие паратеры как Passwd для модуля Elasticsearch и ApiKey для модуля
		// TheHive устанавливаются в конфиге приложения только через переме5нные окружения
		os.Setenv("GO_HIVEHOOK_THAPIKEY", theHiveApiKey)
		os.Setenv("GO_HIVEHOOK_ESPASSWD", elasticsearchPasswd)

		conf, errConf = confighandler.NewConfig(rootDir)
		confWebHookServer = conf.GetApplicationWebHookServer()
	})

	Context("Тест 1. Чтение конфигурационного файла config_prod.yaml (если не задано GO_HIVEHOOK_MAIN=development)", func() {
		It("При чтение конфигурационного файла не должно быть ошибок", func() {
			Expect(errConf).ShouldNot(HaveOccurred())
		})
	})

	Context("Тест 2. Проверка работы WebHookServer", func() {
		BeforeAll(func() {
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

			webHookServer, errServer = webhookserver.New(ctx, confWebHookServer.Host, confWebHookServer.Port)
		})

		It("Ошибок при инициализации сервера быть не должно", func() {
			Expect(errServer).ShouldNot(HaveOccurred())
		})

		It("Работоспособность сервера", func() {
			var err error
			if err = webHookServer.Start(); err != nil {

				fmt.Println("____ ERROR:", err.Error())

				if err.Error() == "http: Server closed" {
					err = nil
				}
			}

			Expect(err).ShouldNot(HaveOccurred())
		})
	})

	AfterAll(func() {
		os.Unsetenv("GO_HIVEHOOK_THAPIKEY")
		os.Unsetenv("GO_HIVEHOOK_ESPASSWD")
	})
})
