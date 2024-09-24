package webhookserver_test

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"reflect"
	"strings"
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

	str, err := NewReadReflectJSONSprint(bodyByte)
	if err != nil {
		fmt.Println("ERROR:", err.Error())

		return
	}

	fmt.Printf("__________\n%s____________\n\n", str)

	/*data, err := json.MarshalIndent(str, "", "  ")
	if err != nil {
		fmt.Println("ERROR: ", err)

		return
	}

	fmt.Println(string(data))*/
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

// NewReadReflectJSONSprint функция выполняет вывод JSON сообщения в виде текста
// Для данной функции не требуется описание текст, так как обработка JSON сообщения
// осуществляется с помощью пакета reflect
func NewReadReflectJSONSprint(b []byte) (string, error) {
	var str string
	errSrc := "error decoding the json file, it may be empty"

	listMap := map[string]interface{}{}
	if err := json.Unmarshal(b, &listMap); err == nil {
		if len(listMap) == 0 {
			return str, fmt.Errorf(errSrc)
		}

		return readReflectMapSprint(listMap, 0), err
	}

	listSlice := []interface{}{}
	if err := json.Unmarshal(b, &listSlice); err == nil {
		if len(listSlice) == 0 {
			return str, fmt.Errorf(errSrc)
		}

		return readReflectSliceSprint(listSlice, 0), err
	}

	return str, fmt.Errorf("the contents of the file are not Map or Slice")
}

func readReflectAnyTypeSprint(name interface{}, anyType interface{}, num int) string {
	var (
		nameStr string
		str     strings.Builder = strings.Builder{}
	)

	r := reflect.TypeOf(anyType)
	ws := GetWhitespace(num)

	if n, ok := name.(int); ok {
		nameStr = fmt.Sprintf("%s%v.", ws, n+1)
	} else if n, ok := name.(string); ok {
		nameStr = fmt.Sprintf("%s\"%s\":", ws, n)
	}

	if r == nil {
		return str.String()
	}

	switch r.Kind() {
	case reflect.String:
		dataStr := reflect.ValueOf(anyType).String()

		if nameStr == "description" {
			dataStr = strings.ReplaceAll(dataStr, "\t", "")
			dataStr = strings.ReplaceAll(dataStr, "\n", "")
		}

		str.WriteString(fmt.Sprintf("%s \"%s\"\n", nameStr, dataStr))

	case reflect.Int, reflect.Int16, reflect.Int32, reflect.Int64:
		str.WriteString(fmt.Sprintf("%s %d\n", nameStr, reflect.ValueOf(anyType).Int()))

	case reflect.Uint, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		str.WriteString(fmt.Sprintf("%s %d\n", nameStr, reflect.ValueOf(anyType).Uint()))

	case reflect.Float32, reflect.Float64:
		str.WriteString(fmt.Sprintf("%s %v\n", nameStr, int(reflect.ValueOf(anyType).Float())))

	case reflect.Bool:
		str.WriteString(fmt.Sprintf("%s %v\n", nameStr, reflect.ValueOf(anyType).Bool()))
	}

	return str.String()
}

func readReflectMapSprint(list map[string]interface{}, num int) string {
	var str strings.Builder = strings.Builder{}
	ws := GetWhitespace(num)

	for k, v := range list {
		r := reflect.TypeOf(v)

		if r == nil {
			return str.String()
		}

		switch r.Kind() {
		case reflect.Map:
			if v, ok := v.(map[string]interface{}); ok {
				str.WriteString(fmt.Sprintf("%s%s:\n", ws, k))
				str.WriteString(readReflectMapSprint(v, num+1))
			}

		case reflect.Slice:
			if v, ok := v.([]interface{}); ok {
				str.WriteString(fmt.Sprintf("%s%s:\n", ws, k))
				str.WriteString(readReflectSliceSprint(v, num+1))
			}

		case reflect.Array:
			str.WriteString(fmt.Sprintf("%s: %s (it is array)\n", k, reflect.ValueOf(v).String()))

		default:
			str.WriteString(readReflectAnyTypeSprint(k, v, num))
		}
	}

	return str.String()
}

func readReflectSliceSprint(list []interface{}, num int) string {
	var str strings.Builder = strings.Builder{}
	ws := GetWhitespace(num)

	for k, v := range list {
		r := reflect.TypeOf(v)

		if r == nil {
			return str.String()
		}

		switch r.Kind() {
		case reflect.Map:
			if v, ok := v.(map[string]interface{}); ok {
				str.WriteString(fmt.Sprintf("%s%d.\n", ws, k+1))
				str.WriteString(readReflectMapSprint(v, num+1))
			}

		case reflect.Slice:
			if v, ok := v.([]interface{}); ok {
				str.WriteString(fmt.Sprintf("%s%d.\n", ws, k+1))
				str.WriteString(readReflectSliceSprint(v, num+1))
			}

		case reflect.Array:
			str.WriteString(fmt.Sprintf("%d. %s (it is array)\n", k, reflect.ValueOf(v).String()))

		default:
			str.WriteString(readReflectAnyTypeSprint(k, v, num))
		}
	}

	return str.String()
}

// GetWhitespace возвращает необходимое количество пробелов
func GetWhitespace(num int) string {
	var str string

	if num == 0 {
		return str
	}

	for i := 0; i < num; i++ {
		str += "  "
	}

	return str
}
