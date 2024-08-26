package main

import (
	"github.com/av-belyakov/simplelogger"
	"github.com/av-belyakov/thehivehook_go_package/internal/confighandler"
)

func getLoggerSettings(cls []confighandler.LogSet) []simplelogger.MessageTypeSettings {
	loggerConf := make([]simplelogger.MessageTypeSettings, 0, len(cls))

	for _, v := range cls {
		loggerConf = append(loggerConf, simplelogger.MessageTypeSettings{
			MsgTypeName:   v.MsgTypeName,
			WritingFile:   v.WritingFile,
			PathDirectory: v.PathDirectory,
			WritingStdout: v.WritingStdout,
			MaxFileSize:   v.MaxFileSize,
		})
	}

	return loggerConf
}
