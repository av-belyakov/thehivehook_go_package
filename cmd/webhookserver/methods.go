package webhookserver

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/av-belyakov/thehivehook_go_package/internal/logginghandler"
)

func New(ctx context.Context, logging *logginghandler.LoggingChan, opts ...webHookServerOptions) (*WebHookServer, <-chan ChanFormWebHookServer, error) {
	chanOutput := make(chan ChanFormWebHookServer)

	whs := &WebHookServer{
		ctx:       ctx,
		name:      "gcm",
		version:   "0.1.1",
		host:      "127.0.0.1",
		port:      7575,
		ttl:       10,
		logger:    logging,
		chanInput: chanOutput,
	}

	for _, opt := range opts {
		opt(whs)
	}

	whts, err := NewWebHookTemporaryStorage(whs.ttl)
	if err != nil {
		return whs, chanOutput, err
	}
	whs.storage = whts

	return whs, chanOutput, nil
}

func (wh *WebHookServer) Start() {
	defer func() {
		wh.Shutdown(context.Background())
	}()
	routers := map[string]func(http.ResponseWriter, *http.Request){
		"/":        wh.RouteIndex,
		"/webhook": wh.RouteWebHook,
	}

	mux := http.NewServeMux()
	for k, v := range routers {
		mux.HandleFunc(k, v)
	}

	server := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", wh.host, wh.port),
		Handler: mux,
	}

	go func() {
		if errServer := server.ListenAndServe(); errServer != nil {
			log.Fatal(errServer)
		}
	}()

	msg := fmt.Sprintf("server 'WebHookServer' was successfully launched, %s:%d", wh.host, wh.port)
	log.Println(msg)
	wh.logger.Send("info", msg)

	<-wh.ctx.Done()
	close(wh.chanInput)
}

func (wh *WebHookServer) Shutdown(ctx context.Context) {
	wh.server.Shutdown(ctx)
}

func WithTTL(v int) webHookServerOptions {
	return func(whs *WebHookServer) {
		whs.ttl = v
	}
}

func WithPort(v int) webHookServerOptions {
	return func(whs *WebHookServer) {
		whs.port = v
	}
}

func WithHost(v string) webHookServerOptions {
	return func(whs *WebHookServer) {
		whs.host = v
	}
}

func WithName(v string) webHookServerOptions {
	return func(whs *WebHookServer) {
		whs.name = v
	}
}

func WithVersion(v string) webHookServerOptions {
	return func(whs *WebHookServer) {
		whs.version = v
	}
}
