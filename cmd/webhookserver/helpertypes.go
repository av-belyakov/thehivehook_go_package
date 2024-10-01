package webhookserver

import "github.com/av-belyakov/thehivehook_go_package/cmd/commoninterfaces"

type ResponsTheHive struct {
	StatusCode int
	RequestId  string
	Data       []byte
}

type RequestTheHive struct {
	RequestId  string
	RootId     string
	Command    string
	ChanOutput chan commoninterfaces.ChannelResponser
}
