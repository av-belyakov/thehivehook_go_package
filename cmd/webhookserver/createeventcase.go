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
func CreateEvenCase(ctx context.Context, rootId string, caseId int, chanInput chan<- ChanFromWebHookServer) (ReadyMadeEventCase, error) {
	rmec := ReadyMadeEventCase{}

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
			//
			// для того что бы исключить ошибку типа
			// 2025-02-14 14:46:24 ERR - thehivehook_go_package - context deadline
			// exceeded /home/artemij/go/src/thehivehook_go_package/cmd/webhookserver/routes.go:53
			// можно убрать return ctx.Err() и оставить просто return
			// вот только надо ли, пока не знаю
			//
			//return ctx.Err()

			return nil

		case res := <-chanResObservable:
			msg := []interface{}{}
			if err := json.Unmarshal(res.GetData(), &msg); err != nil {
				return err
			}

			rmec.Observables = msg
		}

		return nil
	})
	g.Go(func() error {
		select {
		case <-ctx.Done():
			//
			// для того что бы исключить ошибку типа
			// 2025-02-14 14:46:24 ERR - thehivehook_go_package - context deadline
			// exceeded /home/artemij/go/src/thehivehook_go_package/cmd/webhookserver/routes.go:53
			// можно убрать return ctx.Err() и оставить просто return
			// вот только надо ли, пока не знаю
			//
			//return ctx.Err()

			return nil

		case res := <-chanResTTL:
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

	//что бы данную гроутину не держала ссылка на объекты
	reqObservable = NewChannelRequest()
	reqTTP = NewChannelRequest()

	return rmec, err
}
