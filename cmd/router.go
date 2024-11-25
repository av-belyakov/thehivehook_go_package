package main

import (
	"context"

	"github.com/av-belyakov/thehivehook_go_package/cmd/commoninterfaces"
	"github.com/av-belyakov/thehivehook_go_package/cmd/webhookserver"
)

func router(
	ctx context.Context,
	fromWebHook <-chan webhookserver.ChanFromWebHookServer,
	fromNatsAPI <-chan commoninterfaces.ChannelRequester,
	toTheHiveAPI chan<- commoninterfaces.ChannelRequester,
	toNatsAPI chan<- commoninterfaces.ChannelRequester) {

	for {
		select {
		case <-ctx.Done():
			return

		case msg := <-fromWebHook:
			switch msg.ForSomebody {
			case "to thehive":
				toTheHiveAPI <- msg.Data

			case "to nats":
				toNatsAPI <- msg.Data
			}

		case msg := <-fromNatsAPI:
			toTheHiveAPI <- msg

		}
	}
}
