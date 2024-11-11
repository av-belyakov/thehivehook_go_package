package thehiveapi

import (
	"encoding/json"

	"github.com/av-belyakov/thehivehook_go_package/cmd/commoninterfaces"
)

// RequestChannelTheHive тип применяется для передачи запроса в модуль thehiveapi от сторонних модулей
type RequestChannelTheHive struct {
	RequestId  string                                 //UUID идентификатор запроса
	RootId     string                                 //идентификатор по которому в TheHive будет выполнятся поиск
	CaseId     string                                 //идентификатор кейса в TheHive
	Command    string                                 //команда
	Data       interface{}                            //какие то данные
	ChanOutput chan commoninterfaces.ChannelResponser //канал ответа реализующий интерфейс commoninterfaces.ChannelResponser
}

// ResponseChannelTheHive структура ответа от API TheHive
type ResponseChannelTheHive struct {
	StatusCode int    //статус кода ответа
	RequestId  string //UUID идентификатор ответа (соответствует идентификатору запроса)
	Err        error  //объект ошибки
	Data       []byte //набор данных
}

// RequestCommand структура с командами для обработки модулем
type RequestCommand struct {
	Service   string          `json:"service"`    //наименование сервиса
	Command   string          `json:"command"`    //команда
	RootId    string          `json:"root_id"`    //основной id, как правило это rootId case или alert
	CaseId    string          `json:"case_id"`    //id кейса
	Username  string          `json:"username"`   //имя пользователя, необходим если нужно указать пользователя выполнившего действие
	FieldName string          `json:"field_name"` //некое ключевое поле
	Value     string          `json:"value"`      //устанавливаемое значение
	ByteData  json.RawMessage `json:"byte_ data"` //набор данных в бинарном виде которые обрабатываются отдельно
}
