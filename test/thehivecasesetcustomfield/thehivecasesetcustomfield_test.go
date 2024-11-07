package thehivecasesetcustomfield_test

import (
	"context"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/av-belyakov/thehivehook_go_package/cmd/commoninterfaces"
	"github.com/av-belyakov/thehivehook_go_package/cmd/natsapi"
	"github.com/av-belyakov/thehivehook_go_package/cmd/thehiveapi"
	"github.com/av-belyakov/thehivehook_go_package/internal/confighandler"
	"github.com/av-belyakov/thehivehook_go_package/internal/logginghandler"
)

type CustomFieldParameters struct {
	Type     string
	Value    string
	Username string
}

func (cfp CustomFieldParameters) GetType() string {
	return cfp.Type
}

func (cfp CustomFieldParameters) GetValue() string {
	return cfp.Value
}

func (cfp CustomFieldParameters) GetUsername() string {
	return cfp.Username
}

var _ = Describe("Testthehivecasesetcustomfield", Ordered, func() {
	var (
		chApiTheHive chan<- commoninterfaces.ChannelRequester
		caseId       string = "39100"
		rootId       string = "~88678416456" //это мой тестовый кейс с id 39100

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

	Context("Тест 1. Добавление CustomFiled к заданному кейсы TheHive", func() {
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

			chApiTheHive, err = apiTheHive.Start(context.Background())
			Expect(err).ShouldNot(HaveOccurred())

			req := natsapi.NewChannelRequest()
			req.SetCommand("send command")
			req.SetOrder("add case custom fields")
			req.SetRootId(rootId)
			req.SetCaseId(caseId)
			req.SetData(CustomFieldParameters{
				Type:  "attack-type.string",                   //Type:  "misp-event-id.string",
				Value: "выполнение произвольного SQL запроса", //Value: "73f",
			})

			chApiTheHive <- req

			msg := <-logging.GetChan()
			fmt.Println("Type:", msg.GetType(), " LOG:", msg.GetMessage())

			Expect(msg.GetType()).ShouldNot(Equal("error"))
		})
	})
})
