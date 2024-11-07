package thehivecasesettask_test

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

type TaskParameters struct {
	Type     string
	Value    string
	Username string
}

func (tp TaskParameters) GetType() string {
	return tp.Type
}

func (tp TaskParameters) GetValue() string {
	return tp.Value
}

func (tp TaskParameters) GetUsername() string {
	return tp.Username
}

var _ = Describe("Testthehivecasesettask", Ordered, func() {
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
			req.SetOrder("add case task")
			req.SetRootId(rootId)
			req.SetCaseId(caseId)
			req.SetData(TaskParameters{
				Type:     "Developers",
				Value:    "new filtration data",
				Username: "a.belyakov@cloud.gcm",
			})

			chApiTheHive <- req

			msg := <-logging.GetChan()
			fmt.Println("Type:", msg.GetType(), " LOG:", msg.GetMessage())

			Expect(msg.GetType()).ShouldNot(Equal("error"))
		})
	})
})
