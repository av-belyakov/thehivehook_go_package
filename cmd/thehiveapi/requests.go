package thehiveapi

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/av-belyakov/thehivehook_go_package/internal/datamodels"
	"github.com/av-belyakov/thehivehook_go_package/internal/supportingfunctions"
)

// GetAlert формирует запрос на получения из TheHive объекта типа 'alert'
func (api *apiTheHiveModule) GetAlert(ctx context.Context, rootId string) ([]byte, int, error) {
	ctxTimeout, ctxCancel := context.WithTimeout(ctx, 7*time.Second)
	defer ctxCancel()

	res, statusCode, err := api.query(ctxTimeout, fmt.Sprintf("/api/alert/%s", rootId), []byte{}, "GET")
	if err != nil {
		return nil, statusCode, supportingfunctions.CustomError(err)
	}

	return res, statusCode, err
}

// GetObservables формирует запрос на получения из TheHive объекта типа 'observables'
func (api *apiTheHiveModule) GetObservables(ctx context.Context, rootId string) ([]byte, int, error) {
	req, err := json.Marshal(Querys{
		Query: []Query{
			{Name: "getCase", IDOrName: rootId},
			{Name: "observables"},
		},
	})
	if err != nil {
		return nil, 0, supportingfunctions.CustomError(err)
	}

	ctxTimeout, ctxCancel := context.WithTimeout(ctx, 7*time.Second)
	defer ctxCancel()

	res, statusCode, err := api.query(ctxTimeout, "/api/v1/query?name=case-observables", req, "POST")
	if err != nil {
		return nil, statusCode, supportingfunctions.CustomError(err)
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
		return nil, 0, supportingfunctions.CustomError(err)
	}

	ctxTimeout, ctxCancel := context.WithTimeout(ctx, 7*time.Second)
	defer ctxCancel()

	res, statusCode, err := api.query(ctxTimeout, "/api/v1/query?name=case-procedures", req, "POST")
	if err != nil {
		return nil, statusCode, supportingfunctions.CustomError(err)
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
		return nil, 0, supportingfunctions.CustomError(err)
	}

	ctxTimeout, ctxCancel := context.WithTimeout(ctx, 30*time.Second)
	defer ctxCancel()

	res, statusCode, err := api.query(ctxTimeout, "/api/v1/query?name=case", req, "POST")
	if err != nil {
		return nil, statusCode, supportingfunctions.CustomError(err)
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

	fmt.Println("func 'apiTheHiveModule.AddCaseTags', START")

	//получаем информацию по кейсу
	bodyByte, statusCode, err := api.GetCaseEvent(ctx, rc.RootId)
	if err != nil {
		return nil, statusCode, supportingfunctions.CustomError(err)
	}
	if statusCode != http.StatusOK {
		return nil, statusCode, supportingfunctions.CustomError(fmt.Errorf("'when executing the Get Case Event request, the response status is received %d'", statusCode))
	}

	bcee := []datamodels.BaseCaseEventElement{}
	err = json.Unmarshal(bodyByte, &bcee)
	if err != nil {
		return nil, statusCode, supportingfunctions.CustomError(err)
	}

	if len(bcee) == 0 {
		return nil, statusCode, supportingfunctions.CustomError(fmt.Errorf("'no events were found in TheHive by rootId %s'", rc.RootId))
	}

	//получаем список тегов которых нет bcee[0].Tags, но есть в tags
	listUniqTags := supportingfunctions.CompareTwoSlices(bcee[0].Tags, []string{rc.Value})
	//если listUniqTags пустой то команда на добавление в TheHive дополнительных
	//тегов бессмыслена, так как либо список tags пустой, либо в bcee[0].Tags есть все
	//значения из tags

	fmt.Println("func 'apiTheHiveModule.AddCaseTags', получаем список тегов которых нет bcee[0].Tags, но есть в tags", listUniqTags)

	if len(listUniqTags) == 0 {
		api.logger.Send("info", fmt.Sprintf("the command to add the tag '%s' to TheHive for rootId '%s' was not executed", rc.Value, rc.RootId))

		return nil, statusCode, nil
	}

	//формируем тело запроса из новых тегов и уже существующих
	newTagsQuery := tagsQuery{Tags: bcee[0].Tags}
	newTagsQuery.Tags = append(newTagsQuery.Tags, listUniqTags...)

	fmt.Println("func 'apiTheHiveModule.AddCaseTags', формируем тело запроса из новых тегов и уже существующих, newTagsQuery", newTagsQuery)

	req, err := json.Marshal(newTagsQuery)
	if err != nil {
		return nil, statusCode, supportingfunctions.CustomError(err)
	}

	ctxTimeout, ctxCancel := context.WithTimeout(ctx, 15*time.Second)
	defer ctxCancel()

	fmt.Println("func 'apiTheHiveModule.AddCaseTags', делаем запрос ->")

	res, statusCode, err := api.query(ctxTimeout, fmt.Sprintf("/api/case/%s", rc.RootId), req, "PATCH")
	if err != nil {
		return nil, statusCode, supportingfunctions.CustomError(err)
	}

	return res, statusCode, nil
}

// AddCaseCustomFields просто добавляет новые customFields в объект Case TheHive без
// какого либо предварительного поиска и сверки customFields
func (api *apiTheHiveModule) AddCaseCustomFields(ctx context.Context, rc RequestCommand) ([]byte, int, error) {
	req := []byte(fmt.Sprintf(`{"customFields.%s": %q}`, rc.FieldName, rc.Value))

	ctxTimeout, ctxCancel := context.WithTimeout(ctx, 15*time.Second)
	defer ctxCancel()

	res, statusCode, err := api.query(ctxTimeout, fmt.Sprintf("/api/case/%s", rc.RootId), req, "PATCH")
	if err != nil {
		return nil, statusCode, supportingfunctions.CustomError(err)
	}

	return res, statusCode, nil
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

	ctxTimeout, ctxCancel := context.WithTimeout(ctx, 15*time.Second)
	defer ctxCancel()

	res, statusCode, err := api.query(ctxTimeout, fmt.Sprintf("/api/case/%s/task", rc.RootId), req, "POST")
	if err != nil {
		return nil, statusCode, supportingfunctions.CustomError(err)
	}

	return res, statusCode, nil
}

// query функция реализующая непосредственно сам HTTP запрос
func (api *apiTheHiveModule) query(ctx context.Context, reqpath string, query []byte, method string) ([]byte, int, error) {
	var (
		res *http.Response
		err error
	)

	apiKey := "Bearer " + api.apiKey
	url := fmt.Sprintf("http://%s:%d%s", api.host, api.port, reqpath)

	req, err := http.NewRequestWithContext(ctx, method, url, bytes.NewBuffer(query))
	if err != nil {
		return nil, 0, err
	}

	req.Header.Add("Authorization", apiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	res, err = client.Do(req)
	if err != nil {
		return nil, 0, fmt.Errorf("%w, %w", err, datamodels.ConnectionError)
	}
	defer func(body io.ReadCloser, err error) {
		if errClose := body.Close(); errClose != nil {
			errors.Join(err, errClose)
		}
	}(res.Body, err)

	if res.StatusCode == http.StatusOK || res.StatusCode == http.StatusCreated {
		resBody, err := io.ReadAll(res.Body)
		if err != nil {
			return nil, res.StatusCode, err
		}

		return resBody, res.StatusCode, nil
	}

	var msg string
	if m, err := supportingfunctions.GetDetailedMessage(res.Body); err == nil {
		msg = m
	}

	return nil, res.StatusCode, fmt.Errorf("error request, status is '%s' %s", res.Status, msg)
}
