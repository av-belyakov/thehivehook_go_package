package webhookserver

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"runtime"

	"github.com/av-belyakov/thehivehook_go_package/internal/supportingfunctions"
)

func (wh *WebHookServer) RouteIndex(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)

		return
	}

	io.WriteString(w, "Hello WebHookServer version "+wh.version)
}

func (wh *WebHookServer) RouteWebHook(w http.ResponseWriter, r *http.Request) {
	bodyByte, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		_, f, l, _ := runtime.Caller(0)
		wh.logger.Send("error", fmt.Sprintf(" '%s' %s:%d", err.Error(), f, l-3))

		return
	}

	//-------------------------------------------------------------------
	//----------- ЗАПИСЬ в файл ЭТО ТОЛЬКО ДЛЯ ТЕСТОВ -------------------
	//-------------------------------------------------------------------
	if str, err := supportingfunctions.NewReadReflectJSONSprint(bodyByte); err == nil {
		wh.logger.Send("log_for_test", str)
	}

	eventElement := EventElement{}
	err = json.Unmarshal(bodyByte, &eventElement)
	if err != nil {
		_, f, l, _ := runtime.Caller(0)
		wh.logger.Send("error", fmt.Sprintf(" '%s' %s:%d", err.Error(), f, l-2))

		return
	}

	switch eventElement.ObjectType {
	case "case":

	case "case_artifact":
		fmt.Println("Reseived object with object type:", eventElement.ObjectType)

	case "case_task":
		fmt.Println("Reseived object with object type:", eventElement.ObjectType)

	case "case_task_log":
		fmt.Println("Reseived object with object type:", eventElement.ObjectType)

	case "alert":
		fmt.Println("Reseived object with object type:", eventElement.ObjectType)

	}

	log.Println("Recived JSON size =", len(bodyByte))
}
