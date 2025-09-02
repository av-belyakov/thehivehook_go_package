package natsapi

import (
	"context"

	"github.com/av-belyakov/thehivehook_go_package/cmd/commoninterfaces"
)

// RequestFromNats структура запроса из модуля
type RequestFromNats struct {
	Data        any    //набор данных
	RequestId   string //id запроса
	ElementType string //тип элемента
	RootId      string //идентификатор по которому в TheHive будет выполнятся поиск
	CaseId      string //идентификатор кейса в TheHive
	Command     string //команда
	Order       string //распоряжение
	ctx         context.Context
	ChanOutput  chan commoninterfaces.ChannelResponser //канал ответа реализующий интерфейс commoninterfaces.ChannelResponser
}

// ResponsToNats структура ответа в модуля
type ResponsToNats struct {
	RequestId  string //UUID идентификатор ответа (соответствует идентификатору запроса)
	Data       []byte //набор данных
	StatusCode int    //статус кода ответа
}

// RequestCommand структура с командами для обработки модулем
type RequestCommand struct {
	Service string `json:"service"` //наименование сервиса
	Command string `json:"command"` //команда
	RootId  string `json:"root_id"` //основной id, как правило это rootId case или alert
	CaseId  string `json:"case_id"` //id кейса
	//Username  string          `json:"username"`   //имя пользователя, необходим если нужно указать пользователя выполнившего действие
	//FieldName string          `json:"field_name"` //некое ключевое поле
	//Value     string          `json:"value"`      //устанавливаемое значение
	//ByteData  json.RawMessage `json:"byte_ data"` //набор данных в бинарном виде которые обрабатываются отдельно
}
