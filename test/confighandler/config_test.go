package confighandler_test

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"

	"github.com/av-belyakov/thehivehook_go_package/internal/confighandler"
	"github.com/av-belyakov/thehivehook_go_package/internal/constants"
)

func TestConfigHandler(t *testing.T) {
	var (
		listTesting []TestOptions
		testOptions TestOptions

		cfg *confighandler.ConfigApp
		err error
	)

	unsetAllEnviromentEnvAny()

	if err := godotenv.Load("../../.env"); err != nil {
		log.Fatalln(err)
	}

	testOptions = TestOptions{
		name: "Чтение конфигурационного файла (по умолчанию config_prod.yaml)",
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
	/*
		написать дальше остальные тесты
	*/

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

	/*
		t.Run("Тест 1.3. Все пораметры конфигурационного файла 'config_prod.yaml' для THEHIVE должны быть успешно получены", func(t *testing.T) {
			cfgTheHive := cfg.GetApplicationTheHive()

			assert.Equal(t, cfgTheHive.Host, "thehive.cloud.gcm")
			assert.Equal(t, cfgTheHive.Port, 9000)
			assert.Equal(t, cfgTheHive.CacheTTL, 43200)
			assert.Equal(t, cfgTheHive.ApiKey, 0)
		})

		t.Run("Тест 1.4. Все пораметры конфигрурационного файла 'config_prod.yaml' для WEBHOOKSERVER должны быть успешно получены", func(t *testing.T) {
			cfgWebHS := cfg.GetApplicationWebHookServer()

			assert.Equal(t, cfgWebHS.Host, "192.168.13.3")
			assert.Equal(t, cfgWebHS.Port, 5000)
			assert.Equal(t, cfgWebHS.StorageTTL, 180)
			assert.Equal(t, cfgWebHS.StorageDelayToSend, 30)
			assert.Equal(t, cfgWebHS.Name, "gcm")
		})

		t.Run("Тест 1.5. Все пораметры конфигурационного файла 'config_prod.yaml' для DATABASEWRITELOG должны быть успешно получены", func(t *testing.T) {
			cfgWriteLog := cfg.GetApplicationWriteLogDB()

			assert.Equal(t, cfgWriteLog.Host, "datahook.cloud.gcm")
			assert.Equal(t, cfgWriteLog.Port, 9200)
			assert.Equal(t, cfgWriteLog.NameDB, "")
			assert.Equal(t, cfgWriteLog.StorageNameDB, "thehivehook_go_package")
			assert.Equal(t, cfgWriteLog.User, "log_writer")
			assert.Equal(t, len(cfgWriteLog.Passwd), 0)
		})

		t.Run("Тест 2. Чтение конфигурационного файла config_dev.yaml", func(t *testing.T) {
			os.Setenv("GO_HIVEHOOK_MAIN", "development")

			cfg, err = confighandler.NewConfig(constants.Root_Dir)
			if err != nil {
				log.Fatalln(err)
			}

			cfgCommonInfo := cfg.GetCommonInfo()
			assert.Equal(t, cfgCommonInfo.FileName, "config_dev")
		})

		t.Run("Тест 2.1. Все пораметры конфигурационного файла 'config_dev.yaml' для NATS должны быть успешно получены", func(t *testing.T) {
			cfgNATS := cfg.GetApplicationNATS()

			fmt.Println("Application NATS config:")
			fmt.Println(cfgNATS)

			// cfgNATS
			assert.Equal(t, cfgNATS.Host, "nats.cloud.gcm")
			assert.Equal(t, cfgNATS.Port, 4222)
			assert.Equal(t, cfgNATS.CacheTTL, 3600)
			assert.Equal(t, cfgNATS.Subscriptions.SenderCase, "object.casetype")
			assert.Equal(t, cfgNATS.Subscriptions.SenderAlert, "object.alerttype")
			assert.Equal(t, cfgNATS.Subscriptions.ListenerCommand, "object.commandstype")

			Expect(cn.Host).Should(Equal("nats.cloud.gcm"))
			Expect(cn.Port).Should(Equal(4222))
			Expect(cn.CacheTTL).Should(Equal(3600))
			Expect(cn.Subscriptions.SenderCase).Should(Equal("object.casetype"))
			Expect(cn.Subscriptions.SenderAlert).Should(Equal("object.alerttype"))
			Expect(cn.Subscriptions.ListenerCommand).Should(Equal("object.commandstype"))
		})

		t.Run("Тест 2.2. Все пораметры конфигрурационного файла 'config_dev.yaml' для THEHIVE должны быть успешно получены", func(t *testing.T) {
			cth := conf.GetApplicationTheHive()
			Expect(cth.Host).Should(Equal("thehive.cloud.gcm"))
			Expect(cth.Port).Should(Equal(9001))
			Expect(cth.CacheTTL).Should(Equal(3600))
			Expect(len(cth.ApiKey)).ShouldNot(Equal(0))
		})

		t.Run("Тест 2.3. Все пораметры конфигрурационного файла 'config_dev.yaml' для WEBHOOKSERVER должны быть успешно получены", func(t *testing.T) {
			chs := conf.GetApplicationWebHookServer()
			Expect(chs.Host).Should(Equal("127.0.0.1"))
			Expect(chs.Port).Should(Equal(5000))
			Expect(chs.StorageTTL).Should(Equal(180))
			Expect(chs.StorageDelayToSend).Should(Equal(30))
			Expect(chs.Name).Should(Equal("rcmsml"))
		})

		t.Run("Тест 2.4. Все пораметры конфигурационного файла 'config_dev.yaml' для DATABASEWRITELOG должны быть успешно получены", func(t *testing.T) {
			cwl := conf.GetApplicationWriteLogDB()
			Expect(cwl.Host).Should(Equal("datahook.cloud.gcm"))
			Expect(cwl.Port).Should(Equal(9200))
			Expect(cwl.NameDB).Should(Equal("nameDB"))
			Expect(cwl.StorageNameDB).Should(Equal("thehivehook_go_package"))
			Expect(cwl.User).Should(Equal("log_writer"))
			Expect(len(cwl.Passwd)).ShouldNot(Equal(0))
		})*/

	/*t.Run("", func(t *testing.T) {

	})*/

	t.Cleanup(func() {
		unsetAllEnviromentEnvAny()
	})
}

func unsetAllEnviromentEnvAny() {
	os.Unsetenv("GO_HIVEHOOK_MAIN")

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
