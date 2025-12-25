package commoninterfaces

import (
	"context"
)

//************ каналы *************

type ChannelResponser interface {
	RequestIdHandler
	GetStatusCode() int
	SetStatusCode(int)
	GetError() error
	SetError(error)
	GetSource() string
	SetSource(string)
	GetData() []byte
	SetData([]byte)
}

type ChannelRequester interface {
	RequestIdHandler
	CommandHandler
	ElementTypeHandler
	RootIdHandler
	CaseIdHandler
	OrderHandler
	GetData() any
	SetData(any)
	GetContext() context.Context
	SetContext(v context.Context)
	GetChanOutput() chan ChannelResponser
	SetChanOutput(chan ChannelResponser)
}

type CaseIdHandler interface {
	GetCaseId() string
	SetCaseId(string)
}

type RequestIdHandler interface {
	GetRequestId() string
	SetRequestId(string)
}

type RootIdHandler interface {
	GetRootId() string
	SetRootId(string)
}

type OrderHandler interface {
	GetOrder() string
	SetOrder(string)
}

type ElementTypeHandler interface {
	GetElementType() string
	SetElementType(string)
}

type CommandHandler interface {
	GetCommand() string
	SetCommand(string)
}

//************** параметры типа CustomField TheHive ****************

type ParametersCustomFieldKeeper interface {
	GetType() string
	GetValue() string
	GetUsername() string
}

//************** логирование ***************

type Logger interface {
	GetChan() <-chan Messager
	Send(msgType, msgData string)
}

type Messager interface {
	GetType() string
	SetType(v string)
	GetMessage() string
	SetMessage(v string)
}

type WriterLoggingData interface {
	Write(typeLogFile, str string) bool
}
