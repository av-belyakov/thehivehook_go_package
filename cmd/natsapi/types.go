package natsapi

import (
	"github.com/nats-io/nats.go"

	"github.com/av-belyakov/thehivehook_go_package/cmd/commoninterfaces"
	"github.com/av-belyakov/thehivehook_go_package/internal/datamodels"
)

// apiNatsSettings настройки для API NATS
type apiNatsModule struct {
	natsConnection     *nats.Conn
	logger             commoninterfaces.Logger
	host               string
	nameRegionalObject string
	subscriptions      subscription
	//receivingChannel   chan commoninterfaces.ChannelRequester
	receivingChannel chan datamodels.RequestChan
	//sendingChannel     chan commoninterfaces.ChannelRequester
	sendingChannel chan datamodels.RequestChan
	cachettl       int
	port           int
}

type subscription struct {
	senderCase      string
	senderAlert     string
	listenerCommand string
}

// NatsApiOptions функциональные опции
type NatsApiOptions func(*apiNatsModule) error

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
