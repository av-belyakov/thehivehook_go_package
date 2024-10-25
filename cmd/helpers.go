package main

import (
	"github.com/av-belyakov/simplelogger"
	"github.com/av-belyakov/thehivehook_go_package/internal/confighandler"
)

func getLoggerSettings(cls []confighandler.LogSet) []simplelogger.Options {
	loggerConf := make([]simplelogger.Options, 0, len(cls))

	for _, v := range cls {
		loggerConf = append(loggerConf, simplelogger.Options{
			MsgTypeName:     v.MsgTypeName,
			WritingToFile:   v.WritingFile,
			PathDirectory:   v.PathDirectory,
			WritingToStdout: v.WritingStdout,
			MaxFileSize:     v.MaxFileSize,
		})
	}

	return loggerConf
}
