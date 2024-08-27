package elasticsearchapi

import (
	"github.com/elastic/go-elasticsearch/v8"
)

type SettingsInputChan struct {
	UUID    string
	Section string
	Command string
	Data    interface{}
}

type SettingsOutputChan struct {
	Section string
	Command string
	Data    interface{}
}

// ElasticSearchModule инициализированный модуль
// ChanInputModule - канал для отправки данных В модуль
// ChanOutputModule - канал для принятия данных ИЗ модуля
type ElasticSearchModule struct {
	ChanInputModule  chan SettingsInputChan
	ChanOutputModule chan SettingsOutputChan
}

type HandlerSendData struct {
	Client   *elasticsearch.Client
	Settings SettingsHandler
}

type SettingsHandler struct {
	Port   int
	Host   string
	User   string
	Passwd string
}

type listSensorId struct {
	sensors []string
}
