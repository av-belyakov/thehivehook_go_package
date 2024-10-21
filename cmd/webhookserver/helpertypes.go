package webhookserver

import "github.com/av-belyakov/thehivehook_go_package/cmd/commoninterfaces"

// RequestFromWebHook структура запроса из модуля
type RequestFromWebHook struct {
	RequestId  string                                 //id запроса
	RootId     string                                 //идентификатор по которому в TheHive будет выполнятся поиск
	Command    string                                 //команда
	Data       []byte                                 //набор данных
	ChanOutput chan commoninterfaces.ChannelResponser //канал ответа реализующий интерфейс commoninterfaces.ChannelResponser
}

// ResponsToWebHook структура ответа в модуля
type ResponsToWebHook struct {
	StatusCode int    //статус кода ответа
	RequestId  string //UUID идентификатор ответа (соответствует идентификатору запроса)
	Data       []byte //набор данных
}

// EventElement типовой элемент описывающий события приходящие из TheHive
type EventElement struct {
	Operation  string              `json:"operation"`  //тип операции
	ObjectType string              `json:"objectType"` //тип объекта
	RootId     string              `json:"rootId"`     //основной идентификатор объекта
	Object     ObjectEventElement  `json:"object"`     //частичная информация по объекту
	Details    DetailsEventElement `json:"details"`    //частичные детали по объекту
}

// ObjectEventElement содержит информацию из поля 'object' приходящего из TheHive элемента
type ObjectEventElement struct {
	CreatedAt int64 `json:"createdAt"`
}

// DetailsEventElement содержит информацию из поля 'details'
type DetailsEventElement struct {
	Status string `json:"status"`
}
