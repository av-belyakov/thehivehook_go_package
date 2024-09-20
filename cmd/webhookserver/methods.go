package webhookserver

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
)

func New(ctx context.Context, host string, port int) (*WebHookServer, error) {
	wh := &WebHookServer{version: "1.1.0"}

	if host == "" {
		return wh, errors.New("the value of 'host' cannot be empty")
	}

	if port == 0 || port > 65535 {
		return wh, errors.New("an incorrect network port value was received")
	}

	wh.ctx = ctx
	wh.host = host
	wh.port = port

	return wh, nil
}

func (wh *WebHookServer) Start() error {
	var err error
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
	wh.server = server

	go func() {
		if errServer := server.ListenAndServe(); errServer != nil {
			err = errServer
		}
	}()

	log.Printf("server 'WebHookServer' was successfully launched, ip:%s, port:%d", wh.host, wh.port)
	<-wh.ctx.Done()

	return err
}

func (wh *WebHookServer) Shutdown(ctx context.Context) {
	wh.server.Shutdown(ctx)
}
