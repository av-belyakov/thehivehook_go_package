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

	"github.com/av-belyakov/thehivehook_go_package/cmd/commoninterfaces"
	"github.com/av-belyakov/thehivehook_go_package/internal/datamodels"
	"github.com/av-belyakov/thehivehook_go_package/internal/supportingfunctions"
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

	ctxTimeout, ctxCancel := context.WithTimeout(ctx, 5*time.Second)
	defer ctxCancel()

	res, statusCode, err := api.query(ctxTimeout, "/api/v1/query?name=case", req, "POST")
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

	//получаем информацию по кейсу
	bodyByte, statusCode, err := api.GetCaseEvent(ctx, rootId)
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

	//получаем список тегов которых нет bcee[0].Tags, но есть в tags
	listUniqTags := supportingfunctions.CompareTwoSlices(bcee[0].Tags, tags)
	//если listUniqTags пустой то команда на добавление в TheHive дополнительных
	//тегов бессмыслена, так как либо список tags пустой, либо в bcee[0].Tags есть все
	//значения из tags
	if len(listUniqTags) == 0 {
		api.logger.Send("info", fmt.Sprintf("the command to add the tags list '%v' to TheHive for rootId '%s' was not executed", tags, rootId))

		fmt.Println("func 'AddCaseTags' 555.5, listUniqTags:", listUniqTags)

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

	res, statusCode, err := api.query(ctxTimeout, fmt.Sprintf("/api/case/%s", rootId), req, "PATCH")
	if err != nil {
		_, f, l, _ := runtime.Caller(0)
		return nil, 0, fmt.Errorf("%w %s:%d", err, f, l-2)
	}

	return res, statusCode, err
}

// AddCaseTags добавляет новые customFields в TheHive
func (api *apiTheHiveSettings) AddCaseCustomFields(ctx context.Context, rootId string, i interface{}) ([]byte, int, error) {
	cfp, ok := i.(commoninterfaces.ParametersCustomFieldKeeper)
	if !ok {
		_, f, l, _ := runtime.Caller(0)
		return nil, 0, fmt.Errorf("'it is not possible to convert a value to a commoninterfaces.ParametersCustomFieldKeeper interface' %s:%d", f, l-2)
	}

	req := []byte(fmt.Sprintf(`{"customFields.%s": %q}`, cfp.GetType(), cfp.GetValue()))
	ctxTimeout, ctxCancel := context.WithTimeout(ctx, 5*time.Second)
	defer ctxCancel()

	res, statusCode, err := api.query(ctxTimeout, fmt.Sprintf("/api/case/%s", rootId), req, "PATCH")
	if err != nil {
		_, f, l, _ := runtime.Caller(0)
		return nil, 0, fmt.Errorf("%w %s:%d", err, f, l-2)
	}

	return res, statusCode, err
}

// AddCaseTask добавляет новую задачу в TheHive
func (api *apiTheHiveSettings) AddCaseTask(ctx context.Context, rootId string, i interface{}) ([]byte, int, error) {
	tp, ok := i.(commoninterfaces.ParametersCustomFieldKeeper)
	if !ok {
		_, f, l, _ := runtime.Caller(0)
		return nil, 0, fmt.Errorf("'it is not possible to convert a value to a commoninterfaces.ParametersCustomFieldKeeper interface' %s:%d", f, l-2)
	}

	var req []byte
	if tp.GetUsername() == "" {
		req = []byte(fmt.Sprintf(`{"status":"Waiting","group":%q,"title":%q}`, tp.GetType(), tp.GetValue()))
	} else {
		req = []byte(fmt.Sprintf(`{"status":"Waiting","group":%q,"title":%q,"owner":%q}`, tp.GetType(), tp.GetValue(), tp.GetUsername()))
	}

	ctxTimeout, ctxCancel := context.WithTimeout(ctx, 5*time.Second)
	defer ctxCancel()

	res, statusCode, err := api.query(ctxTimeout, fmt.Sprintf("/api/case/%s/task", rootId), req, "POST")
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
