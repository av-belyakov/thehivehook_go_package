package webhookserver

import (
	"context"
)

// CreateEvenAlert генератор алерта содержащего в себе дополнительную информацию
func CreateEvenAlert(ctx context.Context, rootId string, chanInput chan<- ChanFromWebHookServer) (ReadyMadeEventAlert, error) {
	var rmea ReadyMadeEventAlert = ReadyMadeEventAlert{}

	/*
		chanRes := make(chan commoninterfaces.ChannelResponser)
		defer close(chanRes)

		ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
		defer cancel()

		var wg sync.WaitGroup

		go func() {
			select {
			case <-ctx.Done():
				createCaseErr.Type = "context"
				createCaseErr.Err = ctx.Err()

				return createCaseErr

			case res := <-chanRes:

			}
		}()

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

		wg.Wait()

		//попробовать запросить весь alert
		//http://192.168.9.38:9000/api/v1/alert/~74465718400

		/*
		   			query
		   :
		   [{_name: "getAlert", idOrName: "~76432666760"},…]
		   0
		   :
		   {_name: "getAlert", idOrName: "~76432666760"}
		   1
		   :
		   {_name: "similarCases", caseFilter: {_and: [{_field: "status", _value: "Open"},…]}}
		   caseFilter
		   :
		   {_and: [{_field: "status", _value: "Open"},…]}
	*/

	return rmea, nil
}
