package helpers

import "github.com/av-belyakov/thehivehook_go_package/v2/internal/interfaces"

// NewLoggingForTest тестовый логгер
func NewLoggingForTest() *LoggingForTest {
	return &LoggingForTest{
		chMessage: make(chan interfaces.Messager),
	}
}
