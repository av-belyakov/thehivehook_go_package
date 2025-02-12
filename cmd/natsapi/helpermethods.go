package natsapi

import (
	"context"

	"github.com/av-belyakov/thehivehook_go_package/cmd/commoninterfaces"
)

//********************* Response ********************

// NewChannelRespons конструктор ответа
func NewChannelRespons() *ResponsToNats {
	return &ResponsToNats{}
}

// GetStatusCode метод возвращает статус кода ответа
func (r *ResponsToNats) GetStatusCode() int {
	return r.StatusCode
}

// SetStatusCode метод устанавливает статус кода ответа
func (r *ResponsToNats) SetStatusCode(v int) {
	r.StatusCode = v
}

// GetRequestId метод возвращает уникальный идентификатор запроса
func (r *ResponsToNats) GetRequestId() string {
	return r.RequestId
}

// SetRequestId метод устанавливает уникальный идентификатор запроса
func (r *ResponsToNats) SetRequestId(v string) {
	r.RequestId = v
}

// GetData метод возвращает данные
func (r *ResponsToNats) GetData() []byte {
	return r.Data
}

// SetData метод устанавливает определенные данные
func (r *ResponsToNats) SetData(v []byte) {
	r.Data = v
}

//******************* Request *********************

// NewChannelRequest конструктор формирующий структуру для выполнения запросов к модулю apithehive
func NewChannelRequest() *RequestFromNats {
	return &RequestFromNats{}
}

// GetContext возвращает контекст
func (r *RequestFromNats) GetContext() context.Context {
	return r.ctx
}

// SetContext устанавливает контекст
func (r *RequestFromNats) SetContext(v context.Context) {
	r.ctx = v
}

// GetRequestId метод возвращает уникальный идентификатор запроса
func (r *RequestFromNats) GetRequestId() string {
	return r.RequestId
}

// SetRequestId метод устанавливает уникальный идентификатор запроса
func (r *RequestFromNats) SetRequestId(v string) {
	r.RequestId = v
}

// GetElementType метод возвращает тип элемента
func (r *RequestFromNats) GetElementType() string {
	return r.ElementType
}

// SetElementType метод устанавливает тип элемента
func (r *RequestFromNats) SetElementType(v string) {
	r.ElementType = v
}

// GetRootId метод возвращает основной идентификатор кейса или алерта
func (r *RequestFromNats) GetRootId() string {
	return r.RootId
}

// SetRootId метод устанавливает основной идентификатор кейса или алерта
func (r *RequestFromNats) SetRootId(v string) {
	r.RootId = v
}

// GetCaseId метод возвращает идентификатор кейса
func (r *RequestFromNats) GetCaseId() string {
	return r.CaseId
}

// SetCaseId метод устанавливает идентификатор кейса
func (r *RequestFromNats) SetCaseId(v string) {
	r.CaseId = v
}

// GetCommand метод возвращает команду, на основе которой выполняются определенные действия
func (r *RequestFromNats) GetCommand() string {
	return r.Command
}

// SetCommand метод устанавливает команду, на основе которой выполняются определенные действия
func (r *RequestFromNats) SetCommand(v string) {
	r.Command = v
}

// GetOrder метод возвращает распоряжение
func (r *RequestFromNats) GetOrder() string {
	return r.Order
}

// SetOrder метод устанавливает распоряжение
func (r *RequestFromNats) SetOrder(v string) {
	r.Order = v
}

// GetData метод возвращает некий набор данных
func (r *RequestFromNats) GetData() interface{} {
	return r.Data
}

// SetData метод устанавливает некий набор данных
func (r *RequestFromNats) SetData(i interface{}) {
	r.Data = i
}

// GetChanOutput метод возвращает канал через который ответ от модуля apithehive передается
// источнику запроса
func (r *RequestFromNats) GetChanOutput() chan commoninterfaces.ChannelResponser {
	return r.ChanOutput
}

// SetChanOutput метод устанавливает канал через который ответ от модуля apithehive передается
// источнику запроса
func (r *RequestFromNats) SetChanOutput(v chan commoninterfaces.ChannelResponser) {
	r.ChanOutput = v
}
