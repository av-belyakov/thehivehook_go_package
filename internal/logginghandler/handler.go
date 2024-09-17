package logginghandler

import (
	"context"
	"fmt"

	"github.com/av-belyakov/simplelogger"
	"github.com/av-belyakov/thehivehook_go_package/cmd/zabbixapi"
)

// LoggingHandler обработчик и распределитель логов
func LoggingHandler(
	ctx context.Context,
	channelZabbix chan<- zabbixapi.MessageSettings,
	sl simplelogger.SimpleLoggerSettings,
	logging <-chan MessageLogging) {

	for {
		select {
		case <-ctx.Done():
			return

		case msg := <-logging:
			//**********************************************************************
			//здесь так же может быть вывод в консоль, с индикацией цветов соответствующих
			//определенному типу сообщений но для этого надо включить вывод на stdout
			//в конфигурационном фале
			_ = sl.WriteLoggingData(msg.MsgData, msg.MsgType)

			if msg.MsgType == "error" || msg.MsgType == "warning" {
				channelZabbix <- zabbixapi.MessageSettings{
					EventType: "error",
					Message:   fmt.Sprintf("%s: %s", msg.MsgType, msg.MsgData),
				}
			}

			if msg.MsgType == "info" {
				channelZabbix <- zabbixapi.MessageSettings{
					EventType: "info",
					Message:   msg.MsgData,
				}
			}
		}
	}
}
