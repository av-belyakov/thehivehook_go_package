package testthehivegetcase_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/av-belyakov/thehivehook_go_package/cmd/thehiveapi"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Testthehivegetcase", func() {
	apiTheHive := func(apiKey, host string, port int) ([]byte, int, error) {
		query, err := json.Marshal(thehiveapi.RootQuery{
			Query: []thehiveapi.Query{
				{Name: "getCase", IDOrName: "~86676517008"},
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

	Context("Тест 1. Запрос кейса по его номеру", func() {
		It("Запрос должен быть успешно выполнен", func() {
			b, statusCode, err := apiTheHive(os.Getenv("GO_HIVEHOOK_THAPIKEY"), "thehive.cloud.gcm", 9000)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(statusCode).Should(Equal(200))

			caseEvent := []interface{}{}
			err = json.Unmarshal(b, &caseEvent)
			Expect(err).ShouldNot(HaveOccurred())

			b, err = json.MarshalIndent(caseEvent, "", " ")
			Expect(err).ShouldNot(HaveOccurred())

			fmt.Println(string(b))

			Expect(true).Should(BeTrue())

			/********************************

			Этот тест с observables работает, теперь надо попробовать
			с помощью него получить кейс

			*********************************/
		})
	})
})
