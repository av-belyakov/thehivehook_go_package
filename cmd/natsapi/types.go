package natsapi

// ModuleNATS инициализированный модуль
// ChanOutputMISP - канал для отправки полученных данных из модуля
type ModuleNATS struct {
	chanOutputNATS chan SettingsOutputChan
}

// SettingsOutputChan канал вывода данных из модуля
// MsgId - id сообщения
// SubjectType - тип подписки
// Data - данные
type SettingsOutputChan struct {
	MsgId       string
	SubjectType string
	Data        []byte
}

// SettingsInputChan канал для приема данных в модуль
type SettingsInputChan struct {
	Command, EventId, TaskId string
}

func (mnats *ModuleNATS) GetDataReceptionChannel() <-chan SettingsOutputChan {
	return mnats.chanOutputNATS
}

func (mnats *ModuleNATS) SendingData(data SettingsOutputChan) {
	mnats.chanOutputNATS <- data
}
