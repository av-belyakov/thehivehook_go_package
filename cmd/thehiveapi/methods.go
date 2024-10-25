// Модуль для взаимодействия с API TheHive
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
	temporarystoarge "github.com/av-belyakov/thehivehook_go_package/cmd/thehiveapi/temporarystorage"
)

// New настраивает модуль взаимодействия с API TheHive
func New(logger commoninterfaces.Logger, opts ...theHiveAPIOptions) (*apiTheHiveSettings, error) {
	ts, err := temporarystoarge.NewTemporaryStorage(30)
	if err != nil {
		return &apiTheHiveSettings{}, err
	}

	api := &apiTheHiveSettings{
		logger:           logger,
		receivingChannel: make(chan commoninterfaces.ChannelRequester),
		temporaryStorage: ts,
	}

	for _, opt := range opts {
		if err := opt(api); err != nil {
			return api, err
		}
	}

	return api, nil
}

// Start инициализирует новый модуль взаимодействия с API TheHive
// при инициализации возращается канал для взаимодействия с модулем, все
// запросы к модулю выполняются через данный канал
func (api *apiTheHiveSettings) Start(ctx context.Context) chan<- commoninterfaces.ChannelRequester {
	go func() {
		for {
			select {
			case <-ctx.Done():
				return

			case data := <-api.receivingChannel:
				switch data.GetCommand() {
				case "get_observables":
					res, statusCode, err := api.GetObservables(ctx, data.GetRootId())
					if err != nil {
						api.logger.Send("error", err.Error())

						continue
					}

					newRes := NewChannelRespons()
					newRes.SetRequestId(data.GetRequestId())
					newRes.SetStatusCode(statusCode)
					newRes.SetData(res)

					data.GetChanOutput() <- newRes
					close(data.GetChanOutput())

				case "get_ttp":
					res, statusCode, err := api.GetTTP(ctx, data.GetRootId())
					if err != nil {
						api.logger.Send("error", err.Error())

						continue
					}

					newRes := NewChannelRespons()
					newRes.SetRequestId(data.GetRequestId())
					newRes.SetStatusCode(statusCode)
					newRes.SetData(res)

					data.GetChanOutput() <- newRes
					close(data.GetChanOutput())

				case "send command":
					// Вот здесь нужно использовать temporaryStorage как кеширующее
					// хранилище команд корторые нужно отправить в TheHive и которые
					// из-за, по какой то причине, недоступности TheHive отправить сразу
					// не получается
				}
			}
		}
	}()

	return api.receivingChannel
}

// GetObservables формирует запрос на получения из TheHive объекта типа 'observables'
func (api *apiTheHiveSettings) GetObservables(ctx context.Context, rootId string) ([]byte, int, error) {
	req, err := json.Marshal(Querys{
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

// GetTTP формирует запрос на получения из TheHive объекта типа 'ttp'
func (api *apiTheHiveSettings) GetTTP(ctx context.Context, rootId string) ([]byte, int, error) {
	req, err := json.Marshal(&Querys{
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

// query функция реализующая непосредственно сам HTTP запрос
func (api *apiTheHiveSettings) query(ctx context.Context, reqpath string, query []byte, method string) ([]byte, int, error) {
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

// WithAPIKey метод устанавливает идентификатор-ключ для API
func WithAPIKey(v string) theHiveAPIOptions {
	return func(th *apiTheHiveSettings) error {
		if v == "" {
			return errors.New("the value of 'apiKey' cannot be empty")
		}

		th.apiKey = v

		return nil
	}
}

// WithHost метод устанавливает имя или ip адрес хоста API
func WithHost(v string) theHiveAPIOptions {
	return func(th *apiTheHiveSettings) error {
		if v == "" {
			return errors.New("the value of 'host' cannot be empty")
		}

		th.host = v

		return nil
	}
}

// WithPort метод устанавливает порт API
func WithPort(v int) theHiveAPIOptions {
	return func(th *apiTheHiveSettings) error {
		if v <= 0 || v > 65535 {
			return errors.New("an incorrect network port value was received")
		}

		th.port = v

		return nil
	}
}
