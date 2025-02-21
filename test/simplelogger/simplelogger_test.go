package simplelogger_test

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"

	"github.com/av-belyakov/simplelogger"
	"github.com/av-belyakov/thehivehook_go_package/cmd/elasticsearchapi"
	"github.com/av-belyakov/thehivehook_go_package/internal/confighandler"
	"github.com/av-belyakov/thehivehook_go_package/internal/supportingfunctions"
)

const Root_Dir string = "thehivehook_go_package"

var (
	conf         *confighandler.ConfigApp
	simpleLogger *simplelogger.SimpleLoggerSettings
	esc          *elasticsearchapi.ElasticsearchDB

	err error
)

func TestMain(m *testing.M) {
	os.Unsetenv("GO_HIVEHOOK_DBWLOGPASSWD")

	//чтение переменных окружения
	if err := godotenv.Load("../../.env"); err != nil {
		log.Fatalln(err)
	}

	//чтение конфигурационных файлов
	conf, err = confighandler.NewConfig(Root_Dir)
	if err != nil {
		log.Fatalln(err)
	}

	//инициализация логгера
	var listLog []simplelogger.OptionsManager
	for k, v := range conf.GetListLogs() {
		fmt.Printf("%d. values:'%+v'\n", k, v)
		listLog = append(listLog, v)
	}

	opts := simplelogger.CreateOptions(listLog...)
	simpleLogger, err = simplelogger.NewSimpleLogger(context.Background(), Root_Dir, opts)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println("Version:", simplelogger.GetVersion())
	fmt.Println("___ OPTIONS:", opts)

	confDB := conf.GetApplicationWriteLogDB()
	if esc, err := elasticsearchapi.NewElasticsearchConnect(elasticsearchapi.Settings{
		Port:               confDB.Port,
		Host:               confDB.Host,
		User:               confDB.User,
		Passwd:             confDB.Passwd,
		IndexDB:            confDB.StorageNameDB,
		NameRegionalObject: conf.Name,
	}); err != nil {
		_ = simpleLogger.Write("error", supportingfunctions.CustomError(err).Error())
	} else {
		simpleLogger.SetDataBaseInteraction(esc)
	}

	os.Exit(m.Run())
}

func TestSimpleLogger(t *testing.T) {
	msg := supportingfunctions.ReplaceCommaCharacter("Post \"http://thehive.cloud.gcm:9000/api/v1/query?name=case-procedures\": context deadline exceeded, network connection error")

	ok := simpleLogger.Write("error", msg)
	assert.True(t, ok)

	ok = simpleLogger.Write("info", "test info message")
	assert.True(t, ok)
}
