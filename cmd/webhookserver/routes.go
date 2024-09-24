package webhookserver

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"runtime"
)

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
		_, f, l, _ := runtime.Caller(0)
		wh.logger.Send("error", fmt.Sprintf(" '%s' %s:%d", err.Error(), f, l-3))

		return
	}

	strData, err := base64.StdEncoding.DecodeString(string(bodyByte))
	if err != nil {
		_, f, l, _ := runtime.Caller(0)
		wh.logger.Send("error", fmt.Sprintf(" '%s' %s:%d", err.Error(), f, l-3))

		return
	}

	data, err := json.MarshalIndent(strData, "", "  ")
	if err != nil {
		fmt.Println("Error: ", err)

		return
	}

	fmt.Println(string(data))
}
