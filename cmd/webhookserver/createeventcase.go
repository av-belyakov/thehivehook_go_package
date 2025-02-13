package webhookserver

import (
	"encoding/json"
	"fmt"

	"golang.org/x/sync/errgroup"

	"github.com/av-belyakov/thehivehook_go_package/cmd/commoninterfaces"
)

// CreateEvenCase создает новый объект case, содержащий дополнительную информацию типа объектов observables
// и ttp информацию по которым дополнительно запрашивают из TheHive
func CreateEvenCase(rootId string, chanInput chan<- ChanFromWebHookServer) (ReadyMadeEventCase, error) {
	fmt.Printf("func 'CreateEvenCase', rootId '%s' START...\n", rootId)

	var (
		g    errgroup.Group
		rmec ReadyMadeEventCase = ReadyMadeEventCase{}

		chanResObservable chan commoninterfaces.ChannelResponser = make(chan commoninterfaces.ChannelResponser)
		chanResTTL        chan commoninterfaces.ChannelResponser = make(chan commoninterfaces.ChannelResponser)
	)

	g.Go(func() error {
		for res := range chanResObservable {
			msg := []interface{}{}
			if err := json.Unmarshal(res.GetData(), &msg); err != nil {
				return err
			}

			rmec.Observables = msg
		}

		fmt.Println("__________________ goroutine g.Go Observable STOP")

		return nil
	})
	g.Go(func() error {
		for res := range chanResTTL {
			msg := []interface{}{}
			if err := json.Unmarshal(res.GetData(), &msg); err != nil {
				return err
			}

			rmec.TTPs = msg
		}

		fmt.Println("__________________ goroutine g.Go TTP STOP")

		return nil
	})

	//запрос на поиск дополнительной информации об Observables
	reqObservable := NewChannelRequest()
	reqObservable.SetRootId(rootId)
	reqObservable.SetCommand("get_observables")
	reqObservable.SetChanOutput(chanResObservable)
	chanInput <- ChanFromWebHookServer{
		ForSomebody: "to thehive",
		Data:        reqObservable,
	}

	//запрос на поиск дополнительной информации об TTL
	reqTTP := NewChannelRequest()
	reqTTP.SetRootId(rootId)
	reqTTP.SetCommand("get_ttp")
	reqTTP.SetChanOutput(chanResTTL)
	chanInput <- ChanFromWebHookServer{
		ForSomebody: "to thehive",
		Data:        reqTTP,
	}

	err := g.Wait()

	fmt.Println("func 'CreateEvenCase', STOP...")

	//что бы данную гроутину не держала ссылка на объекты
	reqObservable = NewChannelRequest()
	reqTTP = NewChannelRequest()

	return rmec, err
}
