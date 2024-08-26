package datamodels

import (
	"log/slog"
	"sync"
)

type Config struct {
	baseURL  string
	httpPort int
}

type Application struct {
	config Config
	logger *slog.Logger
	wg     sync.WaitGroup
}
