package main

import (
	"context"
	"fmt"
	"runtime"
	"time"

	"github.com/av-belyakov/simplelogger"
	"github.com/av-belyakov/thehivehook_go_package/cmd/zabbixapi"
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

// interactionZabbix осуществляет взаимодействие с Zabbix
func interactionZabbix(
	ctx context.Context,
	confApp *confighandler.ConfigApp,
	sl simplelogger.SimpleLoggerSettings,
	channelZabbix <-chan zabbixapi.MessageSettings) error {

	connTimeout := time.Duration(7 * time.Second)
	hz, err := zabbixapi.NewZabbixConnection(
		ctx,
		zabbixapi.SettingsZabbixConnection{
			Port:              confApp.Zabbix.NetworkPort,
			Host:              confApp.Zabbix.NetworkHost,
			NetProto:          "tcp",
			ZabbixHost:        confApp.Zabbix.ZabbixHost,
			ConnectionTimeout: &connTimeout,
		})
	if err != nil {
		return err
	}

	et := make([]zabbixapi.EventType, len(confApp.Zabbix.EventTypes))
	for _, v := range confApp.Zabbix.EventTypes {
		et = append(et, zabbixapi.EventType{
			IsTransmit: v.IsTransmit,
			EventType:  v.EventType,
			ZabbixKey:  v.ZabbixKey,
			Handshake:  zabbixapi.Handshake(v.Handshake),
		})
	}

	if err = hz.Handler(et, channelZabbix); err != nil {
		return err
	}

	go func() {
		for err := range hz.GetChanErr() {
			_, f, l, _ := runtime.Caller(0)
			_ = sl.WriteLoggingData(fmt.Sprintf("zabbix module: '%s' %s:%d", err.Error(), f, l-1), "error")
		}
	}()

	return nil
}
