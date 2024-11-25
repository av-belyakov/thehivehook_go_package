package webhookserver

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/av-belyakov/thehivehook_go_package/cmd/commoninterfaces"
	"github.com/google/uuid"
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

	fmt.Println("func 'CreateEvenCase' .......... START")

	wg.Add(1)
	go func(wg *sync.WaitGroup, chRes <-chan commoninterfaces.ChannelResponser) {
		defer wg.Done()

		fmt.Println("GOROUTINE 1, func 'CreateEvenCase' START")

		for res := range chRes {
			fmt.Println("func 'CreateEvenCase' .......... RESPONSE 'observables'")

			msg := []interface{}{}
			if err := json.Unmarshal(res.GetData(), &msg); err != nil {
				mainErr = err

				fmt.Println("func 'CreateEvenCase' .......... RESPONSE 'observables' ERRORR:", err.Error())

				return
			}

			rmec.Observables = msg
		}

		fmt.Println("func 'CreateEvenCase' GOROUTINE 'observables' STOP")

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

		fmt.Println("GOROUTINE 2, func 'CreateEvenCase' START")

		for res := range chRes {
			fmt.Println("func 'CreateEvenCase' .......... RESPONSE 'ttp'")

			msg := []interface{}{}
			if err := json.Unmarshal(res.GetData(), &msg); err != nil {
				mainErr = err

				fmt.Println("func 'CreateEvenCase' .......... RESPONSE 'ttp' ERRORR:", err.Error())

				return
			}

			rmec.TTPs = msg
		}

		fmt.Println("func 'CreateEvenCase' GOROUTINE 'ttp' STOP")

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
