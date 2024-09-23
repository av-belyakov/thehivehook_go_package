package webhookserver

import (
	"context"
	"net/http"

	"github.com/av-belyakov/thehivehook_go_package/internal/logginghandler"
)

type WebHookServer struct {
	port    int
	host    string
	version string
	ctx     context.Context
	server  *http.Server
	logger  *logginghandler.LoggingChan
}
