package webhookserver

import (
	"encoding/json"
	"sync"

	"github.com/av-belyakov/thehivehook_go_package/cmd/thehiveapi"
)

func CreateEvenCase(uuidStorage, rootId string, chanTheHiveAPI chan<- thehiveapi.RequestChannelTheHive) (ReadyMadeEventCase, error) {
	var (
		wg   sync.WaitGroup
		rmec ReadyMadeEventCase = ReadyMadeEventCase{}

		mainErr error

		chanResObservable chan thehiveapi.ResponseChannelTheHive = make(chan thehiveapi.ResponseChannelTheHive)
		chanResTTL        chan thehiveapi.ResponseChannelTheHive = make(chan thehiveapi.ResponseChannelTheHive)
	)

	wg.Add(2)

	go func() {
		for res := range chanResObservable {
			msg := []interface{}{}
			if err := json.Unmarshal(res.Data, &msg); err != nil {
				mainErr = err

				return
			}

			rmec.Observables = msg
		}

		wg.Done()
	}()
	go func() {
		for res := range chanResTTL {
			msg := []interface{}{}
			if err := json.Unmarshal(res.Data, &msg); err != nil {
				mainErr = err

				return
			}

			rmec.TTPs = msg
		}

		wg.Done()
	}()

	chanTheHiveAPI <- thehiveapi.RequestChannelTheHive{
		RequestId:  uuidStorage,
		RootId:     rootId,
		Command:    "get_observables",
		ChanOutput: chanResObservable,
	}

	chanTheHiveAPI <- thehiveapi.RequestChannelTheHive{
		RequestId:  uuidStorage,
		RootId:     rootId,
		Command:    "get_ttp",
		ChanOutput: chanResTTL,
	}

	wg.Wait()

	return rmec, mainErr
}
