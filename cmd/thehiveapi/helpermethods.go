package thehiveapi

import "github.com/av-belyakov/thehivehook_go_package/cmd/commoninterfaces"

//********************* Response ********************

func NewChannelRespons() *ResponseChannelTheHive {
	return &ResponseChannelTheHive{}
}

func (r *ResponseChannelTheHive) GetStatusCode() int {
	return r.StatusCode
}

func (r *ResponseChannelTheHive) SetStatusCode(v int) {
	r.StatusCode = v
}

func (r *ResponseChannelTheHive) GetRequestId() string {
	return r.RequestId
}

func (r *ResponseChannelTheHive) SetRequestId(v string) {
	r.RequestId = v
}

func (r *ResponseChannelTheHive) GetData() []byte {
	return r.Data
}

func (r *ResponseChannelTheHive) SetData(v []byte) {
	r.Data = v
}

//******************* Request *********************

func NewChannelRequest() *RequestChannelTheHive {
	return &RequestChannelTheHive{}
}

func (r *RequestChannelTheHive) GetRequestId() string {
	return r.RequestId
}

func (r *RequestChannelTheHive) SetRequestId(v string) {
	r.RequestId = v
}

func (r *RequestChannelTheHive) GetRootId() string {
	return r.RootId
}

func (r *RequestChannelTheHive) SetRootId(v string) {
	r.RootId = v
}

func (r *RequestChannelTheHive) GetCommand() string {
	return r.Command
}

func (r *RequestChannelTheHive) SetCommand(v string) {
	r.Command = v
}

func (r *RequestChannelTheHive) GetChanOutput() chan commoninterfaces.ChannelResponser {
	return r.ChanOutput
}

func (r *RequestChannelTheHive) SetChanOutput(v chan commoninterfaces.ChannelResponser) {
	r.ChanOutput = v
}
