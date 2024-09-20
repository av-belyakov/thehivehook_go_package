package webhookserver_test

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"testing"

	"github.com/av-belyakov/thehivehook_go_package/cmd/webhookserver"
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

	webHookServer, errServer := webhookserver.New(ctx, "192.168.13.3", 5000)
	if errServer != nil {
		t.Fatal("create new server %w", errServer)
	}

	if err := webHookServer.Start(); err != nil {
		t.Fatal("fatal start web hook server: %w", err)
	}
}
