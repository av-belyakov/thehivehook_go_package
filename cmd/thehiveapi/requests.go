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
	"github.com/av-belyakov/thehivehook_go_package/internal/supportingfunctions"
)

// GetObservables формирует запрос на получения из TheHive объекта типа 'observables'
func (api *apiTheHiveModule) GetObservables(ctx context.Context, rootId string) ([]byte, int, error) {
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

	ctxTimeout, ctxCancel := context.WithTimeout(ctx, 5*time.Second)
	defer ctxCancel()

	res, statusCode, err := api.query(ctxTimeout, "/api/v1/query?name=case-observables", req, "POST")
	if err != nil {
		_, f, l, _ := runtime.Caller(0)
		return nil, 0, fmt.Errorf("%w %s:%d", err, f, l-2)
	}

	return res, statusCode, err
}

// GetTTP формирует запрос на получения из TheHive объекта типа 'ttp'
func (api *apiTheHiveModule) GetTTP(ctx context.Context, rootId string) ([]byte, int, error) {
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

	ctxTimeout, ctxCancel := context.WithTimeout(ctx, 5*time.Second)
	defer ctxCancel()

	res, statusCode, err := api.query(ctxTimeout, "/api/v1/query?name=case-procedures", req, "POST")
	if err != nil {
		_, f, l, _ := runtime.Caller(0)
		return nil, 0, fmt.Errorf("%w %s:%d", err, f, l-2)
	}

	return res, statusCode, err
}

// GetCaseEvent формирует запрос на получения из TheHive объекта типа 'event' являющегося
// основой Case
func (api *apiTheHiveModule) GetCaseEvent(ctx context.Context, rootId string) ([]byte, int, error) {
	req, err := json.Marshal(&Querys{
		Query: []Query{
			{Name: "getCase", IDOrName: rootId},
		},
	})
	if err != nil {
		_, f, l, _ := runtime.Caller(0)
		return nil, 0, fmt.Errorf("%w %s:%d", err, f, l-2)
	}

	ctxTimeout, ctxCancel := context.WithTimeout(ctx, 5*time.Second)
	defer ctxCancel()

	res, statusCode, err := api.query(ctxTimeout, "/api/v1/query?name=case", req, "POST")
	if err != nil {
		_, f, l, _ := runtime.Caller(0)
		return nil, 0, fmt.Errorf("%w %s:%d", err, f, l-2)
	}

	return res, statusCode, err
}

// AddCaseTags добавляет новые теги в объект Case, при этом получает от TheHive уже
// существующие теги и объединяет их с новыми, если добавляемые теги уже существуют
// то ничего не делает
func (api *apiTheHiveModule) AddCaseTags(ctx context.Context, rc RequestCommand) ([]byte, int, error) {
	type tagsQuery struct {
		Tags []string `json:"tags"`
	}

	//получаем информацию по кейсу
	bodyByte, statusCode, err := api.GetCaseEvent(ctx, rc.RootId)
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
		return nil, 0, fmt.Errorf("'no events were found in TheHive by rootId %s' %s:%d", rc.RootId, f, l-1)
	}

	//получаем список тегов которых нет bcee[0].Tags, но есть в tags
	listUniqTags := supportingfunctions.CompareTwoSlices(bcee[0].Tags, []string{rc.Value})
	//если listUniqTags пустой то команда на добавление в TheHive дополнительных
	//тегов бессмыслена, так как либо список tags пустой, либо в bcee[0].Tags есть все
	//значения из tags
	if len(listUniqTags) == 0 {
		api.logger.Send("info", fmt.Sprintf("the command to add the tag '%s' to TheHive for rootId '%s' was not executed", rc.Value, rc.RootId))

		return nil, 0, nil
	}

	//формируем тело запроса из новых тегов и уже существующих
	newTagsQuery := tagsQuery{Tags: bcee[0].Tags}
	newTagsQuery.Tags = append(newTagsQuery.Tags, listUniqTags...)

	req, err := json.Marshal(newTagsQuery)
	if err != nil {
		_, f, l, _ := runtime.Caller(0)
		return nil, 0, fmt.Errorf("%w %s:%d", err, f, l-2)
	}

	ctxTimeout, ctxCancel := context.WithTimeout(ctx, 5*time.Second)
	defer ctxCancel()

	res, statusCode, err := api.query(ctxTimeout, fmt.Sprintf("/api/case/%s", rc.RootId), req, "PATCH")
	if err != nil {
		_, f, l, _ := runtime.Caller(0)
		return nil, 0, fmt.Errorf("%w %s:%d", err, f, l-2)
	}

	return res, statusCode, err
}

