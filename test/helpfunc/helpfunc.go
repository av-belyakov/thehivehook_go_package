package helperfunc

import (
	"github.com/av-belyakov/thehivehook_go_package/internal/interfaces"
	"github.com/av-belyakov/thehivehook_go_package/internal/logginghandler"
)

type LoggingForTest struct {
	ch chan interfaces.Messager
}

func NewLoggingForTest() *LoggingForTest {
	return &LoggingForTest{
		ch: make(chan interfaces.Messager),
	}
}

func (l *LoggingForTest) GetChan() <-chan interfaces.Messager {
	return l.ch
}

func (l *LoggingForTest) Send(msgType, message string) {
	ms := logginghandler.NewMessageLogging()
	ms.SetType(msgType)
	ms.SetMessage(message)

	l.ch <- ms
}
