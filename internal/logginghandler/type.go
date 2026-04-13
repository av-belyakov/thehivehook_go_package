package logginghandler

import "github.com/av-belyakov/thehivehook_go_package/v2/internal/interfaces"

type LoggingChan struct {
	dataWriter interfaces.WriterLoggingData
	//запись в систему логирования
	chanSystemMonitoring chan<- interfaces.Messager
	//канал отправки в систему мониторинга, например, Zabbix
	chanLogging chan interfaces.Messager
	//основной канал приёма логов
}

// MessageLogging содержит информацию используемую при логировании
type MessageLogging struct {
	Message, Type string
}
