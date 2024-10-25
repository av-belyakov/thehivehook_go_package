package commoninterfaces

type ChannelResponser interface {
	GetStatusCode() int
	SetStatusCode(int)
	GetRequestId() string
	SetRequestId(string)
	GetData() []byte
	SetData([]byte)
}

type ChannelRequester interface {
	GetRequestId() string
	SetRequestId(string)
	GetRootId() string
	SetRootId(string)
	GetCommand() string
	SetCommand(string)
	GetData() []byte
	SetData([]byte)
	GetChanOutput() chan ChannelResponser
	SetChanOutput(chan ChannelResponser)
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
