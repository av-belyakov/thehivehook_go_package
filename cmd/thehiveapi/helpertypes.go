package thehiveapi

import "github.com/av-belyakov/thehivehook_go_package/cmd/commoninterfaces"

// RequestChannelTheHive тип применяется для передачи запроса в модуль thehiveapi от сторонних модулей
type RequestChannelTheHive struct {
	RequestId  string                                 //UUID идентификатор запроса
	RootId     string                                 //идентификатор по которому в TheHive будет выполнятся поиск
	Command    string                                 //команда
	Data       []byte                                 //набор данных
	ChanOutput chan commoninterfaces.ChannelResponser //канал ответа реализующий интерфейс commoninterfaces.ChannelResponser
}

// ResponseChannelTheHive структура ответа от API TheHive
type ResponseChannelTheHive struct {
	StatusCode int    //статус кода ответа
	RequestId  string //UUID идентификатор ответа (соответствует идентификатору запроса)
	Data       []byte //набор данных
}
