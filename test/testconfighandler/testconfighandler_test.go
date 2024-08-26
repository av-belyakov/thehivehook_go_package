package testconfighandler_test

import (
	"os"
	"strconv"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/av-belyakov/thehivehook_go_package/internal/confighandler"
)

var _ = Describe("Testconfighandler", Ordered, func() {
	var (
		rootDir       string = "thehivehook_go_package"
		theHiveApiKey string = "70e97faa558d188822c55ec9e00744fd"

		conf *confighandler.ConfigApp

		err error
	)

	unSetEnvAny := func() {
		os.Unsetenv("GO_HIVEHOOK_MAIN")
		os.Unsetenv("GO_HIVEHOOK_NHOST")
		os.Unsetenv("GO_HIVEHOOK_NPORT")
		os.Unsetenv("GO_HIVEHOOK_SUBJECTCASE")
		os.Unsetenv("GO_HIVEHOOK_SUBJECTALERT")

		os.Unsetenv("GO_HIVEHOOK_THHOST")
		os.Unsetenv("GO_HIVEHOOK_THPORT")
		os.Unsetenv("GO_HIVEHOOK_THUNAME")
	}

	BeforeAll(func() {
		os.Setenv("GO_HIVEHOOK_THAPIKEY", theHiveApiKey)
	})

	AfterAll(func() {
		os.Unsetenv("GO_HIVEHOOK_THAPIKEY")
	})

	BeforeEach(func() {
		unSetEnvAny()
	})

	Context("Тест 1. Чтение конфигурационного файла (по умолчанию config_prod.yaml)", func() {
		BeforeAll(func() {
			conf, err = confighandler.NewConfig(rootDir)
			Expect(err).ShouldNot(HaveOccurred())
		})

		It("Должен быть прочитан файл config_prod", func() {
			cinfo := conf.GetCommonInfo()
			Expect(cinfo.FileName).Should(Equal("config_prod"))
		})

		It("Все пораметры конфигрурационного файла 'config_prod.yaml' для NATS должны быть успешно получены", func() {
			cn := conf.GetApplicationNATS()
			Expect(cn.Host).Should(Equal("nats.cloud.gcm"))
			Expect(cn.Port).Should(Equal(4222))
			Expect(cn.SubjectCase).Should(Equal("main_caseupdate"))
			Expect(cn.SubjectAlert).Should(Equal("main_alertupdate"))
		})

		It("Все пораметры конфигрурационного файла 'config_prod.yaml' для THEHIVE должны быть успешно получены", func() {
			cth := conf.GetApplicationTheHive()
			Expect(cth.Host).Should(Equal("192.168.42.10"))
			Expect(cth.Port).Should(Equal(9000))
			Expect(cth.UserName).Should(Equal("test"))
			Expect(cth.ApiKey).Should(Equal(theHiveApiKey))
		})
	})

	Context("Тест 2. Чтение конфигурационного файла config_dev.yaml", func() {
		BeforeAll(func() {
			os.Setenv("GO_HIVEHOOK_MAIN", "development")

			conf, err = confighandler.NewConfig(rootDir)
			Expect(err).ShouldNot(HaveOccurred())
		})

		It("Должен быть прочитан файл config_dev", func() {
			cinfo := conf.GetCommonInfo()
			Expect(cinfo.FileName).Should(Equal("config_dev"))
		})

		It("Все пораметры конфигурационного файла 'config_dev.yaml' для NATS должны быть успешно получены", func() {
			cn := conf.GetApplicationNATS()
			Expect(cn.Host).Should(Equal("nats.cloud.gcmtest"))
			Expect(cn.Port).Should(Equal(4223))
			Expect(cn.SubjectCase).Should(Equal("main_caseupdate_test"))
			Expect(cn.SubjectAlert).Should(Equal("main_alertupdate_test"))
		})

		It("Все пораметры конфигрурационного файла 'config_dev.yaml' для THEHIVE должны быть успешно получены", func() {
			cth := conf.GetApplicationTheHive()
			Expect(cth.Host).Should(Equal("192.168.42.10"))
			Expect(cth.Port).Should(Equal(9001))
			Expect(cth.UserName).Should(Equal("testtest"))
			Expect(cth.ApiKey).Should(Equal(theHiveApiKey))
		})
	})

	Context("Тест 3. Проверяем установленные для NATS значения переменных окружения", func() {
		const (
			NATS_HOST         = "nats.cloud.gcm.test.test"
			NATS_PORT         = 4545
			NATS_SUBJECTCASE  = "main_CASE_update"
			NATS_SUBJECTALERT = "main_ALERT_update"
		)

		BeforeAll(func() {
			os.Setenv("GO_HIVEHOOK_NHOST", NATS_HOST)
			os.Setenv("GO_HIVEHOOK_NPORT", strconv.Itoa(NATS_PORT))
			os.Setenv("GO_HIVEHOOK_SUBJECTCASE", NATS_SUBJECTCASE)
			os.Setenv("GO_HIVEHOOK_SUBJECTALERT", NATS_SUBJECTALERT)

			conf, err = confighandler.NewConfig(rootDir)
			Expect(err).ShouldNot(HaveOccurred())
		})

		It("Все параметры конфигурации для NATS должны быть успешно установлены через соответствующие переменные окружения", func() {
			cn := conf.GetApplicationNATS()

			Expect(cn.Host).Should(Equal(NATS_HOST))
			Expect(cn.Port).Should(Equal(NATS_PORT))
			Expect(cn.SubjectCase).Should(Equal(NATS_SUBJECTCASE))
			Expect(cn.SubjectAlert).Should(Equal(NATS_SUBJECTALERT))
		})
	})

	Context("Тест 4. Проверяем установленные для THEHIVE значения переменных окружения", func() {
		const (
			THEHIVE_HOST    = "thehive.cloud.gcm.test"
			THEHIVE_PORT    = 1122
			THEHIVE_THUNAME = "test_hive_name"
		)

		BeforeAll(func() {
			os.Setenv("GO_HIVEHOOK_THHOST", THEHIVE_HOST)
			os.Setenv("GO_HIVEHOOK_THPORT", strconv.Itoa(THEHIVE_PORT))
			os.Setenv("GO_HIVEHOOK_THUNAME", THEHIVE_THUNAME)

			conf, err = confighandler.NewConfig(rootDir)
			Expect(err).ShouldNot(HaveOccurred())
		})

		It("Все параметры конфигурации для TheHive должны быть успешно установлены через соответствующие переменные окружения", func() {
			cth := conf.GetApplicationTheHive()

			Expect(cth.Host).Should(Equal(THEHIVE_HOST))
			Expect(cth.Port).Should(Equal(THEHIVE_PORT))
			Expect(cth.UserName).Should(Equal(THEHIVE_THUNAME))
			Expect(cth.ApiKey).Should(Equal(theHiveApiKey))
		})
	})

	Context("Тест 5. Проверяем обработку файда 'config.yaml'", func() {
		It("Должно быть получено содержимое общего файла 'config.yaml'", func() {
			confApp, err := confighandler.NewConfig(rootDir)
			Expect(err).ShouldNot(HaveOccurred())

			commonApp := confApp.GetCommonApplication()

			//*** настройки Zabbix ***
			Expect(commonApp.Zabbix.NetworkHost).Should(Equal("192.168.9.45"))
			Expect(commonApp.Zabbix.NetworkPort).Should(Equal(10051))
			Expect(commonApp.Zabbix.ZabbixHost).Should(Equal("test-uchet-db.cloud.gcm"))
			Expect(len(commonApp.Zabbix.EventTypes)).Should(Equal(3))

			Expect(commonApp.Zabbix.EventTypes[0].EventType).Should(Equal("error"))
			Expect(commonApp.Zabbix.EventTypes[0].ZabbixKey).Should(Equal("shaper_stix.error"))
			Expect(commonApp.Zabbix.EventTypes[0].IsTransmit).Should(BeTrue())
			Expect(commonApp.Zabbix.EventTypes[0].Handshake.TimeInterval).Should(Equal(0))
			Expect(commonApp.Zabbix.EventTypes[0].Handshake.Message).Should(Equal(""))
			Expect(commonApp.Zabbix.EventTypes[1].EventType).Should(Equal("info"))
			Expect(commonApp.Zabbix.EventTypes[1].ZabbixKey).Should(Equal("shaper_stix.info"))
			Expect(commonApp.Zabbix.EventTypes[1].IsTransmit).Should(BeTrue())
			Expect(commonApp.Zabbix.EventTypes[2].EventType).Should(Equal("handshake"))
			Expect(commonApp.Zabbix.EventTypes[2].ZabbixKey).Should(Equal("shaper_stix.handshake"))
			Expect(commonApp.Zabbix.EventTypes[2].IsTransmit).Should(BeTrue())

			//*** настройки логирования ***
			Expect(len(confApp.GetListLogs())).Should(Equal(3))
		})
	})

	/*Context("Тест 4. Проверяем работу функции NewConfig с разными значениями переменной окружения GO_HIVEHOOK_MAIN", func() {
		It("", func() {

		})
	})*/
})
