package simplelogger_test

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"

	"github.com/av-belyakov/simplelogger"
	"github.com/av-belyakov/thehivehook_go_package/cmd/elasticsearchapi"
	"github.com/av-belyakov/thehivehook_go_package/internal/confighandler"
)

var (
	rootDir      string = "thehivehook_go_package"
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
	conf, err = confighandler.NewConfig(rootDir)
	if err != nil {
		log.Fatalln(err)
	}

	//инициализация модуля взаимодействия с БД
	confDB := conf.GetApplicationWriteLogDB()
	esc, err = elasticsearchapi.NewElasticsearchConnect(elasticsearchapi.Settings{
		Port:    confDB.Port,
		Host:    confDB.Host,
		User:    confDB.User,
		Passwd:  confDB.Passwd,
		IndexDB: confDB.StorageNameDB,
	})
	if err != nil {
		log.Fatalln(err)
	}

	//инициализация логера
	simpleLogger, err = simplelogger.NewSimpleLogger(context.Background(), rootDir, getLoggerSettings(conf.GetListLogs()))
	if err != nil {
		log.Fatalln(err)
	}

	//инициализация в логере возможности взаимодействия с БД
	simpleLogger.SetDataBaseInteraction(esc)

	os.Exit(m.Run())
}

func TestSimpleLogger(t *testing.T) {
	ok := simpleLogger.Write("error", "test error message for check write to database")
	assert.True(t, ok)

	ok = simpleLogger.Write("info", "test info message")
	assert.True(t, ok)
}

func getLoggerSettings(cls []confighandler.LogSet) []simplelogger.Options {
	loggerConf := make([]simplelogger.Options, 0, len(cls))

	for _, v := range cls {
		loggerConf = append(loggerConf, simplelogger.Options{
			WritingToDB:     v.WritingDB,
			WritingToFile:   v.WritingFile,
			WritingToStdout: v.WritingStdout,
			MsgTypeName:     v.MsgTypeName,
			PathDirectory:   v.PathDirectory,
			MaxFileSize:     v.MaxFileSize,
		})
	}

	return loggerConf
}
