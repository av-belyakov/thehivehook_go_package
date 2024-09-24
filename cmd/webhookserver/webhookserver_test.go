package webhookserver_test

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"testing"
)

type WebHookServer struct {
	port    int
	host    string
	version string
	ctx     context.Context
	server  *http.Server
}

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
	wh.server = server

	go func() {
		if errServer := server.ListenAndServe(); errServer != nil {
			log.Fatal(errServer)
		}
	}()

	log.Printf("server 'WebHookServer' was successfully launched, ip:%s, port:%d", wh.host, wh.port)
	<-wh.ctx.Done()
}

func (wh *WebHookServer) Shutdown(ctx context.Context) {
	wh.server.Shutdown(ctx)
}

func (wh *WebHookServer) RouteIndex(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)

		return
	}

	io.WriteString(w, "Hello WebHookServer version "+wh.version)
}

func (wh *WebHookServer) RouteWebHook(w http.ResponseWriter, r *http.Request) {
	fmt.Println("func 'RouteWebHook'")
	fmt.Println("Header:", r.Header)

	bodyByte, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		fmt.Println("ERROR:", err.Error())

		return
	}

	dst := make([]byte, base64.StdEncoding.DecodedLen(len(bodyByte)))
	strData, err := base64.StdEncoding.Decode(dst, bodyByte)
	if err != nil {
		fmt.Println("ERROR:", err.Error())

		return
	}

	data, err := json.MarshalIndent(strData, "", "  ")
	if err != nil {
		fmt.Println("ERROR: ", err)

		return
	}

	fmt.Println(string(data))
}

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

	webHookServer, errServer := New(ctx, "192.168.9.208", 5000)
	if errServer != nil {
		t.Fatal("create new server %w", errServer)
	}

	webHookServer.Start()
}
