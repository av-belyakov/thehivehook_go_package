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

// GetError метод возвращает объект ошибки
func (r *ResponseChannelTheHive) GetError() error {
	return r.Err
}

// SetError метод устанавливает объект ошибки
func (r *ResponseChannelTheHive) SetError(e error) {
	r.Err = e
}

// GetData метод возвращает данные
func (r *ResponseChannelTheHive) GetData() []byte {
	return r.Data
}

// SetData метод устанавливает определенные данные
func (r *ResponseChannelTheHive) SetData(v []byte) {
	r.Data = v
}

// SendToChan отправляет ответ через полученый канал соответвтующий интерфейсу ChannelResponser
func (r *ResponseChannelTheHive) SendToChan(ch chan<- commoninterfaces.ChannelResponser) {
	if ch != nil {
		ch <- r
	}
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

// GetCaseId метод возвращает идентификатор кейса
func (r *RequestChannelTheHive) GetCaseId() string {
	return r.CaseId
}

// SetCaseId метод устанавливает идентификатор кейса
func (r *RequestChannelTheHive) SetCaseId(v string) {
	r.CaseId = v
}

// GetCommand метод возвращает команду, на основе которой выполняются определенные действия
func (r *RequestChannelTheHive) GetCommand() string {
	return r.Command
}

// SetCommand метод устанавливает, на основе которой выполняются определенные действия
func (r *RequestChannelTheHive) SetCommand(v string) {
	r.Command = v
}

// GetData метод возвращает данные
func (r *RequestChannelTheHive) GetData() interface{} {
	return r.Data
}

// SetData метод устанавливает определенные данные
func (r *RequestChannelTheHive) SetData(i interface{}) {
	r.Data = i
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

//******************* Различные вспомогательные методы *********************

// NewLogWrite создаёт вспомогательный тип для логирования
func NewLogWrite(l commoninterfaces.Logger) *LogWrite {
	return &LogWrite{logger: l}
}

// Write логирование данных
func (lw *LogWrite) Write(msgType, msg string) bool {
	lw.logger.Send(msgType, msg)

	return true
}

// NewSpecialObjectForCache конструктор вспомогательного типа реализующий интерфейс CacheStorageFuncHandler[T any]
func NewSpecialObjectForCache[T any]() *SpecialObjectForCache[T] {
	return &SpecialObjectForCache[T]{}
}

func (o *SpecialObjectForCache[T]) SetID(v string) {
	o.id = v
}

func (o *SpecialObjectForCache[T]) GetID() string {
	return o.id
}

func (o *SpecialObjectForCache[T]) SetObject(v T) {
	o.object = v
}

func (o *SpecialObjectForCache[T]) GetObject() T {
	return o.object
}

func (o *SpecialObjectForCache[T]) SetFunc(f func(int) bool) {
	o.handlerFunc = f
}

func (o *SpecialObjectForCache[T]) GetFunc() func(int) bool {
	return o.handlerFunc
}

// Comparison сравнение содержимого объектов
// В данном случае сравнение нет, это простая заглушка.
// Для того что бы не досить thehive метод всегда будет возвращать TRUE.
// Соответственно не будет заменять объект в работе.
func (o *SpecialObjectForCache[T]) Comparison(objFromCache T) bool {
	return true
}

// MatchingAndReplacement сопоставление элементов объекта и замена этих значений
// в объекте который уже находится в кеше
// В данном случае простая заглушка, так как метод Comparison тоже является заглушкой.
// Этот метод даже не будет вызыватся, потому метод Comparison возвращает всегда true.
func (o *SpecialObjectForCache[T]) MatchingAndReplacement(objFromCache T) T {
	return objFromCache
}
