package webhookserver

import (
	"encoding/json"
	"sync"

	"github.com/google/uuid"

	"github.com/av-belyakov/thehivehook_go_package/cmd/commoninterfaces"
)

// CreateEvenCase создает новый объект case, содержащий дополнительную информацию типа объектов observables
// и ttp информацию по которым дополнительно запрашивают из TheHive
func CreateEvenCase(rootId string, chanInput chan<- ChanFromWebHookServer) (ReadyMadeEventCase, error) {
	var (
		wg   sync.WaitGroup
		rmec ReadyMadeEventCase = ReadyMadeEventCase{}

		uuidObservable string = uuid.NewString()
		uuidTTP        string = uuid.NewString()

		mainErr error

		chanResObservable chan commoninterfaces.ChannelResponser = make(chan commoninterfaces.ChannelResponser)
		chanResTTL        chan commoninterfaces.ChannelResponser = make(chan commoninterfaces.ChannelResponser)
	)

	wg.Add(1)
	go func(wg *sync.WaitGroup, chRes <-chan commoninterfaces.ChannelResponser) {
		defer wg.Done()

		for res := range chRes {
			msg := []interface{}{}
			if err := json.Unmarshal(res.GetData(), &msg); err != nil {
				mainErr = err

				return
			}

			rmec.Observables = msg
		}
	}(&wg, chanResObservable)

	//запрос на поиск дополнительной информации об Observables
	reqObservable := NewChannelRequest()
	reqObservable.SetRequestId(uuidObservable)
	reqObservable.SetRootId(rootId)
	reqObservable.SetCommand("get_observables")
	reqObservable.SetChanOutput(chanResObservable)
	chanInput <- ChanFromWebHookServer{
		ForSomebody: "to thehive",
		Data:        reqObservable,
	}

	wg.Add(1)
	go func(wg *sync.WaitGroup, chRes <-chan commoninterfaces.ChannelResponser) {
		defer wg.Done()

		for res := range chRes {
			msg := []interface{}{}
			if err := json.Unmarshal(res.GetData(), &msg); err != nil {
				mainErr = err

				return
			}

			rmec.TTPs = msg
		}
	}(&wg, chanResTTL)

	//запрос на поиск дополнительной информации об TTL
	reqTTP := NewChannelRequest()
	reqTTP.SetRequestId(uuidTTP)
	reqTTP.SetRootId(rootId)
	reqTTP.SetCommand("get_ttp")
	reqTTP.SetChanOutput(chanResTTL)
	chanInput <- ChanFromWebHookServer{
		ForSomebody: "to thehive",
		Data:        reqTTP,
	}

	wg.Wait()

	return rmec, mainErr
}
