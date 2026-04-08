package thehiveapi

import (
	"github.com/av-belyakov/thehivehook_go_package/internal/interfaces"
)

// SpecialObjectForCache является вспомогательным типом который реализует интерфейс
// CacheStorageFuncHandler[T any] где в методе Comparison(objFromCache T) bool необходимо
// реализовать подробное сравнение объекта типа T.
// Нужен для пакета cachingstoragewithqueue
type SpecialObjectForCache[T any] struct {
	object      T
	handlerFunc func(int) bool
	id          string
}

// LogWrite вспомогательный тип применяемый для логирования
type LogWrite struct {
	logger interfaces.Logger
}
