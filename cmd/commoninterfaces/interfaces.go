package commoninterfaces

//************ каналы *************

type ChannelResponser interface {
	RequestIdHandler
	GetStatusCode() int
	SetStatusCode(int)
	GetError() error
	SetError(error)
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
	GetData() interface{}
	SetData(interface{})
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
	WriteLoggingData(str, typeLogFile string) bool
}

//************** кэширование функций *****************

type CacheFuncRunner interface {
	SetMethod(id string, f func(int) bool) string
	GetMethod(id string) (func(int) bool, bool)
	DeleteElement(id string)
}
