package webhookserver

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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

	resBodyByte, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		fmt.Println("Error: ", err)

		return
	}

	data, err := json.MarshalIndent(resBodyByte, "", "  ")
	if err != nil {
		fmt.Println("Error: ", err)

		return
	}

	fmt.Println(string(data))
}
