package webhookserver

import (
	"context"

	"github.com/av-belyakov/thehivehook_go_package/cmd/commoninterfaces"
)

// RequestFromWebHook структура запроса из модуля
type RequestFromWebHook struct {
	Data        any //набор данных
	ctx         context.Context
	RequestId   string                                 //id запроса
	ElementType string                                 //тип элемента
	RootId      string                                 //идентификатор по которому в TheHive будет выполнятся поиск
	CaseId      string                                 //идентификатор кейса в TheHive
	Command     string                                 //команда
	Order       string                                 //распоряжение
	ChanOutput  chan commoninterfaces.ChannelResponser //канал ответа реализующий интерфейс commoninterfaces.ChannelResponser
}

// ResponsToWebHook структура ответа в модуля
type ResponsToWebHook struct {
	Data       []byte //набор данных
	RequestId  string //UUID идентификатор ответа (соответствует идентификатору запроса)
	StatusCode int    //статус кода ответа
}
