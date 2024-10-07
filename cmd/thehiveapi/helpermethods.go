package thehiveapi

import "github.com/av-belyakov/thehivehook_go_package/cmd/commoninterfaces"

//********************* Response ********************

// NewChannelRespons конструктор ответа от TheHive
func NewChannelRespons() *ResponseChannelTheHive {
	return &ResponseChannelTheHive{}
}

// GetStatusCode метод возвращает статус кода ответа
func (r *ResponseChannelTheHive) GetStatusCode() int {
	return r.StatusCode
}

// SetStatusCode метод устанавливает статус кода ответа
func (r *ResponseChannelTheHive) SetStatusCode(v int) {
	r.StatusCode = v
}

// GetRequestId метод возвращает уникальный идентификатор запроса
func (r *ResponseChannelTheHive) GetRequestId() string {
	return r.RequestId
}

// SetRequestId метод устанавливает уникальный идентификатор запроса
func (r *ResponseChannelTheHive) SetRequestId(v string) {
	r.RequestId = v
}

// GetData метод возвращает данные
func (r *ResponseChannelTheHive) GetData() []byte {
	return r.Data
}

// SetData метод устанавливает определенные данные
func (r *ResponseChannelTheHive) SetData(v []byte) {
	r.Data = v
}

//******************* Request *********************

// NewChannelRequest конструктор формирующий структуру для выполнения запросов к модулю apithehive
func NewChannelRequest() *RequestChannelTheHive {
	return &RequestChannelTheHive{}
}

// GetRequestId метод возвращает уникальный идентификатор запроса
func (r *RequestChannelTheHive) GetRequestId() string {
	return r.RequestId
}

// SetRequestId метод устанавливает уникальный идентификатор запроса
func (r *RequestChannelTheHive) SetRequestId(v string) {
	r.RequestId = v
}

// GetRootId метод возвращает основной идентификатор кейса или алерта
func (r *RequestChannelTheHive) GetRootId() string {
	return r.RootId
}

// SetRootId метод устанавливает основной идентификатор кейса или алерта
func (r *RequestChannelTheHive) SetRootId(v string) {
	r.RootId = v
}

// GetCommand метод возвращает команду, на основе которой выполняются определенные действия
func (r *RequestChannelTheHive) GetCommand() string {
	return r.Command
}

// SetCommand метод устанавливает, на основе которой выполняются определенные действия
func (r *RequestChannelTheHive) SetCommand(v string) {
	r.Command = v
}

// GetChanOutput метод возвращает канал через который ответ от модуля apithehive передается
// источнику запроса
func (r *RequestChannelTheHive) GetChanOutput() chan commoninterfaces.ChannelResponser {
	return r.ChanOutput
}

// SetChanOutput метод устанавливает канал через который ответ от модуля apithehive передается
// источнику запроса
func (r *RequestChannelTheHive) SetChanOutput(v chan commoninterfaces.ChannelResponser) {
	r.ChanOutput = v
}
