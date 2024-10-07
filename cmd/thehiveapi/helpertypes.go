package thehiveapi

import "github.com/av-belyakov/thehivehook_go_package/cmd/commoninterfaces"

// RequestChannelTheHive тип применяется для передачи запроса в
// модуль thehiveapi от сторонних модулей
// RequestId - UUID идентификатор запроса
// RootId - идентификатор по которому в TheHive будет выполнятся поиск
// Command - команда
// ChanOutput - канал ответа реализующий интерфейс commoninterfaces.ChannelResponser
type RequestChannelTheHive struct {
	RequestId  string
	RootId     string
	Command    string
	ChanOutput chan commoninterfaces.ChannelResponser
}

// ResponseChannelTheHive структура ответа от API TheHive
type ResponseChannelTheHive struct {
	StatusCode int
	RequestId  string
	Data       []byte
}
