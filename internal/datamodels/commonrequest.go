package datamodels

import "context"

type RequestChan struct {
	Data        any
	Context     context.Context
	Order       string
	RootId      string
	Command     string
	RequestId   string
	ElementType string
	CaseId      int
	ChOutput    chan ResponseChan
}

type ResponseChan struct {
	Data       []byte
	Error      error
	RequestId  string
	StatusCode int
}
