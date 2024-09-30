package thehiveapi

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"runtime"

	"github.com/av-belyakov/thehivehook_go_package/cmd/commoninterfaces"
	"github.com/av-belyakov/thehivehook_go_package/internal/logginghandler"
)

// New инициализирует новый модуль взаимодействия с API TheHive
// при инициализации возращается канал для взаимодействия с модулем, все
// запросы к модулю выполняются через данный канал
func New(ctx context.Context, apiKey, host string, port int, logging *logginghandler.LoggingChan) (chan<- commoninterfaces.ChannelRequester, error) {
	receivingChannel := make(chan commoninterfaces.ChannelRequester)

	if apiKey == "" {
		return receivingChannel, errors.New("the value of 'apiKey' cannot be empty")
	}

	if host == "" {
		return receivingChannel, errors.New("the value of 'host' cannot be empty")
	}

	if port == 0 || port > 65535 {
		return receivingChannel, errors.New("an incorrect network port value was received")
	}

	apiTheHive := &apiTheHive{
		apiKey: apiKey,
		host:   host,
		port:   port,
	}

	go func() {
		for {
			select {
			case <-ctx.Done():
				return

			case data := <-receivingChannel:
				switch data.GetCommand() {
				case "get_observables":
					res, statusCode, err := apiTheHive.GetObservables(ctx, data.GetRootId())
					if err != nil {
						logging.Send("error", err.Error())

						continue
					}

					newRes := NewChannelRespons()
					newRes.SetRequestId(data.GetRequestId())
					newRes.SetStatusCode(statusCode)
					newRes.SetData(res)

					data.GetChanOutput() <- newRes
					close(data.GetChanOutput())

					/*data.GetChanOutput() <- ResponseChannelTheHive{
						RequestId:  data.GetRequestId(),
						StatusCode: statusCode,
						Data:       res,
					}*/

				case "get_ttp":
					res, statusCode, err := apiTheHive.GetTTP(ctx, data.GetRootId())
					if err != nil {
						logging.Send("error", err.Error())

						continue
					}

					newRes := NewChannelRespons()
					newRes.SetRequestId(data.GetRequestId())
					newRes.SetStatusCode(statusCode)
					newRes.SetData(res)

					data.GetChanOutput() <- newRes
					close(data.GetChanOutput())

					/*
						data.ChanOutput <- ResponseChannelTheHive{
							RequestId:  data.GetRequestId(),
							StatusCode: statusCode,
							Data:       res,
						}

						close(data.ChanOutput)
					*/

				case "":
				}
			}
		}
	}()

	return receivingChannel, nil
}

func (api *apiTheHive) GetObservables(ctx context.Context, rootId string) ([]byte, int, error) {
	req, err := json.Marshal(RootQuery{
		Query: []Query{
			{Name: "getCase", IDOrName: rootId},
			{Name: "observables"},
		},
	})
	if err != nil {
		_, f, l, _ := runtime.Caller(0)

		return nil, 0, fmt.Errorf("%s %s:%d", err.Error(), f, l-2)
	}

	res, statusCode, err := api.query(ctx, "/api/v1/query?name=case-observables", req, "POST")
	if err != nil {
		_, f, l, _ := runtime.Caller(0)

		return nil, 0, fmt.Errorf("%s %s:%d", err.Error(), f, l-2)
	}

	return res, statusCode, err
}

func (api *apiTheHive) GetTTP(ctx context.Context, rootId string) ([]byte, int, error) {
	req, err := json.Marshal(&RootQuery{
		Query: []Query{
			{Name: "getCase", IDOrName: rootId},
			{Name: "procedures"},
			{
				Name: "page",
				From: 0,
				To:   999,
				ExtraData: []string{
					"pattern",
					"patternParent",
				},
			},
		},
	})
	if err != nil {
		_, f, l, _ := runtime.Caller(0)

		return nil, 0, fmt.Errorf("%s %s:%d", err.Error(), f, l-2)
	}

	res, statusCode, err := api.query(ctx, "/api/v1/query?name=case-procedures", req, "POST")
	if err != nil {
		_, f, l, _ := runtime.Caller(0)

		return nil, 0, fmt.Errorf("%s %s:%d", err.Error(), f, l-2)
	}

	return res, statusCode, err
}

func (api *apiTheHive) query(ctx context.Context, reqpath string, query []byte, method string) ([]byte, int, error) {
	apiKey := "Bearer " + api.apiKey
	url := fmt.Sprintf("http://%s:%d%s", api.host, api.port, reqpath)

	req, err := http.NewRequestWithContext(ctx, method, url, bytes.NewBuffer(query))
	if err != nil {
		return nil, 0, err
	}

	req.Header.Add("Authorization", apiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	res, err := client.Do(req)
	defer func(body io.ReadCloser) {
		body.Close()
	}(res.Body)
	if err != nil {
		return nil, 0, err
	}

	if res.StatusCode != http.StatusOK { //|| res.StatusCode != http.StatusCreated
		return nil, res.StatusCode, fmt.Errorf("error request, status is '%s'", res.Status)
	}

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, 0, err
	}

	return resBody, res.StatusCode, nil
}
