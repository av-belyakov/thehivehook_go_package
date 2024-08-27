package zabbixapi

import (
	"context"
	"fmt"
	"runtime"
	"time"

	"github.com/av-belyakov/simplelogger"
	"github.com/av-belyakov/thehivehook_go_package/internal/confighandler"
)

// InteractionZabbix осуществляет взаимодействие с Zabbix
func InteractionZabbix(
	ctx context.Context,
	confApp confighandler.ConfigApp,
	sl simplelogger.SimpleLoggerSettings,
	channelZabbix <-chan MessageSettings) error {

	connTimeout := time.Duration(7 * time.Second)
	hz, err := NewZabbixConnection(
		ctx,
		SettingsZabbixConnection{
			Port:              confApp.Zabbix.NetworkPort,
			Host:              confApp.Zabbix.NetworkHost,
			NetProto:          "tcp",
			ZabbixHost:        confApp.Zabbix.ZabbixHost,
			ConnectionTimeout: &connTimeout,
		})
	if err != nil {
		return err
	}

	et := make([]EventType, len(confApp.Zabbix.EventTypes))
	for _, v := range confApp.Zabbix.EventTypes {
		et = append(et, EventType{
			IsTransmit: v.IsTransmit,
			EventType:  v.EventType,
			ZabbixKey:  v.ZabbixKey,
			Handshake:  Handshake(v.Handshake),
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
