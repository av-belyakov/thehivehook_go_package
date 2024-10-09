package webhookserver

import (
	"encoding/json"
	"sync"

	"github.com/av-belyakov/thehivehook_go_package/cmd/commoninterfaces"
)

// CreateEvenCase генератор кейса содержащего в себе дополнительную информацию, такую как
// перечень значений observables и ttp
func CreateEvenCase(uuidStorage, rootId string, chanInput chan<- ChanFormWebHookServer) (ReadyMadeEventCase, error) {
	var (
		wg   sync.WaitGroup
		rmec ReadyMadeEventCase = ReadyMadeEventCase{}

		mainErr error

		chanResObservable chan commoninterfaces.ChannelResponser = make(chan commoninterfaces.ChannelResponser)
		chanResTTL        chan commoninterfaces.ChannelResponser = make(chan commoninterfaces.ChannelResponser)
	)

	wg.Add(2)

	go func() {
		for res := range chanResObservable {
			msg := []interface{}{}
			if err := json.Unmarshal(res.GetData(), &msg); err != nil {
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
			if err := json.Unmarshal(res.GetData(), &msg); err != nil {
				mainErr = err

				return
			}

			rmec.TTPs = msg
		}

		wg.Done()
	}()

	//запрос на поиск дополнительной информации об Observables
	reqObservable := NewChannelRequest()
	reqObservable.SetRequestId(uuidStorage)
	reqObservable.SetRootId(rootId)
	reqObservable.SetCommand("get_observables")
	reqObservable.SetChanOutput(chanResObservable)
	chanInput <- ChanFormWebHookServer{
		ForSomebody: "for thehive",
		Data:        reqObservable,
	}

	//запрос на поиск дополнительной информации об TTL
	reqTTP := NewChannelRequest()
	reqTTP.SetRequestId(uuidStorage)
	reqTTP.SetRootId(rootId)
	reqTTP.SetCommand("get_ttp")
	reqTTP.SetChanOutput(chanResTTL)
	chanInput <- ChanFormWebHookServer{
		ForSomebody: "for thehive",
		Data:        reqTTP,
	}

	wg.Wait()

	return rmec, mainErr
}
