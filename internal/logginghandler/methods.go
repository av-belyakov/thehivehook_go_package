package logginghandler

func New() *LoggingChan {
	return &LoggingChan{
		logChan: make(chan MessageLogging),
	}
}

func (l *LoggingChan) GetChan() <-chan MessageLogging {
	return l.logChan
}

func (l *LoggingChan) Send(msgType, msgData string) {
	l.logChan <- MessageLogging{
		MsgType: msgType,
		MsgData: msgData,
	}
}

func (l *LoggingChan) Close() {
	close(l.logChan)
}
