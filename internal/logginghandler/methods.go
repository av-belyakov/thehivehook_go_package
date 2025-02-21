package logginghandler

import "github.com/av-belyakov/thehivehook_go_package/cmd/commoninterfaces"

func (l *LoggingChan) GetChan() <-chan commoninterfaces.Messager {
	return l.chanLogging
}

func (l *LoggingChan) Send(msgType, message string) {
	ms := NewMessageLogging()
	ms.SetType(msgType)
	ms.SetMessage(message)

	l.chanLogging <- ms
}

func (l *LoggingChan) Close() {
	close(l.chanLogging)
}

func NewMessageLogging() *MessageLogging {
	return &MessageLogging{}
}

// GetMessage возвращает сообщение
func (ml *MessageLogging) GetMessage() string {
	return ml.Message
}

// SetMessage устанавливает сообщение
func (ml *MessageLogging) SetMessage(v string) {
	ml.Message = v
}

// GetType возвращает тип сообщения
func (ml *MessageLogging) GetType() string {
	return ml.Type
}

// SetType устанавливает тип сообщения
func (ml *MessageLogging) SetType(v string) {
	ml.Type = v
}
