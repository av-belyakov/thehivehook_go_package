package logginghandler

type LoggingChan struct {
	logChan chan MessageLogging
}

// MessageLogging содержит информацию используемую при логировании
// MsgData - сообщение
// MsgType - тип сообщения
type MessageLogging struct {
	MsgData, MsgType string
}
