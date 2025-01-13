// Основной модуль приложения
package webhookserver

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/av-belyakov/thehivehook_go_package/cmd/commoninterfaces"
	"github.com/av-belyakov/thehivehook_go_package/internal/appname"
	"github.com/av-belyakov/thehivehook_go_package/internal/appversion"
)

// New конструктор webhookserver принимает функциональные опции для настройки модуля перед запуском
func New(logger commoninterfaces.Logger, opts ...webHookServerOptions) (*WebHookServer, <-chan ChanFromWebHookServer, error) {
	chanOutput := make(chan ChanFromWebHookServer)

	whs := &WebHookServer{
		name:      "nobody",
		version:   "0.1.1",
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
func (wh *WebHookServer) Start(ctx context.Context) {
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

	appStatus := fmt.Sprintf("%vproduction%v", ansiBrightBlue, ansiReset)
	envValue, ok := os.LookupEnv("GO_HIVEHOOK_MAIN")
	if ok && envValue == "development" {
		appStatus = fmt.Sprintf("%v%s%v", ansiBrightRed, envValue, ansiReset)
	}

	msg := fmt.Sprintf("Application '%s' v%s was successfully launched", appname.GetName(), appversion.GetVersion())
	log.Printf("%v%v%v%s%v\n", ansiDarkGreenBackground, boldFont, ansiWhite, msg, ansiReset)
	log.Printf("%vApplication status is '%s'.%v", ansiBrightGreen, appStatus, ansiReset)
	log.Printf("%vWebhook server settings:%v", ansiBrightGreen, ansiReset)
	log.Printf("%v  name: %v%s%v", ansiBrightGreen, ansiBrightDark, wh.name, ansiReset)
	log.Printf("%v  ip: %v%s%v", ansiBrightGreen, ansiBrightBlue, wh.host, ansiReset)
	log.Printf("%v  net port: %v%d%v", ansiBrightGreen, ansiBrightMagenta, wh.port, ansiReset)
	wh.logger.Send("info", strings.ToLower(msg))

	<-ctx.Done()
	close(wh.chanInput)
}

// Shutdown завершает работу модуля
func (wh *WebHookServer) Shutdown(ctx context.Context) {
	wh.server.Shutdown(ctx)
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
