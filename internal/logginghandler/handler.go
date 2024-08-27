package logginghandler

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/lmittmann/tint"

	"github.com/av-belyakov/simplelogger"
	"github.com/av-belyakov/thehivehook_go_package/cmd/zabbixapi"
)

// LoggingHandler обработчик и распределитель логов
func LoggingHandler(
	channelZabbix chan<- zabbixapi.MessageSettings,
	sl simplelogger.SimpleLoggerSettings,
	logging <-chan MessageLogging) {
	loggerColor := slog.New(tint.NewHandler(os.Stdout, &tint.Options{Level: slog.LevelDebug}))

	for msg := range logging {
		_ = sl.WriteLoggingData(msg.MsgData, msg.MsgType)

		switch msg.MsgType {
		case "error":
			loggerColor.Error(msg.MsgData)

		case "warning":
			loggerColor.Warn(msg.MsgData)

		case "info":
			loggerColor.Info(msg.MsgData)
		}

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
