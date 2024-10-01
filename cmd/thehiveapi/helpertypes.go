package thehiveapi

import "github.com/av-belyakov/thehivehook_go_package/cmd/commoninterfaces"

type RequestChannelTheHive struct {
	RequestId  string
	RootId     string
	Command    string
	ChanOutput chan commoninterfaces.ChannelResponser
}

type ResponseChannelTheHive struct {
	StatusCode int
	RequestId  string
	Data       []byte
}
