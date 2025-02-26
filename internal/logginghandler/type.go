package logginghandler

import "github.com/av-belyakov/thehivehook_go_package/cmd/commoninterfaces"

type LoggingChan struct {
	dataWriter commoninterfaces.WriterLoggingData
	//запись в систему логирования
	chanSystemMonitoring chan<- commoninterfaces.Messager
	//канал отправки в систему мониторинга, например, Zabbix
	chanLogging chan commoninterfaces.Messager
	//основной канал приёма логов
}

// MessageLogging содержит информацию используемую при логировании
type MessageLogging struct {
	Message, Type string
}
