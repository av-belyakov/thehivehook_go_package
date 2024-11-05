package natsapi

import (
	"github.com/av-belyakov/thehivehook_go_package/cmd/commoninterfaces"
	temporarystoarge "github.com/av-belyakov/thehivehook_go_package/cmd/natsapi/temporarystorage"
)

// apiNatsSettings настройки для API NATS
type apiNatsModule struct {
	cachettl             int
	port                 int
	host                 string
	subscribers          []SubscriberNATS
	logger               commoninterfaces.Logger
	receivingChannel     chan commoninterfaces.ChannelRequester
	temporaryStorage     *temporarystoarge.TemporaryStorage
	cacheRunningFunction commoninterfaces.CacheFuncRunner
}

// SubscriberNATS абоненты NATS
type SubscriberNATS struct {
	Event      string   `validate:"required" yaml:"event"`
	Responders []string `yaml:"responders"`
}

// NatsApiOptions функциональные опции
type NatsApiOptions func(*apiNatsModule) error

// ModuleNATS инициализированный модуль
type ModuleNATS struct {
	chanOutputNATS chan SettingsOutputChan //канал для отправки полученных данных из модуля
}

// SettingsOutputChan канал вывода данных из модуля
type SettingsOutputChan struct {
	MsgId       string //id сообщения
	SubjectType string //тип подписки
	Data        []byte //набор данных
}

// SettingsInputChan канал для приема данных в модуль
type SettingsInputChan struct {
	Command, EventId, TaskId string
}
