package webhookserver

import (
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

// GetElementType метод возвращает тип элемента
func (r *RequestFromWebHook) GetElementType() string {
	return r.ElementType
}

// SetElementType метод устанавливает тип элемента
func (r *RequestFromWebHook) SetElementType(v string) {
	r.ElementType = v
}

// GetRootId метод возвращает основной идентификатор кейса или алерта
func (r *RequestFromWebHook) GetRootId() string {
	return r.RootId
}

// SetRootId метод устанавливает основной идентификатор кейса или алерта
func (r *RequestFromWebHook) SetRootId(v string) {
	r.RootId = v
}

// GetCaseId метод возвращает идентификатор кейса
func (r *RequestFromWebHook) GetCaseId() string {
	return r.CaseId
}

// SetCaseId метод устанавливает идентификатор кейса
func (r *RequestFromWebHook) SetCaseId(v string) {
	r.CaseId = v
}

// GetCommand метод возвращает команду, на основе которой выполняются определенные действия
func (r *RequestFromWebHook) GetCommand() string {
	return r.Command
}

// SetCommand метод устанавливает команду, на основе которой выполняются определенные действия
func (r *RequestFromWebHook) SetCommand(v string) {
	r.Command = v
}

// GetOrder метод возвращает распоряжение
func (r *RequestFromWebHook) GetOrder() string {
	return r.Order
}

// SetOrder метод устанавливает распоряжение
func (r *RequestFromWebHook) SetOrder(v string) {
	r.Order = v
}

// GetData метод возвращает некий набор данных
func (r *RequestFromWebHook) GetData() interface{} {
	return r.Data
}

// SetData метод устанавливает некий набор данных
func (r *RequestFromWebHook) SetData(i interface{}) {
	r.Data = i
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
