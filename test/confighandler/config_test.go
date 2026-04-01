package confighandler_test

import (
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/av-belyakov/thehivehook_go_package/internal/confighandler"
	"github.com/av-belyakov/thehivehook_go_package/internal/constants"
)

func TestConfigHandler(t *testing.T) {
	const (
		THEHIVE_APIKEY          = "ghu9dffhvbx2237ads2f3fsf"
		DATABASEWRITELOG_PASSWD = "cjis8w-dff0w0-fy2y3"
	)

	var (
		listTesting []TestOptions
		testOptions TestOptions

		cfg *confighandler.ConfigApp
		err error
	)

	unsetAllEnviromentEnvAny()

	//if err := godotenv.Load("../../.env"); err != nil {
	//	log.Fatalln(err)
	//}

	os.Setenv("GO_HIVEHOOK_THAPIKEY", THEHIVE_APIKEY)
	os.Setenv("GO_HIVEHOOK_DBWLOGPASSWD", DATABASEWRITELOG_PASSWD)

	// --- Общие настройки (из config.yml) ---
	testOptions = TestOptions{
		name: "Общие настройки получаемые из файла 'config.yaml'",
		function: func() {
			cfg, err = confighandler.NewConfig(constants.Root_Dir)
			testOptions.err = err
			testOptions.items = []TestParametrs{
				{
					inputParameters:    TestTypeElements{valueString: cfg.GetCommonApplication().Zabbix.NetworkHost},
					expectedParameters: TestTypeElements{valueString: "192.168.9.45"},
				},
				{
					inputParameters:    TestTypeElements{valueInt: cfg.GetCommonApplication().Zabbix.NetworkPort},
					expectedParameters: TestTypeElements{valueInt: 10051},
				},
				{
					inputParameters:    TestTypeElements{valueString: cfg.GetCommonApplication().Zabbix.ZabbixHost},
					expectedParameters: TestTypeElements{valueString: "test-uchet-db.cloud.gcm"},
				},
				// для отслеживания ошибок
				{
					inputParameters:    TestTypeElements{valueString: cfg.GetCommonApplication().Zabbix.EventTypes[0].EventType},
					expectedParameters: TestTypeElements{valueString: "error"},
				},
				{
					inputParameters:    TestTypeElements{valueString: cfg.GetCommonApplication().Zabbix.EventTypes[0].ZabbixKey},
					expectedParameters: TestTypeElements{valueString: "shaper_stix.error"},
				},
				{
					inputParameters:    TestTypeElements{valueBool: cfg.GetCommonApplication().Zabbix.EventTypes[0].IsTransmit},
					expectedParameters: TestTypeElements{valueBool: true},
				},
				{
					inputParameters:    TestTypeElements{valueInt: cfg.GetCommonApplication().Zabbix.EventTypes[0].Handshake.TimeInterval},
					expectedParameters: TestTypeElements{valueInt: 0},
				},
				{
					inputParameters:    TestTypeElements{valueString: cfg.GetCommonApplication().Zabbix.EventTypes[0].Handshake.Message},
					expectedParameters: TestTypeElements{valueString: ""},
				},
				// для информационных сообщений о выполненной работе
				{
					inputParameters:    TestTypeElements{valueString: cfg.GetCommonApplication().Zabbix.EventTypes[1].EventType},
					expectedParameters: TestTypeElements{valueString: "info"},
				},
				{
					inputParameters:    TestTypeElements{valueString: cfg.GetCommonApplication().Zabbix.EventTypes[1].ZabbixKey},
					expectedParameters: TestTypeElements{valueString: "shaper_stix.info"},
				},
				{
					inputParameters:    TestTypeElements{valueBool: cfg.GetCommonApplication().Zabbix.EventTypes[1].IsTransmit},
					expectedParameters: TestTypeElements{valueBool: true},
				},
				{
					inputParameters:    TestTypeElements{valueInt: cfg.GetCommonApplication().Zabbix.EventTypes[1].Handshake.TimeInterval},
					expectedParameters: TestTypeElements{valueInt: 0},
				},
				{
					inputParameters:    TestTypeElements{valueString: cfg.GetCommonApplication().Zabbix.EventTypes[1].Handshake.Message},
					expectedParameters: TestTypeElements{valueString: "I'm still alive"},
				},
				// для регулярного отстукивания что модуль еще работает
				{
					inputParameters:    TestTypeElements{valueString: cfg.GetCommonApplication().Zabbix.EventTypes[2].EventType},
					expectedParameters: TestTypeElements{valueString: "handshake"},
				},
				{
					inputParameters:    TestTypeElements{valueString: cfg.GetCommonApplication().Zabbix.EventTypes[2].ZabbixKey},
					expectedParameters: TestTypeElements{valueString: "shaper_stix.handshake"},
				},
				{
					inputParameters:    TestTypeElements{valueBool: cfg.GetCommonApplication().Zabbix.EventTypes[2].IsTransmit},
					expectedParameters: TestTypeElements{valueBool: true},
				},
				{
					inputParameters:    TestTypeElements{valueInt: cfg.GetCommonApplication().Zabbix.EventTypes[2].Handshake.TimeInterval},
					expectedParameters: TestTypeElements{valueInt: 1},
				},
				{
					inputParameters:    TestTypeElements{valueString: cfg.GetCommonApplication().Zabbix.EventTypes[2].Handshake.Message},
					expectedParameters: TestTypeElements{valueString: "0"},
				},
			}
		},
	}
	testOptions.function()
	listTesting = append(listTesting, testOptions)

	// --- Настройки NATS (для config_prod.yml) ---
	testOptions = TestOptions{
		name: "Настройки NATS (чтение конфигурационного файла по умолчанию config_prod.yaml)",
		function: func() {
			cfg, err = confighandler.NewConfig(constants.Root_Dir)

			testOptions.err = err
			testOptions.items = []TestParametrs{
				{
					inputParameters:    TestTypeElements{valueString: cfg.GetCommonInfo().FileName},
					expectedParameters: TestTypeElements{valueString: "config production"},
				},
				{
					inputParameters:    TestTypeElements{valueString: cfg.GetApplicationNATS().Host},
					expectedParameters: TestTypeElements{valueString: "nats.cloud.gcm"},
				},
				{
					inputParameters:    TestTypeElements{valueInt: cfg.GetApplicationNATS().Port},
					expectedParameters: TestTypeElements{valueInt: 4222},
				},
				{
					inputParameters:    TestTypeElements{valueInt: cfg.GetApplicationNATS().CacheTTL},
					expectedParameters: TestTypeElements{valueInt: 3600},
				},
				{
					inputParameters:    TestTypeElements{valueString: cfg.GetApplicationNATS().Subscriptions.SenderCase},
					expectedParameters: TestTypeElements{valueString: "object.casetype"},
				},
				{
					inputParameters:    TestTypeElements{valueString: cfg.GetApplicationNATS().Subscriptions.SenderAlert},
					expectedParameters: TestTypeElements{valueString: "object.alerttype"},
				},
				{
					inputParameters:    TestTypeElements{valueString: cfg.GetApplicationNATS().Subscriptions.ListenerCommand},
					expectedParameters: TestTypeElements{valueString: "object.commandstype"},
				},
			}
		},
	}
	testOptions.function()
	listTesting = append(listTesting, testOptions)

	// --- Настройки TheHive (для config_prod.yml) ---
	testOptions = TestOptions{
		name: "Настройки TheHive (чтение конфигурационного файла по умолчанию config_prod.yaml)",
		function: func() {
			cfg, err = confighandler.NewConfig(constants.Root_Dir)

			testOptions.err = err
			testOptions.items = []TestParametrs{
				{
					inputParameters:    TestTypeElements{valueString: cfg.GetCommonInfo().FileName},
					expectedParameters: TestTypeElements{valueString: "config production"},
				},
				{
					inputParameters:    TestTypeElements{valueString: cfg.GetApplicationTheHive().Host},
					expectedParameters: TestTypeElements{valueString: "thehive.cloud.gcm"},
				},
				{
					inputParameters:    TestTypeElements{valueInt: cfg.GetApplicationTheHive().Port},
					expectedParameters: TestTypeElements{valueInt: 9000},
				},
				{
					inputParameters:    TestTypeElements{valueInt: cfg.GetApplicationTheHive().CacheTTL},
					expectedParameters: TestTypeElements{valueInt: 43_200},
				},
				{
					inputParameters:    TestTypeElements{valueString: cfg.GetApplicationTheHive().ApiKey},
					expectedParameters: TestTypeElements{valueString: THEHIVE_APIKEY},
				},
			}
		},
	}
	testOptions.function()
	listTesting = append(listTesting, testOptions)

	// --- Настройки WEBHOOKSERVER (для config_prod.yml) ---
	testOptions = TestOptions{
		name: "Настройки WEBHOOKSERVER (чтение конфигурационного файла по умолчанию config_prod.yaml)",
		function: func() {
			cfg, err = confighandler.NewConfig(constants.Root_Dir)

			testOptions.err = err
			testOptions.items = []TestParametrs{
				{
					inputParameters:    TestTypeElements{valueString: cfg.GetCommonInfo().FileName},
					expectedParameters: TestTypeElements{valueString: "config production"},
				},
				{
					inputParameters:    TestTypeElements{valueString: cfg.GetApplicationWebHookServer().Name},
					expectedParameters: TestTypeElements{valueString: "gcm"},
				},
				{
					inputParameters:    TestTypeElements{valueString: cfg.GetApplicationWebHookServer().Host},
					expectedParameters: TestTypeElements{valueString: "192.168.9.53"},
				},
				{
					inputParameters:    TestTypeElements{valueInt: cfg.GetApplicationWebHookServer().Port},
					expectedParameters: TestTypeElements{valueInt: 5001},
				},
				{
					inputParameters:    TestTypeElements{valueInt: cfg.GetApplicationWebHookServer().StorageTTL},
					expectedParameters: TestTypeElements{valueInt: 180},
				},
				{
					inputParameters:    TestTypeElements{valueInt: cfg.GetApplicationWebHookServer().StorageDelayToSend},
					expectedParameters: TestTypeElements{valueInt: 30},
				},
			}
		},
	}
	testOptions.function()
	listTesting = append(listTesting, testOptions)

	// --- Настройки WRITELOGDB (для config_prod.yml) ---
	testOptions = TestOptions{
		name: "Настройки WRITELOGDB (чтение конфигурационного файла по умолчанию config_prod.yaml)",
		function: func() {
			cfg, err = confighandler.NewConfig(constants.Root_Dir)

			testOptions.err = err
			testOptions.items = []TestParametrs{
				{
					inputParameters:    TestTypeElements{valueString: cfg.GetCommonInfo().FileName},
					expectedParameters: TestTypeElements{valueString: "config production"},
				},
				{
					inputParameters:    TestTypeElements{valueString: cfg.GetApplicationWriteLogDB().Host},
					expectedParameters: TestTypeElements{valueString: "datahook.cloud.gcm"},
				},
				{
					inputParameters:    TestTypeElements{valueInt: cfg.GetApplicationWriteLogDB().Port},
					expectedParameters: TestTypeElements{valueInt: 9200},
				},
				{
					inputParameters:    TestTypeElements{valueString: cfg.GetApplicationWriteLogDB().NameDB},
					expectedParameters: TestTypeElements{valueString: ""},
				},
				{
					inputParameters:    TestTypeElements{valueString: cfg.GetApplicationWriteLogDB().StorageNameDB},
					expectedParameters: TestTypeElements{valueString: "thehivehook_go_package"},
				},
				{
					inputParameters:    TestTypeElements{valueString: cfg.GetApplicationWriteLogDB().User},
					expectedParameters: TestTypeElements{valueString: "log_writer"},
				},
				{
					inputParameters:    TestTypeElements{valueString: cfg.GetApplicationWriteLogDB().Passwd},
					expectedParameters: TestTypeElements{valueString: DATABASEWRITELOG_PASSWD},
				},
			}
		},
	}
	testOptions.function()
	listTesting = append(listTesting, testOptions)

	// --- Настройки NATS (для config_dev.yml) ---
	testOptions = TestOptions{
		name: "Настройки NATS (чтение конфигурационного файла по умолчанию config_dev.yaml)",
		function: func() {
			os.Setenv("GO_HIVEHOOK_MAIN", "development")

			cfg, err = confighandler.NewConfig(constants.Root_Dir)
			testOptions.err = err
			testOptions.items = []TestParametrs{
				{
					inputParameters:    TestTypeElements{valueString: cfg.GetCommonInfo().FileName},
					expectedParameters: TestTypeElements{valueString: "config_development"},
				},
				{
					inputParameters:    TestTypeElements{valueString: cfg.GetApplicationNATS().Host},
					expectedParameters: TestTypeElements{valueString: "192.168.9.208"},
				},
				{
					inputParameters:    TestTypeElements{valueInt: cfg.GetApplicationNATS().Port},
					expectedParameters: TestTypeElements{valueInt: 4222},
				},
				{
					inputParameters:    TestTypeElements{valueInt: cfg.GetApplicationNATS().CacheTTL},
					expectedParameters: TestTypeElements{valueInt: 3600},
				},
				{
					inputParameters:    TestTypeElements{valueString: cfg.GetApplicationNATS().Subscriptions.SenderCase},
					expectedParameters: TestTypeElements{valueString: "object.casetype.test"},
				},
				{
					inputParameters:    TestTypeElements{valueString: cfg.GetApplicationNATS().Subscriptions.SenderAlert},
					expectedParameters: TestTypeElements{valueString: "object.alerttype.test"},
				},
				{
					inputParameters:    TestTypeElements{valueString: cfg.GetApplicationNATS().Subscriptions.ListenerCommand},
					expectedParameters: TestTypeElements{valueString: "object.commandstype.test"},
				},
			}
		},
	}
	testOptions.function()
	listTesting = append(listTesting, testOptions)

	// --- Настройки TheHive (для config_dev.yml) ---
	testOptions = TestOptions{
		name: "Настройки TheHive (чтение конфигурационного файла по умолчанию config_dev.yaml)",
		function: func() {
			os.Setenv("GO_HIVEHOOK_MAIN", "development")

			cfg, err = confighandler.NewConfig(constants.Root_Dir)
			testOptions.err = err
			testOptions.items = []TestParametrs{
				{
					inputParameters:    TestTypeElements{valueString: cfg.GetCommonInfo().FileName},
					expectedParameters: TestTypeElements{valueString: "config_development"},
				},
				{
					inputParameters:    TestTypeElements{valueString: cfg.GetApplicationTheHive().Host},
					expectedParameters: TestTypeElements{valueString: "thehive.cloud.gcm"},
				},
				{
					inputParameters:    TestTypeElements{valueInt: cfg.GetApplicationTheHive().Port},
					expectedParameters: TestTypeElements{valueInt: 9000},
				},
				{
					inputParameters:    TestTypeElements{valueInt: cfg.GetApplicationTheHive().CacheTTL},
					expectedParameters: TestTypeElements{valueInt: 43_200},
				},
				{
					inputParameters:    TestTypeElements{valueString: cfg.GetApplicationTheHive().ApiKey},
					expectedParameters: TestTypeElements{valueString: THEHIVE_APIKEY},
				},
			}
		},
	}
	testOptions.function()
	listTesting = append(listTesting, testOptions)

	// --- Настройки WEBHOOKSERVER (для config_dev.yml) ---
	testOptions = TestOptions{
		name: "Настройки WEBHOOKSERVER (чтение конфигурационного файла по умолчанию config_dev.yaml)",
		function: func() {
			os.Setenv("GO_HIVEHOOK_MAIN", "development")

			cfg, err = confighandler.NewConfig(constants.Root_Dir)
			testOptions.err = err
			testOptions.items = []TestParametrs{
				{
					inputParameters:    TestTypeElements{valueString: cfg.GetCommonInfo().FileName},
					expectedParameters: TestTypeElements{valueString: "config_development"},
				},
				{
					inputParameters:    TestTypeElements{valueString: cfg.GetApplicationWebHookServer().Name},
					expectedParameters: TestTypeElements{valueString: "gcm-dev-test"},
				},
				{
					inputParameters:    TestTypeElements{valueString: cfg.GetApplicationWebHookServer().Host},
					expectedParameters: TestTypeElements{valueString: "192.168.9.208"},
				},
				{
					inputParameters:    TestTypeElements{valueInt: cfg.GetApplicationWebHookServer().Port},
					expectedParameters: TestTypeElements{valueInt: 5000},
				},
				{
					inputParameters:    TestTypeElements{valueInt: cfg.GetApplicationWebHookServer().StorageTTL},
					expectedParameters: TestTypeElements{valueInt: 180},
				},
				{
					inputParameters:    TestTypeElements{valueInt: cfg.GetApplicationWebHookServer().StorageDelayToSend},
					expectedParameters: TestTypeElements{valueInt: 30},
				},
			}
		},
	}
	testOptions.function()
	listTesting = append(listTesting, testOptions)

	// --- Настройки WRITELOGDB (для config_dev.yml) ---
	testOptions = TestOptions{
		name: "Настройки WRITELOGDB (чтение конфигурационного файла по умолчанию config_dev.yaml)",
		function: func() {
			os.Setenv("GO_HIVEHOOK_MAIN", "development")

			cfg, err = confighandler.NewConfig(constants.Root_Dir)
			testOptions.err = err
			testOptions.items = []TestParametrs{
				{
					inputParameters:    TestTypeElements{valueString: cfg.GetCommonInfo().FileName},
					expectedParameters: TestTypeElements{valueString: "config_development"},
				},
				{
					inputParameters:    TestTypeElements{valueString: cfg.GetApplicationWriteLogDB().Host},
					expectedParameters: TestTypeElements{valueString: "datahook.cloud.gcm"},
				},
				{
					inputParameters:    TestTypeElements{valueInt: cfg.GetApplicationWriteLogDB().Port},
					expectedParameters: TestTypeElements{valueInt: 9200},
				},
				{
					inputParameters:    TestTypeElements{valueString: cfg.GetApplicationWriteLogDB().NameDB},
					expectedParameters: TestTypeElements{valueString: ""},
				},
				{
					inputParameters:    TestTypeElements{valueString: cfg.GetApplicationWriteLogDB().StorageNameDB},
					expectedParameters: TestTypeElements{valueString: "thehivehook_go_package"},
				},
				{
					inputParameters:    TestTypeElements{valueString: cfg.GetApplicationWriteLogDB().User},
					expectedParameters: TestTypeElements{valueString: "log_writer"},
				},
				{
					inputParameters:    TestTypeElements{valueString: cfg.GetApplicationWriteLogDB().Passwd},
					expectedParameters: TestTypeElements{valueString: DATABASEWRITELOG_PASSWD},
				},
			}
		},
	}
	testOptions.function()
	listTesting = append(listTesting, testOptions)

	// --- Настройки NATS (через переменные окружения) ---
	testOptions = TestOptions{
		name: "Настройки NATS (через переменные окружения)",
		function: func() {
			const (
				NATS_HOST               = "nats.cloud.gcm.test.test"
				NATS_PORT               = 4545
				NATS_CACHETTL           = 3600
				NATS_SUBSENDERCASE      = "sender.case"
				NATS_SUBSENDERALERT     = "sender.alert"
				NATS_SUBLISTENERCOMMAND = "listener.command"
			)

			os.Setenv("GO_HIVEHOOK_NHOST", NATS_HOST)
			os.Setenv("GO_HIVEHOOK_NPORT", strconv.Itoa(NATS_PORT))
			os.Setenv("GO_HIVEHOOK_NCACHETTL", strconv.Itoa(NATS_CACHETTL))
			os.Setenv("GO_HIVEHOOK_NSUBSENDERCASE", NATS_SUBSENDERCASE)
			os.Setenv("GO_HIVEHOOK_NSUBSENDERALERT", NATS_SUBSENDERALERT)
			os.Setenv("GO_HIVEHOOK_NSUBLISTENERCOMMAND", NATS_SUBLISTENERCOMMAND)

			cfg, err = confighandler.NewConfig(constants.Root_Dir)
			testOptions.err = err
			testOptions.items = []TestParametrs{
				{
					inputParameters:    TestTypeElements{valueString: cfg.GetApplicationNATS().Host},
					expectedParameters: TestTypeElements{valueString: NATS_HOST},
				},
				{
					inputParameters:    TestTypeElements{valueInt: cfg.GetApplicationNATS().Port},
					expectedParameters: TestTypeElements{valueInt: NATS_PORT},
				},
				{
					inputParameters:    TestTypeElements{valueInt: cfg.GetApplicationNATS().CacheTTL},
					expectedParameters: TestTypeElements{valueInt: NATS_CACHETTL},
				},
				{
					inputParameters:    TestTypeElements{valueString: cfg.GetApplicationNATS().Subscriptions.SenderCase},
					expectedParameters: TestTypeElements{valueString: NATS_SUBSENDERCASE},
				},
				{
					inputParameters:    TestTypeElements{valueString: cfg.GetApplicationNATS().Subscriptions.SenderAlert},
					expectedParameters: TestTypeElements{valueString: NATS_SUBSENDERALERT},
				},
				{
					inputParameters:    TestTypeElements{valueString: cfg.GetApplicationNATS().Subscriptions.ListenerCommand},
					expectedParameters: TestTypeElements{valueString: NATS_SUBLISTENERCOMMAND},
				},
			}
		},
	}
	testOptions.function()
	listTesting = append(listTesting, testOptions)

	// --- Настройки TheHive (через переменные окружения) ---
	testOptions = TestOptions{
		name: "Настройки TheHive (через переменные окружения)",
		function: func() {
			const (
				THEHIVE_HOST     = "thehive.cloud.gcm.test"
				THEHIVE_PORT     = 1122
				THEHIVE_CACHETTL = 43200
			)

			os.Setenv("GO_HIVEHOOK_THHOST", THEHIVE_HOST)
			os.Setenv("GO_HIVEHOOK_THPORT", strconv.Itoa(THEHIVE_PORT))
			os.Setenv("GO_HIVEHOOK_THCACHETTL", strconv.Itoa(THEHIVE_CACHETTL))

			cfg, err = confighandler.NewConfig(constants.Root_Dir)
			testOptions.err = err
			testOptions.items = []TestParametrs{
				{
					inputParameters:    TestTypeElements{valueString: cfg.GetApplicationTheHive().Host},
					expectedParameters: TestTypeElements{valueString: THEHIVE_HOST},
				},
				{
					inputParameters:    TestTypeElements{valueInt: cfg.GetApplicationTheHive().Port},
					expectedParameters: TestTypeElements{valueInt: THEHIVE_PORT},
				},
				{
					inputParameters:    TestTypeElements{valueInt: cfg.GetApplicationTheHive().CacheTTL},
					expectedParameters: TestTypeElements{valueInt: THEHIVE_CACHETTL},
				},
				{
					inputParameters:    TestTypeElements{valueString: cfg.GetApplicationTheHive().ApiKey},
					expectedParameters: TestTypeElements{valueString: THEHIVE_APIKEY},
				},
			}
		},
	}
	testOptions.function()
	listTesting = append(listTesting, testOptions)

	// --- Настройки WEBHOOKSERVER (через переменные окружения) ---
	testOptions = TestOptions{
		name: "Настройки WEBHOOKSERVER (через переменные окружения)",
		function: func() {
			const (
				HIVEHOOK_WEBHHOST = "11.0.11.10"
				HIVEHOOK_WEBHPORT = 7822
				HIVEHOOK_WEBTTL   = 13
				HIVEHOOK_WEBDS    = 55
				HIVEHOOK_WEBHNAME = "gcm-rcm"
			)

			os.Setenv("GO_HIVEHOOK_WEBHNAME", HIVEHOOK_WEBHNAME)
			os.Setenv("GO_HIVEHOOK_WEBHHOST", HIVEHOOK_WEBHHOST)
			os.Setenv("GO_HIVEHOOK_WEBHPORT", strconv.Itoa(HIVEHOOK_WEBHPORT))
			os.Setenv("GO_HIVEHOOK_WEBHSTORAGETTL", strconv.Itoa(HIVEHOOK_WEBTTL))
			os.Setenv("GO_HIVEHOOK_WEBHSTORAGDS", strconv.Itoa(HIVEHOOK_WEBDS))

			cfg, err = confighandler.NewConfig(constants.Root_Dir)
			testOptions.err = err
			testOptions.items = []TestParametrs{
				{
					inputParameters:    TestTypeElements{valueString: cfg.GetApplicationWebHookServer().Name},
					expectedParameters: TestTypeElements{valueString: HIVEHOOK_WEBHNAME},
				},
				{
					inputParameters:    TestTypeElements{valueString: cfg.GetApplicationWebHookServer().Host},
					expectedParameters: TestTypeElements{valueString: HIVEHOOK_WEBHHOST},
				},
				{
					inputParameters:    TestTypeElements{valueInt: cfg.GetApplicationWebHookServer().Port},
					expectedParameters: TestTypeElements{valueInt: HIVEHOOK_WEBHPORT},
				},
				{
					inputParameters:    TestTypeElements{valueInt: cfg.GetApplicationWebHookServer().StorageTTL},
					expectedParameters: TestTypeElements{valueInt: HIVEHOOK_WEBTTL},
				},
				{
					inputParameters:    TestTypeElements{valueInt: cfg.GetApplicationWebHookServer().StorageDelayToSend},
					expectedParameters: TestTypeElements{valueInt: HIVEHOOK_WEBDS},
				},
			}
		},
	}
	testOptions.function()
	listTesting = append(listTesting, testOptions)

	// --- Настройки WRITELOGDB (через переменные окружения) ---
	testOptions = TestOptions{
		name: "Настройки WRITELOGDB (через переменные окружения)",
		function: func() {
			const (
				HIVEHOOK_DBWLOGHOST        = "45.10.32.1"
				HIVEHOOK_DBWLOGPORT        = 11123
				HIVEHOOK_DBWLOGNAME        = "log_db"
				HIVEHOOK_DBWLOGUSER        = "nreuser"
				HIVEHOOK_DBWLOGSTORAGENAME = "thehivehookgolog"
			)

			os.Setenv("GO_HIVEHOOK_DBWLOGHOST", HIVEHOOK_DBWLOGHOST)
			os.Setenv("GO_HIVEHOOK_DBWLOGPORT", strconv.Itoa(HIVEHOOK_DBWLOGPORT))
			os.Setenv("GO_HIVEHOOK_DBWLOGNAME", HIVEHOOK_DBWLOGNAME)
			os.Setenv("GO_HIVEHOOK_DBWLOGUSER", HIVEHOOK_DBWLOGUSER)
			os.Setenv("GO_HIVEHOOK_DBWLOGSTORAGENAME", HIVEHOOK_DBWLOGSTORAGENAME)

			cfg, err = confighandler.NewConfig(constants.Root_Dir)
			testOptions.err = err
			testOptions.items = []TestParametrs{
				{
					inputParameters:    TestTypeElements{valueString: cfg.GetApplicationWriteLogDB().Host},
					expectedParameters: TestTypeElements{valueString: HIVEHOOK_DBWLOGHOST},
				},
				{
					inputParameters:    TestTypeElements{valueInt: cfg.GetApplicationWriteLogDB().Port},
					expectedParameters: TestTypeElements{valueInt: HIVEHOOK_DBWLOGPORT},
				},
				{
					inputParameters:    TestTypeElements{valueString: cfg.GetApplicationWriteLogDB().NameDB},
					expectedParameters: TestTypeElements{valueString: HIVEHOOK_DBWLOGNAME},
				},
				{
					inputParameters:    TestTypeElements{valueString: cfg.GetApplicationWriteLogDB().StorageNameDB},
					expectedParameters: TestTypeElements{valueString: HIVEHOOK_DBWLOGSTORAGENAME},
				},
				{
					inputParameters:    TestTypeElements{valueString: cfg.GetApplicationWriteLogDB().User},
					expectedParameters: TestTypeElements{valueString: HIVEHOOK_DBWLOGUSER},
				},
				{
					inputParameters:    TestTypeElements{valueString: cfg.GetApplicationWriteLogDB().Passwd},
					expectedParameters: TestTypeElements{valueString: DATABASEWRITELOG_PASSWD},
				},
			}
		},
	}
	testOptions.function()
	listTesting = append(listTesting, testOptions)

	//--------------------------
	// --- Выполнение тестов ---
	for testNum, tt := range listTesting {
		t.Run(fmt.Sprintf("Тест %d. %s", testNum+1, tt.name), func(t *testing.T) {
			assert.NoError(t, tt.err)

			for _, v := range tt.items {
				/*fmt.Printf(
				`%d.
				expectedParameters.valueInt:'%d' = '%d':valueInt.inputParameters
				expectedParameters.valueString:'%s' = '%s':valueString.inputParameters`,
				k+1,
				v.expectedParameters.valueInt, v.inputParameters.valueInt,
				v.expectedParameters.valueString, v.inputParameters.valueString)*/

				assert.Equal(t, v.expectedParameters.valueInt, v.inputParameters.valueInt)
				assert.Equal(t, v.expectedParameters.valueString, v.inputParameters.valueString)
			}
		})
	}

	/*t.Run("", func(t *testing.T) {

	})*/

	t.Cleanup(func() {
		unsetAllEnviromentEnvAny()
	})
}

