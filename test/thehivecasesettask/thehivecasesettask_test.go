package thehivecasesettask_test

import (
	"context"
	"fmt"
	"os"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/av-belyakov/thehivehook_go_package/cmd/commoninterfaces"
	"github.com/av-belyakov/thehivehook_go_package/cmd/natsapi"
	"github.com/av-belyakov/thehivehook_go_package/cmd/thehiveapi"
	"github.com/av-belyakov/thehivehook_go_package/internal/confighandler"
	helperfunc "github.com/av-belyakov/thehivehook_go_package/test/helpfunc"
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
		requestId    string = uuid.New().String()
		//caseId       string = "39100"
		//rootId       string = "~88678416456" //это мой тестовый кейс с id 39100

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
			logging := helperfunc.NewLoggingForTest()

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
			req.SetCommand("send_command")
			req.SetOrder("add_case_task")
			req.SetRequestId(requestId)
			req.SetData([]byte(`{
              "service": "MISP",
	          "command": "add_case_task",
	          "root_id": "~88678416456",
	          "case_id": "39100",
	          "field_name": "NewGroup",
	          "value": "analise new attack type - SQL-injection",
			  "username": "a.belyakov@cloud.gcm"
            }`))
			//"field_name": "Developers"
			chApiTheHive <- req

			msg := <-logging.GetChan()
			fmt.Println("Type:", msg.GetType(), " LOG:", msg.GetMessage())

			Expect(msg.GetType()).ShouldNot(Equal("error"))
		})
	})
})
