package testconfighandler_test

import (
	"fmt"
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

		//настройки NATS
		os.Unsetenv("GO_HIVEHOOK_NPREFIX")
		os.Unsetenv("GO_HIVEHOOK_NHOST")
		os.Unsetenv("GO_HIVEHOOK_NPORT")
		os.Unsetenv("GO_HIVEHOOK_NCACHETTL")
		os.Unsetenv("GO_HIVEHOOK_NSUBSCRIBERS")

		//настройки TheHive
		os.Unsetenv("GO_HIVEHOOK_THHOST")
		os.Unsetenv("GO_HIVEHOOK_THPORT")
		os.Unsetenv("GO_HIVEHOOK_THCACHETTL")

		//настройки WebHook сервера
		os.Unsetenv("GO_HIVEHOOK_WEBHNAME")
		os.Unsetenv("GO_HIVEHOOK_WEBHHOST")
		os.Unsetenv("GO_HIVEHOOK_WEBHPORT")
		os.Unsetenv("GO_HIVEHOOK_WEBHTTLTMPINFO")
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
			if err != nil {
				fmt.Println("ERROR:", err)
			}
			Expect(err).ShouldNot(HaveOccurred())
		})

		It("Должен быть прочитан файл config_prod", func() {
			cinfo := conf.GetCommonInfo()
			Expect(cinfo.FileName).Should(Equal("config_prod"))
		})

		It("Все пораметры конфигурационного файла 'config_prod.yaml' для NATS должны быть успешно получены", func() {
			cn := conf.GetApplicationNATS()

			fmt.Println("Application NATS config:")
			fmt.Println(cn)
			Expect(cn.Prefix).Should(Equal("test"))
			Expect(cn.Host).Should(Equal("nats.cloud.gcm"))
			Expect(cn.Port).Should(Equal(4222))
			Expect(cn.CacheTTL).Should(Equal(43200))
			Expect(len(cn.Subscribers)).Should(Equal(2))
			for _, v := range cn.Subscribers {
				if v.Event == "caseupdate" {
					Expect(len(v.Responders)).Should(Equal(5))
				}
				if v.Event == "alertupdate" {
					Expect(len(v.Responders)).Should(Equal(1))
				}
			}
		})

		It("Все пораметры конфигурационного файла 'config_prod.yaml' для THEHIVE должны быть успешно получены", func() {
			cth := conf.GetApplicationTheHive()
			Expect(cth.Host).Should(Equal("thehive.cloud.gcm"))
			Expect(cth.Port).Should(Equal(9000))
			Expect(cth.CacheTTL).Should(Equal(43200))
			Expect(cth.ApiKey).Should(Equal(theHiveApiKey))
		})

		It("Все пораметры конфигрурационного файла 'config_prod.yaml' для WEBHOOKSERVER должны быть успешно получены", func() {
			chs := conf.GetApplicationWebHookServer()
			Expect(chs.Host).Should(Equal("192.168.13.3"))
			Expect(chs.Port).Should(Equal(5000))
			Expect(chs.TTLTmpInfo).Should(Equal(10))
			Expect(chs.Name).Should(Equal("gcm"))
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
			Expect(cn.Prefix).Should(Equal("test"))
			Expect(cn.Host).Should(Equal("nats.cloud.gcmtest"))
			Expect(cn.Port).Should(Equal(4223))
			Expect(cn.CacheTTL).Should(Equal(3600))
			Expect(len(cn.Subscribers)).Should(Equal(2))
			for _, v := range cn.Subscribers {
				if v.Event == "caseupdate" {
					Expect(len(v.Responders)).Should(Equal(5))
				}
				if v.Event == "alertupdate" {
					Expect(len(v.Responders)).Should(Equal(3))
				}
			}
		})

		It("Все пораметры конфигрурационного файла 'config_dev.yaml' для THEHIVE должны быть успешно получены", func() {
			cth := conf.GetApplicationTheHive()
			Expect(cth.Host).Should(Equal("thehive.cloud.gcm"))
			Expect(cth.Port).Should(Equal(9001))
			Expect(cth.CacheTTL).Should(Equal(3600))
			Expect(cth.ApiKey).Should(Equal(theHiveApiKey))
		})

		It("Все пораметры конфигрурационного файла 'config_dev.yaml' для WEBHOOKSERVER должны быть успешно получены", func() {
			chs := conf.GetApplicationWebHookServer()
			Expect(chs.Host).Should(Equal("127.0.0.1"))
			Expect(chs.Port).Should(Equal(5000))
			Expect(chs.TTLTmpInfo).Should(Equal(12))
			Expect(chs.Name).Should(Equal("rcmsml"))
		})
	})

	Context("Тест 3. Проверяем установленные для NATS значения переменных окружения", func() {
		const (
			NATS_PREFIX      = "main"
			NATS_HOST        = "nats.cloud.gcm.test.test"
			NATS_PORT        = 4545
			NATS_CACHETTL    = 3600
			NATS_SUBSCRIBERS = "main_caseupdate:gcm,rcmmsk,rcmnvs;main_alertupdate:gcm"
		)

		BeforeAll(func() {
			os.Setenv("GO_HIVEHOOK_NPREFIX", NATS_PREFIX)
			os.Setenv("GO_HIVEHOOK_NHOST", NATS_HOST)
			os.Setenv("GO_HIVEHOOK_NPORT", strconv.Itoa(NATS_PORT))
			os.Setenv("GO_HIVEHOOK_NCACHETTL", strconv.Itoa(NATS_CACHETTL))
			os.Setenv("GO_HIVEHOOK_NSUBSCRIBERS", NATS_SUBSCRIBERS)

			conf, err = confighandler.NewConfig(rootDir)
			Expect(err).ShouldNot(HaveOccurred())
		})

		It("Все параметры конфигурации для NATS должны быть успешно установлены через соответствующие переменные окружения", func() {
			cn := conf.GetApplicationNATS()

			Expect(cn.Prefix).Should(Equal(NATS_PREFIX))
			Expect(cn.Host).Should(Equal(NATS_HOST))
			Expect(cn.Port).Should(Equal(NATS_PORT))
			for _, v := range cn.Subscribers {
				if v.Event == "main_caseupdate" {
					Expect(len(v.Responders)).Should(Equal(3))
				}
				if v.Event == "main_alertupdate" {
					Expect(len(v.Responders)).Should(Equal(1))
				}
			}
		})
	})

	Context("Тест 4. Проверяем установленные для THEHIVE значения переменных окружения", func() {
		const (
			THEHIVE_HOST     = "thehive.cloud.gcm.test"
			THEHIVE_PORT     = 1122
			THEHIVE_CACHETTL = 3636
			THEHIVE_THUNAME  = "test_hive_name"
		)

		BeforeAll(func() {
			os.Setenv("GO_HIVEHOOK_THHOST", THEHIVE_HOST)
			os.Setenv("GO_HIVEHOOK_THPORT", strconv.Itoa(THEHIVE_PORT))
			os.Setenv("GO_HIVEHOOK_THCACHETTL", strconv.Itoa(THEHIVE_CACHETTL))

			conf, err = confighandler.NewConfig(rootDir)
			Expect(err).ShouldNot(HaveOccurred())
		})

		It("Все параметры конфигурации для TheHive должны быть успешно установлены через соответствующие переменные окружения", func() {
			cth := conf.GetApplicationTheHive()

			Expect(cth.Host).Should(Equal(THEHIVE_HOST))
			Expect(cth.Port).Should(Equal(THEHIVE_PORT))
			Expect(cth.ApiKey).Should(Equal(theHiveApiKey))
		})
	})

	Context("Тест 5. Проверяем установленные для WEBHOOKSERVER значения переменных окружения", func() {
		const (
			HIVEHOOK_WEBHHOST = "11.0.11.10"
			HIVEHOOK_WEBHPORT = 7822
			HIVEHOOK_WEBTTL   = 13
			HIVEHOOK_WEBHNAME = "gcm-rcm"
		)

		BeforeAll(func() {
			os.Setenv("GO_HIVEHOOK_WEBHNAME", HIVEHOOK_WEBHNAME)
			os.Setenv("GO_HIVEHOOK_WEBHHOST", HIVEHOOK_WEBHHOST)
			os.Setenv("GO_HIVEHOOK_WEBHPORT", strconv.Itoa(HIVEHOOK_WEBHPORT))
			os.Setenv("GO_HIVEHOOK_WEBHTTLTMPINFO", strconv.Itoa(HIVEHOOK_WEBTTL))

			conf, err = confighandler.NewConfig(rootDir)
			Expect(err).ShouldNot(HaveOccurred())
		})

		It("Все параметры конфигурации для WEBHOOKSERVER должны быть успешно установлены через соответствующие переменные окружения", func() {
			whookserver := conf.GetApplicationWebHookServer()

			Expect(whookserver.Name).Should(Equal(HIVEHOOK_WEBHNAME))
			Expect(whookserver.Host).Should(Equal(HIVEHOOK_WEBHHOST))
			Expect(whookserver.Port).Should(Equal(HIVEHOOK_WEBHPORT))
			Expect(whookserver.TTLTmpInfo).Should(Equal(HIVEHOOK_WEBTTL))
		})
	})

	Context("Тест 6. Проверяем обработку файда 'config.yaml'", func() {
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
			Expect(len(confApp.GetListLogs())).Should(Equal(4))
		})
	})

	/*Context("Тест 4. Проверяем работу функции NewConfig с разными значениями переменной окружения GO_HIVEHOOK_MAIN", func() {
		It("", func() {

		})
	})*/
})
