package webhookserver

import (
	"context"
	"encoding/json"
	"time"

	"golang.org/x/sync/errgroup"

	"github.com/av-belyakov/thehivehook_go_package/cmd/commoninterfaces"
	"github.com/av-belyakov/thehivehook_go_package/internal/datamodels"
)

// CreateEvenAlert генератор события типа 'alert' содержащего в себе дополнительную информацию
func CreateEvenAlert(ctx context.Context, rootId string, chanInput chan<- ChanFromWebHookServer) (*ReadyMadeEventAlert, error) {
	rmea := &ReadyMadeEventAlert{}
	customError := &datamodels.CustomError{}

	chanResAlert := make(chan commoninterfaces.ChannelResponser)
	defer close(chanResAlert)

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	var g errgroup.Group
	g.Go(func() error {
		select {
		case <-ctx.Done():
			customError.Type = "context"
			customError.Err = ctx.Err()

			return customError

		case res := <-chanResAlert:
			msg := map[string]any{}
			if err := json.Unmarshal(res.GetData(), &msg); err != nil {
				customError.Type = "json"
				customError.Err = err

				return customError
			}

			rmea.Alert = msg

			return nil
		}
	})

	//запрос на поиск дополнительной информации об Alert
	reqObservable := NewChannelRequest()
	reqObservable.SetContext(ctx)
	reqObservable.SetRootId(rootId)
	reqObservable.SetCommand("get_alert")
	reqObservable.SetChanOutput(chanResAlert)
	chanInput <- ChanFromWebHookServer{
		ForSomebody: "to thehive",
		Data:        reqObservable,
	}

	err := g.Wait()

	return rmea, err
}
