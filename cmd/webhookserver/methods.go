package webhookserver

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/av-belyakov/thehivehook_go_package/internal/logginghandler"
)

func New(ctx context.Context, opts WebHookServerOptions, logging *logginghandler.LoggingChan) (*WebHookServer, <-chan ChanFormWebHookServer, error) {
	chanOutput := make(chan ChanFormWebHookServer)

	wh := &WebHookServer{
		name:      opts.Name,
		version:   opts.Version,
		ctx:       ctx,
		logger:    logging,
		chanInput: chanOutput,
	}

	whts, err := NewWebHookTemporaryStorage(opts.TTL)
	if err != nil {
		return wh, chanOutput, err
	}
	wh.storage = whts

	if opts.Host == "" {
		return wh, chanOutput, errors.New("the value of 'host' cannot be empty")
	}
	wh.host = opts.Host

	if opts.Port == 0 || opts.Port > 65535 {
		return wh, chanOutput, errors.New("an incorrect network port value was received")
	}
	wh.port = opts.Port

	return wh, chanOutput, nil
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

	msg := fmt.Sprintf("server 'WebHookServer' was successfully launched, ip:%s, port:%d", wh.host, wh.port)
	log.Println(msg)
	wh.logger.Send("info", msg)

	<-wh.ctx.Done()
	close(wh.chanInput)
}

func (wh *WebHookServer) Shutdown(ctx context.Context) {
	wh.server.Shutdown(ctx)
}
