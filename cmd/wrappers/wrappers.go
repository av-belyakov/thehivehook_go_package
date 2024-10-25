package wrappers

import (
	"context"
	"fmt"
	"runtime"
	"time"

	"github.com/av-belyakov/thehivehook_go_package/cmd/commoninterfaces"
	"github.com/av-belyakov/thehivehook_go_package/cmd/zabbixapi"
)

// WrappersZabbixInteraction обертка для взаимодействия с модулем zabbixapi
func WrappersZabbixInteraction(
	ctx context.Context,
	writerLoggingData commoninterfaces.WriterLoggingData,
	settings WrappersZabbixInteractionSettings,
	channelZabbix <-chan commoninterfaces.Messager) error {

	connTimeout := time.Duration(7 * time.Second)
	hz, err := zabbixapi.New(zabbixapi.SettingsZabbixConnection{
		Port:              settings.NetworkPort,
		Host:              settings.NetworkHost,
		NetProto:          "tcp",
		ZabbixHost:        settings.ZabbixHost,
		ConnectionTimeout: &connTimeout,
	})
	if err != nil {
		return err
	}

	et := make([]zabbixapi.EventType, len(settings.EventTypes))
	for _, v := range settings.EventTypes {
		et = append(et, zabbixapi.EventType{
			IsTransmit: v.IsTransmit,
			EventType:  v.EventType,
			ZabbixKey:  v.ZabbixKey,
			Handshake:  zabbixapi.Handshake(v.Handshake),
		})
	}

	recipient := make(chan zabbixapi.Messager)
	if err = hz.Start(ctx, et, recipient); err != nil {
		return err
	}

	go func() {
		for {
			select {
			case <-ctx.Done():
				return

			case msg := <-channelZabbix:
				newMessageSettings := &zabbixapi.MessageSettings{}
				newMessageSettings.SetType(msg.GetType())
				newMessageSettings.SetMessage(msg.GetMessage())

				recipient <- newMessageSettings
			}
		}
	}()

	go func() {
		for err := range hz.GetChanErr() {
			_, f, l, _ := runtime.Caller(0)
			writerLoggingData.WriteLoggingData(fmt.Sprintf("zabbix module: '%s' %s:%d", err.Error(), f, l-1), "error")
		}
	}()

	return nil
}