func unsetAllEnviromentEnvAny() {
	os.Unsetenv("GO_HIVEHOOK_MAIN")

	//авторизационные данные
	os.Unsetenv("GO_HIVEHOOK_THAPIKEY")
	os.Unsetenv("GO_HIVEHOOK_DBWLOGPASSWD")

	//настройки NATS
	os.Unsetenv("GO_HIVEHOOK_NHOST")
	os.Unsetenv("GO_HIVEHOOK_NPORT")
	os.Unsetenv("GO_HIVEHOOK_NCACHETTL")
	os.Unsetenv("GO_HIVEHOOK_NSUBSENDERCASE")
	os.Unsetenv("GO_HIVEHOOK_NSUBSENDERALERT")
	os.Unsetenv("GO_HIVEHOOK_NSUBLISTENERCOMMAND")

	//настройки TheHive
	os.Unsetenv("GO_HIVEHOOK_THHOST")
	os.Unsetenv("GO_HIVEHOOK_THPORT")
	os.Unsetenv("GO_HIVEHOOK_THCACHETTL")

	//настройки WebHook сервера
	os.Unsetenv("GO_HIVEHOOK_WEBHNAME")
	os.Unsetenv("GO_HIVEHOOK_WEBHHOST")
	os.Unsetenv("GO_HIVEHOOK_WEBHPORT")
	os.Unsetenv("GO_HIVEHOOK_WEBHSTORAGDS")
	os.Unsetenv("GO_HIVEHOOK_WEBHSTORAGETTL")

	//настройки доступа к БД в которую будут записыватся логи
	os.Unsetenv("GO_HIVEHOOK_DBWLOGHOST")
	os.Unsetenv("GO_HIVEHOOK_DBWLOGPORT")
	os.Unsetenv("GO_HIVEHOOK_DBWLOGNAME")
	os.Unsetenv("GO_HIVEHOOK_DBWLOGUSER")
	os.Unsetenv("GO_HIVEHOOK_DBWLOGSTORAGENAME")
}

type TestTypeElements struct {
	valueString string
	valueInt    int
	valueBool   bool
}

type TestParametrs struct {
	inputParameters    TestTypeElements
	expectedParameters TestTypeElements
}

type TestOptions struct {
	items    []TestParametrs
	function func()
	err      error
	name     string
}
