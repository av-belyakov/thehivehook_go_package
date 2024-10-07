package webhookserver

import "github.com/av-belyakov/thehivehook_go_package/cmd/commoninterfaces"

//********************* Response ********************

// NewChannelRespons конструктор ответа от TheHive
func NewChannelRespons() *ResponsTheHive {
	return &ResponsTheHive{}
}

// GetStatusCode метод возвращает статус кода ответа
func (r *ResponsTheHive) GetStatusCode() int {
	return r.StatusCode
}

// SetStatusCode метод устанавливает статус кода ответа
func (r *ResponsTheHive) SetStatusCode(v int) {
	r.StatusCode = v
}

// GetRequestId метод возвращает уникальный идентификатор запроса
func (r *ResponsTheHive) GetRequestId() string {
	return r.RequestId
}

// SetRequestId метод устанавливает уникальный идентификатор запроса
func (r *ResponsTheHive) SetRequestId(v string) {
	r.RequestId = v
}

// GetData метод возвращает данные
func (r *ResponsTheHive) GetData() []byte {
	return r.Data
}

// SetData метод устанавливает определенные данные
func (r *ResponsTheHive) SetData(v []byte) {
	r.Data = v
}

//******************* Request *********************

// NewChannelRequest конструктор формирующий структуру для выполнения запросов к модулю apithehive
func NewChannelRequest() *RequestTheHive {
	return &RequestTheHive{}
}

// GetRequestId метод возвращает уникальный идентификатор запроса
func (r *RequestTheHive) GetRequestId() string {
	return r.RequestId
}

// SetRequestId метод устанавливает уникальный идентификатор запроса
func (r *RequestTheHive) SetRequestId(v string) {
	r.RequestId = v
}

// GetRootId метод возвращает основной идентификатор кейса или алерта
func (r *RequestTheHive) GetRootId() string {
	return r.RootId
}

// SetRootId метод устанавливает основной идентификатор кейса или алерта
func (r *RequestTheHive) SetRootId(v string) {
	r.RootId = v
}

// GetCommand метод возвращает команду, на основе которой выполняются определенные действия
func (r *RequestTheHive) GetCommand() string {
	return r.Command
}

// SetCommand метод устанавливает, на основе которой выполняются определенные действия
func (r *RequestTheHive) SetCommand(v string) {
	r.Command = v
}

// GetChanOutput метод возвращает канал через который ответ от модуля apithehive передается
// источнику запроса
func (r *RequestTheHive) GetChanOutput() chan commoninterfaces.ChannelResponser {
	return r.ChanOutput
}

// SetChanOutput метод устанавливает канал через который ответ от модуля apithehive передается
// источнику запроса
func (r *RequestTheHive) SetChanOutput(v chan commoninterfaces.ChannelResponser) {
	r.ChanOutput = v
}
