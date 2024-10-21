// Основной модуль приложения
package webhookserver

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/av-belyakov/thehivehook_go_package/internal/logginghandler"
	"github.com/av-belyakov/thehivehook_go_package/internal/versionandname"
)

// New конструктор webhookserver принимает функциональные опции для настройки модуля перед запуском
func New(ctx context.Context, logging *logginghandler.LoggingChan, opts ...webHookServerOptions) (*WebHookServer, <-chan ChanFromWebHookServer, error) {
	chanOutput := make(chan ChanFromWebHookServer)

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

	whts, err := temporarystorage.NewWebHookTemporaryStorage(whs.ttl)
	if err != nil {
		return whs, chanOutput, err
	}
	whs.storage = whts

	return whs, chanOutput, nil
}

// Start выполняет запуск модуля
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

	msg := fmt.Sprintf("Application '%s' v%s was successfully launched, %s:%d", versionandname.GetName(), versionandname.GetVersion(), wh.host, wh.port)
	log.Println(msg)
	wh.logger.Send("info", msg)

	<-wh.ctx.Done()
	close(wh.chanInput)
}

// Shutdown завершает работу модуля
func (wh *WebHookServer) Shutdown(ctx context.Context) {
	wh.server.Shutdown(ctx)
}

//******************** функциональные настройки webhookserver ***********************

// WithTTL устанавливает время TimeToLive для временного хранилища информации в модуле
func WithTTL(v int) webHookServerOptions {
	return func(whs *WebHookServer) {
		whs.ttl = v
	}
}

// WithPort устанавливает порт для взаимодействия с модулем
func WithPort(v int) webHookServerOptions {
	return func(whs *WebHookServer) {
		whs.port = v
	}
}

// WithHost устанавливает хост для взаимодействия с модулем
func WithHost(v string) webHookServerOptions {
	return func(whs *WebHookServer) {
		whs.host = v
	}
}

// WithName устанавливает наименование модуля (обязательно). Наименование основывается
// на имени организации или подразделения эксплуатирующем модуль. Например, gcm, rcmslx и т.д.
func WithName(v string) webHookServerOptions {
	return func(whs *WebHookServer) {
		whs.name = v
	}
}

// WithVersion устанавливает версию модуля (опционально)
func WithVersion(v string) webHookServerOptions {
	return func(whs *WebHookServer) {
		whs.version = v
	}
}
