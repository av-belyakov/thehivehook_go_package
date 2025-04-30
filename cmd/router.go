package main

import (
	"context"

	"github.com/av-belyakov/thehivehook_go_package/cmd/webhookserver"
	"github.com/av-belyakov/thehivehook_go_package/internal/datamodels"
)

func router(
	ctx context.Context,
	fromWebHook <-chan webhookserver.ChanFromWebHookServer,
	//fromNatsAPI <-chan commoninterfaces.ChannelRequester,
	fromNatsAPI <-chan datamodels.RequestChan,
	//toTheHiveAPI chan<- commoninterfaces.ChannelRequester,
	toTheHiveAPI chan<- datamodels.RequestChan,
	//toNatsAPI chan<- commoninterfaces.ChannelRequester,
	toNatsAPI chan<- datamodels.RequestChan) {

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
