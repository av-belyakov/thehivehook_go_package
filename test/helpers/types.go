package helpers

import "github.com/av-belyakov/thehivehook_go_package/v2/internal/interfaces"

// LoggingForTest тестовая структура для логирования
type LoggingForTest struct {
	chMessage chan interfaces.Messager
}

// MessageForTest тестовая структура для описания сообщения
type MessageForTest struct {
	msgType, msgData string
}
