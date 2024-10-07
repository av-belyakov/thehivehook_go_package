package testthehiverequest_test

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/av-belyakov/thehivehook_go_package/cmd/commoninterfaces"
	"github.com/av-belyakov/thehivehook_go_package/cmd/thehiveapi"
	"github.com/av-belyakov/thehivehook_go_package/cmd/webhookserver"
	"github.com/av-belyakov/thehivehook_go_package/internal/confighandler"
	"github.com/av-belyakov/thehivehook_go_package/internal/logginghandler"
)

var _ = Describe("Testthehiverequest", Ordered, func() {
	var (
		rootDir string = "thehivehook_go_package"

		conf           *confighandler.ConfigApp
		chanTheHiveAPI chan<- commoninterfaces.ChannelRequester

		chanDone chan struct{}

		errConf, errTheHiveApi error
	)

	BeforeAll(func() {
		//
		// ВАЖНО!!!
		//
		//перед запуском теста установите переменную окружения GO_HIVEHOOK_THAPIKEY
		//с ключем-идентификатором, необходимым для авторизации в API TheHive,
		//командой export GO_HIVEHOOK_THAPIKEY=<api_key> или воспользоватся godotenv
		if err := godotenv.Load("../../.env"); err != nil {
			log.Fatalln(err)
		}

		chanDone = make(chan struct{})

		conf, errConf = confighandler.NewConfig(rootDir)
		confTheHive := conf.GetApplicationTheHive()

		logging := logginghandler.New()

		go func() {
			for {
				select {
				case <-chanDone:
					return

				case msg := <-logging.GetChan():
					fmt.Println("Log:", msg)
				}
			}
		}()

		chanTheHiveAPI, errTheHiveApi = thehiveapi.New(
			context.Background(),
			logging,
			thehiveapi.WithAPIKey(os.Getenv("GO_HIVEHOOK_THAPIKEY")),
			thehiveapi.WithHost(confTheHive.Host),
			thehiveapi.WithPort(confTheHive.Port))
	})

	Context("Тест 0. Чтение конфигурационного файла", func() {
		It("При чтении конфигурационного файла не должно быть ошибок", func() {
			Expect(errConf).ShouldNot(HaveOccurred())
		})
	})

	Context("Тест 1. Инициализация модуля взаимодействия с API TheHive", func() {
		It("При инициализации модуля не должно быть ошибок", func() {
			Expect(errTheHiveApi).ShouldNot(HaveOccurred())
		})
	})

	Context("Тест 2. Выполнение запросов к TheHive", func() {
		It("Запрос на получения Observable должен быть успешно выполнен", func() {
			var (
				statusCode int
				rootId     string = "~86676517008" //caseId:35144
				myUuid     string = uuid.New().String()
				myUuidRes  string
				wg         sync.WaitGroup

				chanResObservable chan commoninterfaces.ChannelResponser = make(chan commoninterfaces.ChannelResponser)
				chanResTTL        chan commoninterfaces.ChannelResponser = make(chan commoninterfaces.ChannelResponser)
			)

			wg.Add(2)

			go func() {
				for res := range chanResObservable {
					myUuidRes = res.GetRequestId()
					statusCode = res.GetStatusCode()

					fmt.Println("--------- Observable ----------")
					fmt.Println("Resived Response")
					fmt.Println("RequestId:", res.GetRequestId())

					msg := []interface{}{}
					err := json.Unmarshal(res.GetData(), &msg)
					Expect(err).ShouldNot(HaveOccurred())

					b, err := json.MarshalIndent(msg, "", " ")
					Expect(err).ShouldNot(HaveOccurred())

					fmt.Println("DATA:", string(b))
				}

				wg.Done()
			}()
			go func() {
				for res := range chanResTTL {
					myUuidRes = res.GetRequestId()
					statusCode = res.GetStatusCode()

					fmt.Println("--------- TTL ----------")
					fmt.Println("Resived Response")
					fmt.Println("RequestId:", res.GetRequestId())

					msg := []interface{}{}
					err := json.Unmarshal(res.GetData(), &msg)
					Expect(err).ShouldNot(HaveOccurred())

					b, err := json.MarshalIndent(msg, "", " ")
					Expect(err).ShouldNot(HaveOccurred())

					fmt.Println("DATA:", string(b))
				}

				wg.Done()
			}()

			reqObservable := webhookserver.NewChannelRequest()
			reqObservable.SetRequestId(myUuid)
			reqObservable.SetRootId(rootId)
			reqObservable.SetCommand("get_observables")
			reqObservable.SetChanOutput(chanResObservable)
			chanTheHiveAPI <- reqObservable

			reqTTP := webhookserver.NewChannelRequest()
			reqTTP.SetRequestId(myUuid)
			reqTTP.SetRootId(rootId)
			reqTTP.SetCommand("get_ttp")
			reqTTP.SetChanOutput(chanResTTL)
			chanTheHiveAPI <- reqTTP

			wg.Wait()

			chanDone <- struct{}{}

			Expect(statusCode).Should(Equal(200))
			Expect(myUuidRes).Should(Equal(myUuid))
		})
		/*It("Запрос на получения TTP должен быть успешно выполнен", func() {

		})*/
	})

	/*
		Context("", func(){
			It("", func(){

			})
		})
	*/
})
