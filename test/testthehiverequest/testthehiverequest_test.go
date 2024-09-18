package testthehiverequest_test

import (
	"os"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/av-belyakov/thehivehook_go_package/cmd/thehiveapi"
	"github.com/av-belyakov/thehivehook_go_package/internal/confighandler"
)

var _ = Describe("Testthehiverequest", Ordered, func() {
	var (
		rootDir string = "thehivehook_go_package"

		conf       *confighandler.ConfigApp
		theHiveAPI *thehiveapi.TheHiveAPI

		errConf, errTheHiveApi error
	)

	BeforeAll(func() {
		conf, errConf = confighandler.NewConfig(rootDir)
		confTheHive := conf.GetApplicationTheHive()

		//перед запуском теста установите переменную окружения GO_HIVEHOOK_APIKEY
		//с ключем-идентификатором, необходимым для авторизации в API TheHive,
		//командой export GO_HIVEHOOK_APIKEY=<api_key>

		theHiveAPI, errTheHiveApi = thehiveapi.New(os.Getenv("GO_HIVEHOOK_APIKEY"), confTheHive.Host, confTheHive.Port)
	})

	Context("Тест 0. Чтение конфигурационного файла", func() {
		It("При чтении конфигурационного файла не должно быть ошибок", func() {
			Expect(errConf).ShouldNot(HaveOccurred())
		})
	})

	Context("Тест 1. Инициализация модуля взаимодействия с API TheHive", func() {
		It("При инициализации модуля не должно быть ошибок", func() {
			Expect(errTheHiveApi).ShouldNot(HaveOccurred())
		})
	})

	Context("Тест 2. Выполнение запросов к TheHive", func() {
		It("При выполнения запроса на получения кейсов ошибок быть не должно", func() {

		})
	})

	/*
		Context("", func(){
			It("", func(){

			})
		})
	*/
})
