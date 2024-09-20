package webhookserver

import (
	"context"
	"net/http"
)

type WebHookServer struct {
	port    int
	host    string
	version string
	ctx     context.Context
	server  *http.Server
}
