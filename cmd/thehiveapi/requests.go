package thehiveapi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"runtime"
	"time"

	"github.com/av-belyakov/thehivehook_go_package/internal/datamodels"
)

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
		return nil, 0, fmt.Errorf("%w %s:%d", err, f, l-2)
	}

	res, statusCode, err := api.query(ctx, "/api/v1/query?name=case-observables", req, "POST")
	if err != nil {
		_, f, l, _ := runtime.Caller(0)
		return nil, 0, fmt.Errorf("%w %s:%d", err, f, l-2)
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
		return nil, 0, fmt.Errorf("%w %s:%d", err, f, l-2)
	}

	res, statusCode, err := api.query(ctx, "/api/v1/query?name=case-procedures", req, "POST")
	if err != nil {
		_, f, l, _ := runtime.Caller(0)
		return nil, 0, fmt.Errorf("%w %s:%d", err, f, l-2)
	}

	return res, statusCode, err
}

// GetCaseEvent формирует запрос на получения из TheHive объекта типа 'event' являющегося
// основой Case
func (api *apiTheHiveSettings) GetCaseEvent(ctx context.Context, rootId string) ([]byte, int, error) {
	req, err := json.Marshal(&Querys{
		Query: []Query{
			{Name: "getCase", IDOrName: rootId},
		},
	})
	if err != nil {
		_, f, l, _ := runtime.Caller(0)
		return nil, 0, fmt.Errorf("%w %s:%d", err, f, l-2)
	}

	res, statusCode, err := api.query(ctx, "/api/v1/query?name=case", req, "POST")
	if err != nil {
		_, f, l, _ := runtime.Caller(0)
		return nil, 0, fmt.Errorf("%w %s:%d", err, f, l-2)
	}

	return res, statusCode, err
}

// AddCaseTags добавляет новые теги, при этом получает от TheHive уже существующие теги и объединяет их с новыми
func (api *apiTheHiveSettings) AddCaseTags(ctx context.Context, rootId string, i interface{}) ([]byte, int, error) {
	type tagsQuery struct {
		Tags []string `json:"tags"`
	}

	tags, ok := i.([]string)
	if !ok {
		_, f, l, _ := runtime.Caller(0)
		return nil, 0, fmt.Errorf("'it is not possible to convert a value to a []string' %s:%d", f, l-2)
	}

	ctxTimeout, ctxCancel := context.WithTimeout(ctx, 5*time.Second)
	defer ctxCancel()
	//получаем информацию по кейсу
	bodyByte, statusCode, err := api.GetCaseEvent(ctxTimeout, rootId)
	_, f, l, _ := runtime.Caller(0)
	if err != nil {
		return nil, 0, fmt.Errorf("%w %s:%d", err, f, l-1)
	}
	if statusCode != http.StatusOK {
		return nil, 0, fmt.Errorf("'when executing the Get Case Event request, the response status is received %d' %s:%d", statusCode, f, l-1)
	}

	bcee := []datamodels.BaseCaseEventElement{}
	err = json.Unmarshal(bodyByte, &bcee)
	if err != nil {
		_, f, l, _ := runtime.Caller(0)
		return nil, 0, fmt.Errorf("%w %s:%d", err, f, l-2)
	}

	if len(bcee) == 0 {
		return nil, 0, fmt.Errorf("'no events were found in TheHive by rootId %s' %s:%d", rootId, f, l-1)
	}

	//формируем тело запроса из новых тегов и уже существующих
	newTagsQuery := tagsQuery{Tags: bcee[0].Tags}
	newTagsQuery.Tags = append(newTagsQuery.Tags, tags...)

	req, err := json.Marshal(newTagsQuery)
	if err != nil {
		_, f, l, _ := runtime.Caller(0)
		return nil, 0, fmt.Errorf("%w %s:%d", err, f, l-2)
	}

	res, statusCode, err := api.query(ctx, fmt.Sprintf("/api/case/%s", rootId), req, "PATCH")
	if err != nil {
		_, f, l, _ := runtime.Caller(0)
		return nil, 0, fmt.Errorf("%w %s:%d", err, f, l-2)
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
	if err != nil {
		return nil, 0, fmt.Errorf("%w, %w", err, datamodels.ConnectionError) //errors.Join(err, datamodels.ConnectionError)
	}
	defer func(body io.ReadCloser) {
		body.Close()
	}(res.Body)

	if res.StatusCode != http.StatusOK {
		return nil, res.StatusCode, fmt.Errorf("error request, status is '%s'", res.Status)
	}

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, 0, err
	}

	return resBody, res.StatusCode, nil
}
