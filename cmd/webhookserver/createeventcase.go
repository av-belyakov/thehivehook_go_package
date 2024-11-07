package webhookserver

import (
	"encoding/json"
	"sync"

	"github.com/av-belyakov/thehivehook_go_package/cmd/commoninterfaces"
)

// CreateEvenCase создает новый объект case, содержащий дополнительную информацию типа объектов observables
// и ttp информацию по которым дополнительно запрашивают из TheHive
func CreateEvenCase(rootId string, chanInput chan<- ChanFromWebHookServer) (ReadyMadeEventCase, error) {
	var (
		wg   sync.WaitGroup
		rmec ReadyMadeEventCase = ReadyMadeEventCase{}

		mainErr error

		chanResObservable chan commoninterfaces.ChannelResponser = make(chan commoninterfaces.ChannelResponser)
		chanResTTL        chan commoninterfaces.ChannelResponser = make(chan commoninterfaces.ChannelResponser)
	)

	wg.Add(2)

	go func(wg *sync.WaitGroup, chRes <-chan commoninterfaces.ChannelResponser) {
		for res := range chRes {
			msg := []interface{}{}
			if err := json.Unmarshal(res.GetData(), &msg); err != nil {
				mainErr = err

				return
			}

			rmec.Observables = msg
		}

		wg.Done()
	}(&wg, chanResObservable)
	go func(wg *sync.WaitGroup, chRes <-chan commoninterfaces.ChannelResponser) {
		for res := range chRes {
			msg := []interface{}{}
			if err := json.Unmarshal(res.GetData(), &msg); err != nil {
				mainErr = err

				return
			}

			rmec.TTPs = msg
		}

		wg.Done()
	}(&wg, chanResTTL)

	//запрос на поиск дополнительной информации об Observables
	reqObservable := NewChannelRequest()
	reqObservable.SetRootId(rootId)
	reqObservable.SetCommand("get_observables")
	reqObservable.SetChanOutput(chanResObservable)
	chanInput <- ChanFromWebHookServer{
		ForSomebody: "for thehive",
		Data:        reqObservable,
	}

	//запрос на поиск дополнительной информации об TTL
	reqTTP := NewChannelRequest()
	reqTTP.SetRootId(rootId)
	reqTTP.SetCommand("get_ttp")
	reqTTP.SetChanOutput(chanResTTL)
	chanInput <- ChanFromWebHookServer{
		ForSomebody: "for thehive",
		Data:        reqTTP,
	}

	wg.Wait()

	return rmec, mainErr
}
