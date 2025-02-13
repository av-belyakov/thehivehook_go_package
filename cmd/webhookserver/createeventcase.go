package webhookserver

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"golang.org/x/sync/errgroup"

	"github.com/av-belyakov/thehivehook_go_package/cmd/commoninterfaces"
)

// CreateEvenCase создает новый объект case, содержащий дополнительную информацию типа объектов observables
// и ttp информацию по которым дополнительно запрашивают из TheHive
func CreateEvenCase(ctx context.Context, rootId string, chanInput chan<- ChanFromWebHookServer) (ReadyMadeEventCase, error) {
	var (
		g    errgroup.Group
		rmec ReadyMadeEventCase = ReadyMadeEventCase{}
	)

	chanResObservable := make(chan commoninterfaces.ChannelResponser)
	defer close(chanResObservable)

	chanResTTL := make(chan commoninterfaces.ChannelResponser)
	defer close(chanResTTL)

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	fmt.Printf("!!! func 'CreateEvenCase', root id:'%s' START\n", rootId)

	g.Go(func() error {
		select {
		case <-ctx.Done():
			return ctx.Err()

		case res := <-chanResObservable:
			fmt.Println("func 'CreateEventCase', goroutine 'observable' received data")

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
		select {
		case <-ctx.Done():
			return ctx.Err()

		case res := <-chanResTTL:
			fmt.Println("func 'CreateEventCase', goroutine 'observable' received data")

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
	reqObservable.SetContext(ctx)
	reqObservable.SetRootId(rootId)
	reqObservable.SetCommand("get_observables")
	reqObservable.SetChanOutput(chanResObservable)
	chanInput <- ChanFromWebHookServer{
		ForSomebody: "to thehive",
		Data:        reqObservable,
	}

	//запрос на поиск дополнительной информации об TTL
	reqTTP := NewChannelRequest()
	reqTTP.SetContext(ctx)
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
