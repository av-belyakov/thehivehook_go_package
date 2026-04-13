package webhookserver

import (
	"context"
	"fmt"
	"log"
	"os/signal"
	"syscall"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"

	"github.com/av-belyakov/thehivehook_go_package/cmd/webhookserver"
	"github.com/av-belyakov/thehivehook_go_package/internal/confighandler"
	"github.com/av-belyakov/thehivehook_go_package/internal/constants"
	"github.com/av-belyakov/thehivehook_go_package/test/helpers"
)

func TestWebhookServer(t *testing.T) {
	//
	// ВАЖНО!!!
	//
	//перед запуском теста установите переменную окружения GO_HIVEHOOK_THAPIKEY
	//с ключем-идентификатором, необходимым для авторизации в API TheHive,
	//командой export GO_HIVEHOOK_THAPIKEY=<api_key> или воспользоватся godotenv
	if err := godotenv.Load("../../.env"); err != nil {
		log.Fatalln(err)
	}

	cfg, err := confighandler.NewConfig(constants.Root_Dir)
	if err != nil {
		log.Fatalln(err)
	}

	cfg.GetApplicationWebHookServer().Host = "localhost"

	logging := helpers.NewLoggingForTest()
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)

	go func() {
		for {
			select {
			case <-ctx.Done():
				cancel()

				return

			case msg := <-logging.GetChan():
				fmt.Printf("Log message: type:'%s', message:'%s'\n", msg.GetType(), msg.GetMessage())
			}
		}
	}()

	webHookServer, chanOutputWebHookServer, err := webhookserver.New(
		logging,
		webhookserver.WithVersion("0.0.1-for-test"),
		webhookserver.WithName(cfg.GetApplicationWebHookServer().Name),
		webhookserver.WithHost(cfg.GetApplicationWebHookServer().Host),
		webhookserver.WithPort(cfg.GetApplicationWebHookServer().Port),
	)
	if err != nil {
		log.Fatalln(err)
	}

	go func() {
		for {
			select {
			case <-ctx.Done():
				cancel()

				return

			case msg := <-chanOutputWebHookServer:
				fmt.Printf(
					"Received message from webhook server, object type:'%s', root id:'%s', data:'%+v'\n",
					msg.ObjectType,
					msg.RootId,
					msg.Data,
				)
			}
		}
	}()

	assert.NoError(t, webHookServer.Start(ctx))

	t.Cleanup(func() {
		cancel()
	})
}
