package webhookserver

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"runtime"
	"strconv"

	"github.com/av-belyakov/thehivehook_go_package/internal/datamodels"
	"github.com/av-belyakov/thehivehook_go_package/internal/supportingfunctions"
)

// RouteIndex маршрут при обращении к '/'
func (wh *WebHookServer) RouteIndex(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)

		return
	}

	io.WriteString(w, "Hello WebHookServer version "+wh.version)
}

// RouteWebHook маршрут при обращении к '/webhook'
func (wh *WebHookServer) RouteWebHook(w http.ResponseWriter, r *http.Request) {
	bodyByte, err := io.ReadAll(r.Body)
	if err != nil {
		_, f, l, _ := runtime.Caller(0)
		wh.logger.Send("error", fmt.Sprintf(" '%s' %s:%d", err.Error(), f, l-3))

		return
	}
	defer r.Body.Close()

	//-------------------------------------------------------------------
	//----------- ЗАПИСЬ в файл ЭТО ТОЛЬКО ДЛЯ ТЕСТОВ -------------------
	//-------------------------------------------------------------------
	if str, err := supportingfunctions.NewReadReflectJSONSprint(bodyByte); err == nil {
		wh.logger.Send("log_for_test", fmt.Sprintf("\n%s\n", str))
	}
	//-------------------------------------------------------------------

	eventElement := datamodels.CaseEventElement{}
	err = json.Unmarshal(bodyByte, &eventElement)
	if err != nil {
		_, f, l, _ := runtime.Caller(0)
		wh.logger.Send("error", fmt.Sprintf(" '%s' %s:%d", err.Error(), f, l-2))

		return
	}

	fmt.Println("Received object with object type:", eventElement.ObjectType)
	fmt.Println("Received JSON size =", len(bodyByte))

	switch eventElement.ObjectType {
	case "case":
		//формируем запрос на поиск дополнительной информации о кейсе, такой как observables
		//и ttp через модуль взаимодействия с API TheHive в TheHive
		go func() {
			fmt.Println("1111111 ------ func 'RouteWebHook' ------- CASE -----")

			readyMadeEventCase, err := CreateEvenCase(eventElement.RootId, wh.chanInput)
			if err != nil {
				fmt.Println("------ func 'RouteWebHook' ------- CASE ----- ERROR 1:", err.Error())

				_, f, l, _ := runtime.Caller(0)
				wh.logger.Send("error", fmt.Sprintf(" '%s' %s:%d", err.Error(), f, l-2))

				return
			}

			fmt.Println("1222222 ------ func 'RouteWebHook' ------- CASE -----")

			caseEvent := map[string]interface{}{}
			if err := json.Unmarshal(bodyByte, &caseEvent); err != nil {
				fmt.Println("------ func 'RouteWebHook' ------- CASE ----- ERROR 2:", err.Error())

				_, f, l, _ := runtime.Caller(0)
				wh.logger.Send("error", fmt.Sprintf(" '%s' %s:%d", err.Error(), f, l-1))

				return
			}

			fmt.Println("1333333 ------ func 'RouteWebHook' ------- CASE -----")

			readyMadeEventCase.Source = wh.name
			readyMadeEventCase.Case = caseEvent

			fmt.Println("------ func 'RouteWebHook' ------- CASE -----")
			if res, err := json.MarshalIndent(readyMadeEventCase, "", " "); err == nil {
				fmt.Println(string(res))
			}
			fmt.Println("------ func 'RouteWebHook' ------- CASE -----")

			ec, err := json.Marshal(readyMadeEventCase)
			if err != nil {
				_, f, l, _ := runtime.Caller(0)
				wh.logger.Send("error", fmt.Sprintf(" '%s' %s:%d", err.Error(), f, l-1))

				return
			}

			//отправка данных в NATS
			sendData := NewChannelRequest()
			sendData.SetRootId(eventElement.RootId)
			sendData.SetElementType(eventElement.ObjectType)
			sendData.SetCaseId(strconv.Itoa(eventElement.Object.CaseId))
			sendData.SetCommand("send case")
			sendData.SetData(ec)

			wh.chanInput <- ChanFromWebHookServer{
				ForSomebody: "for nats",
				Data:        sendData,
			}
		}()

	case "case_artifact":
	case "case_task":
	case "case_task_log":
	case "alert":
		if eventElement.Operation == "delete" {
			return
		}

		//*****************************************************************
		//ВНИМАНИЕ!!! На данный момент этот модуль еще ничего не обогащает
		//нужно ли делать модуль обогатитель пока не ясно
		//пока до решения этого впроса я еще не дошёл
		readyMadeEventAlert, err := CreateEvenAlert(eventElement.RootId, wh.chanInput)
		if err != nil {
			_, f, l, _ := runtime.Caller(0)
			wh.logger.Send("error", fmt.Sprintf(" '%s' %s:%d", err.Error(), f, l-2))

			return
		}

		//формируем запрос на поиск дополнительной информации по алерту такой как aler
		alertEvent := map[string]interface{}{}
		if err := json.Unmarshal(bodyByte, &alertEvent); err != nil {
			_, f, l, _ := runtime.Caller(0)
			wh.logger.Send("error", fmt.Sprintf(" '%s' %s:%d", err.Error(), f, l-1))

			return
		}

		readyMadeEventAlert.Source = wh.name
		readyMadeEventAlert.Event = alertEvent

		//fmt.Println("------ func 'RouteWebHook' ------- ALERT -----")
		//if res, err := json.MarshalIndent(readyMadeEventAlert, "", " "); err == nil {
		//	fmt.Println(string(res))
		//}
		//fmt.Println("------ func 'RouteWebHook' ------- ALERT -----")

		ea, err := json.Marshal(readyMadeEventAlert)
		if err != nil {
			_, f, l, _ := runtime.Caller(0)
			wh.logger.Send("error", fmt.Sprintf(" '%s' %s:%d", err.Error(), f, l-1))

			return
		}

		//отправка данных в NATS
		sendData := NewChannelRequest()
		sendData.SetRootId(eventElement.RootId)
		sendData.SetElementType(eventElement.ObjectType)
		sendData.SetCommand("send alert")
		sendData.SetData(ea)
		//sendData.SetChanOutput(chanResObservable)

		wh.chanInput <- ChanFromWebHookServer{
			ForSomebody: "for nats",
			Data:        sendData,
		}
	}
}
