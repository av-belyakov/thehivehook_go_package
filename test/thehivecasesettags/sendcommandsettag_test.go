package thehivecasesettags

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/joho/godotenv"

	"github.com/av-belyakov/thehivehook_go_package/cmd/constants"
	"github.com/av-belyakov/thehivehook_go_package/cmd/natsapi"
	"github.com/av-belyakov/thehivehook_go_package/cmd/thehiveapi"
	"github.com/av-belyakov/thehivehook_go_package/internal/confighandler"
	"github.com/av-belyakov/thehivehook_go_package/internal/supportingfunctions"
	helperfunc "github.com/av-belyakov/thehivehook_go_package/test/helpfunc"
)

func TestSendCommandSetTag(t *testing.T) {
	if err := godotenv.Load("../../.env"); err != nil {
		log.Fatalln(err)
	}

	os.Setenv("GO_HIVEHOOK_MAIN", "test")

	rootPath, err := supportingfunctions.GetRootPath(constants.Root_Dir)
	if err != nil {
		t.Fatalf("Не удалось получить корневую директорию: %v", err)
	}

	fmt.Println("path:", rootPath)

	conf, err := confighandler.NewConfig(rootPath)
	if err != nil {
		t.Fatalf("Не удалось прочитать конфигурационный файл: %v", err)
	}

	/*conf := confighandler.AppConfigTheHive{
		Port:   9000,
		Host:   "thehive.cloud.gcm",
		ApiKey: os.Getenv("GO_HIVEHOOK_THAPIKEY"),
	}*/

	ctx, cancel := context.WithCancel(context.Background())

	logging := helperfunc.NewLoggingForTest()
	go func() {
		for {
			select {
			case <-ctx.Done():
				return

			case msg := <-logging.GetChan():
				fmt.Printf("Log message: type:'%s', message:'%s'\n", msg.GetType(), msg.GetMessage())
			}
		}
	}()

	apiTheHive, err := thehiveapi.New(
		logging,
		thehiveapi.WithAPIKey(conf.ApiKey),
		thehiveapi.WithHost(conf.GetApplicationTheHive().Host),
		thehiveapi.WithPort(conf.GetApplicationTheHive().Port),
		thehiveapi.WithNameRegionalObject(conf.GetApplicationWebHookServer().Name),
	)
	if err != nil {
		log.Fatalln(err)
	}

	chApiTheHive, err := apiTheHive.Start(context.Background())
	if err != nil {
		log.Fatalln(err)
	}

	t.Run("Тест 1. Отправить команду на установку тега", func(t *testing.T) {
		req := natsapi.NewChannelRequest()
		req.SetCommand("send_command")
		req.SetOrder("add_case_tag")
		req.SetRequestId(uuid.New().String())
		req.SetData([]byte(`{
			  "service": "MISP",
  			  "command": "add_case_tag",
			  "for_regional_object": "gcm-test",
  			  "root_id": "~88678416456",
  			  "case_id": "39100",
  			  "value": "Webhook: send=\"WEBKOOK_Elasticsearch TEST new tag\""
			}`))

		chApiTheHive <- req

		time.Sleep(2 * time.Second)

		//msg := <-logging.GetChan()
		//fmt.Println("Type:", msg.GetType(), " LOG:", msg.GetMessage())

		//Expect(msg.GetType()).ShouldNot(Equal("error"))
	})

	t.Cleanup(func() {
		cancel()

		os.Unsetenv("GO_HIVEHOOK_THAPIKEY")
		os.Unsetenv("GO_HIVEHOOK_DBWLOGPASSWD")
	})
}
