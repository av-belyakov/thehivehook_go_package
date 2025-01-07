package main

import (
	"github.com/av-belyakov/thehivehook_go_package/cmd/commoninterfaces"
	"github.com/av-belyakov/thehivehook_go_package/cmd/webhookserver"
)

func router(
	fromWebHook <-chan webhookserver.ChanFormWebHookServer,
	toTheHiveAPI chan<- commoninterfaces.ChannelRequester) {

	for msg := range fromWebHook {
		switch msg.ForSomebody {
		case "for thehive":
			toTheHiveAPI <- msg.Data

		case "for nats":
		}
	}
}
