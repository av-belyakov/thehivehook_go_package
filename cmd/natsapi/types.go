package natsapi

// apiNatsSettings настройки для API NATS
type apiNatsSettings struct {
	port        int
	host        string
	subscribers []SubscriberNATS
}

// SubscriberNATS абоненты NATS
type SubscriberNATS struct {
	Event      string   `validate:"required" yaml:"event"`
	Responders []string `yaml:"responders"`
}

// NatsAPIOptions функциональные опции
type NatsAPIOptions func(*apiNatsSettings) error

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
