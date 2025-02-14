package thehiveapi

import (
	"encoding/json"

	"github.com/av-belyakov/thehivehook_go_package/cmd/commoninterfaces"
)

// RequestChannelTheHive тип применяется для передачи запроса в модуль thehiveapi от сторонних модулей
type RequestChannelTheHive struct {
	Data       interface{}                            //какие то данные
	RequestId  string                                 //UUID идентификатор запроса
	RootId     string                                 //идентификатор по которому в TheHive будет выполнятся поиск
	CaseId     string                                 //идентификатор кейса в TheHive
	Command    string                                 //команда
	ChanOutput chan commoninterfaces.ChannelResponser //канал ответа реализующий интерфейс commoninterfaces.ChannelResponser
}

// ResponseChannelTheHive структура ответа от API TheHive
type ResponseChannelTheHive struct {
	Err        error  //объект ошибки
	RequestId  string //UUID идентификатор ответа (соответствует идентификатору запроса)
	Data       []byte //набор данных
	StatusCode int    //статус кода ответа
}

// RequestCommand структура с командами для обработки модулем
type RequestCommand struct {
	ByteData  json.RawMessage `json:"byte_ data"` //набор данных в бинарном виде которые обрабатываются отдельно
	Service   string          `json:"service"`    //наименование сервиса
	Command   string          `json:"command"`    //команда
	RootId    string          `json:"root_id"`    //основной id, как правило это rootId case или alert
	CaseId    string          `json:"case_id"`    //id кейса
	Username  string          `json:"username"`   //имя пользователя, необходим если нужно указать пользователя выполнившего действие
	FieldName string          `json:"field_name"` //некое ключевое поле
	Value     string          `json:"value"`      //устанавливаемое значение
}

// SpecialObjectForCache является вспомогательным типом который реализует интерфейс
// CacheStorageFuncHandler[T any] где в методе Comparison(objFromCache T) bool необходимо
// реализовать подробное сравнение объекта типа T.
// Нужен для пакета cachingstoragewithqueue
type SpecialObjectForCache[T any] struct {
	object      T
	handlerFunc func(int) bool
	id          string
}

// LogWrite вспомогательный тип применяемый для логирования
type LogWrite struct {
	logger commoninterfaces.Logger
}
