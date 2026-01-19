package natssendcommands_test

import (
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/joho/godotenv"
	"github.com/nats-io/nats.go"
	"github.com/stretchr/testify/assert"

	"github.com/av-belyakov/thehivehook_go_package/cmd/constants"
	"github.com/av-belyakov/thehivehook_go_package/internal/confighandler"
	"github.com/av-belyakov/thehivehook_go_package/internal/supportingfunctions"
)

func TestCommandSendTags(t *testing.T) {
	if err := godotenv.Load("../../.env"); err != nil {
		log.Fatalln(err)
	}

	os.Setenv("GO_HIVEHOOK_MAIN", "prod")

	rootPath, err := supportingfunctions.GetRootPath(constants.Root_Dir)
	if err != nil {
		t.Fatalf("Не удалось получить корневую директорию: %v", err)
	}

	fmt.Println("path:", rootPath)

	conf, err := confighandler.NewConfig(rootPath)
	if err != nil {
		t.Fatalf("Не удалось прочитать конфигурационный файл: %v", err)
	}

	fmt.Printf(
		"Nats host:'%s', port:'%d'. Regional name:'%s'\n",
		conf.GetApplicationNATS().Host,
		conf.GetApplicationNATS().Port,
		conf.GetApplicationWebHookServer().Name,
	)

	nc, err := nats.Connect(fmt.Sprintf("%s:%d", conf.GetApplicationNATS().Host, conf.GetApplicationNATS().Port),
		nats.MaxReconnects(-1),
		nats.ReconnectWait(3*time.Second))
	if err != nil {
		log.Fatalln(err)
	}

	// обработка разрыва соединения с NATS
	nc.SetDisconnectErrHandler(func(c *nats.Conn, err error) {
		log.Println(err)
	})

	// обработка переподключения к NATS
	nc.SetReconnectHandler(func(c *nats.Conn) {
		log.Println(err)
	})

	t.Run("Тест 1. Отправляем команду 'add_case_tag'", func(t *testing.T) {
		err := nc.Publish(conf.GetApplicationNATS().Subscriptions.ListenerCommand,
			fmt.Appendf(
				nil,
				`{
					          "service": "MISP",
					          "command": "add_case_tag",
					  		  "for_regional_object": "%s",
					          "root_id": "%s",
					          "case_id": "%s",
					          "value": "Webhook: send=\"___ MISP ___ProductioN\""
					        }`,
				conf.GetApplicationWebHookServer().Name,
				"~88678416456",
				"39100",
			))
		assert.NoError(t, err)
		nc.Flush()

	})

	t.Cleanup(func() {
		nc.Drain()
	})
}
