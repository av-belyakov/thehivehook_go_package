package webhookserver

import "github.com/av-belyakov/thehivehook_go_package/cmd/commoninterfaces"

// RequestFromWebHook структура запроса из модуля
type RequestFromWebHook struct {
	RequestId  string                                 //id запроса
	RootId     string                                 //идентификатор по которому в TheHive будет выполнятся поиск
	CaseId     string                                 //идентификатор кейса в TheHive
	Command    string                                 //команда
	Order      string                                 //распоряжение
	Data       interface{}                            //набор данных
	ChanOutput chan commoninterfaces.ChannelResponser //канал ответа реализующий интерфейс commoninterfaces.ChannelResponser
}

// ResponsToWebHook структура ответа в модуля
type ResponsToWebHook struct {
	StatusCode int    //статус кода ответа
	RequestId  string //UUID идентификатор ответа (соответствует идентификатору запроса)
	Data       []byte //набор данных
}
