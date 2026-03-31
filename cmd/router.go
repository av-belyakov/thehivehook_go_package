package main

import (
	"context"

	"github.com/av-belyakov/thehivehook_go_package/cmd/webhookserver"
	"github.com/av-belyakov/thehivehook_go_package/internal/interfaces"
)

func router(
	ctx context.Context,
	fromWebHook <-chan webhookserver.ChanFromWebHookServer,
	fromNatsAPI <-chan interfaces.ChannelRequester,
	toTheHiveAPI chan<- interfaces.ChannelRequester,
	toNatsAPI chan<- interfaces.ChannelRequester) {

	go func() {
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
	}()
}
