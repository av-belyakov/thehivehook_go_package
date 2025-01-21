package webhookserver

import (
	"encoding/json"

	"github.com/google/uuid"
	"golang.org/x/sync/errgroup"

	"github.com/av-belyakov/thehivehook_go_package/cmd/commoninterfaces"
)

// CreateEvenCase создает новый объект case, содержащий дополнительную информацию типа объектов observables
// и ttp информацию по которым дополнительно запрашивают из TheHive
func CreateEvenCase(rootId string, chanInput chan<- ChanFromWebHookServer) (ReadyMadeEventCase, error) {
	var (
		g    errgroup.Group
		rmec ReadyMadeEventCase = ReadyMadeEventCase{}

		uuidObservable string = uuid.NewString()
		uuidTTP        string = uuid.NewString()

		chanResObservable chan commoninterfaces.ChannelResponser = make(chan commoninterfaces.ChannelResponser)
		chanResTTL        chan commoninterfaces.ChannelResponser = make(chan commoninterfaces.ChannelResponser)
	)

	g.Go(func() error {
		for res := range chanResObservable {

			//fmt.Println("func 'CreateEventCase', goroutine 'observable' received data")

			msg := []interface{}{}
			if err := json.Unmarshal(res.GetData(), &msg); err != nil {
				return err
			}

			rmec.Observables = msg
		}

		return nil
	})
	g.Go(func() error {
		for res := range chanResTTL {

			//fmt.Println("func 'CreateEventCase', goroutine 'ttl' received data")

			msg := []interface{}{}
			if err := json.Unmarshal(res.GetData(), &msg); err != nil {
				return err
			}

			rmec.TTPs = msg
		}

		return nil
	})

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

	err := g.Wait()

	return rmec, err
}
