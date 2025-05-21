package webhookserver

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"golang.org/x/sync/errgroup"

	"github.com/av-belyakov/thehivehook_go_package/cmd/commoninterfaces"
	"github.com/av-belyakov/thehivehook_go_package/internal/datamodels"
)

// CreateEvenCase создает новый объект case, содержащий дополнительную информацию типа объектов observables
// и ttp информацию по которым дополнительно запрашивают из TheHive
func CreateEvenCase(ctx context.Context, rootId string, caseId int, chanInput chan<- ChanFromWebHookServer) (*ReadyMadeEventCase, error) {
	rmec := &ReadyMadeEventCase{}
	customError := &datamodels.CustomError{}

	chanResObservable := make(chan commoninterfaces.ChannelResponser)
	defer close(chanResObservable)

	chanResTTL := make(chan commoninterfaces.ChannelResponser)
	defer close(chanResTTL)

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	var g errgroup.Group
	g.Go(func() error {
		select {
		case <-ctx.Done():
			customError.Type = "context"
			customError.Err = ctx.Err()

			return customError

		case res := <-chanResObservable:
			msg := []any{}
			if err := json.Unmarshal(res.GetData(), &msg); err != nil {
				customError.Type = "json"
				customError.Err = err

				return customError
			}

			rmec.Observables = msg

			return nil
		}
	})
	g.Go(func() error {
		select {
		case <-ctx.Done():
			customError.Type = "context"
			customError.Err = ctx.Err()

			return customError

		case res := <-chanResTTL:
			msg := []any{}
			if err := json.Unmarshal(res.GetData(), &msg); err != nil {
				customError.Type = "json"
				customError.Err = err

				return customError
			}

			rmec.TTPs = msg

			return nil
		}
	})

	//запрос на поиск дополнительной информации об Observables
	reqObservable := NewChannelRequest()
	reqObservable.SetContext(ctx)
	reqObservable.SetRootId(rootId)
	reqObservable.SetCaseId(fmt.Sprint(caseId))
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
	reqTTP.SetCaseId(fmt.Sprint(caseId))
	reqTTP.SetCommand("get_ttp")
	reqTTP.SetChanOutput(chanResTTL)
	chanInput <- ChanFromWebHookServer{
		ForSomebody: "to thehive",
		Data:        reqTTP,
	}

	err := g.Wait()

	return rmec, err
}
