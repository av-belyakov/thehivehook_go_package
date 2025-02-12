package webhookserver

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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
		wh.logger.Send("error", supportingfunctions.CustomError(err).Error())

		return
	}
	defer r.Body.Close()

	eventElement := datamodels.CaseEventElement{}
	err = json.Unmarshal(bodyByte, &eventElement)
	if err != nil {
		wh.logger.Send("error", supportingfunctions.CustomError(err).Error())

		return
	}

	switch eventElement.ObjectType {
	case "case":
		//формируем запрос на поиск дополнительной информации о кейсе, такой как observables
		//и ttp через модуль взаимодействия с API TheHive в TheHive
		go func() {
			wh.logger.Send("log_to_db", fmt.Sprintf("received case id '%d', a request is being sent for additional information about observable and ttl", eventElement.Object.CaseId))

			readyMadeEventCase, err := CreateEvenCase(r.Context(), eventElement.RootId, wh.chanInput)
			if err != nil {
				wh.logger.Send("error", supportingfunctions.CustomError(err).Error())

				return
			}

			caseEvent := map[string]interface{}{}
			if err := json.Unmarshal(bodyByte, &caseEvent); err != nil {
				wh.logger.Send("error", supportingfunctions.CustomError(err).Error())

				return
			}

			readyMadeEventCase.Source = wh.name
			readyMadeEventCase.Case = caseEvent

			wh.logger.Send("log_to_db", fmt.Sprintf("additional information on case id '%d' has been successfully received", eventElement.Object.CaseId))

			ec, err := json.Marshal(readyMadeEventCase)
			if err != nil {
				wh.logger.Send("error", supportingfunctions.CustomError(err).Error())

				return
			}

			//-------------------------------------------------------------------
			//----------- ЗАПИСЬ в файл ЭТО ТОЛЬКО ДЛЯ ТЕСТОВ -------------------
			//-------------------------------------------------------------------
			if str, err := supportingfunctions.NewReadReflectJSONSprint(ec); err == nil {
				wh.logger.Send("log_for_test", fmt.Sprintf("\n%s\n", str))
			}
			//-------------------------------------------------------------------
			if eventElement.Object.CaseId == 39100 {
				fmt.Println("func 'Routers' resived case id 39100")
			}

			//отправка данных в NATS
			sendData := NewChannelRequest()
			sendData.SetRootId(eventElement.RootId)
			sendData.SetElementType(eventElement.ObjectType)
			sendData.SetCaseId(strconv.Itoa(eventElement.Object.CaseId))
			sendData.SetCommand("send case")
			sendData.SetData(ec)

			wh.chanInput <- ChanFromWebHookServer{
				ForSomebody: "to nats",
				Data:        sendData,
			}

			wh.logger.Send("log_to_db", fmt.Sprintf("information on case id '%d' sending to NATS", eventElement.Object.CaseId))
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
			wh.logger.Send("error", supportingfunctions.CustomError(err).Error())

			return
		}

		//формируем запрос на поиск дополнительной информации по алерту такой как aler
		alertEvent := map[string]interface{}{}
		if err := json.Unmarshal(bodyByte, &alertEvent); err != nil {
			wh.logger.Send("error", supportingfunctions.CustomError(err).Error())

			return
		}

		readyMadeEventAlert.Source = wh.name
		readyMadeEventAlert.Event = alertEvent

		ea, err := json.Marshal(readyMadeEventAlert)
		if err != nil {
			wh.logger.Send("error", supportingfunctions.CustomError(err).Error())

			return
		}

		//-------------------------------------------------------------------
		//----------- ЗАПИСЬ в файл ЭТО ТОЛЬКО ДЛЯ ТЕСТОВ -------------------
		//-------------------------------------------------------------------
		if str, err := supportingfunctions.NewReadReflectJSONSprint(ea); err == nil {
			wh.logger.Send("log_for_test", fmt.Sprintf("\n%s\n", str))
		}
		//-------------------------------------------------------------------

		//fmt.Println("------ func 'RouteWebHook' ------- ALERT -----")
		//if res, err := json.MarshalIndent(readyMadeEventAlert, "", " "); err == nil {
		//	fmt.Println(string(res))
		//}
		//fmt.Println("------ func 'RouteWebHook' ------- ALERT -----")

		//отправка данных в NATS
		sendData := NewChannelRequest()
		sendData.SetRootId(eventElement.RootId)
		sendData.SetElementType(eventElement.ObjectType)
		sendData.SetCommand("send alert")
		sendData.SetData(ea)
		//sendData.SetChanOutput(chanResObservable)

		wh.chanInput <- ChanFromWebHookServer{
			ForSomebody: "to nats",
			Data:        sendData,
		}
	}
}
