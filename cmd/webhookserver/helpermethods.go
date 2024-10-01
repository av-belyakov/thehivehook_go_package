package webhookserver

import "github.com/av-belyakov/thehivehook_go_package/cmd/commoninterfaces"

//********************* Response ********************

func NewChannelRespons() *ResponsTheHive {
	return &ResponsTheHive{}
}

func (r *ResponsTheHive) GetStatusCode() int {
	return r.StatusCode
}

func (r *ResponsTheHive) SetStatusCode(v int) {
	r.StatusCode = v
}

func (r *ResponsTheHive) GetRequestId() string {
	return r.RequestId
}

func (r *ResponsTheHive) SetRequestId(v string) {
	r.RequestId = v
}

func (r *ResponsTheHive) GetData() []byte {
	return r.Data
}

func (r *ResponsTheHive) SetData(v []byte) {
	r.Data = v
}

//******************* Request *********************

func NewChannelRequest() *RequestTheHive {
	return &RequestTheHive{}
}

func (r *RequestTheHive) GetRequestId() string {
	return r.RequestId
}

func (r *RequestTheHive) SetRequestId(v string) {
	r.RequestId = v
}

func (r *RequestTheHive) GetRootId() string {
	return r.RootId
}

func (r *RequestTheHive) SetRootId(v string) {
	r.RootId = v
}

func (r *RequestTheHive) GetCommand() string {
	return r.Command
}

func (r *RequestTheHive) SetCommand(v string) {
	r.Command = v
}

func (r *RequestTheHive) GetChanOutput() chan commoninterfaces.ChannelResponser {
	return r.ChanOutput
}

func (r *RequestTheHive) SetChanOutput(v chan commoninterfaces.ChannelResponser) {
	r.ChanOutput = v
}
