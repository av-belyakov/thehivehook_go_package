package natsapi

import (
	"github.com/nats-io/nats.go"

	"github.com/av-belyakov/thehivehook_go_package/cmd/commoninterfaces"
)

// apiNatsSettings настройки для API NATS
type apiNatsModule struct {
	host             string
	subscriptions    subscription
	logger           commoninterfaces.Logger
	natsConnection   *nats.Conn
	receivingChannel chan commoninterfaces.ChannelRequester
	sendingChannel   chan commoninterfaces.ChannelRequester
	cachettl         int
	port             int
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
