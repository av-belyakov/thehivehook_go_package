package webhookserver

import (
	"fmt"

	"github.com/av-belyakov/thehivehook_go_package/cmd/commoninterfaces"
)

//********************* Response ********************

// NewChannelRespons конструктор ответа
func NewChannelRespons() *ResponsToWebHook {
	return &ResponsToWebHook{}
}

// GetStatusCode метод возвращает статус кода ответа
func (r *ResponsToWebHook) GetStatusCode() int {
	return r.StatusCode
}

// SetStatusCode метод устанавливает статус кода ответа
func (r *ResponsToWebHook) SetStatusCode(v int) {
	r.StatusCode = v
}

// GetRequestId метод возвращает уникальный идентификатор запроса
func (r *ResponsToWebHook) GetRequestId() string {
	return r.RequestId
}

// SetRequestId метод устанавливает уникальный идентификатор запроса
func (r *ResponsToWebHook) SetRequestId(v string) {
	r.RequestId = v
}

// GetData метод возвращает данные
func (r *ResponsToWebHook) GetData() []byte {
	return r.Data
}

// SetData метод устанавливает определенные данные
func (r *ResponsToWebHook) SetData(v []byte) {
	r.Data = v
}

//******************* Request *********************

// NewChannelRequest конструктор формирующий структуру для выполнения запросов к модулю apithehive
func NewChannelRequest() *RequestFromWebHook {
	return &RequestFromWebHook{}
}

// GetRequestId метод возвращает уникальный идентификатор запроса
func (r *RequestFromWebHook) GetRequestId() string {
	return r.RequestId
}

// SetRequestId метод устанавливает уникальный идентификатор запроса
func (r *RequestFromWebHook) SetRequestId(v string) {
	r.RequestId = v
}

// GetRootId метод возвращает основной идентификатор кейса или алерта
func (r *RequestFromWebHook) GetRootId() string {
	return r.RootId
}

// SetRootId метод устанавливает основной идентификатор кейса или алерта
func (r *RequestFromWebHook) SetRootId(v string) {
	r.RootId = v
}

// GetCommand метод возвращает команду, на основе которой выполняются определенные действия
func (r *RequestFromWebHook) GetCommand() string {
	return r.Command
}

// SetCommand метод устанавливает, на основе которой выполняются определенные действия
func (r *RequestFromWebHook) SetCommand(v string) {
	r.Command = v
}

// GetData метод возвращает некий набор данных
func (r *RequestFromWebHook) GetData() []byte {
	return r.Data
}

// SetData метод устанавливает некий набор данных
func (r *RequestFromWebHook) SetData(v []byte) {
	r.Data = v
}

// GetChanOutput метод возвращает канал через который ответ от модуля apithehive передается
// источнику запроса
func (r *RequestFromWebHook) GetChanOutput() chan commoninterfaces.ChannelResponser {
	return r.ChanOutput
}

// SetChanOutput метод устанавливает канал через который ответ от модуля apithehive передается
// источнику запроса
func (r *RequestFromWebHook) SetChanOutput(v chan commoninterfaces.ChannelResponser) {
	r.ChanOutput = v
}

//**************************** вспомогательные методы ****************************

// GetEventId возвращает уникальный id элемента основанный на комбинации некоторых значений EventElement
func (e EventElement) GetEventId() string {
	return fmt.Sprintf("%s:%d:%s", e.ObjectType, e.Object.CreatedAt, e.RootId)
}
