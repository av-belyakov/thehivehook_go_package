package webhookserver

import (
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
}
