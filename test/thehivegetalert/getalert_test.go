package thehivegetalert

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"testing"

	"github.com/av-belyakov/thehivehook_go_package/cmd/thehiveapi"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

func apiTheHive(apiKey, host string, port int, rootId string) ([]byte, int, error) {
	query, err := json.Marshal(thehiveapi.Querys{
		Query: []thehiveapi.Query{
			{Name: "getAlert", IDOrName: rootId},
			//{Name: "observables"},
		},
	})
	if err != nil {
		return nil, 0, err
	}

	url := fmt.Sprintf("http://%s:%d%s", host, port, "/api/v1/query?name=case")
	req, err := http.NewRequestWithContext(context.Background(), "POST", url, bytes.NewBuffer(query))
	if err != nil {
		return nil, 0, err
	}

	req.Header.Add("Authorization", "Bearer "+apiKey)
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

func TestGetAlert(t *testing.T) {
	err := godotenv.Load("../../.env")
	assert.NoError(t, err)

	apiKey := os.Getenv("GO_HIVEHOOK_THAPIKEY")
	fmt.Println("API KEY:", apiKey)

	b, statusCode, err := apiTheHive(apiKey, "thehive.cloud.gcm", 9000, "~91355930856")
	assert.NoError(t, err)
	assert.Equal(t, statusCode, 200)
	t.Log(string(b))
}
