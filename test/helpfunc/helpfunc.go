package helperfunc

import (
	"github.com/av-belyakov/thehivehook_go_package/cmd/commoninterfaces"
	"github.com/av-belyakov/thehivehook_go_package/internal/logginghandler"
)

type LoggingForTest struct {
	ch chan commoninterfaces.Messager
}

func NewLoggingForTest() *LoggingForTest {
	return &LoggingForTest{
		ch: make(chan commoninterfaces.Messager),
	}
}

func (l *LoggingForTest) GetChan() <-chan commoninterfaces.Messager {
	return l.ch
}

func (l *LoggingForTest) Send(msgType, message string) {
	ms := logginghandler.NewMessageLogging()
	ms.SetType(msgType)
	ms.SetMessage(message)

	l.ch <- ms
}
