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
	//-------------------------------------------------------------------

	eventElement := EventElement{}
	err = json.Unmarshal(bodyByte, &eventElement)
	if err != nil {
		_, f, l, _ := runtime.Caller(0)
		wh.logger.Send("error", fmt.Sprintf(" '%s' %s:%d", err.Error(), f, l-2))

		return
	}

	//******* выполняем проверку было подобное событие получено ранее *******
	//
	// !!!!!!!!!!!!!!!!!!!!!
	// здесь надо продумать как предотвращать цикличное взаимодействие между
	// элементами находящимися за NATS и генерирующие команды на изменение
	// кейса TheHive и постоянными событиями являющимися результатом этих изменений
	// !!!!!!!!!!!!!!!!!!!!!
	//
	//_, isExistElement := wh.storage.GetValue(eventElement.GetEventId())
	//exception := eventElement.Details.Status != "Resolved"
	//if exception || !isExistElement {
	//
	//	fmt.Println("!!! Received repeated TheHive element with rootId =", eventElement.RootId)
	//
	//	return
	//}

	//записываем информацию о событии полученном из TheHive
	//idStorage := wh.storage.SetValue(eventElement.GetEventId(), "first")

	fmt.Println("Received object with object type:", eventElement.ObjectType)
	log.Println("Received JSON size =", len(bodyByte))

	switch eventElement.ObjectType {
	case "case":
		//формируем запрос на поиск дополнительной информации о кейсе, такой как observables
		//и ttp через модуль взаимодействия с API TheHive в TheHive
		go func() {
			readyMadeEventCase, err := CreateEvenCase(idStorage, eventElement.RootId, wh.chanInput)
			if err != nil {
				_, f, l, _ := runtime.Caller(0)
				wh.logger.Send("error", fmt.Sprintf(" '%s' %s:%d", err.Error(), f, l-2))

				return
			}

			caseEvent := map[string]interface{}{}
			if err := json.Unmarshal(bodyByte, &caseEvent); err != nil {
				_, f, l, _ := runtime.Caller(0)
				wh.logger.Send("error", fmt.Sprintf(" '%s' %s:%d", err.Error(), f, l-1))

				return
			}

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
			sendData.SetRequestId(idStorage)
			sendData.SetRootId(eventElement.RootId)
			sendData.SetCommand("send case")
			sendData.SetData(ec)
			//sendData.SetChanOutput(chanResObservable)
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

		readyMadeEventAlert, err := CreateEvenAlert(idStorage, eventElement.RootId, wh.chanInput)
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

		fmt.Println("------ func 'RouteWebHook' ------- ALERT -----")
		if res, err := json.MarshalIndent(readyMadeEventAlert, "", " "); err == nil {
			fmt.Println(string(res))
		}
		fmt.Println("------ func 'RouteWebHook' ------- ALERT -----")

		ea, err := json.Marshal(readyMadeEventAlert)
		if err != nil {
			_, f, l, _ := runtime.Caller(0)
			wh.logger.Send("error", fmt.Sprintf(" '%s' %s:%d", err.Error(), f, l-1))

			return
		}

		//отправка данных в NATS
		sendData := NewChannelRequest()
		sendData.SetRequestId(idStorage)
		sendData.SetRootId(eventElement.RootId)
		sendData.SetCommand("send alert")
		sendData.SetData(ea)
		//sendData.SetChanOutput(chanResObservable)
		wh.chanInput <- ChanFromWebHookServer{
			ForSomebody: "for nats",
			Data:        sendData,
		}
	}
}
