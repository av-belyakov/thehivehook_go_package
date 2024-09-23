package webhookserver_test

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"testing"

	"github.com/av-belyakov/thehivehook_go_package/cmd/webhookserver"
	"github.com/av-belyakov/thehivehook_go_package/internal/logginghandler"
)

func TestWebhookServer(t *testing.T) {
	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	go func() {
		sigChan := make(chan os.Signal, 1)
		osCall := <-sigChan
		log.Printf("system call:%+v", osCall)

		cancel()
	}()

	logging := logginghandler.New()

	go func() {
		for msg := range logging.GetChan() {
			fmt.Println("Logging:", msg)
		}
	}()

	webHookServer, errServer := webhookserver.New(ctx, "192.168.9.208", 5000, logging)
	if errServer != nil {
		t.Fatal("create new server %w", errServer)
	}

	webHookServer.Start()
}
