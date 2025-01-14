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
func (wh *WebHookServer) Start(ctx context.Context) error {
	defer func() {
		close(wh.chanInput)
		wh.server.Shutdown(ctx)
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
			log.Fatal(errServer)
		}
	}()

	appStatus := fmt.Sprintf("%vproduction%v", Ansi_Bright_Blue, Ansi_Reset)
	envValue, ok := os.LookupEnv("GO_HIVEHOOK_MAIN")
	if ok && envValue == "development" {
		appStatus = fmt.Sprintf("%v%s%v", Ansi_Bright_Red, envValue, Ansi_Reset)
	}

	msg := fmt.Sprintf("Application '%s' v%s was successfully launched", appname.GetName(), appversion.GetVersion())
	log.Printf("%v%v%s%v\n", Ansi_Dark_Green_Background, Ansi_White, msg, Ansi_Reset)
	fmt.Printf("%v%vApplication status is '%s'.%v\n", Underlining, Ansi_Bright_Green, appStatus, Ansi_Reset)
	fmt.Printf("%vWebhook server settings:%v\n", Ansi_Bright_Green, Ansi_Reset)
	fmt.Printf("%v  name: %v%s%v\n", Ansi_Bright_Green, Ansi_Bright_Dark, wh.name, Ansi_Reset)
	fmt.Printf("%v  ip: %v%s%v\n", Ansi_Bright_Green, Ansi_Bright_Blue, wh.host, Ansi_Reset)
	fmt.Printf("%v  net port: %v%d%v\n", Ansi_Bright_Green, Ansi_Bright_Magenta, wh.port, Ansi_Reset)
	wh.logger.Send("info", strings.ToLower(msg))

	<-ctx.Done()

	return ctx.Err()
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
