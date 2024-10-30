package testthehivegetcase_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/av-belyakov/thehivehook_go_package/cmd/thehiveapi"
	"github.com/av-belyakov/thehivehook_go_package/internal/confighandler"
	"github.com/av-belyakov/thehivehook_go_package/internal/datamodels"
	"github.com/av-belyakov/thehivehook_go_package/internal/logginghandler"
)

var _ = Describe("Testthehivegetcase", Ordered, func() {
	apiTheHive := func(apiKey, host string, port int) ([]byte, int, error) {
		query, err := json.Marshal(thehiveapi.Querys{
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

	var errLoadEnv error

	BeforeAll(func() {
		errLoadEnv = godotenv.Load("../../.env")
		fmt.Println("ERROR env:", errLoadEnv)
		fmt.Println("API KEY:", os.Getenv("GO_HIVEHOOK_THAPIKEY"))
	})

	Context("Тест 0. Чтение переменных окружения", func() {
		It("При чтении переменных окружения не должно быть ошибок", func() {
			Expect(errLoadEnv).ShouldNot(HaveOccurred())
		})
	})

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
		})
	})

	Context("Тест 2. Запрос кейса по его номеру через метод GetCaseEvent", func() {
		It("При выполнении запроса на получении объекта 'event' Case ошибок быть не должно", func() {
			logging := logginghandler.New()
			conf := confighandler.AppConfigTheHive{
				Port:   9000,
				Host:   "1thehive.cloud.gcm",
				ApiKey: os.Getenv("GO_HIVEHOOK_THAPIKEY"),
			}
			apiTheHive, err := thehiveapi.New(
				logging,
				thehiveapi.WithAPIKey(conf.ApiKey),
				thehiveapi.WithHost(conf.Host),
				thehiveapi.WithPort(conf.Port))
			Expect(err).ShouldNot(HaveOccurred())

			_, err = apiTheHive.Start(context.Background())
			Expect(err).ShouldNot(HaveOccurred())

			b, code, err := apiTheHive.GetCaseEvent(context.Background(), "~88678416456" /*"~88325656792"*/)
			fmt.Println("ERROR:", err)
			fmt.Println("My error is exist:", errors.Is(err, datamodels.ConnectionError))
			Expect(err).ShouldNot(HaveOccurred())
			Expect(code).Should(Equal(http.StatusOK))

			//caseEvent := []map[string]interface{}{}
			caseEvent := []datamodels.BaseCaseEventElement(nil)
			err = json.Unmarshal(b, &caseEvent)
			Expect(err).ShouldNot(HaveOccurred())

			b, err = json.MarshalIndent(caseEvent, "", " ")
			Expect(err).ShouldNot(HaveOccurred())

			fmt.Println(string(b))

			Expect(true).Should(BeTrue())
		})
	})
})
