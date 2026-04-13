package natsapi

import (
	"github.com/nats-io/nats.go"

	"github.com/av-belyakov/thehivehook_go_package/v2/internal/interfaces"
)

// NatsApi модуль API NATS
type NatsApi struct {
	subscriptions      subscription
	host               string
	nameRegionalObject string
	cachettl           int
	port               int
	logger             interfaces.Logger
	natsConnection     *nats.Conn
	chanInputModule    chan InputData
	chanOutputModule   chan OutputData
}

type subscription struct {
	senderCase      string
	senderAlert     string
	listenerCommand string
}

// InputData входящие в модуль данные (alert и case)
type InputData struct {
	Data        []byte
	RootId      string
	CaseId      string
	ElementType string
}

// OutputData исходящие из модуля данные (команды на добавление tags, custom fields)
type OutputData struct {
	Data       []byte        //набор данных
	ChanDone   chan struct{} //канал информирующий о закрытии канала ChanOutput
	ChanOutput chan []byte   //канал ответа реализующий интерфейс commoninterfaces.ChannelResponser
}

// ApiApiOptions функциональные опции
type NatsApiOptions func(*NatsApi) error
