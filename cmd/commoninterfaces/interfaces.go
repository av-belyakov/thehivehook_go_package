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
