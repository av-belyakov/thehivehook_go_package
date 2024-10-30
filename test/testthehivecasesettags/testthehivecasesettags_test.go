package testthehivecasesettags_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/av-belyakov/thehivehook_go_package/cmd/thehiveapi"
	"github.com/av-belyakov/thehivehook_go_package/internal/confighandler"
	"github.com/av-belyakov/thehivehook_go_package/internal/logginghandler"
)

var _ = Describe("Testthehivecasesettags", Ordered, func() {
	var (
		rootId string = "~88678416456" //это мой тестовый кейс с id 39100

		errLoadEnv error
	)

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

	Context("Тест 1. Добавление тегов к заданному кейсы TheHive", func() {
		It("При выполнении запроса ошибок быть не должно", func() {
			logging := logginghandler.New()
			conf := confighandler.AppConfigTheHive{
				Port:   9000,
				Host:   "thehive.cloud.gcm",
				ApiKey: os.Getenv("GO_HIVEHOOK_THAPIKEY"),
			}
			apiTheHive, err := thehiveapi.New(
				logging,
				thehiveapi.WithAPIKey(conf.ApiKey),
				thehiveapi.WithHost(conf.Host),
				thehiveapi.WithPort(conf.Port))
			Expect(err).ShouldNot(HaveOccurred())

			_ = apiTheHive.Start(context.Background())

			b, code, err := apiTheHive.AddCaseTags(context.Background(), rootId, []string{"Webhook:send=\"WEBKOOK_mytest\""})
			Expect(err).ShouldNot(HaveOccurred())
			Expect(code).Should(Equal(http.StatusOK))

			event := map[string]interface{}{}
			err = json.Unmarshal(b, &event)
			Expect(err).ShouldNot(HaveOccurred())

			b, err = json.MarshalIndent(event, "", " ")
			Expect(err).ShouldNot(HaveOccurred())

			fmt.Println(string(b))

			Expect(true).Should(BeTrue())
		})
	})
})
