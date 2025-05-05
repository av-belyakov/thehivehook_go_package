package webhookserver

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/av-belyakov/thehivehook_go_package/internal/supportingfunctions"
)

// RouteIndex маршрут при обращении к '/'
func (wh *WebHookServer) RouteIndex(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)

		return
	}

	status := "production"
	if os.Getenv("GO_HIVEHOOK_MAIN") == "development" {
		status = os.Getenv("GO_HIVEHOOK_MAIN")
	}

	numberHours := int(time.Since(wh.timeStart).Hours())

	io.WriteString(w,
		fmt.Sprintf("Hello, WebHookServer version %s, application status:'%s'. %d hours have passed since the launch of the application.\n\n%s\n",
			wh.version,
			status,
			numberHours,
			printMemStats()))
}

// RouteWebHook маршрут при обращении к '/webhook'
func (wh *WebHookServer) RouteWebHook(w http.ResponseWriter, r *http.Request) {
	eventElement := map[string]any{}
	bodyByte, err := io.ReadAll(r.Body)
	if err != nil {
		wh.logger.Send("error", supportingfunctions.CustomError(err).Error())

		return
	}
	defer func() {
		r.Body.Close()
		bodyByte = []byte(nil)
	}()

	//----------------------------------------------------------------------
	//----------- запись в файл принятых в обработку объектов --------------
	//----------------------------------------------------------------------
	go func(d []byte) {
		if str, err := supportingfunctions.NewReadReflectJSONSprint(d); err == nil {
			if str != "" {
				wh.logger.Send("accepted_objects", fmt.Sprintf("\n%s\n", str))
			}
		}
	}(bodyByte)

	err = json.Unmarshal(bodyByte, &eventElement)
	if err != nil {
		wh.logger.Send("error", supportingfunctions.CustomError(err).Error())

		return
	}

	objectType, err := GetObjectType(eventElement)
	if err != nil {
		wh.logger.Send("error", supportingfunctions.CustomError(err).Error())

		return
	}

	rootId, err := GetRootId(eventElement)
	if err != nil {
		wh.logger.Send("error", supportingfunctions.CustomError(err).Error())

		return
	}

	operation, err := GetOperation(eventElement)
	if err != nil {
		wh.logger.Send("error", supportingfunctions.CustomError(err).Error())
	}

	switch objectType {
	case "case":
		caseId, err := GetCaseId(eventElement)
		if err != nil {
			wh.logger.Send("error", supportingfunctions.CustomError(err).Error())
		}

		wh.logger.Send("info", fmt.Sprintf("received caseId:'%d', rootId:'%s', operation:'%s', a request is being sent for additional information about 'observable' and 'ttl' objects", caseId, rootId, operation))

		//формируем запрос на поиск дополнительной информации о кейсе, такой как observables
		//и ttp через модуль взаимодействия с API TheHive в TheHive
		readyMadeEventCase, err := CreateEvenCase(r.Context(), rootId, caseId, wh.chanInput)
		if err != nil {
			if !errors.Is(err, &CreateCaseError{Type: "context"}) {
				wh.logger.Send("error", supportingfunctions.CustomError(err).Error())
			}

			return
		}

		readyMadeEventCase.Source = wh.name
		readyMadeEventCase.Case = eventElement

		wh.logger.Send("info", fmt.Sprintf("additional information on case id '%d' has been successfully received", caseId))

		ec, err := json.Marshal(readyMadeEventCase)
		if err != nil {
			wh.logger.Send("error", supportingfunctions.CustomError(err).Error())

			return
		}

		//передача в NATS
		sendData := NewChannelRequest()
		sendData.SetRootId(rootId)
		sendData.SetElementType(objectType)
		sendData.SetCaseId(strconv.Itoa(caseId))
		sendData.SetCommand("send case")
		sendData.SetData(ec)

		wh.chanInput <- ChanFromWebHookServer{
			ForSomebody: "to nats",
			Data:        sendData,
		}

		wh.logger.Send("info", fmt.Sprintf("information on case id '%d' sending to NATS", caseId))

	case "case_artifact":
	case "case_task":
	case "case_task_log":
	case "alert":
		wh.logger.Send("info", fmt.Sprintf("received alert rootId:'%s', operation:'%s'", rootId, operation))

		if operation == "delete" {
			return
		}

		//*****************************************************************
		//ВНИМАНИЕ!!! На данный момент этот модуль еще ничего не обогащает
		//нужно ли делать модуль обогатитель пока не ясно
		//пока до решения этого впроса я еще не дошёл
		readyMadeEventAlert, err := CreateEvenAlert(r.Context(), rootId, wh.chanInput)
		if err != nil {
			wh.logger.Send("error", supportingfunctions.CustomError(err).Error())

			return
		}

		readyMadeEventAlert.Source = wh.name
		readyMadeEventAlert.Event = eventElement

		ea, err := json.Marshal(readyMadeEventAlert)
		if err != nil {
			wh.logger.Send("error", supportingfunctions.CustomError(err).Error())

			return
		}

		//передача в NATS
		sendData := NewChannelRequest()
		sendData.SetRootId(rootId)
		sendData.SetElementType(objectType)
		sendData.SetCommand("send alert")
		sendData.SetData(ea)

		wh.chanInput <- ChanFromWebHookServer{
			ForSomebody: "to nats",
			Data:        sendData,
		}
	}
}
