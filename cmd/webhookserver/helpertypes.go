package webhookserver

type ResponsTheHive struct {
	StatusCode int
	RequestId  string
	Data       []byte
}

type RequestTheHive struct {
	RequestId  string
	RootId     string
	Command    string
	ChanOutput chan ResponsTheHive
}
