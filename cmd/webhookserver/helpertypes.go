package webhookserver

import "github.com/av-belyakov/thehivehook_go_package/cmd/commoninterfaces"

// RequestChannelTheHive структура запроса в модуль
// RequestId - UUID идентификатор запроса
// RootId - идентификатор по которому в TheHive будет выполнятся поиск
// Command - команда
// ChanOutput - канал ответа реализующий интерфейс commoninterfaces.ChannelResponser
type RequestTheHive struct {
	RequestId  string
	RootId     string
	Command    string
	ChanOutput chan commoninterfaces.ChannelResponser
}

// ResponseChannelTheHive структура ответа от модуля
type ResponsTheHive struct {
	StatusCode int
	RequestId  string
	Data       []byte
}