// AddCaseCustomFields просто добавляет новые customFields в объект Case TheHive без
// какого либо предварительного поиска и сверки customFields
func (api *apiTheHiveModule) AddCaseCustomFields(ctx context.Context, rc RequestCommand) ([]byte, int, error) {
	req := []byte(fmt.Sprintf(`{"customFields.%s": %q}`, rc.FieldName, rc.Value))
	ctxTimeout, ctxCancel := context.WithTimeout(ctx, 5*time.Second)
	defer ctxCancel()

	res, statusCode, err := api.query(ctxTimeout, fmt.Sprintf("/api/case/%s", rc.RootId), req, "PATCH")
	if err != nil {
		_, f, l, _ := runtime.Caller(0)
		return nil, 0, fmt.Errorf("%w %s:%d", err, f, l-2)
	}

	return res, statusCode, err
}

// AddCaseTask просто добавляет новую задачу в объект Case TheHive без какого либо
// поиска и сравнения схожей задачи
func (api *apiTheHiveModule) AddCaseTask(ctx context.Context, rc RequestCommand) ([]byte, int, error) {
	var req []byte
	if rc.Username == "" {
		req = []byte(fmt.Sprintf(`{"status":"Waiting","group":%q,"title":%q}`, rc.FieldName, rc.Value))
	} else {
		req = []byte(fmt.Sprintf(`{"status":"Waiting","group":%q,"title":%q,"owner":%q}`, rc.FieldName, rc.Value, rc.Username))
	}

	ctxTimeout, ctxCancel := context.WithTimeout(ctx, 5*time.Second)
	defer ctxCancel()

	res, statusCode, err := api.query(ctxTimeout, fmt.Sprintf("/api/case/%s/task", rc.RootId), req, "POST")
	if err != nil {
		_, f, l, _ := runtime.Caller(0)
		return nil, 0, fmt.Errorf("%w %s:%d", err, f, l-2)
	}

	return res, statusCode, err
}

// query функция реализующая непосредственно сам HTTP запрос
func (api *apiTheHiveModule) query(ctx context.Context, reqpath string, query []byte, method string) ([]byte, int, error) {
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
		return nil, 0, fmt.Errorf("%w, %w", err, datamodels.ConnectionError)
	}
	defer func(body io.ReadCloser) {
		body.Close()
	}(res.Body)

	if res.StatusCode == http.StatusOK || res.StatusCode == http.StatusCreated {
		resBody, err := io.ReadAll(res.Body)
		if err != nil {
			return nil, 0, err
		}

		return resBody, res.StatusCode, nil
	}

	/*

	   надо понять почему приходят нулевые status code
	   и разобратся с ошибкой соединения с TheHive

	   ОШИБКИ! Возможно когда много кейсов

	   func 'automaticExecutionMethods' new tick: cacheFunc, id: 5dceea4e-bcf2-4084-b734-f0bd05562d79
	   func 'automaticExecutionMethods' new tick: cacheFunc, id: 5b5c371a-997f-4f86-aa67-0fd299e8562c
	   2024-11-26 15:20:13 ERR - thehivehook_go_package - Post "http://thehive.cloud.gcm:9000/api/v1/query?name=case-observables": context deadline exceeded, network connection error /home/artemij/go/src/thehivehook_go_package/cmd/thehiveapi/requests.go:33
	   func 'automaticExecutionMethods' new tick: cacheFunc, id: 2fe81847-3572-4475-82d3-f3fd52e6d6de
	   func 'automaticExecutionMethods' new tick: cacheFunc, id: b082e476-1cf4-44c4-a898-b41e5d6af6a6
	   func 'automaticExecutionMethods' new tick: cacheFunc, id: d011369e-9f85-48ed-9bbe-40baca020b48
	   2024-11-26 15:20:13 ERR - thehivehook_go_package - Post "http://thehive.cloud.gcm:9000/api/v1/query?name=case-procedures": context deadline exceeded, network connection error /home/artemij/go/src/thehivehook_go_package/cmd/thehiveapi/requests.go:67
	   2024-11-26 15:20:13 ERR - thehivehook_go_package - Post "http://thehive.cloud.gcm:9000/api/v1/query?name=case-observables": context deadline exceeded, network connection error /home/artemij/go/src/thehivehook_go_package/cmd/thehiveapi/requests.go:33
	   2024-11-26 15:20:13 ERR - thehivehook_go_package - Post "http://thehive.cloud.gcm:9000/api/v1/query?name=case-procedures": context deadline exceeded, network connection error /home/artemij/go/src/thehivehook_go_package/cmd/thehiveapi/requests.go:67


	*/

	var msg string
	if m, err := supportingfunctions.GetDetailedMessage(res.Body); err == nil {
		msg = m
	}

	return nil, res.StatusCode, fmt.Errorf("error request, status is '%s' %s", res.Status, msg)
}
