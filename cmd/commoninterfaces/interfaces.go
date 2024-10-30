package commoninterfaces

type ChannelResponser interface {
	RequestIdHandler
	DataHandler
	GetStatusCode() int
	SetStatusCode(int)
}

type ChannelRequester interface {
	RequestIdHandler
	RootIdHandler
	CommandHandler
	OrderHandler
	DataHandler
	GetChanOutput() chan ChannelResponser
	SetChanOutput(chan ChannelResponser)
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

type CommandHandler interface {
	GetCommand() string
	SetCommand(string)
}

type DataHandler interface {
	GetData() []byte
	SetData([]byte)
}

type Logger interface {
	GetChan() <-chan Messager
	Send(msgType, msgData string)
}

type Messager interface {
	GetType() string
	GetMessage() string
	SetType(v string)
	SetMessage(v string)
}

type WriterLoggingData interface {
	WriteLoggingData(str, typeLogFile string) bool
}
