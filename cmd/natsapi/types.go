package natsapi

import (
	"github.com/nats-io/nats.go"

	"github.com/av-belyakov/thehivehook_go_package/cmd/natsapi/storage"
	"github.com/av-belyakov/thehivehook_go_package/internal/interfaces"
)

// apiNatsSettings настройки для API NATS
type apiNatsModule struct {
	logger             interfaces.Logger
	subscriptions      subscription
	host               string
	nameRegionalObject string
	cachettl           int
	port               int
	storageCache       *storage.StorageAcceptedCommands
	natsConnection     *nats.Conn
	receivingChannel   chan interfaces.ChannelRequester
	sendingChannel     chan interfaces.ChannelRequester
}

// NatsApi модуль NATS API
type NatsApi struct {
	subscriptions      subscription
	host               string
	nameRegionalObject string
	logger             interfaces.Logger
	natsConnection     *nats.Conn
	cachettl           int
	port               int
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

// NatsApiOptions функциональные опции
type NatsApiOptions func(*apiNatsModule) error

// NatsApiOptions функциональные опции
type NatsApiOptions_New func(*NatsApi) error

// ModuleNATS инициализированный модуль
type ModuleNATS struct {
	chanOutputNATS chan SettingsOutputChan //канал для отправки полученных данных из модуля
}

// SettingsOutputChan канал вывода данных из модуля
type SettingsOutputChan struct {
	Data        []byte //набор данных
	MsgId       string //id сообщения
	SubjectType string //тип подписки
}

// SettingsInputChan канал для приема данных в модуль
type SettingsInputChan struct {
	Command, EventId, TaskId string
}
