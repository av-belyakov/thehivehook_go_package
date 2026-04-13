package main

import (
	"github.com/av-belyakov/thehivehook_go_package/v2/internal/confighandler"
	"github.com/av-belyakov/thehivehook_go_package/v2/internal/interfaces"
)

type majorRouter struct {
	cfg    confighandler.ConfigApp
	logger interfaces.Logger
}
