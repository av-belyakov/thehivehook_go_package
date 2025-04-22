// Основной модуль приложения
package webhookserver

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"net/http/pprof"
	"os"
	"strings"
	"time"

	"github.com/av-belyakov/thehivehook_go_package/cmd/commoninterfaces"
	"golang.org/x/sync/errgroup"
)

// New конструктор webhookserver принимает функциональные опции для настройки модуля перед запуском
func New(logger commoninterfaces.Logger, opts ...webHookServerOptions) (*WebHookServer, <-chan ChanFromWebHookServer, error) {
	chanOutput := make(chan ChanFromWebHookServer)
	whs := &WebHookServer{
		name:      "nobody",
		version:   "0.1.1",
		timeStart: time.Now(),
		host:      "127.0.0.1",
		port:      7575,
		ttl:       10,
		logger:    logger,
		chanInput: chanOutput,
	}

	for _, opt := range opts {
		if err := opt(whs); err != nil {
			return whs, chanOutput, err
		}
	}

	return whs, chanOutput, nil
}

// Start выполняет запуск модуля
func (wh *WebHookServer) Start(ctx context.Context) error {
	defer func() {
		close(wh.chanInput)
	}()

	routers := map[string]func(http.ResponseWriter, *http.Request){
		"/":        wh.RouteIndex,
		"/webhook": wh.RouteWebHook,
	}

	//для отладки через pprof (только для теста)
	//http://confWebHook.Host:confWebHook.Port/debug/pprof/
	//go tool pprof http://confWebHook.Host:confWebHook.Port/debug/pprof/heap
	//go tool pprof http://confWebHook.Host:confWebHook.Port/debug/pprof/allocs
	//go tool pprof http://confWebHook.Host:confWebHook.Port/debug/pprof/goroutine
	if os.Getenv("GO_HIVEHOOK_MAIN") == "test" {
		routers["/debug/pprof/"] = pprof.Index
	}

	mux := http.NewServeMux()
	for k, v := range routers {
		mux.HandleFunc(k, v)
	}

	wh.server = &http.Server{
		Addr:    fmt.Sprintf("%s:%d", wh.host, wh.port),
		Handler: mux,
		BaseContext: func(_ net.Listener) context.Context {
			return ctx
		},
	}

	g, gCtx := errgroup.WithContext(ctx)
	g.Go(func() error {
		// вывод информационного сообщения при старте приложения
		infoMsg := getInformationMessage(wh.name, wh.host, wh.port)
		wh.logger.Send("info", strings.ToLower(infoMsg))

		return wh.server.ListenAndServe()
	})
	g.Go(func() error {
		<-gCtx.Done()

		return wh.server.Shutdown(context.Background())
	})

	return g.Wait()
}

//******************** функциональные настройки webhookserver ***********************

// WithTTL устанавливает время TimeToLive для временного хранилища информации в модуле
func WithTTL(v int) webHookServerOptions {
	return func(whs *WebHookServer) error {
		whs.ttl = v

		return nil
	}
}

// WithPort устанавливает порт для взаимодействия с модулем
func WithPort(v int) webHookServerOptions {
	return func(whs *WebHookServer) error {
		whs.port = v

		return nil
	}
}

// WithHost устанавливает хост для взаимодействия с модулем
func WithHost(v string) webHookServerOptions {
	return func(whs *WebHookServer) error {
		whs.host = v

		return nil
	}
}

// WithName устанавливает наименование модуля (обязательно). Наименование основывается
// на имени организации или подразделения эксплуатирующем модуль. Например, gcm, rcmslx и т.д.
func WithName(v string) webHookServerOptions {
	return func(whs *WebHookServer) error {
		whs.name = v

		return nil
	}
}

// WithVersion устанавливает версию модуля (опционально)
func WithVersion(v string) webHookServerOptions {
	return func(whs *WebHookServer) error {
		whs.version = v

		return nil
	}
}
