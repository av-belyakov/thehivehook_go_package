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
)

func New(apiKey, host string, port int) (*TheHiveAPI, error) {
	hapi := &TheHiveAPI{}

	if apiKey == "" {
		return hapi, errors.New("the value of 'apiKey' cannot be empty")
	}

	if host == "" {
		return hapi, errors.New("the value of 'host' cannot be empty")
	}

	if port == 0 || port > 65535 {
		return hapi, errors.New("an incorrect network port value was received")
	}

	hapi.apiKey = apiKey
	hapi.host = host
	hapi.port = port

	return hapi, nil
}

func (h *TheHiveAPI) GetObservables(ctx context.Context, id string) ([]byte, error) {
	query := &RootQuery{
		Query: []Query{
			{Name: "getCase", IDOrName: id},
			{Name: "observables"},
		},
	}

	path := "/api/v1/query?name=case-observables"

	reqbody, err := json.Marshal(query)
	if err != nil {
		_, f, l, _ := runtime.Caller(0)

		return nil, fmt.Errorf("%w %s:%d", err, f, l-2)
	}

	resp, err := h.query(ctx, path, reqbody, "POST")
	if err != nil {
		_, f, l, _ := runtime.Caller(0)

		return nil, fmt.Errorf("%w %s:%d", err, f, l-2)
	}

	return resp, nil
}

func (h *TheHiveAPI) query(ctx context.Context, path string, query []byte, method string) ([]byte, error) {
	bearer := "Bearer " + h.apiKey
	url := fmt.Sprintf("http://%s:%d%s", h.host, h.port, path)

	req, err := http.NewRequestWithContext(ctx, method, url, bytes.NewBuffer(query))
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", bearer)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	res, err := client.Do(req)
	defer func(body io.ReadCloser) {
		body.Close()
	}(res.Body)
	if err != nil {
		return nil, err
	}

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	/*var errAns ErrorAnswer
	_ = json.Unmarshal(resBody, &errAns)
	if errAns.Error() != "" { //nolint:wsl //b
		return resBody, &errAns
	}*/
	if res.StatusCode != http.StatusOK {
		_, f, l, _ := runtime.Caller(0)

		return resBody, fmt.Errorf("error sending the request, response status is %s %s:%d", res.Status, f, l-1)
	}

	return resBody, nil
}